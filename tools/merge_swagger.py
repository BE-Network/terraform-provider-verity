#!/usr/bin/env python3
"""
Merge Multiple Swagger/OpenAPI Specification Files

This script merges multiple Swagger JSON files into a unified specification,
combining all endpoints and taking the union of fields for overlapping endpoints.

Usage:
    python3 merge_swagger.py <base_file> <file1> [<file2> ...]

Arguments:
    base_file: Path to the base Swagger JSON file
    file1, file2, etc.: Additional Swagger JSON files to merge

Output:
    A merged Swagger file named "swagger_merged.json" in the current directory

Example:
    python3 merge_swagger.py swagger.json swagger_65_dc.json swagger_65_campus.json
"""
import json
import sys
import argparse
from collections.abc import Mapping
from copy import deepcopy


def deep_merge(base, overlay):
    """
    Deep merge two dictionaries recursively.
    For overlapping keys, if both values are:
    - Dictionaries: merge them recursively
    - Lists: combine them (with deduplication for primitive types)
    - Other: prefer the overlay value
    """
    result = deepcopy(base)
    
    for key, value in overlay.items():
        if key in result:
            # If both values are dictionaries, merge them recursively
            if isinstance(result[key], Mapping) and isinstance(value, Mapping):
                result[key] = deep_merge(result[key], value)
            # If both are lists, combine them (with deduplication for simple types)
            elif isinstance(result[key], list) and isinstance(value, list):
                # For simple types (strings, numbers, etc.), deduplicate
                if all(not isinstance(item, (dict, list)) for item in result[key] + value):
                    result[key] = list(set(result[key] + value))
                # For complex types (dicts, lists), deduplicate by converting to JSON strings
                else:
                    merged_list = result[key].copy()
                    for item in value:
                        # Skip if this item is a duplicate
                        if not any(json.dumps(item, sort_keys=True) == json.dumps(existing, sort_keys=True) 
                                  for existing in merged_list):
                            merged_list.append(deepcopy(item))
                    result[key] = merged_list
            # For other types, prefer the overlay value
            else:
                result[key] = deepcopy(value)
        else:
            # Key doesn't exist in base, add it
            result[key] = deepcopy(value)
    
    return result


def merge_params(base_params, overlay_params):
    """
    Merge parameters lists from multiple Swagger files, deduplicating by name and location.
    """
    if not base_params:
        return deepcopy(overlay_params) if overlay_params else []
    if not overlay_params:
        return deepcopy(base_params)
    
    result = deepcopy(base_params)
    param_keys = set()
    
    # Create keys for existing parameters
    for param in result:
        if 'name' in param and 'in' in param:
            param_keys.add(f"{param['name']}:{param['in']}")
    
    # Add non-duplicate parameters from overlay
    for param in overlay_params:
        if 'name' in param and 'in' in param:
            param_key = f"{param['name']}:{param['in']}"
            if param_key not in param_keys:
                result.append(deepcopy(param))
                param_keys.add(param_key)
    
    return result


def merge_security(base_sec, overlay_sec):
    """
    Merge security requirements, deduplicating by name.
    """
    if not base_sec:
        return deepcopy(overlay_sec) if overlay_sec else []
    if not overlay_sec:
        return deepcopy(base_sec)
    
    result = []
    seen = set()
    
    # First add from base, tracking what we've seen
    for sec in base_sec:
        key = json.dumps(sec, sort_keys=True)
        if key not in seen:
            result.append(deepcopy(sec))
            seen.add(key)
    
    # Then add from overlay if not a duplicate
    for sec in overlay_sec:
        key = json.dumps(sec, sort_keys=True)
        if key not in seen:
            result.append(deepcopy(sec))
            seen.add(key)
    
    return result


def merge_servers(base_servers, overlay_servers):
    """
    Merge servers lists, deduplicating by URL.
    """
    if not base_servers:
        return deepcopy(overlay_servers) if overlay_servers else []
    if not overlay_servers:
        return deepcopy(base_servers)
    
    result = deepcopy(base_servers)
    server_urls = {server['url'] for server in result}
    
    # Add non-duplicate servers from overlay
    for server in overlay_servers:
        if server['url'] not in server_urls:
            result.append(deepcopy(server))
            server_urls.add(server['url'])
    
    return result


def merge_paths(base_paths, overlay_paths):
    """
    Merge paths from multiple Swagger files.
    
    For paths that exist in both:
    - Merge their HTTP methods (GET, POST, etc.)
    - For methods that exist in both, merge their properties
    """
    result = deepcopy(base_paths)
    
    for path, methods in overlay_paths.items():
        if path in result:
            # Path exists in both, merge methods
            for method, details in methods.items():
                if method in result[path]:
                    # Method exists in both, merge details
                    merged_details = deepcopy(result[path][method])
                    
                    # Special handling for parameters
                    if 'parameters' in details and 'parameters' in merged_details:
                        merged_details['parameters'] = merge_params(merged_details['parameters'], details['parameters'])
                    elif 'parameters' in details:
                        merged_details['parameters'] = deepcopy(details['parameters'])
                    
                    # Special handling for security
                    if 'security' in details and 'security' in merged_details:
                        merged_details['security'] = merge_security(merged_details['security'], details['security'])
                    elif 'security' in details:
                        merged_details['security'] = deepcopy(details['security'])
                    
                    # Merge the rest normally
                    for key, value in details.items():
                        if key not in ['parameters', 'security']:
                            if key in merged_details:
                                if isinstance(merged_details[key], Mapping) and isinstance(value, Mapping):
                                    merged_details[key] = deep_merge(merged_details[key], value)
                                elif isinstance(merged_details[key], list) and isinstance(value, list):
                                    # Deduplicate lists
                                    merged_list = merged_details[key].copy()
                                    for item in value:
                                        if item not in merged_list:
                                            merged_list.append(item)
                                    merged_details[key] = merged_list
                                else:
                                    # Prefer overlay for non-container types
                                    merged_details[key] = deepcopy(value)
                            else:
                                merged_details[key] = deepcopy(value)
                    
                    result[path][method] = merged_details
                else:
                    # Method only in overlay, add it
                    result[path][method] = deepcopy(details)
        else:
            # Path only in overlay, add it
            result[path] = deepcopy(methods)
    
    return result


def merge_components(base_components, overlay_components):
    """
    Merge components section (schemas, securitySchemes, etc.)
    """
    if not base_components and not overlay_components:
        return {}
    if not base_components:
        return deepcopy(overlay_components)
    if not overlay_components:
        return deepcopy(base_components)
    
    result = deepcopy(base_components)
    
    for section, items in overlay_components.items():
        if section in result:
            # Section exists in both (schemas, securitySchemes, etc.)
            for name, definition in items.items():
                if name in result[section]:
                    # Definition exists in both, merge them
                    result[section][name] = deep_merge(result[section][name], definition)
                else:
                    # Definition only in overlay, add it
                    result[section][name] = deepcopy(definition)
        else:
            # Section only in overlay, add it
            result[section] = deepcopy(items)
    
    return result


def merge_swagger_files(base_file, overlay_files):
    """
    Merge multiple Swagger JSON files into one unified specification.
    
    Args:
        base_file: Path to the base Swagger file
        overlay_files: List of paths to additional Swagger files to merge
    
    Returns:
        Merged Swagger specification as a dictionary
    """
    try:
        with open(base_file, 'r') as f:
            base_swagger = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Error loading base file '{base_file}': {e}")
        sys.exit(1)
    
    result = deepcopy(base_swagger)
    
    for overlay_file in overlay_files:
        try:
            with open(overlay_file, 'r') as f:
                overlay_swagger = json.load(f)
        except (FileNotFoundError, json.JSONDecodeError) as e:
            print(f"Error loading overlay file '{overlay_file}': {e}")
            continue
        
        # Merge paths section
        if 'paths' in overlay_swagger:
            if 'paths' not in result:
                result['paths'] = {}
            result['paths'] = merge_paths(result['paths'], overlay_swagger['paths'])
        
        # Merge components section
        if 'components' in overlay_swagger:
            if 'components' not in result:
                result['components'] = {}
            result['components'] = merge_components(result['components'], overlay_swagger['components'])
        
        # Merge servers section with deduplication
        if 'servers' in overlay_swagger:
            if 'servers' not in result:
                result['servers'] = []
            result['servers'] = merge_servers(result['servers'], overlay_swagger['servers'])
        
        # Merge other top-level properties
        for key, value in overlay_swagger.items():
            if key not in ['paths', 'components', 'servers']:
                if key in result:
                    # If both values are lists, combine them
                    if isinstance(result[key], list) and isinstance(value, list):
                        result[key] = list(set(result[key] + value)) if isinstance(value[0], (str, int, float, bool)) else result[key] + value
                    # If both values are dicts, merge them
                    elif isinstance(result[key], dict) and isinstance(value, dict):
                        result[key] = deep_merge(result[key], value)
                    # For other types, keep the base value
                else:
                    # Key only in overlay, add it
                    result[key] = deepcopy(value)
    
    return result


def main():
    parser = argparse.ArgumentParser(description='Merge multiple Swagger/OpenAPI specification files')
    parser.add_argument('base_file', help='Base Swagger JSON file')
    parser.add_argument('overlay_files', nargs='+', help='Additional Swagger files to merge')
    parser.add_argument('--output', '-o', default='swagger_merged.json', help='Output file name')
    args = parser.parse_args()
    
    merged_swagger = merge_swagger_files(args.base_file, args.overlay_files)
    
    with open(args.output, 'w') as f:
        json.dump(merged_swagger, f, indent=2)
    
    print(f"Merged Swagger specification saved to {args.output}")
    print(f"Base file: {args.base_file}")
    print(f"Overlay files: {', '.join(args.overlay_files)}")


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Error: You must provide at least two Swagger JSON files to merge.")
        print("Usage: python3 merge_swagger.py <base_file> <file1> [<file2> ...]")
        sys.exit(1)
    
    main()
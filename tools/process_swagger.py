#!/usr/bin/env python3
"""
Unified Swagger Processing Pipeline

This script combines four operations into one seamless workflow:
1. Merge multiple Swagger/OpenAPI specification files
2. Remove unwanted endpoints from the merged specification
3. Make all number and integer fields nullable
4. Transform the schema for proper Go SDK generation

Usage:
    python3 process_swagger.py <file1.json> <file2.json> [options]

Arguments:
    file1.json: First Swagger JSON file (base file)
    file2.json: Second Swagger JSON file (overlay file)

Options:
    --output, -o: Output filename (default: merged_transformed.json)
    --keep-intermediate: Keep intermediate files (merged, cleaned, nullable versions)

Output:
    A fully processed Swagger file ready for Go SDK generation

Example:
    python3 process_swagger.py vncgenie.json vnckrait.json
    python3 process_swagger.py vncgenie.json vnckrait.json --output final_swagger.json
"""
import json
import sys
import argparse
import os
from collections.abc import Mapping
from copy import deepcopy
from pathlib import Path


# ============================================================================
# SECTION 1: MERGE SWAGGER FILES
# ============================================================================

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


def merge_swagger_files(base_file, overlay_file):
    """
    Merge two Swagger JSON files into one unified specification.
    
    Args:
        base_file: Path to the base Swagger file
        overlay_file: Path to the overlay Swagger file to merge
    
    Returns:
        Merged Swagger specification as a dictionary
    """
    try:
        with open(base_file, 'r') as f:
            base_swagger = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Error loading base file '{base_file}': {e}")
        sys.exit(1)
    
    try:
        with open(overlay_file, 'r') as f:
            overlay_swagger = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Error loading overlay file '{overlay_file}': {e}")
        sys.exit(1)
    
    result = deepcopy(base_swagger)
    
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


# ============================================================================
# SECTION 2: REMOVE UNWANTED ENDPOINTS
# ============================================================================

def remove_endpoints(swagger_data, endpoints_to_remove):
    """
    Remove specified endpoints from the swagger data.
    
    Args:
        swagger_data (dict): The swagger data structure
        endpoints_to_remove (list): List of endpoint paths to remove
    
    Returns:
        tuple: (modified swagger_data, list of removed endpoints)
    """
    # Check if 'paths' exists in the swagger data
    if 'paths' not in swagger_data:
        print("Warning: No 'paths' section found in swagger data")
        return swagger_data, []
    
    # Remove the specified endpoints
    removed_endpoints = []
    for endpoint in endpoints_to_remove:
        if endpoint in swagger_data['paths']:
            del swagger_data['paths'][endpoint]
            removed_endpoints.append(endpoint)
            print(f"  ✓ Removed endpoint: {endpoint}")
        else:
            print(f"  - Endpoint not found (skipping): {endpoint}")
    
    return swagger_data, removed_endpoints


# ============================================================================
# SECTION 3: TRANSFORM SCHEMA FOR GO SDK GENERATION
# ============================================================================

def transform_schema(swagger_data):
    """
    Transform the swagger schema for proper Go SDK generation.
    Adds additionalProperties fields where needed.
    
    Args:
        swagger_data (dict): The swagger data structure
    
    Returns:
        dict: Transformed swagger data
    """
    transformation_count = 0
    
    for path, methods in swagger_data['paths'].items():
        if path.startswith('/config'):
            continue
            
        for method_name in ['put', 'patch']:
            if method_name in methods:
                method = methods[method_name]
                if 'requestBody' in method and 'content' in method['requestBody'] and 'application/json' in method['requestBody']['content']:
                    content = method['requestBody']['content']['application/json']
                    
                    if 'schema' in content and 'properties' in content['schema']:
                        for prop_name, prop_value in content['schema']['properties'].items():
                            if 'properties' in prop_value:
                                # Get the first property from the inner properties object
                                inner_properties = prop_value['properties']
                                
                                if inner_properties and len(inner_properties) == 1:
                                    # Get the first property key (like gateway_profile_name)
                                    inner_prop_name = next(iter(inner_properties))
                                    inner_prop_value = inner_properties[inner_prop_name]
                                    
                                    # Create the additionalProperties structure
                                    prop_value['additionalProperties'] = inner_prop_value
                                    
                                    # Empty the original properties
                                    prop_value['properties'] = {}
                                    
                                    transformation_count += 1
                                    print(f"  ✓ Transformed: {path} [{method_name}] - {prop_name}")
    
    return swagger_data, transformation_count


# ============================================================================
# SECTION 4: MAKE NUMBER AND INTEGER FIELDS NULLABLE
# ============================================================================

def make_numeric_fields_nullable(obj, path_context="", field_name=""):
    """
    Recursively traverse the schema and add 'nullable: true' to all number 
    and integer type fields that don't already have it.
    Skips fields named 'index' as they should not be nullable.
    
    Args:
        obj: The object to traverse (dict, list, or other)
        path_context: String representing current path in schema (for tracking)
        field_name: The name of the current field being processed
    
    Returns:
        tuple: (modified object, count of fields made nullable)
    """
    nullable_count = 0
    
    if isinstance(obj, dict):
        # Check if this is a schema property with type number or integer
        # Skip if the field name is 'index'
        if 'type' in obj and obj['type'] in ['number', 'integer']:
            if field_name != 'index' and ('nullable' not in obj or obj['nullable'] != True):
                obj['nullable'] = True
                nullable_count += 1
        
        # Recursively process all nested dictionaries
        for key, value in obj.items():
            new_path = f"{path_context}.{key}" if path_context else key
            if isinstance(value, (dict, list)):
                _, count = make_numeric_fields_nullable(value, new_path, key)
                nullable_count += count
    
    elif isinstance(obj, list):
        # Recursively process all items in lists
        for idx, item in enumerate(obj):
            if isinstance(item, (dict, list)):
                _, count = make_numeric_fields_nullable(item, f"{path_context}[{idx}]", field_name)
                nullable_count += count
    
    return obj, nullable_count


def process_nullable_fields(swagger_data):
    """
    Process the entire Swagger data structure to make all number and integer
    fields nullable.
    
    Args:
        swagger_data (dict): The swagger data structure
    
    Returns:
        tuple: (modified swagger_data, total count of fields made nullable)
    """
    total_count = 0
    
    # Process paths
    if 'paths' in swagger_data:
        for path, methods in swagger_data['paths'].items():
            for method_name, method_data in methods.items():
                # Process request body schemas
                if 'requestBody' in method_data:
                    context = f"{path} [{method_name}] requestBody"
                    _, count = make_numeric_fields_nullable(method_data['requestBody'], context)
                    total_count += count
                
                # Process response schemas
                if 'responses' in method_data:
                    for status_code, response_data in method_data['responses'].items():
                        context = f"{path} [{method_name}] response {status_code}"
                        _, count = make_numeric_fields_nullable(response_data, context)
                        total_count += count
                
                # Process parameters
                if 'parameters' in method_data:
                    context = f"{path} [{method_name}] parameters"
                    _, count = make_numeric_fields_nullable(method_data['parameters'], context)
                    total_count += count
    
    # Process components/schemas
    if 'components' in swagger_data and 'schemas' in swagger_data['components']:
        for schema_name, schema_data in swagger_data['components']['schemas'].items():
            context = f"components.schemas.{schema_name}"
            _, count = make_numeric_fields_nullable(schema_data, context)
            total_count += count
    
    return swagger_data, total_count


# ============================================================================
# MAIN PIPELINE
# ============================================================================

def main():
    parser = argparse.ArgumentParser(
        description='Unified Swagger processing pipeline: merge, clean, and transform',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 process_swagger.py vncgenie.json vnckrait.json
  python3 process_swagger.py vncgenie.json vnckrait.json --output final.json
  python3 process_swagger.py vncgenie.json vnckrait.json --keep-intermediate
        """
    )
    parser.add_argument('file1', help='First Swagger JSON file (base)')
    parser.add_argument('file2', help='Second Swagger JSON file (overlay)')
    parser.add_argument('--output', '-o', default='merged_transformed.json', 
                        help='Output filename (default: merged_transformed.json)')
    parser.add_argument('--keep-intermediate', action='store_true',
                        help='Keep intermediate files (merged and cleaned versions)')
    
    args = parser.parse_args()
    
    # Validate input files
    if not Path(args.file1).exists():
        print(f"Error: File '{args.file1}' not found")
        sys.exit(1)
    
    if not Path(args.file2).exists():
        print(f"Error: File '{args.file2}' not found")
        sys.exit(1)
    
    print("=" * 70)
    print("SWAGGER PROCESSING PIPELINE")
    print("=" * 70)
    print(f"Base file:    {args.file1}")
    print(f"Overlay file: {args.file2}")
    print(f"Output file:  {args.output}")
    print("=" * 70)
    
    # Step 1: Merge Swagger files
    print("\n[STEP 1/4] Merging Swagger files...")
    merged_swagger = merge_swagger_files(args.file1, args.file2)
    print(f"  ✓ Successfully merged {args.file1} and {args.file2}")
    
    # Save intermediate merged file if requested
    if args.keep_intermediate:
        intermediate_merged = args.output.replace('.json', '_merged_only.json')
        with open(intermediate_merged, 'w') as f:
            json.dump(merged_swagger, f, indent=2)
        print(f"  ✓ Intermediate file saved: {intermediate_merged}")
    
    # Step 2: Remove unwanted endpoints
    print("\n[STEP 2/4] Removing unwanted endpoints...")
    endpoints_to_remove = [
        "/alarms/mask",
        "/config",
        "/changesets", 
        "/readmode",
        "/request",
        "/snmp",
        "/syslog",
        "/backups",
        "/timetraveler"
    ]
    
    cleaned_swagger, removed = remove_endpoints(merged_swagger, endpoints_to_remove)
    print(f"  ✓ Removed {len(removed)} endpoints")
    
    # Save intermediate cleaned file if requested
    if args.keep_intermediate:
        intermediate_cleaned = args.output.replace('.json', '_cleaned_only.json')
        with open(intermediate_cleaned, 'w') as f:
            json.dump(cleaned_swagger, f, indent=2)
        print(f"  ✓ Intermediate file saved: {intermediate_cleaned}")
    
    # Step 3: Make numeric fields nullable
    print("\n[STEP 3/4] Making number and integer fields nullable...")
    nullable_swagger, nullable_count = process_nullable_fields(cleaned_swagger)
    print(f"  ✓ Made {nullable_count} numeric fields nullable")
    
    # Save intermediate nullable file if requested
    if args.keep_intermediate:
        intermediate_nullable = args.output.replace('.json', '_nullable_only.json')
        with open(intermediate_nullable, 'w') as f:
            json.dump(nullable_swagger, f, indent=2)
        print(f"  ✓ Intermediate file saved: {intermediate_nullable}")
    
    # Step 4: Transform schema
    print("\n[STEP 4/4] Transforming schema for Go SDK generation...")
    transformed_swagger, transform_count = transform_schema(nullable_swagger)
    print(f"  ✓ Applied {transform_count} schema transformations")
    
    # Save final output
    print(f"\n[FINAL] Writing output to {args.output}...")
    with open(args.output, 'w') as f:
        json.dump(transformed_swagger, f, indent=2)
    
    print("\n" + "=" * 70)
    print("✓ PIPELINE COMPLETED SUCCESSFULLY!")
    print("=" * 70)
    print(f"Final output: {args.output}")
    print(f"Total paths: {len(transformed_swagger.get('paths', {}))}")
    print(f"Endpoints removed: {len(removed)}")
    print(f"Numeric fields made nullable: {nullable_count}")
    print(f"Schema transformations: {transform_count}")
    print("=" * 70)


if __name__ == "__main__":
    if len(sys.argv) < 2 or (len(sys.argv) < 3 and '--help' not in sys.argv and '-h' not in sys.argv):
        print("Error: You must provide two Swagger JSON files to process.")
        print("Usage: python3 process_swagger.py <file1.json> <file2.json> [options]")
        print("Try --help for more information.")
        sys.exit(1)
    
    main()

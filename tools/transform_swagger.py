#!/usr/bin/env python3
"""
Transform Swagger JSON Schema for Go SDK Generation
This script modifies Swagger JSON schema to adapt it for proper Go SDK generation.
Specifically, it transforms property structures by adding additionalProperties fields where needed.
This enables the Go SDK generator to correctly handle maps and dynamic property structures.
Usage:
    python3 transform_swagger.py <json_file>
Arguments:
    input_file: Path to the input Swagger JSON file
Output:
    A transformed Swagger file named with "_transformed" suffix in the same directory
Example:
    python3 transform_swagger.py swagger.json
"""
import json
import os
import sys
import argparse



def transform_schema(input_file, output_file):
    try:
        with open(input_file, 'r') as f:
            swagger = json.load(f)
    except FileNotFoundError:
        print(f"Error: Input file '{input_file}' not found.")
        sys.exit(1)
    except json.JSONDecodeError:
        print(f"Error: '{input_file}' is not a valid JSON file.")
        sys.exit(1)
    
    for path, methods in swagger['paths'].items():
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
                                    
                                    print(f"Transformed {path} - {method_name} - {prop_name}")
    
    with open(output_file, 'w') as f:
        json.dump(swagger, f, indent=2)
    
    print(f"Transformation complete. Output saved to {output_file}")

def main():
    parser = argparse.ArgumentParser(description='Transform Swagger JSON to add additionalProperties')
    parser.add_argument('input_file', help='Input Swagger JSON file (required)')
    args = parser.parse_args()
    
    filename, ext = os.path.splitext(args.input_file)
    output_file = f"{filename}_transformed{ext}"
    
    transform_schema(args.input_file, output_file)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Error: You must provide a valid JSON file name as an argument.")
        print("Usage: python3 transform_swagger.py <json_file>")
        sys.exit(1)
    
    main()
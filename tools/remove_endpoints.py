#!/usr/bin/env python3
"""
Script to remove unwanted endpoint definitions from swagger.json
"""

import json
import sys
from pathlib import Path

def remove_endpoints(swagger_file_path, endpoints_to_remove):
    """
    Remove specified endpoints from the swagger JSON file
    
    Args:
        swagger_file_path (str): Path to the swagger JSON file
        endpoints_to_remove (list): List of endpoint paths to remove
    """
    
    # Read the swagger file
    try:
        with open(swagger_file_path, 'r') as f:
            swagger_data = json.load(f)
    except FileNotFoundError:
        print(f"Error: File {swagger_file_path} not found")
        return False
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON in {swagger_file_path}: {e}")
        return False
    
    # Check if 'paths' exists in the swagger data
    if 'paths' not in swagger_data:
        print("Error: No 'paths' section found in swagger file")
        return False
    
    # Remove the specified endpoints
    removed_endpoints = []
    for endpoint in endpoints_to_remove:
        if endpoint in swagger_data['paths']:
            del swagger_data['paths'][endpoint]
            removed_endpoints.append(endpoint)
            print(f"Removed endpoint: {endpoint}")
        else:
            print(f"Warning: Endpoint {endpoint} not found in swagger file")
    
    # Create backup of original file
    backup_path = swagger_file_path + '.backup'
    try:
        with open(backup_path, 'w') as f:
            with open(swagger_file_path, 'r') as original:
                f.write(original.read())
        print(f"Backup created: {backup_path}")
    except Exception as e:
        print(f"Warning: Could not create backup: {e}")
    
    # Write the modified swagger data back to file
    try:
        with open(swagger_file_path, 'w') as f:
            json.dump(swagger_data, f, indent=2)
        print(f"Successfully updated {swagger_file_path}")
        print(f"Removed {len(removed_endpoints)} endpoints: {removed_endpoints}")
        return True
    except Exception as e:
        print(f"Error writing to file: {e}")
        return False

def main():
    # Check if filename argument is provided
    if len(sys.argv) != 2:
        print("Usage: python3 remove_endpoints.py <swagger_file.json>")
        print("Example: python3 remove_endpoints.py swagger.json")
        sys.exit(1)
    
    # Get the swagger file path from command line argument
    swagger_file = sys.argv[1]
    
    # Define the endpoints to remove
    endpoints_to_remove = [
        "/config",
        "/changesets", 
        "/readmode",
        "/request",
        "/snmp",
        "/syslog",
        "/backups",
        "/timetraveler"
    ]
    
    # Check if file exists
    if not Path(swagger_file).exists():
        print(f"Error: {swagger_file} not found")
        sys.exit(1)
    
    print(f"Removing unwanted endpoints from {swagger_file}...")
    print(f"Endpoints to remove: {endpoints_to_remove}")
    print()
    
    # Remove the endpoints
    success = remove_endpoints(swagger_file, endpoints_to_remove)
    
    if success:
        print("\nOperation completed successfully!")
        print("Remember to regenerate your OpenAPI SDK after this change.")
    else:
        print("\nOperation failed!")
        sys.exit(1)

if __name__ == "__main__":
    main()

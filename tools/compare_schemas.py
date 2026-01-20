#!/usr/bin/env python3
"""
This script compares OpenAPI schemas for datacenter and campus API modes,
identifying common fields, datacenter-specific fields, and campus-specific fields.

Usage:
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --endpoint /services
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --output report.json
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --generate-go
"""

import argparse
import json
import sys
from pathlib import Path
from typing import Any, Dict, List, Optional, Set, Tuple
from dataclasses import dataclass, field
from enum import Enum


class FieldMode(Enum):
    COMMON = "common"
    DATACENTER_ONLY = "datacenter_only"
    CAMPUS_ONLY = "campus_only"


# Endpoints to exclude from comparison
EXCLUDED_ENDPOINTS = {
    "/auth",
    "/alarms",
    "/alarms/mask",
    "/backups",
    "/changesets",
    "/config",
    "/imageupdatesets",
    "/readmode",
    "/request",
    "/snmp",
    "/syslog",
    "/switchpoints/currentconfig",
    "/switchpoints/markoutofservice",
    "/switchpoints/upgrade",
    "/timetraveler",
    "/version",
}


@dataclass
class FieldInfo:
    """Information about a field in the schema."""
    name: str
    field_type: str
    default: Any
    nullable: bool
    description: str
    mode: FieldMode
    nested_fields: Optional[Dict[str, 'FieldInfo']] = None
    path: str = ""
    
    def to_dict(self) -> dict:
        """Convert to dictionary for JSON serialization."""
        result = {
            "name": self.name,
            "type": self.field_type,
            "default": self.default,
            "nullable": self.nullable,
            "description": self.description,
            "mode": self.mode.value,
            "path": self.path
        }
        if self.nested_fields:
            result["nested_fields"] = {
                k: v.to_dict() for k, v in self.nested_fields.items()
            }
        return result


@dataclass
class EndpointComparison:
    """Comparison result for a single endpoint."""
    endpoint: str
    common_fields: Dict[str, FieldInfo] = field(default_factory=dict)
    datacenter_only_fields: Dict[str, FieldInfo] = field(default_factory=dict)
    campus_only_fields: Dict[str, FieldInfo] = field(default_factory=dict)
    
    def to_dict(self) -> dict:
        """Convert to dictionary for JSON serialization."""
        return {
            "endpoint": self.endpoint,
            "common_fields": {k: v.to_dict() for k, v in self.common_fields.items()},
            "datacenter_only_fields": {k: v.to_dict() for k, v in self.datacenter_only_fields.items()},
            "campus_only_fields": {k: v.to_dict() for k, v in self.campus_only_fields.items()}
        }


def load_schema(filepath: str) -> dict:
    """Load a JSON schema file."""
    path = Path(filepath)
    if not path.exists():
        raise FileNotFoundError(f"Schema file not found: {filepath}")
    
    with open(path, 'r', encoding='utf-8') as f:
        return json.load(f)


def get_endpoints(schema: dict) -> Set[str]:
    """Extract all endpoint paths from the schema, excluding non-resource endpoints."""
    paths = schema.get("paths", {})
    return set(p for p in paths.keys() if p not in EXCLUDED_ENDPOINTS)


def extract_field_info(
    field_name: str,
    field_schema: dict,
    mode: FieldMode,
    parent_path: str = ""
) -> FieldInfo:
    """Extract field information from a schema field definition."""
    field_type = field_schema.get("type", "unknown")
    default = field_schema.get("default", None)
    description = field_schema.get("description", "")
    
    # Integer and number (float) fields are always nullable by definition
    if field_type in ("integer", "number"):
        nullable = True
    else:
        nullable = field_schema.get("nullable", False)
    
    current_path = f"{parent_path}.{field_name}" if parent_path else field_name
    
    nested_fields = None
    
    # Handle array types with nested objects
    if field_type == "array":
        items = field_schema.get("items", {})
        if items.get("type") == "object":
            nested_props = items.get("properties", {})
            if nested_props:
                nested_fields = {}
                for nested_name, nested_schema in nested_props.items():
                    nested_fields[nested_name] = extract_field_info(
                        nested_name, 
                        nested_schema, 
                        mode,
                        f"{current_path}[]"
                    )
    
    # Handle object types with nested properties
    elif field_type == "object":
        nested_props = field_schema.get("properties", {})
        if nested_props:
            nested_fields = {}
            for nested_name, nested_schema in nested_props.items():
                nested_fields[nested_name] = extract_field_info(
                    nested_name, 
                    nested_schema, 
                    mode,
                    current_path
                )
    
    return FieldInfo(
        name=field_name,
        field_type=field_type,
        default=default,
        nullable=nullable,
        description=description,
        mode=mode,
        nested_fields=nested_fields,
        path=current_path
    )


def get_request_body_properties(endpoint_data: dict, method: str = "put") -> dict:
    """Extract properties from the request body of an endpoint."""
    method_data = endpoint_data.get(method, {})
    if not method_data:
        # Try patch if put is not available
        method_data = endpoint_data.get("patch", {})
    if not method_data:
        return {}
    
    request_body = method_data.get("requestBody", {})
    content = request_body.get("content", {})
    json_content = content.get("application/json", {})
    schema = json_content.get("schema", {})
    properties = schema.get("properties", {})
    
    # Get the actual field properties
    all_fields = {}
    for resource_key, resource_schema in properties.items():
        if resource_schema.get("type") == "object":
            # First check additionalProperties - this is where most fields are defined
            # The swagger structure is: properties.service.additionalProperties.properties
            additional_props = resource_schema.get("additionalProperties", {})
            if additional_props.get("type") == "object":
                add_props_fields = additional_props.get("properties", {})
                if add_props_fields:
                    all_fields.update(add_props_fields)
            
            # Also check direct properties
            resource_props = resource_schema.get("properties", {})
            for obj_key, obj_schema in resource_props.items():
                if obj_schema.get("type") == "object":
                    obj_props = obj_schema.get("properties", {})
                    all_fields.update(obj_props)
                else:
                    all_fields[obj_key] = obj_schema
        else:
            all_fields[resource_key] = resource_schema
    
    return all_fields


def compare_fields(
    dc_fields: dict,
    campus_fields: dict,
    parent_path: str = ""
) -> Tuple[Dict[str, FieldInfo], Dict[str, FieldInfo], Dict[str, FieldInfo]]:
    """
    Compare fields between datacenter and campus schemas.
    
    Returns:
        Tuple of (common_fields, datacenter_only_fields, campus_only_fields)
    """
    dc_keys = set(dc_fields.keys())
    campus_keys = set(campus_fields.keys())
    
    common_keys = dc_keys & campus_keys
    dc_only_keys = dc_keys - campus_keys
    campus_only_keys = campus_keys - dc_keys
    
    common_fields = {}
    dc_only_fields = {}
    campus_only_fields = {}
    
    # Process common fields
    for key in common_keys:
        dc_schema = dc_fields[key]
        campus_schema = campus_fields[key]
        
        # Use datacenter schema as reference for common fields
        field_info = extract_field_info(key, dc_schema, FieldMode.COMMON, parent_path)
        
        # Check if nested fields need comparison
        if field_info.nested_fields and dc_schema.get("type") in ("array", "object"):
            dc_nested = {}
            campus_nested = {}
            
            if dc_schema.get("type") == "array":
                dc_nested = dc_schema.get("items", {}).get("properties", {})
                campus_nested = campus_schema.get("items", {}).get("properties", {})
            elif dc_schema.get("type") == "object":
                dc_nested = dc_schema.get("properties", {})
                campus_nested = campus_schema.get("properties", {})
            
            if dc_nested or campus_nested:
                nested_common, nested_dc_only, nested_campus_only = compare_fields(
                    dc_nested, 
                    campus_nested,
                    f"{field_info.path}[]" if dc_schema.get("type") == "array" else field_info.path
                )
                
                # Update nested fields with comparison results
                all_nested = {}
                all_nested.update(nested_common)
                all_nested.update(nested_dc_only)
                all_nested.update(nested_campus_only)
                field_info.nested_fields = all_nested
        
        common_fields[key] = field_info
    
    # Process datacenter-only fields
    for key in dc_only_keys:
        dc_only_fields[key] = extract_field_info(
            key, dc_fields[key], FieldMode.DATACENTER_ONLY, parent_path
        )
    
    # Process campus-only fields
    for key in campus_only_keys:
        campus_only_fields[key] = extract_field_info(
            key, campus_fields[key], FieldMode.CAMPUS_ONLY, parent_path
        )
    
    return common_fields, dc_only_fields, campus_only_fields


def compare_endpoint(
    endpoint: str,
    dc_schema: dict,
    campus_schema: dict,
    endpoint_mode: str = "common"  # "common", "datacenter_only", or "campus_only"
) -> EndpointComparison:
    """Compare a single endpoint between datacenter and campus schemas.
    
    Args:
        endpoint: The endpoint path (e.g., "/services")
        dc_schema: The datacenter schema
        campus_schema: The campus schema
        endpoint_mode: Whether this endpoint exists in both, datacenter only, or campus only
    """
    dc_endpoint = dc_schema.get("paths", {}).get(endpoint, {})
    campus_endpoint = campus_schema.get("paths", {}).get(endpoint, {})
    
    dc_fields = get_request_body_properties(dc_endpoint)
    campus_fields = get_request_body_properties(campus_endpoint)
    
    if endpoint_mode == "datacenter_only":
        # All fields are datacenter-only
        dc_only_fields = {}
        for key, schema in dc_fields.items():
            dc_only_fields[key] = extract_field_info(key, schema, FieldMode.DATACENTER_ONLY)
        return EndpointComparison(
            endpoint=endpoint,
            common_fields={},
            datacenter_only_fields=dc_only_fields,
            campus_only_fields={}
        )
    elif endpoint_mode == "campus_only":
        # All fields are campus-only
        campus_only_fields = {}
        for key, schema in campus_fields.items():
            campus_only_fields[key] = extract_field_info(key, schema, FieldMode.CAMPUS_ONLY)
        return EndpointComparison(
            endpoint=endpoint,
            common_fields={},
            datacenter_only_fields={},
            campus_only_fields=campus_only_fields
        )
    else:
        # Common endpoint - compare fields
        common, dc_only, campus_only = compare_fields(dc_fields, campus_fields)
        return EndpointComparison(
            endpoint=endpoint,
            common_fields=common,
            datacenter_only_fields=dc_only,
            campus_only_fields=campus_only
        )


def format_default_value(value: Any) -> str:
    """Format a default value for display."""
    if value is None:
        return "null"
    elif isinstance(value, bool):
        return str(value).lower()
    elif isinstance(value, str):
        return f'"{value}"' if value else '""'
    else:
        return str(value)


def format_field_type(field_info: FieldInfo) -> str:
    """Format field type with nullable indicator."""
    type_str = field_info.field_type
    if field_info.nullable:
        type_str += " (nullable)"
    return type_str


def print_field_tree(
    fields: Dict[str, FieldInfo],
    indent: int = 0,
    mode_filter: Optional[FieldMode] = None
) -> None:
    """Print fields in a tree structure."""
    indent_str = "  " * indent
    
    for name, field_info in sorted(fields.items()):
        if mode_filter and field_info.mode != mode_filter:
            continue
        
        mode_indicator = ""
        if field_info.mode == FieldMode.DATACENTER_ONLY:
            mode_indicator = " [DC]"
        elif field_info.mode == FieldMode.CAMPUS_ONLY:
            mode_indicator = " [CAMPUS]"
        
        default_str = format_default_value(field_info.default)
        type_str = format_field_type(field_info)
        
        print(f"{indent_str}├─ {name}: {type_str}{mode_indicator}")
        print(f"{indent_str}│    Default: {default_str}")
        
        if field_info.description:
            # Truncate long descriptions
            desc = field_info.description
            if len(desc) > 80:
                desc = desc[:77] + "..."
            print(f"{indent_str}│    Description: {desc}")
        
        if field_info.nested_fields:
            print(f"{indent_str}│    Nested fields:")
            print_field_tree(field_info.nested_fields, indent + 2)


def print_comparison_report(comparison: EndpointComparison) -> None:
    """Print a formatted comparison report for an endpoint."""
    print("\n" + "=" * 80)
    
    # Determine endpoint type
    has_common = len(comparison.common_fields) > 0
    has_dc_only = len(comparison.datacenter_only_fields) > 0
    has_campus_only = len(comparison.campus_only_fields) > 0
    
    if has_common or (has_dc_only and has_campus_only):
        endpoint_type = "COMMON ENDPOINT"
    elif has_dc_only and not has_campus_only:
        endpoint_type = "DATACENTER-ONLY ENDPOINT"
    elif has_campus_only and not has_dc_only:
        endpoint_type = "CAMPUS-ONLY ENDPOINT"
    else:
        endpoint_type = "ENDPOINT"
    
    print(f"{endpoint_type}: {comparison.endpoint}")
    print("=" * 80)
    
    # Summary counts
    common_count = len(comparison.common_fields)
    dc_count = len(comparison.datacenter_only_fields)
    campus_count = len(comparison.campus_only_fields)
    
    print(f"\nSummary:")
    print(f"  Common fields: {common_count}")
    print(f"  Datacenter-only fields: {dc_count}")
    print(f"  Campus-only fields: {campus_count}")
    
    if comparison.common_fields:
        print(f"\n{'─' * 40}")
        print("COMMON FIELDS (exist in both schemas):")
        print("─" * 40)
        print_field_tree(comparison.common_fields)
    
    if comparison.datacenter_only_fields:
        print(f"\n{'─' * 40}")
        print("DATACENTER-ONLY FIELDS:")
        print("─" * 40)
        print_field_tree(comparison.datacenter_only_fields)
    
    if comparison.campus_only_fields:
        print(f"\n{'─' * 40}")
        print("CAMPUS-ONLY FIELDS:")
        print("─" * 40)
        print_field_tree(comparison.campus_only_fields)


def generate_go_mode_map(comparisons: List[EndpointComparison]) -> str:
    """Generate Go code for mode field mappings."""
    lines = [
        "// Auto-generated mode field mappings",
        "// Generated by compare_schemas.py",
        "//",
        "// Usage: python3 tools/compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --generate-go > internal/utils/schema.go",
        "",
        "package utils",
        "",
        "type FieldMode string",
        "",
        "const (",
        '    FieldModeBoth       FieldMode = "both"',
        '    FieldModeDatacenter FieldMode = "datacenter"',
        '    FieldModeCampus     FieldMode = "campus"',
        ")",
        "",
        "// FieldAppliesToMode checks if a field applies to the given mode.",
        "// Returns true if the field should be populated for the given mode.",
        "// If the field is not found in ModeFields, it defaults to true (applies to both modes).",
        "func FieldAppliesToMode(resourceType, fieldName, mode string) bool {",
        '    resourceFields, ok := ModeFields[resourceType]',
        '    if !ok {',
        '        // Resource not found in mode map, assume field applies to all modes',
        '        return true',
        '    }',
        '',
        '    fieldMode, ok := resourceFields[fieldName]',
        '    if !ok {',
        '        // Field not found in mode map, assume it applies to all modes',
        '        return true',
        '    }',
        '',
        '    switch fieldMode {',
        '    case FieldModeBoth:',
        '        return true',
        '    case FieldModeDatacenter:',
        '        return mode == "datacenter"',
        '    case FieldModeCampus:',
        '        return mode == "campus"',
        '    default:',
        '        return true',
        '    }',
        "}",
        "",
        "// ModeFields maps field names to their applicable mode",
        "var ModeFields = map[string]map[string]FieldMode{",
    ]
    
    for comparison in comparisons:
        resource_name = comparison.endpoint.strip("/").replace("-", "_")
        lines.append(f'    "{resource_name}": {{')
        
        def add_fields(fields: Dict[str, FieldInfo], prefix: str = ""):
            for name, info in sorted(fields.items()):
                full_name = f"{prefix}{name}" if prefix else name
                mode = "FieldModeBoth"
                if info.mode == FieldMode.DATACENTER_ONLY:
                    mode = "FieldModeDatacenter"
                elif info.mode == FieldMode.CAMPUS_ONLY:
                    mode = "FieldModeCampus"
                lines.append(f'        "{full_name}": {mode},')
                
                if info.nested_fields:
                    add_fields(info.nested_fields, f"{full_name}.")
        
        add_fields(comparison.common_fields)
        add_fields(comparison.datacenter_only_fields)
        add_fields(comparison.campus_only_fields)
        
        lines.append("    },")
    
    lines.append("}")
    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(
        description="Compare datacenter and campus API schemas",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
    # Compare all endpoints
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json

    # Compare a specific endpoint
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --endpoint /services

    # Output to JSON file
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --output report.json

    # Generate Go code for mode mappings
    python compare_schemas.py --datacenter vncgenie.json --campus vnckrait.json --generate-go
        """
    )
    
    parser.add_argument(
        "--datacenter", "-d",
        required=True,
        help="Path to datacenter (vncgenie) JSON schema file"
    )
    parser.add_argument(
        "--campus", "-c",
        required=True,
        help="Path to campus (vnckrait) JSON schema file"
    )
    parser.add_argument(
        "--endpoint", "-e",
        help="Specific endpoint to compare (e.g., /services). If not specified, compares all common endpoints."
    )
    parser.add_argument(
        "--output", "-o",
        help="Output file for JSON report"
    )
    parser.add_argument(
        "--generate-go", "-g",
        action="store_true",
        help="Generate Go code for mode field mappings"
    )
    parser.add_argument(
        "--list-endpoints", "-l",
        action="store_true",
        help="List all endpoints and exit"
    )
    parser.add_argument(
        "--show-common-only",
        action="store_true",
        help="Only show endpoints that exist in both schemas"
    )
    
    args = parser.parse_args()
    
    # Load schemas
    try:
        print(f"Loading datacenter schema: {args.datacenter}", file=sys.stderr)
        dc_schema = load_schema(args.datacenter)
        print(f"Loading campus schema: {args.campus}", file=sys.stderr)
        campus_schema = load_schema(args.campus)
    except FileNotFoundError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON: {e}", file=sys.stderr)
        sys.exit(1)
    
    # Get endpoints
    dc_endpoints = get_endpoints(dc_schema)
    campus_endpoints = get_endpoints(campus_schema)
    common_endpoints = dc_endpoints & campus_endpoints
    
    if args.list_endpoints:
        print("\nDatacenter-only endpoints:")
        for ep in sorted(dc_endpoints - campus_endpoints):
            print(f"  {ep}")
        
        print("\nCampus-only endpoints:")
        for ep in sorted(campus_endpoints - dc_endpoints):
            print(f"  {ep}")
        
        print("\nCommon endpoints:")
        for ep in sorted(common_endpoints):
            print(f"  {ep}")
        
        print(f"\nTotal: {len(dc_endpoints)} datacenter, {len(campus_endpoints)} campus, {len(common_endpoints)} common")
        sys.exit(0)
    
    # Determine which endpoints to compare
    dc_only_endpoints = dc_endpoints - campus_endpoints
    campus_only_endpoints = campus_endpoints - dc_endpoints
    
    if args.endpoint:
        if args.endpoint not in dc_endpoints and args.endpoint not in campus_endpoints:
            print(f"Error: Endpoint {args.endpoint} not found in either schema", file=sys.stderr)
            sys.exit(1)
        endpoints_to_compare = [args.endpoint]
    elif args.show_common_only:
        endpoints_to_compare = sorted(common_endpoints)
    else:
        # Include all endpoints from both schemas
        endpoints_to_compare = sorted(dc_endpoints | campus_endpoints)
    
    # Compare endpoints
    comparisons = []
    for endpoint in endpoints_to_compare:
        if endpoint in common_endpoints:
            comparison = compare_endpoint(endpoint, dc_schema, campus_schema, "common")
            comparisons.append(comparison)
        elif endpoint in dc_only_endpoints:
            comparison = compare_endpoint(endpoint, dc_schema, campus_schema, "datacenter_only")
            comparisons.append(comparison)
        elif endpoint in campus_only_endpoints:
            comparison = compare_endpoint(endpoint, dc_schema, campus_schema, "campus_only")
            comparisons.append(comparison)
    
    # Output results
    if args.output:
        output_data = {
            "datacenter_schema": args.datacenter,
            "campus_schema": args.campus,
            "comparisons": [c.to_dict() for c in comparisons]
        }
        with open(args.output, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, indent=2)
        print(f"Report written to: {args.output}")
    elif args.generate_go:
        go_code = generate_go_mode_map(comparisons)
        print(go_code)
    else:
        for comparison in comparisons:
            print_comparison_report(comparison)
    
    # Print summary
    if not args.output and not args.generate_go:
        print("\n" + "=" * 80)
        print("OVERALL SUMMARY")
        print("=" * 80)
        
        # Count endpoints by type
        common_endpoint_count = 0
        dc_only_endpoint_count = 0
        campus_only_endpoint_count = 0
        
        for c in comparisons:
            has_common = len(c.common_fields) > 0
            has_dc = len(c.datacenter_only_fields) > 0
            has_campus = len(c.campus_only_fields) > 0
            
            if has_common or (has_dc and has_campus):
                common_endpoint_count += 1
            elif has_dc and not has_campus:
                dc_only_endpoint_count += 1
            elif has_campus and not has_dc:
                campus_only_endpoint_count += 1
        
        total_common = sum(len(c.common_fields) for c in comparisons)
        total_dc_only = sum(len(c.datacenter_only_fields) for c in comparisons)
        total_campus_only = sum(len(c.campus_only_fields) for c in comparisons)
        
        print(f"\nEndpoints:")
        print(f"  Common endpoints: {common_endpoint_count}")
        print(f"  Datacenter-only endpoints: {dc_only_endpoint_count}")
        print(f"  Campus-only endpoints: {campus_only_endpoint_count}")
        print(f"  Total endpoints: {len(comparisons)}")
        
        print(f"\nFields:")
        print(f"  Common fields: {total_common}")
        print(f"  Datacenter-only fields: {total_dc_only}")
        print(f"  Campus-only fields: {total_campus_only}")


if __name__ == "__main__":
    main()

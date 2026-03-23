#!/usr/bin/env python3
"""
Convert Terraform add.tf / modify.tf files to JSON PUT / PATCH request bodies for Swagger API testing.

Usage:
    # Single resource type folder:
    python scripts/tf_to_json.py tests/integration/test_cases/campus/badges

    # All resource types in a mode directory:
    python scripts/tf_to_json.py tests/integration/test_cases/campus

    # Filter to a single resource by name:
    python scripts/tf_to_json.py tests/integration/test_cases/campus/ethportprofiles --resource eth_port_profile_test_script1
"""

import argparse
import json
import os
import sys

import hcl2


# Folder name -> JSON wrapper key mapping
FOLDER_TO_WRAPPER_KEY = {
    "tenants":                 "tenant",
    "gateways":                "gateway",
    "gatewayprofiles":         "gateway_profile",
    "ethportprofiles":         "eth_port_profile_",
    "lags":                    "lag",
    "sflowcollectors":         "sflow_collector",
    "diagnosticsprofiles":     "diagnostics_profile",
    "diagnosticsportprofiles": "diagnostics_port_profile",
    "services":                "service",
    "ethportsettings":         "eth_port_settings",
    "bundles":                 "endpoint_bundle",
    "badges":                  "badge",
    "authenticatedethports":   "authenticated_eth_port",
    "devicevoicesettings":     "device_voice_settings",
    "packetbroker":            "pb_egress_profile",
    "packetqueues":            "packet_queue",
    "serviceportprofiles":     "service_port_profile",
    "voiceportprofiles":       "voice_port_profiles",
    "switchpoints":            "switchpoint",
    "devicecontrollers":       "device_controller",
    "aspathaccesslists":       "as_path_access_list",
    "communitylists":          "community_list",
    "devicesettings":          "eth_device_profiles",
    "extendedcommunitylists":  "extended_community_list",
    "ipv4lists":               "ipv4_list_filter",
    "ipv4prefixlists":         "ipv4_prefix_list",
    "ipv6lists":               "ipv6_list_filter",
    "ipv6prefixlists":         "ipv6_prefix_list",
    "routemapclauses":         "route_map_clause",
    "routemaps":               "route_map",
    "sfpbreakouts":            "sfp_breakouts",
    "sites":                   "site",
    "pods":                    "pod",
    "spineplanes":             "spine_plane",
    "policybasedroutingacl":   "pb_routing_acl",
    "policybasedrouting":      "pb_routing",
    "portacls":                "port_acl",
    "groupingrules":           "grouping_rules",
    "thresholdgroups":         "threshold_group",
    "thresholds":              "threshold",
    "aclsipv4":                "ip_filter",
    "aclsipv6":                "ip_filter",
}

# Folder names for ACL resources that need ip_version query parameter
ACL_FOLDERS = {
    "acls_ipv4": "4",
    "acls_ipv6": "6",
}

# Resources that cannot be created via PUT
NON_CREATABLE = {"sites", "sfpbreakouts"}

# Nested blocks that should be unwrapped from a list to a single object
SINGLE_OBJECT_BLOCKS = {"object_properties"}

# Terraform meta-arguments to strip from resource attributes
META_ARGUMENTS = {"depends_on", "provider", "lifecycle", "count", "for_each"}


def normalize_folder_name(folder_name: str) -> str:
    """Normalize folder name by removing underscores for lookup in the mapping dict."""
    return folder_name.replace("_", "")


def get_wrapper_key(folder_name: str) -> str | None:
    """Get the JSON wrapper key for a folder name."""
    normalized = normalize_folder_name(folder_name)
    return FOLDER_TO_WRAPPER_KEY.get(normalized)


def process_value(val):
    """Recursively process a parsed HCL value to clean it for JSON output."""
    if isinstance(val, list):
        return [process_value(item) for item in val]
    elif isinstance(val, dict):
        return {k: process_value(v) for k, v in val.items()}
    else:
        return val


def transform_resource(attrs: dict) -> dict:
    """Transform a single resource's attributes from HCL parsed format to API JSON format.

    - Strips meta-arguments (depends_on, lifecycle, etc.)
    - Unwraps object_properties from list to single object
    - Keeps other nested block arrays as-is
    """
    result = {}
    for key, value in attrs.items():
        if key in META_ARGUMENTS:
            continue

        if key in SINGLE_OBJECT_BLOCKS:
            if isinstance(value, list) and len(value) > 0 and isinstance(value[0], dict):
                result[key] = process_value(value[0])
            else:
                result[key] = process_value(value)
        else:
            result[key] = process_value(value)

    return result


def parse_tf_file(filepath: str) -> list[tuple[str, str, dict]]:
    """Parse a .tf file and return list of (resource_type, resource_name, attributes)."""
    with open(filepath, "r") as f:
        parsed = hcl2.load(f)

    resources = []
    for resource_block in parsed.get("resource", []):
        for resource_type, instances in resource_block.items():
            for instance_name, attrs in instances.items():
                resources.append((resource_type, instance_name, attrs))

    return resources


def process_add_tf(folder_path: str, resource_filter: str | None = None) -> dict | None:
    """Process an add.tf file in the given folder and return the PUT body JSON structure.

    Returns None if no add.tf found or folder is non-creatable.
    """
    folder_name = os.path.basename(folder_path)

    if folder_name in NON_CREATABLE:
        print(f"# SKIPPED: {folder_name} — cannot be created via PUT (read-only/PATCH-only)",
              file=sys.stderr)
        return None

    add_tf_path = os.path.join(folder_path, "add.tf")
    if not os.path.isfile(add_tf_path):
        return None

    wrapper_key = get_wrapper_key(folder_name)
    if not wrapper_key:
        print(f"# WARNING: No wrapper key mapping found for folder '{folder_name}'",
              file=sys.stderr)
        return None

    resources = parse_tf_file(add_tf_path)
    if not resources:
        return None

    resource_map = {}
    for _resource_type, _instance_name, attrs in resources:
        transformed = transform_resource(attrs)
        resource_name = transformed.get("name", _instance_name)

        if resource_filter and resource_name != resource_filter:
            continue

        resource_map[resource_name] = transformed

    if not resource_map:
        return None

    return {wrapper_key: resource_map}


def process_modify_tf(folder_path: str, resource_filter: str | None = None) -> dict | None:
    """Process a modify.tf file in the given folder and return the PATCH body JSON structure.

    modify.tf resources have no 'name' attribute — the Terraform instance name is the resource name.
    Only changed fields are present (partial update).

    Returns None if no modify.tf found.
    """
    folder_name = os.path.basename(folder_path)

    modify_tf_path = os.path.join(folder_path, "modify.tf")
    if not os.path.isfile(modify_tf_path):
        return None

    wrapper_key = get_wrapper_key(folder_name)
    if not wrapper_key:
        print(f"# WARNING: No wrapper key mapping found for folder '{folder_name}'",
              file=sys.stderr)
        return None

    resources = parse_tf_file(modify_tf_path)
    if not resources:
        return None

    resource_map = {}
    for _resource_type, instance_name, attrs in resources:
        transformed = transform_resource(attrs)

        if resource_filter and instance_name != resource_filter:
            continue

        resource_map[instance_name] = transformed

    if not resource_map:
        return None

    return {wrapper_key: resource_map}


def main():
    parser = argparse.ArgumentParser(
        description="Convert Terraform add.tf / modify.tf files to JSON PUT / PATCH request bodies for Swagger"
    )
    parser.add_argument(
        "path",
        help="Path to a resource type folder (e.g., test_cases/campus/badges) "
             "or a parent folder (e.g., test_cases/campus) to process all resource types"
    )
    parser.add_argument(
        "--resource",
        help="Filter to a single resource by its 'name' attribute value",
        default=None,
    )
    args = parser.parse_args()

    target_path = os.path.abspath(args.path)

    if not os.path.isdir(target_path):
        print(f"Error: '{args.path}' is not a directory", file=sys.stderr)
        sys.exit(1)

    has_tf = (os.path.isfile(os.path.join(target_path, "add.tf"))
              or os.path.isfile(os.path.join(target_path, "modify.tf")))
    if has_tf:
        folders = [target_path]
    else:
        folders = sorted([
            os.path.join(target_path, d)
            for d in os.listdir(target_path)
            if os.path.isdir(os.path.join(target_path, d))
        ])

    FOLDER_TO_ENDPOINT = {
        "acls_ipv4": "/acls",
        "acls_ipv6": "/acls",
    }

    first = True
    for folder in folders:
        folder_name = os.path.basename(folder)

        put_result = process_add_tf(folder, args.resource)
        patch_result = process_modify_tf(folder, args.resource)

        if put_result is None and patch_result is None:
            continue

        endpoint = FOLDER_TO_ENDPOINT.get(folder_name, f"/{normalize_folder_name(folder_name)}")
        acl_suffix = ""
        if folder_name in ACL_FOLDERS:
            acl_suffix = f"?ip_version={ACL_FOLDERS[folder_name]}"

        if put_result is not None:
            if not first:
                print()
            first = False
            print(f"=== {folder_name} → PUT {endpoint}{acl_suffix} ===")
            print(json.dumps(put_result, indent=2))

        if patch_result is not None:
            if not first:
                print()
            first = False
            print(f"=== {folder_name} → PATCH {endpoint}{acl_suffix} ===")
            print(json.dumps(patch_result, indent=2))


if __name__ == "__main__":
    main()

# Verity Terraform Provider Documentation


## 1. Provider Configuration

The Verity provider offers flexible configuration options, allowing you to specify credentials through provider configuration blocks or (recommended) environment variables.

### Recommended: Using Environment Variables

Export the following environment variables before running Terraform:

```bash
# For Linux/macOS
export TF_VAR_uri="<your-verity-uri>"
export TF_VAR_username="<your-username>"
export TF_VAR_password="<your-password>"
```
```powershell
# For Windows PowerShell
$env:TF_VAR_uri="<your-verity-uri>"
$env:TF_VAR_username="<your-username>"
$env:TF_VAR_password="<your-password>"
```

Then use a minimal provider block:

```hcl
terraform {
  required_providers {
    verity = {
      source  = "BE-Network/verity"
      version = "6.4.0" # Replace with the desired release version
    }
  }
}

provider "verity" {
  mode = "datacenter" # Valid values: "datacenter" or "campus"
}
```

> Replace `6.4.0` with the desired release version. Set `mode` to match your Verity deployment type.

If a configuration value is not specified in the provider block, the provider will automatically look for it in the corresponding environment variable. For security, do not write sensitive values (like username and password) directly in your configuration files.


### Required CLI Parallelism

The Verity provider requires the following environment variable to be set at all times to allow the provider to make bulk requests to the API. If this variable is not set, Terraform will use the default parallelism of 10, which will significantly slow down the provider:

#### Unix-based Systems
```bash
export TF_CLI_ARGS_apply="-parallelism=500"
```

#### Windows
```powershell
$env:TF_CLI_ARGS_apply="-parallelism=500"
```

Make sure to set these environment variables before running any Terraform commands.

## 2. Resource Types

The provider supports the following resource types:

- `verity_acl_v4`
- `verity_acl_v6`
- `verity_as_path_access_list`
- `verity_authenticated_eth_port`
- `verity_badge`
- `verity_bundle`
- `verity_community_list`
- `verity_device_controller`
- `verity_device_settings`
- `verity_device_voice_settings`
- `verity_diagnostics_port_profile`
- `verity_diagnostics_profile`
- `verity_eth_port_profile`
- `verity_eth_port_settings`
- `verity_extended_community_list`
- `verity_gateway`
- `verity_gateway_profile`
- `verity_grouping_rule`
- `verity_ipv4_list`
- `verity_ipv4_prefix_list`
- `verity_ipv6_list`
- `verity_ipv6_prefix_list`
- `verity_lag`
- `verity_operation_stage`
- `verity_packet_broker`
- `verity_packet_queue`
- `verity_pb_routing`
- `verity_pb_routing_acl`
- `verity_pod`
- `verity_port_acl`
- `verity_route_map`
- `verity_route_map_clause`
- `verity_service`
- `verity_service_port_profile`
- `verity_sflow_collector`
- `verity_sfp_breakout`
- `verity_site`
- `verity_spine_plane`
- `verity_switchpoint`
- `verity_tenant`
- `verity_threshold`
- `verity_threshold_group`
- `verity_voice_port_profile`

Each resource type has specific attributes and configurations. See the resource documentation for detailed usage.

## 3. State Importer

The provider includes a state importer data source that helps you import existing Verity configurations into your Terraform state.
You don't have to invoke this data source manually; you can instead run the `import_verity_state` scripts (see Tools section) to automate the export and import workflow.

```hcl
data "verity_state_importer" "import" {
  output_dir = "/path/to/directory"  # defaults to current working directory if not specified
}
```

The state importer workflow:

1. **Configuration Export**: The importer connects to your Verity instance and exports the current configuration
2. **Resource Generation**: It automatically generates Terraform resource files (`.tf`) that map your current Verity configuration to Terraform resources
3. **Import Blocks Generation**: An `import_blocks.tf` file is automatically generated containing import blocks for all resources.
4. **Import Process**: Run `terraform apply` to import all resources at once using the generated import blocks

### Generated Files
The importer generates the following Terraform resource files:
- `acls_ipv4.tf`
- `acls_ipv6.tf`
- `aspathaccesslists.tf`
- `authenticatedethports.tf`
- `badges.tf`
- `bundles.tf`
- `communitylists.tf`
- `devicecontrollers.tf`
- `devicesettings.tf`
- `devicevoicesettings.tf`
- `diagnosticsportprofiles.tf`
- `diagnosticsprofiles.tf`
- `ethportprofiles.tf`
- `ethportsettings.tf`
- `extendedcommunitylists.tf`
- `gatewayprofiles.tf`
- `gateways.tf`
- `groupingrules.tf`
- `ipv4lists.tf`
- `ipv4prefixlists.tf`
- `ipv6lists.tf`
- `ipv6prefixlists.tf`
- `lags.tf`
- `packetbroker.tf`
- `packetqueues.tf`
- `policybasedrouting.tf`
- `policybasedroutingacl.tf`
- `pods.tf`
- `portacls.tf`
- `routemapclauses.tf`
- `routemaps.tf`
- `services.tf`
- `serviceportprofiles.tf`
- `sflowcollectors.tf`
- `sfpbreakouts.tf`
- `sites.tf`
- `spineplanes.tf`
- `switchpoints.tf`
- `tenants.tf`
- `thresholdgroups.tf`
- `thresholds.tf`
- `voiceportprofiles.tf`
- `import_blocks.tf` - Import blocks for all resources in the correct dependency order
- `stages.tf` - Resource dependency ordering


### Resource Dependency Management

The import process creates a special `stages.tf` file that defines explicit dependency ordering for resources. This uses the `verity_operation_stage` resource, which helps to:

1. Establish a clear sequence for creating, updating, and destroying resources
2. Prevent dependency conflicts between resource types
3. Ensure that resources are processed in the optimal order for the Verity API

Each imported resource is configured with the appropriate `depends_on` attribute referring to its corresponding stage. This prevents Terraform from attempting to create resources before their dependencies are ready, which is particularly important when working with the Verity API's interdependent resources.

Since API version 6.5, the provider supports two modes: **campus** and **datacenter**. Each mode has its own resource dependency ordering for creation and update operations:

**Order for CAMPUS:**
1. IPv4 Lists
2. IPv6 Lists
3. ACLs (IPv4)
4. ACLs (IPv6)
5. PB Routing ACL
6. PB Routing
7. Port ACLs
8. Services
9. Eth Port Profiles
10. SFlow Collectors
11. Packet Queues
12. Service Port Profiles
13. Diagnostics Port Profiles
14. Device Voice Settings
15. Authenticated Eth-Ports
16. Diagnostics Profiles
17. Eth Port Settings
18. Voice-Port Profiles
19. Device Settings
20. Lags
21. Bundles
22. Badges
23. Switchpoints
24. Thresholds
25. Grouping Rules
26. Threshold Groups
27. Sites
28. Device Controllers

**Order for DATACENTER:**
1. SFP Breakouts
2. IPv6 Prefix Lists
3. Community Lists
4. IPv4 Prefix Lists
5. Extended Community Lists
6. AS Path Access Lists
7. Route Map Clauses
8. ACLs (IPv6)
9. ACLs (IPv4)
10. Route Maps
11. PB Routing ACL
12. Tenants
13. PB Routing
14. IPv4 Lists
15. IPv6 Lists
16. Services
17. Port ACLs
18. Packet Broker
19. Eth Port Profiles
20. Packet Queues
21. SFlow Collectors
22. Gateways
23. Lags
24. Eth Port Settings
25. Diagnostics Profiles
26. Gateway Profiles
27. Diagnostics Port Profiles
28. Bundles
29. Pods
30. Badges
31. Spine Planes
32. Switchpoints
33. Device Settings
34. Thresholds
35. Grouping Rules
36. Threshold Groups
37. Sites
38. Device Controllers

For delete operations, the order is automatically reversed to ensure proper dependency handling when removing resources.

#### Creating New Resources

When manually creating new resources (not through import), it's strongly recommended to follow the same pattern and include the appropriate `depends_on` attribute referring to the corresponding stage. For example:

```hcl
resource "verity_tenant" "example" {
  name = "example-tenant"
  // other attributes...
  depends_on = [verity_operation_stage.tenant_stage]
}

resource "verity_service" "example" {
  name = "example-service"
  // other attributes...
  depends_on = [verity_operation_stage.service_stage]
}
```

This ensures proper ordering of operations and helps avoid dependency issues when managing your infrastructure.


## 4. Tools

When you run `terraform init`, the provider binary and a `tools` folder are placed alongside the plugin in the `.terraform/providers` directory. This `tools` folder contains:

- `import_verity_state.sh`
- `import_verity_state.ps1`

To import your existing Verity system state into Terraform, run the appropriate script from your Terraform project directory:

**Linux/macOS**:

```bash
# Production (Linux/macOS)
.terraform/providers/registry.terraform.io/be-network/verity/<VERSION>/<OS>_<ARCH>/tools/import_verity_state.sh
```

**Windows PowerShell**:
```powershell
# Production (Windows)
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1
```
> **Tip:** You can always locate the `tools` folder inside the provider directory (where the provider binary is installed). Use your file browser or terminal to navigate to the correct folder, then right-click the script and choose "Copy Path" (Linux/macOS) or "Copy as path" (Windows) to avoid manually typing the full path.

> **Note:** Replace:
> - `<VERSION>` with the actual provider version (e.g. `6.4.0`)
> - `<OS>` with your operating system (e.g. `linux`, `windows`, `darwin`)
> - `<ARCH>` with your CPU architecture (e.g. `amd64`, `arm64`)

## 5. Handling Auto-Assigned Fields

When you change an auto-assigned field's flag (such as `auto_assigned_vni`, `auto_assigned_vlan`, etc.) from `false` to `true`, you must remove the corresponding field (such as `vni`, `vlan`, etc.) from your Terraform resource block. Leaving the field present will cause issues, as the backend will automatically assign its value and may overwrite or ignore the value you specify in Terraform.

Our `data_source_state_importer` is designed to check if a field has a corresponding auto-assigned flag. If the flag is set to `true`, the importer will not write that field in the generated Terraform resource file — only the auto-assigned flag will be present. This ensures your configuration matches the backend's behavior and avoids conflicts.

**Best Practice:**
Whenever you enable auto-assignment for a field, always remove the manually specified value for that field from your `.tf` resource block.
# Verity Terraform Provider Documentation

## Supported Versions

| Provider version | Supported Verity API versions |
| --- | --- |
| 6.6.100 | 6.6.0.269 and up |

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
      version = "6.6.0" # Replace with the desired release version
    }
  }
}

provider "verity" {
  mode = "datacenter" # Valid values: "datacenter" or "campus"
}
```

> Replace `6.6.0` with the desired release version. Set `mode` to match your Verity deployment type.

If a configuration value is not specified in the provider block, the provider will automatically look for it in the corresponding environment variable. For security, do not write sensitive values (like username and password) directly in your configuration files.


### Parallelism Configuration (Important)

The Verity provider uses a **bulk operations architecture** — all resources of a given type are collected and sent to the API in a single request. For this to work correctly, Terraform's parallelism must be set **higher than the total number of resources affected in a single `terraform apply` run** (creates + updates + deletes combined).

For example, if one apply run adds 300 resources, updates 100, and deletes 50, the total is 450 — so `parallelism=500` is sufficient. This ensures all affected resources start concurrently, queue their operations, and each type is sent to the API in a single request.

**Recommended: `parallelism=500`** — sufficient for most deployments (up to ~500 total affected resources per apply). For larger environments, use `parallelism=1000` or `parallelism=2000`.

> **Why is high parallelism safe?** Terraform parallelism controls Go goroutines. Each resource goroutine blocks on a wait channel (sleeping) until the bulk operation for its type executes — consuming negligible CPU and memory. Values of 1000–2000 are perfectly safe even on modest hardware.
>
> **What happens with low parallelism?** If parallelism is lower than the total number of affected resources, resources are processed in waves. This means some types may be split across multiple API calls. While this is usually harmless, it can slow down execution and may cause issues if resources of the same type reference each other across batches. Cross-type ordering is always preserved regardless of parallelism.

If this variable is not set, Terraform will use the default parallelism of 10, which will significantly slow down the provider and may cause batch-splitting issues:

#### Unix-based Systems
```bash
export TF_CLI_ARGS_apply="-parallelism=500"
```

#### Windows
```powershell
$env:TF_CLI_ARGS_apply="-parallelism=500"
```

For large deployments (more than 500 affected resources per apply):
```bash
export TF_CLI_ARGS_apply="-parallelism=2000"
```

Make sure to set these environment variables before running any Terraform commands.

## 2. Resource Types

The provider supports the following resource types:

- `verity_aaa_profile`
- `verity_acl_v4`
- `verity_acl_v6`
- `verity_as_path_access_list`
- `verity_authenticated_eth_port`
- `verity_badge`
- `verity_bundle`
- `verity_community_list`
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
- `verity_ldap_profile`
- `verity_mac_filter`
- `verity_operation_stage`
- `verity_packet_broker`
- `verity_packet_queue`
- `verity_pair`
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
- `verity_ssp_group`
- `verity_fabric`
- `verity_plane`
- `verity_rack`
- `verity_spine_plane`
- `verity_switchpoint`
- `verity_su`
- `verity_tacacs_profile`
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
- `deviceaaaprofiles.tf`
- `ldapprofiles.tf`
- `acls_ipv4.tf`
- `acls_ipv6.tf`
- `aspathaccesslists.tf`
- `authenticatedethports.tf`
- `badges.tf`
- `bundles.tf`
- `communitylists.tf`
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
- `macfilters.tf`
- `packetbroker.tf`
- `packetqueues.tf`
- `pairs.tf`
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
- `fabrics.tf`
- `planes.tf`
- `racks.tf`
- `spineplanes.tf`
- `sspgroups.tf`
- `switchpoints.tf`
- `sus.tf`
- `tacacsprofiles.tf`
- `tenants.tf`
- `thresholdgroups.tf`
- `thresholds.tf`
- `voiceportprofiles.tf`
- `import_blocks.tf` - Import blocks for all resources in the correct dependency order
- `stages.tf` - Resource dependency ordering


### Resource Dependency Management

The import process creates a special `stages.tf` file that defines explicit dependency ordering for resources. This uses the `verity_operation_stage` resource, which acts as an **active barrier** between resource type groups:

1. Establish a clear sequence for creating, updating, and destroying resources
2. Prevent dependency conflicts between resource types
3. Ensure that resources are processed in the optimal order for the Verity API
4. Wait for all operations from the current type group to complete before allowing the next group to start

Each imported resource is configured with the appropriate `depends_on` attribute referring to its corresponding stage. When a stage's `Create` is executed, it actively waits for its sibling resources to queue their operations, flushes them to the API, and only returns once all operations for that type group are complete. This guarantees sequential, ordered API execution regardless of Terraform's internal scheduling.

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
> - `<VERSION>` with the actual provider version (e.g. `6.6.0`)
> - `<OS>` with your operating system (e.g. `linux`, `windows`, `darwin`)
> - `<ARCH>` with your CPU architecture (e.g. `amd64`, `arm64`)

## 5. Handling Auto-Assigned Fields

When you change an auto-assigned field's flag (such as `auto_assigned_vni`, `auto_assigned_vlan`, etc.) from `false` to `true`, you must remove the corresponding field (such as `vni`, `vlan`, etc.) from your Terraform resource block. Leaving the field present will cause issues, as the backend will automatically assign its value and may overwrite or ignore the value you specify in Terraform.

Our `data_source_state_importer` is designed to check if a field has a corresponding auto-assigned flag. If the flag is set to `true`, the importer will not write that field in the generated Terraform resource file — only the auto-assigned flag will be present. This ensures your configuration matches the backend's behavior and avoids conflicts.

**Best Practice:**
Whenever you enable auto-assignment for a field, always remove the manually specified value for that field from your `.tf` resource block.

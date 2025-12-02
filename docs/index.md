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
      version = "6.4.0"
    }
  }
}

provider "verity" {
  mode = "datacenter" # API 6.4 supports datacenter mode only
}
```

If a configuration value is not specified in the provider block, the provider will automatically look for it in the corresponding environment variable. For security, do not write sensitive values (like username and password) directly in your configuration files.


### Required CLI Parallelism

The Verity provider requires the following environment variable to be set at all times to allow the provider to make bulk requests to the API. If this variable is not set, Terraform will use the default parallelism of 10, which will significantly slow down the provider:

#### Unix-based Systems
```bash
export TF_CLI_ARGS_apply="-parallelism=250"
```

#### Windows
```powershell
$env:TF_CLI_ARGS_apply="-parallelism=250"
```

Make sure to set these environment variables before running any Terraform commands.

## 2. Resource Types

The provider supports the following resource types for API version 6.4 (datacenter mode only):

- `verity_bundle`
- `verity_eth_port_profile`
- `verity_eth_port_settings`
- `verity_gateway`
- `verity_gateway_profile`
- `verity_lag`
- `verity_service`
- `verity_tenant`
- `verity_operation_stage`

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
- `bundles.tf`
- `ethportprofiles.tf`
- `ethportsettings.tf`
- `gatewayprofiles.tf`
- `gateways.tf`
- `lags.tf`
- `services.tf`
- `tenants.tf`
- `import_blocks.tf` - Import blocks for all resources in the correct dependency order
- `stages.tf` - Resource dependency ordering


### Resource Dependency Management

The import process creates a special `stages.tf` file that defines explicit dependency ordering for resources. This uses the `verity_operation_stage` resource, which helps to:

1. Establish a clear sequence for creating, updating, and destroying resources
2. Prevent dependency conflicts between resource types
3. Ensure that resources are processed in the optimal order for the Verity API

Each imported resource is configured with the appropriate `depends_on` attribute referring to its corresponding stage. This prevents Terraform from attempting to create resources before their dependencies are ready, which is particularly important when working with the Verity API's interdependent resources.

API version 6.4 supports **datacenter mode only**. The resource dependency ordering for creation and update operations is as follows:

**Order for DATACENTER (API 6.4):**
1. Tenants
2. Gateways
3. Gateway Profiles
4. Services
5. Eth Port Profiles
6. Eth Port Settings
7. Lags
8. Bundles

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
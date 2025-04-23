# Verity Terraform Provider Documentation

## 1. Provider Configuration

The provider requires basic authentication parameters in your Terraform configuration:

```hcl
terraform {
  required_providers {
    verity = {
      source  = "BE-Network/verity"
      version = "<VERSION>"
    }
  }
}

provider "verity" {
  uri      = var.uri
  username = var.username
  password = var.password
}
```
> **Note:** Replace `<VERSION>` with the actual provider version (e.g. `1.0.3`)

Required parameters:
- **uri**: Base URL for the API
- **username** and **password**: For authentication

You can export these variables using environment variables:

```bash
# For Linux/MacOS
export TF_VAR_uri="https://your-verity-instance"
export TF_VAR_username="your-username"
export TF_VAR_password="your-password"

# For Windows PowerShell
$env:TF_VAR_uri="https://your-verity-instance"
$env:TF_VAR_username="your-username"
$env:TF_VAR_password="your-password"
```

### CLI Parallelism

You can adjust Terraform’s parallelism for apply operations:

```bash
export TF_CLI_ARGS_apply="-parallelism=250"
```

For Windows PowerShell:

```powershell
$env:TF_CLI_ARGS_apply="-parallelism=250"
```

## 2. Resource Types

The provider supports the following resource types:

- `verity_tenant`: Manage Verity tenants
- `verity_gateway_profile`: Configure gateway profiles
- `verity_eth_port_profile`: Manage Ethernet port profiles
- `verity_lag`: Configure Link Aggregation Groups
- `verity_service`: Manage services
- `verity_eth_port_settings`: Configure Ethernet port settings
- `verity_bundle`: Manage bundles
- `verity_gateway`: Configure gateways

Each resource type has specific attributes and configurations. See the examples section below for detailed usage.

## 3. State Importer

The provider includes a state importer data source that helps you import existing Verity configurations into your Terraform state.
You don't have to invoke this data source manually; you can instead run the `import_verity_state` scripts (see Tools section) to automate the export and import workflow.

```hcl
data "verity_state_importer" "import" {
  output_dir = var.config_dir  # defaults to current working directory if not specified
}
```

The state importer workflow:

1. **Configuration Export**: The importer connects to your Verity instance and exports the current configuration
2. **Resource Generation**: It automatically generates Terraform resource files (`.tf`) that map your current Verity configuration to Terraform resources
3. **Import Blocks Generation**: An `import_blocks.tf` file is automatically generated containing import blocks for all resources.
4. **Import Process**: Run `terraform apply` to import all resources at once using the generated import blocks

### Generated Files
The importer generates the following Terraform resource files:
- `bundles.tf` - Bundle resources
- `ethportprofiles.tf` - Ethernet port profile resources
- `ethportsettings.tf` - Ethernet port settings resources
- `gatewayprofiles.tf` - Gateway profile resources
- `gateways.tf` - Gateway resources
- `lags.tf` - LAG resources
- `services.tf` - Service resources
- `tenants.tf` - Tenant resources
- `import_blocks.tf` - Import blocks for all resources in the correct dependency order


## 4. Tools

When you run `terraform init`, the provider binary and a `tools` folder are placed alongside the plugin in the `.terraform/providers` directory. This `tools` folder contains:

- `import_verity_state.sh`
- `import_verity_state.ps1`

To import your existing Verity system state into Terraform, run the appropriate script from your Terraform project directory:

**Linux/macOS**:
```bash
.terraform/providers/registry.terraform.io/be-network/verity/<VERSION>/<OS>_<ARCH>/tools/import_verity_state.sh
```

**Windows PowerShell**:
```powershell
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1
```
> **Note:** Replace:
> - `<VERSION>` with the actual provider version (e.g. `1.0.3`)
> - `<OS>` with your operating system (e.g. `linux`, `windows`, `darwin`)
> - `<ARCH>` with your CPU architecture (e.g. `amd64`, `arm64`)

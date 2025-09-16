# Local Verity Terraform Provider Setup

## Building the Provider

To compile the provider binary, run the following command in the root directory of the project:

 For macOS/Linux:
```bash
go build -o terraform-provider-verity
```

For Windows:

```bash
go build -o terraform-provider-verity.exe
```


## Provider Configuration

### Custom Provider Binary

To use a local development version of the provider, you need to configure Terraform to use your custom provider binary instead of downloading it from the registry. This is done using a `.tfrc` configuration file. Here's an example (`dev.tfrc`):


#### Example for macOS/Linux
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/local/verity" = "/home/<user>/terraform-provider-verity"
  }
  direct {}
}
```

#### Example for Windows
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/local/verity" = "C:\\Users\\<user>\\terraform-provider-verity.exe"
  }
  direct {}
}
```

Then in your Terraform configuration file, specify the local provider:

```hcl
terraform {
  required_providers {
    verity = {
      source = "registry.terraform.io/local/verity"
    }
  }
}

provider "verity" {
  mode = mode = "datacenter" # Valid values: "datacenter" or "campus"
}
```


To use this configuration, set the `TF_CLI_CONFIG_FILE` environment variable to point to your custom `.tfrc` file:
  
  For macOS/Linux
  ```bash
  export TF_CLI_CONFIG_FILE=/home/<user>/terraform-provider-verity/examples/dev.tfrc
  ```
  For Windows:
  ```powershell
  $env:TF_CLI_CONFIG_FILE="C:\path\to\terraform-provider-verity\examples\dev.tfrc"
  ```

### Required Environment Variables

The Verity provider offers flexible configuration options, allowing you to specify credentials through provider configuration blocks, variable files, or environment variables:

- `TF_VAR_uri`: The base URL of the Verity API
- `TF_VAR_username`: Your Verity API username
- `TF_VAR_password`: Your Verity API password


You can configure the Verity provider in two ways:

1. **Recommended: Export environment variables**
  Export the following environment variables before running Terraform:
  ```bash
  export TF_VAR_uri="<your-verity-uri>"
  export TF_VAR_username="<your-username>"
  export TF_VAR_password="<your-password>"
  ```
  Then use a minimal provider block:
  ```terraform
  provider "verity" {}
  ```
  All configuration is taken from environment variables.

2. **Alternative: Specify fields directly in the provider block**
  ```terraform
  provider "verity" {
    uri = "<your-verity-uri>"
    # username and password should NOT be written in plain text here
    # prefer environment variables for sensitive values
  }
  ```
  You may specify any or all fields directly. If a field is not specified, the provider will look for it in the corresponding environment variable. For security, do not write sensitive values (like username and password) directly in your configuration files.

For Linux and macOS, use the following commands to set environment variables:

```bash
export TF_VAR_uri="<your-verity-uri>"
export TF_VAR_username="<your-username>"
export TF_VAR_password="<your-password>"
```

For Windows, use the following commands to set environment variables:

```powershell
$env:TF_VAR_uri="<your-verity-uri>"
$env:TF_VAR_username="<your-username>"
$env:TF_VAR_password="<your-password>"
```

Additionally, the Verity provider requires the following environment variable to be set at all times to allow the provider to make bulk requests to the API. If this variable is not set, Terraform will use the default parallelism of 10, which will significantly slow down the provider:

#### Unix-based Systems
```bash
export TF_CLI_ARGS_apply="-parallelism=250"
```

#### Windows
```powershell
$env:TF_CLI_ARGS_apply="-parallelism=250"
```


Make sure to set these environment variables before running any Terraform commands.


## Production Setup


For production use, configure Terraform to download the provider from the registry. You can find all available versions on the HashiCorp Registry:

**Provider Versions:** [View all releases on HashiCorp Registry](https://registry.terraform.io/providers/BE-Network/verity/latest)

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



## Regenerating the OpenAPI Go SDK

To regenerate the OpenAPI SDK for Go:

1. Get the latest Swagger JSON file(s)
2. **Remove unnecessary endpoints**: Before processing the Swagger file, use the `remove_endpoints.py` script to remove API endpoints that are not needed for the Terraform provider. This script removes endpoints like `/config`, `/changesets`, `/readmode`, `/request`, `/snmp`, `/syslog`, `/backups`, and `/timetraveler` which are not relevant for Terraform resource management:
   ```bash
   python3 tools/remove_endpoints.py swagger.json
   ```
   This will create a backup of your original file (as `swagger.json.backup`) and modify the original file to exclude the unnecessary endpoints.

3. If you have multiple Swagger JSON files that need to be combined (e.g., from different API versions), use the `merge_swagger.py` script to merge them:
   ```bash
   python3 tools/merge_swagger.py <file1> [<file2> ...] --output swagger_merged.json
   ```
   This will combine multiple Swagger specifications into a unified file, merging paths, components, and other sections. For example:
   ```bash
   python3 tools/merge_swagger.py swagger_65_dc.json swagger_65_campus.json --output swagger_merged.json
   ```
   **Note**: If you used the merge script, make sure to also run the `remove_endpoints.py` script on the merged file:
   ```bash
   python3 tools/remove_endpoints.py swagger_merged.json
   ```

4. Transform the Swagger file to be compatible with the Go SDK generator:
   ```bash
   python3 tools/transform_swagger.py swagger.json
   ```
   Or if you used the merge script:
   ```bash
   python3 tools/transform_swagger.py swagger_merged.json
   ```

5. Install the OpenAPI Generator CLI (if not already installed):
   ```bash
   npm install @openapitools/openapi-generator-cli -g
   ```

6. Generate the Go SDK using openapi-generator-cli:
   ```bash
   openapi-generator-cli generate -i swagger_transformed.json -g go -o ./openapi
   ```

7. Replace the existing openapi folder with the newly generated one


### Updating Provider Resource Files

After regenerating the SDK, you need to update the provider resource files:

- For fields deleted from the API: Remove them from the corresponding provider resource files
- For new fields added to the API: Add them to the appropriate provider resource files

## Using the State Import Scripts

The provider includes scripts to help import existing Verity resources into Terraform state. These scripts automate the process of creating resource files and importing existing resources.

### Resource Dependency Management

The import process creates a special `stages.tf` file that defines explicit dependency ordering for resources. This uses the `verity_operation_stage` resource, which helps to:

1. Establish a clear sequence for creating, updating, and destroying resources
2. Prevent dependency conflicts between resource types
3. Ensure that resources are processed in the optimal order for the Verity API

Each imported resource is configured with the appropriate `depends_on` attribute referring to its corresponding stage. This prevents Terraform from attempting to create resources before their dependencies are ready, which is particularly important when working with the Verity API's interdependent resources.

The operation stages maintain the following order for creation and update operations:

1. tenant_stage
2. gateway_stage
3. gateway_profile_stage
4. service_stage
5. eth_port_profile_stage
6. eth_port_settings_stage
7. lag_stage
8. bundle_stage

For delete operations, this order is automatically reversed (8→1) to ensure proper dependency handling when removing resources.

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

### What the Scripts Do

1. Find the main Terraform file with the Verity provider
2. Add the `verity_state_importer` data source if it doesn't exist
3. Run a first `terraform apply` to generate resource files and import blocks
4. Run a second `terraform apply` to import the resources into your state
5. Clean up temporary files

### Running the Scripts

#### Linux and macOS


```bash
# Production
.terraform/providers/registry.terraform.io/be-network/verity/<VERSION>/<OS>_<ARCH>/tools/import_verity_state.sh

# Local (skip terraform init)
../tools/import_verity_state.sh --local
```

> **Tip:** For production, you can always locate the `tools` folder inside the provider directory (where the provider binary is installed). Use your file browser or terminal to navigate to the correct folder, then right-click the script and choose "Copy Path" to avoid manually typing the full path.

#### Windows PowerShell

```powershell
# Production
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1

# Local (skip terraform init)
..\tools\import_verity_state.ps1 -Local
```

> **Tip:** On Windows, you can use File Explorer to navigate to the provider's `tools` folder, then right-click the script and select "Copy as path" to get the exact path for your command.

> **Note:** Replace:
> - `<VERSION>` with the actual provider version (e.g. `6.4.0`)
> - `<OS>` with your operating system (e.g. `linux`, `windows`, `darwin`)
> - `<ARCH>` with your CPU architecture (e.g. `amd64`, `arm64`)

### Prerequisites

- Terraform must be installed and in your PATH
- Your Terraform files must include a Verity provider configuration
- Environment variables for authentication must be set (see "Required Environment Variables" section)
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
  mode = "datacenter" # Valid values: "datacenter" or "campus"
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

1. Get the latest Swagger JSON files from CAMPUS and DATACENTER systems.
2. Use `tools/process_swagger.py` script to remove unnecessary endpoints, merge both swagger files, and transform the result into a file ready for the Go SDK generator.

    ```bash
    python3 tools/process_swagger.py datacenter.json campus.json --output swagger_transformed.json
    ```

3. Install the OpenAPI Generator CLI (if not already installed):
  ```bash
  npm install @openapitools/openapi-generator-cli -g
  ```

4. Remove "openapi" folder completely from the project.

5. Generate the Go SDK using openapi-generator-cli:
  ```bash
  openapi-generator-cli generate -i swagger_transformed.json -g go -o ./openapi
  ```


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

Since API version 6.5, the provider supports two modes: **campus** and **datacenter**. Each mode has its own resource dependency ordering for creation and update operations:

**Order for CAMPUS:**
1. PB Routing ACL
2. PB Routing
3. Services
4. Eth Port Profiles
5. Authenticated Eth-Ports
6. Device Voice Settings
7. Packet Queues
8. Service Port Profiles
9. Voice-Port Profiles
10. Eth Port Settings
11. Device Settings
12. Lags
13. SFlow Collectors
14. Diagnostics Profiles
15. Diagnostics Port Profiles
16. Bundles
17. ACLs
18. IPv4 Lists
19. IPv6 Lists
20. Port ACLs
21. Badges
22. Switchpoints
23. Device Controllers
24. Sites

**Order for DATACENTER:**
1. Tenants
2. Gateways
3. Gateway Profiles
4. PB Routing ACL
5. PB Routing
6. Services
7. Packet Queues
8. Eth Port Profiles
9. Eth Port Settings
10. Device Settings
11. Lags
12. SFlow Collectors
13. Diagnostics Profiles
14. Diagnostics Port Profiles
15. Bundles
16. ACLs
17. IPv4 Prefix Lists
18. IPv6 Prefix Lists
19. IPv4 Lists
20. IPv6 Lists
21. PacketBroker
22. Port ACLs
23. Badges
24. Pods
25. Spine Planes
26. Switchpoints
27. Device Controllers
28. AS Path Access Lists
29. Community Lists
30. Extended Community Lists
31. Route Map Clauses
32. Route Maps
33. SFP Breakouts
34. Sites

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

# Local development
../tools/import_verity_state.sh
```

> **Tip:** For production, you can always locate the `tools` folder inside the provider directory (where the provider binary is installed). Use your file browser or terminal to navigate to the correct folder, then right-click the script and choose "Copy Path" to avoid manually typing the full path.

#### Windows PowerShell

```powershell
# Production
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1

# Local development
..\tools\import_verity_state.ps1
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

## Handling Auto-Assigned Fields

When you change an auto-assigned field's flag (such as `auto_assigned_vni`, `auto_assigned_vlan`, etc.) from `false` to `true`, you must remove the corresponding field (such as `vni`, `vlan`, etc.) from your Terraform resource block. Leaving the field present will cause issues, as the backend will automatically assign its value and may overwrite or ignore the value you specify in Terraform.

Our `data_source_state_importer` is designed to check if a field has a corresponding auto-assigned flag. If the flag is set to `true`, the importer will not write that field in the generated Terraform resource file — only the auto-assigned flag will be present. This ensures your configuration matches the backend's behavior and avoids conflicts.

**Best Practice:**
Whenever you enable auto-assignment for a field, always remove the manually specified value for that field from your `.tf` resource block.
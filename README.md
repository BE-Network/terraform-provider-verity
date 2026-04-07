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

The import process creates a special `stages.tf` file that defines explicit dependency ordering for resources. This uses the `verity_operation_stage` resource, which acts as an **active barrier** between resource type groups:

1. Establish a clear sequence for creating, updating, and destroying resources
2. Prevent dependency conflicts between resource types
3. Ensure that resources are processed in the optimal order for the Verity API
4. Wait for all operations from the current type group to complete before allowing the next group to start

Each imported resource is configured with the appropriate `depends_on` attribute referring to its corresponding stage. When a stage's `Create` is executed, it actively waits for its sibling resources to queue their operations, flushes them to the API, and only returns once all operations for that type group are complete. This guarantees sequential, ordered API execution regardless of Terraform's internal scheduling.

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
27. Device Settings
28. Diagnostics Port Profiles
29. Bundles
30. Pods
31. Badges
32. Spine Planes
33. Switchpoints
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

## Unit Tests

The provider includes a unit test suite that runs fully offline using a mock HTTP server — no real Verity API or Terraform state is required.

### Test packages

**`tests/unit/bulkops/`** — Tests for the bulk operations manager:
- Delete batching: large delete sets are split into batches of ≤100; each batch contains the correct resource names with none missing or duplicated across batches; a batch failure aborts remaining batches immediately; ACL header parameters handling
- Execution ordering: correct PUT/PATCH/DELETE sequencing for datacenter and campus modes, circular reference resolution, mixed operations, resource types with no queued operations generate no API calls; ACL v4 and v6 operations are dispatched as two separate PUT calls each carrying the correct `ip_version` query param; a PUT/PATCH/DELETE API failure stops all subsequent operations in the ordered sequence — resources scheduled after the failing type are never sent to the API, while those that already executed are unaffected

**`tests/unit/lifecycle/`** — Generic resource lifecycle tests run against every registered provider resource:
- Schema discovery: all resources expose a `name` attribute and discoverable fields/blocks
- PUT body completeness: all schema fields appear in the initial create request, including both the ref field and its `*_ref_type_` companion; integer `0`, bool `false`, and empty string values are present rather than silently omitted
- PUT body boundaries: a name-only create (only `name` provided in HCL) produces no unexpected extra fields beyond `name` and any auto-assigned flags
- PATCH correctness: enable field toggling, single string field updates, ref field pairs, nested block updates
- Nullable field transitions: explicit null vs omitted field handling
- Auto-assigned field exclusion: when a boolean `*_auto_assigned_` flag is set to `true`, the corresponding value field (e.g. `layer_3_vni`) is omitted from the PUT body — the backend assigns the value instead
- Mode field exclusion: datacenter-only fields absent in campus mode and vice versa
- Required query params: ACL `ip_version` param sent correctly for v4/v6
- Delete and import: resource removal and `terraform import` paths

### Running locally

```bash
# Export those env vars to speed up tests
export VERITY_DEFAULT_BATCH_DELAY=100ms
export VERITY_BATCH_COLLECTION_WINDOW=100ms
export VERITY_MAX_BATCH_DELAY=200ms
export VERITY_RESPONSE_PROCESSOR_DELAY=0s
export VERITY_DEBOUNCE_DELAY=100ms

go test ./tests/unit/lifecycle/ -count=1 -timeout 5m
go test ./tests/unit/bulkops/ -count=1 -timeout 2m
```

`-count=1` disables Go's test result cache, ensuring tests always execute rather than reusing a previous result.

### CI

Tests run automatically on every PR and push to `main` via `.github/workflows/test.yml`.
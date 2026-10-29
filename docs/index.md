# Verity Terraform Provider Documentation

## Supported Versions

Support is listed for specific provider and Verity API release pairs. A listed
pair does not imply support for newer API releases.

| Provider release | Supported Verity API release |
| --- | --- |
| `6.6.100` | `6.6.0.269` |

The runtime version check compares API major/minor, not the complete build
number. Passing that check does not establish support for an unlisted API build.
For breaking release upgrades, follow the
[fresh-import upgrade guide](guides/fresh-import-upgrades.md).

## Provider Configuration

The provider requires `uri`, `username`, `password`, and `mode`. Set them in the
provider block or through `TF_VAR_uri`, `TF_VAR_username`, `TF_VAR_password`, and
`TF_VAR_mode`. Nonempty provider attributes take precedence over environment
variables. Mode must be `datacenter` or `campus`; it has no default.

### Recommended: Using Environment Variables

Export the following environment variables before running Terraform:

For Linux/macOS:

```bash
export TF_VAR_uri="<your-verity-uri>"
export TF_VAR_username="<your-username>"
export TF_VAR_password="<your-password>"
export TF_VAR_mode="datacenter"
```

For Windows PowerShell:

```powershell
$env:TF_VAR_uri="<your-verity-uri>"
$env:TF_VAR_username="<your-username>"
$env:TF_VAR_password="<your-password>"
$env:TF_VAR_mode="datacenter"
```

With all four environment variables set, use:

```hcl
terraform {
  required_providers {
    verity = {
      source  = "BE-Network/verity"
      version = "6.6.100"
    }
  }
}

provider "verity" {}
```

Select the provider release that supports your Verity API release using the table
above. Set `TF_VAR_mode` to match your deployment.

Keep credentials out of committed Terraform configuration. See the
[provider configuration guide](../README.md#provider-configuration) for local
development overrides and additional examples.

### Parallelism Configuration

The provider groups queued operations by resource type. Higher parallelism can
let more resources queue together, while lower parallelism may split a type
across multiple API requests. Batching windows and the dependency graph also
affect grouping, so parallelism does not guarantee one request per type.

Example apply settings for Linux/macOS:

```bash
export TF_CLI_ARGS_apply="-parallelism=500"
```

For Windows PowerShell:

```powershell
$env:TF_CLI_ARGS_apply="-parallelism=500"
```

Adjust the value for your environment and retain the generated stage
dependencies. See [parallelism configuration](../README.md#parallelism-configuration).

### Reference Field Updates

An unset reference field and its `*_ref_type_` companion are omitted from unrelated
updates. For example, changing `vlan` does not clear an unmanaged
`policy_based_routing` reference. Configure the reference and its companion as
required by the resource schema when changing it; an explicit empty string can
clear a configured reference.

## State Importer

The state importer data source generates Terraform configuration and import
blocks from the current Verity API state. The `import_verity_state` scripts
automate generation and adoption; see [Tools](#tools).

```hcl
data "verity_state_importer" "import" {
  output_dir = "/path/to/directory"
}
```

If `output_dir` is omitted, generation uses the current working directory.
The scripts run two applies: the first generates configuration and import blocks,
and the second imports those objects and the generated operation stages into
Terraform state. Importing stages avoids their creation waits without changing
bulk batching settings or dependencies. Stages already in state are unchanged.

Only resources supported by the selected mode and API compatibility policy are
exported. Unsupported arguments, including nested arguments, are removed before
writing resource files and reported in a Terraform warning and
`unsupported_arguments.txt`.

Import blocks use identities recorded during generation, with IDs escaped as
literal HCL. `imported_files` lists only Terraform files written by the current
run, without duplicates. Existing user-authored files are not scanned for imports.

If a known resource output file exists but was not regenerated, generation fails
with its path and instructions to review and move or remove it. The scripts stop
before the import apply, preventing stale configuration from recreating deleted
objects. Files are not automatically deleted. User-authored files with other
names still participate in Terraform plans and need review before applying.

See the [state importer documentation](data-sources/verity_state_importer.md)
for its schema, warnings, and output behavior.

### Resource Dependency Management

The importer creates a mode-specific `stages.tf` dependency chain. API resources
use `depends_on` to refer to their corresponding `verity_operation_stage`.
Stage creation and deletion wait for and flush pending bulk operations; the
manager applies the operation orders from the registry.

When adding resources manually, reuse that chain and the appropriate
`depends_on`. See [resource dependency management](../README.md#resource-dependency-management)
and the [operation stage resource](resources/verity_operation_stage.md).

## Tools

Published provider archives include the import scripts in a `tools` folder next
to the binary. After `terraform init` installs the provider, this folder contains:

- `import_verity_state.sh`
- `import_verity_state.ps1`

Run the appropriate script from your Terraform project directory, which must
contain a configured Verity provider. Terraform must be in PATH; the shell script
also requires Bash and `awk`.

**Linux/macOS**:

```bash
.terraform/providers/registry.terraform.io/be-network/verity/<VERSION>/<OS>_<ARCH>/tools/import_verity_state.sh
```

**Windows PowerShell**:

```powershell
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1
```

Replace `<VERSION>` with the installed provider release, `<OS>` with your
operating system, and `<ARCH>` with your CPU architecture. The scripts remove
import blocks after a successful import and display any unsupported-argument
warning. For local development, see [running the scripts](../README.md#running-the-scripts).

## Handling Auto-Assigned Fields

When enabling a `<field>_auto_assigned_` flag, leave its value field unset in your
resource configuration. For example, with `vni_auto_assigned_ = true` on a
datacenter `verity_service`, omit the configured `vni` value.

The importer already omits the value when the API reports that the auto-assigned
flag is true. Keep that pattern when modifying generated configuration.

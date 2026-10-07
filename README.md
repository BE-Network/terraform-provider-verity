# Verity Terraform Provider

Terraform provider for Verity datacenter and campus deployments. See the
[provider documentation](docs/index.md) for resource examples and the
[state importer documentation](docs/data-sources/verity_state_importer.md) for
import output details.

## Table of Contents

- [Building the Provider](#building-the-provider)
- [Provider Configuration](#provider-configuration)
  - [Custom Provider Binary](#custom-provider-binary)
  - [Configuration and Environment Variables](#configuration-and-environment-variables)
  - [Parallelism Configuration](#parallelism-configuration)
- [Production Setup](#production-setup)
- [Regenerating the provider](#regenerating-the-provider)
  - [Preparing a new API release](#preparing-a-new-api-release)
  - [Breaking API changes and existing Terraform state](#breaking-api-changes-and-existing-terraform-state)
  - [Bulk bindings and scheduling](#bulk-bindings-and-scheduling)
  - [SDK generation and logging](#sdk-generation-and-logging)
  - [Adding or Changing a Resource](#adding-or-changing-a-resource)
- [Using the State Import Scripts](#using-the-state-import-scripts)
  - [Resource Dependency Management](#resource-dependency-management)
  - [What the Scripts Do](#what-the-scripts-do)
  - [Running the Scripts](#running-the-scripts)
  - [Prerequisites](#prerequisites)
- [Handling Auto-Assigned Fields](#handling-auto-assigned-fields)
- [Unit Tests](#unit-tests)
  - [Test packages](#test-packages)
  - [Compatibility evidence](#compatibility-evidence)
  - [Running locally](#running-locally)
  - [CI](#ci)
- [Security Toolchain](#security-toolchain)

## Building the Provider

Install the Go version required by [go.mod](go.mod), then build from the
repository root.

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

Build the provider, then create a `dev.tfrc` file with a development override.
The override must point to the directory containing the binary, on both platforms.
See [Terraform development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).

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
    "registry.terraform.io/local/verity" = "C:\\Users\\<user>\\terraform-provider-verity"
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

Set `TF_CLI_CONFIG_FILE` to the configuration file you created:

For macOS/Linux:

```bash
export TF_CLI_CONFIG_FILE=/path/to/dev.tfrc
```

For Windows PowerShell:

```powershell
$env:TF_CLI_CONFIG_FILE="C:\path\to\dev.tfrc"
```

Use `terraform plan` or `terraform apply` with the override. `terraform init`
still attempts to install a published provider, so it cannot download the
unpublished `local/verity` development provider.

### Configuration and Environment Variables

Set these values in the provider block or through the corresponding environment
variables. A nonempty provider attribute takes precedence over its environment
variable. All four values are required; mode has no default.

| Provider attribute | Environment variable | Value |
| --- | --- | --- |
| `uri` | `TF_VAR_uri` | Verity API base URL |
| `username` | `TF_VAR_username` | Verity API username |
| `password` | `TF_VAR_password` | Verity API password |
| `mode` | `TF_VAR_mode` | `datacenter` or `campus` |

For Linux and macOS:

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
provider "verity" {}
```

You can also set non-secret attributes directly, while supplying credentials
through environment variables:

```hcl
provider "verity" {
  uri  = "<your-verity-uri>"
  mode = "datacenter"
}
```

Keep credentials out of committed Terraform configuration.

### Parallelism Configuration

The provider groups queued operations by resource type. Higher Terraform
parallelism allows more resources to queue together and can reduce the number of
bulk API requests. Lower parallelism processes resources in waves, potentially
splitting a resource type across multiple requests; batching windows and the
dependency graph also affect how operations are grouped.

For larger applies, the existing `500` setting is a starting point. Adjust it to
your environment rather than treating it as a requirement or a guarantee that
all operations fit into one request.

For Linux and macOS:

```bash
export TF_CLI_ARGS_apply="-parallelism=500"
```

For Windows PowerShell:

```powershell
$env:TF_CLI_ARGS_apply="-parallelism=500"
```

Terraform's [apply parallelism](https://developer.hashicorp.com/terraform/cli/commands/apply#apply-options)
default is 10. Keep the stage dependencies described under
[Resource Dependency Management](#resource-dependency-management).

## Production Setup

For production use, configure Terraform to download the provider from the registry. You can find all available versions on the HashiCorp Registry:

**Provider Versions:** [View all releases on HashiCorp Registry](https://registry.terraform.io/providers/BE-Network/verity/latest)

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

## Regenerating the provider

The selected API version is `api_version` in `specs/overrides.yaml`. Generation
reads canonical inputs from `specs/openapi/<selected version>/`. The generated
registry carries that version into the provider's runtime compatibility check
and the mock server's default version response; no separate version constants
need updating.

Use the same command locally and in CI to verify all generated artifacts:

```bash
tools/generate_provider.sh --check
```

After updating and reviewing inputs and overrides, regenerate everything with:

```bash
tools/generate_provider.sh --write
```

The command verifies canonical formatting, checksums, and the selected input
version, then processes the SDK, extraction manifest, registry and embedded copy,
transport adapters, bulk bindings and operation orders, and resource
documentation. Write mode prepares all outputs in temporary directories before
copying them into the repository. Generators are
compiled before any output is replaced, so SDK type changes do not prevent the
adapter generator from running. Check mode leaves generated files unchanged.

Resource metadata lives in `specs/overrides.yaml`, which contains reviewed field
entries and deviations from defaults. Every extracted field requires an entry;
a conventional field can often declare just `api_name`. Lifecycle profiles supply
policies, modes and version ranges inherit from the resource, Terraform names
follow API names, and descriptions come from committed OpenAPI inputs. Explicit
values override defaults. `specs/generated_registry.json` and its embedded copy
in `internal/registry/registry.json` are generated together and checked in CI.

### Preparing a new API release

1. Update `api_version` in `specs/overrides.yaml` and review compatibility ranges.
   Ranges express supported releases and remain explicit policy.
2. Obtain both mode-specific exports, preserve the raw files outside the
   repository, and normalize them with the actual export date and provenance:

   ```bash
   task_api_version=$(go run ./tools/specgen version --overrides specs/overrides.yaml)
   go run ./tools/specgen normalize \
     --version "$task_api_version" \
     --datacenter /path/to/datacenter.json \
     --campus /path/to/campus.json \
     --output-dir "specs/openapi/$task_api_version" \
     --source-export-date YYYY-MM-DD \
     --provenance "Verity API export source"
   ```

3. Review endpoint and field changes and adjust override entries, policies,
   mode-specific metadata, and the reviewed bulk operation order. Classify
   removals, renames, type changes, and semantic changes using the compatibility
   process below before refreshing any fixtures.
4. Run `tools/generate_provider.sh --write`, review the diff, then run
   `tools/generate_provider.sh --check` and the tests. Update mock data and golden
   fixtures deliberately. Incompatible state changes require a reviewed fresh-import
   declaration and upgrade guide; an in-place migration is optional future work.
5. Validate against a lab system in both modes before releasing. Review README
   release examples and Dependabot branch targets when the supported release
   branches change. `.github/dependabot.yml` intentionally follows branch policy,
   independently of the selected API version. Historical tests and fixtures retain
   the version they describe.

### Breaking API changes and existing Terraform state

A new API release may add, remove, rename, or reinterpret endpoints and fields.
The supported workflow for a breaking release is to start with fresh Terraform
configuration and empty state, then run `import_verity_state` to adopt the current
server configuration. See [Fresh-import upgrades](docs/guides/fresh-import-upgrades.md).
An in-place upgrade of old state is not promised for a release declared this way;
production state upgraders are not required for that workflow.

Review the previous and new canonical API inputs in both modes. Generation
updates the SDK and metadata; maintainers still decide whether renamed endpoints
represent the same object, whether fields changed meaning or units, and what the
new provider can manage. Preserve Terraform names where practical. Removed
endpoints must be documented as no longer managed, rather than interpreted as
instructions to delete objects from the server. Review semantic changes even when
field names and types stay unchanged; the state-shape guard cannot detect them.

`specs/state_transitions.json` records reviewed fresh-import releases. It currently
contains no exceptions. Each release entry names:

- `api_version`, matching the selected version in the embedded registry;
- `baseline_sha256`, matching the pinned historical baseline;
- `strategy: "fresh_import"`;
- `upgrade_guide`, an existing, nonempty Markdown file under `docs/guides/` that
  explains the release changes and links to the fresh-import procedure;
- `changes`, with each affected `terraform_type`, a `reason`, and
  `expected_state_sha256` for its complete reviewed target representation, or
  `"removed"` for a retired resource.

After regeneration, print a declaration proposal without changing files:

```bash
go test ./tests/unit/lifecycle/ -run '^TestStateCompatibilityReport$' -count=1 -v
```

Review the reported changes, write the release guide, fill in the guide path and
reasons, and add the entry to `releases` in `specs/state_transitions.json`. Then run
the guard, importer and lifecycle tests, and only then update the ordinary schema
snapshot. The report grants no exceptions and does not edit the policy or baseline.

Declarations are scoped to the selected API version and exact resource targets.
A missing declaration, an unlisted resource, or a later target-shape change still
fails. A blanket release exemption or environment-variable bypass is not supported.
The pinned baseline remains unchanged. Preserve a new historical baseline when
publishing later releases and extend coverage to it so newly introduced resources
and fields are protected in subsequent upgrades.

Fresh import adopts server configuration; unapplied changes from the old
Terraform configuration do not carry over automatically.

The API version and each resource's schema version are independent. Fresh import
creates state from the current schema; incrementing schema versions or implementing
old-state conversions is not a prerequisite for a declared fresh-import release.
Supporting in-place upgrades later would require production state-upgrade or
resource-state-move integration and separate migration tests.

### Bulk bindings and scheduling

`tools/specgen bulk` generates `internal/bulkops/generated_registry.go` from the
registry and SDK declarations. The unified generation command invokes it for both
write and check modes. SDK services, request types, body setters, delete query
parameters, typed-map request preparers, response collection keys, cache keys,
and fixed split parameters are discovered and checked; missing or ambiguous
bindings fail generation. Pre-existence checks use the generated GET binding
and collection key.

Each resource declares `bulk_order` in `specs/overrides.yaml`:

```yaml
bulk_order:
  datacenter: {put: 110, patch: 120, delete: 330}
```

Positive ranks run in ascending order within their operation and mode. Leave an
unsupported operation absent. Gaps between ranks allow inserting a resource
without renumbering its neighbors. Every supported mutation needs a rank; ranks
cannot collide across bulk keys, and variants sharing a key must agree on order.
PATCH order is independent of PUT order. Validation requires DELETE order to
reverse PUT order in each mode, and `sfp_breakout` has only a PATCH rank that
runs first in both modes. ACL variants share `acl` and retain their `ip_version` query split
and response collection selection. Review ordering metadata when adding or
changing an endpoint; conventional bindings require no handwritten Go callback
or operation-list edits. Bulk orders are separate from importer stage orders.

### SDK generation and logging

OpenAPI Generator `v7.25.0` is pinned by Docker image digest. The generation
command requires Docker, Python 3, Go, and `rsync` for `--write`.
`tools/generate_openapi_sdk.sh --check` and `--write` remain available for SDK-only
work and use the same selected version and preparation pipeline. The complete
provider command also updates the dependent adapters and metadata.

Generation retains only Go SDK sources, excludes generated tests and module
files, and formats output. `tools/sdkgen/call_api.go.tmpl` supplies `callAPI` auth
body, authorization/session header, and URL user-information redaction. Actual
traffic remains intact. Ordinary resource bodies remain logged in full and may
contain resource secrets. Edit the customization instead of restoring `client.go` by hand.
See [SDK preparation](tools/sdkgen/README.md) for the retention policy.

CI runs the complete generated-artifact check and debug-log regression tests on
every PR and push covered by the test workflow, using Docker and Python 3 on the
Ubuntu runner.

### Adding or Changing a Resource

Every API-backed resource is served by one generic engine
(`internal/genericresource`) from the resource registry; there are no
per-resource lifecycle files. `verity_operation_stage` is the one bespoke
resource.

To add a resource:

1. **OpenAPI input.** Make sure the endpoint is in the committed inputs under
   `specs/openapi/<selected version>/`, and run the generated-artifact check above.
2. **Override entry.** Add the resource to `specs/overrides.yaml`: `path`,
   `terraform_type`, `description`, `modes`, `identity_path`, the `api` keys
   (`bulk_key`, `cache_key`, `response_collection_key`, `delete_parameter`),
   and its `fields`. Declare only what OpenAPI cannot express:
   - `import_stage`, the state importer's stage name and order for each mode
     the resource supports;
   - `bulk_order`, independent PUT/PATCH/DELETE ranks for each supported mode;
     these orders express API dependencies and differ from import stages;
   - `auto_assignment.recomputed_when`, if the server reassigns an
     auto-assigned value when another field changes (see `vni` under
     `verity_service`);
   - a nested list's `collection`;
   - `modes` on fields that exist in only one mode;
   - a policy that differs from the field's profile: `create_null`,
     `update_clear`, `unknown_plan`, `response_absence`, `state_ownership`.

   Reference and auto-assignment pairs need no declaration. A field with a
   `<field>_ref_type_` companion becomes a reference pair whose allowed types
   are the companion's enum. A field with a `<field>_auto_assigned_` companion
   becomes an auto-assignment pair.
3. **Regenerate** with `tools/generate_provider.sh --write`, which updates the
   SDK, manifest, registry and embedded copy, adapters, bulk bindings/orders, and
   docs together. The adapter generator prints why it cannot serve an unsupported
   resource.
4. **Tests.** Coverage tests read the resource path, wrapper, mode, fixed
   headers, and supported operations from the registry. Add a mock response in
   `tests/unit/testdata/responses/<mode>/`, record the golden fixtures and the
   schema snapshot, review the diff, then run the suites (see [Unit Tests](#unit-tests)):

   ```bash
   UPDATE_GOLDEN=1 go test ./tests/unit/lifecycle/ -run TestGoldenWireFixtures -count=1
   UPDATE_SCHEMA_SNAPSHOT=1 go test ./tests/unit/lifecycle/ -run TestSchemaGolden -count=1
   ```

To add or change a field on an existing resource, update its override entry and
regenerate as in step 3. Defaults supply conventional shapes and policies, but
every extracted field still requires an entry. Review
the regenerated registry, adapter, and docs, then update the golden fixtures
and schema snapshot as in step 4.

The pages in `docs/resources` are generated from one template,
`tools/specgen/templates/resource.md.tmpl`. Change the template or the
registry, never a page. `verity_operation_stage.md` is the one handwritten page.
In CI, `tools/generate_provider.sh --check` verifies every generated artifact,
so a registry, adapter, SDK file, or doc page that differs from its inputs fails
the build.

## Using the State Import Scripts

The provider includes scripts to help import existing Verity resources into Terraform state. These scripts automate the process of creating resource files and importing existing resources.

### Resource Dependency Management

The importer generates `stages.tf` with a dependency chain of
`verity_operation_stage` resources. Each API resource depends on its corresponding
stage. Stage creation and deletion wait for and flush pending bulk operations;
the manager applies the PUT/PATCH/DELETE orders from the registry.

#### Creating New Resources

When adding resources manually, reuse the generated stages and the appropriate
`depends_on`. A minimal example for one resource type is:

```hcl
resource "verity_operation_stage" "service_stage" {}

resource "verity_service" "example" {
  name       = "example-service"
  depends_on = [verity_operation_stage.service_stage]
}
```

For multiple resource types, retain the mode-specific chain in the generated
`stages.tf`. See [operation stages](docs/resources/verity_operation_stage.md).

### What the Scripts Do

1. Find the main Terraform file with the Verity provider
2. Add the `verity_state_importer` data source if it doesn't exist
3. Run a first `terraform apply` to generate resource files and import blocks
4. Run a second `terraform apply` to import the resources into your state
5. Remove the generated import blocks after a successful import and clean up
   temporary files

Before writing resource files, the importer removes arguments unsupported by
the provider schema, including nested arguments. It reports them in a Terraform
warning and `unsupported_arguments.txt`; the scripts print that file after the
import. Resource names and values are escaped as literal HCL, names have a
deterministic natural order, and ACL v4/v6 use separate files.

Import blocks use the resource addresses and API names recorded while generating
configuration. IDs are escaped as literal HCL, and blocks retain importer stage
order. Generated operation stages also receive import blocks with ID `"stage"`,
so importing them avoids the waits in stage creation. Normal batching settings
and stage dependencies are unchanged. Stages already in state are not imported
again. The `imported_files` result lists only Terraform files written by the
current run, including `stages.tf` and `import_blocks.tf` once each. Existing
user-authored or stale files are not scanned for imports or included in that
result. If a known resource output file exists but was not regenerated, generation
fails with its path and instructions to review and move or remove it. This also
covers resource files excluded by the current mode or API compatibility policy.
The scripts stop before the import apply, preventing stale generated configuration
from recreating deleted objects. Files are not automatically deleted.

User-authored files with other names remain on disk and still participate in
Terraform plans. They no longer receive automatically generated import blocks;
review them before applying the directory.

Configuration generation stops at the first fetch, render, or write error. The
scripts stop when that first apply fails and do not run the import apply with an
incomplete set of generated resources. Files written for earlier resource types
may remain in the output directory; correct the reported error and rerun generation
before importing.

### Running the Scripts

#### Linux and macOS

```bash
# Production
.terraform/providers/registry.terraform.io/be-network/verity/<VERSION>/<OS>_<ARCH>/tools/import_verity_state.sh

# Local development
../tools/import_verity_state.sh
```

#### Windows PowerShell

```powershell
# Production
.terraform\providers\registry.terraform.io\be-network\verity\<VERSION>\<OS>_<ARCH>\tools\import_verity_state.ps1

# Local development
..\tools\import_verity_state.ps1
```

> **Note:** Replace:
> - `<VERSION>` with the actual provider version (e.g. `6.6.0`)
> - `<OS>` with your operating system (e.g. `linux`, `windows`, `darwin`)
> - `<ARCH>` with your CPU architecture (e.g. `amd64`, `arm64`)

### Prerequisites

- Terraform with configuration-driven import blocks must be installed and in PATH.
- Run the scripts from the Terraform working directory after `terraform init`
  for a registry-installed provider.
- A `.tf` file must contain a configured Verity provider, with credentials and
  mode supplied as described in [Configuration and Environment Variables](#configuration-and-environment-variables).
- The shell script requires Bash and `awk`; Windows uses PowerShell.

## Handling Auto-Assigned Fields

When enabling a `<field>_auto_assigned_` flag, leave the corresponding value
unset in your resource configuration. For example, with
`verity_service.vni_auto_assigned_ = true`, omit the configured `vni` value.

The importer already omits the value when the API reports that its auto-assigned
flag is true. Keep that pattern when modifying generated configuration.

## Unit Tests

Tests use mock HTTP servers and do not require a live Verity system. Lifecycle
tests run the Terraform CLI with temporary configuration and state; install
Terraform before running them.

### Test packages

| Package | Coverage |
| --- | --- |
| `internal/...` | Provider contracts, transport conversion, importer rendering and stale-file checks, bulk bindings and ordering, registry validation, and utility behavior |
| `tools/...` | Spec extraction, generators, SDK preparation, and reproducibility |
| `tests/unit/lifecycle/` | Create/read/update/delete/import, null and auto-assigned fields, references, mode exclusions, request goldens, importer output, and state compatibility |
| `tests/unit/bulkops/` | DELETE batching, operation ordering, ACL query parameters, and failure handling |
| `tests/unit/sdk/` | Auth-body and session-header log redaction without changing real traffic |
| `tests/unit/utils/` | HCL configuration parsing |

### Compatibility evidence

`internal/provider` checks that registered provider schemas, modes, and fields
match the current registry contract. That is an integration check of current
code, rather than a comparison with a retired provider implementation.

Mode-exclusion expectations use `specs/generated_manifest.json`, extracted from
canonical datacenter/campus OpenAPI inputs before overrides are applied. The test
lookup rejects unknown endpoints, nested field paths, and modes. Schema-derived
field paths must exist in that independent source; a typo cannot silently skip
an assertion. Any future Terraform field alias needs an explicit mapping to its
API field for these checks.

Independent baselines include the lifecycle schema snapshot, recorded legacy
requests, historical field-policy and cache-key fixtures, bulk wire/order/cache-key
goldens, and importer output goldens. The `generic_*_differential_test.go` tests
compare current behavior with recorded requests; they do not run a retired provider.
Review baseline changes explicitly. Historical legacy-request fixtures have no
automatic update path, and updating a schema snapshot does not provide a migration.

`TestStateCompatibility` separately checks existing field names, types, collection
element types, and block nesting against
`tests/unit/testdata/state_compatibility_v6_6.json`. This baseline was projected
from the checked-in 6.6 schema snapshot at commit
`71b2718753bd1297595afc34636d1d78fa718448`; its provenance and checksum are pinned.
The maintainer review verified that all 51 state representations and schema
versions also match released provider `v6.6.100`.
It excludes descriptions, validators, plan modifiers, and other schema metadata.
New resources and added fields are allowed when existing representations remain
compatible. Removing a resource or field, renaming a field, changing a type or
nesting, or changing a schema version fails unless it matches an explicit reviewed
fresh-import declaration for the selected API version.

`TestSchemaGolden` runs this guard before writing, including when
`UPDATE_SCHEMA_SNAPSHOT=1`. Updating the ordinary snapshot cannot replace the
historical state baseline or bypass compatibility checks. Run the fast guard with:

```bash
go test ./tests/unit/lifecycle/ -run '^TestStateCompatibility' -count=1
```

The guard detects undeclared breaks rather than requiring automatic migration.
Approved fresh-import declarations allow resource retirement and incompatible
field changes without a production upgrader. They must match the complete target
representation, including newly added fields and the schema version. A subsequent
unreviewed change invalidates that declaration. The declaration loader rejects
missing or empty upgrade guides, malformed policies, and unknown resource names.
Future support for in-place migrations remains separate from this release policy.

### Running locally

```bash
# Use the CI delay settings for local tests
export VERITY_DEFAULT_BATCH_DELAY=100ms
export VERITY_BATCH_COLLECTION_WINDOW=100ms
export VERITY_MAX_BATCH_DELAY=200ms
export VERITY_RESPONSE_PROCESSOR_DELAY=0s
export VERITY_POST_OPERATION_VERIFICATION_BACKOFF=0s
export VERITY_DEBOUNCE_DELAY=100ms

go test ./internal/... ./tools/... -count=1 -timeout 5m
go test ./tests/unit/lifecycle/ -count=1 -timeout 15m
go test ./tests/unit/bulkops/ -count=1 -timeout 2m
go test ./tests/unit/sdk/ -count=1 -timeout 2m
go test ./tests/unit/utils/ -count=1
```

`-count=1` disables Go's test result cache, ensuring tests always execute rather than reusing a previous result.

### CI

Tests and generated-artifact checks run automatically for PRs and pushes targeting
`main`, `release/**`, and `dev/**` through `.github/workflows/test.yml`.

## Security Toolchain

Security checks are defined in
[security-baseline.yml](.github/workflows/security-baseline.yml),
[codeql.yml](.github/workflows/codeql.yml), and
[scorecards.yml](.github/workflows/scorecards.yml).

The shared standard is documented in the
[Security Toolchain Standard](https://github.com/BE-Network/verity-monitoring/blob/main/SECURITY_TOOLCHAIN_STANDARD.md).

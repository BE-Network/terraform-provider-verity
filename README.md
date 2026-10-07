# Local Verity Terraform Provider Setup

## Security Toolchain

This repository uses a standardized security baseline implemented via GitHub Actions workflows:

- `.github/workflows/security-baseline.yml`
- `.github/workflows/codeql.yml`

Central standard:

- https://github.com/BE-Network/verity-monitoring/blob/main/SECURITY_TOOLCHAIN_STANDARD.md

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

For each target release, update overrides, API policies, importer coverage, mock
responses and fixtures for added, removed, renamed and retyped fields/endpoints.
Exercise fresh import from an empty state against target API mocks in both modes,
check the resulting plan and HTTP behavior, then validate on a lab system. Document
server objects the new API no longer exposes, and any configuration that must be
supplied manually. This workflow adopts server configuration, so unapplied changes
from the old Terraform configuration do not carry over automatically.

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
contain resource secrets; registry-based sensitive-field redaction is a separate
follow-up. Edit the customization instead of restoring `client.go` by hand.
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
   schema snapshot, review the diff, then run the suites (see "Unit Tests"):

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
> - `<VERSION>` with the actual provider version (e.g. `6.6.0`)
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
- Execution ordering: correct PUT/PATCH/DELETE sequencing for datacenter and campus modes, mixed operations, resource types with no queued operations generate no API calls; ACL v4 and v6 operations are dispatched as two separate PUT calls each carrying the correct `ip_version` query param; a PUT/PATCH/DELETE API failure stops all subsequent operations in the ordered sequence — resources scheduled after the failing type are never sent to the API, while those that already executed are unaffected

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
- Recorded legacy comparisons (`generic_*_differential_test.go`): `runLifecycle`
  executes the current provider, while `legacyReference` reads the retired
  handwritten provider's recorded requests from `testdata/legacy_reference`.
  These tests do not run a second implementation. Historical test names stay
  stable because they identify the fixture filenames. Intended semantic changes
  are asserted explicitly against the recorded behavior.

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

Independent evidence remains checked in: the lifecycle schema snapshot, recorded
legacy requests, `internal/provider/testdata/legacy_field_policies.json`, the
historical resource-to-cache-key fixture, bulk wire/order/cache-key goldens, and
importer output goldens. The cache-key fixture preserves the former handwritten
test mapping; it is not generated from the current registry. Review changes to
these baselines explicitly. The schema snapshot update command above refreshes
that snapshot; it does not implement a state migration or prove that one is
unnecessary. Historical legacy-request fixtures have no automatic update path.

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

Release requirements and generation instructions live in tracked documentation.
Private `refactor/` notes are not required inputs, CI checks, or release gates.

### Running locally

```bash
# Export those env vars to speed up tests
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
```

`-count=1` disables Go's test result cache, ensuring tests always execute rather than reusing a previous result.

### CI

Tests and generated-artifact checks run automatically for PRs and pushes targeting
`main`, `release/**`, and `dev/**` through `.github/workflows/test.yml`.

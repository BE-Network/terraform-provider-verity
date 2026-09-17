# Refactoring plan: schema-driven Verity resources

- Status: Phases 0, 1, and 2 closed; Phase 3 in progress (16 opt-in generic resources, singleton objects included)
- Prepared: 2026-09-08
- Scope: API-backed Terraform resources in `internal/provider`, their field handling, resource registration, bulk-operation metadata, and schema/OpenAPI tooling.

## Executive recommendation

Build one generic resource lifecycle engine around an enriched, declarative `ResourceSpec`. Generate most of each spec from the datacenter and campus OpenAPI documents, keep a small reviewed override file for behavior that OpenAPI cannot express, and embed the resulting registry in the provider binary.

The desired end state is close to “add a schema and the resource works,” but a plain Terraform or OpenAPI schema is not enough. It knows field types and nesting, but it does not completely describe:

- Terraform names, required/optional/computed behavior, replacement rules, and state compatibility;
- the API path, request wrapper, response collection key, delete parameter, cache key, and mode;
- whether an array is a value, an indexed child collection, or a singleton object exposed as a block;
- how an indexed child is added, updated, deleted, reordered, or assigned an index by the server;
- paired value and `*_ref_type_` behavior;
- paired value and `*_auto_assigned_` behavior;
- nullable-field reset behavior and the difference between API omission and API `null`;
- update-only resources, endpoint headers such as ACL `ip_version`, response filtering, operation ordering, and circular dependency workarounds.

Therefore, the source of truth should be:

```text
datacenter OpenAPI + campus OpenAPI
                  |
                  v
       deterministic spec generator <--- small explicit overrides
                  |
                  v
       validated, embedded ResourceSpec registry
          |             |              |
          v             v              v
 Terraform schema   generic CRUD   generated docs/tests
                         |
                         v
              existing bulk operation manager
```

This should be an incremental refactor. Do not regenerate all 41 current API-backed resource implementations in one change, and do not rewrite the bulk-operation scheduler at the same time as the field codec.

Phase 0 can begin immediately. Work on the production generic lifecycle engine must wait until all of these gates pass:

1. versioned, mode-specific OpenAPI inputs are committed and reproducible;
2. API-version applicability and full field lifecycle policies are represented and validated in the spec;
3. a Terraform Framework spike proves generic plan/config reads and state writes without losing null or unknown values;
4. a generated transport adapter proves that canonical codec output can traverse the current typed OpenAPI/bulk boundary with exact wire JSON;
5. resource registration is fixed before configuration and mode-independent; and
6. the version-zero state contract and future state-upgrader mechanism are tested.

## Repository assessment

### What exists today

- There are 41 `resource_verity_*.go` files totaling about 33,070 lines. Every one implements `Schema`, `Create`, `Read`, `Update`, `Delete`, `ImportState`, and `ModifyPlan` separately.
- Those files contain 72 `ListNestedBlock` declarations; 35 resources have at least one nested list block and 28 model `object_properties` as a list block.
- All resources already use common operation, fetch/retry, state-mapping, and mode-nullification helpers. This is a useful base for migration.
- Twenty-six resources use `ProcessIndexedArrayUpdates`, 22 parse HCL files for nullable-field presence, 24 contain reference-type fields, and four contain auto-assignment fields.
- `internal/bulkops/registry.go` is already a transport registry for API request types and endpoint functions, but it repeats resource metadata also found in provider constructors, compatibility maps, importer maps, cache keys, and resource files.
- `tools/process_swagger.py` produces the Go SDK input, while `tools/compare_schemas.py` separately generates `internal/utils/schema.go` for mode-specific fields. Schema generation is therefore already partially adopted, but only for mode metadata.
- The lifecycle test suite introspects resource schemas and covers common PUT, PATCH, delete, import, reference, nullable, auto-assigned, nested-block, and mode cases. It can become the parity harness for the migration.

### Main maintainability problems

1. **A field is declared in too many places.** A typical field appears in a Go model, Terraform schema, create mapping, update comparison, read mapping, plan nullifier, generated mode map, documentation, and tests. Nested and nullable fields require still more declarations. Missing one location can create a silent state or PATCH bug.

2. **Resource identity has several aliases.** For example, one resource can have a Terraform type, endpoint/plural name, bulk-operation key, request wrapper, response key, delete parameter, cache key, and importer key. These conventions are not fully consistent, so they are maintained in multiple maps and string literals.

3. **Common behavior is helper-based rather than data-driven.** Helpers reduce individual statements but every resource still assembles lists of string, bool, integer, number, nullable, nested, and mode-aware fields. Adding a field remains a manual lifecycle edit.

4. **Nested collection behavior is repeated in closures.** Each indexed block supplies create, compare/update, and delete-marker callbacks even though the dominant API protocol is the same: match by `index`, send changed/new elements, and represent removal using an index-only element.

5. **Special field semantics are inconsistently centralized.** Generic reference-pair helpers are used directly by only a small part of the resource set; other references are handled as ordinary strings or by resource-specific code. Auto-assigned fields have repeated validation and payload-suppression logic.

6. **The nullable-field solution depends on Terraform source files.** `ParseResourceConfiguredAttributes` scans `*.tf` in the provider working directory to distinguish omission from explicit `null`. That cannot reliably represent all valid Terraform configurations, including child modules, `.tf.json`, expressions whose value is null, generated configurations, renamed resource labels, or clients that do not expose local source files to the provider. Filesystem parsing should not be part of the final lifecycle engine. *(Superseded: the scan is kept as the permanent nullable contract; see the decision under section 6, which also corrects this list — a null from an expression is handled, because the scan records the attribute whatever its value.)*

7. **Mode metadata fails open.** `FieldAppliesToMode` treats unknown resources and fields as available in both modes. This avoids breakage but lets a newly added field bypass mode review. A generated spec should validate every field and fail generation/CI when mode is unresolved.

8. **Registration and ordering are fragile to extend.** Adding a resource currently affects the provider constructor list, compatibility map, bulk registry, importer mappings, mock-test table, and usually hard-coded operation ordering. `FilterResourcesByMode` also recovers Terraform names by matching Go concrete type names rather than consuming canonical metadata.

9. **Singleton objects are represented as unrestricted lists.** Code generally reads only `object_properties[0]`. The schema should express cardinality explicitly and reject or eliminate extra elements.

10. **OpenAPI updates are not end-to-end.** The README currently instructs maintainers to regenerate the SDK and then manually add or remove fields in provider resource files. This is exactly the high-risk work the new generator should remove.

## Refactoring options considered

| Option | What it improves | What remains | Recommendation |
| --- | --- | --- | --- |
| Shared lifecycle shell with per-resource callbacks | Removes repeated authentication, cache, operation submission, refresh, delete, and import code | Models plus create/update/read field lists and nested handlers remain manual | Useful as a temporary extraction if needed, but not the end state |
| Declarative field tables with typed resource models | Centralizes scalar mapping, mode, references, and nullability while retaining compile-time models | Every field still needs model declarations and generated-SDK binding code | Viable intermediate step toward the generic codec |
| Generate complete resource Go files from templates | Makes new resources faster to scaffold | Produces large generated files, duplicates behavior, and makes engine fixes regenerate the whole provider | Better than manual copying, but inferior to a runtime generic engine |
| Embedded `ResourceSpec` plus generic lifecycle engine | Makes schema, mapping, diffing, nesting, and normal CRUD data-driven | Requires a carefully designed value codec and explicit exception policies | **Recommended target** |
| Fetch schema dynamically from the configured server | Appears to eliminate release-time regeneration | Conflicts with Terraform's stable, pre-configuration schema contract and makes behavior server-dependent | Reject as a provider runtime design |

The first two options can reduce risk during extraction, but every intermediate change should move metadata into the eventual `ResourceSpec` rather than create another permanent registry.

## Target design

### 1. A canonical `ResourceSpec`

Introduce a provider-owned intermediate representation. It should be independent of the generated OpenAPI Go struct names so regenerated SDK naming changes do not leak into lifecycle code.

Conceptually:

```go
type ResourceSpec struct {
    TerraformType string
    Description   string
    Modes         []Mode
    Versions      VersionRange
    IdentityPath  string
    SchemaVersion int64

    API APIResourceSpec
    Operations OperationSpec
    Fields []FieldSpec
    Dependencies DependencySpec

    Hooks ResourceHooks // empty for the normal case
}

type FieldSpec struct {
    TerraformName string
    APIName       string
    Kind          FieldKind // string, bool, int64, number, object, list
    Access        Access    // required, optional, computed, optional+computed
    Description   string
    Nullable      bool
    Sensitive     bool
    Replace       bool
    Modes         []Mode
    Versions      VersionRange
    Default       *LiteralSpec
    Validators    []ValidatorSpec

    ResponseAbsence ResponseAbsencePolicy
    CreateNull      CreateNullPolicy
    UpdateClear     UpdateClearPolicy
    UnknownPlan     UnknownPlanPolicy
    StateOwnership  StateOwnershipPolicy

    Fields         []FieldSpec
    Collection     *CollectionSpec
    Reference      *ReferenceSpec
    AutoAssignment *AutoAssignmentSpec
}

type LiteralSpec struct {
    Kind  LiteralKind // null, string, bool, int64, decimal
    Value string      // canonical lexical value; empty for null
}
```

Use a provider-owned API version type such as `{Major, Minor int}` and an inclusive-minimum/exclusive-maximum `VersionRange`. Do not compare versions as floats or strings. Every resource and field must have an explicit range; the generator should reject an unresolved or contradictory range. Nested fields inherit their parent's range only when the generated manifest records that inheritance explicitly.

The current provider binary targets API 6.6 exactly. Its compiled registry and Terraform schema must therefore be generated for 6.6 even if ignored 6.5 source files happen to exist in a maintainer's checkout. A later API-version provider build selects its matching registry at generation/build time. Version metadata is not permission to change Terraform schemas dynamically after provider configuration. If a future binary intentionally supports several API minors, its public Terraform schema must be the stable union and version-specific lifecycle behavior must be validated separately.

`LiteralSpec` replaces `any` defaults so generated output is deterministic and serializable. Strings remain exact strings, booleans use `true`/`false`, integers use canonical base-10 text, decimals retain canonical decimal text rather than passing through `float64`, and null has its own kind. Spec validation must type-check every literal against its field.

The five field policies make request/state ownership explicit:

- `ResponseAbsence`: set Terraform null, apply the declared default, preserve prior state, or raise a diagnostic;
- `CreateNull`: omit, send API null, send the declared default/literal, or reject;
- `UpdateClear`: send API null, empty string, zero, false, declared default/literal, omit as unmanaged, or reject;
- `UnknownPlan`: omit and resolve by Read, preserve prior state, or reject when the operation requires a known value;
- `StateOwnership`: configuration-owned, server-owned/computed, or configuration-or-server (`Optional+Computed`).

Unknown values are never serialized accidentally. Any policy that substitutes state/default for unknown must say so explicitly. The spec validator should also reject incompatible combinations, such as configuration-owned fields with `ResponseAbsence=preserve` but no prior-state context, or computed-only fields with a user-driven update-clear policy.

`APIResourceSpec` should contain all current aliases in one place: endpoint path, bulk key, request wrapper key, response collection key, delete query parameter, cache key, supported methods, and fixed headers/query parameters. All provider registration, importer mapping, mock setup, and bulk registry generation should consume this same object.

`CollectionSpec` should explicitly select a strategy rather than infer all arrays as equivalent:

- `singleton`: one API object, exposed as a single nested block or a list constrained to one during compatibility migration;
- `indexed_patch`: identity is `index`; create/update/delete use the existing partial-array protocol;
- `indexed_server_assigned`: zero or missing index means create and must not be collapsed with other new items;
- `replace`: send the complete array when anything changes;
- `computed_subset`: API returns hardware/global entries, but state preserves only configured identities;
- `ordered` or `unordered`: controls diff and state normalization.

The spec validator must reject an object/list without an explicit cardinality and update strategy.

### 2. OpenAPI extraction plus reviewed overrides

#### Reproducible input decision

Commit the canonical mode-specific generator inputs to this repository. CI must not download schemas from a live Verity system or depend on ignored files in a maintainer's checkout.

The current repository tracks `openapi/api/openapi.yaml`, which is the merged/transformed SDK input, while files named `dc6*.json`, `campus6*.json`, and transformed variants are ignored. The tracked merged file is not sufficient for resource-spec generation because mode provenance has already been lost. Adopt this layout:

```text
specs/openapi/
  6.6/
    datacenter.json
    campus.json
    manifest.json
```

`manifest.json` records the API version, mode, canonical file SHA-256, source/export date, normalization-tool version, and licensing/provenance note. The committed JSON should be a deterministic normalization of the exported source: sorted object keys, stable formatting, and removal only of an explicitly documented allowlist of volatile metadata. The normalization command must produce identical bytes on repeated runs.

For the 6.6 provider, only the committed 6.6 pair is authoritative. Local exports for other API versions are development inputs until a corresponding source set and manifest are deliberately committed on the appropriate provider version line. `openapi/api/openapi.yaml`, the Go SDK, the embedded resource registry, mode/version reports, and generated documentation become derived outputs of those committed inputs plus overrides.

CI runs fully offline for this pipeline:

1. verify manifest checksums;
2. regenerate the merged SDK input and resource manifest into a temporary directory;
3. regenerate code/docs as applicable; and
4. fail on any diff.

Acquiring new exports remains a maintainer/release action outside CI. This keeps builds reproducible even when the source Verity systems are private or unavailable.

Generate, rather than hand-author, everything the source API documents describe reliably:

- endpoint methods and parameters;
- API field name and primitive/object/array type;
- nested structure;
- description, enum, minimum/maximum, length constraints, format, and default;
- nullability as declared by the source API schema;
- datacenter/campus presence;
- first/last API version in which each resource and field exists;
- likely reference pairs from `<field>` plus `<field>_ref_type_`, including the allowed reference types from the enum;
- likely auto-assignment pairs from `<field>` plus `<field>_auto_assigned_`;
- likely indexed arrays from an `index` child field.

Keep ambiguous Terraform and Verity behavior in a small override document. Overrides must be explicit and validated, not arbitrary Go snippets. Expected examples include:

- canonical Terraform type and exceptional wrapper/delete names;
- `name` as identity and `RequiresReplace`;
- optional/computed policy when an API default exists;
- update-only `site` and `sfp_breakout` operations;
- ACL v4/v6 variants and `ip_version` parameters;
- hardware-response filtering;
- singleton versus repeated blocks;
- collection patch strategy and server-assigned indices;
- dependency ordering and the route-map-clause/tenant circular-reference strategy;
- rare field-specific normalization or diff suppression.

The generator should ingest the committed pair for the selected API version directly. When more than one version is committed, it should compare adjacent manifests to derive candidate version ranges and require review for removals, reintroductions, or type changes. Do not derive resource behavior from only the merged/transformed SDK input, because merging loses source-mode provenance and the current transform intentionally changes numeric nullability for SDK generation.

Generated files should contain a header and be reproducible. CI should run the generator in `--check` mode and fail when regeneration creates a diff.

### 3. One recursive field codec

#### Framework value-bridge gate

Before implementing the codec, build a test-only spike that answers whether the Terraform Plugin Framework can support the proposed generic value path safely. Do not migrate a production resource as part of this spike.

The spike schema should contain:

- required and `Optional+Computed` string/bool/int64/number attributes;
- a nullable value and a computed unknown value;
- a singleton nested object represented in the current state shape;
- an indexed list of nested objects; and
- at least two nesting levels.

Exercise both whole-value and path-based APIs (`Get`, `GetAttribute`, `Set`, and `SetAttribute`) using `attr.Value`, `types.Object`, and `types.List` rather than a concrete resource model. Run it through the protocol test server, not only direct unit calls. The test matrix must cover known, null, and unknown values in config/plan/state, diagnostics from partially unknown collections, exact numeric round trips, and state written from decoded API data.

The gate passes only if:

- no concrete per-resource Go struct is required;
- null and unknown remain distinguishable at every supported nesting level;
- state can be written using a schema-derived type without conversion diagnostics;
- list element types and object attribute types are deterministic from `FieldSpec`;
- prior state can be read for diff/preserve policies without flattening unknown values; and
- import and refresh can construct the same state type from only identity plus API JSON.

If the Framework cannot safely support the recursive runtime tree, stop and choose the intermediate design from the options table: generate thin typed model bindings while keeping schema, field policies, and lifecycle traversal declarative. Do not compensate with unchecked reflection.

Create a codec that walks `FieldSpec` and handles Terraform's value states without resource-specific type switches:

- retain the distinction between known, unknown, and null values internally;
- encode known scalar values to API JSON using `APIName`;
- omit unknowns from requests;
- apply a single documented policy for known null values;
- decode JSON numbers without losing integer width or decimal stability;
- recursively encode/decode singleton objects and nested lists;
- apply mode filtering at every path;
- normalize API absence versus null consistently;
- validate and update reference pairs as one atomic unit;
- validate auto-assigned pairs and omit the value when the auto flag is true;
- preserve false, zero, and empty string when those are intentional values.

Use Terraform Framework `attr.Value`/typed values and `path.Path` access in the engine. API payloads use the provider-owned `WireObject` described in the transport contract below. Avoid reflection over generated OpenAPI model fields and avoid `interface{}` switches scattered across lifecycle code.

The first implementation should keep typed OpenAPI request adapters behind the transport registry. Generate those adapters from the spec/OpenAPI naming information. Moving to a fully generic HTTP transport can be evaluated later, after lifecycle parity; it should not be a prerequisite for eliminating resource duplication.

### 4. Generated transport adapter contract

The field codec and bulk manager need an explicit typed boundary. Today `createRequestPreparer` type-asserts each queued value to a concrete generated OpenAPI value type and wraps it in a concrete request type. A codec-produced `map[string]any` cannot be queued directly without panics or loss of explicit-null semantics.

Define a provider-owned wire representation with deterministic JSON semantics, for example a `WireValue` union for null/string/bool/int64/decimal/object/list and `WireObject map[string]WireValue`. Decimal and integer values must retain their canonical lexical representation. This is the only output of the field codec.

Generate one adapter per API operation/resource implementing a contract equivalent to:

```go
type TransportAdapter interface {
    BuildRequest(operation Operation, resources map[string]WireObject) (any, error)
    Execute(ctx context.Context, client *openapi.APIClient, operation Operation,
        request any, parameters map[string]string) (*http.Response, error)
    ResponseCollection(raw map[string]json.RawMessage,
        parameters map[string]string) (map[string]json.RawMessage, error)
}
```

Contract boundaries:

- the codec owns field names, field values, null/omission, and nested encoding;
- the adapter owns the outer name-keyed resource map, request wrapper, generated OpenAPI request/value types, method invocation, and fixed endpoint parameters;
- the bulk manager owns queueing, batching, ordering, waiting, retries, and response caching;
- no adapter may apply Terraform diff/default/reference behavior;
- conversion failures return diagnostics/errors and must never be unchecked type assertions;
- adapters must preserve explicit JSON null separately from omitted fields.

The initial bulk-manager change should replace the large request-preparer type switch with `config.Adapter.BuildRequest`, while leaving scheduling and execution order intact. The adapter may use generated assignment code or a validated JSON-to-SDK conversion, but the choice is accepted only after wire-format tests demonstrate equivalence.

#### Transport gate

Generate an adapter for `verity_ipv4_list` and test PUT and PATCH before building `GenericResource`. For each operation, compare the marshaled request with the current implementation's semantic JSON. Include name wrapper, false, empty string, omitted fields, endpoint execution, response extraction, and error handling. Add a focused synthetic nullable type test to prove that explicit null survives the chosen adapter conversion even though IPv4 List itself has no nullable field.

The gate passes when the bulk manager can queue `WireObject`, build the exact concrete SDK request, invoke the mock endpoint, and cache the response without resource-specific type switches or panics.

### 5. One generic lifecycle resource

Implement `GenericResource{spec *ResourceSpec}` once. Provider registration becomes a stable constructor closure per spec in the build-selected registry.

#### Schema

Compile `FieldSpec` recursively into Framework attributes/blocks. Apply descriptions, validators, defaults, plan modifiers, sensitivity, mode behavior, and exact nested cardinality from the spec.

The schema must be available before provider configuration. Specs should therefore be generated/embedded at build time or loaded from files embedded in the binary. Fetching an arbitrary schema from the configured Verity server at runtime is not a safe target: Terraform requests stable provider schemas before normal resource CRUD, and installed provider behavior must not change based on server availability.

#### Provider registration contract

`verityProvider.Resources()` must return the same deterministic constructor list before and after `Configure`. It should enumerate every resource in the build-selected API-version registry in canonical Terraform-name order. Remove context lookup, concrete-type-name matching, configured-mode filtering, and the “return all resources” fallback from registration.

Mode is a runtime compatibility rule, not a registration rule. A datacenter/campus provider binary exposes the same resource schemas; `ResourceSpec.Modes` is enforced by the generic resource during plan/apply with an actionable diagnostic when an incompatible resource is actually used. Per-field mode policies remain part of plan/encode/decode. Existing incompatible resources already present in state must still be readable enough to produce a clear migration diagnostic rather than disappearing from the provider schema.

API version is selected at build/generation time for the current single-version release. Server-version validation during provider configuration remains, but it does not add, remove, or reshape registered resources.

#### Create

1. Read plan/config into a generic typed value tree.
2. Run generated field and cross-field validation.
3. Encode all managed known values, applying reference, auto-assignment, nullable, mode, and nested strategies.
4. Submit through the bulk manager using `spec.API.BulkKey` and fixed parameters.
5. Set minimal identity state, then populate from the operation response or generic Read.

#### Read

1. Fetch through a spec-driven GET/cache adapter.
2. Extract `spec.API.ResponseCollectionKey` and find the object by the identity field.
3. Recursively decode it into state.
4. Apply configured-subset filtering only when the spec requests it.
5. Remove state when the object is not found.

#### Update

1. Diff plan and state by walking the spec.
2. Group paired reference/auto-assignment fields before scalar diffing.
3. Build the minimal PATCH object for ordinary fields.
4. Apply the selected collection strategy for nested fields.
5. Skip the operation for an empty patch; otherwise submit it and refresh state.

For `indexed_patch`, implement a tested generic algorithm with three explicit cases: new elements, changed existing elements, and index-only deletion markers. Handle multiple zero/server-assigned indices without storing them in a map keyed only by zero.

#### Delete and import

Use the canonical identity, delete parameter, endpoint parameters, and supported-operation flags from the spec. Most imports remain passthrough to the identity field. Update-only resources should report consistent diagnostics from the same operation policy.

#### Plan modification

Replace each resource's manual lists passed to `ModeFieldNullifier` with a recursive spec walk. The same pass should handle unavailable-mode paths, auto-assigned values, reference pair validation, and the chosen nullable reset semantics.

### 6. Resolve nullable/reset semantics explicitly

This deserves a design spike before migration. Terraform's decoded configuration does not generally preserve a robust distinction between an omitted optional attribute and an attribute explicitly assigned `null`; both commonly arrive as null. The existing source-file scanner recovers syntax in some cases, but it is outside Terraform's provider protocol and is not reliable enough for the target architecture.

Evaluate these options with protocol-level tests and document one provider-wide contract:

1. **Preferred where API defaults are stable:** model the default in Terraform, treat removal as reset to that default/null, and send a PATCH when desired state differs from prior state. Omit null fields only during create when the server should choose the default.
2. **For unmanaged computed values:** mark them computed-only, so the provider never treats omission as a request to clear.
3. **For a true three-way API distinction that Terraform cannot express:** expose an explicit generic reset mechanism, for example `reset_fields = ["peer_link_vlan"]`, or a narrowly scoped per-field reset flag. This is more honest and module-safe than parsing `.tf` files.
4. Retain HCL scanning only as a temporary compatibility adapter, behind one interface and with deprecation tests. It must not be required by newly generic resources. *(Adopted as the permanent contract instead; see the decision below.)*

Do not silently choose “null always means false/zero/empty.” The current scalar comparison helpers demonstrate why each kind needs a defined clear representation.

**Decision (2026-09-17): keep the `.tf` scanning as the permanent contract.** Terraform cannot tell the provider whether a nullable numeric was written as `null` or left out, and users need both: `x = null` must send an explicit API null that clears the value, and an absent `x` must leave the server's value alone. The existing configured-attribute scan is the mechanism that makes that distinction, so the generic engine preserves it rather than replacing it with a reset field or a reset-on-removal rule. Option 4 above is therefore adopted without the "temporary" qualification: the scan stays behind the engine's `Runtime` interface and is used for nullable fields on create, update, and plan.

The scan answers only two questions — whether the attribute is written, and what its configuration value is. What a null or unknown then sends is still decided by the field's declared `create_null`, `update_clear`, and `unknown_plan` policies, exactly as for every other field.

The scan records an attribute by name whatever its right-hand side is, so `x = var.y` counts as written, and a variable that resolves to null is treated exactly like a literal `null`: the field's policies decide whether that sends an API null, omits, defaults, or rejects, and with today's registry policies it sends an API null. Its real limits are where it looks and how it matches a resource:

- it reads only `*.tf` files directly in the provider working directory, so child modules in other directories and `.tf.json` files are not seen;
- it matches a resource by its `name` attribute when that is a literal string, and otherwise by block label, so a resource whose `name` comes from an expression and whose label differs from that name — typical with `count` or `for_each` — is not found;
- a run without the configuration files in the working directory finds nothing.

These are accepted rather than solved. When the scan does not find a resource, every nullable numeric on it is treated as not written: neither a value nor a null is sent for it, and the server keeps what it has. That is the handwritten resources' behavior today, so migrating a resource does not change it.

### 7. Consolidate registries and scheduling metadata

Generate these existing structures from `ResourceSpec`, then delete the hand-maintained copies after parity:

- `getAllResources` and resource compatibility;
- `ModeFields`;
- bulk `ResourceConfig` endpoint adapters;
- request/response/importer JSON-key maps;
- cache keys and operation keys;
- lifecycle mock metadata and basic documentation tables.

Move create/update/delete ordering to declarative dependency metadata consumed by the bulk manager. Compute a stable topological order per mode, and normally reverse it for delete. Keep an explicit phase/override mechanism when delete order is not simply the reverse.

Circular references are graph behavior, not field codec behavior. Preserve the existing route-map-clause/tenant workaround initially as a named bulk-manager hook. Later, describe the cycle and its break/restore strategy declaratively, but do not force all scheduler exceptions into generic field callbacks.

The non-API `verity_operation_stage` barrier should remain a small bespoke resource. It is not evidence that the generic API resource model needs arbitrary hooks.

### 8. Controlled escape hatches

The engine should cover the normal case completely while allowing named, testable exceptions:

```go
type ResourceHooks struct {
    Validate        ValidateHook
    BeforeEncode    EncodeHook
    NormalizeRead   NormalizeHook
    CustomizeDiff  DiffHook
}
```

Rules for hooks:

- empty by default and registered by name, not anonymous closures inside every spec;
- operate on provider-owned value/payload types, not generated SDK structs;
- cannot replace the complete lifecycle method;
- each hook needs a focused test and a comment explaining why schema metadata is insufficient;
- add a new generic strategy when the same hook appears a second time.

Likely initial exceptions are ACL endpoint variants, site/SFP response subset filtering, and dependency-cycle handling. Most current resource-specific code should not become a hook.

## Suggested package and file layout

```text
specs/
  openapi/6.6/
    datacenter.json               # committed canonical generator input
    campus.json                   # committed canonical generator input
    manifest.json                 # versions, provenance, SHA-256 checksums
  overrides.yaml                 # reviewed provider behavior only
  generated_manifest.json        # optional review/debug artifact

internal/spec/
  model.go                       # ResourceSpec and FieldSpec types
  validate.go                    # invariants and actionable errors
  registry_gen.go                # generated embedded registry

internal/genericresource/
  resource.go                    # Framework lifecycle implementation
  schema.go                      # FieldSpec -> Framework schema
  values.go                      # typed Terraform value tree
  encode.go                      # Terraform -> API payload
  decode.go                      # API response -> Terraform state
  diff.go                        # scalar/object recursive diff
  collections.go                 # indexed/replace/subset strategies
  plan.go                        # mode, references, auto assignment, reset
  hooks.go                       # named exception interfaces

internal/transport/
  registry_gen.go                # generated SDK endpoint adapters

tools/specgen/
  main.go                        # deterministic OpenAPI/spec generator
  normalize.go                   # canonicalizes maintainer-provided exports
```

The exact names are flexible. The important boundary is that resource definitions are data, lifecycle logic is centralized code, and OpenAPI SDK types exist only behind transport adapters.

## Example declarative resource

An illustrative override/spec for LAG could look like this; most descriptions, types, enums, modes, and nullability would be generated rather than repeated manually:

```yaml
resource:
  terraform_type: verity_lag
  modes: [datacenter, campus]
  versions: {min_inclusive: "6.6", max_exclusive: "6.7"}
  schema_version: 0
  identity: name
  api:
    path: /lags
    bulk_key: lag
    request_key: lag
    response_key: lag
    delete_parameter: lag_name

fields:
  name:
    access: required
    replace: true

  peer_link_vlan:
    ownership: configuration_or_server
    response_absence: terraform_null
    create_null: omit
    update_clear: api_null
    unknown_plan: preserve_state

  eth_port_profile:
    reference:
      type_field: eth_port_profile_ref_type_
      allowed_types: [eth_port_profile_, service_port_profile, pb_egress_profile]

  object_properties:
    collection:
      strategy: singleton
```

Adding an ordinary scalar or nested field should require only an OpenAPI update and, when inference is unambiguous, regeneration. Adding a conventional resource should require its OpenAPI endpoint plus a small resource-level override. Handwritten Go should be necessary only for a genuinely new lifecycle strategy.

## Migration plan

### Status: Phases 0, 1, and 2 are closed; Phase 3 is in progress (2026-09-17)

Both exit criteria are met, and every pre-migration gate in the executive
recommendation passes. Phase 2 closed as an opt-in pilot: its demonstrated
scalar parity permits the Phase 3 scalar-policy work to proceed. Phase 3 has not
met its exit criterion.

Evidence, all runnable from a clean checkout:

- `go test ./...` is green, with the CI delay variables from
  `.github/workflows/test.yml`.
- All five `specgen` verification/generation commands report no drift (four use
  `--check`), so the committed OpenAPI inputs, the registry, and the five
  generated metadata consumers regenerate offline and byte-identically.
- One registry covers all 50 API-backed resources. `verity_operation_stage` is
  excluded by design: it has no endpoint.
- 3,415 lifecycle-policy assertions compare the registry against evidence
  extracted from the legacy resource source, with 135 fields excluded by
  declared category. Four of the five policies are checked that way; the fifth,
  `state_ownership`, is a pure function of `access` and is enforced in
  validation instead.
- 149 golden wire fixtures pin PUT, PATCH, and post-apply state per resource,
  plus the schema golden file. They are the parity baseline Phase 2 is judged
  against.
- A live 6.6 run with `VERITY_GENERIC_RESOURCES=verity_ipv4_list` confirmed the
  generic IPv4 List implementation works: it registered, constructed the
  expected `ipv4_list_filter` request, and batched two creates into one
  `PUT /api/ipv4lists`.
- Every case the Phase 0 list names is covered, including the last one,
  "multiple new index-zero items"; see the carry-forward below for what it found.

Four things are carried forward rather than closed. None blocks Phase 2, whose
pilot is `verity_ipv4_list` — three scalar fields, no collection, nothing
nullable — so none of them is on its path:

1. **Adding an indexed child without naming its index does not work.** The API
   documents index zero as the way to append to a collection, and `index` is
   Optional and Computed in all 54 collections, but neither spelling survives an
   apply: `index = 0` is a known value the server then contradicts, and an
   omitted index is dropped by the guarded setter so the entry reaches the wire
   unnamed. Characterized by `tests/unit/lifecycle/index_zero_test.go`, which
   asserts today's behavior so Phase 4 has to change it deliberately. This is
   pre-existing and belongs to Phase 4's collection engine.
2. **`verity_tenant.vrf_name` needs a validator** and the override format cannot
   express one. `FieldSpec.Validators` exists and is validated; no override
   populates it. The one field-policy gap with no current representation.
3. **`unknown_plan: omit_and_read` is a create-path claim.** The six
   `CompareAndSet*Field` helpers are unguarded for every kind, so the update path
   does not implement it. The generic engine should implement the policy in both
   directions; see status.md.
4. **`3e76e76` is not independently green.** It regenerated the policy evidence
   ahead of the registry derivations that caught up in `e1c54fc`. Both are
   pushed, so the squash was declined rather than force-push a shared branch.

### Phase 0: characterize and freeze behavior — CLOSED

- Commit normalized datacenter/campus OpenAPI inputs and checksum manifests for API 6.6; make the existing SDK spec and generated files reproducible from them.
- Record current schema snapshots and request/response golden fixtures for every resource.
- Extend field-coverage tests to verify Read and PATCH coverage, not only schema discovery and PUT presence.
- Add cases for unknown, null, false, zero, empty string, reference pair transitions, auto-assignment transitions, multiple new index-zero items, nested deletion, mode exclusion, import, and missing API objects.
- Inventory every current resource exception and classify it as field policy, collection policy, operation policy, dependency policy, or true hook.
- Approve the API version/range and complete field-policy design, including deterministic literals, so Phase 1 can implement and validate it without reopening semantics.
- Complete the Framework value-bridge spike and retain its protocol-level tests.
- Complete the generated IPv4 List transport-adapter spike and exact wire-format tests, including a synthetic explicit-null case.
- Snapshot version-zero state types and implement/test the generic state-upgrader registration mechanism without changing current schemas.
- Lock the provider registration contract with a test that calls `Resources()` before and after configuration for both modes and asserts identical names, schemas, and ordering.
- Decide and document nullable/reset semantics before implementing the production generic codec.

Exit criterion: existing behavior is measurable, intended behavior changes are separated from refactoring changes, committed inputs regenerate offline, the field-policy design is approved, and the Framework, transport, registration, and state-contract spikes pass. Failure of a spike changes the implementation approach before the spec registry or production generic engine is built.

### Phase 1: introduce and validate the spec registry — CLOSED

- Add `ResourceSpec`/`FieldSpec` types and strict validation.
- Build the OpenAPI extractor and override merge.
- Generate current mode/version, compatibility, JSON key, bulk adapter, and constructor metadata without changing lifecycle implementations.
- Add `specgen --check` to CI.
- Report OpenAPI fields that are unrepresented, ambiguous, or overridden; do not silently guess.
- Validate every current resource and field against explicit API-version applicability and complete lifecycle policies.

Exit criterion: one generated registry accounts for every existing resource, field path, alias, operation, mode, version, lifecycle policy, and nested strategy while old resources still run, and every pre-migration gate in the executive recommendation passes. Only then may Phase 2 start the production generic-resource migration.

### Phase 2: generic scalar lifecycle pilot — CLOSED

- Implement schema compilation plus create/read/update/delete/import for scalar fields only.
- Migrate `verity_ipv4_list` behind a feature/build switch or in a parity test. It has only `name`, `enable`, and `ipv4_list`, with no nullable-source parser or nested block.
- Differentially run old and generic implementations against the same mock responses and compare schemas, requests, diagnostics, and final state.
- Keep the existing bulk manager unchanged except for accepting canonical spec metadata.

Exit criterion: the pilot has behavior parity and no handwritten model, mapper, nullifier, or lifecycle methods.

Where it stands: the engine is built and serves `verity_ipv4_list` behind the
`VERITY_GENERIC_RESOURCES` switch, reproducing the golden fixtures captured from
the handwritten resource byte for byte — the PUT, the PATCH, and the state after
apply. Schema compilation, the scalar codec, all five lifecycle methods, and the
mode nullifier are spec-driven, and the reviewed registry is embedded in the
binary. The bulk manager is unchanged, reached through a per-resource transport
adapter.

The switch remains off by default, so the provider registers what it always did.
That is intentional: Phase 2 closes the opt-in pilot, not the later decision to
make it the default or retire the handwritten
`resource_verity_ipv4_list.go`. A live 6.6 run additionally confirmed generic
registration, batching, adapter encoding, and transport. See
[status.md](status.md) for the evidence behind each item.

### Phase 3: shared semantic field policies — IN PROGRESS

- Implement reference pairs at top level and in nested objects/lists.
- Implement auto-assignment pairs.
- Implement the chosen nullable/reset policy for integers and numbers.
- Add singleton-object handling while preserving each resource's version-zero list-nested state shape.
- Migrate `verity_badge` only after nullable/reset and singleton policies pass their contract tests.
- Migrate several resources containing combinations of these policies, for example LAG plus a nullable/reference-heavy resource.
- Deprecate direct use of `Set*Fields`, `CompareAndSet*`, object-property switches, and manual mode nullifier lists in migrated resources.

Exit criterion: adding any supported scalar kind or semantic pair requires a spec change only.

Where it stands: the engine serves 16 resources behind the opt-in switch — six
scalar-only and ten whose only nested shape is a singleton `object_properties`
block. One rule, `genericresource.Supported`, decides that set, and the typed
transport adapters are generated for exactly what it accepts rather than written
per resource. All 16 have schema parity, blocks included, and golden-fixture
parity with the handwritten resources.

Implemented: reference pairs at the top level and inside a singleton, from
`FieldSpec.Reference`; nullable numerics under the configured-attribute scan
decided in section 6, with the declared policies deciding what is sent; and
singleton objects compiled to the version-zero list block, following the
handwritten create, update, and read rules. Differential tests compare both
implementations for each.

Remaining: auto-assignment pairs (which unblock `verity_service`), fixed headers
on the write path (the ACLs), and nullable members inside a singleton. Comparing
singleton updates also found three handwritten cases that do not converge, which
the engine reproduces for parity; whether one of them is a live defect depends on
whether the API replaces or merges a partially PATCHed `object_properties`. See
status.md.

### Phase 4: indexed collections

- Implement and property-test `indexed_patch`, server-assigned index, full replacement, ordering, and subset strategies.
- Start with a simple single-list resource such as an access/community/prefix list.
- Continue with resources containing multiple indexed lists.
- Verify list add/update/delete and API-assigned index refresh in both plan and state.

Exit criterion: indexed nested blocks no longer require resource-specific handler closures.

### Phase 5: complex and exceptional resources

- Migrate multi-block resources such as bundle, gateway, packet queue/broker, eth port profile/settings, and switchpoint.
- Migrate deep nesting (`site.object_properties.system_graphs`) only after recursive path behavior has dedicated tests.
- Add operation policies for update-only site and SFP breakout resources and configured-subset response filtering.
- Add ACL variants with fixed endpoint parameters and any response extractor metadata.
- Keep `verity_operation_stage` bespoke.

Exit criterion: all API-backed resources use the generic engine; only named hooks and the operation-stage resource remain handwritten.

### Phase 6: remove legacy duplication

- Delete migrated resource models and lifecycle files.
- Remove generated `ModeFields`, compatibility/key maps, handwritten provider constructor list, and redundant bulk registry entries after all consumers use `ResourceSpec`.
- Remove the HCL source parser after its compatibility window, or isolate it only for legacy state versions if required.
- Generate resource documentation and a new-resource checklist from the registry.
- Add a repository check that forbids new handwritten API-backed lifecycle resources without an approved exception.

Exit criterion: a normal OpenAPI field addition regenerates cleanly, a normal resource addition is declarative, and CI proves all declared fields participate in schema, encode, diff, and decode.

## Test strategy

### Spec and generator tests

- committed source checksums and fully offline regeneration from the selected API version;
- deterministic output and clean `--check` generation;
- every endpoint/resource has one canonical identity and all required aliases;
- all resource/field paths have explicit mode and API-version ranges plus complete lifecycle policies;
- all default/reset literals have deterministic representation and validate against the field kind;
- reference and auto-assignment pairs point to existing same-scope fields;
- nested collections declare identity, cardinality, and patch strategy;
- overrides fail if their target disappears or changes type in OpenAPI;
- generated Terraform names are unique and stable.

### Codec contract tests

Use table/property-style tests independent of any individual resource:

- every scalar kind across known/null/unknown and mode states;
- integer-width and decimal round-trip behavior;
- create omission versus update clearing;
- reference pair add/change/clear and one-type versus multi-type validation;
- auto-assigned true/false transitions;
- singleton object presence and absence;
- indexed collection add/change/delete/reorder, duplicate/zero index handling, and API index assignment;
- nested mode-specific fields and recursive blocks.

The Framework value-bridge spike is a prerequisite to these tests and remains in the suite as a contract test for the Framework version in `go.mod`.

### Transport adapter tests

- canonical `WireObject` to the exact generated OpenAPI request/value type;
- semantically exact PUT/PATCH wire JSON, including wrapper keys and resource-name maps;
- distinction between omitted and explicit-null fields;
- preservation of false, zero, empty string, int64 width, and decimal text;
- fixed parameters/header splitting such as ACL `ip_version`;
- response collection extraction and errors for malformed/unexpected response shapes;
- conversion errors returned normally, with no type-assertion panic.

### Resource parity tests

- schema snapshots before and after migration;
- golden PUT/PATCH/DELETE requests and Read state for all specs;
- differential old-versus-generic tests during each phase;
- mock-server lifecycle coverage generated from the registry instead of a second manual `allResourceTests` table;
- acceptance/integration coverage for at least one resource in every behavior class and both modes.

### State compatibility

- Treat every current resource schema as version 0. Before the first migration, capture a machine-readable snapshot containing attribute/block names, Terraform types, nesting modes, access flags, validators/plan modifiers, and the resulting Framework state type.
- The initial generic implementation must publish the identical version-0 schema and state type for a migrated resource. Swapping lifecycle implementation alone does not justify a schema-version increment.
- Add a generic implementation of `ResourceWithUpgradeState` before making any representation change. `ResourceSpec` should carry `SchemaVersion` plus named prior-version descriptors. Generated registration must associate each prior version with its exact Framework `PriorSchema`/state type and a tested upgrader function.
- State-upgrader tests must start from raw version-0 state produced by the last released provider, run the Framework upgrade path, then refresh against mock API data and assert no unintended diff.
- Preserve resource addresses, Terraform type names, identity values, attribute names, and nested list ordering during the lifecycle migration.
- Keep `object_properties` as its current list-nested representation during generic migration. Converting it to a single-nested block is a separate post-parity change: increment that resource's schema version, define list-empty/list-one/null conversions, decide what to do with invalid multi-element legacy state, and test configuration syntax/state migration.
- CI must reject a generated schema/state-type diff unless the resource version is incremented and an upgrader plus old-state fixture is present.

## Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| OpenAPI is incomplete or inaccurate | Use strict inference plus small reviewed overrides; never silently infer ambiguous lifecycle semantics. |
| Raw inputs exist only on developer machines | Commit canonical mode/version inputs and checksums; make CI regeneration offline. |
| A generic engine hides behavior and becomes hard to debug | Include resource and field path in every diagnostic/log; make encoded payload and selected strategy observable at debug level. |
| Framework generic values lose null/unknown or nested type information | Require the protocol-level value-bridge spike before production codec work; fall back to generated thin typed bindings if it fails. |
| Reflection or untyped maps cause runtime failures | Use typed Terraform values, validate specs at startup/test time, and limit maps to the JSON boundary. |
| Codec output is incompatible with concrete OpenAPI request types | Use generated transport adapters and exact wire-format tests before migrating the first resource. |
| Large state churn after schema generation | Snapshot schemas, preserve names/types first, and use explicit state upgraders for later improvements. |
| PATCH semantics differ by array/resource | Require an explicit collection strategy; add hooks only for proven exceptions. |
| Generator changes every resource at once | Pin source API version, make output deterministic, and review generated manifest diffs separately from engine changes. |
| Bulk scheduling refactor destabilizes lifecycle migration | Generate adapters first and preserve the current manager; move ordering to metadata in a later isolated phase. |
| Nullable behavior cannot be represented by omission/null alone | Keep the configured-attribute scan as the documented contract (section 6), behind the `Runtime` interface, with its matching limits stated and the declared lifecycle policies deciding what a null sends. |
| Runtime server schema changes Terraform's public contract | Embed a build-time schema registry and continue matching provider releases to supported API versions. |

## Approaches not recommended as the final design

### Generate another full resource Go file per resource

Templating the current 400–1,900 line files would reduce typing but keep lifecycle behavior duplicated in generated output and make template changes produce very large diffs. Code generation is appropriate for the spec and small transport adapters; behavior should live in the generic engine.

### Use only reflection over generated OpenAPI structs

This couples Terraform behavior to generator-specific Go names and pointer wrappers, moves mistakes to runtime, and still does not provide collection, mode, reference, or reset semantics. Generated SDK structs should stay behind adapters.

### Load schemas dynamically from each configured server

Terraform expects a stable provider schema before CRUD and often before usable provider configuration. A dynamic server schema would also make the same provider version expose different contracts and prevent offline validation. Generate from server/OpenAPI inputs during provider development or release, then embed the result.

### Keep adding narrowly typed helpers

More `SetXFields`, `MapXWithModeNested`, and `CompareAndSetX` helpers shorten repeated code but do not remove the need to enumerate each field in every lifecycle stage. The field list must become data traversed by common logic.

## Definition of done

- A conventional scalar, reference pair, auto-assigned pair, singleton block, or indexed nested field is implemented by OpenAPI regeneration plus at most a declarative override.
- The provider can be regenerated offline from committed, checksummed datacenter/campus inputs for its selected API version.
- A conventional new resource requires no handwritten lifecycle, model, state population, mode-nullification, registry, importer-key, or mock-table code.
- One validated registry is the source of every resource and field alias/mode/operation decision.
- Resource and field API-version ranges plus create/read/update/unknown/ownership policies are explicit and validated.
- `Resources()` and all published schemas are stable before/after provider configuration and across configured modes.
- The Framework generic value bridge and generated transport boundary have permanent protocol/wire contract tests.
- Generic encode, decode, diff, plan, and collection behavior is unit-tested independently of resources.
- CI detects an API field added, removed, renamed, retyped, or left without a Terraform policy.
- Existing state upgrades and refreshes without unintended diffs.
- API-backed resources share the generic engine; remaining custom behavior is a short, named, tested list.

## Recommended first implementation slice

Start with Phase 0 only, split into reviewable changes:

1. commit and verify normalized 6.6 mode-specific inputs, then capture schema/state and wire baselines;
2. approve the versioned resource/field policy model and its strict validation rules;
3. prove the generic Terraform Framework value bridge in tests;
4. generate and prove the IPv4 List transport adapter against exact wire JSON; and
5. lock stable provider registration and the state-upgrader contract.

Only after those gates pass, introduce the production generic lifecycle engine and use `verity_ipv4_list` as its first resource. It is the clean scalar pilot: `name`, `enable`, and `ipv4_list`, with no HCL nullable-source parsing and no nested collection. Move `verity_badge` to Phase 3 after reset/null and singleton-list behavior are proven. Follow Badge with a simple indexed resource such as `verity_community_list`.

This sequence proves the central promise early: after the codec supports a field category, future fields of that category require schema/spec data, not repeated Create/Read/Update/ModifyPlan edits.

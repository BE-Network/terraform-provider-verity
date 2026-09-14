# Resource exception inventory

Phase 0 of [refactor_plan.md](../refactor_plan.md) asks for every current resource
exception, classified as field policy, collection policy, operation policy,
dependency policy, or true hook. The classification decides where each one lives
after the refactor: the first four become declarative spec data, and only a true
hook stays as handwritten code.

An exception here means a resource deviating from the contract the other 49
follow. Behavior every resource shares is not an exception, however unusual it
looks: `WithPostOperationFallback`, `SetPostOperationFallbackState`,
`HasPendingOrRecentOperations`, and bulk-operation notification appear in all 49
API-backed resources and are the common lifecycle, not deviations from it.

Counts are derived from `specs/generated_registry.json` and the resource
implementations; see [status.md](../status.md) for how each was verified.

## Field policy

Expressible as `FieldSpec` data. No code needed after migration.

| Exception | Scope | Detail |
| --- | --- | --- |
| Reference pairs | 92 pairs across 29 resources | A value field and its `*_ref_type_` companion must be written together. Modelled as `FieldSpec.Reference`, with allowed types read from the companion's OpenAPI enum. |
| Multiple permitted reference types | 5 resources | `verity_packet_broker`, `verity_gateway`, `verity_bundle`, `verity_lag`, `verity_threshold_group` use `HandleMultipleRefTypesSupported`; their enums list more than one target type. Same `Reference` shape, longer `AllowedTypes`. |
| Auto-assignment pairs | 15 pairs across 4 resources | `verity_fabric`, `verity_service`, `verity_tenant`, `verity_switchpoint`. A boolean `*_auto_assigned_` flag makes the server choose the value. Modelled as `FieldSpec.AutoAssignment`. |
| Numeric nullability | all numerics except `index` | `tools/process_swagger.py` marks every number and integer nullable except one named `index`. Drives `update_clear` and `create_null`; the committed documents predate the transform. |
| Singleton members clear by omission | 38 fields | A member of an object sent whole is cleared by dropping its key, not by sending a zero value. Required adding `UpdateClearOmit`, which the Phase 0 vocabulary lacked. |
| Nullable member inside a singleton | `verity_switchpoint.object_properties.number_of_multipoints` | Sends an explicit null where its siblings omit. |
| Guarded `Set*Fields` in a collection | 2 fields, both on `verity_as_path_access_list.lists` | Cleared by omission where the same field on sibling resources sends a zero value. |
| API object declared with no properties | 7 endpoints | `object_properties` is an empty object. Two resources ship it as an empty block; five do not surface it at all, recorded as `unmanaged`. |
| Field the API documents with an empty description | `verity_gateway.fabric_interconnect` | The OpenAPI document itself carries `"description": ""`. Recorded with `undocumented: true` rather than inventing text. |
| Value constrained beyond its type | `verity_tenant.vrf_name` | Rejects a generated placeholder; the coverage harness supplies a valid value. Needs a validator in the spec, which the format does not yet express. |

## Collection policy

Expressible as `CollectionSpec`.

| Exception | Scope | Detail |
| --- | --- | --- |
| Indexed child collections | 54 collections | Entries are addressed by an `index` the server assigns. `index` is the collection identity: always sent, never compared or cleared. |
| Singleton objects exposed as blocks | 27 collections | OpenAPI declares an object; every resource models it as a single-entry `ListNestedBlock`. The generated object and the legacy list block are reconciled explicitly. |
| Blocks nested inside blocks | `verity_fabric.object_properties.system_graphs` | An indexed collection inside a singleton. The only two-level nesting in the provider. |

## Operation policy

Expressible as `OperationSpec` plus transport metadata.

| Exception | Scope | Detail |
| --- | --- | --- |
| One endpoint, two resources | `verity_acl_v4`, `verity_acl_v6` | `/acls` backs both, discriminated by a required `ip_version` header. Needs header-aware bulk PUT/PATCH/DELETE/GET, a `HeaderResponseExtractor` because the response key differs by version (`ipv4_filter` / `ipv6_filter`), distinct cache keys, and the header as a required delete query parameter. Modelled as `FixedHeaders`; bulk key plus discriminator is what must stay unique. |
| Neither create nor delete | `verity_sfp_breakout` | Both return a hard diagnostic: the resource represents existing hardware that can only be read and updated. It is brought into state by import for its read fixture, and its update fixture is recorded by `TestUpdateOnlyResourceGoldenPatch`, which drives the bulk manager directly because a test case that must destroy what it creates cannot hold this resource across an update. It has no create fixture and never will. |
| Field with no update path | `verity_packet_broker.ipv6_permit.enable` | The collection's `UpdateExisting` handles only its reference pair and index, so the flag cannot be changed on an existing entry. A legacy gap, recorded rather than replicated. |
| DELETE with no required identifying parameter | `/alarms/mask` | Supports delete but exposes only optional array parameters, so nothing names the objects to remove. Unrepresented. |
| PATCH-only endpoints | `/sdlcs/upgrade`, `/switchpoints/upgrade`, `/systemconfig` | No GET, so no read path. Unrepresented. |
| Endpoint with no mutable fields | `/readmode` | Unrepresented. |

## Dependency policy

Expressible as `DependencySpec` and consumed by the bulk manager.

| Exception | Scope | Detail |
| --- | --- | --- |
| Operation ordering | all resources | PUT follows a fixed order per mode, PATCH the same order with SFP Breakouts first, DELETE the reverse. Recorded in `order_of_operations.txt`, sourced from `internal/bulkops/execution.go`. |
| ACL batch ordering | `verity_acl_v4`, `verity_acl_v6` | Always IPv6 before IPv4 for PUT and PATCH, reversed for DELETE. |
| Import stage ordering | importer | Generated stages place Fabric before Gateway; campus stages exclude datacenter-only resources. |

## True hook

Genuinely handwritten after migration.

| Exception | Detail |
| --- | --- |
| `verity_operation_stage` | Not API-backed: no endpoint, no registry entry. The plan keeps it bespoke (lines 402, 571, 573). It is the only resource excluded from the registry by design. |

## What this inventory changes

Nothing in the four policy categories needs handwritten code after migration;
each is already modelled in `internal/spec` or has a named place to go. Two items
are the exceptions to that:

- **`verity_tenant.vrf_name`** needs a validator, and the override format has no
  way to express one. `FieldSpec.Validators` exists and is validated, but no
  override populates it. This is the one field-policy gap with no current
  representation.
- **`verity_packet_broker.ipv6_permit.enable`** is most likely a defect rather
  than a policy. Replicating it in the generic engine would preserve a bug; the
  registry records the absence so the decision is deliberate.

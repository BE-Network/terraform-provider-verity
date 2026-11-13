# Upgrading to Verity API 6.6.1.39

This provider release supports Verity API 6.6.1.39. That build no longer
accepts an object's name inside its body, so this release does not support
earlier 6.6 builds. The provider's version check compares only the
major and minor API version and cannot detect an earlier 6.6 build.

## Fresh import required

The API removes fields that existing Terraform state for these resources still
records:

| Resource | Removed |
| --- | --- |
| `verity_switchpoint` | `tenant`, `tenant_ref_type_` |
| `verity_route_map_clause` | `object_properties.match_fields_shown` |

Follow the [fresh-import procedure](fresh_import_upgrades.md) with this provider
release. Remove the removed arguments from any configuration you carry over.

## Object names

The API no longer repeats an object's name inside its body. The provider now uses
the `name` attribute only as the object's key in requests and responses. Existing
configuration keeps its `name` arguments unchanged. Names keep their exact
spelling; only generated Terraform resource labels are sanitized.

## Other changes

- `verity_eth_port_settings`: `max_allowed_unit` accepts `kbps` instead of
  `Kpps`. Update configurations that set `Kpps`. The new
  `priority_flow_control_cable_length` argument is available.
- `verity_route_map_clause`: new `set_*` arguments, including the
  `set_community_list_delete` reference.
- `verity_switchpoint`: `eths` entries accept `eth_num_role`.
- `verity_tenant`: new `hgx_hostnames` block (datacenter).
- New `verity_static_ip` resource (campus).

## Behaviour changes in this provider release

These changes come from fixes for destroy and re-create. Some need changes to
existing configurations; the simplest path is the
[fresh-import procedure](fresh_import_upgrades.md), which generates
configuration that already follows all of them.

### Configuration changes

- **Plaintext password fields are write-only.** `verity_switchpoint`
  `password`, `private_password`, `passphrase`, `ssh_key_or_password` and
  `verity_gateway` `md5_password` are `sensitive`, and the API returns them
  empty, so Terraform keeps the configured value instead of reporting an
  "inconsistent result" after setting a password. The importer doesn't write
  them. Their `_encrypted` companions are `sensitive` too; on 6.6 the API
  returns the real hash, so they are imported and restored as before. To set
  or change a password, put the plaintext in the configuration and remove the
  imported `_encrypted` line; setting both is rejected at plan time, because
  the API computes a new hash.
- **`port_name` is read-only** in `verity_bundle.eth_port_paths` and
  `verity_switchpoint.eths`. Ports are identified by `index`; the API fills the
  name in from the device. Remove `port_name` lines from existing
  configurations (Terraform rejects them as "Invalid Configuration for
  Read-Only Attribute"), or re-import.
- **One stage for route maps, policy-based routing and tenants** (datacenter).
  The importer now puts `verity_route_map_clause`, `verity_pb_routing_acl`,
  `verity_route_map`, `verity_pb_routing` and `verity_tenant` behind one
  `verity_operation_stage`, `route_map_tenant_stage`. Configurations generated
  before this release have a separate stage for each type, and re-creating
  route map clauses that reference a tenant then fails ("Cannot resolve object
  reference"). Regenerate the configuration, or make all five types depend on
  one stage.

### Behaviour changes without configuration changes

- **Destroy:** predefined objects (for example `(Packet Queue)`,
  `(Device Settings)`) are removed from Terraform state with a warning instead
  of failing. A blocked object no longer fails a whole DELETE request; the
  others are deleted and the blocked one fails with the API's reason. Route map
  clauses deleted together with the tenant they reference get `match_vrf`
  cleared first.
- **`verity_sfp_breakout`:** the API can neither create nor delete it. Destroy
  removes it from state and leaves its settings applied; applying the
  configuration again adopts the existing object and sends only the settings
  that differ. No import needed.
- **Route map clauses created together with their tenant** are created without
  `match_vrf`, which is set right after the tenant exists. A failed restore can
  leave the clause without `match_vrf`; the next apply repairs it.
- **Indexed lists** (for example `breakout`, `eths`) keep the order of the
  configuration; reordering entries no longer leaves a permanent diff.
- **Errors** show the API's reason instead of only "400 Bad Request".
- **Parallelism:** set it above the number of resources in the configuration,
  for `apply` and `destroy` (`TF_CLI_ARGS_destroy`). See the README.

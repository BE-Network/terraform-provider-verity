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

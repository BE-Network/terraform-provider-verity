# Upgrading to Verity API 6.7

This provider release supports Verity API 6.7 only. The provider checks the
server's major and minor API version and refuses a 6.6 server. Use a 6.6
provider release with 6.6 servers.

## Fresh import required

API 6.7 removes fields that existing Terraform state for these resources still
records:

| Resource | Removed |
| --- | --- |
| `verity_switchpoint` | `tenant`, `tenant_ref_type_` |
| `verity_route_map_clause` | `object_properties.match_fields_shown` |
| `verity_fabric` | `enable` |
| `verity_pair` | `enable` |
| `verity_plane` | `enable` |
| `verity_pod` | `enable` |
| `verity_rack` | `enable` |
| `verity_spine_plane` | `enable` |
| `verity_ssp_group` | `enable` |
| `verity_su` | `enable` |

API 6.7.0.10 also removes `enable` from `verity_fabric_collection`. A fabric
collection now has no arguments besides `name`: it is an empty object that
fabrics join through `verity_fabric.domain_for_fabric`.

Follow the [fresh-import procedure](fresh_import_upgrades.md) with this provider
release. Remove the removed arguments from any configuration you carry over.

## Object names

API 6.7 no longer repeats an object's name inside its body. The provider now uses
the `name` attribute only as the object's key in requests and responses. Existing
configuration keeps its `name` arguments unchanged. Names keep their exact
spelling; only generated Terraform resource labels are sanitized.

## Other changes

- `verity_device_settings`: `ntp_vrf_tenant` and `ntp_vrf_tenant_ref_type_` are
  datacenter-only. In campus mode, `ntp_vrf` no longer accepts `tenant`.
- `verity_fabric`: new `underlay_link_addressing` argument (`unnumbered`,
  `ipv4_numbered` or `ipv6_numbered`; the API default is `unnumbered`).
- `verity_bundle` and `verity_switchpoint`: `port_name` in `eth_port_paths` and
  `eths` is read-only. The API fills it in from the device. Remove `port_name`
  lines from carried-over configuration.
- Password fields (`verity_switchpoint`: `password`, `private_password`,
  `passphrase`, `ssh_key_or_password`; `verity_gateway`: `md5_password`; each
  with its `_encrypted` pair) are never imported, because the API doesn't return
  them. Set them in configuration after a fresh import.
- `verity_eth_port_settings`: `max_allowed_unit` accepts `kbps` instead of
  `Kpps`. Update configurations that set `Kpps`. The new
  `priority_flow_control_cable_length` argument is available.
- `verity_route_map_clause`: new `set_*` arguments, including the
  `set_community_list_delete` reference.
- `verity_switchpoint`: `eths` entries accept `eth_num_role_icon`.
- `verity_tenant`: new `hgx_hostnames` block (datacenter).
- New `verity_static_ip` resource (campus).

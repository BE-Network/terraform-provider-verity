# verity_route_map_clause (Resource)

Manages a Verity Route Map Clause.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_route_map_clause" "example" {
  name = "example"
  enable = false
  match_as_path_access_list = ""
  match_as_path_access_list_ref_type_ = "as_path_access_list"
  match_community_list = ""
  match_community_list_ref_type_ = "community_list"
  match_evpn_route_type = ""
  match_evpn_route_type_default = false
  match_extended_community_list = ""
  match_extended_community_list_ref_type_ = "extended_community_list"
  match_interface_number = null
  match_interface_vlan = null
  match_ipv4_address_ip_prefix_list = ""
  match_ipv4_address_ip_prefix_list_ref_type_ = "ipv4_prefix_list"
  match_ipv4_next_hop_ip_prefix_list = ""
  match_ipv4_next_hop_ip_prefix_list_ref_type_ = "ipv4_prefix_list"
  match_ipv6_address_ipv6_prefix_list = ""
  match_ipv6_address_ipv6_prefix_list_ref_type_ = "ipv6_prefix_list"
  match_ipv6_next_hop_ipv6_prefix_list = ""
  match_ipv6_next_hop_ipv6_prefix_list_ref_type_ = "ipv6_prefix_list"
  match_local_preference = null
  match_metric = null
  match_origin = ""
  match_peer_interface = null
  match_peer_ip_address = ""
  match_peer_vlan = null
  match_source_protocol = ""
  match_tag = null
  match_vni = null
  match_vrf = ""
  match_vrf_ref_type_ = "tenant"
  permit_deny = ""

  object_properties {
    match_fields_shown = ""
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable flag of this provisioning object.
* `match_as_path_access_list` (String) - Match AS Path Access List. Set together with `match_as_path_access_list_ref_type_`.
* `match_as_path_access_list_ref_type_` (String) - Object type for match_as_path_access_list field.
* `match_community_list` (String) - Match Community List. Set together with `match_community_list_ref_type_`.
* `match_community_list_ref_type_` (String) - Object type for match_community_list field.
* `match_evpn_route_type` (String) - Match based on the indicated EVPN Route Type.
* `match_evpn_route_type_default` (Boolean) - Match based on the type of EVPN Route Type being Default".
* `match_extended_community_list` (String) - Match Extended Community List. Set together with `match_extended_community_list_ref_type_`.
* `match_extended_community_list_ref_type_` (String) - Object type for match_extended_community_list field.
* `match_interface_number` (Integer) - Match Interface Number. Set it to `null` to clear it.
* `match_interface_vlan` (Integer) - Match Interface VLAN. Set it to `null` to clear it.
* `match_ipv4_address_ip_prefix_list` (String) - Match IPv4 Address IPv4 Prefix List. Set together with `match_ipv4_address_ip_prefix_list_ref_type_`.
* `match_ipv4_address_ip_prefix_list_ref_type_` (String) - Object type for match_ipv4_address_ip_prefix_list field.
* `match_ipv4_next_hop_ip_prefix_list` (String) - Match IPv4 Next Hop IPv4 Prefix List. Set together with `match_ipv4_next_hop_ip_prefix_list_ref_type_`.
* `match_ipv4_next_hop_ip_prefix_list_ref_type_` (String) - Object type for match_ipv4_next_hop_ip_prefix_list field.
* `match_ipv6_address_ipv6_prefix_list` (String) - Match IPv4 Address IPv6 Prefix List. Set together with `match_ipv6_address_ipv6_prefix_list_ref_type_`.
* `match_ipv6_address_ipv6_prefix_list_ref_type_` (String) - Object type for match_ipv6_address_ipv6_prefix_list field.
* `match_ipv6_next_hop_ipv6_prefix_list` (String) - Match IPv6 Next Hop IPv6 Prefix List. Set together with `match_ipv6_next_hop_ipv6_prefix_list_ref_type_`.
* `match_ipv6_next_hop_ipv6_prefix_list_ref_type_` (String) - Object type for match_ipv6_next_hop_ipv6_prefix_list field.
* `match_local_preference` (Integer) - Match BGP Local Preference value on the route. Set it to `null` to clear it.
* `match_metric` (Integer) - Match Metric of the IP route entry. Set it to `null` to clear it.
* `match_origin` (String) - Match routes based on the value of the BGP Origin attribute.
* `match_peer_interface` (Integer) - Match BGP Peer port the route was learned from. Set it to `null` to clear it.
* `match_peer_ip_address` (String) - Match BGP Peer IP Address the route was learned from.
* `match_peer_vlan` (Integer) - Match BGP Peer VLAN over which the route was learned. Set it to `null` to clear it.
* `match_source_protocol` (String) - Match Routing Protocol the route originated from.
* `match_tag` (Integer) - Match routes that have this value for a Tag attribute. Set it to `null` to clear it.
* `match_vni` (Integer) - Match based on the VNI value. Set it to `null` to clear it.
* `match_vrf` (String) - Match VRF the route is associated with. Set together with `match_vrf_ref_type_`.
* `match_vrf_ref_type_` (String) - Object type for match_vrf field.
* `object_properties` (Block) - Object properties for the Route Map Clause. At most one block.
  * `match_fields_shown` (String) - Match fields shown.
  * `notes` (String) - User Notes.
* `permit_deny` (String) - Action upon match of Community Strings.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `match_as_path_access_list` | `match_as_path_access_list_ref_type_` | `as_path_access_list` |
| `match_community_list` | `match_community_list_ref_type_` | `community_list` |
| `match_extended_community_list` | `match_extended_community_list_ref_type_` | `extended_community_list` |
| `match_ipv4_address_ip_prefix_list` | `match_ipv4_address_ip_prefix_list_ref_type_` | `ipv4_prefix_list` |
| `match_ipv4_next_hop_ip_prefix_list` | `match_ipv4_next_hop_ip_prefix_list_ref_type_` | `ipv4_prefix_list` |
| `match_ipv6_address_ipv6_prefix_list` | `match_ipv6_address_ipv6_prefix_list_ref_type_` | `ipv6_prefix_list` |
| `match_ipv6_next_hop_ipv6_prefix_list` | `match_ipv6_next_hop_ipv6_prefix_list_ref_type_` | `ipv6_prefix_list` |
| `match_vrf` | `match_vrf_ref_type_` | `tenant` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_route_map_clause.<resource_name> <name>
```

# Route Map Clause Resource

Provides a Verity Route Map Clause resource. Route Map Clauses are used to define match and permit/deny rules for routing policies.

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
  match_evpn_route_type_default = null
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
  permit_deny = "permit"

  object_properties {
    match_fields_shown = ""
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable flag of this provisioning object
* `permit_deny` (String) - Action upon match of Community Strings
* `match_as_path_access_list` (String) - Match AS Path Access List
* `match_as_path_access_list_ref_type_` (String) - Object type for match_as_path_access_list field
* `match_community_list` (String) - Match Community List
* `match_community_list_ref_type_` (String) - Object type for match_community_list field
* `match_extended_community_list` (String) - Match Extended Community List
* `match_extended_community_list_ref_type_` (String) - Object type for match_extended_community_list field
* `match_interface_number` (Integer) - Match Interface Number
* `match_interface_vlan` (Integer) - Match Interface VLAN
* `match_ipv4_address_ip_prefix_list` (String) - Match IPv4 Address IPv4 Prefix List
* `match_ipv4_address_ip_prefix_list_ref_type_` (String) - Object type for match_ipv4_address_ip_prefix_list field
* `match_ipv4_next_hop_ip_prefix_list` (String) - Match IPv4 Next Hop IPv4 Prefix List
* `match_ipv4_next_hop_ip_prefix_list_ref_type_` (String) - Object type for match_ipv4_next_hop_ip_prefix_list field
* `match_local_preference` (Integer) - Match BGP Local Preference value on the route
* `match_metric` (Integer) - Match Metric of the IP route entry
* `match_origin` (String) - Match routes based on the value of the BGP Origin attribute
* `match_peer_ip_address` (String) - Match BGP Peer IP Address the route was learned from
* `match_peer_interface` (Integer) - Match BGP Peer port the route was learned from
* `match_peer_vlan` (Integer) - Match BGP Peer VLAN over which the route was learned
* `match_source_protocol` (String) - Match Routing Protocol the route originated from
* `match_vrf` (String) - Match VRF the route is associated with
* `match_vrf_ref_type_` (String) - Object type for match_vrf field
* `match_tag` (Integer) - Match routes that have this value for a Tag attribute
* `match_evpn_route_type_default` (Boolean) - Match based on the type of EVPN Route Type being Default
* `match_evpn_route_type` (String) - Match based on the indicated EVPN Route Type
* `match_vni` (Integer) - Match based on the VNI value
* `object_properties` (Object) - Additional object properties
  * `notes` (String) - User Notes
  * `match_fields_shown` (String) - Match fields shown
* `match_ipv6_address_ipv6_prefix_list` (String) - Match IPv4 Address IPv6 Prefix List
* `match_ipv6_address_ipv6_prefix_list_ref_type_` (String) - Object type for match_ipv6_address_ipv6_prefix_list field
* `match_ipv6_next_hop_ipv6_prefix_list` (String) - Match IPv6 Next Hop IPv6 Prefix List
* `match_ipv6_next_hop_ipv6_prefix_list_ref_type_` (String) - Object type for match_ipv6_next_hop_ipv6_prefix_list field

## Import

Route Map Clauses can be imported using the name:

```sh
terraform import verity_route_map_clause.<resource_name> <name>
```

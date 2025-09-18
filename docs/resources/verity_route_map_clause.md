# Route Map Clause Resource

Provides a Verity Route Map Clause resource. Route Map Clauses are used to define match and permit/deny rules for routing policies.

## Example Usage

```hcl
resource "verity_route_map_clause" "test1" {
  name = "test1"
  depends_on = [verity_operation_stage.route_map_clause_stage]
  object_properties {
    match_fields_shown = ""
    notes = ""
  }
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
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable flag of this provisioning object. Default: `false`.
- `permit_deny` (String, Optional) — Action upon match of Community Strings. Allowed values: `"permit"`, `"deny"`. Default: `"permit"`.
- `match_as_path_access_list` (String, Optional) — Match AS Path Access List. Default: `""`.
- `match_as_path_access_list_ref_type_` (String, Optional) — Object type for match_as_path_access_list field. Allowed value: `"as_path_access_list"`.
- `match_community_list` (String, Optional) — Match Community List. Default: `""`.
- `match_community_list_ref_type_` (String, Optional) — Object type for match_community_list field. Allowed value: `"community_list"`.
- `match_extended_community_list` (String, Optional) — Match Extended Community List. Default: `""`.
- `match_extended_community_list_ref_type_` (String, Optional) — Object type for match_extended_community_list field. Allowed value: `"extended_community_list"`.
- `match_interface_number` (Integer, Optional) — Match Interface Number. Default: `null`. Min: `1`. Max: `256`.
- `match_interface_vlan` (Integer, Optional) — Match Interface VLAN. Default: `null`. Min: `1`. Max: `4094`.
- `match_ipv4_address_ip_prefix_list` (String, Optional) — Match IPv4 Address IPv4 Prefix List. Default: `""`.
- `match_ipv4_address_ip_prefix_list_ref_type_` (String, Optional) — Object type for match_ipv4_address_ip_prefix_list field. Allowed value: `"ipv4_prefix_list"`.
- `match_ipv4_next_hop_ip_prefix_list` (String, Optional) — Match IPv4 Next Hop IPv4 Prefix List. Default: `""`.
- `match_ipv4_next_hop_ip_prefix_list_ref_type_` (String, Optional) — Object type for match_ipv4_next_hop_ip_prefix_list field. Allowed value: `"ipv4_prefix_list"`.
- `match_local_preference` (Integer, Optional) — Match BGP Local Preference value on the route. Default: `null`. Max: `4294967295`.
- `match_metric` (Integer, Optional) — Match Metric of the IP route entry. Default: `null`. Min: `1`. Max: `4294967295`.
- `match_origin` (String, Optional) — Match routes based on the value of the BGP Origin attribute. Allowed values: `""`, `"egp"`, `"igp"`, `"incomplete"`. Default: `""`.
- `match_peer_ip_address` (String, Optional) — Match BGP Peer IP Address the route was learned from. Default: `""`.
- `match_peer_interface` (Integer, Optional) — Match BGP Peer port the route was learned from. Default: `null`. Min: `1`. Max: `256`.
- `match_peer_vlan` (Integer, Optional) — Match BGP Peer VLAN over which the route was learned. Default: `null`. Min: `1`. Max: `4094`.
- `match_source_protocol` (String, Optional) — Match Routing Protocol the route originated from. Allowed values: `""`, `"bgp"`, `"connected"`, `"ospf"`, `"static"`. Default: `""`.
- `match_vrf` (String, Optional) — Match VRF the route is associated with. Default: `""`.
- `match_vrf_ref_type_` (String, Optional) — Object type for match_vrf field. Allowed value: `"tenant"`.
- `match_tag` (Integer, Optional) — Match routes that have this value for a Tag attribute. Default: `null`. Min: `1`. Max: `4294967295`.
- `match_evpn_route_type_default` (Boolean, Optional) — Match based on the type of EVPN Route Type being Default. Default: `null`.
- `match_evpn_route_type` (String, Optional) — Match based on the indicated EVPN Route Type. Allowed values: `""`, `"macip"`, `"multicast"`, `"prefix"`. Default: `""`.
- `match_vni` (Integer, Optional) — Match based on the VNI value. Default: `null`. Min: `1`. Max: `16777215`.
- `object_properties` (Block, Optional) —
  - `notes` (String, Optional) — User Notes. Default: `""`.
  - `match_fields_shown` (String, Optional) — Match fields shown. Default: `""`.
- `match_ipv6_address_ipv6_prefix_list` (String, Optional) — Match IPv4 Address IPv6 Prefix List. Default: `""`.
- `match_ipv6_address_ipv6_prefix_list_ref_type_` (String, Optional) — Object type for match_ipv6_address_ipv6_prefix_list field. Allowed value: `"ipv6_prefix_list"`.
- `match_ipv6_next_hop_ipv6_prefix_list` (String, Optional) — Match IPv6 Next Hop IPv6 Prefix List. Default: `""`.
- `match_ipv6_next_hop_ipv6_prefix_list_ref_type_` (String, Optional) — Object type for match_ipv6_next_hop_ipv6_prefix_list field. Allowed value: `"ipv6_prefix_list"`.

## Import

Route Map Clauses can be imported using the name:

```sh
terraform import verity_route_map_clause.<resource_name> <name>
```

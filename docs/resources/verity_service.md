# Service Resource

`verity_service` manages service resources in Verity, which define network service configurations.

## Example Usage

```hcl
resource "verity_service" "example" {
  name = "example"
  anycast_ipv4_mask = ""
  anycast_ipv6_mask = ""
  dhcp_server_ipv4 = ""
  dhcp_server_ipv6 = ""
  mtu = 1500
  tenant = ""
  tenant_ref_type_ = "tenant"
  policy_based_routing = ""
  policy_based_routing_ref_type_ = "pb_routing"
  vni_auto_assigned_ = true
  enable = false
  vlan = 4007
  vni = 104007
  max_upstream_rate_mbps = null
  max_downstream_rate_mbps = null
  packet_priority = "0"
  multicast_management_mode = "flooding"
  tagged_packets = false
  tls = false
  allow_local_switching = true
  act_as_multicast_querier = false
  block_unknown_unicast_flood = false
  block_downstream_dhcp_server = true
  is_management_service = false
  use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
  allow_fast_leave = false
  mst_instance = 0

  object_properties {
    group = ""
    on_summary = true
    warn_on_no_external_source = true
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `vlan` (Integer) - Layer 2 Virtual Network Identifier. A Value between 1 and 4096
* `vni` (Integer) - Identifies the service within the VXLAN fabric
* `vni_auto_assigned_` (Boolean) - Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit
* `tenant` (String) - Tenant
* `tenant_ref_type_` (String) - Object type for tenant field
* `anycast_ipv4_mask` (String) - Comma separated list of Static anycast gateway addresses(IPv4) for service
* `anycast_ipv6_mask` (String) - Comma separated list of Static anycast gateway addresses(IPv6) for service
* `dhcp_server_ipv4` (String) - IPv4 address(s) of the DHCP server for service. May have up to four separated by commas
* `dhcp_server_ipv6` (String) - IPv6 address(s) of the DHCP server for service. May have up to four separated by commas
* `mtu` (Integer) - MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets
* `object_properties` (Object) - Additional object properties
  * `group` (String) - Group
  * `on_summary` (Boolean) - Show on the summary view
  * `warn_on_no_external_source` (Boolean) - Warn if there is not outbound path for service in SD-Router or a Service Port Profile
* `policy_based_routing` (String) - Policy Based Routing
* `policy_based_routing_ref_type_` (String) - Object type for policy_based_routing field
* `max_upstream_rate_mbps` (Integer) - Bandwidth allocated per port in the upstream direction. (Max 10000 Mbps)
* `max_downstream_rate_mbps` (Integer) - Bandwidth allocated per port in the downstream direction. (Max 10000 Mbps)
* `packet_priority` (String) - Priority untagged packets will be tagged with on ingress to the network. If the network is flooded packets of lower priority will be dropped
* `multicast_management_mode` (String) - Determines how undefined handle multicast packet for Service
* `tagged_packets` (Boolean) - Overrides priority bits on incoming tagged packets. Always done for untagged packets
* `tls` (Boolean) - Is a Transparent LAN Service?
* `allow_local_switching` (Boolean) - Allow Edge Devices to communicate with each other. Disabling this forces upstream traffic to the router
* `act_as_multicast_querier` (Boolean) - Multicast managment through IGMP requires a multicast querier. Check this box if SD LAN should provide a multicast querier
* `block_unknown_unicast_flood` (Boolean) - Block unknown unicast traffic flooding and only permits egress traffic with MAC addresses that are known to exit on the port
* `block_downstream_dhcp_server` (Boolean) - Block inbound packets sent by Downstream DHCP servers
* `is_management_service` (Boolean) - Denotes a Management Service
* `use_dscp_to_p_bit_mapping_for_l3_packets_if_available` (Boolean) - use DSCP to p-bit Mapping for L3 packets if available
* `allow_fast_leave` (Boolean) - The Fast Leave feature causes the switch to immediately remove a port from the forwarding list for a IGMP multicast group when the port receives a leave message. Not recommended unless there is only a single receiver present on every point in the VLAN
* `mst_instance` (Integer) - MST Instance ID (0-4094)

## Import

Service resources can be imported using the `name` attribute:

```sh
terraform import verity_service.<resource_name> <name>
```
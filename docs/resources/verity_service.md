# verity_service (Resource)

Manages a Service resource.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_service" "example" {
  name = "example"
  act_as_multicast_querier = false
  allow_fast_leave = false
  allow_local_switching = false
  block_downstream_dhcp_server = false
  block_unknown_unicast_flood = false
  enable = false
  is_management_service = false
  max_downstream_rate_mbps = null
  max_upstream_rate_mbps = null
  mst_instance = null
  multicast_management_mode = ""
  packet_priority = ""
  tagged_packets = false
  tls = false
  use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
  vlan = null

  object_properties {
    warn_on_no_external_source = false
  }
}
```

### Datacenter mode

```hcl
resource "verity_service" "example" {
  name = "example"
  anycast_ipv4_mask = ""
  anycast_ipv6_mask = ""
  dhcp_server_ipv4 = ""
  dhcp_server_ipv6 = ""
  enable = false
  ip_attach_host_advertise = null
  mtu = null
  policy_based_routing = ""
  policy_based_routing_ref_type_ = "pb_routing"
  tenant = ""
  tenant_ref_type_ = "tenant"
  vlan = null
  vni_auto_assigned_ = true
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `act_as_multicast_querier` (Boolean) - Multicast managment through IGMP requires a multicast querier. Check this box if SD LAN should provide a multicast querier. Campus mode only.
* `allow_fast_leave` (Boolean) - The Fast Leave feature causes the switch to immediately remove a port from the forwarding list for a IGMP multicast group when the port receives a leave message. Not recommended unless there is only a single receiver present on every point in the VLAN. Campus mode only.
* `allow_local_switching` (Boolean) - Allow Edge Devices to communicate with each other. Disabling this forces upstream traffic to the router. Campus mode only.
* `anycast_ipv4_mask` (String) - Comma separated list of Static anycast gateway addresses(IPv4) for service. Datacenter mode only.
* `anycast_ipv6_mask` (String) - Comma separated list of Static anycast gateway addresses(IPv6) for service. Datacenter mode only.
* `block_downstream_dhcp_server` (Boolean) - Block inbound packets sent by Downstream DHCP servers. Campus mode only.
* `block_unknown_unicast_flood` (Boolean) - Block unknown unicast traffic flooding and only permits egress traffic with MAC addresses that are known to exit on the port. Campus mode only.
* `dhcp_server_ipv4` (String) - IPv4 address(s) of the DHCP server for service. May have up to four separated by commas. Datacenter mode only.
* `dhcp_server_ipv6` (String) - IPv6 address(s) of the DHCP server for service. May have up to four separated by commas. Datacenter mode only.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `ip_attach_host_advertise` (Integer) - Converts IP addresses learned via ARP and ND into host routes and adds them to the routing table. Sets the administrative distance of the attached route (range 0-250). Datacenter mode only. Set it to `null` to clear it.
* `is_management_service` (Boolean) - Denotes a Management Service. Campus mode only.
* `max_downstream_rate_mbps` (Integer) - Bandwidth allocated per port in the downstream direction. (Max 10000 Mbps). Campus mode only. Set it to `null` to clear it.
* `max_upstream_rate_mbps` (Integer) - Bandwidth allocated per port in the upstream direction. (Max 10000 Mbps). Campus mode only. Set it to `null` to clear it.
* `mst_instance` (Integer) - MST Instance ID (0-4094). Campus mode only. Set it to `null` to clear it.
* `mtu` (Integer) - MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. Datacenter mode only. Set it to `null` to clear it.
* `multicast_management_mode` (String) - Determines how undefined handle multicast packet for Service "Multicast Flooding (Normal)" Multicast packets are broadcast; "Multicast Flooding (AVB/PTP/Cobranet)" Multicast packets are broadcast with special treatment for critical latency packets such as used by AVB, PTP, and Cobranet; "IPTV Filtering (IGMP Snooping)" Multicast packets are propagated via IGMP Snooping; "IPTV Filtering (IGMP Report/Leave Flooding)" Multicast packets are propagated via IGMP Snooping. except that IGMP Report/Leave packets are broadcast. Campus mode only.
* `object_properties` (Block) - Object properties for the service. At most one block.
  * `warn_on_no_external_source` (Boolean) - Warn if there is not outbound path for service in SD-Router or a Service Port Profile. Campus mode only.
* `packet_priority` (String) - Priority untagged packets will be tagged with on ingress to the network. If the network is flooded packets of lower priority will be dropped. Campus mode only.
* `policy_based_routing` (String) - Policy Based Routing. Datacenter mode only. Set together with `policy_based_routing_ref_type_`.
* `policy_based_routing_ref_type_` (String) - Object type for policy_based_routing field. Datacenter mode only.
* `tagged_packets` (Boolean) - Overrides priority bits on incoming tagged packets. Always done for untagged packets. Campus mode only.
* `tenant` (String) - Tenant. Datacenter mode only. Set together with `tenant_ref_type_`.
* `tenant_ref_type_` (String) - Object type for tenant field. Datacenter mode only.
* `tls` (Boolean) - Is a Transparent LAN Service? Campus mode only.
* `use_dscp_to_p_bit_mapping_for_l3_packets_if_available` (Boolean) - use DSCP to p-bit Mapping for L3 packets if available. Campus mode only.
* `vlan` (Integer) - Layer 2 Virtual Network Identifier. A Value between 1 and 4096. Some switches have reserved values within the range. Set it to `null` to clear it.
* `vni` (Integer) - Identifies the service within the VXLAN fabric - Range is 1-16777215. If not using auto, must be outside of the reserved range settings. Datacenter mode only. Set it to `null` to clear it. Assigned by the server while `vni_auto_assigned_` is `true`; it cannot be set then.
* `vni_auto_assigned_` (Boolean) - Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit. Datacenter mode only.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `policy_based_routing` | `policy_based_routing_ref_type_` | `pb_routing` |
| `tenant` | `tenant_ref_type_` | `tenant` |

## Auto-Assigned Fields

While a flag is `true`, the server assigns the paired value and the value cannot be set in the configuration. Set the flag to `false` to set the value yourself; if you leave the value unset, the value the server assigned is kept.

| Field | Flag | Reassigned when |
| --- | --- | --- |
| `vni` | `vni_auto_assigned_` | `vlan` changes |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_service.<resource_name> <name>
```

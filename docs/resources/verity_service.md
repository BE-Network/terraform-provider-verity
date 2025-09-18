# Service Resource

`verity_service` manages service resources in Verity, which define network service configurations.

## Example Usage

```hcl
resource "verity_service" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  anycast_ip_mask = ""
  dhcp_server_ip = ""
  mtu = 1500
  tenant = ""
  tenant_ref_type_ = ""
  vni_auto_assigned_ = true
  enable = false
  vlan = 4007
  vni = 104007
}
```

## Argument Reference

* `name` - Unique identifier for the service
* `enable` - Enable this service. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
  * `on_summary` - Show on the summary view (available as of API version 6.5)
  * `warn_on_no_external_source` - Warn if there is not outbound path for service in SD-Router or a Service Port Profile (available as of API version 6.5)
* `vlan` - VLAN ID for the service
* `vni` - VNI (VXLAN Network Identifier) for the service
* `vni_auto_assigned_` - Whether the VNI value is automatically assigned
* `tenant` - Reference to a tenant resource
* `tenant_ref_type_` - Object type for tenant reference
* `anycast_ip_mask` - Static anycast gateway address for service
* `dhcp_server_ip` - IP address(es) of the DHCP server for service. May have up to four separated by commas.
* `dhcp_server_ipv4` - IPv4 address(es) of the DHCP server for service (available as of API version 6.5)
* `dhcp_server_ipv6` - IPv6 address(es) of the DHCP server for service (available as of API version 6.5)
* `mtu` - MTU (Maximum Transmission Unit) - the size used by a switch to determine when large packets must be broken up for delivery
* `anycast_ipv4_mask` - Static anycast gateway addresses (IPv4) for service (available as of API version 6.5)
* `anycast_ipv6_mask` - Static anycast gateway addresses (IPv6) for service (available as of API version 6.5)
* `max_upstream_rate_mbps` - Bandwidth allocated per port in the upstream direction (available as of API version 6.5)
* `max_downstream_rate_mbps` - Bandwidth allocated per port in the downstream direction (available as of API version 6.5)
* `packet_priority` - Priority untagged packets will be tagged with on ingress to the network (available as of API version 6.5)
* `multicast_management_mode` - Determines how to handle multicast packets for Service (available as of API version 6.5)
* `tagged_packets` - Overrides priority bits on incoming tagged packets (available as of API version 6.5)
* `tls` - Is a Transparent LAN Service (available as of API version 6.5)
* `allow_local_switching` - Allow Edge Devices to communicate with each other (available as of API version 6.5)
* `act_as_multicast_querier` - Multicast management through IGMP requires a multicast querier (available as of API version 6.5)
* `block_unknown_unicast_flood` - Block unknown unicast traffic flooding (available as of API version 6.5)
* `block_downstream_dhcp_server` - Block inbound packets sent by Downstream DHCP servers (available as of API version 6.5)
* `is_management_service` - Denotes a Management Service (available as of API version 6.5)
* `use_dscp_to_p_bit_mapping_for_l3_packets_if_available` - Use DSCP to p-bit Mapping for L3 packets if available (available as of API version 6.5)
* `allow_fast_leave` - The Fast Leave feature causes the switch to immediately remove a port from the forwarding list (available as of API version 6.5)
* `mst_instance` - MST Instance ID (0-4094) (available as of API version 6.5)

## Import

Service resources can be imported using the `name` attribute:

```sh
terraform import verity_service.<resource_name> <name>
```
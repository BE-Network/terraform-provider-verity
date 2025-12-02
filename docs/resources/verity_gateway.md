# Gateway Resource

`verity_gateway` manages gateway resources in Verity, which define network gateway configurations.

## Example Usage

```hcl
resource "verity_gateway" "example" {
  name = "example"
  enable = false
  tenant = ""
  tenant_ref_type_ = ""
  neighbor_ip_address = ""
  neighbor_as_number = null
  fabric_interconnect = false
  keepalive_timer = 60
  hold_timer = 180
  connect_timer = 120
  advertisement_interval = 30
  ebgp_multihop = 255
  egress_vlan = null
  source_ip_address = ""
  anycast_ip_mask = ""
  md5_password = ""
  import_route_map = ""
  import_route_map_ref_type_ = ""
  export_route_map = ""
  export_route_map_ref_type_ = ""
  gateway_mode = "Static BGP"
  local_as_number = null
  local_as_no_prepend = false
  replace_as = false
  max_local_as_occurrences = 0
  dynamic_bgp_subnet = ""
  dynamic_bgp_limits = 0
  helper_hop_ip_address = ""
  enable_bfd = false
  bfd_receive_interval = 300
  bfd_transmission_interval = 300
  bfd_detect_multiplier = 3
  next_hop_self = false
  default_originate = false
  bfd_multihop = false

  static_routes {
    index = 1
    enable = false
    ipv4_route_prefix = ""
    next_hop_ip_address = ""
    ad_value = null
  }

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `tenant` (String) - Tenant.
* `tenant_ref_type_` (String) - Object type for tenant field.
* `neighbor_ip_address` (String) - IP address of remote BGP peer.
* `neighbor_as_number` (Integer) - Autonomous System Number of remote BGP peer.
* `fabric_interconnect` (Boolean) - Fabric interconnect setting.
* `keepalive_timer` (Integer) - Interval in seconds between Keepalive messages sent to remote BGP peer.
* `hold_timer` (Integer) - Time, in seconds, used to determine failure of session Keepalive messages received from remote BGP peer.
* `connect_timer` (Integer) - Time in seconds between sucessive attempts to Establish BGP session.
* `advertisement_interval` (Integer) - The minimum time in seconds between sending route updates to BGP neighbor.
* `ebgp_multihop` (Integer) - Allows external BGP neighbors to establish peering session multiple network hops away.
* `egress_vlan` (Integer) - VLAN used to carry BGP TCP session.
* `source_ip_address` (String) - Source IP address used to override the default source address calculation for BGP TCP session.
* `anycast_ip_mask` (String) - The Anycast Address will be used to enable an IP routing redundancy mechanism designed to allow for transparent failover across a leaf pair at the first-hop IP router.
* `md5_password` (String) - MD5 password.
* `import_route_map` (String) - A route-map applied to routes imported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes.
* `import_route_map_ref_type_` (String) - Object type for import_route_map field.
* `export_route_map` (String) - A route-map applied to routes exported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes.
* `export_route_map_ref_type_` (String) - Object type for export_route_map field.
* `gateway_mode` (String) - Gateway Mode. Can be BGP, Static, or Default.
* `local_as_number` (Integer) - Local AS Number.
* `local_as_no_prepend` (Boolean) - Do not prepend the local-as number to the AS-PATH for routes advertised through this BGP gateway. The Local AS Number must be set for this to be able to be set.
* `replace_as` (Boolean) - Prepend only Local AS in updates to EBGP peers.
* `max_local_as_occurrences` (Integer) - Allow routes with the local AS number in the AS-path, specifying the maximum occurrences permitted before declaring a routing loop. Leave blank or '0' to disable.
* `dynamic_bgp_subnet` (String) - Dynamic BGP Subnet.
* `dynamic_bgp_limits` (Integer) - Dynamic BGP Limits.
* `helper_hop_ip_address` (String) - Helper Hop IP Address.
* `enable_bfd` (Boolean) - Enable BFD(Bi-Directional Forwarding).
* `bfd_receive_interval` (Integer) - Configure the minimum interval during which the system can receive BFD control packets.
* `bfd_transmission_interval` (Integer) - Configure the minimum transmission interval during which the system can send BFD control packets.
* `bfd_detect_multiplier` (Integer) - Configure the detection multiplier to determine packet loss.
* `next_hop_self` (Boolean) - Optional attribute that disables the normal BGP calculation of next-hops for advertised routes and instead sets the next-hops for advertised routes to the IP address of the switch itself.
* `default_originate` (Boolean) - Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor.
* `bfd_multihop` (Boolean) - Enable BFD Multi-Hop for Neighbor. This is used to detect failures in the forwarding path between the BGP peers.
* `static_routes` (Array) - 
  * `enable` (Boolean) - Enable of this static route.
  * `ipv4_route_prefix` (String) - IPv4 unicast IP address followed by a subnet mask length.
  * `next_hop_ip_address` (String) - Next Hop IP Address. Must be a unicast IP address.
  * `ad_value` (Integer) - Administrative distancing value, also known as route preference - values from 0-255.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - Additional object properties.
  * `group` (String) - Group.

## Import

Gateway resources can be imported using the `name` attribute:

```sh
terraform import verity_gateway.<resource_name> <name>
```

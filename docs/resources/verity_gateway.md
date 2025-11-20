# Gateway Resource

`verity_gateway` manages gateway resources in Verity, which define network gateway configurations.

## Example Usage

```hcl
resource "verity_gateway" "example" {
  name = "example"
  tenant = ""
  md5_password = ""
  local_as_no_prepend = false
  dynamic_bgp_subnet = ""
  import_route_map_ref_type_ = ""
  connect_timer = 120
  anycast_ip_mask = ""
  helper_hop_ip_address = ""
  export_route_map_ref_type_ = ""
  neighbor_as_number = 515
  bfd_transmission_interval = 366
  bfd_detect_multiplier = 3
  next_hop_self = false
  tenant_ref_type_ = ""
  ebgp_multihop = 255
  max_local_as_occurrences = 0
  enable_bfd = false
  source_ip_address = ""
  local_as_number = 351
  dynamic_bgp_limits = 0
  keepalive_timer = 60
  hold_timer = 180
  enable = false
  advertisement_interval = 30
  egress_vlan = null
  import_route_map = ""
  gateway_mode = "Dynamic BGP"
  bfd_receive_interval = 307
  bfd_multihop = false
  fabric_interconnect = false
  export_route_map = ""
  replace_as = false
  neighbor_ip_address = ""

  static_routes {
    index = 1
    ad_value = 2
    enable = false
    ipv4_route_prefix = ""
    next_hop_ip_address = ""
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
* `fabric_interconnect` (Boolean) - .
* `keepalive_timer` (Integer) - Interval in seconds between Keepalive messages sent to remote BGP peer.
* `hold_timer` (Integer) - Time, in seconds, used to determine failure of session Keepalive messages received from remote BGP peer.
* `connect_timer` (Integer) - Time in seconds between sucessive attempts to Establish BGP session.
* `advertisement_interval` (Integer) - The minimum time in seconds between sending route updates to BGP neighbor.
* `ebgp_multihop` (Integer) - Allows external BGP neighbors to establish peering session multiple network hops away.
* `egress_vlan` (Integer) - VLAN used to carry BGP TCP session.
* `source_ip_address` (String) - Source IP address used to override the default source address calculation for BGP TCP session.
* `anycast_ip_mask` (String) - The Anycast Address can be used to enable an IP routing redundancy mechanism designed to allow for transparent failover across a leaf pair at the first-hop IP router.
* `md5_password` (String) - MD5 Password used in the BGP session.
* `import_route_map` (String) - A Route Map applied to routes imported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes.
* `import_route_map_ref_type_` (String) - Object type for import_route_map field.
* `export_route_map` (String) - A route-map applied to routes exported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes.
* `export_route_map_ref_type_` (String) - Object type for export_route_map field.
* `gateway_mode` (String) - Gateway Mode is the method used for defining routes for the Tenant.
* `local_as_number` (Integer) - Local AS Number to use as an override to switch AS number.
* `local_as_no_prepend` (Boolean) - Do not prepend the local-as number to the AS-PATH for routes advertised through this BGP gateway. The Local AS Number must be set for this to be able to be set.
* `replace_as` (Boolean) - Prepend only Local AS in updates to EBGP peers.
* `max_local_as_occurrences` (Integer) - Allow routes with the local AS number in the AS-path, specifying the maximum occurrences permitted before declaring a routing loop. Leave blank or '0' to disable.
* `dynamic_bgp_subnet` (String) - Dynamic BGP Subnet.
* `dynamic_bgp_limits` (Integer) - Dynamic BGP Limits.
* `helper_hop_ip_address` (String) - Neighbor Next Hop IP Address is used as the next hop to reach the BGP peer in the case it is not a direct connection.
* `enable_bfd` (Boolean) - Enable BFD(Bi-Directional Forwarding).
* `bfd_receive_interval` (Integer) - Configure the minimum interval during which the system can receive BFD control packets.
* `bfd_transmission_interval` (Integer) - Configure the minimum transmission interval during which the system can send BFD control packets.
* `bfd_detect_multiplier` (Integer) - Configure the detection multiplier to determine packet loss.
* `bfd_multihop` (Boolean) - Enable BFD Multi-Hop for Neighbor. This is used to detect failures in the forwarding path between the BGP peers.
* `next_hop_self` (Boolean) - Optional attribute that disables the normal BGP calculation of next-hops for advertised routes and instead sets the next-hops for advertised routes to the IP address of the switch itself.
* `static_routes` (Array) - 
  * `enable` (Boolean) - Enable of this static route.
  * `ipv4_route_prefix` (String) - IPv4 unicast IP address followed by a subnet mask length.
  * `next_hop_ip_address` (String) - Next Hop IP Address. Must be a unicast IP address.
  * `ad_value` (Integer) - Administrative distancing value, also known as route preference - values from 0-255.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.
* `switch_encrypted_md5_password` (Boolean) - Indicates the entered password is a switch encrypted password.
* `md5_password_encrypted` (String) - MD5 Password Encrypted used in the BGP session.
* `default_originate` (Boolean) - Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor.

## Import

Gateway resources can be imported using the `name` attribute:

```sh
terraform import verity_gateway.<resource_name> <name>
```

# verity_gateway (Resource)

Manages a Verity Gateway.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_gateway" "example" {
  name = "example"
  advertisement_interval = null
  allowas_in_origin = false
  anycast_ip_mask = ""
  bfd_detect_multiplier = null
  bfd_multihop = false
  bfd_receive_interval = null
  bfd_transmission_interval = null
  bgp_instance_as_number = null
  connect_timer = null
  default_originate = false
  dynamic_bgp_limits = null
  dynamic_bgp_subnet = ""
  ebgp_multihop = null
  egress_vlan = null
  enable = false
  enable_bfd = false
  export_route_map = ""
  export_route_map_ref_type_ = "route_map"
  fabric = ""
  fabric_interconnect = false
  fabric_ref_type_ = "fabric"
  gateway_mode = ""
  helper_hop_ip_address = ""
  hold_timer = null
  import_route_map = ""
  import_route_map_ref_type_ = "route_map"
  keepalive_timer = null
  local_as_no_prepend = false
  local_as_number = null
  max_local_as_occurrences = null
  md5_password = ""
  md5_password_encrypted = ""
  neighbor_as_number = null
  neighbor_ip_address = ""
  next_hop_self = false
  remove_private_as = false
  replace_as = false
  source_ip_address = ""
  switch_encrypted_md5_password = false
  tenant = ""
  tenant_ref_type_ = "tenant"
  type = ""

  static_routes {
    index = 1
    ad_value = null
    enable = false
    ipv4_route_prefix = ""
    next_hop_ip_address = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object Name. Must be unique. Changing it replaces the resource.

### Optional

* `advertisement_interval` (Integer) - The minimum time in seconds between sending route updates to BGP neighbor. Set it to `null` to clear it.
* `allowas_in_origin` (Boolean) - Only accept the current AS in the as-path if the route was originated in the Local AS.
* `anycast_ip_mask` (String) - The Anycast Address can be used to enable an IP routing redundancy mechanism designed to allow for transparent failover across a leaf pair at the first-hop IP router.
* `bfd_detect_multiplier` (Integer) - Configure the detection multiplier to determine packet loss. Set it to `null` to clear it.
* `bfd_multihop` (Boolean) - Enable BFD Multi-Hop for Neighbor. This is used to detect failures in the forwarding path between the BGP peers.
* `bfd_receive_interval` (Integer) - Configure the minimum interval during which the system can receive BFD control packets. Set it to `null` to clear it.
* `bfd_transmission_interval` (Integer) - Configure the minimum transmission interval during which the system can send BFD control packets. Set it to `null` to clear it.
* `bgp_instance_as_number` (Integer) - Override the switch's AS number used in the Tenant router definition where this Gateway is applied. Set it to `null` to clear it.
* `connect_timer` (Integer) - Time in seconds between successive attempts to establish BGP session. Set it to `null` to clear it.
* `default_originate` (Boolean) - Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor.
* `dynamic_bgp_limits` (Integer) - Dynamic BGP Limits. Set it to `null` to clear it.
* `dynamic_bgp_subnet` (String) - Dynamic BGP Subnet.
* `ebgp_multihop` (Integer) - Allows external BGP neighbors to establish peering session multiple network hops away. Set it to `null` to clear it.
* `egress_vlan` (Integer) - VLAN used to carry BGP TCP session. Set it to `null` to clear it.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `enable_bfd` (Boolean) - Enable BFD (Bi-Directional Forwarding).
* `export_route_map` (String) - A route-map applied to routes exported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes. Set together with `export_route_map_ref_type_`.
* `export_route_map_ref_type_` (String) - Object type for export_route_map field.
* `fabric` (String) - Fabric this Gateway is assigned to. Set together with `fabric_ref_type_`.
* `fabric_interconnect` (Boolean)
* `fabric_ref_type_` (String) - Object type for fabric field.
* `gateway_mode` (String) - Gateway Mode is the method used for defining routes for the Tenant.
* `helper_hop_ip_address` (String) - Neighbor Next Hop IP Address is used as the next hop to reach the BGP peer in the case it is not a direct connection.
* `hold_timer` (Integer) - Time, in seconds, used to determine failure of session Keepalive messages received from remote BGP peer. Set it to `null` to clear it.
* `import_route_map` (String) - A Route Map applied to routes imported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes. Set together with `import_route_map_ref_type_`.
* `import_route_map_ref_type_` (String) - Object type for import_route_map field.
* `keepalive_timer` (Integer) - Interval in seconds between Keepalive messages sent to remote BGP peer. Set it to `null` to clear it.
* `local_as_no_prepend` (Boolean) - Do not prepend the local-as number to the AS-PATH for routes advertised through this BGP gateway. The Local AS Number must be set for this to be able to be set.
* `local_as_number` (Integer) - Local AS Number to use as an override to switch AS number. Set it to `null` to clear it.
* `max_local_as_occurrences` (Integer) - Allow routes with the local AS number in the AS-path, specifying the maximum occurrences permitted before declaring a routing loop. Leave blank or '0' to disable. Set it to `null` to clear it.
* `md5_password` (String) - MD5 Password used in the BGP session.
* `md5_password_encrypted` (String) - MD5 Password Encrypted used in the BGP session.
* `neighbor_as_number` (Integer) - Autonomous System Number of remote BGP peer. Set it to `null` to clear it.
* `neighbor_ip_address` (String) - IP address of remote BGP peer.
* `next_hop_self` (Boolean) - Optional attribute that disables the normal BGP calculation of next-hops for advertised routes and instead sets the next-hops for advertised routes to the IP address of the switch itself.
* `remove_private_as` (Boolean) - Remove all private AS numbers from AS-PATH attributes for routes advertised through this BGP gateway.
* `replace_as` (Boolean) - Prepend only Local AS in updates to EBGP peers.
* `source_ip_address` (String) - Source IP address used to override the default source address calculation for BGP TCP session.
* `static_routes` (Block List) - List of static routes. Entries are matched by `index`.
  * `ad_value` (Integer) - Administrative distancing value, also known as route preference - values from 0-255. Set it to `null` to clear it.
  * `enable` (Boolean) - Enable of this static route.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `ipv4_route_prefix` (String) - IPv4 unicast IP address followed by a subnet mask length.
  * `next_hop_ip_address` (String) - Next Hop IP Address. Must be a unicast IP address.
* `switch_encrypted_md5_password` (Boolean) - Indicates the entered password is a switch encrypted password.
* `tenant` (String) - Tenant. Set together with `tenant_ref_type_`.
* `tenant_ref_type_` (String) - Object type for tenant field.
* `type` (String) - Gateway classification.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `export_route_map` | `export_route_map_ref_type_` | `route_map` |
| `fabric` | `fabric_ref_type_` | `fabric` |
| `import_route_map` | `import_route_map_ref_type_` | `route_map` |
| `tenant` | `tenant_ref_type_` | `tenant` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_gateway.<resource_name> <name>
```

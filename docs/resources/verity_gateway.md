# Gateway Resource

`verity_gateway` manages gateway resources in Verity, which define network gateway configurations.

## Example Usage

```hcl
resource "verity_gateway" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  tenant = ""
  md5_password = ""
  local_as_no_prepend = false
  dynamic_bgp_subnet = ""
  import_route_map_ref_type_ = ""
  connect_timer = 120
  anycast_ip_mask = ""
  helper_hop_ip_address = ""
  static_routes {
    index = 1
    ad_value = 2
    enable = false
    ipv4_route_prefix = ""
    next_hop_ip_address = ""
  }
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
}
```

## Argument Reference

* `name` - Unique identifier for the gateway
* `enable` - Enable this gateway. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
* `tenant` - Reference to a tenant resource
* `tenant_ref_type_` - Object type for tenant reference
* `neighbor_ip_address` - IP address of remote BGP peer
* `neighbor_as_number` - AS number of remote BGP peer
* `fabric_interconnect` - Whether this is a fabric interconnect
* `keepalive_timer` - BGP keepalive timer in seconds
* `hold_timer` - BGP hold timer in seconds
* `connect_timer` - BGP connect timer in seconds
* `advertisement_interval` - BGP advertisement interval in seconds
* `ebgp_multihop` - EBGP multihop count
* `static_routes` - List of static route blocks
  * `enable` - Enable this static route
  * `ipv4_route_prefix` - IPv4 route prefix in CIDR notation
  * `next_hop_ip_address` - Next hop IP address
  * `ad_value` - Administrative distance value (0-255)
  * `index` - Index identifying the object
* `egress_vlan` - BGP egress VLAN
* `source_ip_address` - Source IP address (available as of API version 6.5)
* `anycast_ip_mask` - Anycast IP mask (available as of API version 6.5)
* `md5_password` - MD5 password for BGP authentication (available as of API version 6.5)
* `import_route_map` - Import route map reference (available as of API version 6.5)
* `import_route_map_ref_type_` - Object type for import route map reference (available as of API version 6.5)
* `export_route_map` - Export route map reference (available as of API version 6.5)
* `export_route_map_ref_type_` - Object type for export route map reference (available as of API version 6.5)
* `gateway_mode` - Gateway mode (e.g., "Dynamic BGP") (available as of API version 6.5)
* `local_as_number` - Local AS number (available as of API version 6.5)
* `local_as_no_prepend` - Whether to not prepend local AS number (available as of API version 6.5)
* `replace_as` - Whether to replace AS (available as of API version 6.5)
* `max_local_as_occurrences` - Maximum local AS occurrences (available as of API version 6.5)
* `dynamic_bgp_subnet` - Dynamic BGP subnet (available as of API version 6.5)
* `dynamic_bgp_limits` - Dynamic BGP limits (available as of API version 6.5)
* `helper_hop_ip_address` - Helper hop IP address (available as of API version 6.5)
* `enable_bfd` - Enable BFD (Bidirectional Forwarding Detection) (available as of API version 6.5)
* `bfd_receive_interval` - BFD receive interval (available as of API version 6.5)
* `bfd_transmission_interval` - BFD transmission interval (available as of API version 6.5)
* `bfd_detect_multiplier` - BFD detect multiplier (available as of API version 6.5)
* `bfd_multihop` - BFD multihop (available as of API version 6.5)
* `next_hop_self` - Use next hop self (available as of API version 6.5)
* `default_originate` - Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor (available as of API version 6.5)

## Import

Gateway resources can be imported using the `name` attribute:

```sh
terraform import verity_gateway.<resource_name> <name>
```

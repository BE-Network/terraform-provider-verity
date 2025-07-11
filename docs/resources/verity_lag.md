# Link Aggregation Group (LAG) Resource

`verity_lag` manages Link Aggregation Groups in Verity, which combine multiple network connections for increased throughput and redundancy.

## Example Usage

```hcl
resource "verity_lag" "example" {
  name = "example"
  object_properties {}
  is_peer_link = false
  peer_link_vlan = null
  fallback = true
  fast_rate = false
  enable = true
  color = "chardonnay"
  lacp = true
  eth_port_profile = ""
  uplink = false
  eth_port_profile_ref_type_ = ""
}
```

## Argument Reference

* `name` - Unique identifier for the LAG
* `enable` - Enable this LAG. Default is `false`
* `object_properties` - Object properties block
* `is_peer_link` - Whether this LAG is a peer link
* `color` - Color identifier for visual representation
* `lacp` - Enable LACP (Link Aggregation Control Protocol)
* `eth_port_profile` - Reference to an Ethernet port profile
* `eth_port_profile_ref_type_` - Object type for Ethernet port profile reference
* `peer_link_vlan` - VLAN ID used for peer link
* `fallback` - Enable fallback mode
* `fast_rate` - Enable fast rate transmission
* `uplink` - Whether this LAG is an uplink

## Import

LAG resources can be imported using the `name` attribute:

```
$ terraform import verity_lag.example example
```
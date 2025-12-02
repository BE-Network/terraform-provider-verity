# Link Aggregation Group (LAG) Resource

`verity_lag` manages Link Aggregation Groups in Verity, which combine multiple network connections for increased throughput and redundancy.

## Example Usage

```hcl
resource "verity_lag" "example" {
  name = "example"
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

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `is_peer_link` (Boolean) - Indicates this LAG is used for peer-to-peer Peer-LAG/IDS link.
* `color` (String) - Choose the color to display the connectors on the network view.
* `lacp` (Boolean) - LACP.
* `eth_port_profile` (String) - Choose an Eth Port Profile.
* `eth_port_profile_ref_type_` (String) - Object type for eth_port_profile field.
* `peer_link_vlan` (Integer) - For peer-peer LAGs. The VLAN used for control.
* `fallback` (Boolean) - Allows an active member interface to establish a connection with a peer interface before the port channel receives the LACP protocol negotiation from the peer.
* `fast_rate` (Boolean) - Send LACP packets every second (if disabled, packets are sent every 30 seconds).
* `object_properties` (Object) - Additional object properties.
* `uplink` (Boolean) - Indicates this LAG is designated as an uplink in the case of a spineless pod. Link State Tracking will be applied to BGP Egress VLANs/Interfaces and the MCLAG Peer Link VLAN.

## Import

LAG resources can be imported using the `name` attribute:

```sh
terraform import verity_lag.<resource_name> <name>
```
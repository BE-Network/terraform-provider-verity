# verity_lag (Resource)

Manages a Link Aggregation Group (LAG).

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_lag" "example" {
  name = "example"
  color = ""
  crc_failure_threshold = null
  enable = false
  eth_port_profile = ""
  eth_port_profile_ref_type_ = ""
  fallback = false
  fast_rate = false
  is_peer_link = false
  lacp = false
  peer_link_vlan = null
  uplink = false

  object_properties {
    fabric = ""
    fabric_ref_type_ = "fabric"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `color` (String) - Choose the color to display the connectors on the network view.
* `crc_failure_threshold` (Integer) - Threshold in Errors per second that when met will disable this LAG's links. Set it to `null` to clear it.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `eth_port_profile` (String) - Choose an Eth Port Profile. Set together with `eth_port_profile_ref_type_`.
* `eth_port_profile_ref_type_` (String) - Object type for eth_port_profile field.
* `fallback` (Boolean) - Allows an active member interface to establish a connection with a peer interface before the port channel receives the LACP protocol negotiation from the peer.
* `fast_rate` (Boolean) - Send LACP packets every second (if disabled, packets are sent every 30 seconds).
* `is_peer_link` (Boolean) - Indicates this LAG is used for peer-to-peer Peer-LAG/IDS link.
* `lacp` (Boolean) - LACP.
* `object_properties` (Block) - Object properties. At most one block.
  * `fabric` (String) - Choose a Fabric. Set together with `fabric_ref_type_`.
  * `fabric_ref_type_` (String) - Object type for fabric field.
* `peer_link_vlan` (Integer) - For peer-peer LAGs. The VLAN used for control. Set it to `null` to clear it.
* `uplink` (Boolean) - Indicates this LAG is designated as an uplink in the case of a spineless pod. Link State Tracking will be applied to BGP Egress VLANs/Interfaces and the MCLAG Peer Link VLAN.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `eth_port_profile` | `eth_port_profile_ref_type_` | `eth_port_profile_`, `pb_egress_profile`, `service_port_profile` |
| `object_properties.fabric` | `fabric_ref_type_` | `fabric` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_lag.<resource_name> <name>
```

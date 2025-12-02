# Ethernet Port Settings Resource

`verity_eth_port_settings` manages Ethernet port settings in Verity, which define specific configurations for physical Ethernet ports.

## Example Usage

```hcl
resource "verity_eth_port_settings" "example" {
  name = "example"
  enable = false
  auto_negotiation = true
  max_bit_rate = "-1"
  duplex_mode = "Auto"
  stp_enable = false
  fast_learning_mode = true
  bpdu_guard = false
  bpdu_filter = false
  guard_loop = false
  poe_enable = false
  priority = "High"
  allocated_power = "0.0"
  bsp_enable = false
  broadcast = true
  multicast = true
  max_allowed_value = 1000
  max_allowed_unit = "pps"
  action = "Protect"
  fec = "unaltered"
  single_link = false

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String, Required) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `auto_negotiation` (Boolean) - Indicates if port speed and duplex mode should be auto negotiated.
* `max_bit_rate` (String) - Maximum Bit Rate allowed.
* `duplex_mode` (String) - Duplex Mode.
* `stp_enable` (Boolean) - Enable Spanning Tree on the port. Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings.
* `fast_learning_mode` (Boolean) - Enable Immediate Transition to Forwarding.
* `bpdu_guard` (Boolean) - Block port on BPDU Receive.
* `bpdu_filter` (Boolean) - Drop all Rx and Tx BPDUs.
* `guard_loop` (Boolean) - Enable Cisco Guard Loop.
* `poe_enable` (Boolean) - PoE Enable.
* `priority` (String) - Priority given when assigning power in a limited power situation.
* `allocated_power` (String) - Power the PoE system will attempt to allocate on this port.
* `bsp_enable` (Boolean) - Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU.
* `broadcast` (Boolean) - Broadcast.
* `multicast` (Boolean) - Multicast.
* `max_allowed_value` (Integer, Nullable) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action. Range: 1-9999.
* `max_allowed_unit` (String) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action.
* `action` (String) - Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max.
* `fec` (String) - FEC is Forward Error Correction which is error correction on the fiber link.
* `single_link` (Boolean) - Ports with this setting will be disabled when link state tracking takes effect.
* `object_properties` (Object) - Additional object properties.
  * `group` (String) - Group.

## Import

Ethernet port settings resources can be imported using the `name` attribute:

```sh
terraform import verity_eth_port_settings.<resource_name> <name>
```
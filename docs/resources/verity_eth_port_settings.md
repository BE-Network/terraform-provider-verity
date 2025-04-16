# Ethernet Port Settings Resource

`verity_eth_port_settings` manages Ethernet port settings in Verity, which define specific configurations for physical Ethernet ports.

## Example Usage

```hcl
resource "verity_eth_port_settings" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  max_allowed_value = 1000
  fec = "unaltered"
  single_link = true
  allocated_power = "0.0"
  bsp_enable = false
  fast_learning_mode = true
  bpdu_guard = false
  enable = true
  duplex_mode = "Auto"
  stp_enable = false
  bpdu_filter = false
  poe_enable = false
  priority = "High"
  multicast = true
  max_allowed_unit = "pps"
  auto_negotiation = true
  max_bit_rate = "-1"
  action = "Protect"
  broadcast = true
  guard_loop = false
}
```

## Argument Reference

* `name` - Unique identifier for the Ethernet port settings
* `enable` - Enable these settings. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
* `auto_negotiation` - Enable auto-negotiation
* `max_bit_rate` - Maximum bit rate (e.g., "10G", "25G", "100G")
* `duplex_mode` - Duplex mode ("full", "half")
* `stp_enable` - Enable Spanning Tree Protocol
* `fast_learning_mode` - Enable fast learning mode
* `bpdu_guard` - Enable BPDU guard
* `bpdu_filter` - Enable BPDU filter
* `guard_loop` - Enable loop guard
* `poe_enable` - Enable Power over Ethernet
* `priority` - Priority level ("low", "medium", "high", "critical")
* `allocated_power` - Power allocation value
* `bsp_enable` - Enable BSP
* `broadcast` - Enable broadcast
* `multicast` - Enable multicast
* `max_allowed_value` - Maximum allowed value for storm control
* `max_allowed_unit` - Unit for maximum allowed value
* `action` - Action to take when maximum is exceeded
* `fec` - Forward Error Correction setting
* `single_link` - Whether this is a single link

## Import

Ethernet port settings resources can be imported using the `name` attribute:


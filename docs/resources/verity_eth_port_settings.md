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
  * `overridden_object` - Overridden object (available as of API version 6.5)
  * `overridden_object_ref_type_` - Object type for overridden_object field (available as of API version 6.5)
  * `isdefault` - Default object (available as of API version 6.5)
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
* `minimum_wred_threshold` - A value between 1 to 12480(in KiloBytes) (available as of API version 6.5)
* `maximum_wred_threshold` - A value between 1 to 12480(in KiloBytes) (available as of API version 6.5)
* `wred_drop_probability` - A value between 0 to 100 (available as of API version 6.5)
* `priority_flow_control_watchdog_action` - Action when link state tracking takes effect (available as of API version 6.5)
* `priority_flow_control_watchdog_detect_time` - A value between 100 to 5000 (available as of API version 6.5)
* `priority_flow_control_watchdog_restore_time` - A value between 100 to 60000 (available as of API version 6.5)
* `packet_queue` - Packet Queue reference (available as of API version 6.5)
* `packet_queue_ref_type_` - Object type for packet_queue field (available as of API version 6.5)
* `enable_wred_tuning` - Enables custom tuning of WRED values (available as of API version 6.5)
* `enable_ecn` - Enables Explicit Congestion Notification for WRED (available as of API version 6.5)
* `enable_watchdog_tuning` - Enables custom tuning of Watchdog values (available as of API version 6.5)
* `cli_commands` - CLI Commands (available as of API version 6.5)
* `detect_bridging_loops` - Enable Detection of Bridging Loops (available as of API version 6.5)
* `unidirectional_link_detection` - Enable Detection of Unidirectional Link (available as of API version 6.5)
* `mac_security_mode` - MAC security mode (available as of API version 6.5)
* `mac_limit` - Between 1-1000 (available as of API version 6.5)
* `security_violation_action` - Security violation action (available as of API version 6.5)
* `aging_type` - Limit MAC authentication based on inactivity or absolute time (available as of API version 6.5)
* `aging_time` - In minutes, how long the client will stay authenticated (available as of API version 6.5)
* `lldp_enable` - LLDP enable (available as of API version 6.5)
* `lldp_mode` - LLDP mode. Enables LLDP Rx and/or LLDP Tx (available as of API version 6.5)
* `lldp_med_enable` - LLDP med enable (available as of API version 6.5)

* `lldp_med` - LLDP MED configurations (available as of API version 6.5)
  * `lldp_med_row_num_enable` - Per LLDP Med row enable
  * `lldp_med_row_num_advertised_applicatio` - Advertised application
  * `lldp_med_row_num_dscp_mark` - LLDP DSCP Mark
  * `lldp_med_row_num_priority` - LLDP Priority
  * `lldp_med_row_num_service` - LLDP Service
  * `lldp_med_row_num_service_ref_type_` - Object type for lldp_med_row_num_service field
  * `index` - The index identifying the object

## Import

Ethernet port settings resources can be imported using the `name` attribute:

```
$ terraform import verity_eth_port_settings.example example
```
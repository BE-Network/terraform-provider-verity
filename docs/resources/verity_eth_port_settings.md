# Ethernet Port Settings Resource

`verity_eth_port_settings` manages Ethernet port settings in Verity, which define specific configurations for physical Ethernet ports.

## Example Usage

```hcl
resource "verity_eth_port_settings" "example" {
  name = "example"
  enable = true
  auto_negotiation = true
  enable_speed_control = true
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
  minimum_wred_threshold = 1
  maximum_wred_threshold = 1
  wred_drop_probability = 0
  priority_flow_control_watchdog_action = "DROP"
  priority_flow_control_watchdog_detect_time = 100
  priority_flow_control_watchdog_restore_time = 100
  packet_queue = ""
  packet_queue_ref_type_ = ""
  enable_wred_tuning = false
  enable_ecn = true
  enable_watchdog_tuning = false
  cli_commands = ""
  detect_bridging_loops = false
  unidirectional_link_detection = false
  mac_security_mode = "disabled"
  mac_limit = 1000
  security_violation_action = "protect"
  aging_type = "absolute"
  aging_time = 0
  lldp_enable = true
  lldp_mode = "RxAndTx"
  lldp_med_enable = false

  lldp_med {
    lldp_med_row_num_enable = false
    lldp_med_row_num_advertised_applicatio = ""
    lldp_med_row_num_dscp_mark = 0
    lldp_med_row_num_priority = 0
    lldp_med_row_num_service = ""
    lldp_med_row_num_service_ref_type_ = ""
    index = 1
  }

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `auto_negotiation` (Boolean) - Indicates if duplex mode should be auto negotiated.
* `enable_speed_control` (Boolean) - Turns on speed control fields.
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
* `max_allowed_value` (Integer) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action.
* `max_allowed_unit` (String) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action. %: Percentage. kbps: kilobits per second. mbps: megabits per second. gbps: gigabits per second. pps: packet per second. kpps: kilopacket per second.
* `action` (String) - Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops. Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops. Alarm is raised. Alarm automatically clears when rate is below configured threshold. Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service.
* `fec` (String) - FEC is Forward Error Correction which is error correction on the fiber link. Any: Allows switch Negotiation between FC and RS. None: Disables FEC on an interface. FC: Enables FEC on supported interfaces. FC stands for fire code. RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. None: VnetC doesn't alter the Switch Value.
* `single_link` (Boolean) - Ports with this setting will be disabled when link state tracking takes effect.
* `minimum_wred_threshold` (Integer) - A value between 1 to 12480(in KiloBytes). 
* `maximum_wred_threshold` (Integer) - A value between 1 to 12480(in KiloBytes).
* `wred_drop_probability` (Integer) - A value between 0 to 100.
* `priority_flow_control_watchdog_action` (String) - Ports with this setting will be disabled when link state tracking takes effect.
* `priority_flow_control_watchdog_detect_time` (Integer) - A value between 100 to 5000.
* `priority_flow_control_watchdog_restore_time` (Integer) - A value between 100 to 60000.
* `object_properties` (Object) - 
  * `group` (String) - Group.
* `packet_queue` (String) - Packet Queue.
* `packet_queue_ref_type_` (String) - Object type for packet_queue field.
* `enable_wred_tuning` (Boolean) - Enables custom tuning of WRED values. Uncheck to use Switch default values.
* `enable_ecn` (Boolean) - Enables Explicit Congestion Notification for WRED.
* `enable_watchdog_tuning` (Boolean) - Enables custom tuning of Watchdog values. Uncheck to use Switch default values.
* `cli_commands` (String) - CLI Commands.
* `detect_bridging_loops` (Boolean) - Enable Detection of Bridging Loops.
* `unidirectional_link_detection` (Boolean) - Enable Detection of Unidirectional Link.
* `mac_security_mode` (String) - Dynamic - MACs are learned and aged normally up to the limit. Packets will be dropped from clients exceeding the limit. Once a client ages out, a new client can take its slot. When the port goes operationally down (disconnecting or disabling), the MACs will be flushed. Sticky - Semi permenant learning. Packets will be dropped from clients exceeding the limit. Addresses do not age out or move within the same switch. Operationally downing a port (disconnecting) does NOT flush the entries. Learned MACs can only be flushed by administratively taking the port down or rebooting the switch.
* `mac_limit` (Integer) - Between 1-1000.
* `security_violation_action` (String) - Protect - All packets are dropped from clients above the MAC Limit. Exceeding the limit is not alarmed. Restrict - All packets are dropped from clients above the MAC Limit. Alarm is raised while attempts to exceed limit are active (MAC has not aged). Alarm automatically clears. Shutdown - Alarm is raised and port is taken down if attempt to exceed MAC limit is made. User must administratively Disable and Enable the port to restore service.
* `aging_type` (String) - Limit MAC authentication based on inactivity or on absolute time. See Also Aging Time.
* `aging_time` (Integer) - In minutes, how long the client will stay authenticated. See Also Aging Type.
* `lldp_enable` (Boolean) - LLDP enable.
* `lldp_mode` (String) - LLDP mode. Enables LLDP Rx and/or LLDP Tx.
* `lldp_med_enable` (Boolean) - LLDP med enable.
* `lldp_med` (Array) - 
  * `lldp_med_row_num_enable` (Boolean) - Per LLDP Med row enable.
  * `lldp_med_row_num_advertised_applicatio` (String) - Advertised application.
  * `lldp_med_row_num_dscp_mark` (Integer) - LLDP DSCP Mark.
  * `lldp_med_row_num_priority` (Integer) - LLDP Priority.
  * `lldp_med_row_num_service` (String) - LLDP Service.
  * `lldp_med_row_num_service_ref_type_` (String) - Object type for lldp_med_row_num_service field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Ethernet port settings resources can be imported using the `name` attribute:

```sh
terraform import verity_eth_port_settings.<resource_name> <name>
```
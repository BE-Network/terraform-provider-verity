# verity_eth_port_settings (Resource)

Manages Ethernet port settings.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_eth_port_settings" "example" {
  name = "example"
  action = ""
  aging_time = null
  aging_type = ""
  allocated_power = ""
  auto_negotiation = false
  bpdu_filter = false
  bpdu_guard = false
  broadcast = false
  bsp_enable = false
  cli_commands = ""
  detect_bridging_loops = false
  duplex_mode = ""
  enable = false
  enable_speed_control = false
  fast_learning_mode = false
  fec = ""
  guard_loop = false
  lldp_enable = false
  lldp_med_enable = false
  lldp_mode = ""
  mac_limit = null
  mac_security_mode = ""
  max_allowed_unit = ""
  max_allowed_value = null
  max_bit_rate = ""
  mtu = null
  multicast = false
  packet_queue = ""
  packet_queue_ref_type_ = "packet_queue"
  poe_enable = false
  priority = ""
  security_violation_action = ""
  standalone_link_training = false
  stp_enable = false
  unidirectional_link_detection = false

  lldp_med {
    index = 1
    lldp_med_row_num_advertised_applicatio = ""
    lldp_med_row_num_dscp_mark = null
    lldp_med_row_num_enable = false
    lldp_med_row_num_priority = null
    lldp_med_row_num_service = ""
    lldp_med_row_num_service_ref_type_ = "service"
  }
}
```

### Datacenter mode

```hcl
resource "verity_eth_port_settings" "example" {
  name = "example"
  action = ""
  allocated_power = ""
  auto_negotiation = false
  bpdu_filter = false
  bpdu_guard = false
  broadcast = false
  bsp_enable = false
  cli_commands = ""
  duplex_mode = ""
  enable = false
  enable_ecn = false
  enable_speed_control = false
  enable_watchdog_tuning = false
  enable_wred_tuning = false
  fast_learning_mode = false
  fec = ""
  guard_loop = false
  max_allowed_unit = ""
  max_allowed_value = null
  max_bit_rate = ""
  maximum_wred_threshold = null
  minimum_wred_threshold = null
  mtu = null
  multicast = false
  packet_queue = ""
  packet_queue_ref_type_ = "packet_queue"
  poe_enable = false
  priority = ""
  priority_flow_control_watchdog_action = ""
  priority_flow_control_watchdog_detect_time = null
  priority_flow_control_watchdog_restore_time = null
  single_link = false
  standalone_link_training = false
  stp_enable = false
  wred_drop_probability = null
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `action` (String) - Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: <div class="tab"> Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops. Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops. Alarm is raised . Alarm automatically clears when rate is below configured threshold. Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service. </div>.
* `aging_time` (Integer) - In minutes, how long the client will stay authenticated. See Also Aging Type. Campus mode only. Set it to `null` to clear it.
* `aging_type` (String) - Limit MAC authentication based on inactivity or on absolute time. See Also Aging Time. Campus mode only.
* `allocated_power` (String) - Power the PoE system will attempt to allocate on this port.
* `auto_negotiation` (Boolean) - Indicates if duplex mode should be auto negotiated.
* `bpdu_filter` (Boolean) - Drop all Rx and Tx BPDUs.
* `bpdu_guard` (Boolean) - Block port on BPDU Receive.
* `broadcast` (Boolean) - Broadcast.
* `bsp_enable` (Boolean) - Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU.
* `cli_commands` (String) - CLI Commands.
* `detect_bridging_loops` (Boolean) - Enable Detection of Bridging Loops. Campus mode only.
* `duplex_mode` (String) - Duplex Mode.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `enable_ecn` (Boolean) - Enables Explicit Congestion Notification for WRED. Datacenter mode only.
* `enable_speed_control` (Boolean) - Turns on speed control fields.
* `enable_watchdog_tuning` (Boolean) - Enables custom tuning of Watchdog values. Uncheck to use Switch default values. Datacenter mode only.
* `enable_wred_tuning` (Boolean) - Enables custom tuning of WRED values. Uncheck to use Switch default values. Datacenter mode only.
* `fast_learning_mode` (Boolean) - Enable Immediate Transition to Forwarding.
* `fec` (String) - FEC is Forward Error Correction which is error correction on the fiber link. <div class="tab"> Any: Allows switch Negotiation between FC and RS None: Disables FEC on an interface. FC: Enables FEC on supported interfaces. FC stands for fire code. RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. None: VnetC doesn't alter the Switch Value. </div>.
* `guard_loop` (Boolean) - Enable Cisco Guard Loop.
* `lldp_enable` (Boolean) - LLDP enable. Campus mode only.
* `lldp_med` (Block List) - LLDP MED configurations. Entries are matched by `index`. Campus mode only.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `lldp_med_row_num_advertised_applicatio` (String) - Advertised application.
  * `lldp_med_row_num_dscp_mark` (Integer) - Defines egress LLDP sent when a device is connected to this Eth-Port Settings allowing the device to auto-provision its DSCP marking. Set it to `null` to clear it.
  * `lldp_med_row_num_enable` (Boolean) - Per LLDP Med row enable.
  * `lldp_med_row_num_priority` (Integer) - LLDP Priority. Set it to `null` to clear it.
  * `lldp_med_row_num_service` (String) - LLDP Service. Set together with `lldp_med_row_num_service_ref_type_`.
  * `lldp_med_row_num_service_ref_type_` (String) - Object type for lldp_med_row_num_service field.
* `lldp_med_enable` (Boolean) - LLDP med enable. Campus mode only.
* `lldp_mode` (String) - LLDP mode. Enables LLDP Rx and/or LLDP Tx. Campus mode only.
* `mac_limit` (Integer) - Between 1-1000. Campus mode only. Set it to `null` to clear it.
* `mac_security_mode` (String) - Dynamic - MACs are learned and aged normally up to the limit. <div class="tab"> Packets will be dropped from clients exceeding the limit. Once a client ages out, a new client can take its slot. When the port goes operationally down (disconnecting or disabling), the MACs will be flushed. </div> Sticky - Semi permenant learning. <div class="tab"> Packets will be dropped from clients exceeding the limit. Addresses do not age out or move within the same switch. Operationally downing a port (disconnecting) does NOT flush the entries. Learned MACs can only be flushed by administratively taking the port down or rebooting the switch. </div>. Campus mode only.
* `max_allowed_unit` (String) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action <div class="tab"> %: Percentage. kbps: kilobits per second mbps: megabits per second gbps: gigabits per second pps: packet per second kpps: kilopacket per second </div>.
* `max_allowed_value` (Integer) - Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action. Set it to `null` to clear it.
* `max_bit_rate` (String) - Maximum Bit Rate allowed.
* `maximum_wred_threshold` (Integer) - A value between 1 to 12480(in KiloBytes). Datacenter mode only. Set it to `null` to clear it.
* `minimum_wred_threshold` (Integer) - A value between 1 to 12480(in KiloBytes). Datacenter mode only. Set it to `null` to clear it.
* `mtu` (Integer) - MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. Set it to `null` to clear it.
* `multicast` (Boolean) - Multicast.
* `packet_queue` (String) - Packet Queue. Set together with `packet_queue_ref_type_`.
* `packet_queue_ref_type_` (String) - Object type for packet_queue field.
* `poe_enable` (Boolean) - Enable PoE on the port.
* `priority` (String) - Priority given when assigning power in a limited power situation.
* `priority_flow_control_watchdog_action` (String) - Ports with this setting will be disabled when link state tracking takes effect. Datacenter mode only.
* `priority_flow_control_watchdog_detect_time` (Integer) - A value between 100 to 5000. Datacenter mode only. Set it to `null` to clear it.
* `priority_flow_control_watchdog_restore_time` (Integer) - A value between 100 to 60000. Datacenter mode only. Set it to `null` to clear it.
* `security_violation_action` (String) - Protect - All packets are dropped from clients above the MAC Limit. <div class="tab"> Exceeding the limit is not alarmed. </div> Restrict - All packets are dropped from clients above the MAC Limit. <div class="tab"> Alarm is raised while attempts to exceed limit are active (MAC has not aged). Alarm automatically clears. </div> Shutdown - Alarm is raised and port is taken down if attempt to exceed MAC limit is made. <div class="tab"> User must administratively Disable and Enable the port to restore service. </div>. Campus mode only.
* `single_link` (Boolean) - Ports with this setting will be disabled when link state tracking takes effect. Datacenter mode only.
* `standalone_link_training` (Boolean) - For use when the port speed/FEC are manually fixed, but the physical link still needs SerDes tuning most commonly high-speed passive DAC/copper links.
* `stp_enable` (Boolean) - Enable Spanning Tree on the port. Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Fabric Settings.
* `unidirectional_link_detection` (Boolean) - Enable Detection of Unidirectional Link. Campus mode only.
* `wred_drop_probability` (Integer) - A value between 0 to 100. Datacenter mode only. Set it to `null` to clear it.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `lldp_med.lldp_med_row_num_service` | `lldp_med_row_num_service_ref_type_` | `service` |
| `packet_queue` | `packet_queue_ref_type_` | `packet_queue` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_eth_port_settings.<resource_name> <name>
```

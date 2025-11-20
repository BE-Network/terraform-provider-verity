# Device Settings Resource

`verity_device_settings` manages Device Settings resources in Verity, which define configuration profiles for Ethernet devices and switches.

## Example Usage

```hcl
resource "verity_device_settings" "example" {
  name = "example"
  enable = false
  mode = "IEEE 802.3af"
  usage_threshold = 0.99
  external_battery_power_available = 40
  external_power_available = 75
  disable_tcp_udp_learned_packet_acceleration = false
  packet_queue = "packet_queue|(Packet Queue)|"
  packet_queue_ref_type_ = "packet_queue"
  security_audit_interval = 60
  commit_to_flash_interval = 60
  rocev2 = false
  cut_through_switching = false
  hold_timer = 0
  mac_aging_timer_override = null
  spanning_tree_priority = "byLevel"

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `mode` (String) - Mode.
* `usage_threshold` (Number) - Usage Threshold.
* `external_battery_power_available` (Integer) - External Battery Power Available.
* `external_power_available` (Integer) - External Power Available.
* `disable_tcp_udp_learned_packet_acceleration` (Boolean) - Required for AVB, PTP and Cobranet Support for ONT Devices.
* `packet_queue` (String) - Packet Queue for device.
* `packet_queue_ref_type_` (String) - Object type for packet_queue field.
* `security_audit_interval` (Integer) - Frequency in minutes of rereading this Switch running configuration and comparing it to expected values. If the value is blank, audit will use default switch settings. If the value is 0, audit will be turned off.
* `commit_to_flash_interval` (Integer) - Time delay in minutes to write the Switch configuration to flash after a change is made. If the value is blank, commit will use default switch settings of 12 hours. If the value is 0, commit will be turned off.
* `rocev2` (Boolean) - Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs.
* `cut_through_switching` (Boolean) - Enable Cut-through Switching on all Switches.
* `object_properties` (Object) - 
  * `group` (String) - Group.
* `hold_timer` (Integer) - Hold Timer.
* `mac_aging_timer_override` (Integer) - Blank uses the Device's default; otherwise an integer between 1 to 1,000,000 seconds.
* `spanning_tree_priority` (String) - STP per switch, priority are in 4096 increments, the lower the number, the higher the priority.

## Import

Device Settings resources can be imported using the `name` attribute:

```sh
terraform import verity_device_settings.<resource_name> <name>
```

# Device Settings Resource

`verity_device_settings` manages Device Settings resources in Verity, which define configuration profiles for Ethernet devices and switches.

## Example Usage

```hcl
resource "verity_device_settings" "test1" {
  name = "test1"
  depends_on = [verity_operation_stage.device_settings_stage]
  object_properties {
    group = ""
    isdefault = false
  }
  commit_to_flash_interval = 60
  cut_through_switching = false
  disable_tcp_udp_learned_packet_acceleration = false
  enable = false
  external_battery_power_available = 40
  external_power_available = 75
  mode = "IEEE 802.3af"
  rocev2 = false
  security_audit_interval = 60
  usage_threshold = 0.99
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the Device Settings profile.
* `enable` - (Optional) Enable this Device Settings profile. Default is `false`.
* `mode` - (Optional) Mode. Allowed values: `Manual`, `IEEE 802.3af`. Default is `IEEE 802.3af`.
* `usage_threshold` - (Optional) Usage Threshold. Default is `0.99`.
* `external_battery_power_available` - (Optional) External Battery Power Available. Default is `40`. Maximum is `2000`.
* `external_power_available` - (Optional) External Power Available. Default is `75`. Maximum is `2000`.
* `disable_tcp_udp_learned_packet_acceleration` - (Optional) Required for AVB, PTP and Cobranet Support. Default is `false`.
* `security_audit_interval` - (Optional) Frequency in minutes of rereading this Switch running configuration and comparing it to expected values. Default is `60`. Maximum is `1440`.
* `commit_to_flash_interval` - (Optional) Frequency in minutes to write the Switch configuration to flash. Default is `60`. Maximum is `1440`.
* `rocev2` - (Optional) Enable RDMA over Converged Ethernet version 2 network protocol. Default is `false`.
* `cut_through_switching` - (Optional) Enable Cut-through Switching on all Switches. Default is `false`.
* `object_properties` - (Optional) Object properties configuration:
  * `group` - (Optional) Group name. Default is `""`.
  * `isdefault` - (Optional) Default object. Default is `false`.
* `hold_timer` - (Optional) Hold Timer. Default is `0`. Maximum is `86400`.
* `mac_aging_timer_override` - (Optional) Blank uses the Device's default; otherwise an integer between 1 to 1,000,000 seconds.
* `spanning_tree_priority` - (Optional) STP per switch, priority are in 4096 increments, the lower the number, the higher the priority. Allowed values: `8192`, `32768`, `36864`, `20480`, `4096`, `57344`, `45056`, `61440`, `byLevel`, `16384`, `40960`, `24576`, `49152`, `0`, `12288`, `28672`, `53248`. Default is `byLevel`.
* `packet_queue_id` - (Optional) Packet Queue for device. Default is `packet_queue|(Packet Queue)|`.
* `packet_queue_id_ref_type_` - (Optional) Object type for packet_queue_id field. Allowed value: `packet_queue`.

## Import

Device Settings resources can be imported using the `name` attribute:

```
$ terraform import verity_device_settings.<resource_name> <name>
```

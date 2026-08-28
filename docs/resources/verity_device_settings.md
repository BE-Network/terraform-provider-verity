# Device Settings Resource

`verity_device_settings` manages Device Settings resources in Verity, which define Ethernet device profile settings.

## Example Usage

```hcl
resource "verity_device_settings" "example" {
  name = "example"
  enable = true
  cli_commands = ""
  mode = "IEEE 802.3af"
  usage_threshold = 0.99
  external_battery_power_available = 40
  external_power_available = 75
  disable_tcp_udp_learned_packet_acceleration = false
  packet_queue = "queue-1"
  packet_queue_ref_type_ = "packet_queue"
  security_audit_interval = 60
  commit_to_flash_interval = 60
  rocev2 = false
  cut_through_switching = false
  login_banner = ""
  device_aaa_profile = "aaa-default"
  device_aaa_profile_ref_type_ = "device_aaa_profile"
  hold_timer = 0
  mac_aging_timer_override = null
  spanning_tree_priority = "byLevel"
  ntp_vrf = "tenant"
  ntp_vrf_tenant = "tenant-a"
  ntp_vrf_tenant_ref_type_ = "tenant"

  dns_servers {
    index = 1
    enabled = true
    server = "8.8.8.8"
  }

  ntp_servers {
    index = 1
    enabled = true
    server = "pool.ntp.org"
  }

  syslog_servers {
    index = 1
    enabled = true
    scheme = "udp_bsd"
    server = "syslog.example.com"
    port = "514"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `cli_commands` (String) - CLI Commands.
* `mode` (String) - Mode.
* `usage_threshold` (Number) - Usage Threshold.
* `external_battery_power_available` (Integer) - External Battery Power Available.
* `external_power_available` (Integer) - External Power Available.
* `security_audit_interval` (Integer) - Frequency in minutes of rereading this switch running configuration and comparing it to expected values.
* `commit_to_flash_interval` (Integer) - Frequency in minutes to write the switch configuration to flash.
* `rocev2` (Boolean) - Enable RDMA over Converged Ethernet version 2 network protocol.
* `cut_through_switching` (Boolean) - Enable Cut-through Switching on all switches.
* `login_banner` (String) - Banner message displayed at login.
* `device_aaa_profile` (String) - Device AAA Profile for authentication settings.
* `device_aaa_profile_ref_type_` (String) - Object type for `device_aaa_profile` field.
* `disable_tcp_udp_learned_packet_acceleration` (Boolean) - Required for AVB, PTP and Cobranet Support.
* `mac_aging_timer_override` (Integer) - Blank uses the device default; otherwise an integer between 1 and 1,000,000 seconds.
* `spanning_tree_priority` (String) - STP per switch priority.
* `packet_queue` (String) - Packet Queue for device.
* `packet_queue_ref_type_` (String) - Object type for `packet_queue` field.
* `hold_timer` (Integer) - Hold Timer.
* `ntp_vrf` (String) - VRF used for NTP servers. Valid values are `default/underlay`, `mgmt` and `tenant`. Defaults to `mgmt`.
* `ntp_vrf_tenant` (String) - Tenant used for NTP servers. Applies when `ntp_vrf` is `tenant`.
* `ntp_vrf_tenant_ref_type_` (String) - Object type for `ntp_vrf_tenant` field.
* `dns_servers` (Array) - DNS servers.
  * `enabled` (Boolean) - Enable DNS name server.
  * `server` (String) - IPv4 or IPv6 DNS name server address.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ntp_servers` (Array) - NTP servers.
  * `enabled` (Boolean) - Enable NTP server.
  * `server` (String) - IPv4, IPv6, or DNS name for NTP server.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `syslog_servers` (Array) - Syslog servers.
  * `enabled` (Boolean) - Enable syslog server.
  * `scheme` (String) - Syslog connection scheme.
  * `server` (String) - IPv4, IPv6, or DNS name for syslog server.
  * `port` (String) - Syslog server port.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Device Settings resources can be imported using the `name` attribute:

```sh
terraform import verity_device_settings.<resource_name> <name>
```

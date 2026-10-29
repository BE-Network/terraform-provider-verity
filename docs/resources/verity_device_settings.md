# verity_device_settings (Resource)

Manages a Verity Device Settings (Eth Device Profile).

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_device_settings" "example" {
  name = "example"
  cli_commands = ""
  commit_to_flash_interval = null
  cut_through_switching = false
  device_aaa_profile = ""
  device_aaa_profile_ref_type_ = "device_aaa_profile"
  disable_tcp_udp_learned_packet_acceleration = false
  enable = false
  external_battery_power_available = null
  external_power_available = null
  hold_timer = null
  login_banner = ""
  mac_aging_timer_override = null
  mode = ""
  ntp_vrf = ""
  ntp_vrf_tenant = ""
  ntp_vrf_tenant_ref_type_ = "tenant"
  packet_queue = ""
  packet_queue_ref_type_ = "packet_queue"
  rocev2 = false
  security_audit_interval = null
  spanning_tree_priority = ""
  usage_threshold = null

  dns_servers {
    index = 1
    enabled = false
    server = ""
  }

  ntp_servers {
    index = 1
    enabled = false
    server = ""
  }

  object_properties {
  }

  syslog_servers {
    index = 1
    enabled = false
    port = ""
    scheme = ""
    server = ""
  }
}
```

### Datacenter mode

```hcl
resource "verity_device_settings" "example" {
  name = "example"
  cli_commands = ""
  commit_to_flash_interval = null
  cut_through_switching = false
  device_aaa_profile = ""
  device_aaa_profile_ref_type_ = "device_aaa_profile"
  disable_tcp_udp_learned_packet_acceleration = false
  enable = false
  external_battery_power_available = null
  external_power_available = null
  login_banner = ""
  mode = ""
  ntp_vrf = ""
  ntp_vrf_tenant = ""
  ntp_vrf_tenant_ref_type_ = "tenant"
  packet_queue = ""
  packet_queue_ref_type_ = "packet_queue"
  rocev2 = false
  security_audit_interval = null
  usage_threshold = null

  dns_servers {
    index = 1
    enabled = false
    server = ""
  }

  ntp_servers {
    index = 1
    enabled = false
    server = ""
  }

  object_properties {
  }

  syslog_servers {
    index = 1
    enabled = false
    port = ""
    scheme = ""
    server = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object Name. Must be unique. Changing it replaces the resource.

### Optional

* `cli_commands` (String) - CLI Commands.
* `commit_to_flash_interval` (Integer) - Frequency in minutes to write the Switch configuration to flash. If the value is blank, commit will use default switch settings. If the value is 0, commit will be turned off. (maximum: 1440). Set it to `null` to clear it.
* `cut_through_switching` (Boolean) - Enable Cut-through Switching on all Switches.
* `device_aaa_profile` (String) - Device AAA Profile for authentication settings. Set together with `device_aaa_profile_ref_type_`.
* `device_aaa_profile_ref_type_` (String) - Object type for device_aaa_profile field.
* `disable_tcp_udp_learned_packet_acceleration` (Boolean) - Required for AVB, PTP and Cobranet Support.
* `dns_servers` (Block List) - DNS servers. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable DNS name server.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `server` (String) - IPv4 or IPv6 DNS name server address.
* `enable` (Boolean) - Enable object.
* `external_battery_power_available` (Integer) - External Battery Power Available (maximum: 2000). Set it to `null` to clear it.
* `external_power_available` (Integer) - External Power Available (maximum: 2000). Set it to `null` to clear it.
* `hold_timer` (Integer) - Hold Timer (maximum: 86400). Campus mode only. Set it to `null` to clear it.
* `login_banner` (String) - Banner message displayed at login.
* `mac_aging_timer_override` (Integer) - Blank uses the Device's default; otherwise an integer between 1 to 1,000,000 seconds (minimum: 1, maximum: 1000000). Campus mode only. Set it to `null` to clear it.
* `mode` (String) - Mode.
* `ntp_servers` (Block List) - NTP servers. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable NTP server.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `server` (String) - IPv4, IPv6, or DNS name for NTP server.
* `ntp_vrf` (String) - VRF used for NTP Servers. Valid values are `default/underlay`, `mgmt` and `tenant`.
* `ntp_vrf_tenant` (String) - Tenant used for NTP Servers. Set together with `ntp_vrf_tenant_ref_type_`.
* `ntp_vrf_tenant_ref_type_` (String) - Object type for ntp_vrf_tenant field.
* `object_properties` (Block) - Object properties for the Device Settings. At most one block.
* `packet_queue` (String) - Packet Queue for device. Set together with `packet_queue_ref_type_`.
* `packet_queue_ref_type_` (String) - Object type for packet_queue field.
* `rocev2` (Boolean) - Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs.
* `security_audit_interval` (Integer) - Frequency in minutes of rereading this Switch running configuration and comparing it to expected values. If the value is blank, audit will use default switch settings. If the value is 0, audit will be turned off. (maximum: 1440). Set it to `null` to clear it.
* `spanning_tree_priority` (String) - STP per switch, priority are in 4096 increments, the lower the number, the higher the priority. Campus mode only.
* `syslog_servers` (Block List) - Syslog servers. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable syslog server.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port` (String) - Syslog server port.
  * `scheme` (String) - Syslog connection scheme.
  * `server` (String) - IPv4, IPv6, or DNS name for syslog server.
* `usage_threshold` (Number) - Usage Threshold. Set it to `null` to clear it.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `device_aaa_profile` | `device_aaa_profile_ref_type_` | `device_aaa_profile` |
| `ntp_vrf_tenant` | `ntp_vrf_tenant_ref_type_` | `tenant` |
| `packet_queue` | `packet_queue_ref_type_` | `packet_queue` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_device_settings.<resource_name> <name>
```

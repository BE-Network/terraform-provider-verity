# Device Controller Resource

`verity_device_controller` manages device controllers in Verity, which define configurations for network device management.

## Example Usage

```hcl
resource "verity_device_controller" "example" {
  name = "example_controller"
  enable = true
  ip_source = "static"
  controller_ip_and_mask = "192.168.1.10/24"
  gateway = "192.168.1.1"
  comm_type = "snmp"
  snmp_community_string = "public"
  communication_mode = "out-of-band"
  cli_access_mode = "ssh"
  username = "admin"
  password = "password"
  device_managed_as = "controller"
  switchpoint = "main_switchpoint"
  switchpoint_ref_type_ = "switchpoint"
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the device controller.
* `enable` - (Optional) Enable this device controller. Default is `false`.
* `ip_source` - (Optional) Source of IP address assignment (e.g., `static`, `dhcp`).
* `controller_ip_and_mask` - (Optional) Controller IP address and subnet mask in CIDR notation.
* `gateway` - (Optional) Gateway IP address for the controller.
* `switch_ip_and_mask` - (Optional) Switch IP address and subnet mask in CIDR notation.
* `switch_gateway` - (Optional) Gateway IP address for the switch.
* `comm_type` - (Optional) Communication type (e.g., `snmp`, `cli`).
* `snmp_community_string` - (Optional) SNMP community string for SNMP communication.
* `uplink_port` - (Optional) Uplink port identifier.
* `lldp_search_string` - (Optional) LLDP search string for device discovery.
* `ztp_identification` - (Optional) Zero Touch Provisioning identification method.
* `located_by` - (Optional) Method used to locate the device.
* `power_state` - (Optional) Power state of the device.
* `communication_mode` - (Optional) Communication mode (e.g., `in-band`, `out-of-band`).
* `cli_access_mode` - (Optional) CLI access mode (e.g., `ssh`, `telnet`).
* `username` - (Optional) Username for device authentication.
* `password` - (Optional) Password for device authentication.
* `enable_password` - (Optional) Enable password for privileged mode.
* `ssh_key_or_password` - (Optional) SSH key or password for SSH authentication.
* `managed_on_native_vlan` - (Optional) Whether the device is managed on native VLAN. Default is `false`.
* `sdlc` - (Optional) Software-defined lifecycle configuration.
* `switchpoint` - (Optional) Reference to a switchpoint resource.
* `switchpoint_ref_type_` - (Optional) Object type for switchpoint reference.
* `security_type` - (Optional) Security type configuration.
* `snmpv3_username` - (Optional) SNMPv3 username for SNMPv3 authentication.
* `authentication_protocol` - (Optional) Authentication protocol for SNMPv3.
* `passphrase` - (Optional) Passphrase for SNMPv3 authentication.
* `private_protocol` - (Optional) Privacy protocol for SNMPv3.
* `private_password` - (Optional) Privacy password for SNMPv3.
* `device_managed_as` - (Optional) How the device is managed (e.g., `controller`, `switch`).
* `switch` - (Optional) Reference to a switch resource.
* `switch_ref_type_` - (Optional) Object type for switch reference.
* `connection_service` - (Optional) Reference to a connection service.
* `connection_service_ref_type_` - (Optional) Object type for connection service reference.
* `port` - (Optional) Port identifier.
* `sfp_mac_address_or_sn` - (Optional) SFP MAC address or serial number.
* `uses_tagged_packets` - (Optional) Whether the device uses tagged packets. Default is `false`.

## Import

Device Controller resources can be imported using the `name` attribute:

```
$ terraform import verity_device_controller.example example_controller
```

# Device Controller Resource

`verity_device_controller` manages device controllers in Verity, which define configurations for network device management.

## Example Usage

```hcl
resource "verity_device_controller" "example" {
  name = "example"
  enable = false
  ip_source = "dhcp"
  controller_ip_and_mask = ""
  gateway = ""
  switch_ip_and_mask = ""
  switch_gateway = ""
  comm_type = "snmpv2"
  snmp_community_string = ""
  uplink_port = ""
  lldp_search_string = ""
  ztp_identification = ""
  located_by = "LLDP"
  power_state = "on"
  communication_mode = "generic_snmp"
  cli_access_mode = "SSH"
  username = ""
  password = ""
  enable_password = ""
  ssh_key_or_password = ""
  managed_on_native_vlan = false
  sdlc = ""
  switchpoint = ""
  switchpoint_ref_type_ = "switchpoint"
  security_type = "noAuthNoPriv"
  snmpv3_username = ""
  authentication_protocol = "MD5"
  passphrase = ""
  private_protocol = "DES"
  private_password = ""
  password_encrypted = ""
  enable_password_encrypted = ""
  ssh_key_or_password_encrypted = ""
  passphrase_encrypted = ""
  private_password_encrypted = ""
  device_managed_as = "switch"
  switch = ""
  switch_ref_type_ = "switchpoint"
  connection_service = ""
  connection_service_ref_type_ = "service"
  port = ""
  sfp_mac_address_or_sn = ""
  uses_tagged_packets = true
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `ip_source` (String) - IP Source.
* `controller_ip_and_mask` (String) - Controller IP and Mask.
* `gateway` (String) - Gateway.
* `switch_ip_and_mask` (String) - Switch IP and Mask.
* `switch_gateway` (String) - Gateway of Managed Device.
* `comm_type` (String) - Comm Type.
* `snmp_community_string` (String) - Comm Credentials.
* `uplink_port` (String) - Uplink Port of Managed Device.
* `lldp_search_string` (String) - Optional unless Located By is "LLDP" or Device managed as "Active SFP". Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used.
* `ztp_identification` (String) - Service Tag or Serial Number to identify device for Zero Touch Provisioning.
* `located_by` (String) - Controls how the system locates this Device within its LAN.
* `power_state` (String) - Power state of Switch Controller.
* `communication_mode` (String) - Communication Mode.
* `cli_access_mode` (String) - CLI Access Mode.
* `username` (String) - Username.
* `password` (String) - Password.
* `enable_password` (String) - Enable Password - to enable privileged CLI operations.
* `ssh_key_or_password` (String) - SSH Key or Password.
* `managed_on_native_vlan` (Boolean) - Managed on native VLAN.
* `sdlc` (String) - SDLC that Device Controller belongs to.
* `switchpoint` (String) - Switchpoint reference.
* `switchpoint_ref_type_` (String) - Object type for switchpoint field.
* `security_type` (String) - Security level.
* `snmpv3_username` (String) - Username.
* `authentication_protocol` (String) - Protocol.
* `passphrase` (String) - Passphrase.
* `private_protocol` (String) - Protocol.
* `private_password` (String) - Password.
* `password_encrypted` (String) - Password.
* `enable_password_encrypted` (String) - Enable Password - to enable privileged CLI operations.
* `ssh_key_or_password_encrypted` (String) - SSH Key or Password.
* `passphrase_encrypted` (String) - Passphrase.
* `private_password_encrypted` (String) - Password.
* `device_managed_as` (String) - Device managed as.
* `switch` (String) - Switchpoint locating the Switch to be controlled.
* `switch_ref_type_` (String) - Object type for switch field.
* `connection_service` (String) - Connect a Service.
* `connection_service_ref_type_` (String) - Object type for connection_service field.
* `port` (String) - Port locating the Switch to be controlled.
* `sfp_mac_address_or_sn` (String) - SFP MAC Address or SN.
* `uses_tagged_packets` (Boolean) - Indicates if the direct interface expects tagged or untagged packets.

## Import

Device Controller resources can be imported using the `name` attribute:

```sh
terraform import verity_device_controller.<resource_name> <name>
```

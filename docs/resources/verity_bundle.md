# Bundle Resource

`verity_bundle` manages bundle resources in Verity, which define collections of configurations.

## Example Usage

```hcl
resource "verity_bundle" "example" {
  name = "example"
  enable = false
  protocol = "SIP"
  device_settings = "eth_device_profile|(Device Settings)|"
  device_settings_ref_type_ = "eth_device_profiles"
  cli_commands = ""
  diagnostics_profile = ""
  diagnostics_profile_ref_type_ = "diagnostics_profile"
  device_voice_settings = "voice_device_profile|(SIP Voice Device)|"
  device_voice_settings_ref_type_ = "device_voice_settings"

  eth_port_paths {
    eth_port_num_eth_port_profile = ""
    eth_port_num_eth_port_profile_ref_type_ = ""
    eth_port_num_eth_port_settings = ""
    eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
    eth_port_num_gateway_profile = ""
    eth_port_num_gateway_profile_ref_type_ = ""
    diagnostics_port_profile_num_diagnostics_port_profile = ""
    diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = "diagnostics_port_profile"
    index = 1
    port_name = ""
  }

  user_services {
    row_app_enable = false
    row_app_connected_service = ""
    row_app_connected_service_ref_type_ = "service"
    row_app_cli_commands = ""
    row_ip_mask = ""
    index = 1
  }

  voice_port_profile_paths {
    voice_port_num_voice_port_profiles = ""
    voice_port_num_voice_port_profiles_ref_type_ = "voice_port_profiles"
    index = 0
  }

  object_properties {
    group = ""
    is_for_switch = false
    is_public = false
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `protocol` (String) - Voice Protocol: MGCP or SIP.
* `device_settings` (String) - Device Settings for device.
* `device_settings_ref_type_` (String) - Object type for device_settings field.
* `cli_commands` (String) - CLI Commands.
* `diagnostics_profile` (String) - Diagnostics Profile for device.
* `diagnostics_profile_ref_type_` (String) - Object type for diagnostics_profile field.
* `device_voice_settings` (String) - Device Voice Settings for device.
* `device_voice_settings_ref_type_` (String) - Object type for device_voice_settings field.
* `eth_port_paths` (Array) - 
  * `eth_port_num_eth_port_profile` (String) - Eth Port Profile Or LAG for Eth Port.
  * `eth_port_num_eth_port_profile_ref_type_` (String) - Object type for eth_port_num_eth_port_profile field.
  * `eth_port_num_eth_port_settings` (String) - Choose an Eth Port Settings.
  * `eth_port_num_eth_port_settings_ref_type_` (String) - Object type for eth_port_num_eth_port_settings field.
  * `eth_port_num_gateway_profile` (String) - Gateway Profile or LAG for Eth Port.
  * `eth_port_num_gateway_profile_ref_type_` (String) - Object type for eth_port_num_gateway_profile field.
  * `diagnostics_port_profile_num_diagnostics_port_profile` (String) - Diagnostics Port Profile for port.
  * `diagnostics_port_profile_num_diagnostics_port_profile_ref_type_` (String) - Object type for diagnostics_port_profile_num_diagnostics_port_profile field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port_name` (String) - The name identifying the port. Used for reference only, it won't actually change the port name.
* `user_services` (Array) - 
  * `row_app_enable` (Boolean) - Enable of this User application.
  * `row_app_connected_service` (String) - Service connected to this User application.
  * `row_app_connected_service_ref_type_` (String) - Object type for row_app_connected_service field.
  * `row_app_cli_commands` (String) - CLI Commands of this User application.
  * `row_ip_mask` (String) - IP/Mask.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.
  * `is_for_switch` (Boolean) - Denotes a Switch Bundle.
  * `is_public` (Boolean) - Denotes a shared Switch Bundle.
* `voice_port_profile_paths` (Array) - 
  * `voice_port_num_voice_port_profiles` (String) -  Voice Port Settings for Voice Port.
  * `voice_port_num_voice_port_profiles_ref_type_` (String) - Object type for voice_port_num_voice_port_profiles field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Bundle resources can be imported using the `name` attribute:

```sh
terraform import verity_bundle.<resource_name> <name>
```


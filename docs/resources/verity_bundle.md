# Bundle Resource

`verity_bundle` manages bundle resources in Verity, which define collections of configurations.

## Example Usage

```hcl
resource "verity_bundle" "example" {
  name = "example"
  enable = true
  protocol = "SIP"
  device_settings = "Default"
  device_settings_ref_type_ = "eth_device_profiles"
  device_voice_settings = "Voice_Default"
  device_voice_settings_ref_type_ = "device_voice_settings"
  cli_commands = ""
  
  object_properties {
    is_for_switch = true
    group = "network_devices"
    is_public = false
  }
  
  eth_port_paths {
    port_name = "1/1"
    index = 1
    eth_port_num_eth_port_settings = "Default"
    eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
    eth_port_num_eth_port_profile = ""
    eth_port_num_eth_port_profile_ref_type_ = ""
    eth_port_num_gateway_profile = ""
    eth_port_num_gateway_profile_ref_type_ = ""
  }
  
  user_services {
    index = 1
    row_app_enable = false
    row_app_connected_service = ""
    row_app_connected_service_ref_type_ = ""
    row_app_cli_commands = ""
    row_ip_mask = ""
  }
  
  rg_services {
    index = 0
    row_app_enable = true
    row_app_connected_service = "data_service"
    row_app_connected_service_ref_type_ = "service"
    row_app_type = "data"
    row_ip_mask = "192.168.1.0/24"
  }
  
  voice_port_profile_paths {
    index = 0
    voice_port_num_voice_port_profiles = "voice_profile"
    voice_port_num_voice_port_profiles_ref_type_ = "voice_port_profile"
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the bundle
* `enable` - (Optional) Enable object (available as of API version 6.5)
* `device_settings` - (Optional) Reference to device settings
* `device_settings_ref_type_` - (Optional) Object type for device settings reference
* `cli_commands` - (Optional) CLI commands to execute
* `protocol` - (Optional) Voice Protocol: MGCP or SIP (available as of API version 6.5)
* `device_voice_settings` - (Optional) Device Voice Settings for device (available as of API version 6.5)
* `device_voice_settings_ref_type_` - (Optional) Object type for device_voice_settings field (available as of API version 6.5)
* `object_properties` - (Optional) Object properties block
  * `is_for_switch` - (Optional) Whether this bundle is for a switch
  * `group` - (Optional) Group name (available as of API version 6.5)
  * `is_public` - (Optional) Denotes a shared Switch Bundle (available as of API version 6.5)
* `eth_port_paths` - (Optional) List of Ethernet port path blocks
  * `port_name` - (Optional) Physical port name
  * `index` - (Optional) Index value for ordering
  * `eth_port_num_eth_port_settings` - (Optional) Reference to Ethernet port settings
  * `eth_port_num_eth_port_settings_ref_type_` - (Optional) Object type for Ethernet port settings reference
  * `eth_port_num_eth_port_profile` - (Optional) Reference to an Ethernet port profile
  * `eth_port_num_eth_port_profile_ref_type_` - (Optional) Object type for Ethernet port profile reference
  * `eth_port_num_gateway_profile` - (Optional) Reference to a gateway profile
  * `eth_port_num_gateway_profile_ref_type_` - (Optional) Object type for gateway profile reference
* `user_services` - (Optional) List of user service blocks
  * `index` - (Optional) Index value for ordering
  * `row_app_enable` - (Optional) Enable this service
  * `row_app_connected_service` - (Optional) Reference to a connected service
  * `row_app_connected_service_ref_type_` - (Optional) Object type for connected service reference
  * `row_app_cli_commands` - (Optional) CLI commands specific to this service
  * `row_ip_mask` - (Optional) IP mask in CIDR notation
* `rg_services` - (Optional) List of RG services configurations (available as of API version 6.5)
  * `index` - (Optional) Index identifying this configuration
  * `row_app_enable` - (Optional) Enable this ONT application
  * `row_app_connected_service` - (Optional) Service connected to this ONT application
  * `row_app_connected_service_ref_type_` - (Optional) Object type for connected service reference
  * `row_app_type` - (Optional) Type of ONT Application
  * `row_ip_mask` - (Optional) IP/Mask in IPv4 format
* `voice_port_profile_paths` - (Optional) List of voice port profile configurations (available as of API version 6.5)
  * `index` - (Optional) Index identifying this configuration
  * `voice_port_num_voice_port_profiles` - (Optional) Voice Port Profile for Voice Port
  * `voice_port_num_voice_port_profiles_ref_type_` - (Optional) Object type for voice port profiles reference

## Import

Bundle resources can be imported using the `name` attribute:
```
terraform import verity_bundle.example example
```


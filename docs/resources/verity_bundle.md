# Bundle Resource

`verity_bundle` manages bundle resources in Verity, which define collections of configurations.

## Example Usage

```hcl
resource "verity_bundle" "example" {
  name = "example"
  
  object_properties {
    is_for_switch = true
  }
  
  device_settings = "Default"
  device_settings_ref_type_ = "eth_device_profiles"
  cli_commands = ""
  
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
}
```

## Argument Reference

* `name` - Unique identifier for the bundle
* `device_settings` - Reference to device settings
* `device_settings_ref_type_` - Object type for device settings reference
* `cli_commands` - CLI commands to execute
* `object_properties` - Object properties block
  * `is_for_switch` - Whether this bundle is for a switch
* `eth_port_paths` - List of Ethernet port path blocks
  * `port_name` - Physical port name
  * `index` - Index value for ordering
  * `eth_port_num_eth_port_settings` - Reference to Ethernet port settings
  * `eth_port_num_eth_port_settings_ref_type_` - Object type for Ethernet port settings reference
  * `eth_port_num_eth_port_profile` - Reference to an Ethernet port profile
  * `eth_port_num_eth_port_profile_ref_type_` - Object type for Ethernet port profile reference
  * `eth_port_num_gateway_profile` - Reference to a gateway profile
  * `eth_port_num_gateway_profile_ref_type_` - Object type for gateway profile reference
* `user_services` - List of user service blocks
  * `index` - Index value for ordering
  * `row_app_enable` - Enable this service
  * `row_app_connected_service` - Reference to a connected service
  * `row_app_connected_service_ref_type_` - Object type for connected service reference
  * `row_app_cli_commands` - CLI commands specific to this service
  * `row_ip_mask` - IP mask in CIDR notation

## Import

Bundle resources can be imported using the `name` attribute:
````


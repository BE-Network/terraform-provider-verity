# Bundle Resource

`verity_bundle` manages bundle resources in Verity, which define collections of configurations.

> **Note:** In API 6.4, bundle resources only support **update (PATCH)** operations. Create (PUT) and delete operations are not supported. Bundles must be created through other means before they can be managed with Terraform.

## Example Usage

```hcl
resource "verity_bundle" "example" {
  name = "example"
  device_settings = "eth_device_profile|(Device Settings)|"
  device_settings_ref_type_ = "eth_device_profiles"
  cli_commands = ""

  eth_port_paths {
    eth_port_num_eth_port_profile = ""
    eth_port_num_eth_port_profile_ref_type_ = ""
    eth_port_num_eth_port_settings = ""
    eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
    eth_port_num_gateway_profile = ""
    eth_port_num_gateway_profile_ref_type_ = ""
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

  object_properties {
    is_for_switch = false
  }
}
```

## Argument Reference

* `name` (String, Required) - Object Name. Must be unique.
* `device_settings` (String) - Device Settings for device.
* `device_settings_ref_type_` (String) - Object type for device_settings field.
* `cli_commands` (String) - CLI Commands.
* `eth_port_paths` (Array) - Configuration for ethernet port paths.
  * `eth_port_num_eth_port_profile` (String) - Eth Port Profile Or LAG for Eth Port.
  * `eth_port_num_eth_port_profile_ref_type_` (String) - Object type for eth_port_num_eth_port_profile field.
  * `eth_port_num_eth_port_settings` (String) - Choose an Eth Port Settings.
  * `eth_port_num_eth_port_settings_ref_type_` (String) - Object type for eth_port_num_eth_port_settings field.
  * `eth_port_num_gateway_profile` (String) - Gateway Profile or LAG for Eth Port.
  * `eth_port_num_gateway_profile_ref_type_` (String) - Object type for eth_port_num_gateway_profile field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port_name` (String) - The name identifying the port. Used for reference only, it won't actually change the port name.
* `user_services` (Array) - User service configurations.
  * `row_app_enable` (Boolean) - Enable of this User application.
  * `row_app_connected_service` (String) - Service connected to this User application.
  * `row_app_connected_service_ref_type_` (String) - Object type for row_app_connected_service field.
  * `row_app_cli_commands` (String) - CLI Commands of this User application.
  * `row_ip_mask` (String) - IP/Mask.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - Additional object properties.
  * `is_for_switch` (Boolean) - Denotes a Switch Bundle.

## Import

Bundle resources can be imported using the `name` attribute:

```sh
terraform import verity_bundle.<resource_name> <name>
```


# verity_bundle (Resource)

Manages a Verity Bundle.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_bundle" "example" {
  name = "example"
  cli_commands = ""
  device_settings = ""
  device_settings_ref_type_ = "eth_device_profiles"
  device_voice_settings = ""
  device_voice_settings_ref_type_ = "device_voice_settings"
  diagnostics_profile = ""
  diagnostics_profile_ref_type_ = "diagnostics_profile"
  enable = false
  protocol = ""

  eth_port_paths {
    index = 1
    eth_port_num_diagnostics_port_profile = ""
    eth_port_num_diagnostics_port_profile_ref_type_ = "diagnostics_port_profile"
    eth_port_num_eth_port_profile = ""
    eth_port_num_eth_port_profile_ref_type_ = ""
    eth_port_num_eth_port_settings = ""
    eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
    port_name = ""
  }

  object_properties {
    is_for_switch = false
    is_public = false
  }

  rg_services {
    index = 1
    row_app_connected_service = ""
    row_app_connected_service_ref_type_ = "service"
    row_app_enable = false
    row_app_type = ""
    row_ip_mask = ""
  }

  user_services {
    index = 1
    row_app_cli_commands = ""
    row_app_connected_service = ""
    row_app_connected_service_ref_type_ = "service"
    row_app_enable = false
    row_ip_mask = ""
  }

  voice_port_profile_paths {
    index = 1
    voice_port_num_voice_port_profiles = ""
    voice_port_num_voice_port_profiles_ref_type_ = "voice_port_profiles"
  }
}
```

### Datacenter mode

```hcl
resource "verity_bundle" "example" {
  name = "example"
  cli_commands = ""
  device_settings = ""
  device_settings_ref_type_ = "eth_device_profiles"
  diagnostics_profile = ""
  diagnostics_profile_ref_type_ = "diagnostics_profile"
  enable = false
  protocol = ""

  eth_port_paths {
    index = 1
    eth_port_num_diagnostics_port_profile = ""
    eth_port_num_diagnostics_port_profile_ref_type_ = "diagnostics_port_profile"
    eth_port_num_eth_port_profile = ""
    eth_port_num_eth_port_profile_ref_type_ = ""
    eth_port_num_eth_port_settings = ""
    eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
    eth_port_num_gateway_profile = ""
    eth_port_num_gateway_profile_ref_type_ = ""
    port_name = ""
  }

  object_properties {
    is_for_switch = false
  }

  user_services {
    index = 1
    row_app_cli_commands = ""
    row_app_connected_service = ""
    row_app_connected_service_ref_type_ = "service"
    row_app_enable = false
    row_ip_mask = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `cli_commands` (String) - CLI Commands.
* `device_settings` (String) - Device Settings for device. Set together with `device_settings_ref_type_`.
* `device_settings_ref_type_` (String) - Object type for device_settings field.
* `device_voice_settings` (String) - Device Voice Settings for device. Campus mode only. Set together with `device_voice_settings_ref_type_`.
* `device_voice_settings_ref_type_` (String) - Object type for device_voice_settings field. Campus mode only.
* `diagnostics_profile` (String) - Diagnostics Profile for device. Set together with `diagnostics_profile_ref_type_`.
* `diagnostics_profile_ref_type_` (String) - Object type for diagnostics_profile field.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `eth_port_paths` (Block List) - List of ethernet port configurations. Entries are matched by `index`.
  * `eth_port_num_diagnostics_port_profile` (String) - Diagnostics Port Profile for port. Set together with `eth_port_num_diagnostics_port_profile_ref_type_`.
  * `eth_port_num_diagnostics_port_profile_ref_type_` (String) - Object type for eth_port_num_diagnostics_port_profile field.
  * `eth_port_num_eth_port_profile` (String) - Eth Port Profile Or LAG for Eth Port. Set together with `eth_port_num_eth_port_profile_ref_type_`.
  * `eth_port_num_eth_port_profile_ref_type_` (String) - Object type for eth_port_num_eth_port_profile field.
  * `eth_port_num_eth_port_settings` (String) - Choose an Eth Port Settings. Set together with `eth_port_num_eth_port_settings_ref_type_`.
  * `eth_port_num_eth_port_settings_ref_type_` (String) - Object type for eth_port_num_eth_port_settings field.
  * `eth_port_num_gateway_profile` (String) - Gateway Profile for Eth Port. Datacenter mode only. Set together with `eth_port_num_gateway_profile_ref_type_`.
  * `eth_port_num_gateway_profile_ref_type_` (String) - Object type for eth_port_num_gateway_profile field. Datacenter mode only.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port_name` (String) - The name identifying the port. Used for reference only, it won't actually change the port name.
* `object_properties` (Block) - Object properties for the bundle. At most one block.
  * `is_for_switch` (Boolean) - Denotes a Switch Bundle.
  * `is_public` (Boolean) - Denotes a shared Switch Bundle. Campus mode only.
* `protocol` (String) - Voice Protocol: MGCP or SIP.
* `rg_services` (Block List) - List of RG services configurations. Entries are matched by `index`. Campus mode only.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `row_app_connected_service` (String) - Service connected to this ONT application. Set together with `row_app_connected_service_ref_type_`.
  * `row_app_connected_service_ref_type_` (String) - Object type for row_app_connected_service field.
  * `row_app_enable` (Boolean) - Enable of this ONT application.
  * `row_app_type` (String) - Type of ONT Application.
  * `row_ip_mask` (String) - IP/Mask.
* `user_services` (Block List) - List of user services configurations. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `row_app_cli_commands` (String) - CLI Commands of this User application.
  * `row_app_connected_service` (String) - Service connected to this User application. Set together with `row_app_connected_service_ref_type_`.
  * `row_app_connected_service_ref_type_` (String) - Object type for row_app_connected_service field.
  * `row_app_enable` (Boolean) - Enable of this User application.
  * `row_ip_mask` (String) - IP/Mask.
* `voice_port_profile_paths` (Block List) - List of voice port profile configurations. Entries are matched by `index`. Campus mode only.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `voice_port_num_voice_port_profiles` (String) - Voice Port Settings for Voice Port. Set together with `voice_port_num_voice_port_profiles_ref_type_`.
  * `voice_port_num_voice_port_profiles_ref_type_` (String) - Object type for voice_port_num_voice_port_profiles field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `device_settings` | `device_settings_ref_type_` | `eth_device_profiles` |
| `device_voice_settings` | `device_voice_settings_ref_type_` | `device_voice_settings` |
| `diagnostics_profile` | `diagnostics_profile_ref_type_` | `diagnostics_profile` |
| `eth_port_paths.eth_port_num_diagnostics_port_profile` | `eth_port_num_diagnostics_port_profile_ref_type_` | `diagnostics_port_profile` |
| `eth_port_paths.eth_port_num_eth_port_profile` | `eth_port_num_eth_port_profile_ref_type_` | `authenticated_eth_port`, `eth_port_profile_`, `lag`, `nac_port_profile`, `pb_egress_profile`, `service_port_profile`, `type` |
| `eth_port_paths.eth_port_num_eth_port_settings` | `eth_port_num_eth_port_settings_ref_type_` | `eth_port_settings` |
| `eth_port_paths.eth_port_num_gateway_profile` | `eth_port_num_gateway_profile_ref_type_` | `gateway_profile`, `lag` |
| `rg_services.row_app_connected_service` | `row_app_connected_service_ref_type_` | `service` |
| `user_services.row_app_connected_service` | `row_app_connected_service_ref_type_` | `service` |
| `voice_port_profile_paths.voice_port_num_voice_port_profiles` | `voice_port_num_voice_port_profiles_ref_type_` | `voice_port_profiles` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_bundle.<resource_name> <name>
```

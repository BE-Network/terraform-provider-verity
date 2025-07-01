# Authenticated Eth-Port Resource

`verity_authenticated_eth_port` manages authenticated Ethernet port resources in Verity, which define authentication settings for Ethernet ports.

## Example Usage

```hcl
resource "verity_authenticated_eth_port" "example" {
  name = "example"
  enable = true
  connection_mode = "dot1x"
  reauthorization_period_sec = 3600
  allow_mac_based_authentication = true
  mac_authentication_holdoff_sec = 30
  trusted_port = false

  eth_ports {
    eth_port_profile_num_enable = true
    eth_port_profile_num_eth_port = "port-1"
    eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile"
    eth_port_profile_num_walled_garden_set = false
    eth_port_profile_num_radius_filter_id = "filter-1"
    index = 1
  }

  object_properties {
    group = "auth-ports"
    port_monitoring = "monitor-1"
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the authenticated Ethernet port.
* `enable` - (Optional) Enable this authenticated Ethernet port. Default is `false`.
* `connection_mode` - (Optional) Connection mode for authentication. Possible values include `"dot1x"`.
* `reauthorization_period_sec` - (Optional) Period for reauthorization in seconds.
* `allow_mac_based_authentication` - (Optional) Whether to allow MAC-based authentication.
* `mac_authentication_holdoff_sec` - (Optional) MAC authentication holdoff time in seconds.
* `trusted_port` - (Optional) Whether this is a trusted port.
* `eth_ports` - (Optional) List of Ethernet port configurations:
  * `eth_port_profile_num_enable` - (Optional) Enable this Ethernet port profile.
  * `eth_port_profile_num_eth_port` - (Optional) Ethernet port reference.
  * `eth_port_profile_num_eth_port_ref_type_` - (Optional) Reference type for the Ethernet port.
  * `eth_port_profile_num_walled_garden_set` - (Optional) Whether walled garden is set for this port.
  * `eth_port_profile_num_radius_filter_id` - (Optional) RADIUS filter ID.
  * `index` - (Optional) Index identifying this Ethernet port configuration.
* `object_properties` - (Optional) Object properties configuration:
  * `group` - (Optional) Group name.
  * `port_monitoring` - (Optional) Port monitoring configuration.

## Import

Authenticated Ethernet Port resources can be imported using the `name` attribute:

```
$ terraform import verity_authenticated_eth_port.example example
```

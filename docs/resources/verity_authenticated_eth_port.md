# Authenticated Eth-Port Resource

`verity_authenticated_eth_port` manages authenticated Ethernet port resources in Verity, which define authentication settings for Ethernet ports.

## Example Usage

```hcl
resource "verity_authenticated_eth_port" "example" {
  name = "example"
  enable = false
  connection_mode = "PortMode"
  reauthorization_period_sec = 3600
  allow_mac_based_authentication = false
  mac_authentication_holdoff_sec = 60
  trusted_port = false

  eth_ports {
    eth_port_profile_num_enable = false
    eth_port_profile_num_eth_port = ""
    eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
    eth_port_profile_num_walled_garden_set = false
    eth_port_profile_num_radius_filter_id = ""
    index = 1
  }

  object_properties {
    group = ""
    port_monitoring = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `connection_mode` (String) - Choose connection mode for Authenticated Eth-Port. Port Mode: Standard mode. The last authenticated clients VLAN access is applied. Single Client Mode: MAC filtered client. Only the authenticated clients traffic can pass. No traffic from a second client may pass. Only when the first client deauthenticates can a new authentication take place. Multiple Client Mode: MAC filtered clients. Only authenticated client traffic can pass. Multiple clients can authenticate and gain access to individual service offerings. MAC-based authentication is not supported.
* `reauthorization_period_sec` (Integer) - Amount of time in seconds before 802.1X requires reauthorization of an active session. "0" disables reauthorization (not recommended).
* `allow_mac_based_authentication` (Boolean) - Enables 802.1x to capture the connected MAC address and send it tothe Radius Server instead of requesting credentials.  Useful for printers and similar devices.
* `mac_authentication_holdoff_sec` (Integer) - Amount of time in seconds 802.1X authentication is allowed to run before MAC-based authentication has begun.
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.
* `eth_ports` (Array) - 
  * `eth_port_profile_num_enable` (Boolean) - Enable row.
  * `eth_port_profile_num_eth_port` (String) - Choose an Eth Port Profile.
  * `eth_port_profile_num_eth_port_ref_type_` (String) - Object type for eth_port_profile_num_eth_port field.
  * `eth_port_profile_num_walled_garden_set` (Boolean) - Flag indicating this Eth Port Profile is the Walled Garden.
  * `eth_port_profile_num_radius_filter_id` (String) - The value of filter-id in the RADIUS response which will evoke this Eth Port Profile.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.

## Import

Authenticated Ethernet Port resources can be imported using the `name` attribute:

```sh
terraform import verity_authenticated_eth_port.<resource_name> <name>
```

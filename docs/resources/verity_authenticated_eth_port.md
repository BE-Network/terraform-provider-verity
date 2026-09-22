# verity_authenticated_eth_port (Resource)

Manages a Verity Authenticated Eth-Port.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_authenticated_eth_port" "example" {
  name = "example"
  allow_mac_based_authentication = false
  connection_mode = ""
  enable = false
  mac_authentication_holdoff_sec = null
  reauthorization_period_sec = null
  trusted_port = false

  eth_ports {
    index = 1
    eth_port_profile_num_enable = false
    eth_port_profile_num_eth_port = ""
    eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
    eth_port_profile_num_radius_filter_id = ""
    eth_port_profile_num_walled_garden_set = false
  }

  object_properties {
    port_monitoring = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `allow_mac_based_authentication` (Boolean) - Enables 802.1x to capture the connected MAC address and send it tothe Radius Server instead of requesting credentials. Useful for printers and similar devices.
* `connection_mode` (String) - Choose connection mode for Authenticated Eth-Port Port Mode Standard mode. The last authenticated clients VLAN access is applied. Single Client Mode MAC filtered client. Only the authenticated clients traffic can pass. No traffic from a second client may pass. Only when the first client deauthenticates can a new authentication take place. Multiple Client Mode MAC filtered clients. Only authenticated client traffic can pass. Multiple clients can authenticate and gain access to individual service offerings. MAC-based authentication is not supported.
* `enable` (Boolean) - Enable object.
* `eth_ports` (Block List) - Ethernet port configurations. Entries are matched by `index`.
  * `eth_port_profile_num_enable` (Boolean) - Enable row.
  * `eth_port_profile_num_eth_port` (String) - Choose an Eth Port Profile. Set together with `eth_port_profile_num_eth_port_ref_type_`.
  * `eth_port_profile_num_eth_port_ref_type_` (String) - Object type for eth_port_profile_num_eth_port field.
  * `eth_port_profile_num_radius_filter_id` (String) - The value of filter-id in the RADIUS response which will evoke this Eth Port Profile.
  * `eth_port_profile_num_walled_garden_set` (Boolean) - Flag indicating this Eth Port Profile is the Walled Garden.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `mac_authentication_holdoff_sec` (Integer) - Amount of time in seconds 802.1X authentication is allowed to run before MAC-based authentication has begun. Set it to `null` to clear it.
* `object_properties` (Block) - Object properties for the authenticated eth-port. At most one block.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.
* `reauthorization_period_sec` (Integer) - Amount of time in seconds before 802.1X requires reauthorization of an active session. "0" disables reauthorization (not recommended). Set it to `null` to clear it.
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `eth_ports.eth_port_profile_num_eth_port` | `eth_port_profile_num_eth_port_ref_type_` | `eth_port_profile_` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_authenticated_eth_port.<resource_name> <name>
```

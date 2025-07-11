# Ethernet Port Profile Resource

`verity_eth_port_profile` manages Ethernet port profiles in Verity, which define configuration templates for network ports.

## Example Usage

```hcl
resource "verity_eth_port_profile" "example" {
  name = "example"
  object_properties {
    group = "_4522_4"
    port_monitoring = "high"
  }
  tenant_slice_managed = false
  services {
    index = 1
    row_num_service_ref_type_ = ""
    row_num_enable = true
    row_num_service = ""
    row_num_external_vlan = null
  }
  services {
    index = 2
    row_num_enable = true
    row_num_service = ""
    row_num_external_vlan = null
    row_num_service_ref_type_ = ""
  }
  enable = true
}
```

## Argument Reference

* `name` - Unique identifier for the Ethernet port profile
* `enable` - Enable this Ethernet port profile. Default is `false`
* `tenant_slice_managed` - Whether this profile is tenant slice managed. Default is `false`
* `tls` - Transparent LAN Service Trunk (available as of API version 6.5)
* `tls_service` - Choose a Service supporting Transparent LAN Service (available as of API version 6.5)
* `tls_service_ref_type_` - Object type for tls_service field (available as of API version 6.5)
* `trusted_port` - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping (available as of API version 6.5)
* `object_properties` - Object properties block
  * `group` - Group name
  * `port_monitoring` - Defines importance of Link Down on this port ("high", "critical", or "")
  * `sort_by_name` - Choose to sort by service name or by order of creation (available as of API version 6.5)
  * `label` - Port Label displayed for ports provisioned with this profile but with no Port Label defined in the endpoint (available as of API version 6.5)
  * `icon` - Port Icon displayed for ports provisioned with this profile but with no Port Icon defined in the endpoint (available as of API version 6.5)
* `services` - List of service configuration blocks
  * `index` - Index value for ordering
  * `row_num_enable` - Enable this service configuration
  * `row_num_service` - Reference to a service resource
  * `row_num_service_ref_type_` - Object type for service reference
  * `row_num_external_vlan` - External VLAN ID. A value of 0 will make the VLAN untagged, while null will use service VLAN
  * `row_num_mac_filter` - Choose an access control list (available as of API version 6.5)
  * `row_num_mac_filter_ref_type_` - Object type for mac_filter reference (available as of API version 6.5)
  * `row_num_lan_iptv` - Denotes a LAN or IPTV service (available as of API version 6.5)

## Import

Ethernet port profile resources can be imported using the `name` attribute:

```
$ terraform import verity_eth_port_profile.example example
```

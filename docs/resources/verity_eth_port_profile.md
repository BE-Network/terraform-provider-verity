# Ethernet Port Profile Resource

`verity_eth_port_profile` manages Ethernet port profiles in Verity, which define configuration templates for network ports.

## Example Usage

```hcl
resource "verity_eth_port_profile" "example" {
  name = "example"
  enable = true
  ingress_acl = ""
  ingress_acl_ref_type_ = ""
  egress_acl = ""
  egress_acl_ref_type_ = ""
  tls = false
  tls_service = ""
  tls_service_ref_type_ = ""
  trusted_port = false

  services {
    index = 1
    row_num_enable = true
    row_num_service = ""
    row_num_service_ref_type_ = ""
    row_num_external_vlan = null
    row_num_ingress_acl = ""
    row_num_ingress_acl_ref_type_ = ""
    row_num_egress_acl = ""
    row_num_egress_acl_ref_type_ = ""
    row_num_mac_filter = ""
    row_num_mac_filter_ref_type_ = ""
    row_num_lan_iptv = ""
  }

  object_properties {
    group = ""
    port_monitoring = ""
    sort_by_name = false
    label = ""
    icon = "empty"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `ingress_acl` (String) - Choose an ingress access control list.
* `ingress_acl_ref_type_` (String) - Object type for ingress_acl field.
* `egress_acl` (String) - Choose an egress access control list.
* `egress_acl_ref_type_` (String) - Object type for egress_acl field.
* `services` (Array) - 
  * `row_num_enable` (Boolean) - Enable row.
  * `row_num_service` (String) - Choose a Service to connect.
  * `row_num_service_ref_type_` (String) - Object type for row_num_service field.
  * `row_num_external_vlan` (Integer) - Choose an external vlan. A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service.
  * `row_num_ingress_acl` (String) - Choose an ingress access control list.
  * `row_num_ingress_acl_ref_type_` (String) - Object type for row_num_ingress_acl field.
  * `row_num_egress_acl` (String) - Choose an egress access control list.
  * `row_num_egress_acl_ref_type_` (String) - Object type for row_num_egress_acl field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `row_num_mac_filter` (String) - Choose an access control list.
  * `row_num_mac_filter_ref_type_` (String) - Object type for row_num_mac_filter field.
  * `row_num_lan_iptv` (String) - Denotes a LAN or IPTV service.
* `object_properties` (Object) - 
  * `group` (String) - Group.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.
  * `sort_by_name` (Boolean) - Choose to sort by service name or by order of creation.
  * `label` (String) - Port Label displayed ports provisioned with this Eth Port Profile but with no Port Label defined in the endpoint.
  * `icon` (String) - Port Icon displayed ports provisioned with this Eth Port Profile but with no Port Icon defined in the endpoint.
* `tls` (Boolean) - Transparent LAN Service Trunk.
* `tls_service` (String) - Choose a Service supporting Transparent LAN Service.
* `tls_service_ref_type_` (String) - Object type for tls_service field.
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.

## Import

Ethernet port profile resources can be imported using the `name` attribute:

```sh
terraform import verity_eth_port_profile.<resource_name> <name>
```

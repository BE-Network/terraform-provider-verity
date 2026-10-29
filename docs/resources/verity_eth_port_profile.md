# verity_eth_port_profile (Resource)

Manages an Ethernet Port Profile.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_eth_port_profile" "example" {
  name = "example"
  egress_acl = ""
  egress_acl_ref_type_ = "port_acl"
  enable = false
  ingress_acl = ""
  ingress_acl_ref_type_ = "port_acl"
  tls = false
  tls_service = ""
  tls_service_ref_type_ = "service"
  trusted_port = false

  object_properties {
    icon = ""
    label = ""
    port_monitoring = ""
    sort_by_name = false
  }

  services {
    index = 1
    row_num_egress_acl = ""
    row_num_egress_acl_ref_type_ = "port_acl"
    row_num_enable = false
    row_num_external_vlan = null
    row_num_ingress_acl = ""
    row_num_ingress_acl_ref_type_ = "port_acl"
    row_num_lan_iptv = ""
    row_num_mac_filter = ""
    row_num_mac_filter_ref_type_ = "mac_filter"
    row_num_service = ""
    row_num_service_ref_type_ = "service"
  }
}
```

### Datacenter mode

```hcl
resource "verity_eth_port_profile" "example" {
  name = "example"
  egress_acl = ""
  egress_acl_ref_type_ = "port_acl"
  enable = false
  ingress_acl = ""
  ingress_acl_ref_type_ = "port_acl"

  object_properties {
    port_monitoring = ""
  }

  services {
    index = 1
    row_num_egress_acl = ""
    row_num_egress_acl_ref_type_ = "port_acl"
    row_num_enable = false
    row_num_external_vlan = null
    row_num_ingress_acl = ""
    row_num_ingress_acl_ref_type_ = "port_acl"
    row_num_service = ""
    row_num_service_ref_type_ = "service"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `egress_acl` (String) - Choose an egress access control list. Set together with `egress_acl_ref_type_`.
* `egress_acl_ref_type_` (String) - Object type for egress_acl field.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `ingress_acl` (String) - Choose an ingress access control list. Set together with `ingress_acl_ref_type_`.
* `ingress_acl_ref_type_` (String) - Object type for ingress_acl field.
* `object_properties` (Block) - Object properties for the profile. At most one block.
  * `icon` (String) - Port Icon displayed ports provisioned with this Eth Port Profile but with no Port Icon defined in the endpoint. Campus mode only.
  * `label` (String) - Port Label displayed ports provisioned with this Eth Port Profile but with no Port Label defined in the endpoint. Campus mode only.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.
  * `sort_by_name` (Boolean) - Choose to sort by service name or by order of creation. Campus mode only.
* `services` (Block List) - List of service configurations. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `row_num_egress_acl` (String) - Choose an egress access control list. Set together with `row_num_egress_acl_ref_type_`.
  * `row_num_egress_acl_ref_type_` (String) - Object type for row_num_egress_acl field.
  * `row_num_enable` (Boolean) - Enable row.
  * `row_num_external_vlan` (Integer) - Choose an external vlan A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service. Set it to `null` to clear it.
  * `row_num_ingress_acl` (String) - Choose an ingress access control list. Set together with `row_num_ingress_acl_ref_type_`.
  * `row_num_ingress_acl_ref_type_` (String) - Object type for row_num_ingress_acl field.
  * `row_num_lan_iptv` (String) - Denotes a LAN or IPTV service. Campus mode only.
  * `row_num_mac_filter` (String) - Choose an access control list. Campus mode only. Set together with `row_num_mac_filter_ref_type_`.
  * `row_num_mac_filter_ref_type_` (String) - Object type for row_num_mac_filter field. Campus mode only.
  * `row_num_service` (String) - Choose a Service to connect. Set together with `row_num_service_ref_type_`.
  * `row_num_service_ref_type_` (String) - Object type for row_num_service field.
* `tls` (Boolean) - Transparent LAN Service Trunk. Campus mode only.
* `tls_service` (String) - Choose a Service supporting Transparent LAN Service. Campus mode only. Set together with `tls_service_ref_type_`.
* `tls_service_ref_type_` (String) - Object type for tls_service field. Campus mode only.
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks. Campus mode only.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `egress_acl` | `egress_acl_ref_type_` | `port_acl` |
| `ingress_acl` | `ingress_acl_ref_type_` | `port_acl` |
| `services.row_num_egress_acl` | `row_num_egress_acl_ref_type_` | `port_acl` |
| `services.row_num_ingress_acl` | `row_num_ingress_acl_ref_type_` | `port_acl` |
| `services.row_num_mac_filter` | `row_num_mac_filter_ref_type_` | `mac_filter` |
| `services.row_num_service` | `row_num_service_ref_type_` | `service` |
| `tls_service` | `tls_service_ref_type_` | `service` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_eth_port_profile.<resource_name> <name>
```

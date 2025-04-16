# Service Resource

`verity_service` manages service resources in Verity, which define network service configurations.

## Example Usage

```hcl
resource "verity_service" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  anycast_ip_mask = ""
  dhcp_server_ip = ""
  mtu = 1500
  tenant = ""
  tenant_ref_type_ = ""
  vni_auto_assigned_ = true
  enable = false
  vlan = 4007
  vni = 104007
}
```

## Argument Reference

* `name` - Unique identifier for the service
* `enable` - Enable this service. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
* `vlan` - VLAN ID for the service
* `vni` - VNI (VXLAN Network Identifier) for the service
* `vni_auto_assigned_` - Whether the VNI value is automatically assigned
* `tenant` - Reference to a tenant resource
* `tenant_ref_type_` - Object type for tenant reference
* `anycast_ip_mask` - Static anycast gateway address for service
* `dhcp_server_ip` - IP address(es) of the DHCP server for service. May have up to four separated by commas.
* `mtu` - MTU (Maximum Transmission Unit) - the size used by a switch to determine when large packets must be broken up for delivery

## Import

Service resources can be imported using the `name` attribute:


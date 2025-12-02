# Service Resource

`verity_service` manages service resources in Verity, which define network service configurations.

## Example Usage

```hcl
resource "verity_service" "example" {
  name = "example"
  enable = false
  vlan = null
  vni = null
  vni_auto_assigned_ = false
  tenant = ""
  tenant_ref_type_ = ""
  anycast_ip_mask = ""
  dhcp_server_ip = ""
  mtu = 1500

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `vlan` (Integer) - A Value between 1 and 4096.
* `vni` (Integer) - Indication of the outgoing VLAN layer 2 service.
* `vni_auto_assigned_` (Boolean) - Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit.
* `tenant` (String) - Tenant.
* `tenant_ref_type_` (String) - Object type for tenant field.
* `anycast_ip_mask` (String) - Static anycast gateway address for service.
* `dhcp_server_ip` (String) - IP address(s) of the DHCP server for service. May have up to four separated by commas.
* `mtu` (Integer) - MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets.
* `object_properties` (Object) - 
  * `group` (String) - Group.

## Import

Service resources can be imported using the `name` attribute:

```sh
terraform import verity_service.<resource_name> <name>
```
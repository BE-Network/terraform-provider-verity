# Ethernet Port Profile Resource

`verity_eth_port_profile` manages Ethernet port profiles in Verity, which define configuration templates for network ports.

## Example Usage

```hcl
resource "verity_eth_port_profile" "example" {
  name = "example"
  enable = false
  tenant_slice_managed = false

  services {
    row_num_enable = false
    row_num_service = ""
    row_num_service_ref_type_ = "service"
    row_num_external_vlan = null
    index = 1
  }

  object_properties {
    group = ""
    port_monitoring = "high"
  }
}
```

## Argument Reference

* `name` (String, Required) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran. 
* `tenant_slice_managed` (Boolean) - Profiles that Tenant Slice creates and manages.
* `services` (Array) - Service configurations for the ethernet port profile.
  * `row_num_enable` (Boolean) - Enable row.
  * `row_num_service` (String) - Choose a Service to connect.
  * `row_num_service_ref_type_` (String) - Object type for row_num_service field.
  * `row_num_external_vlan` (Integer, Nullable) - Choose an external vlan. A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service. Range: 2-4096.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - Additional object properties.
  * `group` (String) - Group.
  * `port_monitoring` (String) - Defines importance of Link Down on this port. Valid values: `""`, `"critical"`, `"high"`.

## Import

Ethernet port profile resources can be imported using the `name` attribute:

```sh
terraform import verity_eth_port_profile.<resource_name> <name>
```

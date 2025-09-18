# Gateway Profile Resource

`verity_gateway_profile` manages gateway profiles in Verity, which define gateway configuration templates.

## Example Usage

```hcl
resource "verity_gateway_profile" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  enable = false
  tenant_slice_managed = false
  
  external_gateways {
    index = 1
    enable = false
    gateway = ""
    source_ip_mask = ""
    peer_gw = false
    gateway_ref_type_ = ""
  }
}
```

## Argument Reference

* `name` - Unique identifier for the gateway profile
* `enable` - Enable this gateway profile. Default is `false`
* `tenant_slice_managed` - Whether this profile is tenant slice managed. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
* `external_gateways` - List of external gateway blocks
  * `enable` - Enable this external gateway
  * `gateway` - Reference to a gateway resource
  * `gateway_ref_type_` - Object type for gateway reference
  * `source_ip_mask` - Source IP mask in CIDR notation
  * `peer_gw` - Whether this is a peer gateway
  * `index` - Index value for ordering

## Import

Gateway profile resources can be imported using the `name` attribute:

```sh
terraform import verity_gateway_profile.<resource_name> <name>
```
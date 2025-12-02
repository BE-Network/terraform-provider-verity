# Gateway Profile Resource

`verity_gateway_profile` manages gateway profiles in Verity, which define gateway configuration templates.

## Example Usage

```hcl
resource "verity_gateway_profile" "example" {
  name = "example"
  enable = false
  tenant_slice_managed = false
  
  external_gateways {
    index = 1
    enable = false
    gateway = ""
    gateway_ref_type_ = "gateway"
    source_ip_mask = ""
    peer_gw = false
  }

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String, Required) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `tenant_slice_managed` (Boolean) - Profiles that Tenant Slice creates and manages.
* `external_gateways` (Array) - External gateway configurations.
  * `enable` (Boolean) - Enable row.
  * `gateway` (String) - BGP Gateway referenced for port profile.
  * `gateway_ref_type_` (String) - Object type for gateway field.
  * `source_ip_mask` (String) - Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session.
  * `peer_gw` (Boolean) - Setting for paired switches only. Flag indicating that this gateway is a peer gateway. For each gateway profile referencing a BGP session on a member of a leaf pair, the peer should have a gateway profile entry indicating the IP address for the peers gateway.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - Additional object properties.
  * `group` (String) - Group.

## Import

Gateway profile resources can be imported using the `name` attribute:

```sh
terraform import verity_gateway_profile.<resource_name> <name>
```
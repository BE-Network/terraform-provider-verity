# ipv6 Prefix List Resource

Provides a Verity IPv6 Prefix List resource. IPv6 Prefix Lists are used to match and filter IPv6 routes based on prefix and mask criteria.

## Example Usage

```hcl
resource "verity_ipv6_prefix_list" "example" {
  name = "example"
  enable = false

  lists {
    index = 1
    enable = false
    greater_than_equal_value = null
    ipv6_prefix = ""
    less_than_equal_value = null
    permit_deny = "permit"
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `lists` (Array) - 
  * `enable` (Boolean) - Enable of this IPv6 Prefix List.
  * `permit_deny` (String) - Action upon match of Community Strings.
  * `ipv6_prefix` (String) - IPv6 address and subnet to match against.
  * `greater_than_equal_value` (Integer) - Match IP routes with a subnet mask greater than or equal to the value indicated.
  * `less_than_equal_value` (Integer) - Match IP routes with a subnet mask less than or equal to the value indicated.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

IPv6 Prefix Lists can be imported using the name:

```sh
terraform import verity_ipv6_prefix_list.<resource_name> <name>
```

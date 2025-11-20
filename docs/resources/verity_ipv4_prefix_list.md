# ipv4 Prefix List Resource

Provides a Verity IPv4 Prefix List resource. IPv4 Prefix Lists are used to match and filter IPv4 routes based on prefix and mask criteria.

## Example Usage

```hcl
resource "verity_ipv4_prefix_list" "example" {
  name = "example"
  enable = false

  lists {
    index = 1
    enable = false
    greater_than_equal_value = null
    ipv4_prefix = ""
    less_than_equal_value = null
    permit_deny = "permit"
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `lists` (Array) - 
  * `enable` (Boolean) - Enable of this IPv4 Prefix List.
  * `permit_deny` (String) - Action upon match of Community Strings.
  * `ipv4_prefix` (String) - IPv4 address and subnet to match against.
  * `greater_than_equal_value` (Integer) - Match IP routes with a subnet mask greater than or equal to the value indicated.
  * `less_than_equal_value` (Integer) - Match IP routes with a subnet mask less than or equal to the value indicated.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

IPv4 Prefix Lists can be imported using the name:

```sh
terraform import verity_ipv4_prefix_list.<resource_name> <name>
```

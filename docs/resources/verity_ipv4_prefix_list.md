# verity_ipv4_prefix_list (Resource)

Manages a Verity IPv4 Prefix List.

Supported modes: Datacenter.

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
    permit_deny = ""
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `lists` (Block List) - List of IPv4 Prefix List entries. Entries are matched by `index`.
  * `enable` (Boolean) - Enable of this IPv4 Prefix List.
  * `greater_than_equal_value` (Integer) - Match IP routes with a subnet mask greater than or equal to the value indicated. Set it to `null` to clear it.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `ipv4_prefix` (String) - IPv4 address and subnet to match against.
  * `less_than_equal_value` (Integer) - Match IP routes with a subnet mask less than or equal to the value indicated. Set it to `null` to clear it.
  * `permit_deny` (String) - Action upon match of Community Strings.
* `object_properties` (Block) - Object properties for the IPv4 Prefix List. At most one block.
  * `notes` (String) - User Notes.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_ipv4_prefix_list.<resource_name> <name>
```

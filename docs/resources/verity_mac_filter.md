# verity_mac_filter (Resource)

Manages a Verity MAC Filter.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_mac_filter" "example" {
  name = "example"
  enable = false
  type = ""

  filters {
    index = 1
    filter_num_enable = false
    filter_num_mac = ""
    filter_num_mask = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `filters` (Block List) - List of MAC filter entries. Entries are matched by `index`.
  * `filter_num_enable` (Boolean) - Enable of this MAC Filter.
  * `filter_num_mac` (String) - MAC address descriptor including colons example 01:23:45:67:9a:ab. and * notation accepted example 12:*.
  * `filter_num_mask` (String) - Hexidecimal mask including colons example ff:ff:fe:00:00:00. /n and * notation accepted example /16 or 12:*.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `type` (String) - Black vs White MAC Filter.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_mac_filter.<resource_name> <name>
```

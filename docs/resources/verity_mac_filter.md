# MAC Filter Resource

`verity_mac_filter` manages MAC filter resources in Verity.

## Example Usage

```hcl
resource "verity_mac_filter" "example" {
  name = "example"
  enable = true
  type = "White"

  filters {
    index = 1
    filter_num_mac = "01:23:45:67:89:ab"
    filter_num_mask = "ff:ff:ff:ff:ff:ff"
    filter_num_enable = true
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `type` (String) - Black vs White MAC Filter.
* `filters` (Array) -
  * `filter_num_mac` (String) - MAC address descriptor including colons example `01:23:45:67:9a:ab`. `*` notation is also accepted, for example `12:*`.
  * `filter_num_mask` (String) - Hexidecimal mask including colons example `ff:ff:fe:00:00:00`. `/n` and `*` notation are also accepted, for example `/16` or `12:*`.
  * `filter_num_enable` (Boolean) - Enable of this MAC Filter.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

MAC filter resources can be imported using the `name` attribute:

```sh
terraform import verity_mac_filter.<resource_name> <name>
```

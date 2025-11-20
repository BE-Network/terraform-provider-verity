# SFB Breakout Resource

Manages an SFP Breakout object in the Verity system.

**Note:** Only PATCH operations are supported for this resource. Create and delete operations are prohibited.

## Example Usage

```hcl
resource "verity_sfp_breakout" "example" {
    name = "example"
    enable = true

    breakout {
        index = 1
        breakout = "1x100G"
        enable = false
        part_number = ""
        vendor = ""
    }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `breakout` (Array) - Defines breakout configuration for the SFP port
  * `enable` (Boolean) - Enable
  * `vendor` (String) - Vendor
  * `part_number` (String) - Part Number
  * `breakout` (String) - Breakout definition; defines number of ports of what speed this port is brokenout to
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list

## Import

This resource can be imported using the object name:

```sh
terraform import verity_sfp_breakout.<resource_name> <name>
```

# verity_sfp_breakout (Resource)

Manages a Verity SFP Breakout.

Supported modes: Campus, Datacenter.

This resource can only be updated. The object always exists on the Verity system, so it cannot be created or destroyed through Terraform; import it, then manage its settings.

## Example Usage

```hcl
resource "verity_sfp_breakout" "example" {
  name = "example"

  breakout {
    index = 1
    breakout = ""
    enable = false
    part_number = ""
    vendor = ""
  }

  object_properties {
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object Name. Must be unique. Changing it replaces the resource.

### Optional

* `breakout` (Block List) - List of breakout configurations. Entries are matched by `index`.
  * `breakout` (String) - Breakout definition; defines number of ports of what speed this port is brokenout to.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `part_number` (String) - Part Number.
  * `vendor` (String) - Vendor.
* `object_properties` (Block) - Object properties. At most one block.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_sfp_breakout.<resource_name> <name>
```

# verity_badge (Resource)

Manages a Badge resource.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_badge" "example" {
  name = "example"
  color = ""
  enable = false
  number = null

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `color` (String) - Color of Badge.
* `enable` (Boolean) - Enable object.
* `number` (Integer) - Number of Badge. Set it to `null` to clear it.
* `object_properties` (Block) - Object properties for the badge. At most one block.
  * `notes` (String) - User Notes.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_badge.<resource_name> <name>
```

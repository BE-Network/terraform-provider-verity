# Badge Resource

`verity_badge` manages badge resources in Verity, which define identification badges with colors and numbers.


## Example Usage

```hcl
resource "verity_badge" "example" {
  name = "example"
  enable = true
  color = "red"
  number = 1

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `color` (String) - Color of Badge.
* `number` (Integer) - Number of Badge.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

Badge resources can be imported using the `name` attribute:

```sh
terraform import verity_badge.<resource_name> <name>
```

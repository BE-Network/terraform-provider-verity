# Badge Resource

`verity_badge` manages badge resources in Verity, which define identification badges with colors and numbers.

## Version Compatibility

**This resource requires Verity API version 6.5 or higher.**

## Example Usage

```hcl
resource "verity_badge" "example" {
  name = "example_badge"
  color = "blue"
  number = 42
  
  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the badge.
* `color` - (Optional) Badge color.
* `number` - (Optional) Badge number.
* `object_properties` - (Optional) Object properties configuration:
  * `notes` - (Optional) User notes for the badge.

## Import

Badge resources can be imported using the `name` attribute:

```
$ terraform import verity_badge.example example
```

# Rack Resource

`verity_rack` manages Rack resources in Verity.

## Example Usage

```hcl
resource "verity_rack" "example" {
  name         = "rack-a"
  enable       = true
  position     = 1
  su           = "su-a"
  su_ref_type_ = "su"

  object_properties {
    notes = "Example rack"
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `position` (Number) - Position of the Rack.
* `su` (String) - SU this Rack is assigned to.
* `su_ref_type_` (String) - Object type for `su` field.
* `object_properties` (Object) -
  * `notes` (String) - User Notes.

## Import

```sh
terraform import verity_rack.<resource_name> <name>
```

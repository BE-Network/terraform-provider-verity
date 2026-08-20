# SSP Group Resource

`verity_ssp_group` manages SuperSpine Group resources in Verity.

## Example Usage

```hcl
resource "verity_ssp_group" "example" {
  name = "example"
  enable = true
  fabric = "fabric-a"
  fabric_ref_type_ = "fabric"
  position = 1

  object_properties {
    notes = "Example SSP group"
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `fabric` (String) - Fabric this SuperSpine Group is assigned to.
* `fabric_ref_type_` (String) - Object type for `fabric` field.
* `position` (Number) - Position of the Switch.
* `object_properties` (Object) -
  * `notes` (String) - User Notes.

## Import

SSP group resources can be imported using the `name` attribute:

```sh
terraform import verity_ssp_group.<resource_name> <name>
```

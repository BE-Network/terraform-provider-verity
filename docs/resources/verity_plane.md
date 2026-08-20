# Plane Resource

`verity_plane` manages Plane resources in Verity.

## Example Usage

```hcl
resource "verity_plane" "example" {
  name           = "plane-a"
  enable         = true
  fabric           = "fabric-a"
  fabric_ref_type_ = "fabric"
  position       = 1

  object_properties {
    notes = "Example plane"
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `fabric` (String) - Fabric this Plane is assigned to.
* `fabric_ref_type_` (String) - Object type for `fabric` field.
* `position` (Number) - Position of the Plane.
* `object_properties` (Object) -
  * `notes` (String) - User Notes.

## Import

```sh
terraform import verity_plane.<resource_name> <name>
```

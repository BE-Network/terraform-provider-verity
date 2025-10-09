# Spine Plane Resource

Manages a Spine Plane configuration in Verity. This resource allows you to define named spine plane objects, which can be enabled or disabled and may include additional properties.

## Example Usage

```hcl
resource "verity_spine_plane" "example" {
  name   = "spine_plane_test1"
  enable = true

  object_properties = {
    notes = ""
  }
}
```

## Argument Reference

- `name` (String, Required): The unique name of the spine plane object.
- `enable` (Boolean, Required): Whether this spine plane is enabled.
- `object_properties` (Map(String), Optional): Additional properties for the spine plane. Common fields include:
  - `notes` (String, Optional): Free-form notes or description for the spine plane.


## Import

This resource can be imported using the spine plane name:

```sh
terraform import verity_spine_plane.<resource_name> <name>
```

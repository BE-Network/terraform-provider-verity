# Spine Plane Resource

Manages a Spine Plane configuration in Verity. This resource allows you to define named spine plane objects, which can be enabled or disabled and may include additional properties.

## Example Usage

```hcl
resource "verity_spine_plane" "example" {
  name   = "example"
  enable = true

  object_properties = {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `object_properties` (Object) - Additional properties for the spine plane
  * `notes` (String) - User Notes


## Import

This resource can be imported using the spine plane name:

```sh
terraform import verity_spine_plane.<resource_name> <name>
```

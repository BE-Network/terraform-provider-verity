# verity_spine_plane (Resource)

Manages a Spine Plane resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_spine_plane" "example" {
  name = "example"
  enable = false
  fabric = ""
  fabric_ref_type_ = "fabric"

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `fabric` (String) - Fabric this Spine Plane is assigned to. Set together with `fabric_ref_type_`.
* `fabric_ref_type_` (String) - Object type for fabric field.
* `object_properties` (Block) - Object properties for the spine plane. At most one block.
  * `notes` (String) - User Notes.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `fabric` | `fabric_ref_type_` | `fabric` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_spine_plane.<resource_name> <name>
```

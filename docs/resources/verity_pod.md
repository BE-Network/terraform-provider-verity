# verity_pod (Resource)

Manages a Pod resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_pod" "example" {
  name = "example"
  enable = false
  expected_spine_count = null
  fabric = ""
  fabric_ref_type_ = "fabric"
  position = null

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
* `expected_spine_count` (Integer) - Number of spine switches expected in this pod. Set it to `null` to clear it.
* `fabric` (String) - Fabric this Pod is assigned to. Set together with `fabric_ref_type_`.
* `fabric_ref_type_` (String) - Object type for fabric field.
* `object_properties` (Block) - Object properties for the pod. At most one block.
  * `notes` (String) - User Notes.
* `position` (Number) - Position of the Switch. Set it to `null` to clear it.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `fabric` | `fabric_ref_type_` | `fabric` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_pod.<resource_name> <name>
```

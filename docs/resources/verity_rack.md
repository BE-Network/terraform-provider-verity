# verity_rack (Resource)

Manages a Verity Rack resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_rack" "example" {
  name = "example"
  enable = false
  position = null
  su = ""
  su_ref_type_ = "su"

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
* `object_properties` (Block) - Object properties for the Rack. At most one block.
  * `notes` (String) - User Notes.
* `position` (Number) - Position of the Rack. Set it to `null` to clear it.
* `su` (String) - SU this Rack is assigned to. Set together with `su_ref_type_`.
* `su_ref_type_` (String) - Object type for su field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `su` | `su_ref_type_` | `su` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_rack.<resource_name> <name>
```

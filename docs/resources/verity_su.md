# verity_su (Resource)

Manages a Verity SU resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_su" "example" {
  name = "example"
  enable = false
  pod = ""
  pod_ref_type_ = "pod"
  position = null

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object name. Identifies the object in the API and must be unique within its type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `object_properties` (Block) - Object properties for the SU. At most one block.
  * `notes` (String) - User Notes.
* `pod` (String) - Pod this SU is assigned to. Set together with `pod_ref_type_`.
* `pod_ref_type_` (String) - Object type for pod field.
* `position` (Number) - Position of the Switch. Set it to `null` to clear it.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `pod` | `pod_ref_type_` | `pod` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_su.<resource_name> <name>
```

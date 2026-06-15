# SU Resource

`verity_su` manages SU resources in Verity.

## Example Usage

```hcl
resource "verity_su" "example" {
  name = "example"
  enable = true
  pod = "pod-a"
  pod_ref_type_ = "pod"
  position = 1

  object_properties {
    notes = "Example SU"
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `pod` (String) - Pod this SU is assigned to.
* `pod_ref_type_` (String) - Object type for `pod` field.
* `position` (Number) - Position of the Switch.
* `object_properties` (Object) -
  * `notes` (String) - User Notes.

## Import

SU resources can be imported using the `name` attribute:

```sh
terraform import verity_su.<resource_name> <name>
```

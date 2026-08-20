# Pod Resource

Provides a Verity Pod resource. Pods are logical groupings used for network segmentation and management.

## Example Usage

```hcl
resource "verity_pod" "example" {
  name = "example"
  enable = true
  expected_spine_count = 2
  fabric = ""
  fabric_ref_type_ = "fabric"
  position = 0

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `fabric` (String) - Fabric this Pod is assigned to.
* `fabric_ref_type_` (String) - Object type for fabric field.
* `position` (Number) - Position of the Switch.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.
* `expected_spine_count` (Integer) - Number of spine switches expected in this pod.

## Import

Pods can be imported using the name:

```sh
terraform import verity_pod.<resource_name> <name>
```

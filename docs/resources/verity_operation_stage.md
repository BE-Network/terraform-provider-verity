# Operation Stage Resource

`verity_operation_stage` represents a stage in the operation sequence, used to establish dependencies between resource types in Verity.

## Example Usage

```hcl
resource "verity_operation_stage" "example" {
  # No configuration required
}
```

## Argument Reference

This resource does not require any configuration arguments.

## Attributes Reference

* `id` - The unique identifier for this stage.

## Import

Operation Stage resources can be imported using the `id` attribute:

```
$ terraform import verity_operation_stage.example_stage example_stage
```

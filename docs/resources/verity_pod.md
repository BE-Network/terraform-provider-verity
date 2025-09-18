# Pod Resource

Provides a Verity Pod resource. Pods are logical groupings used for network segmentation and management.

## Example Usage

```hcl
resource "verity_pod" "Pod_1" {
  name         = "Pod 1"
  depends_on   = [verity_operation_stage.pod_stage]
  enable       = true
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `true`.

## Import

Pods can be imported using the name:

```sh
terraform import verity_pod.<resource_name> <name>
```

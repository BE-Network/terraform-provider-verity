# Pair Resource

`verity_pair` manages switch pair resources in Verity.

## Example Usage

```hcl
resource "verity_pair" "example" {
  name = "example"
  enable = true
  switchpoint_1 = "switch-a"
  switchpoint_1_ref_type_ = "switchpoint"
  switchpoint_2 = "switch-b"
  switchpoint_2_ref_type_ = "switchpoint"
  lag = "lag-1"
  lag_ref_type_ = "lag"
  is_whitebox_pair = false
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `switchpoint_1` (String) - Switchpoint.
* `switchpoint_1_ref_type_` (String) - Object type for `switchpoint_1` field.
* `switchpoint_2` (String) - Switchpoint.
* `switchpoint_2_ref_type_` (String) - Object type for `switchpoint_2` field.
* `lag` (String) - LAG.
* `lag_ref_type_` (String) - Object type for `lag` field.
* `is_whitebox_pair` (Boolean) - Is Whitebox Pair.

## Import

Pair resources can be imported using the `name` attribute:

```sh
terraform import verity_pair.<resource_name> <name>
```

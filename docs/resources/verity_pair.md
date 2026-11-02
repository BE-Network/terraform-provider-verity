# verity_pair (Resource)

Manages a Verity Switch Pair.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_pair" "example" {
  name = "example"
  enable = false
  is_whitebox_pair = false
  lag = ""
  lag_ref_type_ = "lag"
  switchpoint_1 = ""
  switchpoint_1_ref_type_ = "switchpoint"
  switchpoint_2 = ""
  switchpoint_2_ref_type_ = "switchpoint"
}
```

## Argument Reference

### Required

* `name` (String) - Object name. Identifies the object in the API and must be unique within its type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `is_whitebox_pair` (Boolean) - Is Whitebox Pair.
* `lag` (String) - LAG. Set together with `lag_ref_type_`.
* `lag_ref_type_` (String) - Object type for lag field.
* `switchpoint_1` (String) - Switchpoint. Set together with `switchpoint_1_ref_type_`.
* `switchpoint_1_ref_type_` (String) - Object type for switchpoint_1 field.
* `switchpoint_2` (String) - Switchpoint. Set together with `switchpoint_2_ref_type_`.
* `switchpoint_2_ref_type_` (String) - Object type for switchpoint_2 field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `lag` | `lag_ref_type_` | `lag` |
| `switchpoint_1` | `switchpoint_1_ref_type_` | `switchpoint` |
| `switchpoint_2` | `switchpoint_2_ref_type_` | `switchpoint` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_pair.<resource_name> <name>
```

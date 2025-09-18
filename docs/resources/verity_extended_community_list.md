# Extended Community List Resource

Provides a Verity Extended Community List resource. Extended Community Lists are used to match BGP extended communities for route filtering and policy control.

## Example Usage

```hcl
resource "verity_extended_community_list" "test1" {
  name         = "test1"
  depends_on   = [verity_operation_stage.extended_community_list_stage]
  object_properties {
    notes = ""
  }
  any_all      = "any"
  enable       = false
  lists {
    index                         = 1
    enable                        = false
    mode                          = "route"
    route_target_expanded_expression = ""
  }
  permit_deny   = "permit"
  standard_expanded = "standard"
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `permit_deny` (String, Optional) — Action upon match of Community Strings. Allowed values: `"permit"`, `"deny"`. Default: `"permit"`.
- `any_all` (String, Optional) — BGP does not advertise any or all routes that do not match the Community String. Allowed values: `"any"`, `"all"`. Default: `"any"`.
- `standard_expanded` (String, Optional) — Used Community String or Expanded Expression. Allowed values: `"standard"`, `"expanded"`. Default: `"standard"`.
- `lists` (Block, Optional) — List of extended community list entries:
  - `index` (Integer, Optional) — The index identifying the object. Zero if you want to add an object to the list.
  - `enable` (Boolean, Optional) — Enable of this Extended Community List. Default: `false`.
  - `mode` (String, Optional) — Mode. Allowed values: `"route"`, `"soo"`. Default: `"route"`.
  - `route_target_expanded_expression` (String, Optional) — Match against a BGP extended community of type Route Target. Default: `""`.
- `object_properties` (Block, Optional) —
  - `notes` (String, Optional) — User Notes. Default: `""`.

## Import

Extended Community Lists can be imported using the name:

```sh
terraform import verity_extended_community_list.<resource_name> <name>
```

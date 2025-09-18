# ipv4 Prefix List Resource

Provides a Verity IPv4 Prefix List resource. IPv4 Prefix Lists are used to match and filter IPv4 routes based on prefix and mask criteria.

## Example Usage

```hcl
resource "verity_ipv4_prefix_list" "test1" {
  name         = "test1"
  depends_on   = [verity_operation_stage.ipv4_prefix_list_stage]
  object_properties {
    notes = ""
  }
  enable       = false
  lists {
    index                    = 1
    enable                   = false
    greater_than_equal_value = null
    ipv4_prefix              = ""
    less_than_equal_value    = null
    permit_deny              = "permit"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `lists` (Block, Optional) — List of IPv4 prefix list entries:
  - `index` (Integer, Optional) — The index identifying the object. Zero if you want to add an object to the list.
  - `enable` (Boolean, Optional) — Enable of this IPv4 Prefix List. Default: `false`.
  - `permit_deny` (String, Optional) — Action upon match of Community Strings. Allowed values: `"permit"`, `"deny"`. Default: `"permit"`.
  - `ipv4_prefix` (String, Optional) — IPv4 address and subnet to match against. Default: `""`.
  - `greater_than_equal_value` (Integer, Optional) — Match IP routes with a subnet mask greater than or equal to the value indicated. Maximum: `32`.
  - `less_than_equal_value` (Integer, Optional) — Match IP routes with a subnet mask less than or equal to the value indicated. Maximum: `32`.
- `object_properties` (Block, Optional) —
  - `notes` (String, Optional) — User Notes. Default: `""`.

## Import

IPv4 Prefix Lists can be imported using the name:

```hcl
terraform import verity_ipv4_prefix_list.<resource_name> <name>
```

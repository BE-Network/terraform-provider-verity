# ipv6 Prefix List Resource

Provides a Verity IPv6 Prefix List resource. IPv6 Prefix Lists are used to match and filter IPv6 routes based on prefix and mask criteria.

## Example Usage

```hcl
resource "verity_ipv6_prefix_list" "test1" {
  name         = "test1"
  depends_on   = [verity_operation_stage.ipv6_prefix_list_stage]
  object_properties {
    notes = ""
  }
  enable       = false
  lists {
    index                    = 1
    enable                   = false
    greater_than_equal_value = null
    ipv6_prefix              = ""
    less_than_equal_value    = null
    permit_deny              = "permit"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `lists` (Block, Optional) — List of IPv6 prefix list entries:
  - `index` (Integer, Optional) — The index identifying the object. Zero if you want to add an object to the list.
  - `enable` (Boolean, Optional) — Enable of this IPv6 Prefix List. Default: `false`.
  - `permit_deny` (String, Optional) — Action upon match of Community Strings. Allowed values: `"permit"`, `"deny"`. Default: `"permit"`.
  - `ipv6_prefix` (String, Optional) — IPv6 address and subnet to match against. Default: `""`.
  - `greater_than_equal_value` (Integer, Optional) — Match IP routes with a subnet mask greater than or equal to the value indicated. Minimum: `1`. Maximum: `128`.
  - `less_than_equal_value` (Integer, Optional) — Match IP routes with a subnet mask less than or equal to the value indicated. Minimum: `1`. Maximum: `128`.
- `object_properties` (Block, Optional) —
  - `notes` (String, Optional) — User Notes. Default: `""`.

## Import

IPv6 Prefix Lists can be imported using the name:

```hcl
terraform import verity_ipv6_prefix_list.<resource_name> <name>
```

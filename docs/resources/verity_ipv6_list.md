# ipv6 List Resource

Provides a Verity IPv6 List resource. IPv6 Lists are used to define sets of IPv6 addresses for filtering and policy application.

## Example Usage

```hcl
resource "verity_ipv6_list" "ipv6_test1" {
  name         = "ipv6_test1"
  depends_on   = [verity_operation_stage.ipv6_list_stage]
  enable       = false
  ipv6_list    = ""
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `ipv6_list` (String, Optional) — Comma separated list of IPv6 addresses. Default: `""`.

## Import

IPv6 Lists can be imported using the name:

```sh
terraform import verity_ipv6_list.<resource_name> <name>
```

# ipv4 List Resource

Provides a Verity IPv4 List resource. IPv4 Lists are used to define sets of IPv4 addresses for filtering and policy application.

## Example Usage

```hcl
resource "verity_ipv4_list" "ipv4_test1" {
  name         = "ipv4_test1"
  depends_on   = [verity_operation_stage.ipv4_list_stage]
  enable       = false
  ipv4_list    = ""
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `ipv4_list` (String, Optional) — Comma separated list of IPv4 addresses. Default: `""`.

## Import

IPv4 Lists can be imported using the name:

```sh
terraform import verity_ipv4_list.<resource_name> <name>
```

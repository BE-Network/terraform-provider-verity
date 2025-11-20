# ipv6 List Resource

Provides a Verity IPv6 List resource. IPv6 Lists are used to define sets of IPv6 addresses for filtering and policy application.

## Example Usage

```hcl
resource "verity_ipv6_list" "example" {
  name = "example"
  enable = false
  ipv6_list = ""
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `ipv6_list` (String) - Comma separated list of IPv6 addresses.

## Import

IPv6 Lists can be imported using the name:

```sh
terraform import verity_ipv6_list.<resource_name> <name>
```

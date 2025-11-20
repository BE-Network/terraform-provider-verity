# ipv4 List Resource

Provides a Verity IPv4 List resource. IPv4 Lists are used to define sets of IPv4 addresses for filtering and policy application.

## Example Usage

```hcl
resource "verity_ipv4_list" "example" {
  name = "example"
  enable = false
  ipv4_list = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `ipv4_list` (String) - Comma separated list of IPv4 addresses.

## Import

IPv4 Lists can be imported using the name:

```sh
terraform import verity_ipv4_list.<resource_name> <name>
```

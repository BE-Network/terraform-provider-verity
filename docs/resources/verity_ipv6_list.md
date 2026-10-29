# verity_ipv6_list (Resource)

Manages a Verity IPv6 List Filter.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_ipv6_list" "example" {
  name = "example"
  enable = false
  ipv6_list = ""
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `ipv6_list` (String) - Comma separated list of IPv6 addresses.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_ipv6_list.<resource_name> <name>
```

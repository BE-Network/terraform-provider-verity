# verity_ipv4_list (Resource)

Manages a Verity IPv4 List Filter.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_ipv4_list" "example" {
  name = "example"
  enable = false
  ipv4_list = ""
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `ipv4_list` (String) - Comma separated list of IPv4 addresses.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_ipv4_list.<resource_name> <name>
```

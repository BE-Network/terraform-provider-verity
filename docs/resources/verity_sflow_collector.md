# verity_sflow_collector (Resource)

Manages a Verity SFlow Collector.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_sflow_collector" "example" {
  name = "example"
  enable = false
  ip = ""
  port = null
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `ip` (String) - IP address of the sFlow Collector.
* `port` (Integer) - Port. Set it to `null` to clear it.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_sflow_collector.<resource_name> <name>
```

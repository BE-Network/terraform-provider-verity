# Sflow Collector Resource

Manages a sFlow Collector object in the Verity system.

## Example Usage

```hcl
resource "verity_sflow_collector" "example" {
    name = "example"
    enable = false
    ip = ""
    port = 6343
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `ip` (String) - IP address of the sFlow Collector
* `port` (Integer) - Port

## Import

This resource can be imported using the object name:

```sh
terraform import verity_sflow_collector.<resource_name> <name>
```

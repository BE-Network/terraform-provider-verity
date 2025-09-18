# Sflow Collector Resource

Manages a sFlow Collector object in the Verity system.

## Example Usage

```hcl
resource "verity_sflow_collector" "sflow_collector_name" {
    name        = "sflow_collector_name"
    depends_on  = [verity_operation_stage.sflow_collector_stage]
    enable      = false
    ip          = ""
    port        = 6343
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `ip` (String, Optional) — IP address of the sFlow Collector. Default: `""`.
- `port` (Integer, Optional) — Port. Default: `6343`. Maximum: `65535`.

## Import

This resource can be imported using the object name:

```sh
terraform import verity_sflow_collector.<resource_name> <name>
```

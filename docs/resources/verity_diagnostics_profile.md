# Diagnostics Profile Resource

Provides a Verity Diagnostics Profile resource. Diagnostics Profiles are used to configure sFlow and polling settings for network diagnostics.

## Example Usage

```hcl
resource "verity_diagnostics_profile" "diagnostics_profile1" {
  name                      = "diagnostics_profile1"
  depends_on                = [verity_operation_stage.diagnostics_profile_stage]
  enable                    = false
  enable_sflow              = true
  flow_collector            = ""
  flow_collector_ref_type_  = "sflow_collector"
  poll_interval             = 20
  vrf_type                  = "management"
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `enable_sflow` (Boolean, Optional) — Enable sFlow for this Diagnostics Profile. Default: `false`.
- `flow_collector` (String, Optional) — Flow Collector for this Diagnostics Profile. Default: `""`.
- `flow_collector_ref_type_` (String, Optional) — Object type for flow_collector field. Allowed value: `"sflow_collector"`.
- `poll_interval` (Integer, Optional) — The sampling rate for sFlow polling (seconds). Default: `20`. Minimum: `5`. Maximum: `300`.
- `vrf_type` (String, Optional) — Management or Underlay. Allowed values: `"management"`, `"underlay"`. Default: `"management"`.

## Import

Diagnostics Profiles can be imported using the name:

```sh
terraform import verity_diagnostics_profile.<resource_name> <name>
```

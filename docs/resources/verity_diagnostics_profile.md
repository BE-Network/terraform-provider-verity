# Diagnostics Profile Resource

Provides a Verity Diagnostics Profile resource. Diagnostics Profiles are used to configure sFlow and polling settings for network diagnostics.

## Example Usage

```hcl
resource "verity_diagnostics_profile" "example" {
  name = "example"
  enable = false
  enable_sflow = true
  flow_collector = ""
  flow_collector_ref_type_ = "sflow_collector"
  poll_interval = 20
  vrf_type = "management"
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `enable_sflow` (Boolean) - Enable sFlow for this Diagnostics Profile.
* `flow_collector` (String) - Flow Collector for this Diagnostics Profile.
* `flow_collector_ref_type_` (String) - Object type for flow_collector field.
* `poll_interval` (Integer) - The sampling rate for sFlow polling (seconds).
* `vrf_type` (String) - Management or Underlay.

## Import

Diagnostics Profiles can be imported using the name:

```sh
terraform import verity_diagnostics_profile.<resource_name> <name>
```

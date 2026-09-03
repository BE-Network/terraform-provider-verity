# Diagnostics Profile Resource

Provides a Verity Diagnostics Profile resource. Diagnostics Profiles are used to configure sFlow and polling settings for network diagnostics.

## Example Usage

```hcl
resource "verity_diagnostics_profile" "example" {
  name = "example"
  enable = false
  enable_sflow = true
  use_internal_collector = false
  flow_collector = ""
  flow_collector_ref_type_ = "sflow_collector"
  poll_interval = 20
  vrf_type = "management"
  monitoring_acl = "example-monitoring-acl"
  monitoring_acl_ref_type_ = "monitoring_acl"
  erspan_destination_ip = "192.0.2.10"
  erspan_dscp = 0
  erspan_ttl = 64
  erspan_gre_type = "0x88BE"
  erspan_queue = 0
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `enable_sflow` (Boolean) - Enable sFlow for this Diagnostics Profile.
* `use_internal_collector` (Boolean) - Use the internal Collector as the flow collector.
* `flow_collector` (String) - Flow Collector for this Diagnostics Profile.
* `flow_collector_ref_type_` (String) - Object type for flow_collector field.
* `poll_interval` (Integer) - The sampling rate for sFlow polling (seconds).
* `vrf_type` (String) - Management or Underlay.
* `monitoring_acl` (String) - Monitoring ACL whose Service VLANs are mirrored to the ERSPAN destination.
* `monitoring_acl_ref_type_` (String) - Object type for `monitoring_acl` field.
* `erspan_destination_ip` (String) - IPv4 address of the remote ERSPAN collector.
* `erspan_dscp` (Integer) - DSCP value for ERSPAN packets (0-63).
* `erspan_ttl` (Integer) - Time-to-live value for ERSPAN packets (0-255).
* `erspan_gre_type` (String) - GRE protocol type as a `0x`-prefixed hexadecimal value.
* `erspan_queue` (Integer) - Output queue for ERSPAN packets (0-63).

## Import

Diagnostics Profiles can be imported using the name:

```sh
terraform import verity_diagnostics_profile.<resource_name> <name>
```

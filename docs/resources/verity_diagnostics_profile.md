# verity_diagnostics_profile (Resource)

Manages a Verity Diagnostics Profile.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_diagnostics_profile" "example" {
  name = "example"
  enable = false
  enable_sflow = false
  erspan_destination_ip = ""
  erspan_dscp = null
  erspan_gre_type = ""
  erspan_queue = null
  erspan_ttl = null
  flow_collector = ""
  flow_collector_ref_type_ = "sflow_collector"
  monitoring_acl = ""
  monitoring_acl_ref_type_ = "monitoring_acl"
  poll_interval = null
  use_internal_collector = false
  vrf_type = ""
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `enable_sflow` (Boolean) - Enable sFlow for this Diagnostics Profile.
* `erspan_destination_ip` (String) - IPv4 address of the remote ERSPAN collector.
* `erspan_dscp` (Integer) - DSCP value for ERSPAN packets (0-63). Set it to `null` to clear it.
* `erspan_gre_type` (String) - GRE protocol type as a 0x-prefixed hexadecimal value.
* `erspan_queue` (Integer) - Output queue for ERSPAN packets (0-63). Set it to `null` to clear it.
* `erspan_ttl` (Integer) - Time-to-live value for ERSPAN packets (0-255). Set it to `null` to clear it.
* `flow_collector` (String) - Flow Collector for this Diagnostics Profile. Set together with `flow_collector_ref_type_`.
* `flow_collector_ref_type_` (String) - Object type for flow_collector field.
* `monitoring_acl` (String) - Monitoring ACL whose Service VLANs are mirrored to the ERSPAN destination. Set together with `monitoring_acl_ref_type_`.
* `monitoring_acl_ref_type_` (String) - Object type for monitoring_acl field.
* `poll_interval` (Integer) - The sampling rate for sFlow polling (seconds). Set it to `null` to clear it.
* `use_internal_collector` (Boolean) - Use the internal Collector as the flow collector.
* `vrf_type` (String) - Management or Underlay.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `flow_collector` | `flow_collector_ref_type_` | `sflow_collector` |
| `monitoring_acl` | `monitoring_acl_ref_type_` | `monitoring_acl` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_diagnostics_profile.<resource_name> <name>
```

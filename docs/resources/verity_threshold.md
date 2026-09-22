# verity_threshold (Resource)

Manages a Verity Threshold.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_threshold" "example" {
  name = "example"
  critical_escalation_value = ""
  enable = false
  error_escalation_value = ""
  escalation_metric = ""
  escalation_operation = ""
  for = ""
  keep_firing_for = ""
  notice_escalation_value = ""
  operation = ""
  severity = ""
  type = ""
  warning_escalation_value = ""

  rules {
    index = 1
    enable = false
    metric = ""
    operation = ""
    threshold = ""
    threshold_ref_type_ = "threshold"
    type = ""
    value = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `critical_escalation_value` (String) - Value to compare the metric to.
* `enable` (Boolean) - Enable object.
* `error_escalation_value` (String) - Value to compare the metric to.
* `escalation_metric` (String) - Metric threshold is on.
* `escalation_operation` (String) - How to compare the metric to the value.
* `for` (String) - Duration in minutes the threshold must be met before firing the alarm.
* `keep_firing_for` (String) - Duration in minutes to keep firing the alarm after the threshold is no longer met.
* `notice_escalation_value` (String) - Value to compare the metric to.
* `operation` (String) - How to combine rules.
* `rules` (Block List) - Rules for the threshold. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `metric` (String) - Metric threshold is on.
  * `operation` (String) - How to compare the metric to the value.
  * `threshold` (String) - Nested threshold to evaluate (when Type is Threshold). Set together with `threshold_ref_type_`.
  * `threshold_ref_type_` (String) - Object type for threshold field.
  * `type` (String) - Use a metric or a nested threshold.
  * `value` (String) - Value to compare the metric to.
* `severity` (String) - Severity of the alarm when the threshold is met.
* `type` (String) - Type of elements threshold applies to.
* `warning_escalation_value` (String) - Value to compare the metric to.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `rules.threshold` | `threshold_ref_type_` | `threshold` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_threshold.<resource_name> <name>
```

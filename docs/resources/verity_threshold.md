# Threshold Resource

`verity_threshold` manages threshold resources in Verity, which define alarm thresholds based on metrics or nested conditions.

## Example Usage

```hcl
resource "verity_threshold" "example" {
  name = "example"
  enable = true
  type = "interface"
  operation = "and"
  severity = "notice"
  for = "5"
  keep_firing_for = "5"

  rules {
    index = 1
    enable = false
    type = "metric"
    metric = "prometheus"
    operation = "eq"
    value = "1"
    threshold = ""
    threshold_ref_type_ = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `type` (String) - Type of elements threshold applies to
* `operation` (String) - How to combine rules
* `severity` (String) - Severity of the alarm when the threshold is met
* `for` (String) - Duration in minutes the threshold must be met before firing the alarm
* `keep_firing_for` (String) - Duration in minutes to keep firing the alarm after the threshold is no longer met
* `escalation_metric` (String) - Metric threshold is on
* `escalation_operation` (String) - How to compare the metric to the value
* `critical_escalation_value` (String) - Value to compare the metric to
* `error_escalation_value` (String) - Value to compare the metric to
* `warning_escalation_value` (String) - Value to compare the metric to
* `notice_escalation_value` (String) - Value to compare the metric to
* `rules` (Array) - List of rule blocks
  * `enable` (Boolean) - Enable
  * `type` (String) - Use a metric or a nested threshold
  * `metric` (String) - Metric threshold is on
  * `operation` (String) - How to compare the metric to the value
  * `value` (String) - Value to compare the metric to
  * `threshold` (String) - How to compare the metric to the value
  * `threshold_ref_type_` (String) - Object type for threshold field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list

## Import

Threshold resources can be imported using the `name` attribute:

```sh
terraform import verity_threshold.<resource_name> <name>
```

# Threshold Resource

`verity_threshold` manages threshold resources in Verity, which define alarm thresholds based on metrics or nested conditions.

## Example Usage

```hcl
resource "verity_threshold" "example" {
  name            = "example"
  enable          = true
  type            = "interface"
  operation       = "and"
  severity        = "notice"
  for             = "5"
  keep_firing_for = "5"
  rules {
    index               = 1
    enable              = false
    type                = "metric"
    metric              = "prometheus_Switch_interfaces_interface_stats_receive_rx_bytes"
    operation           = "eq"
    value               = "1"
    threshold           = ""
    threshold_ref_type_ = ""
  }
}
```

## Argument Reference

* `name` - Object Name. Must be unique
* `enable` - Enable object
* `type` - Type of elements threshold applies to
* `operation` - How to combine rules
* `severity` - Severity of the alarm when the threshold is met
* `for` - Duration in minutes the threshold must be met before firing the alarm
* `keep_firing_for` - Duration in minutes to keep firing the alarm after the threshold is no longer met
* `rules` - List of rule blocks
  * `enable` - Enable
  * `type` - Use a metric or a nested threshold
  * `metric` - Metric threshold is on
  * `operation` - How to compare the metric to the value
  * `value` - Value to compare the metric to
  * `threshold` - How to compare the metric to the value
  * `threshold_ref_type_` - Object type for threshold field
  * `index` - The index identifying the object. Zero if you want to add an object to the list

## Import

Threshold resources can be imported using the `name` attribute:

```sh
terraform import verity_threshold.<resource_name> <name>
```

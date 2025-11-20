# Threshold Group Resource

`verity_threshold_group` manages threshold group resources in Verity, which define groups of thresholds to apply to network elements.

## Example Usage

```hcl
resource "verity_threshold_group" "example" {
  name = "example"
  enable = true
  type = "interface"

  targets {
    index = 1
    enable = true
    type = "grouping_rules"
    grouping_rules = "teste"
    grouping_rules_ref_type_ = "grouping_rules"
    switchpoint = ""
    switchpoint_ref_type_ = ""
    port = ""
  }

  thresholds {
    index = 1
    enable = true
    severity_override = "critical"
    threshold = "_newtest"
    threshold_ref_type_ = "threshold"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `type` (String) - Type of elements to apply thresholds to. Can be `interface` or `device`
* `targets` (Array) - List of target blocks
  * `enable` (Boolean) - Enable
  * `type` (String) - Specific element or Grouping Rules to apply thresholds to
  * `grouping_rules` (String) - Elements to apply thresholds to
  * `grouping_rules_ref_type_` (String) - Object type for grouping_rules field
  * `switchpoint` (String) - Switchpoint to apply thresholds to
  * `switchpoint_ref_type_` (String) - Object type for switchpoint field
  * `port` (String) - Port to apply thresholds to
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `thresholds` (Array) - List of threshold blocks
  * `enable` (Boolean) - Enable
  * `severity_override` (String) - Override the severity defined in the threshold for this group only. Can be `""`, `warning`, `notice`, `error`, or `critical`
  * `threshold` (String) - Threshold to apply to this group
  * `threshold_ref_type_` (String) - Object type for threshold field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list

## Import

Threshold group resources can be imported using the `name` attribute:

```sh
terraform import verity_threshold_group.<resource_name> <name>
```

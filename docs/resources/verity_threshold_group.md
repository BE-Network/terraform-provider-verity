# verity_threshold_group (Resource)

Manages a Verity Threshold Group.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_threshold_group" "example" {
  name = "example"
  enable = false
  type = ""

  targets {
    index = 1
    element = ""
    element_ref_type_ = ""
    enable = false
    grouping_rules = ""
    grouping_rules_ref_type_ = "grouping_rules"
    port = ""
    sdlc = ""
    type = ""
  }

  thresholds {
    index = 1
    enable = false
    severity_override = ""
    threshold = ""
    threshold_ref_type_ = "threshold"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `targets` (Block List) - Targets to apply thresholds to. Entries are matched by `index`.
  * `element` (String) - Element to apply thresholds to. Set together with `element_ref_type_`.
  * `element_ref_type_` (String) - Object type for element field.
  * `enable` (Boolean) - Enable.
  * `grouping_rules` (String) - Elements to apply thresholds to. Set together with `grouping_rules_ref_type_`.
  * `grouping_rules_ref_type_` (String) - Object type for grouping_rules field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port` (String) - Port to apply thresholds to.
  * `sdlc` (String) - SDLC to apply thresholds to.
  * `type` (String) - Specific element or Grouping Rules to apply thresholds to.
* `thresholds` (Block List) - Thresholds to apply to this group. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `severity_override` (String) - Override the severity defined in the thereshold for this group only.
  * `threshold` (String) - Threshold to apply to this group. Set together with `threshold_ref_type_`.
  * `threshold_ref_type_` (String) - Object type for threshold field.
* `type` (String) - Type of elements to apply thresholds to.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `targets.element` | `element_ref_type_` | `ai_service`, `switchpoint`, `type`, `type` |
| `targets.grouping_rules` | `grouping_rules_ref_type_` | `grouping_rules` |
| `thresholds.threshold` | `threshold_ref_type_` | `threshold` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_threshold_group.<resource_name> <name>
```

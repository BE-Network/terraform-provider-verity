# Grouping Rule Resource

`verity_grouping_rule` manages grouping rule resources in Verity, which define rules for grouping network elements.

## Example Usage

```hcl
resource "verity_grouping_rule" "example" {
  name = "example"
  enable = true
  type = "interface"
  operation = "and"

  rules {
    index = 1
    enable = true
    rule_invert = false
    rule_type = "endpoint_type"
    rule_value = "leaf"
    rule_value_path = ""
    rule_value_path_ref_type_ = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `type` (String) - Type of elements to group.
* `operation` (String) - How to combine rules.
* `rules` (Array) - 
  * `enable` (Boolean) - Enable.
  * `rule_invert` (Boolean) - Invert the rule.
  * `rule_type` (String) - Which type of rule to apply.
  * `rule_value` (String) - Value to compare.
  * `rule_value_path` (String) - Object to compare.
  * `rule_value_path_ref_type_` (String) - Object type for rule_value_path field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Grouping rule resources can be imported using the `name` attribute:

```sh
terraform import verity_grouping_rule.<resource_name> <name>
```

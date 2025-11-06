# Grouping Rule Resource

`verity_grouping_rule` manages grouping rule resources in Verity, which define rules for grouping network elements.

## Example Usage

```hcl
resource "verity_grouping_rule" "example" {
  name      = "example"
  enable    = true
  type      = "interface"
  operation = "and"
  rules {
    index                     = 1
    enable                    = true
    rule_invert               = false
    rule_type                 = "endpoint_type"
    rule_value                = "leaf"
    rule_value_path           = ""
    rule_value_path_ref_type_ = ""
  }
}
```

## Argument Reference

* `name` - Object Name. Must be unique
* `enable` - Enable object
* `type` - Type of elements to group. Can be `interface` or `device`
* `operation` - How to combine rules. Can be `and` or `or`
* `rules` - List of rule blocks
  * `enable` - Enable
  * `rule_invert` - Invert the rule
  * `rule_type` - Which type of rule to apply
  * `rule_value` - Value to compare
  * `rule_value_path` - Object to compare
  * `rule_value_path_ref_type_` - Object type for rule_value_path field
  * `index` - The index identifying the object. Zero if you want to add an object to the list

## Import

Grouping rule resources can be imported using the `name` attribute:

```sh
terraform import verity_grouping_rule.<resource_name> <name>
```

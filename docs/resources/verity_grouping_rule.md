# verity_grouping_rule (Resource)

Manages a Verity Grouping Rule.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_grouping_rule" "example" {
  name = "example"
  enable = false
  operation = ""
  type = ""

  rules {
    index = 1
    enable = false
    rule_invert = false
    rule_type = ""
    rule_value = ""
    rule_value_path = ""
    rule_value_path_ref_type_ = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `operation` (String) - How to combine rules.
* `rules` (Block List) - List of rules within the grouping rule. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `rule_invert` (Boolean) - Invert the rule.
  * `rule_type` (String) - Which type of rule to apply.
  * `rule_value` (String) - Value to compare.
  * `rule_value_path` (String) - Object to compare. Set together with `rule_value_path_ref_type_`.
  * `rule_value_path_ref_type_` (String) - Object type for rule_value_path field.
* `type` (String) - Type of elements to group.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `rules.rule_value_path` | `rule_value_path_ref_type_` | `authenticated_eth_port`, `diagnostics_port_profile`, `eth_port_profile_`, `fabric`, `gateway_profile`, `grouping_rules`, `lag`, `nac_port_profile`, `pod`, `service_port_profile`, `type` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_grouping_rule.<resource_name> <name>
```

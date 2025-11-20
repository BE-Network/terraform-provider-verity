# Community List Resource

`verity_community_list` manages Community List resources in Verity, which define rules for matching BGP community strings and expanded expressions.

## Example Usage

```hcl
resource "verity_community_list" "example" {
  name = "example"
  any_all = "any"
  enable = false
  permit_deny = "permit"
  standard_expanded = "standard"

  lists {
    index = 1
    community_string_expanded_expression = ""
    enable = false
    mode = "community"
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `permit_deny` (String) - Action upon match of Community Strings.
* `any_all` (String) - BGP does not advertise any or all routes that do not match the Community String.
* `standard_expanded` (String) - Used Community String or Expanded Expression.
* `lists` (Array) - 
  * `enable` (Boolean) - Enable of this Community List.
  * `mode` (String) - Mode.
  * `community_string_expanded_expression` (String) - Community String in standard mode and Expanded Expression in Expanded mode.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

Community List resources can be imported using the `name` attribute:

```sh
terraform import verity_community_list.<resource_name> <name>
```

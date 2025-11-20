# Extended Community List Resource

Provides a Verity Extended Community List resource. Extended Community Lists are used to match BGP extended communities for route filtering and policy control.

## Example Usage

```hcl
resource "verity_extended_community_list" "example" {
  name = "example"
  any_all = "any"
  enable = false
  permit_deny = "permit"
  standard_expanded = "standard"

  lists {
    index = 1
    enable = false
    mode = "route"
    route_target_expanded_expression = ""
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `permit_deny` (String) - Action upon match of Community Strings.
* `any_all` (String) - BGP does not advertise any or all routes that do not match the Community String.
* `standard_expanded` (String) - Used Community String or Expanded Expression.
* `lists` (Array) - 
  * `enable` (Boolean) - Enable of this Extended Community List.
  * `mode` (String) - Mode.
  * `route_target_expanded_expression` (String) - Match against a BGP extended community of type Route Target.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

Extended Community Lists can be imported using the name:

```sh
terraform import verity_extended_community_list.<resource_name> <name>
```

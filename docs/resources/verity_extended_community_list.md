# verity_extended_community_list (Resource)

Manages a Verity Extended Community List.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_extended_community_list" "example" {
  name = "example"
  any_all = ""
  enable = false
  permit_deny = ""
  standard_expanded = ""

  lists {
    index = 1
    enable = false
    mode = ""
    route_target_expanded_expression = ""
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `any_all` (String) - BGP does not advertise any or all routes that do not match the Community String.
* `enable` (Boolean) - Enable object.
* `lists` (Block List) - List of Extended Community List entries. Entries are matched by `index`.
  * `enable` (Boolean) - Enable of this Extended Community List.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `mode` (String) - Mode.
  * `route_target_expanded_expression` (String) - Match against a BGP extended community of type Route Target.
* `object_properties` (Block) - Object properties for the Extended Community List. At most one block.
  * `notes` (String) - User Notes.
* `permit_deny` (String) - Action upon match of Community Strings.
* `standard_expanded` (String) - Used Community String or Expanded Expression.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_extended_community_list.<resource_name> <name>
```

# verity_community_list (Resource)

Manages a Verity Community List.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_community_list" "example" {
  name = "example"
  any_all = ""
  enable = false
  permit_deny = ""
  standard_expanded = ""

  lists {
    index = 1
    community_string_expanded_expression = ""
    enable = false
    mode = ""
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
* `lists` (Block List) - List of Community List entries. Entries are matched by `index`.
  * `community_string_expanded_expression` (String) - Community String in standard mode and Expanded Expression in Expanded mode.
  * `enable` (Boolean) - Enable of this Community List.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `mode` (String) - Mode.
* `object_properties` (Block) - Object properties for the Community List. At most one block.
  * `notes` (String) - User Notes.
* `permit_deny` (String) - Action upon match of Community Strings.
* `standard_expanded` (String) - Used Community String or Expanded Expression.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_community_list.<resource_name> <name>
```

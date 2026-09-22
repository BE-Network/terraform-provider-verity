# verity_as_path_access_list (Resource)

Manages a Verity AS Path Access List.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_as_path_access_list" "example" {
  name = "example"
  enable = false
  permit_deny = ""

  lists {
    index = 1
    enable = false
    regular_expression = ""
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

* `enable` (Boolean) - Enable object.
* `lists` (Block List) - List of AS Path Access List entries. Entries are matched by `index`.
  * `enable` (Boolean) - Enable this AS Path Access List.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `regular_expression` (String) - Regular Expression to match BGP Community Strings.
* `object_properties` (Block) - Object properties for the AS Path Access List. At most one block.
  * `notes` (String) - User Notes.
* `permit_deny` (String) - Action upon match of Community Strings.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_as_path_access_list.<resource_name> <name>
```

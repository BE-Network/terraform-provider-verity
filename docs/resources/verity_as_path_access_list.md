# AS Path Access List Resource

`verity_as_path_access_list` manages AS Path Access List resources in Verity, which define rules for matching BGP AS paths using regular expressions.

## Example Usage

```hcl
resource "verity_as_path_access_list" "example" {
  name = "example"
  enable = false
  permit_deny = "permit"

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

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `permit_deny` (String) - Action upon match of Community Strings.
* `lists` (Array) - 
  * `enable` (Boolean) - Enable this AS Path Access List.
  * `regular_expression` (String) - Regular Expression to match BGP Community Strings.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

AS Path Access List resources can be imported using the `name` attribute:

```sh
terraform import verity_as_path_access_list.<resource_name> <name>
```

# AS Path Access List Resource

`verity_as_path_access_list` manages AS Path Access List resources in Verity, which define rules for matching BGP AS paths using regular expressions.

## Example Usage

```hcl
resource "verity_as_path_access_list" "test1" {
  name = "test1"
  depends_on = [verity_operation_stage.as_path_access_list_stage]
  enable = false
  permit_deny = "permit"
  object_properties {
    notes = ""
  }
  lists {
    index = 1
    enable = false
    regular_expression = ""
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the AS Path Access List.
* `enable` - (Optional) Enable this AS Path Access List. Default is `false`.
* `permit_deny` - (Optional) Action upon match of Community Strings. Allowed values: `permit`, `deny`. Default is `permit`.
* `object_properties` - (Optional) Object properties configuration:
  * `notes` - (Optional) User Notes. Default is `""`.
* `lists` - (Optional) List of AS Path Access List entries:
  * `index` - (Optional) The index identifying the object. Zero if you want to add an object to the list.
  * `enable` - (Optional) Enable this AS Path Access List entry. Default is `false`.
  * `regular_expression` - (Optional) Regular Expression to match BGP Community Strings. Default is `""`.

## Import

AS Path Access List resources can be imported using the `name` attribute:

```
$ terraform import verity_as_path_access_list.<resource_name> <name>
```

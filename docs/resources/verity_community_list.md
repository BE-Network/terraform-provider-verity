# Community List Resource

`verity_community_list` manages Community List resources in Verity, which define rules for matching BGP community strings and expanded expressions.

## Example Usage

```hcl
resource "verity_community_list" "test1" {
  name = "test1"
  depends_on = [verity_operation_stage.community_list_stage]
  object_properties {
    notes = ""
  }
  any_all = "any"
  enable = false
  lists {
    index = 1
    community_string_expanded_expression = ""
    enable = false
    mode = "community"
  }
  permit_deny = "permit"
  standard_expanded = "standard"
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the Community List.
* `enable` - (Optional) Enable this Community List. Default is `false`.
* `permit_deny` - (Optional) Action upon match of Community Strings. Allowed values: `permit`, `deny`. Default is `permit`.
* `any_all` - (Optional) BGP does not advertise any or all routes that do not match the Community String. Allowed values: `any`, `all`. Default is `any`.
* `standard_expanded` - (Optional) Used Community String or Expanded Expression. Allowed values: `standard`, `expanded`. Default is `standard`.
* `object_properties` - (Optional) Object properties configuration:
  * `notes` - (Optional) User Notes. Default is `""`.
* `lists` - (Optional) List of Community List entries:
  * `index` - (Optional) The index identifying the object. Zero if you want to add an object to the list.
  * `community_string_expanded_expression` - (Optional) Community String in standard mode and Expanded Expression in Expanded mode. Default is `""`.
  * `enable` - (Optional) Enable this Community List entry. Default is `false`.
  * `mode` - (Optional) Mode. Allowed values: `no_advertise`, `local_as`, `no_peer_set`, `community`, `no_export_set`. Default is `community`.

## Import

Community List resources can be imported using the `name` attribute:

```sh
terraform import verity_community_list.<resource_name> <name>
```

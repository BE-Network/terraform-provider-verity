# AAA Profile Resource

`verity_aaa_profile` manages device AAA profile resources in Verity.

## Example Usage

```hcl
resource "verity_aaa_profile" "example" {
  name = "example"
  enable = true
  fail_through = true
  tacacs_profile = "example-tacacs"
  tacacs_profile_ref_type_ = "tacacs_profile"

  login_default {
    index = 1
    enabled = true
    login_method = "local"
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `fail_through` (Boolean) - When enabled, authentication continues to access each server in the method list if an authentication request fails on one server.
* `tacacs_profile` (String) - TACACS+ profile for authentication.
* `tacacs_profile_ref_type_` (String) - Object type for `tacacs_profile` field.
* `login_default` (Array) -
  * `enabled` (Boolean) - Enable this login method.
  * `login_method` (String) - Authentication method for remote access (SSH, etc.).
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

AAA profile resources can be imported using the `name` attribute:

```sh
terraform import verity_aaa_profile.<resource_name> <name>
```

# verity_aaa_profile (Resource)

Manages a Verity Device AAA Profile.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_aaa_profile" "example" {
  name = "example"
  enable = false
  fail_through = false
  ldap_profile = ""
  ldap_profile_ref_type_ = "ldap_profile"
  tacacs_profile = ""
  tacacs_profile_ref_type_ = "tacacs_profile"

  login_default {
    index = 1
    enabled = false
    login_method = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `fail_through` (Boolean) - When enabled, authentication continues to access each server in the method list if an authentication request fails on one server.
* `ldap_profile` (String) - LDAP profile for authentication. Set together with `ldap_profile_ref_type_`.
* `ldap_profile_ref_type_` (String) - Object type for ldap_profile field.
* `login_default` (Block List) - Authentication method list for remote access. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable this login method.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `login_method` (String) - Authentication method for remote access (SSH, etc.).
* `tacacs_profile` (String) - TACACS+ profile for authentication. Set together with `tacacs_profile_ref_type_`.
* `tacacs_profile_ref_type_` (String) - Object type for tacacs_profile field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `ldap_profile` | `ldap_profile_ref_type_` | `ldap_profile` |
| `tacacs_profile` | `tacacs_profile_ref_type_` | `tacacs_profile` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_aaa_profile.<resource_name> <name>
```

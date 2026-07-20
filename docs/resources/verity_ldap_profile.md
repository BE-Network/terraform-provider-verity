# LDAP Profile Resource

`verity_ldap_profile` manages LDAP authentication and directory lookup profiles in Verity.

## Example Usage

```hcl
resource "verity_ldap_profile" "example" {
  name = "example-ldap"
  enable = true
  base_dn = "dc=example,dc=com"
  bind_dn = "cn=reader,dc=example,dc=com"
  bind_password = "reader-password"
  ldap_version = "3"
  ssl_tls_mode = "start-tls"
  search_scope = "sub"

  ldap_servers {
    index = 1
    enabled = true
    server = "ldap.example.com"
    port = 389
    use_type = "all"
    priority = 1
    ssl_tls_mode = ""
    retransmit_attempts = null
  }

  attribute_maps {
    index = 1
    enabled = true
    map_name = "attribute"
    from = "uid"
    to = "sAMAccountName"
  }
}
```

## Argument Reference

All attributes other than `name` are optional and computed.

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `base_dn`, `bind_dn`, `bind_password`, `encrypted_bind_password` (String) - LDAP bind and search settings.
* `ldap_version`, `ssl_tls_mode`, `search_scope` (String) - LDAP protocol and search behavior.
* `default_port`, `search_time_limit`, `bind_time_limit`, `idle_time_limit`, `retransmit_attempts` (Integer) - Nullable connection and lookup limits.
* `nss_base_passwd`, `nss_base_group`, `nss_base_shadow`, `nss_base_netgroup`, `nss_base_sudoers`, `nss_initgroups_ignore_users`, `nss_skip_members` - NSS lookup settings.
* `pam_filter`, `pam_login_attribute`, `pam_group_dn`, `pam_member_attribute` - PAM lookup and authorization settings.
* `sudoers_base`, `sudoers_search_filter` (String) - Sudo LDAP query settings.
* `ldap_servers` (Array) - LDAP server entries with `enabled`, `server`, `port`, `use_type`, `priority`, `ssl_tls_mode`, `retransmit_attempts`, and `index`.
* `attribute_maps` (Array) - Attribute mapping entries with `enabled`, `map_name`, `from`, `to`, and `index`.

## Import

```sh
terraform import verity_ldap_profile.<resource_name> <name>
```

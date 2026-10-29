# verity_ldap_profile (Resource)

Manages a Verity LDAP Profile.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_ldap_profile" "example" {
  name = "example"
  base_dn = ""
  bind_dn = ""
  bind_password = ""
  bind_time_limit = null
  default_port = null
  enable = false
  encrypted_bind_password = ""
  idle_time_limit = null
  ldap_version = ""
  nss_base_group = ""
  nss_base_netgroup = ""
  nss_base_passwd = ""
  nss_base_shadow = ""
  nss_base_sudoers = ""
  nss_initgroups_ignore_users = ""
  nss_skip_members = false
  pam_filter = ""
  pam_group_dn = ""
  pam_login_attribute = ""
  pam_member_attribute = ""
  retransmit_attempts = null
  search_scope = ""
  search_time_limit = null
  ssl_tls_mode = ""
  sudoers_base = ""
  sudoers_search_filter = ""

  attribute_maps {
    index = 1
    enabled = false
    from = ""
    map_name = ""
    to = ""
  }

  ldap_servers {
    index = 1
    enabled = false
    port = null
    priority = null
    retransmit_attempts = null
    server = ""
    ssl_tls_mode = ""
    use_type = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `attribute_maps` (Block List) - LDAP attribute mapping overrides. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable this mapping entry.
  * `from` (String) - Original RFC2307 attribute or class name to map from.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `map_name` (String) - Category of mapping override.
  * `to` (String) - Replacement attribute/class name or value to map to.
* `base_dn` (String) - Base Distinguished Name to use for LDAP searches.
* `bind_dn` (String) - Distinguished Name with which to bind to the LDAP server. Empty value means anonymous bind.
* `bind_password` (String) - Credentials with which to bind to the LDAP server. Only used together with Bind DN.
* `bind_time_limit` (Integer) - Bind/connect time limit, in seconds. Set it to `null` to clear it.
* `default_port` (Integer) - Default LDAP server port (389 for plain/StartTLS, 636 for LDAPS). Set it to `null` to clear it.
* `enable` (Boolean) - Enable object.
* `encrypted_bind_password` (String) - System-generated encrypted version of Bind Password.
* `idle_time_limit` (Integer) - NSS idle connection time limit, in seconds. Set it to `null` to clear it.
* `ldap_servers` (Block List) - LDAP server entries. Entries are matched by `index`.
  * `enabled` (Boolean) - Enable this LDAP server entry.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port` (Integer) - Server port (overrides global default port). Set it to `null` to clear it.
  * `priority` (Integer) - Server priority (1-99, lower = higher priority). Set it to `null` to clear it.
  * `retransmit_attempts` (Integer) - Per-server retransmit attempts (0-10). Set it to `null` to clear it.
  * `server` (String) - IPv4, IPv6, or DNS hostname for LDAP server.
  * `ssl_tls_mode` (String) - Per-server TLS mode (overrides global setting).
  * `use_type` (String) - Which LDAP client(s) use this server.
* `ldap_version` (String) - LDAP protocol version.
* `nss_base_group` (String) - NSS search base for group map.
* `nss_base_netgroup` (String) - NSS search base for netgroup map.
* `nss_base_passwd` (String) - NSS search base for passwd map.
* `nss_base_shadow` (String) - NSS search base for shadow map.
* `nss_base_sudoers` (String) - NSS search base for sudoers map.
* `nss_initgroups_ignore_users` (String) - Comma-separated list of users for which initgroups() lookups are skipped.
* `nss_skip_members` (Boolean) - If true, the group entry is returned without member attributes.
* `pam_filter` (String) - PAM search filter for retrieving user info.
* `pam_group_dn` (String) - DN of a group a user must belong to for login authorization to succeed.
* `pam_login_attribute` (String) - Attribute used to construct the assertion for the user's login name.
* `pam_member_attribute` (String) - Attribute used to test a user's membership of the PAM group DN.
* `retransmit_attempts` (Integer) - Number of retransmit attempts (0-10). Set it to `null` to clear it.
* `search_scope` (String) - Default LDAP search scope.
* `search_time_limit` (Integer) - Search time limit, in seconds. Set it to `null` to clear it.
* `ssl_tls_mode` (String) - Global TLS mode for LDAP connections.
* `sudoers_base` (String) - Base DN for sudo LDAP queries.
* `sudoers_search_filter` (String) - LDAP filter used to restrict records returned for sudo LDAP queries.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_ldap_profile.<resource_name> <name>
```

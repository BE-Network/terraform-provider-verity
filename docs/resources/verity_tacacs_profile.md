# verity_tacacs_profile (Resource)

Manages a Verity TACACS Profile.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_tacacs_profile" "example" {
  name = "example"
  enable = false

  tacacs_servers {
    index = 1
    auth_type = ""
    enabled = false
    enc_secret = ""
    port = ""
    secret = ""
    server = ""
    timeout = null
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `tacacs_servers` (Block List) - List of TACACS+ servers. Entries are matched by `index`.
  * `auth_type` (String) - TACACS+ authentication type.
  * `enabled` (Boolean) - Enable TACACS+ server.
  * `enc_secret` (String) - TACACS+ shared secret (encrypted).
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port` (String) - TACACS+ server port.
  * `secret` (String) - TACACS+ shared secret.
  * `server` (String) - IPv4, IPv6, or DNS name for TACACS+ server.
  * `timeout` (Integer) - TACACS+ server timeout in seconds. Set it to `null` to clear it.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_tacacs_profile.<resource_name> <name>
```

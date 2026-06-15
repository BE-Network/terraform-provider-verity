# TACACS Profile Resource

`verity_tacacs_profile` manages TACACS profile resources in Verity.

## Example Usage

```hcl
resource "verity_tacacs_profile" "example" {
  name = "example"
  enable = true

  tacacs_servers {
    index = 1
    enabled = true
    server = "10.0.0.10"
    auth_type = "pap"
    port = "49"
    timeout = 5
    secret = "shared-secret"
    enc_secret = ""
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `tacacs_servers` (Array) -
  * `enabled` (Boolean) - Enable TACACS+ server.
  * `server` (String) - IPv4, IPv6, or DNS name for TACACS+ server.
  * `auth_type` (String) - TACACS+ authentication type.
  * `port` (String) - TACACS+ server port.
  * `timeout` (Integer) - TACACS+ server timeout in seconds.
  * `secret` (String) - TACACS+ shared secret.
  * `enc_secret` (String) - TACACS+ shared secret (encrypted).
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

TACACS profile resources can be imported using the `name` attribute:

```sh
terraform import verity_tacacs_profile.<resource_name> <name>
```

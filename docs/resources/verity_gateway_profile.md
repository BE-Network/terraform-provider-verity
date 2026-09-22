# verity_gateway_profile (Resource)

Manages a Gateway Profile.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_gateway_profile" "example" {
  name = "example"
  enable = false

  external_gateways {
    index = 1
    enable = false
    gateway = ""
    gateway_ref_type_ = "gateway"
    peer_gw = false
    source_ip_mask = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `external_gateways` (Block List) - List of external gateway configurations. Entries are matched by `index`.
  * `enable` (Boolean) - Enable row.
  * `gateway` (String) - BGP Gateway referenced for port profile. Set together with `gateway_ref_type_`.
  * `gateway_ref_type_` (String) - Object type for gateway field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `peer_gw` (Boolean) - Setting for paired switches only. Flag indicating that this gateway is a peer gateway. For each gateway profile referencing a BGP session on a member of a leaf pair, the peer should have a gateway profile entry indicating the IP address for the peers gateway.
  * `source_ip_mask` (String) - Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `external_gateways.gateway` | `gateway_ref_type_` | `gateway` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_gateway_profile.<resource_name> <name>
```

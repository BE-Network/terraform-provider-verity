# verity_packet_broker (Resource)

Manages a Verity Packet Broker.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_packet_broker" "example" {
  name = "example"
  enable = false

  ipv4_deny {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = ""
  }

  ipv4_permit {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = ""
  }

  ipv6_deny {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = ""
  }

  ipv6_permit {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `ipv4_deny` (Block List) - IPv4 Deny filters. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter. Set together with `filter_ref_type_`.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv4_permit` (Block List) - IPv4 Permit filters. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter. Set together with `filter_ref_type_`.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_deny` (Block List) - IPv6 Deny filters. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter. Set together with `filter_ref_type_`.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_permit` (Block List) - IPv6 Permit filters. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter. Set together with `filter_ref_type_`.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `ipv4_deny.filter` | `filter_ref_type_` | `ipv4_filter`, `ipv4_list_filter` |
| `ipv4_permit.filter` | `filter_ref_type_` | `ipv4_filter`, `ipv4_list_filter` |
| `ipv6_deny.filter` | `filter_ref_type_` | `ipv6_filter`, `ipv6_list_filter` |
| `ipv6_permit.filter` | `filter_ref_type_` | `ipv6_filter`, `ipv6_list_filter` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_packet_broker.<resource_name> <name>
```

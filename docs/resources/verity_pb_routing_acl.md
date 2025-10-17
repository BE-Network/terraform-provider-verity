# PB Routing ACL Resource

Manages a Policy-Based Routing ACL configuration in Verity. This resource allows you to define named ACL objects with IPv4 and IPv6 permit/deny filter rules that can be referenced by Policy-Based Routing policies for advanced traffic matching and routing decisions.

## Example Usage

```hcl
resource "verity_pb_routing_acl" "example" {
  name         = "pbr_acl_ipv4_1"
  enable       = true
  ipv_protocol = "ipv4"
  next_hop_ips = "192.168.1.1"

  ipv4_permit {
    enable             = true
    filter             = "acl_v4_permit"
    filter_ref_type_   = "ipv4_list"
    index              = 1
  }

  ipv4_deny {
    enable             = true
    filter             = "acl_v4_deny"
    filter_ref_type_   = "ipv4_list"
    index              = 2
  }

  ipv6_permit {
    enable             = true
    filter             = "acl_v6_permit"
    filter_ref_type_   = "ipv6_list"
    index              = 1
  }

  ipv6_deny {
    enable             = true
    filter             = "acl_v6_deny"
    filter_ref_type_   = "ipv6_list"
    index              = 2
  }
}
```

## Argument Reference

- `name` (String, Required): The unique name of the policy-based routing ACL object.
- `enable` (Boolean, Optional): Whether this PB Routing ACL object is enabled.
- `ipv_protocol` (String, Optional): IP protocol version - either "ipv4" or "ipv6".
- `next_hop_ips` (String, Optional): Next hop IP addresses for routing decisions.
- `ipv4_permit` (Block, Optional, Repeatable): Defines an IPv4 permit filter entry. Each block supports:
  - `enable` (Boolean, Optional): Whether this filter entry is enabled.
  - `filter` (String, Optional): The name of the referenced filter list to match.
  - `filter_ref_type_` (String, Optional): The reference type for the filter (e.g., `ipv4_list`).
  - `index` (Number, Optional): The order of this filter entry within the ACL.
- `ipv4_deny` (Block, Optional, Repeatable): Defines an IPv4 deny filter entry. Each block supports:
  - `enable` (Boolean, Optional): Whether this filter entry is enabled.
  - `filter` (String, Optional): The name of the referenced filter list to match.
  - `filter_ref_type_` (String, Optional): The reference type for the filter (e.g., `ipv4_list`).
  - `index` (Number, Optional): The order of this filter entry within the ACL.
- `ipv6_permit` (Block, Optional, Repeatable): Defines an IPv6 permit filter entry. Each block supports:
  - `enable` (Boolean, Optional): Whether this filter entry is enabled.
  - `filter` (String, Optional): The name of the referenced filter list to match.
  - `filter_ref_type_` (String, Optional): The reference type for the filter (e.g., `ipv6_list`).
  - `index` (Number, Optional): The order of this filter entry within the ACL.
- `ipv6_deny` (Block, Optional, Repeatable): Defines an IPv6 deny filter entry. Each block supports:
  - `enable` (Boolean, Optional): Whether this filter entry is enabled.
  - `filter` (String, Optional): The name of the referenced filter list to match.
  - `filter_ref_type_` (String, Optional): The reference type for the filter (e.g., `ipv6_list`).
  - `index` (Number, Optional): The order of this filter entry within the ACL.

## Import

This resource can be imported using the PB Routing ACL name:

```sh
terraform import verity_pb_routing_acl.<resource_name> <name>
```

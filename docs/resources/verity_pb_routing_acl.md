# PB Routing ACL Resource

Manages a Policy-Based Routing ACL configuration in Verity. This resource allows you to define named ACL objects with IPv4 and IPv6 permit/deny filter rules that can be referenced by Policy-Based Routing policies for advanced traffic matching and routing decisions.

## Example Usage

```hcl
resource "verity_pb_routing_acl" "example" {
  name   = "example"
  enable = true

  ipv4_permit {
    enable = true
    filter = "acl_v4_permit"
    filter_ref_type_ = "ipv4_filter"
    index = 1
  }

  ipv4_deny {
    enable = true
    filter = "acl_v4_deny"
    filter_ref_type_ = "ipv4_filter"
    index = 1
  }

  ipv6_permit {
    enable = true
    filter = "acl_v6_permit"
    filter_ref_type_ = "ipv6_filter"
    index = 1
  }

  ipv6_deny {
    enable = true
    filter = "acl_v6_deny"
    filter_ref_type_ = "ipv6_filter"
    index = 1
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `ipv4_permit` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv4_deny` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_permit` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_deny` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

This resource can be imported using the PB Routing ACL name:

```sh
terraform import verity_pb_routing_acl.<resource_name> <name>
```

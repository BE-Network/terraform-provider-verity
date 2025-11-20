# Port ACL Resource

Provides a Verity Port ACL resource. Port ACLs are used to control access to network ports using IPv4 and IPv6 filters for permit and deny rules.

## Example Usage

```hcl
resource "verity_port_acl" "example" {
  name = "example"
  enable = false

  ipv4_deny {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = "ipv4_filter"
  }

  ipv4_permit {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = "ipv4_filter"
  }

  ipv6_deny {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = "ipv6_filter"
  }

  ipv6_permit {
    index = 1
    enable = false
    filter = ""
    filter_ref_type_ = "ipv6_filter"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `ipv4_permit` (Array) - List of IPv4 permit rules
  * `enable` (Boolean) - Enable
  * `filter` (String) - Filter
  * `filter_ref_type_` (String) - Object type for filter field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `ipv4_deny` (Array) - List of IPv4 deny rules
  * `enable` (Boolean) - Enable
  * `filter` (String) - Filter
  * `filter_ref_type_` (String) - Object type for filter field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `ipv6_permit` (Array) - List of IPv6 permit rules
  * `enable` (Boolean) - Enable
  * `filter` (String) - Filter
  * `filter_ref_type_` (String) - Object type for filter field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `ipv6_deny` (Array) - List of IPv6 deny rules
  * `enable` (Boolean) - Enable
  * `filter` (String) - Filter
  * `filter_ref_type_` (String) - Object type for filter field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list

## Import

Port ACLs can be imported using the name:

```sh
terraform import verity_port_acl.<resource_name> <name>
```

# Port ACL Resource

Provides a Verity Port ACL resource. Port ACLs are used to control access to network ports using IPv4 and IPv6 filters for permit and deny rules.

## Example Usage

```hcl
resource "verity_port_acl" "port_acl_test1" {
  name         = "port_acl_test1"
  depends_on   = [verity_operation_stage.port_acl_stage]
  enable       = false
  ipv4_deny {
    index             = 1
    enable            = false
    filter            = ""
    filter_ref_type_  = "ipv4_filter"
  }
  ipv4_permit {
    index             = 1
    enable            = false
    filter            = ""
    filter_ref_type_  = "ipv4_filter"
  }
  ipv6_deny {
    index             = 1
    enable            = false
    filter            = ""
    filter_ref_type_  = "ipv6_filter"
  }
  ipv6_permit {
    index             = 1
    enable            = false
    filter            = ""
    filter_ref_type_  = "ipv6_filter"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `ipv4_permit` (Block, Optional) — List of IPv4 permit rules:
  - `enable` (Boolean, Optional) — Enable. Default: `false`.
  - `filter` (String, Optional) — Filter. Default: `""`.
  - `filter_ref_type_` (String, Optional) — Object type for filter field. Allowed value: `"ipv4_filter"`.
- `ipv4_deny` (Block, Optional) — List of IPv4 deny rules:
  - `enable` (Boolean, Optional) — Enable. Default: `false`.
  - `filter` (String, Optional) — Filter. Default: `""`.
  - `filter_ref_type_` (String, Optional) — Object type for filter field. Allowed value: `"ipv4_filter"`.
- `ipv6_permit` (Block, Optional) — List of IPv6 permit rules:
  - `enable` (Boolean, Optional) — Enable. Default: `false`.
  - `filter` (String, Optional) — Filter. Default: `""`.
  - `filter_ref_type_` (String, Optional) — Object type for filter field. Allowed value: `"ipv6_filter"`.
- `ipv6_deny` (Block, Optional) — List of IPv6 deny rules:
  - `enable` (Boolean, Optional) — Enable. Default: `false`.
  - `filter` (String, Optional) — Filter. Default: `""`.
  - `filter_ref_type_` (String, Optional) — Object type for filter field. Allowed value: `"ipv6_filter"`.

## Import

Port ACLs can be imported using the name:

```hcl
terraform import verity_port_acl.<resource_name> <name>
```

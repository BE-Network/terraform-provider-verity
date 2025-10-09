# PB Routing Resource

Manages a Policy-Based Routing (PBR) configuration in Verity. This resource allows you to define named PBR objects, each with a set of ordered policies referencing ACLs for advanced routing decisions.

## Example Usage

```hcl
resource "verity_pb_routing" "example" {
  name   = "pbr1"
  enable = true

  policy {
    enable                  = true
    pb_routing_acl          = "ipv4_1"
    pb_routing_acl_ref_type = "pb_routing_acl"
    index                   = 1
  }

  policy {
    enable                  = true
    pb_routing_acl          = "pbr_acl_ipv6_1"
    pb_routing_acl_ref_type = "pb_routing_acl"
    index                   = 2
  }
}
```

## Argument Reference

- `name` (String, Required): The unique name of the policy-based routing object.
- `enable` (Boolean, Required): Whether this PBR object is enabled.
- `policy` (Block, Optional, Repeatable): Defines a routing policy entry. Each block supports:
  - `enable` (Boolean, Required): Whether this policy entry is enabled.
  - `pb_routing_acl` (String, Required): The name of the referenced ACL to match for this policy.
  - `pb_routing_acl_ref_type` (String, Required): The reference type for the ACL (typically `pb_routing_acl`).
  - `index` (Number, Required): The order of this policy entry within the PBR object.

## Import

This resource can be imported using the PBR name:

```sh
terraform import verity_pb_routing.<resource_name> <name>
```
# PB Routing Resource

Manages a Policy-Based Routing (PBR) configuration in Verity. This resource allows you to define named PBR objects, each with a set of ordered policies referencing ACLs for advanced routing decisions.

## Example Usage

```hcl
resource "verity_pb_routing" "example" {
  name = "example"
  enable = true

  policy {
    enable = true
    pb_routing_acl = "ipv4_1"
    pb_routing_acl_ref_type = "pb_routing_acl"
    index = 1
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `policy` (Array) - 
  * `enable` (Boolean) - Enable.
  * `pb_routing_acl` (String) - Path to the PB Routing ACL.
  * `pb_routing_acl_ref_type` (String) - Object type for pb_routing_acl field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

This resource can be imported using the PBR name:

```sh
terraform import verity_pb_routing.<resource_name> <name>
```
# verity_pb_routing (Resource)

Manages a Policy-Based Routing resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_pb_routing" "example" {
  name = "example"
  enable = false

  policy {
    index = 1
    enable = false
    pb_routing_acl = ""
    pb_routing_acl_ref_type_ = "pb_routing_acl"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `policy` (Block List) - Policy configurations. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `pb_routing_acl` (String) - Path to the PB Routing ACL. Set together with `pb_routing_acl_ref_type_`.
  * `pb_routing_acl_ref_type_` (String) - Object type for pb_routing_acl field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `policy.pb_routing_acl` | `pb_routing_acl_ref_type_` | `pb_routing_acl` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_pb_routing.<resource_name> <name>
```

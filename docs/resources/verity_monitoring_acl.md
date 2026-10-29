# verity_monitoring_acl (Resource)

Manages a Verity Monitoring ACL.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_monitoring_acl" "example" {
  name = "example"
  enable = false

  services {
    index = 1
    enable = false
    service = ""
    service_ref_type_ = "service"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `services` (Block List) - Services selected for traffic monitoring. Entries are matched by `index`.
  * `enable` (Boolean) - Enable this service entry.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `service` (String) - Service whose VLAN traffic will be selected for monitoring. Set together with `service_ref_type_`.
  * `service_ref_type_` (String) - Object type for service field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `services.service` | `service_ref_type_` | `service` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_monitoring_acl.<resource_name> <name>
```

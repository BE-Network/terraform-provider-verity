# verity_static_ip (Resource)

Manages a Verity Static IP.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_static_ip" "example" {
  name = "example"
  allowed_on = ""
  enable = false
  ip = ""
  lag = ""
  lag_ref_type_ = "lag"
  mac = ""
  port = ""
  service = ""
  service_ref_type_ = "service"
  switch = ""
  switch_ref_type_ = "switchpoint"

  object_properties {
    fabric = ""
    fabric_ref_type_ = "fabric"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object name. Identifies the object in the API and must be unique within its type. Changing it replaces the resource.

### Optional

* `allowed_on` (String) - Allowed On.
* `enable` (Boolean) - Enable object.
* `ip` (String) - IP.
* `lag` (String) - LAG. Set together with `lag_ref_type_`.
* `lag_ref_type_` (String) - Object type for lag field.
* `mac` (String) - MAC.
* `object_properties` (Block) - Object properties for the Static IP. At most one block.
  * `fabric` (String) - Choose a Fabric. Set together with `fabric_ref_type_`.
  * `fabric_ref_type_` (String) - Object type for fabric field.
* `port` (String) - Port.
* `service` (String) - Service. Set together with `service_ref_type_`.
* `service_ref_type_` (String) - Object type for service field.
* `switch` (String) - Switchpoint. Set together with `switch_ref_type_`.
* `switch_ref_type_` (String) - Object type for switch field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `lag` | `lag_ref_type_` | `lag` |
| `object_properties.fabric` | `fabric_ref_type_` | `fabric` |
| `service` | `service_ref_type_` | `service` |
| `switch` | `switch_ref_type_` | `switchpoint` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_static_ip.<resource_name> <name>
```

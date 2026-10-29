# verity_fabric_collection (Resource)

Manages a Verity Fabric Collection.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_fabric_collection" "example" {
  name = "example"
  enable = false
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_fabric_collection.<resource_name> <name>
```

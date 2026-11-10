# verity_fabric_collection (Resource)

Manages a Verity Fabric Collection.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_fabric_collection" "example" {
  name = "example"
}
```

## Argument Reference

### Required

* `name` (String) - Object name. Identifies the object in the API and must be unique within its type. Changing it replaces the resource.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_fabric_collection.<resource_name> <name>
```

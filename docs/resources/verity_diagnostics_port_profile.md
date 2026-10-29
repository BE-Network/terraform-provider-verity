# verity_diagnostics_port_profile (Resource)

Manages a Verity Diagnostics Port Profile.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_diagnostics_port_profile" "example" {
  name = "example"
  enable = false
  enable_sflow = false
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `enable_sflow` (Boolean) - Enable sFlow for this Diagnostics Profile.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_diagnostics_port_profile.<resource_name> <name>
```

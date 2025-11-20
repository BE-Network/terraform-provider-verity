# Diagnostics Port Profile Resource

`verity_diagnostics_port_profile` manages Diagnostics Port Profile resources in Verity, which define diagnostics and sFlow settings for switch ports.

## Example Usage

```hcl
resource "verity_diagnostics_port_profile" "example" {
  name = "example"
  enable = false
  enable_sflow = true
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `enable_sflow` (Boolean) - Enable sFlow for this Diagnostics Profile.

## Import

Diagnostics Port Profile resources can be imported using the `name` attribute:

```sh
terraform import verity_diagnostics_port_profile.<resource_name> <name>
```

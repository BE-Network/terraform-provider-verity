# Diagnostics Port Profile Resource

`verity_diagnostics_port_profile` manages Diagnostics Port Profile resources in Verity, which define diagnostics and sFlow settings for switch ports.

## Example Usage

```hcl
resource "verity_diagnostics_port_profile" "diagnostics_port_profile1" {
  name = "diagnostics_port_profile1"
  depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
  enable = false
  enable_sflow = true
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the Diagnostics Port Profile.
* `enable` - (Optional) Enable this Diagnostics Port Profile. Default is `false`.
* `enable_sflow` - (Optional) Enable sFlow for this Diagnostics Profile. Default is `true`.

## Import

Diagnostics Port Profile resources can be imported using the `name` attribute:

```sh
terraform import verity_diagnostics_port_profile.<resource_name> <name>
```

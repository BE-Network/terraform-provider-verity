# SFB Breakout Resource

Manages an SFP Breakout object in the Verity system.

**Note:** Only PATCH operations are supported for this resource. Create and delete operations are prohibited.

## Example Usage

```hcl
resource "verity_sfp_breakout" "SFP_Breakouts" {
    name        = "SFP Breakouts"
    depends_on  = [verity_operation_stage.sfp_breakout_stage]
    breakout {
        index       = 1
        breakout    = "1x100G"
        enable      = false
        part_number = ""
        vendor      = ""
    }
    enable = true
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `breakout` (Block, Optional, repeatable) — Defines breakout configuration for the SFP port. Each block supports:
    - `breakout` (String, Optional) — Breakout definition; defines number of ports of what speed this port is brokenout to. Default: `"1x100G"`. Allowed values: `8x1G`, `2x10G`, `pg200G`, `1x50G`, `pg100G`, `2x400G`, `8x400G`, `4x100G`, `1x40G`, `1x1G`, `2x25G`, `pg1G`, `8x10G`, `2x40G`, `1x200G`, `2x50G`, `4x40G`, `pg800G`, `pg10G`, `4x50G`, `8x40G`, `8x200G`, `8x100G`, `8x50G`, `pg40G`, `1x25G`, `2x800G`, `4x400G`, `1x10G`, `4x25G`, `4x1G`, `2x100G`, `8x800G`, `8x25G`, `pg400G`, `1x800G`, `2x200G`, `4x200G`, `pg25G`, `1x100G`, `2x1G`, `1x400G`, `4x10G`, `4x800G`, `pg50G`.
    - `enable` (Boolean, Optional) — Enable. Default: `false`.
    - `vendor` (String, Optional) — Vendor. Default: `""`.
    - `part_number` (String, Optional) — Part Number. Default: `""`.
- `object_properties` (Block, Optional) — Additional object properties (reserved for future use).

## Import

This resource can be imported using the object name:

```sh
terraform import verity_sfp_breakout.<resource_name> <name>
```

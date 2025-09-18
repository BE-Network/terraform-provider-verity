# Route Map Resource

Provides a Verity Route Map resource. Route Maps are used to define match and set rules for BGP and other routing protocols.

## Example Usage

```hcl
resource "verity_route_map" "test1" {
  name         = "test1"
  depends_on   = [verity_operation_stage.route_map_stage]
  object_properties {
    notes = ""
  }
  enable       = false
  route_map_clauses {
    index                    = 1
    enable                   = false
    route_map_clause         = ""
    route_map_clause_ref_type_ = "route_map_clause"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `false`.
- `route_map_clauses` (Block, Optional) — List of route map clause entries:
  - `index` (Integer, Optional) — The index identifying the object. Zero if you want to add an object to the list.
  - `enable` (Boolean, Optional) — Enable. Default: `false`.
  - `route_map_clause` (String, Optional) — Route Map Clause is a collection match and set rules. Default: `""`.
  - `route_map_clause_ref_type_` (String, Optional) — Object type for route_map_clause field. Allowed value: `"route_map_clause"`.
- `object_properties` (Block, Optional) —
  - `notes` (String, Optional) — User Notes. Default: `""`.

## Import

Route Maps can be imported using the name:

```hcl
terraform import verity_route_map.<resource_name> <name>
```

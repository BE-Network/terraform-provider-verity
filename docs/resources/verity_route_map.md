# verity_route_map (Resource)

Manages a Verity Route Map.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_route_map" "example" {
  name = "example"
  enable = false

  object_properties {
    notes = ""
  }

  route_map_clauses {
    index = 1
    enable = false
    route_map_clause = ""
    route_map_clause_ref_type_ = "route_map_clause"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `object_properties` (Block) - Object properties for the Route Map. At most one block.
  * `notes` (String) - User Notes.
* `route_map_clauses` (Block List) - List of Route Map Clauses. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `route_map_clause` (String) - Route Map Clause is a collection match and set rules. Set together with `route_map_clause_ref_type_`.
  * `route_map_clause_ref_type_` (String) - Object type for route_map_clause field.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `route_map_clauses.route_map_clause` | `route_map_clause_ref_type_` | `route_map_clause` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_route_map.<resource_name> <name>
```

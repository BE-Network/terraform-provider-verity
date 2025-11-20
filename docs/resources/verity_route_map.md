# Route Map Resource

Provides a Verity Route Map resource. Route Maps are used to define match and set rules for BGP and other routing protocols.

## Example Usage

```hcl
resource "verity_route_map" "example" {
  name = "example"
  enable = false

  route_map_clauses {
    index = 1
    enable = false
    route_map_clause = ""
    route_map_clause_ref_type_ = "route_map_clause"
  }

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `route_map_clauses` (Array) - List of route map clause entries
  * `enable` (Boolean) - Enable
  * `route_map_clause` (String) - Route Map Clause is a collection match and set rules
  * `route_map_clause_ref_type_` (String) - Object type for route_map_clause field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `object_properties` (Object) - Additional object properties
  * `notes` (String) - User Notes

## Import

Route Maps can be imported using the name:

```sh
terraform import verity_route_map.<resource_name> <name>
```

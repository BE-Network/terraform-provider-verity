# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_route_map_clause" "route_map_clause_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
}
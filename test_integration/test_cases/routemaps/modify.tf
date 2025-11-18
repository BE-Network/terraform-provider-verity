# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_route_map" "route_map_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	route_map_clauses {
		index = 2
		enable = false
		route_map_clause = "ipv6_clause"
		route_map_clause_ref_type_ = "route_map_clause"
	}
}
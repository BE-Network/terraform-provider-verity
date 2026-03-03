# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_route_map" "route_map_test_script1" {
    name = "route_map_test_script1"
    depends_on = [verity_operation_stage.route_map_stage]
	object_properties {
		notes = ""
	}
	enable = true
	route_map_clauses {
		index = 1
		enable = true
		route_map_clause = "ipv6_clause"
		route_map_clause_ref_type_ = "route_map_clause"
	}
}
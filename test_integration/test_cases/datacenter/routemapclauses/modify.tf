# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_route_map_clause" "route_map_clause_test_script1" {
	object_properties {
		notes = ""
	}
	enable = false
}

resource "verity_route_map_clause" "route_map_clause_test_script2" {
	object_properties {
		notes = "test123"
	}
	enable = false
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_route_map" "route_map_test_script1" {
	object_properties {
		notes = "test123"
	}
	enable = false
	route_map_clauses {
		index = 1
		enable = false
	}
	route_map_clauses {
		index = 2
	}
}

resource "verity_route_map" "route_map_test_script2" {
	object_properties {
		notes = "test"
	}
	enable = true
	route_map_clauses {
		index = 1
		enable = false
	}
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_extended_community_list" "extended_community_list_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 2
		enable = false
		mode = "route"
		route_target_expanded_expression = ""
	}
}
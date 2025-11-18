# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_extended_community_list" "extended_community_list_test_script1" {
    name = "extended_community_list_test_script1"
    depends_on = [verity_operation_stage.extended_community_list_stage]
	object_properties {
		notes = ""
	}
	any_all = "any"
	enable = true
	lists {
		index = 1
		enable = false
		mode = "route"
		route_target_expanded_expression = ""
	}
	permit_deny = "permit"
	standard_expanded = "standard"
}
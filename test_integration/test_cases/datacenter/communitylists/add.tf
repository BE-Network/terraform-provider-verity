# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_community_list" "community_list_test_script1" {
    name = "community_list_test_script1"
    depends_on = [verity_operation_stage.community_list_stage]
	object_properties {
		notes = ""
	}
	any_all = "any"
	enable = true
	lists {
		index = 1
		community_string_expanded_expression = ""
		enable = false
		mode = "community"
	}
	lists {
		index = 2
		community_string_expanded_expression = ""
		enable = false
		mode = "community"
	}
	permit_deny = "permit"
	standard_expanded = "standard"
}
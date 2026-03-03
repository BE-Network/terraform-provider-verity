# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_as_path_access_list" "as_path_access_list_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 4
		enable = false
		regular_expression = ""
	}
}

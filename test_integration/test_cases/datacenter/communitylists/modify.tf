# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_community_list" "community_list_test_script1" {
	object_properties {
		notes = "test1"
	}
	enable = false
	lists {
		index = 3
		community_string_expanded_expression = ""
		enable = false
		mode = "community"
	}
	lists {
		index = 4
		community_string_expanded_expression = ""
		enable = true
		mode = "community"
	}
}

resource "verity_community_list" "community_list_test_script2" {
	object_properties {
		notes = "test123"
	}
	enable = true
	lists {
		index = 2
	}
}
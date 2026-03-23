# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_as_path_access_list" "as_path_access_list_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 2
		enable = true
	}
	lists {
		index = 3
		enable = false
		regular_expression = ""
	}
}

resource "verity_as_path_access_list" "as_path_access_list_test_script2" {
	object_properties {
		notes = "test123"
	}
	enable = true
	lists {
		index = 1
		enable = false
	}
	lists {
		index = 3
	}
	lists {
		index = 4
	}
}

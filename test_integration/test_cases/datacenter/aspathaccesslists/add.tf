# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_as_path_access_list" "as_path_access_list_test_script1" {
    name = "as_path_access_list_test_script1"
    depends_on = [verity_operation_stage.as_path_access_list_stage]
	object_properties {
		notes = ""
	}
	enable = true
	lists {
		index = 1
		enable = false
		regular_expression = ""
	}
	lists {
		index = 2
		enable = false
		regular_expression = ""
	}
	permit_deny = "permit"
}

resource "verity_as_path_access_list" "as_path_access_list_test_script2" {
    name = "as_path_access_list_test_script2"
    depends_on = [verity_operation_stage.as_path_access_list_stage]
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 1
		enable = true
		regular_expression = ""
	}
	lists {
		index = 2
		enable = true
		regular_expression = ""
	}
	lists {
		index = 3
		enable = true
		regular_expression = ""
	}
	lists {
		index = 4
		enable = true
		regular_expression = ""
	}
	permit_deny = "permit"
}


resource "verity_as_path_access_list" "as_path_access_list_test1" {
    name = "as_path_access_list_test1"
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


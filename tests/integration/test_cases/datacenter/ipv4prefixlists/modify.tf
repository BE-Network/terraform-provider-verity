# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_ipv4_prefix_list" "ipv4_prefix_list_test_script1" {
	object_properties {
		notes = ""
	}
	enable = true
	lists {
		index = 1
		enable = true
		greater_than_equal_value = 4
		ipv4_prefix = ""
		less_than_equal_value = 10
	}
	lists {
		index = 2
	}
}

resource "verity_ipv4_prefix_list" "ipv4_prefix_list_test_script2" {
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 2
		enable = true
		greater_than_equal_value = 0
		ipv4_prefix = ""
		less_than_equal_value = 0
		permit_deny = "permit"
	}
}

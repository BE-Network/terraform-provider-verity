# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_ipv4_prefix_list" "ipv4_prefix_list_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	lists {
		index = 2
		enable = false
		greater_than_equal_value = null
		ipv4_prefix = ""
		less_than_equal_value = null
		permit_deny = "permit"
	}
}
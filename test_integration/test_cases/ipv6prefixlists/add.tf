# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_ipv6_prefix_list" "ipv6_prefix_list_test_script1" {
    name = "ipv6_prefix_list_test_script1"
    depends_on = [verity_operation_stage.ipv6_prefix_list_stage]
	object_properties {
		notes = ""
	}
	enable = true
	lists {
		index = 1
		enable = false
		greater_than_equal_value = null
		ipv6_prefix = "0::0/0"
		less_than_equal_value = null
		permit_deny = "permit"
	}
}
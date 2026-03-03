# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_grouping_rule" "grouping_rule_test_script1" {
	enable = false
	rules {
		index = 2
		enable = true
		rule_invert = false
		rule_type = "endpoint_type"
		rule_value = "leaf"
		rule_value_path = ""
		rule_value_path_ref_type_ = ""
	}
}
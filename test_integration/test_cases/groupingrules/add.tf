# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_grouping_rule" "grouping_rule_test_script1" {
    name = "grouping_rule_test_script1"
    depends_on = [verity_operation_stage.grouping_rule_stage]
	enable = true
	operation = "and"
	rules {
		index = 1
		enable = true
		rule_invert = false
		rule_type = "endpoint_type"
		rule_value = "leaf"
		rule_value_path = ""
		rule_value_path_ref_type_ = ""
	}
	type = "interface"
}
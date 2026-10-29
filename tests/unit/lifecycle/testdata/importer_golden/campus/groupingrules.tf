
resource "verity_grouping_rule" "grouping_rules_test1" {
    name = "grouping_rules_test1"
    depends_on = [verity_operation_stage.grouping_rule_stage]
	enable = false
	operation = "and"
	rules {
		index = 1
		enable = false
		rule_invert = false
		rule_type = ""
		rule_value = ""
		rule_value_path = ""
		rule_value_path_ref_type_ = ""
	}
	type = "device"
}


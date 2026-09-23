
resource "verity_grouping_rule" "test" {
    name = "test"
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


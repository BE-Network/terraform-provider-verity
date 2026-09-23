
resource "verity_threshold_group" "threshold_group_test1" {
    name = "threshold_group_test1"
    depends_on = [verity_operation_stage.threshold_group_stage]
	enable = false
	targets {
		index = 1
		enable = false
		grouping_rules = ""
		grouping_rules_ref_type_ = ""
		port = ""
		type = "grouping_rules"
	}
	thresholds {
		index = 1
		enable = false
		severity_override = ""
		threshold = ""
		threshold_ref_type_ = ""
	}
	type = "device"
}


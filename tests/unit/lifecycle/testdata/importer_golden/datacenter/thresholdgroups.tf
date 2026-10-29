
resource "verity_threshold_group" "Test_Alarm" {
    name = "Test Alarm"
    depends_on = [verity_operation_stage.threshold_group_stage]
	enable = true
	targets {
		index = 1
		enable = true
		grouping_rules = ""
		grouping_rules_ref_type_ = ""
		port = "Eth/0.5"
		type = "element"
	}
	thresholds {
		index = 1
		enable = true
		severity_override = "critical"
		threshold = "_newtest"
		threshold_ref_type_ = "threshold"
	}
	type = "interface"
}


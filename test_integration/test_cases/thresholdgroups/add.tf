# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_threshold_group" "threshold_group_test_script1" {
    name = "threshold_group_test_script1"
    depends_on = [verity_operation_stage.threshold_group_stage]
	enable = true
	targets {
		index = 1
		enable = true
		grouping_rules = ""
		grouping_rules_ref_type_ = ""
		port = "Eth/0.5"
		switchpoint = "t6"
		switchpoint_ref_type_ = "switchpoint"
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
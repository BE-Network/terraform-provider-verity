# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_threshold" "threshold_test_script1" {
    name = "threshold_test_script1"
    depends_on = [verity_operation_stage.threshold_stage]
	critical_escalation_value = ""
	enable = true
	error_escalation_value = ""
	escalation_metric = ""
	escalation_operation = "eq"
	for = "5"
	keep_firing_for = "5"
	notice_escalation_value = ""
	operation = "and"
	rules {
		index = 1
		enable = true
		metric = ""
		operation = "eq"
		threshold = ""
		threshold_ref_type_ = ""
		type = "metric"
		value = "1"
	}
	severity = "notice"
	type = "interface"
	warning_escalation_value = ""
}


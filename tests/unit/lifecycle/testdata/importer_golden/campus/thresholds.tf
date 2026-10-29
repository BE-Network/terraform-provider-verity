
resource "verity_threshold" "threshold_name" {
    name = "threshold_name"
    depends_on = [verity_operation_stage.threshold_stage]
	critical_escalation_value = ""
	enable = false
	error_escalation_value = ""
	escalation_metric = ""
	escalation_operation = "eq"
	for = "5"
	keep_firing_for = "5"
	notice_escalation_value = ""
	operation = "and"
	rules {
		index = 1
		enable = false
		metric = ""
		operation = "=="
		threshold = ""
		threshold_ref_type_ = ""
		type = "metric"
		value = ""
	}
	severity = "notice"
	type = "device"
	warning_escalation_value = ""
}


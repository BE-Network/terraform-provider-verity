# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_threshold_group" "threshold_group_test_script1" {
	enable = false
	thresholds {
		index = 2
		enable = true
		severity_override = "critical"
		threshold = "_newtest"
		threshold_ref_type_ = "threshold"
	}
}
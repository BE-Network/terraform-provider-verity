# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_threshold" "threshold_test_script1" {
	enable = false
	rules {
		index = 2
		enable = true
		metric = ""
		operation = "eq"
		threshold = ""
		threshold_ref_type_ = ""
		type = "metric"
		value = "1"
	}
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_lag" "lag_test_script1" {
	enable = false
	fallback = true
}
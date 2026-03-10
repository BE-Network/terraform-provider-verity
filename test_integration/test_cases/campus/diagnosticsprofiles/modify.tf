# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_diagnostics_profile" "diagnostics_profile_test_script1" {
	enable = false
	enable_sflow = false
	poll_interval = 25
}

resource "verity_diagnostics_profile" "diagnostics_profile_test_script2" {
	enable = true
	enable_sflow = true
	poll_interval = 32
}
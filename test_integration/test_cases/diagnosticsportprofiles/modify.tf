# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test_script1" {
	enable = false
	enable_sflow = false
}
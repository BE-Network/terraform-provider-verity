# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test_script1" {
    name = "diagnostics_port_profile_test_script1"
    depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
	enable = true
	enable_sflow = true
}

resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test_script2" {
    name = "diagnostics_port_profile_test_script2"
    depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
	enable = false
	enable_sflow = false
}
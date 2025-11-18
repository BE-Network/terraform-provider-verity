# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_diagnostics_profile" "diagnostics_profile_test_script1" {
    name = "diagnostics_profile_test_script1"
    depends_on = [verity_operation_stage.diagnostics_profile_stage]
	enable = true
	enable_sflow = true
	flow_collector = "Example Collector"
	flow_collector_ref_type_ = "sflow_collector"
	poll_interval = 20
	vrf_type = "management"
}
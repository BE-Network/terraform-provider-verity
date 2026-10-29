
resource "verity_diagnostics_profile" "diagnostics_profile_test1" {
    name = "diagnostics_profile_test1"
    depends_on = [verity_operation_stage.diagnostics_profile_stage]
	enable = false
	enable_sflow = false
	flow_collector = ""
	flow_collector_ref_type_ = ""
	poll_interval = 20
	vrf_type = "management"
}


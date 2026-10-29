
resource "verity_diagnostics_profile" "Example" {
    name = "Example"
    depends_on = [verity_operation_stage.diagnostics_profile_stage]
	enable = true
	enable_sflow = true
	flow_collector = "Example Collector"
	flow_collector_ref_type_ = "sflow_collector"
	poll_interval = 20
	vrf_type = "management"
}


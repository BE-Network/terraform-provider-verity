
resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test1" {
    name = "diagnostics_port_profile_test1"
    depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
	enable = false
	enable_sflow = true
}


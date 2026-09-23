
resource "verity_diagnostics_port_profile" "Example_Diagnostic_Port_Profile" {
    name = "Example Diagnostic Port Profile"
    depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
	enable = true
	enable_sflow = true
}


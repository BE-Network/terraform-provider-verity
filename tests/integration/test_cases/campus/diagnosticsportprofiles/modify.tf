
resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test_script1" {
	enable = false
	enable_sflow = false
}

resource "verity_diagnostics_port_profile" "diagnostics_port_profile_test_script2" {
	enable = true
	enable_sflow = true
}
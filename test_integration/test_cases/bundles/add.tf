# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_bundle" "bundle_test_script1" {
    name = "bundle_test_script1"
    depends_on = [verity_operation_stage.bundle_stage]
	object_properties {
		group = ""
		is_for_switch = true
		is_public = false
	}
	cli_commands = ""
	device_settings = "(Device Settings)"
	device_settings_ref_type_ = "eth_device_profiles"
	device_voice_settings = "(SIP Voice Device)"
	device_voice_settings_ref_type_ = "device_voice_settings"
	diagnostics_profile = ""
	diagnostics_profile_ref_type_ = ""
	enable = true
	eth_port_paths {
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "swp1"
		eth_port_num_eth_port_profile_ref_type_ = "service_port_profile"
		eth_port_num_eth_port_profile = "SD Router"
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 1
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
	}
	protocol = "SIP"
	user_services {
		row_ip_mask = ""
		index = 1
		row_app_connected_service_ref_type_ = ""
		row_app_enable = false
		row_app_connected_service = ""
		row_app_cli_commands = ""
	}
	voice_port_profile_paths {
		voice_port_num_voice_port_profiles = ""
		voice_port_num_voice_port_profiles_ref_type_ = ""
		index = 1
	}
}
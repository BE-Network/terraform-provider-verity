# Test case: Modify existing resources
# Define modified versions of the resources from add.tf


resource "verity_bundle" "bundle_test_script1" {
	object_properties {
		group = "test"
		is_for_switch = false
	}
	enable = false
	eth_port_paths {
		index = 2
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "swp2"
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
	}
	user_services {
		row_ip_mask = ""
		index = 2
		row_app_connected_service_ref_type_ = ""
		row_app_enable = false
		row_app_connected_service = ""
		row_app_cli_commands = ""
	}
	voice_port_profile_paths {
		voice_port_num_voice_port_profiles = ""
		voice_port_num_voice_port_profiles_ref_type_ = ""
		index = 2
	}
}
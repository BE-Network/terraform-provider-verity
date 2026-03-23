# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_bundle" "bundle_test_script1" {
    name = "bundle_test_script1"
    depends_on = [verity_operation_stage.bundle_stage]
	object_properties {
		group = ""
		is_for_switch = false
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
		eth_port_num_eth_port_profile = "data_1"
		eth_port_num_eth_port_settings = "test"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 1
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "vwan0"
		eth_port_num_eth_port_profile_ref_type_ = "eth_port_profile_"
	}
	eth_port_paths {
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 2
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "Eth1"
		eth_port_num_eth_port_profile_ref_type_ = "eth_port_profile_"
		eth_port_num_eth_port_profile = "data_2"
		eth_port_num_eth_port_settings = "(Port Settings)"
	}
	eth_port_paths {
		port_name = "Eth2"
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 3
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
	}
	eth_port_paths {
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 4
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "Eth3"
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
	}
	protocol = "SIP"
	user_services {
		row_app_connected_service_ref_type_ = ""
		row_app_enable = false
		row_app_connected_service = ""
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 1
	}
	voice_port_profile_paths {
		voice_port_num_voice_port_profiles = ""
		voice_port_num_voice_port_profiles_ref_type_ = ""
		index = 1
	}
}


resource "verity_bundle" "bundle_test_script2" {
    name = "bundle_test_script2"
    depends_on = [verity_operation_stage.bundle_stage]
	object_properties {
		group = ""
		is_for_switch = false
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
		eth_port_num_eth_port_profile = "data_1"
		eth_port_num_eth_port_settings = "test"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 1
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "vwan0"
		eth_port_num_eth_port_profile_ref_type_ = "eth_port_profile_"
	}
	eth_port_paths {
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 2
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "Eth1"
		eth_port_num_eth_port_profile_ref_type_ = "eth_port_profile_"
		eth_port_num_eth_port_profile = "data_2"
		eth_port_num_eth_port_settings = "(Port Settings)"
	}
	eth_port_paths {
		port_name = "Eth2"
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 3
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
	}
	eth_port_paths {
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 4
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		port_name = "Eth3"
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
	}
	protocol = "SIP"
	user_services {
		row_app_connected_service_ref_type_ = ""
		row_app_enable = false
		row_app_connected_service = ""
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 1
	}
	voice_port_profile_paths {
		voice_port_num_voice_port_profiles = ""
		voice_port_num_voice_port_profiles_ref_type_ = ""
		index = 1
	}
}
# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_bundle" "bundle_test_script1" {
    name = "bundle_test_script1"
    depends_on = [verity_operation_stage.bundle_stage]
	object_properties {
		group = ""
		is_for_switch = true
	}
	cli_commands = ""
	device_settings = "ChangesDeviceSetting"
	device_settings_ref_type_ = "eth_device_profiles"
	diagnostics_profile = "Example"
	diagnostics_profile_ref_type_ = "diagnostics_profile"
	enable = true
	eth_port_paths {
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_gateway_profile = ""
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 1
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = ""
	}
	eth_port_paths {
		eth_port_num_gateway_profile = ""
		index = 2
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		eth_port_num_eth_port_profile_ref_type_ = ""
	}
	protocol = "SIP"
	user_services {
		row_app_enable = true
		row_app_connected_service = "Apple"
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 1
		row_app_connected_service_ref_type_ = "service"
	}
}


resource "verity_bundle" "bundle_test_script2" {
    name = "bundle_test_script2"
    depends_on = [verity_operation_stage.bundle_stage]
	object_properties {
		group = "test"
		is_for_switch = false
	}
	cli_commands = ""
	device_settings = "ChangesDeviceSetting"
	device_settings_ref_type_ = "eth_device_profiles"
	diagnostics_profile = "Example"
	diagnostics_profile_ref_type_ = "diagnostics_profile"
	enable = false
	eth_port_paths {
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_gateway_profile = ""
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 1
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = "1/1"
	}
	eth_port_paths {
		eth_port_num_gateway_profile = ""
		index = 2
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = "1/2"
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		eth_port_num_eth_port_profile_ref_type_ = ""
	}
	eth_port_paths {
		eth_port_num_gateway_profile = ""
		index = 3
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = "1/3"
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		eth_port_num_eth_port_profile_ref_type_ = ""
	}
	protocol = "SIP"
	user_services {
		row_app_enable = true
		row_app_connected_service = "Apple"
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 1
		row_app_connected_service_ref_type_ = "service"
	}
	user_services {
		row_app_enable = true
		row_app_connected_service = "Apple"
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 2
		row_app_connected_service_ref_type_ = "service"
	}
}
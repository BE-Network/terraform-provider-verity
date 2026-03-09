# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_bundle" "bundle_test_script1" {
	object_properties {
		group = "test1"
	}
	enable = false
	eth_port_paths {
		eth_port_num_eth_port_profile_ref_type_ = ""
		eth_port_num_gateway_profile = ""
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		index = 3
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = "1/3"
	}
	eth_port_paths {
		eth_port_num_gateway_profile = ""
		index = 4
		diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ = ""
		eth_port_num_eth_port_settings_ref_type_ = "eth_port_settings"
		eth_port_num_gateway_profile_ref_type_ = ""
		port_name = "1/4"
		eth_port_num_eth_port_profile = ""
		eth_port_num_eth_port_settings = "(Port Settings)"
		diagnostics_port_profile_num_diagnostics_port_profile = ""
		eth_port_num_eth_port_profile_ref_type_ = ""
	}
	user_services {
		row_app_enable = false
		row_app_connected_service = "Apple"
		row_app_cli_commands = ""
		row_ip_mask = ""
		index = 2
		row_app_connected_service_ref_type_ = "service"
	}
}


resource "verity_bundle" "bundle_test_script2" {
	object_properties {
		is_for_switch = true
	}
	enable = true
	eth_port_paths {
		index = 3
	}
	user_services {
		row_app_enable = false
		index = 1
	}
	user_services {
		index = 2
	}
}
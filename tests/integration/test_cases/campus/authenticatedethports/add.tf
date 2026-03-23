# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script1" {
    name = "authenticated_eth_port_test_script1"
    depends_on = [verity_operation_stage.authenticated_eth_port_stage]
	object_properties {
		group = "test"
		port_monitoring = "test"
	}
	allow_mac_based_authentication = false
	connection_mode = "PortMode"
	enable = true
	eth_ports {
		index = 1
		eth_port_profile_num_enable = true
		eth_port_profile_num_eth_port = "data_2"
		eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
		eth_port_profile_num_radius_filter_id = "data_2"
		eth_port_profile_num_walled_garden_set = false
	}
	eth_ports {
		index = 2
		eth_port_profile_num_enable = true
		eth_port_profile_num_eth_port = ""
		eth_port_profile_num_eth_port_ref_type_ = ""
		eth_port_profile_num_radius_filter_id = ""
		eth_port_profile_num_walled_garden_set = false
	}
	mac_authentication_holdoff_sec = 60
	reauthorization_period_sec = 3600
	trusted_port = false
}

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script2" {
    name = "authenticated_eth_port_test_script2"
    depends_on = [verity_operation_stage.authenticated_eth_port_stage]
	object_properties {
		group = ""
		port_monitoring = ""
	}
	allow_mac_based_authentication = false
	connection_mode = "PortMode"
	enable = true
	eth_ports {
		index = 1
		eth_port_profile_num_enable = true
		eth_port_profile_num_eth_port = "data_2"
		eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
		eth_port_profile_num_radius_filter_id = "data_2"
		eth_port_profile_num_walled_garden_set = false
	}
	eth_ports {
		index = 2
		eth_port_profile_num_enable = true
		eth_port_profile_num_eth_port = ""
		eth_port_profile_num_eth_port_ref_type_ = ""
		eth_port_profile_num_radius_filter_id = ""
		eth_port_profile_num_walled_garden_set = false
	}
	mac_authentication_holdoff_sec = 60
	reauthorization_period_sec = 3600
	trusted_port = false
}

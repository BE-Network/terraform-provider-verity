# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script1" {
    name = "authenticated_eth_port_test_script1"
    depends_on = [verity_operation_stage.authenticated_eth_port_stage]
	object_properties {
		group = ""
		port_monitoring = ""
	}
	allow_mac_based_authentication = false
	connection_mode = "PortMode"
	enable = false
	eth_ports {
		index = 1
		eth_port_profile_num_enable = false
		eth_port_profile_num_eth_port = ""
		eth_port_profile_num_eth_port_ref_type_ = ""
		eth_port_profile_num_radius_filter_id = ""
		eth_port_profile_num_walled_garden_set = false
	}
	mac_authentication_holdoff_sec = 60
	reauthorization_period_sec = 3600
	trusted_port = false
}

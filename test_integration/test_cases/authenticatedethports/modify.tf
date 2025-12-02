# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script1" {
	object_properties {
		group = "test"
	}
	allow_mac_based_authentication = true
	enable = true
	eth_ports {
		index = 2
		eth_port_profile_num_enable = false
		eth_port_profile_num_eth_port = ""
		eth_port_profile_num_eth_port_ref_type_ = ""
		eth_port_profile_num_radius_filter_id = ""
		eth_port_profile_num_walled_garden_set = false
	}
	trusted_port = true
}
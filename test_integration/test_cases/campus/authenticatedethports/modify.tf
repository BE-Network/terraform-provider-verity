# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script1" {
	object_properties {
		group = ""
		port_monitoring = ""
	}
	allow_mac_based_authentication = true
	enable = false
	eth_ports {
		index = 1
		eth_port_profile_num_enable = false
		eth_port_profile_num_walled_garden_set = true
	}
	eth_ports {
		index = 2
	}
	mac_authentication_holdoff_sec = 65
	reauthorization_period_sec = 3700
	trusted_port = true
}

resource "verity_authenticated_eth_port" "authenticated_eth_port_test_script2" {
	object_properties {
		group = "test"
		port_monitoring = "test1"
	}
	allow_mac_based_authentication = true
	enable = false
	eth_ports {
		index = 1
		eth_port_profile_num_enable = false
		eth_port_profile_num_walled_garden_set = true
	}
	eth_ports {
		index = 2
	}
	mac_authentication_holdoff_sec = 62
	reauthorization_period_sec = 3601
	trusted_port = true
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
	object_properties {
		group = "test"
		icon = "test"
		label = "test"
		port_monitoring = "test"
		sort_by_name = true
	}
	enable = true
	services {
		index = 1
		row_num_enable = true
		row_num_external_vlan = 136
	}
	services {
		index = 2
		row_num_enable = true
		row_num_external_vlan = 0
	}
	tls = true
	trusted_port = true
}

resource "verity_eth_port_profile" "eth_port_profile_test_script2" {
	enable = true
	services {
		index = 1
		row_num_enable = true
		row_num_external_vlan = 137
	}
	services {
		index = 4
		row_num_enable = true
		row_num_external_vlan = 0
	}
	tls = true
	trusted_port = true
}
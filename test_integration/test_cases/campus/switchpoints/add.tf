# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_switchpoint" "switchpoint_test_script1" {
    name = "switchpoint_test_script1"
    depends_on = [verity_operation_stage.switchpoint_stage]
	object_properties {
		aggregate = false
		draw_as_edge_device = false
		expected_parent_endpoint = ""
		is_host = false
		number_of_multipoints = 0
		user_notes = ""
	}
	connected_bundle = "fourth"
	connected_bundle_ref_type_ = "endpoint_bundle"
	device_serial_number = ""
	enable = true
	eths {
		index = 1
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "vwan0"
	}
	eths {
		index = 2
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth1"
	}
	eths {
		index = 3
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth2"
	}
	eths {
		index = 4
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth3"
	}
	eths {
		index = 5
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth4"
	}
	eths {
		index = 6
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth5"
	}
	eths {
		index = 7
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth6"
	}
	eths {
		index = 8
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth7"
	}
	eths {
		index = 9
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth8"
	}
	is_fabric = false
	locked = false
	out_of_band_management = false
	read_only_mode = false
}

resource "verity_switchpoint" "switchpoint_test_script2" {
    name = "switchpoint_test_script2"
    depends_on = [verity_operation_stage.switchpoint_stage]
	object_properties {
		aggregate = false
		draw_as_edge_device = false
		expected_parent_endpoint = ""
		is_host = false
		number_of_multipoints = 0
		user_notes = ""
	}
	connected_bundle = "fourth"
	connected_bundle_ref_type_ = "endpoint_bundle"
	device_serial_number = ""
	enable = true
	eths {
		index = 1
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "vwan0"
	}
	eths {
		index = 2
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth1"
	}
	eths {
		index = 3
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth2"
	}
	eths {
		index = 4
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth3"
	}
	eths {
		index = 5
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth4"
	}
	eths {
		index = 6
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth5"
	}
	eths {
		index = 7
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth6"
	}
	eths {
		index = 8
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth7"
	}
	eths {
		index = 9
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = "Eth8"
	}
	is_fabric = false
	locked = false
	out_of_band_management = false
	read_only_mode = false
}
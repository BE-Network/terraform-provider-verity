# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
	object_properties {
		group = "test"
	}
	aging_time = 1
	auto_negotiation = false
	bpdu_filter = true
	bpdu_guard = true
	broadcast = false
	bsp_enable = true
	detect_bridging_loops = true
	enable = false
	enable_speed_control = false
	fast_learning_mode = false
	guard_loop = true
	lldp_enable = false
	lldp_med {
		index = 1
		lldp_med_row_num_dscp_mark = 2
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 2
	}
	lldp_med {
		index = 2
		lldp_med_row_num_dscp_mark = 3
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 3
	}
	lldp_med_enable = false 
	mac_limit = 102
	max_allowed_value = 1002
	multicast = false
	poe_enable = true
	stp_enable = true
	unidirectional_link_detection = true
}

resource "verity_eth_port_settings" "eth_port_settings_test_script2" {
	object_properties {
		group = "test1"
	}
	aging_time = 2
	auto_negotiation = true
	bpdu_filter = true
	bpdu_guard = true
	broadcast = false
	bsp_enable = true
	detect_bridging_loops = true
	enable = false
	enable_speed_control = false
	fast_learning_mode = false
	guard_loop = true
	lldp_enable = false
	lldp_med {
		index = 1
		lldp_med_row_num_dscp_mark = 3
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 3
	}
	lldp_med {
		index = 2
		lldp_med_row_num_dscp_mark = 4
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 4
	}
	lldp_med_enable = false
	mac_limit = 103
	max_allowed_value = 1003
	multicast = false
	poe_enable = true
	stp_enable = true
	unidirectional_link_detection = true
}
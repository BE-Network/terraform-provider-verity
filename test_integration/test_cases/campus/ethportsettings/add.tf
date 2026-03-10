# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
    name = "eth_port_settings_test_script1"
    depends_on = [verity_operation_stage.eth_port_settings_stage]
	object_properties {
		group = ""
	}
	action = "Protect"
	aging_time = 0
	aging_type = "absolute"
	allocated_power = "0.0"
	auto_negotiation = true
	bpdu_filter = false
	bpdu_guard = false
	broadcast = true
	bsp_enable = false
	cli_commands = ""
	detect_bridging_loops = false
	duplex_mode = "Half"
	enable = true
	enable_speed_control = true
	fast_learning_mode = true
	fec = "unaltered"
	guard_loop = false
	lldp_enable = true
	lldp_med {
		index = 1
		lldp_med_row_num_advertised_applicatio = ""
		lldp_med_row_num_dscp_mark = 0
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 0
		lldp_med_row_num_service = ""
		lldp_med_row_num_service_ref_type_ = ""
	}
	lldp_med {
		index = 2
		lldp_med_row_num_advertised_applicatio = ""
		lldp_med_row_num_dscp_mark = 0
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 0
		lldp_med_row_num_service = ""
		lldp_med_row_num_service_ref_type_ = ""
	}
	lldp_med_enable = true
	lldp_mode = "TxOnly"
	mac_limit = 100
	mac_security_mode = "disabled"
	max_allowed_unit = "pps"
	max_allowed_value = 1000
	max_bit_rate = "400000"
	multicast = true
	packet_queue = ""
	packet_queue_ref_type_ = ""
	poe_enable = false
	priority = "High"
	security_violation_action = "protect"
	stp_enable = false
	unidirectional_link_detection = false
}

resource "verity_eth_port_settings" "eth_port_settings_test_script2" {
    name = "eth_port_settings_test_script2"
    depends_on = [verity_operation_stage.eth_port_settings_stage]
	object_properties {
		group = ""
	}
	action = "Protect"
	aging_time = 0
	aging_type = "absolute"
	allocated_power = "0.0"
	auto_negotiation = true
	bpdu_filter = false
	bpdu_guard = false
	broadcast = true
	bsp_enable = false
	cli_commands = ""
	detect_bridging_loops = false
	duplex_mode = "Half"
	enable = true
	enable_speed_control = true
	fast_learning_mode = true
	fec = "unaltered"
	guard_loop = false
	lldp_enable = true
	lldp_med {
		index = 1
		lldp_med_row_num_advertised_applicatio = ""
		lldp_med_row_num_dscp_mark = 0
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 0
		lldp_med_row_num_service = ""
		lldp_med_row_num_service_ref_type_ = ""
	}
	lldp_med {
		index = 2
		lldp_med_row_num_advertised_applicatio = ""
		lldp_med_row_num_dscp_mark = 0
		lldp_med_row_num_enable = false
		lldp_med_row_num_priority = 0
		lldp_med_row_num_service = ""
		lldp_med_row_num_service_ref_type_ = ""
	}
	lldp_med_enable = true
	lldp_mode = "TxOnly"
	mac_limit = 100
	mac_security_mode = "disabled"
	max_allowed_unit = "pps"
	max_allowed_value = 1000
	max_bit_rate = "400000"
	multicast = true
	packet_queue = ""
	packet_queue_ref_type_ = ""
	poe_enable = false
	priority = "High"
	security_violation_action = "protect"
	stp_enable = false
	unidirectional_link_detection = false
}
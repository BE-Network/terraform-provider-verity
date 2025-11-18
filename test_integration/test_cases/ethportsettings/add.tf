# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
    name = "eth_port_settings_test_script1"
    depends_on = [verity_operation_stage.eth_port_settings_stage]
	object_properties {
		group = ""
	}
	action = "Protect"
	allocated_power = "0.0"
	auto_negotiation = true
	bpdu_filter = false
	bpdu_guard = false
	broadcast = true
	bsp_enable = false
	duplex_mode = "Auto"
	enable = false
	enable_ecn = true
	enable_watchdog_tuning = false
	enable_wred_tuning = false
	fast_learning_mode = true
	fec = "unaltered"
	guard_loop = false
	max_allowed_unit = "pps"
	max_allowed_value = 1000
	max_bit_rate = "-1"
	maximum_wred_threshold = 1
	minimum_wred_threshold = 1
	multicast = true
	packet_queue = ""
	packet_queue_ref_type_ = ""
	poe_enable = false
	priority = "High"
	priority_flow_control_watchdog_action = "DROP"
	priority_flow_control_watchdog_detect_time = 100
	priority_flow_control_watchdog_restore_time = 100
	single_link = false
	stp_enable = false
	wred_drop_probability = 0
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
	object_properties {
		group = "test"
	}
	auto_negotiation = false
	bpdu_filter = true
	bpdu_guard = true
	broadcast = false
	bsp_enable = true
	duplex_mode = "Half"
	enable = true
	enable_ecn = false
	enable_speed_control = false
	enable_watchdog_tuning = true
	enable_wred_tuning = true
	fast_learning_mode = false
	guard_loop = true
	max_allowed_value = 990
	maximum_wred_threshold = 3
	minimum_wred_threshold = 2
	multicast = false
	poe_enable = true
	priority_flow_control_watchdog_detect_time = 110
	priority_flow_control_watchdog_restore_time = 110
	single_link = true
	stp_enable = true
	wred_drop_probability = 1
}

resource "verity_eth_port_settings" "eth_port_settings_test_script2" {
	object_properties {
		group = "test"
	}
	auto_negotiation = false
	bpdu_filter = true
	bpdu_guard = true
	broadcast = false
	bsp_enable = true
	duplex_mode = "Half"
	enable = false
	enable_ecn = false
	enable_speed_control = true
	enable_watchdog_tuning = true
	enable_wred_tuning = true
	fast_learning_mode = false
	guard_loop = true
	max_allowed_value = 900
	maximum_wred_threshold = 5
	minimum_wred_threshold = 2
	multicast = false
	poe_enable = true
	priority_flow_control_watchdog_detect_time = 120
	priority_flow_control_watchdog_restore_time = 120
	single_link = true
	stp_enable = true
	wred_drop_probability = 4
}
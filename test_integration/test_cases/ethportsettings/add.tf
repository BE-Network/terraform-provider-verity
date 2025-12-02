# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
    name = "eth_port_settings_test_script1"
    depends_on = [verity_operation_stage.eth_port_settings_stage]
	object_properties {
		group = "test"
	}
	action = "Protect"
	allocated_power = "0.0"
	auto_negotiation = true
	bpdu_filter = false
	bpdu_guard = false
	broadcast = false
	bsp_enable = false
	duplex_mode = "Auto"
	enable = true
	fast_learning_mode = true
	fec = "unaltered"
	guard_loop = false
	max_allowed_unit = "pps"
	max_allowed_value = 1000
	max_bit_rate = "-1"
	multicast = true
	poe_enable = false
	priority = "High"
	single_link = true
	stp_enable = false
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
	object_properties {
		group = "test123"
	}
	bpdu_guard = true
	broadcast = true
	bsp_enable = true
	enable = false
	multicast = false
	poe_enable = false
	single_link = false
	stp_enable = true
}
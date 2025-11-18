# Test case: Modify existing resources
# Define modified versions of the resources from add.tf


resource "verity_eth_port_settings" "eth_port_settings_test_script1" {
	object_properties {
		group = "test"
	}
	enable = true
	enable_ecn = false
	enable_watchdog_tuning = true
}
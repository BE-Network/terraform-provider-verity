# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_settings" "device_settings_test_script1" {
	object_properties {
		group = "test1"
	}
	commit_to_flash_interval = null
	cut_through_switching = true
	disable_tcp_udp_learned_packet_acceleration = true
	enable = false
	external_battery_power_available = 30
	external_power_available = 70
	rocev2 = true
	security_audit_interval = null
	usage_threshold = 0.98
}

resource "verity_device_settings" "device_settings_test_script2" {
	object_properties {
		group = "test123"
	}
	commit_to_flash_interval = 0
	cut_through_switching = false
	disable_tcp_udp_learned_packet_acceleration = false
	enable = true
	external_battery_power_available = 40
	external_power_available = 60
	rocev2 = false
	security_audit_interval = 0
	usage_threshold = 0.90
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_settings" "device_settings_test_script1" {
	object_properties {
		group = "test123"
	}
	commit_to_flash_interval = 60
	cut_through_switching = true
	disable_tcp_udp_learned_packet_acceleration = true
	enable = false
	external_battery_power_available = 42
	external_power_available = 80
	hold_timer = 1
	mac_aging_timer_override = 2
	rocev2 = true
	security_audit_interval = 71
	usage_threshold = 0.96
}

resource "verity_device_settings" "device_settings_test_script2" {
	object_properties {
		group = "test"
	}
	commit_to_flash_interval = 80
	cut_through_switching = true
	disable_tcp_udp_learned_packet_acceleration = true
	enable = false
	external_battery_power_available = 44
	external_power_available = 62
	hold_timer = 2
	mac_aging_timer_override = 4
	rocev2 = true
	security_audit_interval = 61
	usage_threshold = 0.90
}
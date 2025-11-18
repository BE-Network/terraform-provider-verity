# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_device_settings" "device_settings_test_script1" {
    name = "device_settings_test_script1"
    depends_on = [verity_operation_stage.device_settings_stage]
	object_properties {
		group = ""
	}
	commit_to_flash_interval = null
	cut_through_switching = false
	disable_tcp_udp_learned_packet_acceleration = false
	enable = true
	external_battery_power_available = 40
	external_power_available = 75
	mode = "IEEE 802.3af"
	packet_queue = "(Packet Queue)"
	packet_queue_ref_type_ = "packet_queue"
	rocev2 = false
	security_audit_interval = null
	usage_threshold = 0.99
}
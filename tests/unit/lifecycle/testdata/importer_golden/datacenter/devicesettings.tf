
resource "verity_device_settings" "device_settings_test1" {
    name = "device_settings_test1"
    depends_on = [verity_operation_stage.device_settings_stage]
	object_properties {
	}
	commit_to_flash_interval = 0
	cut_through_switching = false
	disable_tcp_udp_learned_packet_acceleration = false
	enable = true
	external_battery_power_available = 40
	external_power_available = 75
	mode = "IEEE 802.3af"
	packet_queue = "(Packet Queue)"
	packet_queue_ref_type_ = "packet_queue"
	rocev2 = false
	security_audit_interval = 0
	usage_threshold = 0.95
}


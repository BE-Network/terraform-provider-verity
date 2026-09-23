
resource "verity_device_settings" "_Device_Settings_" {
    name = "(Device Settings)"
    depends_on = [verity_operation_stage.device_settings_stage]
	object_properties {
	}
	commit_to_flash_interval = null
	cut_through_switching = false
	disable_tcp_udp_learned_packet_acceleration = false
	enable = true
	external_battery_power_available = 0
	external_power_available = 0
	hold_timer = 0
	mac_aging_timer_override = null
	mode = "IEEE 802.3af"
	packet_queue = "(Packet Queue)"
	packet_queue_ref_type_ = "packet_queue"
	rocev2 = false
	security_audit_interval = null
	spanning_tree_priority = "byLevel"
	usage_threshold = 0.99
}


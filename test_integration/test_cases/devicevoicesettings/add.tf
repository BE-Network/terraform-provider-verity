# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_device_voice_settings" "device_voice_settings_test_script1" {
    name = "device_voice_settings_test_script1"
    depends_on = [verity_operation_stage.device_voice_setting_stage]
	object_properties {
		group = ""
	}
	anon_cid_block_activate = ""
	anon_cid_block_deactivate = ""
	bit_rate = "14400"
	call_agent_1 = ""
	call_agent_2 = ""
	call_agent_port_1 = 0
	call_agent_port_2 = 0
	call_forward_on_busy_activate = ""
	call_forward_on_busy_deactivate = ""
	call_forward_on_no_answer_activate = ""
	call_forward_on_no_answer_deactivate = ""
	call_forward_unconditional_activate = ""
	call_forward_unconditional_deactivate = ""
	call_hold = ""
	cancel_call_waiting = ""
	cas_events = 0
	cids_activate = ""
	cids_deactivate = ""
	do_not_disturb_activate = ""
	do_not_disturb_deactivate = ""
	do_not_disturb_pin_change = ""
	domain = ""
	dscp_mark = 0
	dtmf_method = "Inband"
	emergency_service_number = "911"
	enable = true
	event_payload_type = 101
	fax_t38 = false
	intercom_1 = ""
	intercom_2 = ""
	intercom_3 = ""
	local_port_max = 30200
	local_port_min = 30000
	mgcp_dscp_mark = 0
	outbound_proxy = ""
	outbound_proxy_port = 0
	outbound_proxy_secondary = ""
	outbound_proxy_secondary_port = 0
	protocol = "SIP"
	proxy_server = ""
	proxy_server_port = 0
	proxy_server_secondary = ""
	proxy_server_secondary_port = 0
	region = "US"
	register_expires = 3600
	registrar_server = ""
	registrar_server_port = 0
	registrar_server_secondary = ""
	registrar_server_secondary_port = 0
	registration_period = 3240
	rtcp = true
	sip_dscp_mark = 0
	termination_base = ""
	user_agent_domain = ""
	user_agent_port = 0
	user_agent_transport = "UDP"
	voicemail_server = ""
	voicemail_server_expires = 3600
	voicemail_server_port = 0
}
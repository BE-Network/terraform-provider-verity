# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_device_voice_settings" "device_voice_settings_test_script1" {
    name = "device_voice_settings_test_script1"
    depends_on = [verity_operation_stage.device_voice_setting_stage]
	object_properties {
		group = "test"
	}
	anon_cid_block_activate = "*77"
	anon_cid_block_deactivate = "*87"
	bit_rate = "14400"
	call_agent_1 = ""
	call_agent_2 = ""
	call_agent_port_1 = 0
	call_agent_port_2 = 0
	call_forward_on_busy_activate = "*90"
	call_forward_on_busy_deactivate = "*91"
	call_forward_on_no_answer_activate = "*92"
	call_forward_on_no_answer_deactivate = "*93"
	call_forward_unconditional_activate = "*72"
	call_forward_unconditional_deactivate = "*73"
	call_hold = "*9"
	cancel_call_waiting = "*70"
	cas_events = 0
	cids_activate = "*67"
	cids_deactivate = "*82"
	codecs {
		index = 1
		codec_num_enable = true
		codec_num_name = "G.711MuLaw"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 2
		codec_num_enable = true
		codec_num_name = "G.711ALaw"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 3
		codec_num_enable = true
		codec_num_name = "G.726"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 4
		codec_num_enable = true
		codec_num_name = "G.722"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 5
		codec_num_enable = true
		codec_num_name = "G.729"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	do_not_disturb_activate = "*78"
	do_not_disturb_deactivate = "*79"
	do_not_disturb_pin_change = "*10"
	domain = ""
	dscp_mark = 0
	dtmf_method = "Inband"
	emergency_service_number = "911"
	enable = true
	event_payload_type = 101
	fax_t38 = false
	intercom_1 = "*53"
	intercom_2 = "*54"
	intercom_3 = "*55"
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
	termination_base = "aaln/"
	user_agent_domain = ""
	user_agent_port = 0
	user_agent_transport = "UDP"
	voicemail_server = ""
	voicemail_server_expires = 3600
	voicemail_server_port = 0
}


resource "verity_device_voice_settings" "device_voice_settings_test_script2" {
    name = "device_voice_settings_test_script2"
    depends_on = [verity_operation_stage.device_voice_setting_stage]
	object_properties {
		group = ""
	}
	anon_cid_block_activate = "*77"
	anon_cid_block_deactivate = "*87"
	bit_rate = "14400"
	call_agent_1 = ""
	call_agent_2 = ""
	call_agent_port_1 = 0
	call_agent_port_2 = 0
	call_forward_on_busy_activate = "*90"
	call_forward_on_busy_deactivate = "*91"
	call_forward_on_no_answer_activate = "*92"
	call_forward_on_no_answer_deactivate = "*93"
	call_forward_unconditional_activate = "*72"
	call_forward_unconditional_deactivate = "*73"
	call_hold = "*9"
	cancel_call_waiting = "*70"
	cas_events = 0
	cids_activate = "*67"
	cids_deactivate = "*82"
	codecs {
		index = 1
		codec_num_enable = true
		codec_num_name = "G.711MuLaw"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 2
		codec_num_enable = true
		codec_num_name = "G.711ALaw"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 3
		codec_num_enable = true
		codec_num_name = "G.726"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 4
		codec_num_enable = true
		codec_num_name = "G.722"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	codecs {
		index = 5
		codec_num_enable = true
		codec_num_name = "G.729"
		codec_num_packetization_period = "20"
		codec_num_silence_suppression = false
	}
	do_not_disturb_activate = "*78"
	do_not_disturb_deactivate = "*79"
	do_not_disturb_pin_change = "*10"
	domain = ""
	dscp_mark = 0
	dtmf_method = "Inband"
	emergency_service_number = "911"
	enable = true
	event_payload_type = 101
	fax_t38 = false
	intercom_1 = "*53"
	intercom_2 = "*54"
	intercom_3 = "*55"
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
	termination_base = "aaln/"
	user_agent_domain = ""
	user_agent_port = 0
	user_agent_transport = "UDP"
	voicemail_server = ""
	voicemail_server_expires = 3600
	voicemail_server_port = 0
}

# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_voice_port_profile" "voice_port_profile_test_script1" {
    name = "voice_port_profile_test_script1"
    depends_on = [verity_operation_stage.voice_port_profile_stage]
	object_properties {
		format_dial_plan = true
		group = ""
		port_monitoring = ""
	}
	anonymous_call_block_enable = false
	audio_mwi_enable = false
	call_forward_on_busy_enable = false
	call_forward_on_no_answer_ring_count = 4
	call_forward_unconditional_enable = false
	call_hold_enable = false
	call_three_way_enable = false
	call_transfer_enable = false
	call_waiting_caller_id_enable = false
	call_waiting_enable = false
	caller_id_enable = false
	caller_id_name_enable = false
	cid_blocking_enable = false
	cid_name_presentation_status = "Public"
	cid_num_presentation_status = "Public"
	dial_tone_feature_delay = 4
	digit_map = "(T)"
	do_not_disturb_enable = false
	echo_cancellation_enable = true
	enable = false
	hotline_enable = false
	intercom_enable = false
	intercom_transfer_enable = false
	jitter_buffer_max = 180
	jitter_target = 40
	mwi_refresh_timer = 30
	protocol = "SIP"
	receive_gain = -30
	release_timer = 10
	roh_timer = 15
	signaling_code = "LoopStart"
	transmit_gain = -30
	visual_mwi_enable = false
}
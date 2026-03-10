# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_voice_port_profile" "voice_port_profile_test_script1" {
	object_properties {
		format_dial_plan = true
		group = "test"
		port_monitoring = "test"
	}
	anonymous_call_block_enable = true
	audio_mwi_enable = true
	call_forward_on_busy_enable = true
	call_forward_on_no_answer_ring_count = 6
	call_forward_unconditional_enable = true
	call_hold_enable = true
	call_three_way_enable = true
	call_transfer_enable = true
	call_waiting_caller_id_enable = true
	call_waiting_enable = true
	caller_id_enable = true
	caller_id_name_enable = true
	cid_blocking_enable = true
	dial_tone_feature_delay = 6
	do_not_disturb_enable = true
	echo_cancellation_enable = false
	enable = false
	hotline_enable = true
	intercom_enable = true
	intercom_transfer_enable = true
	jitter_buffer_max = 182
	jitter_target = 42
	mwi_refresh_timer = 32
	receive_gain = -28
	release_timer = 12
	roh_timer = 17
	transmit_gain = -28
	visual_mwi_enable = true
}

resource "verity_voice_port_profile" "voice_port_profile_test_script2" {
	object_properties {
		format_dial_plan = false
		group = ""
		port_monitoring = ""
	}
	anonymous_call_block_enable = true
	audio_mwi_enable = true
	call_forward_on_busy_enable = true
	call_forward_on_no_answer_ring_count = 6
	call_forward_unconditional_enable = true
	call_hold_enable = true
	call_three_way_enable = true
	call_transfer_enable = true
	call_waiting_caller_id_enable = true
	call_waiting_enable = true
	caller_id_enable = true
	caller_id_name_enable = true
	cid_blocking_enable = true
	dial_tone_feature_delay = 6
	do_not_disturb_enable = true
	echo_cancellation_enable = false
	enable = false
	hotline_enable = true
	intercom_enable = true
	intercom_transfer_enable = true
	jitter_buffer_max = 182
	jitter_target = 42
	mwi_refresh_timer = 32
	receive_gain = -28
	release_timer = 12
	roh_timer = 17
	transmit_gain = -28
	visual_mwi_enable = true
}
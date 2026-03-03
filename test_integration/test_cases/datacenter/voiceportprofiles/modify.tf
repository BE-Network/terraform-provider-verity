# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_voice_port_profile" "voice_port_profile_test_script1" {
	object_properties {
		format_dial_plan = false
		group = "test"
	}
	anonymous_call_block_enable = true
	audio_mwi_enable = true
}
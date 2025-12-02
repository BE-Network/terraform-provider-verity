# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_voice_settings" "device_voice_settings_test_script1" {
	object_properties {
		group = "test"
	}
	enable = false
	register_expires = 3660
	rtcp = false
}
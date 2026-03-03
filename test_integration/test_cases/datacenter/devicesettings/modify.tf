# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_settings" "device_settings_test_script1" {
	object_properties {
		group = "test"
	}
	enable = false
}
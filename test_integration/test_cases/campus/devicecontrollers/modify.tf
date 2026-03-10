# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_controller" "device_controller_test_script1" {
	enable = true
	managed_on_native_vlan = true
	uses_tagged_packets = false
}

resource "verity_device_controller" "device_controller_test_script2" {
	enable = true
	managed_on_native_vlan = true
	uses_tagged_packets = false
}
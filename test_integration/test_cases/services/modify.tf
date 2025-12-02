# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service" "service_test_script1" {
	object_properties {
		group = "test123"
	}
	enable = false
	mtu = 1501
}
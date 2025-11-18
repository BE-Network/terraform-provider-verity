# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service" "service_test_script1" {
	object_properties {
		group = "test"
		on_summary = false
	}
	enable = false
	mtu = 1501
	vlan = 105
}
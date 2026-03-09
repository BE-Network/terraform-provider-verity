# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service" "service_test_script1" {
	object_properties {
		group = ""
		on_summary = false
	}
	mtu = 1500
	vlan = 119
}

resource "verity_service" "service_test_script2" {
	object_properties {
		group = "test123"
		on_summary = false
	}
	mtu = 1501
	vlan = 120
}
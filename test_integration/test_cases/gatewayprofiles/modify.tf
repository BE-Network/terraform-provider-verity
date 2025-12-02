# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_gateway_profile" "gateway_profile_test_script1" {

	object_properties {
		group = "test"
	}
	enable = true
	external_gateways {
		index = 2
		peer_gw = true
	}
	tenant_slice_managed = true
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_gateway_profile" "gateway_profile_test_script1" {
	object_properties {
		group = "test"
	}
	enable = false
	external_gateways {
		index = 1
		enable = false
		peer_gw = true
	}
	external_gateways {
		index = 2
	}
}

resource "verity_gateway_profile" "gateway_profile_test_script2" {
	object_properties {
		group = ""
	}
	enable = true
	external_gateways {
		index = 3
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = false
		source_ip_mask = ""
	}
}
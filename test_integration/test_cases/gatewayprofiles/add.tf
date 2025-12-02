# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_gateway_profile" "gateway_profile_test_script1" {
    name = "gateway_profile_test_script1"
    depends_on = [verity_operation_stage.gateway_profile_stage]
	object_properties {
		group = ""
	}
	enable = false
	external_gateways {
		index = 1
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = true
		source_ip_mask = ""
	}
	external_gateways {
		index = 2
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = false
		source_ip_mask = ""
	}
	tenant_slice_managed = false
}
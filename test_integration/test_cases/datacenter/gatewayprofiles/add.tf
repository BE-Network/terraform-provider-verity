# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_gateway_profile" "gateway_profile_test_script1" {
    name = "gateway_profile_test_script1"
    depends_on = [verity_operation_stage.gateway_profile_stage]
	object_properties {
		group = ""
	}
	enable = true
	external_gateways {
		index = 1
		enable = true
		gateway = "EG_service_test"
		gateway_ref_type_ = "gateway"
		peer_gw = false
		source_ip_mask = "1.2.3.5/16"
	}
	external_gateways {
		index = 2
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = false
		source_ip_mask = ""
	}
}
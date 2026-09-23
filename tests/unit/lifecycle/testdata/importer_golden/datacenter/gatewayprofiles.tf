
resource "verity_gateway_profile" "test" {
    name = "test"
    depends_on = [verity_operation_stage.gateway_profile_stage]
	enable = true
	external_gateways {
		index = 1
		enable = true
		gateway = "gateway_test1"
		gateway_ref_type_ = "gateway"
		peer_gw = false
		source_ip_mask = "10.2.2.2/24"
	}
	external_gateways {
		index = 2
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = false
		source_ip_mask = ""
	}
	external_gateways {
		index = 3
		enable = false
		gateway = ""
		gateway_ref_type_ = ""
		peer_gw = false
		source_ip_mask = ""
	}
}


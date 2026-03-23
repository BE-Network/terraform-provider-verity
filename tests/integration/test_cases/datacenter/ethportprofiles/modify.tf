# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
	object_properties {
		group = "test"
	}
	enable = false
	services {
		index = 1
		row_num_enable = false
	}
	services {
		index = 2
	}
}

resource "verity_eth_port_profile" "eth_port_profile_test_script2" {
	object_properties {
		group = ""
	}
	enable = true
	services {
		index = 2
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = false
		row_num_external_vlan = null
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = "service"
		row_num_service_ref_type_ = "service"
	}
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
	object_properties {
		group = ""
		port_monitoring = ""
	}
	enable = false
	services {
		index = 3
		row_num_enable = true
		row_num_external_vlan = 14
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	services {
		index = 4
		row_num_enable = false
		row_num_external_vlan = 15
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	tenant_slice_managed = true
}
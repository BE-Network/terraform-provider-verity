# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
    name = "eth_port_profile_test_script1"
    depends_on = [verity_operation_stage.eth_port_profile_stage]
	object_properties {
		group = "test"
		port_monitoring = "high"
	}
	enable = true
	services {
		index = 1
		row_num_enable = true
		row_num_external_vlan = 12
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	services {
		index = 2
		row_num_enable = false
		row_num_external_vlan = 13
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	tenant_slice_managed = false
}
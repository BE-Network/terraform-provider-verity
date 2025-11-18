# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
    name = "eth_port_profile_test_script1"
    depends_on = [verity_operation_stage.eth_port_profile_stage]
	object_properties {
		group = ""
		port_monitoring = "high"
	}
	egress_acl = ""
	egress_acl_ref_type_ = ""
	enable = true
	ingress_acl = ""
	ingress_acl_ref_type_ = ""
	services {
		index = 1
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = true
		row_num_external_vlan = 33
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = "service"
		row_num_service_ref_type_ = "service"
	}
}
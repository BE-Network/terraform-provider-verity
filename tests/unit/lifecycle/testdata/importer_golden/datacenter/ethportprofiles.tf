
resource "verity_eth_port_profile" "TEST" {
    name = "TEST"
    depends_on = [verity_operation_stage.eth_port_profile_stage]
	object_properties {
		port_monitoring = "high"
	}
	egress_acl = ""
	egress_acl_ref_type_ = ""
	enable = false
	ingress_acl = ""
	ingress_acl_ref_type_ = ""
	services {
		index = 1
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = false
		row_num_external_vlan = null
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	services {
		index = 2
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = false
		row_num_external_vlan = null
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	services {
		index = 3
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = false
		row_num_external_vlan = null
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	services {
		index = 4
		row_num_egress_acl = ""
		row_num_egress_acl_ref_type_ = ""
		row_num_enable = false
		row_num_external_vlan = null
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
}


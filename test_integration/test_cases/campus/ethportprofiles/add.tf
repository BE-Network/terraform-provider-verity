# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
    name = "eth_port_profile_test_script1"
    depends_on = [verity_operation_stage.eth_port_profile_stage]
	object_properties {
		group = ""
		icon = ""
		label = ""
		port_monitoring = ""
		sort_by_name = false
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
		row_num_external_vlan = 0
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	tls = false
	tls_service = ""
	tls_service_ref_type_ = ""
	trusted_port = false
}

resource "verity_eth_port_profile" "eth_port_profile_test_script2" {
    name = "eth_port_profile_test_script2"
    depends_on = [verity_operation_stage.eth_port_profile_stage]
	object_properties {
		group = ""
		icon = ""
		label = ""
		port_monitoring = ""
		sort_by_name = false
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
		row_num_external_vlan = 0
		row_num_ingress_acl = ""
		row_num_ingress_acl_ref_type_ = ""
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
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
		row_num_lan_iptv = ""
		row_num_mac_filter = ""
		row_num_mac_filter_ref_type_ = ""
		row_num_service = ""
		row_num_service_ref_type_ = ""
	}
	tls = false
	tls_service = ""
	tls_service_ref_type_ = ""
	trusted_port = false
}
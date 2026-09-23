
resource "verity_service" "service_test1" {
    name = "service_test1"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {}
	anycast_ipv4_mask = "10.12.14.16/24"
	anycast_ipv6_mask = ""
	dhcp_server_ipv4 = ""
	dhcp_server_ipv6 = ""
	enable = false
	mtu = 1501
	policy_based_routing = ""
	policy_based_routing_ref_type_ = ""
	tenant = "Changes_EthPort-Tenant"
	tenant_ref_type_ = "tenant"
	vlan = 108
	vni_auto_assigned_ = true
}


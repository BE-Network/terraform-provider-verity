# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_service" "service_test_script1" {
    name = "service_test_script1"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {
		group = "test"
		on_summary = true
	}
	anycast_ipv4_mask = "10.12.14.16/24"
	anycast_ipv6_mask = ""
	dhcp_server_ipv4 = ""
	dhcp_server_ipv6 = ""
	enable = false
	mtu = null
	policy_based_routing = ""
	policy_based_routing_ref_type_ = ""
	tenant = "Changes_EthPort-Tenant"
	tenant_ref_type_ = "tenant"
	vlan = 117
	vni_auto_assigned_ = true
}

resource "verity_service" "service_test_script2" {
    name = "service_test_script2"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {
		group = "test"
		on_summary = true
	}
	anycast_ipv4_mask = "10.12.14.16/24"
	anycast_ipv6_mask = ""
	dhcp_server_ipv4 = ""
	dhcp_server_ipv6 = ""
	enable = false
	mtu = 0
	policy_based_routing = ""
	policy_based_routing_ref_type_ = ""
	tenant = "Changes_EthPort-Tenant"
	tenant_ref_type_ = "tenant"
	vlan = 118
	vni_auto_assigned_ = true
}
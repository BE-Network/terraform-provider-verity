# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_service" "service_test_script1" {
    name = "service_test_script1"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {
		group = "test"
	}
	anycast_ip_mask = ""
	dhcp_server_ip = ""
	enable = true
	mtu = 1502
	tenant = "test2"
	tenant_ref_type_ = "tenant"
	vlan = 4041
	vni_auto_assigned_ = true
}
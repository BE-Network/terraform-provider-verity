# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_tenant" "tenant_test_script1" {
    name = "tenant_test_script1"
    depends_on = [verity_operation_stage.tenant_stage]
	object_properties {
		group = ""
	}
	default_originate = false
	dhcp_relay_source_ipv4s_subnet = ""
	dhcp_relay_source_ipv6s_subnet = ""
	enable = true
	export_route_map = ""
	export_route_map_ref_type_ = ""
	import_route_map = ""
	import_route_map_ref_type_ = ""
	layer_3_vlan_auto_assigned_ = true
	layer_3_vni_auto_assigned_ = true
	route_distinguisher = ""
	route_target_export = ""
	route_target_import = ""
	route_tenants {
		index = 1
		enable = false
		tenant = ""
	}
	vrf_name_auto_assigned_ = true
}
# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_service" "service_test_script1" {
    name = "service_test_script1"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {
		group = ""
		on_summary = true
		warn_on_no_external_source = false
	}
	act_as_multicast_querier = false
	allow_fast_leave = false
	allow_local_switching = true
	block_downstream_dhcp_server = true
	block_unknown_unicast_flood = false
	enable = true
	is_management_service = false
	max_downstream_rate_mbps = null
	max_upstream_rate_mbps = null
	mst_instance = 0
	multicast_management_mode = "flooding"
	packet_priority = "0"
	policy_based_routing = ""
	policy_based_routing_ref_type_ = ""
	tagged_packets = false
	tls = false
	use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
	vlan = 104
}

resource "verity_service" "service_test_script2" {
    name = "service_test_script2"
    depends_on = [verity_operation_stage.service_stage]
	object_properties {
		group = ""
		on_summary = true
		warn_on_no_external_source = false
	}
	act_as_multicast_querier = false
	allow_fast_leave = false
	allow_local_switching = true
	block_downstream_dhcp_server = true
	block_unknown_unicast_flood = false
	enable = true
	is_management_service = false
	max_downstream_rate_mbps = null
	max_upstream_rate_mbps = null
	mst_instance = 0
	multicast_management_mode = "flooding"
	packet_priority = "0"
	policy_based_routing = ""
	policy_based_routing_ref_type_ = ""
	tagged_packets = false
	tls = false
	use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
	vlan = 105
}

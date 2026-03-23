# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service" "service_test_script1" {
	object_properties {
		group = "test"
		on_summary = false
		warn_on_no_external_source = false
	}
	act_as_multicast_querier = true
	allow_fast_leave = true
	allow_local_switching = false
	block_downstream_dhcp_server = false
	block_unknown_unicast_flood = true
	enable = false
	is_management_service = true
	max_downstream_rate_mbps = 2200
	max_upstream_rate_mbps = 1100
	mst_instance = 1
	tagged_packets = false
	tls = false
	use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
	vlan = 106
}

resource "verity_service" "service_test_script2" {
	object_properties {
		group = "test"
		on_summary = false
		warn_on_no_external_source = false
	}
	act_as_multicast_querier = true
	allow_fast_leave = true
	allow_local_switching = false
	block_downstream_dhcp_server = false
	block_unknown_unicast_flood = true
	enable = false
	is_management_service = true
	max_downstream_rate_mbps = 2000
	max_upstream_rate_mbps = 1000
	mst_instance = 1
	tagged_packets = false
	tls = false
	use_dscp_to_p_bit_mapping_for_l3_packets_if_available = false
	vlan = 107
}
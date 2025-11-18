# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_route_map_clause" "route_map_clause_test_script1" {
    name = "route_map_clause_test_script1"
    depends_on = [verity_operation_stage.route_map_clause_stage]
	object_properties {
		match_fields_shown = "match_ipv6_address_prefix_list_path"
		notes = ""
	}
	enable = true
	match_as_path_access_list = ""
	match_as_path_access_list_ref_type_ = ""
	match_community_list = ""
	match_community_list_ref_type_ = ""
	match_evpn_route_type = ""
	match_evpn_route_type_default = null
	match_extended_community_list = ""
	match_extended_community_list_ref_type_ = ""
	match_interface_number = null
	match_interface_vlan = null
	match_ipv4_address_ip_prefix_list = ""
	match_ipv4_address_ip_prefix_list_ref_type_ = ""
	match_ipv4_next_hop_ip_prefix_list = ""
	match_ipv4_next_hop_ip_prefix_list_ref_type_ = ""
	match_ipv6_address_ipv6_prefix_list = "ipv6_list"
	match_ipv6_address_ipv6_prefix_list_ref_type_ = "ipv6_prefix_list"
	match_ipv6_next_hop_ipv6_prefix_list = ""
	match_ipv6_next_hop_ipv6_prefix_list_ref_type_ = ""
	match_local_preference = null
	match_metric = null
	match_origin = ""
	match_peer_interface = null
	match_peer_ip_address = ""
	match_peer_vlan = null
	match_source_protocol = ""
	match_tag = null
	match_vni = null
	match_vrf = ""
	match_vrf_ref_type_ = ""
	permit_deny = "permit"
}
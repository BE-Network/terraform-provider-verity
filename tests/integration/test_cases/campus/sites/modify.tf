# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_site" "TOR_Complex" {
	object_properties {
        system_graphs {
			graph_num_data = "graph_endpoint_path=::graph_group_name=::graph_lagg_path=::graph_line_set=::graph_processes=::graph_source_mod=::graph_source_path=::graph_sub_path=::graph_time_frame="
			index = 6
		}
	}
	aggressive_reporting = true
	duplicate_address_detection_time = 182
	enable = false
	enable_dhcp_snooping = true
	evpn_mac_holdtime = 1082
	evpn_multihoming_startup_delay = 310
	force_spanning_tree_on_fabric_ports = true
	ip_source_guard = true
	islands {
		index = 3
		toi_switchpoint = "g"
		toi_switchpoint_ref_type_ = "switchpoint"
	}
	link_state_timeout_value = 62
	read_only_mode = true
}
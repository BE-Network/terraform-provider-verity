
resource "verity_fabric" "TOR_Complex" {
    name = "TOR Complex"
    depends_on = [verity_operation_stage.fabric_stage]
	object_properties {
		system_graphs {
			index = 1
		}
		system_graphs {
			index = 2
		}
		system_graphs {
			index = 3
		}
		system_graphs {
			index = 4
		}
		system_graphs {
			index = 5
		}
	}
	aggressive_reporting = false
	dscp_to_p_bit_map = "0000000011111111222222223333333344444444555555556666666677777777"
	duplicate_address_detection_max_number_of_moves = 5
	duplicate_address_detection_time = 180
	enable = true
	enable_dhcp_snooping = false
	evpn_mac_holdtime = 1080
	evpn_multihoming_startup_delay = 300
	force_spanning_tree_on_fabric_ports = false
	ip_source_guard = false
	link_state_timeout_value = 60
	read_only_mode = false
	region_name = ""
	revision = 0
	service_for_fabric = "Management"
	service_for_fabric_ref_type_ = "service"
	spanning_tree_type = "pvst"
}


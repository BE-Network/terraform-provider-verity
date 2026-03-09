# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_switchpoint" "switchpoint_test_script1" {
    name = "switchpoint_test_script1"
    depends_on = [verity_operation_stage.switchpoint_stage]
	object_properties {
		aggregate = false
		expected_parent_endpoint = ""
		is_host = false
		number_of_multipoints = 0
		user_notes = ""
	}
	bgp_as_number_auto_assigned_ = true
	connected_bundle = "Bundle for t4"
	connected_bundle_ref_type_ = "endpoint_bundle"
	device_serial_number = ""
	enable = true
	eths {
		index = 1
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 2
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 3
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 4
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 5
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 6
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 7
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 8
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 9
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 10
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 11
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 12
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 13
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 14
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 15
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 16
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 17
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 18
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 19
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 20
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 21
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 22
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 23
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 24
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	locked = false
	out_of_band_management = true
	pod = "Pod 2"
	pod_ref_type_ = "pod"
	rack = ""
	read_only_mode = true
	spine_plane = ""
	spine_plane_ref_type_ = ""
	switch_router_id_ip_mask_auto_assigned_ = true
	switch_vtep_id_ip_mask_auto_assigned_ = true
	type = "leaf"
}


resource "verity_switchpoint" "switchpoint_test_script2" {
    name = "switchpoint_test_script2"
    depends_on = [verity_operation_stage.switchpoint_stage]
	object_properties {
		aggregate = false
		expected_parent_endpoint = ""
		is_host = false
		number_of_multipoints = 0
		user_notes = ""
	}
	bgp_as_number_auto_assigned_ = true
	connected_bundle = "Bundle for t4"
	connected_bundle_ref_type_ = "endpoint_bundle"
	device_serial_number = ""
	enable = true
	eths {
		index = 1
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 2
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 3
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 4
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 5
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 6
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 7
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 8
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 9
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 10
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 11
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 12
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 13
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 14
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 15
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 16
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 17
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 18
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 19
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 20
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 21
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 22
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 23
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	eths {
		index = 24
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	locked = false
	out_of_band_management = true
	pod = "Pod 2"
	pod_ref_type_ = "pod"
	rack = ""
	read_only_mode = true
	spine_plane = ""
	spine_plane_ref_type_ = ""
	switch_router_id_ip_mask_auto_assigned_ = true
	switch_vtep_id_ip_mask_auto_assigned_ = true
	type = "leaf"
}
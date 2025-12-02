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
	connected_bundle = "Bundle for test"
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
	locked = false
	out_of_band_management = false
	pod = ""
	pod_ref_type_ = ""
	rack = ""
	read_only_mode = true
	spine_plane = ""
	spine_plane_ref_type_ = ""
	switch_router_id_ip_mask_auto_assigned_ = true
	switch_vtep_id_ip_mask_auto_assigned_ = true
	type = "management"
}
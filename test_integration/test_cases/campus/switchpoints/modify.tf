# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_switchpoint" "switchpoint_test_script1" {
	object_properties {
		aggregate = true
		draw_as_edge_device = true
		expected_parent_endpoint = ""
		is_host = true
		number_of_multipoints = 2
		user_notes = "test"
	}
	enable = false
	eths {
		index = 2
		enable = false
	}
	eths {
		index = 9
		enable = false
	}
	is_fabric = true
	locked = true
	out_of_band_management = true
	read_only_mode = true
}

resource "verity_switchpoint" "switchpoint_test_script2" {
	object_properties {
		aggregate = true
		draw_as_edge_device = true
		expected_parent_endpoint = ""
		is_host = true
		number_of_multipoints = 1
		user_notes = "test1"
	}
	enable = false
	eths {
		index = 2
		enable = false
	}
	eths {
		index = 8
		enable = false
	}
	is_fabric = true
	locked = true
	out_of_band_management = true
	read_only_mode = true
}
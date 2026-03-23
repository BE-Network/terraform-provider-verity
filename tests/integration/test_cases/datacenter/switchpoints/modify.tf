# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_switchpoint" "switchpoint_test_script1" {
	enable = false
	eths {
		index = 2
		enable = false
		port_name = "1/2"
	}
	eths {
		index = 24		
		enable = false
		port_name = "1/24"
	}
	locked = true
	out_of_band_management = false
	read_only_mode = false
}

resource "verity_switchpoint" "switchpoint_test_script2" {
	object_properties {
		aggregate = true
		user_notes = "test"
	}
	enable = false
	eths {
		index = 2
		enable = false
		port_name = "1/2"
	}
	eths {
		index = 24
		enable = false
		port_name = "1/24"
	}
	out_of_band_management = false
	read_only_mode = false
}
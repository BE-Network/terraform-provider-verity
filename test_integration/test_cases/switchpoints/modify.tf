# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_switchpoint" "switchpoint_test_script1" {
	object_properties {
		aggregate = true
	}
	enable = false
	eths {
		index = 2
		breakout = ""
		enable = true
		eth_num_icon = "empty"
		eth_num_label = ""
		port_name = ""
	}
	locked = true
	out_of_band_management = true
	read_only_mode = false
}
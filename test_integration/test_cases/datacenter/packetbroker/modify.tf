# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_packet_broker" "packet_broker_test_script1" {
	enable = false
	ipv4_deny {
		index = 2
		enable = true
		filter = "ipv4_test_list"
		filter_ref_type_ = "ipv4_list_filter"
	}
}
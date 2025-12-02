# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_pb_routing_acl" "pb_routing_acl_test_script1" {
	enable = false
	ipv4_deny {
		index = 2
		enable = true
		filter = "filter1"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv4_permit {
		index = 2
		enable = false
		filter = "filter2"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv6_deny {
		index = 2
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_permit {
		index = 2
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
}
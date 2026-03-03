# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_pb_routing_acl" "pb_routing_acl_test_script1" {
    name = "pb_routing_acl_test_script1"
    depends_on = [verity_operation_stage.pb_routing_acl_stage]
	enable = true
	ipv4_deny {
		index = 1
		enable = true
		filter = "filter1"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv4_permit {
		index = 1
		enable = false
		filter = "filter2"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv6_deny {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_permit {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv_protocol = "ipv4"
	next_hop_ips = "20.20.20.20"
}
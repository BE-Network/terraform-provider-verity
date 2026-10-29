
resource "verity_packet_broker" "prof" {
    name = "prof"
    depends_on = [verity_operation_stage.packet_broker_stage]
	enable = true
	ipv4_deny {
		index = 1
		enable = true
		filter = "ipv4_test_list"
		filter_ref_type_ = "ipv4_list_filter"
	}
	ipv4_deny {
		index = 2
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv4_permit {
		index = 1
		enable = true
		filter = "ipv4_test_2"
		filter_ref_type_ = "ipv4_list_filter"
	}
	ipv6_deny {
		index = 1
		enable = true
		filter = "ipv6_list_test"
		filter_ref_type_ = "ipv6_list_filter"
	}
	ipv6_permit {
		index = 1
		enable = true
		filter = "ipv6_list_test"
		filter_ref_type_ = "ipv6_list_filter"
	}
}


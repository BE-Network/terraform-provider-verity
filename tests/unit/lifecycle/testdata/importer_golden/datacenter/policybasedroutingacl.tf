
resource "verity_pb_routing_acl" "pbr_acl_ipv4_2" {
    name = "pbr_acl_ipv4_2"
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
	next_hop_ips = "20.20.20.20"
}


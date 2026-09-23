
resource "verity_port_acl" "Port_acl_1" {
    name = "Port_acl_1"
    depends_on = [verity_operation_stage.port_acl_stage]
	enable = true
	ipv4_deny {
		index = 1
		enable = true
		filter = "ipv4_filter"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv4_permit {
		index = 1
		enable = true
		filter = "test4"
		filter_ref_type_ = "ipv4_filter"
	}
	ipv6_deny {
		index = 1
		enable = false
		filter = "ipv6_te"
		filter_ref_type_ = "ipv6_filter"
	}
	ipv6_permit {
		index = 1
		enable = true
		filter = "ipv6_te"
		filter_ref_type_ = "ipv6_filter"
	}
	ipv6_permit {
		index = 2
		enable = true
		filter = "OINKER"
		filter_ref_type_ = "ipv6_filter"
	}
}


# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_packet_broker" "packet_broker_test_script1" {
	enable = false
	ipv4_deny {
		index = 2
	}
	ipv4_permit {
		index = 2
		enable = true
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_deny {
		index = 2
		enable = true
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_permit {
		index = 2
		enable = true
		filter = ""
		filter_ref_type_ = ""
	}
}

resource "verity_packet_broker" "packet_broker_test_script2" {
	enable = false
	ipv4_deny {
		index = 1
		enable = false
	}
	ipv4_permit {
		index = 1
		enable = false
	}
	ipv6_deny {
		index = 1
		enable = false
	}
	ipv6_permit {
		index = 1
		enable = false
	}
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_pb_routing_acl" "pb_routing_acl_test_script1" {
	enable = false
	ipv4_deny {
		index = 1
		enable = false
	}
	ipv4_deny {
		index = 2
	}
	ipv4_permit {
		index = 1
		enable = false
	}
	ipv6_deny {
		index = 1
		enable = true
	}
	ipv6_permit {
		index = 1
		enable = true
	}
	ipv_protocol = "ipv4"
	next_hop_ips = "10.12.14.16"
}


resource "verity_pb_routing_acl" "pb_routing_acl_test_script2" {
	enable = false
	ipv4_deny {
		index = 1
		enable = false
	}
	ipv4_permit {
		index = 1
		enable = true
	}
	ipv6_deny {
		index = 1
		enable = true
	}
	ipv6_permit {
		index = 1
		enable = true
	}
	ipv_protocol = "ipv4"
	next_hop_ips = "20.20.20.20"
}
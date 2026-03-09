# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_gateway" "gateway_test_script1" {
	object_properties {
		group = "test"
	}
	advertisement_interval = 31
	bfd_detect_multiplier = 3
	bfd_multihop = true
	bfd_receive_interval = 302
	bfd_transmission_interval = 302
	connect_timer = 122
	default_originate = true
	dynamic_bgp_limits = 1
	enable = false
	enable_bfd = true
	fabric_interconnect = true
	hold_timer = 190
	keepalive_timer = 62
	local_as_no_prepend = false
	local_as_number = 12
	max_local_as_occurrences = 1
	neighbor_as_number = 13
	next_hop_self = true
	replace_as = true
	static_routes {
		index = 1
		enable = true
	}
	static_routes {
		index = 2
	}
}

resource "verity_gateway" "gateway_test_script2" {
	object_properties {
		group = ""
	}
	advertisement_interval = 33
	bfd_detect_multiplier = 5
	bfd_multihop = true
	bfd_receive_interval = 303
	bfd_transmission_interval = 303
	connect_timer = 124
	default_originate = false
	dynamic_bgp_limits = 1
	enable = false
	enable_bfd = true
	fabric_interconnect = true
	hold_timer = 170
	keepalive_timer = 130
	local_as_no_prepend = false
	local_as_number = 13
	max_local_as_occurrences = 1
	neighbor_as_number = 13
	next_hop_self = true
	replace_as = true
	static_routes {
		index = 2
		ad_value = 2
		enable = false
		ipv4_route_prefix = ""
		next_hop_ip_address = ""
	}
}
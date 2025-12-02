# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_gateway" "gateway_test_script1" {
	object_properties {
		group = "test"
	}
	advertisement_interval = 32
	bfd_receive_interval = 310
	bfd_transmission_interval = 310
	connect_timer = 122
	default_originate = true
	ebgp_multihop = 261
	enable = true
	enable_bfd = true
	hold_timer = 184
	static_routes {
		index = 2
		ad_value = 1
		enable = true
		ipv4_route_prefix = ""
		next_hop_ip_address = ""
	}
}
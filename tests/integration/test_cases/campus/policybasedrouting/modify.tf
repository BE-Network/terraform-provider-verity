# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_pb_routing" "pb_routing_test_script1" {
	enable = false
	policy {
		index = 1
		enable = false
	}
	policy {
		index = 2
	}
}


resource "verity_pb_routing" "pb_routing_test_script2" {
	enable = false
	policy {
		index = 3
		enable = true
		pb_routing_acl = "pbr_acl_ipv4_2"
		pb_routing_acl_ref_type_ = "pb_routing_acl"
	}
}
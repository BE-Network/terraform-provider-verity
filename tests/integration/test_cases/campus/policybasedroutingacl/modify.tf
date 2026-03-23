# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_pb_routing_acl" "pb_routing_acl_test_script1" {
	enable = true
}

resource "verity_pb_routing_acl" "pb_routing_acl_test_script2" {
	enable = true
}
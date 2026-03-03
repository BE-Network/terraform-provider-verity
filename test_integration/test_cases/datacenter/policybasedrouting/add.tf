# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_pb_routing" "pb_routing_test_script1" {
    name = "pb_routing_test_script1"
    depends_on = [verity_operation_stage.pb_routing_stage]
	enable = true
	policy {
		index = 1
		enable = true
		pb_routing_acl = "pbr_acl_ipv4_2"
		pb_routing_acl_ref_type_ = "pb_routing_acl"
	}
	policy {
		index = 2
		enable = true
		pb_routing_acl = "pbr_acl_ipv6_2"
		pb_routing_acl_ref_type_ = "pb_routing_acl"
	}
}
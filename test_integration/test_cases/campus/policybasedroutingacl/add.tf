# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_pb_routing_acl" "pb_routing_acl_test_script1" {
    name = "pb_routing_acl_test_script1"
    depends_on = [verity_operation_stage.pb_routing_acl_stage]
	enable = false
	ipv_protocol = "ipv4"
	next_hop_ips = ""
}

resource "verity_pb_routing_acl" "pb_routing_acl_test_script2" {
    name = "pb_routing_acl_test_script2"
    depends_on = [verity_operation_stage.pb_routing_acl_stage]
	enable = false
	ipv_protocol = "ipv4"
	next_hop_ips = ""
}

resource "verity_pb_routing" "pbr1" {
    name = "pbr1"
    depends_on = [verity_operation_stage.pb_routing_stage]
	enable = true
	policy {
		index = 1
		enable = true
		pb_routing_acl = "ipv4_1"
		pb_routing_acl_ref_type_ = "pb_routing_acl"
	}
	policy {
		index = 2
		enable = true
		pb_routing_acl = "pbr_acl_ipv6_1"
		pb_routing_acl_ref_type_ = "pb_routing_acl"
	}
}


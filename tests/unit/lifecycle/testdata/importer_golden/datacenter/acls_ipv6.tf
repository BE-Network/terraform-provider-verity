
resource "verity_acl_v6" "test1" {
    name = "test1"
    depends_on = [verity_operation_stage.acl_v6_stage]
	object_properties {
		notes = "test1"
	}
	bidirectional = false
	destination_ip = "2001:2001:2001:0000:0000:2001:2001:2001"
	destination_port_1 = null
	destination_port_2 = null
	destination_port_operator = ""
	enable = true
	protocol = "tcp"
	source_ip = "2001:2001:2001:0000:0000:2001:2001:2001"
	source_port_1 = 1
	source_port_2 = 5
	source_port_operator = "range"
}


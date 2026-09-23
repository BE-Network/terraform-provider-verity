
resource "verity_acl_v4" "hello-world-ipv4-filter" {
    name = "hello-world-ipv4-filter"
    depends_on = [verity_operation_stage.acl_v4_stage]
	object_properties {
		notes = ""
	}
	bidirectional = false
	destination_ip = ""
	destination_port_1 = null
	destination_port_2 = null
	destination_port_operator = ""
	enable = true
	protocol = "tcp"
	source_ip = ""
	source_port_1 = null
	source_port_2 = null
	source_port_operator = ""
}


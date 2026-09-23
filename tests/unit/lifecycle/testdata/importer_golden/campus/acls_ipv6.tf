
resource "verity_acl_v6" "ip_filter_test1" {
    name = "ip_filter_test1"
    depends_on = [verity_operation_stage.acl_v6_stage]
	object_properties {
		notes = ""
	}
	bidirectional = false
	destination_ip = ""
	destination_port_1 = null
	destination_port_2 = null
	destination_port_operator = ""
	enable = false
	protocol = ""
	source_ip = ""
	source_port_1 = null
	source_port_2 = null
	source_port_operator = ""
}


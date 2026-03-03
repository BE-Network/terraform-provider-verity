# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_acl_v6" "acl_v6_test_script1" {
    name = "acl_v6_test_script1"
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
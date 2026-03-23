# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_port_acl" "port_acl_test_script1" {
    name = "port_acl_test_script1"
    depends_on = [verity_operation_stage.port_acl_stage]
	enable = true
}

resource "verity_port_acl" "port_acl_test_script2" {
    name = "port_acl_test_script2"
    depends_on = [verity_operation_stage.port_acl_stage]
	enable = true
}
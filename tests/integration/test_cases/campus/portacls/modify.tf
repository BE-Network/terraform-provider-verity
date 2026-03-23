# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_port_acl" "port_acl_test_script1" {
	enable = false
}

resource "verity_port_acl" "port_acl_test_script2" {
	enable = false
}
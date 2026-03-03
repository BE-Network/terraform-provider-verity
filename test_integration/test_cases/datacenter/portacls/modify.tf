# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_port_acl" "port_acl_test_script1" {
	enable = true
	ipv4_deny {
		index = 2
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_acl_v6" "acl_v6_test_script1" {
	object_properties {
		notes = "test"
	}
	bidirectional = false
	enable = true
	source_port_1 = 5
	source_port_2 = 6
}

resource "verity_acl_v6" "acl_v6_test_script2" {
	object_properties {
		notes = ""
	}
	bidirectional = true
	enable = false
	source_port_1 = 7
	source_port_2 = 8
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_acl_v4" "acl_v4_test_script1" {
	object_properties {
		notes = "test1"
	}
	bidirectional = true
	enable = false
}

resource "verity_acl_v4" "acl_v4_test_script2" {
	object_properties {
		notes = ""
	}
	bidirectional = false
	enable = true
}
# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_acl_v6" "acl_v6_test_script1" {
	object_properties {
		notes = "test"
	}
	bidirectional = true
}
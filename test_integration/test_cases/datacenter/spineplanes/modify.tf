# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_spine_plane" "spine_plane_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
}
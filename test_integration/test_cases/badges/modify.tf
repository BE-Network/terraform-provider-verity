# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_badge" "badge_test_script1" {
	object_properties {
		notes = "test123"
	}
	enable = false
}
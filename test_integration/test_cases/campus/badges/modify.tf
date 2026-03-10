# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_badge" "badge_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	number = 5
}

resource "verity_badge" "badge_test_script2" {
	object_properties {
		notes = "test123"
	}
	enable = false
	number = 6
}
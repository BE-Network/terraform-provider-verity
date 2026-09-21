
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
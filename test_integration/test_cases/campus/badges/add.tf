# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_badge" "badge_test_script1" {
    name = "badge_test_script1"
    depends_on = [verity_operation_stage.badge_stage]
	object_properties {
		notes = ""
	}
	color = "yellow"
	enable = true
	number = 4
}

resource "verity_badge" "badge_test_script2" {
    name = "badge_test_script2"
    depends_on = [verity_operation_stage.badge_stage]
	object_properties {
		notes = "test"
	}
	color = "orange"
	enable = true
	number = 5
}
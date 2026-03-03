# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_badge" "badge_test_script1" {
    name = "badge_test_script1"
    depends_on = [verity_operation_stage.badge_stage]
	object_properties {
		notes = "test"
	}
	color = "blue"
	enable = true
	number = 2
}
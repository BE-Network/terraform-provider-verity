
resource "verity_badge" "badge_test1" {
    name = "badge_test1"
    depends_on = [verity_operation_stage.badge_stage]
	object_properties {
		notes = "test"
	}
	color = "red"
	enable = true
	number = 1
}



resource "verity_badge" "test" {
    name = "test"
    depends_on = [verity_operation_stage.badge_stage]
	object_properties {
		notes = ""
	}
	color = "blue"
	enable = true
	number = 1
}


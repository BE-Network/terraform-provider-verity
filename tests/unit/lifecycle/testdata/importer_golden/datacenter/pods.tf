
resource "verity_pod" "Pod_1" {
    name = "Pod 1"
    depends_on = [verity_operation_stage.pod_stage]
	object_properties {
		notes = ""
	}
	enable = true
	expected_spine_count = 2
}


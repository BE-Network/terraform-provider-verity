
resource "verity_pod" "pod_test_script1" {
	object_properties {
		notes = "test"
	}
	enable = false
	expected_spine_count = 2
}

resource "verity_pod" "pod_test_script2" {
	object_properties {
		notes = ""
	}
	enable = true
	expected_spine_count = 1
}
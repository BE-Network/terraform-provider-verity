# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_pod" "pod_test_script1" {
    name = "pod_test_script1"
    depends_on = [verity_operation_stage.pod_stage]
	object_properties {
		notes = ""
	}
	enable = true
	expected_spine_count = 1
}

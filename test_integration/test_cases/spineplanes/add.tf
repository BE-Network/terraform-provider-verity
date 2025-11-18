# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_spine_plane" "spine_plane_test_script1" {
    name = "spine_plane_test_script1"
    depends_on = [verity_operation_stage.spine_plane_stage]
	object_properties {
		notes = ""
	}
	enable = true
}
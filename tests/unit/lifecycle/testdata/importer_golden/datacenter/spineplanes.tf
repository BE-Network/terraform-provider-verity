
resource "verity_spine_plane" "spine_plane_test1" {
    name = "spine_plane_test1"
    depends_on = [verity_operation_stage.spine_plane_stage]
	object_properties {
		notes = ""
	}
	enable = true
}


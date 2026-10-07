
resource "verity_fabric_collection" "test" {
    name = "test"
    depends_on = [verity_operation_stage.fabric_collection_stage]
	enable = true
}



resource "verity_sflow_collector" "Example_Collector" {
    name = "Example Collector"
    depends_on = [verity_operation_stage.sflow_collector_stage]
	enable = true
	ip = "187.22.22.22"
	port = 6343
}


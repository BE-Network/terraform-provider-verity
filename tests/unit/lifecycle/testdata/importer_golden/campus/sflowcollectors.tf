
resource "verity_sflow_collector" "sflow_collector_test1" {
    name = "sflow_collector_test1"
    depends_on = [verity_operation_stage.sflow_collector_stage]
	enable = false
	ip = ""
	port = 6343
}


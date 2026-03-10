# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing


resource "verity_sflow_collector" "sflow_collector_test_script1" {
    name = "sflow_collector_test_script1"
    depends_on = [verity_operation_stage.sflow_collector_stage]
	enable = true
	ip = ""
	port = 6343
}
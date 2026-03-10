# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_sflow_collector" "sflow_collector_test_script1" {
	enable = false
	port = 6344
}
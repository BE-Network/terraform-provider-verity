# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_ipv4_list" "ipv4_test_script1" {
    name = "ipv4_test_script1"
    depends_on = [verity_operation_stage.ipv4_list_stage]
	enable = true
	ipv4_list = "7.7.7.7, 8.8.8.8"
}

resource "verity_ipv4_list" "ipv4_test_script2" {
    name = "ipv4_test_script2"
    depends_on = [verity_operation_stage.ipv4_list_stage]
	enable = true
	ipv4_list = "9.9.9.9, 10.10.10.10"
}
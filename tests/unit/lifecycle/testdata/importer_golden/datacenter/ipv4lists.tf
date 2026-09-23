
resource "verity_ipv4_list" "ipv4_test_list" {
    name = "ipv4_test_list"
    depends_on = [verity_operation_stage.ipv4_list_stage]
	enable = true
	ipv4_list = "1.1.1.1, 2.2.2.2"
}


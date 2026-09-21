
resource "verity_ipv6_list" "ipv6_list_test_script1" {
    name = "ipv6_list_test_script1"
    depends_on = [verity_operation_stage.ipv6_list_stage]
	enable = true
	ipv6_list = ""
}
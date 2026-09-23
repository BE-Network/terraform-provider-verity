
resource "verity_ipv6_list" "ipv6_list_test" {
    name = "ipv6_list_test"
    depends_on = [verity_operation_stage.ipv6_list_stage]
	enable = true
	ipv6_list = "2004:2004:2004:2004:2004:2004:2004:2004, 2005:2005:2005:2005:2005:2005:2005:2005, 2006:2006:2006:2006:2006:2006:2006:2006"
}


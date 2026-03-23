# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_ipv6_list" "ipv6_list_test_script1" {
	enable = false
}

resource "verity_ipv6_list" "ipv6_list_test_script2" {
	enable = true
}
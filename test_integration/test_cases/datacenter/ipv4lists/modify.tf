# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_ipv4_list" "ipv4_test_script1" {
	enable = false
}
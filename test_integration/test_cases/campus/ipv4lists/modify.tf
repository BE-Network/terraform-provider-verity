# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_ipv4_list" "ipv4_test_script1" {
	enable = false
	ipv4_list = "9.9.9.9, 10.10.10.10"
}

resource "verity_ipv4_list" "ipv4_test_script2" {
	enable = false
	ipv4_list = "7.7.7.7, 8.8.8.8"
}
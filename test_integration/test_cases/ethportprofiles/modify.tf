# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_eth_port_profile" "eth_port_profile_test_script1" {
	object_properties {
		group = "test"
	}
	enable = true
}
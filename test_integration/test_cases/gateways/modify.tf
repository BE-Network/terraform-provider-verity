# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_gateway" "gateway_test_script1" {
	bfd_transmission_interval = 311
	keepalive_timer = 66
}
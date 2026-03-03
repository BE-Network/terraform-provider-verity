# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service_port_profile" "service_port_profile_test_script1" {
	object_properties {
		group = "test"
		on_summary = false
	}
	enable = false
	services {
		index = 1
		row_num_enable = false
	}
	tls_limit_in = 1010
}
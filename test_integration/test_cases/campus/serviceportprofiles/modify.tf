# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_service_port_profile" "service_port_profile_test_script1" {
	object_properties {
		group = "test"
		on_summary = false
		port_monitoring = "test"
	}
	enable = false
	services {
		index = 1
		row_num_enable = false
	}
	services {
		index = 4
		row_num_enable = false
	}
	tls_limit_in = 1001
	trusted_port = true
}

resource "verity_service_port_profile" "service_port_profile_test_script2" {
	object_properties {
		group = "test1"
		on_summary = false
		port_monitoring = "test1"
	}
	enable = false
	services {
		index = 1
		row_num_enable = false
	}
	services {
		index = 4
		row_num_enable = false
	}
	tls_limit_in = 1002
	trusted_port = true
}
# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_service_port_profile" "service_port_profile_test_script1" {
    name = "service_port_profile_test_script1"
    depends_on = [verity_operation_stage.service_port_profile_stage]
	object_properties {
		group = ""
		on_summary = true
		port_monitoring = ""
	}
	enable = true
	ip_mask = ""
	port_type = "up"
	services {
		index = 1
		row_num_enable = true
		row_num_external_vlan = null
		row_num_limit_in = 1000
		row_num_limit_out = 1000
		row_num_service = "Management"
		row_num_service_ref_type_ = "service"
	}
	tls_limit_in = 1000
	tls_service = ""
	tls_service_ref_type_ = ""
	trusted_port = false
}
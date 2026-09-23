
resource "verity_service_port_profile" "test1" {
    name = "test1"
    depends_on = [verity_operation_stage.service_port_profile_stage]
	object_properties {
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
		row_num_limit_in = null
		row_num_limit_out = 1000
		row_num_service = "v100"
		row_num_service_ref_type_ = "service"
	}
	services {
		index = 2
		row_num_enable = true
		row_num_external_vlan = null
		row_num_limit_in = null
		row_num_limit_out = 1000
		row_num_service = "srv3"
		row_num_service_ref_type_ = "service"
	}
	services {
		index = 3
		row_num_enable = true
		row_num_external_vlan = null
		row_num_limit_in = null
		row_num_limit_out = 1000
		row_num_service = "srv2"
		row_num_service_ref_type_ = "service"
	}
	services {
		index = 4
		row_num_enable = true
		row_num_external_vlan = null
		row_num_limit_in = null
		row_num_limit_out = 1000
		row_num_service = "srv1"
		row_num_service_ref_type_ = "service"
	}
	tls_limit_in = 1000
	tls_service = ""
	tls_service_ref_type_ = ""
	trusted_port = false
}


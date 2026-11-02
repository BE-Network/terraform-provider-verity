
resource "verity_static_ip" "test" {
    name = "test"
    depends_on = [verity_operation_stage.static_ip_stage]
	object_properties {
		fabric = ""
		fabric_ref_type_ = ""
	}
	allowed_on = "Static Port"
	enable = true
	ip = "10.0.0.10"
	lag = ""
	lag_ref_type_ = ""
	mac = "00:11:22:33:44:55"
	port = ""
	service = "test"
	service_ref_type_ = "service"
	switch = ""
	switch_ref_type_ = ""
}


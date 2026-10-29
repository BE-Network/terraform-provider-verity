
resource "verity_monitoring_acl" "test" {
    name = "test"
    depends_on = [verity_operation_stage.monitoring_acl_stage]
	enable = true
	services {
		index = 1
		enable = true
		service = "test"
		service_ref_type_ = "service"
	}
}


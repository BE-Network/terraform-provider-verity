
resource "verity_port_acl" "ex" {
    name = "ex"
    depends_on = [verity_operation_stage.port_acl_stage]
	enable = true
}


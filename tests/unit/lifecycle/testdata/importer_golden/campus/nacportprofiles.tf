
resource "verity_nac_port_profile" "test" {
    name = "test"
    depends_on = [verity_operation_stage.nac_port_profile_stage]
	object_properties {
		port_monitoring = "high"
	}
	enable = true
	eth_ports {
		index = 1
		eth_port_profile_num_enable = true
		eth_port_profile_num_eth_port = "test"
		eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
	}
}


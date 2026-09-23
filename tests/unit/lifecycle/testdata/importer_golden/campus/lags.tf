
resource "verity_lag" "x" {
    name = "x"
    depends_on = [verity_operation_stage.lag_stage]
	color = "anakiwa"
	enable = false
	eth_port_profile = ""
	eth_port_profile_ref_type_ = ""
	fallback = false
	fast_rate = false
	is_peer_link = false
	lacp = true
	peer_link_vlan = null
	uplink = false
}


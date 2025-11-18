# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_lag" "lag_test_script1" {
    name = "lag_test_script1"
    depends_on = [verity_operation_stage.lag_stage]
	color = "anakiwa"
	enable = true
	eth_port_profile = ""
	eth_port_profile_ref_type_ = ""
	fallback = false
	fast_rate = false
	is_peer_link = false
	lacp = true
	peer_link_vlan = null
	uplink = false
}
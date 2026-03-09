# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_lag" "lag_test_script1" {
	color = "lavender"
	enable = false
	fallback = true
	fast_rate = true
	is_peer_link = true
	lacp = false
	peer_link_vlan = 101
	uplink = false
}

resource "verity_lag" "lag_test_script2" {
	color = "chardonnay"
	enable = false
	fallback = false
	fast_rate = false
	is_peer_link = false
	lacp = true
	peer_link_vlan = null
	uplink = true
}
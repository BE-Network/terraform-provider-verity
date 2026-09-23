
resource "verity_gateway" "gateway_test1" {
    name = "gateway_test1"
    depends_on = [verity_operation_stage.gateway_stage]
	advertisement_interval = 30
	anycast_ip_mask = ""
	bfd_detect_multiplier = 3
	bfd_multihop = false
	bfd_receive_interval = 300
	bfd_transmission_interval = 300
	bgp_instance_as_number = null
	connect_timer = 120
	default_originate = false
	dynamic_bgp_limits = 0
	dynamic_bgp_subnet = ""
	ebgp_multihop = 255
	egress_vlan = null
	enable = true
	enable_bfd = false
	export_route_map = ""
	export_route_map_ref_type_ = ""
	fabric_interconnect = false
	gateway_mode = "Static BGP"
	helper_hop_ip_address = ""
	hold_timer = 182
	import_route_map = ""
	import_route_map_ref_type_ = ""
	keepalive_timer = 60
	local_as_no_prepend = true
	local_as_number = 10
	max_local_as_occurrences = 0
	md5_password = ""
	md5_password_encrypted = ""
	neighbor_as_number = 12
	neighbor_ip_address = "8.8.8.8"
	next_hop_self = false
	replace_as = false
	source_ip_address = ""
	static_routes {
		index = 1
		ad_value = 1
		enable = false
		ipv4_route_prefix = ""
		next_hop_ip_address = ""
	}
	switch_encrypted_md5_password = false
	tenant = "Visualizator"
	tenant_ref_type_ = "tenant"
}


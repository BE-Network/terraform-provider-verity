# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_device_controller" "device_controller_test_script1" {
    name = "device_controller_test_script1"
    depends_on = [verity_operation_stage.device_controller_stage]
	authentication_protocol = ""
	cli_access_mode = "SSH"
	comm_type = ""
	communication_mode = "generic_advanced_snmp"
	controller_ip_and_mask = ""
	enable = false
	enable_password = ""
	enable_password_encrypted = ""
	gateway = ""
	ip_source = "dhcp"
	lldp_search_string = ""
	located_by = "LLDP"
	managed_on_native_vlan = true
	passphrase = ""
	passphrase_encrypted = ""
	password = ""
	password_encrypted = ""
	power_state = "on"
	private_password = ""
	private_password_encrypted = ""
	private_protocol = "DES"
	sdlc = ""
	security_type = ""
	snmp_community_string = "public"
	snmpv3_username = ""
	ssh_key_or_password = ""
	ssh_key_or_password_encrypted = ""
	switch_gateway = "10.10.10.1"
	switch_ip_and_mask = "10.10.10.9/24"
	switchpoint = "t1"
	switchpoint_ref_type_ = "switchpoint"
	uplink_port = ""
	username = ""
	ztp_identification = ""
}
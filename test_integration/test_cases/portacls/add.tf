# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_port_acl" "port_acl_test_script1" {
    name = "port_acl_test_script1"
    depends_on = [verity_operation_stage.port_acl_stage]
	enable = false
	ipv4_deny {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv4_permit {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_deny {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
	ipv6_permit {
		index = 1
		enable = false
		filter = ""
		filter_ref_type_ = ""
	}
}
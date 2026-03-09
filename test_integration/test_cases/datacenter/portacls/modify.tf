# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_port_acl" "port_acl_test_script1" {
	enable = true
	ipv4_deny {
		index = 1
		enable = true
	}
	ipv4_permit {
		index = 1
		enable = true
	}
	ipv6_deny {
		index = 1
		enable = true
	}
	ipv6_permit {
		index = 1
		enable = true
	}
}

resource "verity_port_acl" "port_acl_test_script2" {
	enable = true
	ipv4_deny {
		index = 1
		enable = true
	}
	ipv4_deny {
		index = 2
		enable = true
	}
	ipv4_permit {
		index = 1
		enable = true
	}
	ipv6_deny {
		index = 1
		enable = true
	}
	ipv6_permit {
		index = 1
		enable = true
	}
	ipv6_permit {
		index = 2
	}
}
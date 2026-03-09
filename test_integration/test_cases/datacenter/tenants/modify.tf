# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_tenant" "tenant_test_script1" {
	object_properties {
		group = "test"
	}
	default_originate = true
	enable = false
	route_tenants {
		index = 1
		enable = true
	}
	route_tenants {
		index = 2
	}
}

resource "verity_tenant" "tenant_test_script2" {
	object_properties {
		group = ""
	}
	default_originate = true
	enable = false
	route_tenants {
		index = 1
		enable = true
	}
	route_tenants {
		index = 2
		enable = true
	}
}
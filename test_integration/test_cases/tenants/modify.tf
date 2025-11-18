# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_tenant" "tenant_test_script1" {
	object_properties {
		group = "test"
	}
	default_originate = true
	enable = false
	route_tenants {
		index = 2
		enable = false
		tenant = ""
	}
}
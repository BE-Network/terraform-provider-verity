
resource "verity_ipv4_prefix_list" "test1" {
    name = "test1"
    depends_on = [verity_operation_stage.ipv4_prefix_list_stage]
	object_properties {
		notes = ""
	}
	enable = true
	lists {
		index = 1
		enable = false
		greater_than_equal_value = null
		ipv4_prefix = ""
		less_than_equal_value = null
		permit_deny = "permit"
	}
	lists {
		index = 2
		enable = false
		greater_than_equal_value = null
		ipv4_prefix = ""
		less_than_equal_value = null
		permit_deny = "permit"
	}
	lists {
		index = 3
		enable = false
		greater_than_equal_value = null
		ipv4_prefix = ""
		less_than_equal_value = null
		permit_deny = "permit"
	}
	lists {
		index = 4
		enable = false
		greater_than_equal_value = null
		ipv4_prefix = ""
		less_than_equal_value = null
		permit_deny = "permit"
	}
}


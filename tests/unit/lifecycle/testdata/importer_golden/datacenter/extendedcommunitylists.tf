
resource "verity_extended_community_list" "Extended_Community_List" {
    name = "Extended Community List"
    depends_on = [verity_operation_stage.extended_community_list_stage]
	object_properties {
		notes = ""
	}
	any_all = "any"
	enable = true
	lists {
		index = 1
		enable = false
		mode = "route"
		route_target_expanded_expression = ""
	}
	permit_deny = "permit"
	standard_expanded = "standard"
}



resource "verity_community_list" "Changes_BGP_CommunityLists" {
    name = "Changes_BGP_CommunityLists"
    depends_on = [verity_operation_stage.community_list_stage]
	object_properties {
		notes = ""
	}
	any_all = "any"
	enable = true
	lists {
		index = 1
		community_string_expanded_expression = ""
		enable = false
		mode = "community"
	}
	lists {
		index = 2
		community_string_expanded_expression = ""
		enable = false
		mode = "community"
	}
	permit_deny = "permit"
	standard_expanded = "standard"
}


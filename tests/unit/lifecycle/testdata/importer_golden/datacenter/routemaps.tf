
resource "verity_route_map" "Changes_BGP-RouteMaps" {
    name = "Changes_BGP-RouteMaps"
    depends_on = [verity_operation_stage.route_map_stage]
	object_properties {
		notes = ""
	}
	enable = true
	route_map_clauses {
		index = 1
		enable = true
		route_map_clause = "Changes_BGP_RMClause"
		route_map_clause_ref_type_ = "route_map_clause"
	}
}


resource "verity_route_map" "Route_Map" {
    name = "Route Map"
    depends_on = [verity_operation_stage.route_map_stage]
	object_properties {
		notes = ""
	}
	enable = true
	route_map_clauses {
		index = 1
		enable = true
		route_map_clause = "ipv6_clause"
		route_map_clause_ref_type_ = "route_map_clause"
	}
}


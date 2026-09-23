
# These resources establish ordering for bulk operations in DATACENTER mode
resource "verity_operation_stage" "sfp_breakout_stage" {
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "community_list_stage" {
  depends_on = [verity_operation_stage.sfp_breakout_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "as_path_access_list_stage" {
  depends_on = [verity_operation_stage.community_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ipv6_prefix_list_stage" {
  depends_on = [verity_operation_stage.as_path_access_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ipv4_prefix_list_stage" {
  depends_on = [verity_operation_stage.ipv6_prefix_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "extended_community_list_stage" {
  depends_on = [verity_operation_stage.ipv4_prefix_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "acl_v6_stage" {
  depends_on = [verity_operation_stage.extended_community_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "acl_v4_stage" {
  depends_on = [verity_operation_stage.acl_v6_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "route_map_clause_stage" {
  depends_on = [verity_operation_stage.acl_v4_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "pb_routing_acl_stage" {
  depends_on = [verity_operation_stage.route_map_clause_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "route_map_stage" {
  depends_on = [verity_operation_stage.pb_routing_acl_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "pb_routing_stage" {
  depends_on = [verity_operation_stage.route_map_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "tenant_stage" {
  depends_on = [verity_operation_stage.pb_routing_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "service_stage" {
  depends_on = [verity_operation_stage.tenant_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "fabric_stage" {
  depends_on = [verity_operation_stage.service_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "tacacs_profile_stage" {
  depends_on = [verity_operation_stage.fabric_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ldap_profile_stage" {
  depends_on = [verity_operation_stage.tacacs_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "port_acl_stage" {
  depends_on = [verity_operation_stage.ldap_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ipv6_list_stage" {
  depends_on = [verity_operation_stage.port_acl_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ipv4_list_stage" {
  depends_on = [verity_operation_stage.ipv6_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "pod_stage" {
  depends_on = [verity_operation_stage.ipv4_list_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "packet_queue_stage" {
  depends_on = [verity_operation_stage.pod_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "device_aaa_profile_stage" {
  depends_on = [verity_operation_stage.packet_queue_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_profile_stage" {
  depends_on = [verity_operation_stage.device_aaa_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "packet_broker_stage" {
  depends_on = [verity_operation_stage.eth_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "sflow_collector_stage" {
  depends_on = [verity_operation_stage.packet_broker_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "gateway_stage" {
  depends_on = [verity_operation_stage.sflow_collector_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "su_stage" {
  depends_on = [verity_operation_stage.gateway_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "diagnostics_port_profile_stage" {
  depends_on = [verity_operation_stage.su_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "device_settings_stage" {
  depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "lag_stage" {
  depends_on = [verity_operation_stage.device_settings_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "diagnostics_profile_stage" {
  depends_on = [verity_operation_stage.lag_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "gateway_profile_stage" {
  depends_on = [verity_operation_stage.diagnostics_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_settings_stage" {
  depends_on = [verity_operation_stage.gateway_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "badge_stage" {
  depends_on = [verity_operation_stage.eth_port_settings_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "plane_stage" {
  depends_on = [verity_operation_stage.badge_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "spine_plane_stage" {
  depends_on = [verity_operation_stage.plane_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "rack_stage" {
  depends_on = [verity_operation_stage.spine_plane_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "bundle_stage" {
  depends_on = [verity_operation_stage.rack_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "ssp_group_stage" {
  depends_on = [verity_operation_stage.bundle_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "grouping_rule_stage" {
  depends_on = [verity_operation_stage.ssp_group_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "switchpoint_stage" {
  depends_on = [verity_operation_stage.grouping_rule_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "threshold_stage" {
  depends_on = [verity_operation_stage.switchpoint_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "threshold_group_stage" {
  depends_on = [verity_operation_stage.threshold_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "pair_stage" {
  depends_on = [verity_operation_stage.threshold_group_stage]
  lifecycle {
    create_before_destroy = true
  }
}


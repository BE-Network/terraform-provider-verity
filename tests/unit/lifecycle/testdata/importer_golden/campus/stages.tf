
# These resources establish ordering for bulk operations in CAMPUS mode
resource "verity_operation_stage" "sfp_breakout_stage" {
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "acl_v6_stage" {
  depends_on = [verity_operation_stage.sfp_breakout_stage]
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

resource "verity_operation_stage" "mac_filter_stage" {
  depends_on = [verity_operation_stage.acl_v4_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "service_stage" {
  depends_on = [verity_operation_stage.mac_filter_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "port_acl_stage" {
  depends_on = [verity_operation_stage.service_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "tacacs_profile_stage" {
  depends_on = [verity_operation_stage.port_acl_stage]
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

resource "verity_operation_stage" "sflow_collector_stage" {
  depends_on = [verity_operation_stage.ldap_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_profile_stage" {
  depends_on = [verity_operation_stage.sflow_collector_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "packet_queue_stage" {
  depends_on = [verity_operation_stage.eth_port_profile_stage]
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

resource "verity_operation_stage" "fabric_stage" {
  depends_on = [verity_operation_stage.device_aaa_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "service_port_profile_stage" {
  depends_on = [verity_operation_stage.fabric_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "diagnostics_profile_stage" {
  depends_on = [verity_operation_stage.service_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "authenticated_eth_port_stage" {
  depends_on = [verity_operation_stage.diagnostics_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "device_settings_stage" {
  depends_on = [verity_operation_stage.authenticated_eth_port_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "voice_port_profile_stage" {
  depends_on = [verity_operation_stage.device_settings_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "lag_stage" {
  depends_on = [verity_operation_stage.voice_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "device_voice_setting_stage" {
  depends_on = [verity_operation_stage.lag_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_settings_stage" {
  depends_on = [verity_operation_stage.device_voice_setting_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "diagnostics_port_profile_stage" {
  depends_on = [verity_operation_stage.eth_port_settings_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "bundle_stage" {
  depends_on = [verity_operation_stage.diagnostics_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "badge_stage" {
  depends_on = [verity_operation_stage.bundle_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "grouping_rule_stage" {
  depends_on = [verity_operation_stage.badge_stage]
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


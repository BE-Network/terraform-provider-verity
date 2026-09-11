// Mode field mappings.
//
// The bulk of this data is now derived from the reviewed spec registry and lives
// in generated_mode_metadata.go, produced by `specgen metadata`. What remains
// here is the hand-maintained residue for endpoints the registry does not
// represent yet, plus the lookup helpers.

package utils

type FieldMode string

const (
	FieldModeBoth       FieldMode = "both"
	FieldModeDatacenter FieldMode = "datacenter"
	FieldModeCampus     FieldMode = "campus"
)

// FieldAppliesToMode checks if a field applies to the given mode.
// Returns true if the field should be populated for the given mode.
// If the field is not found in ModeFields, it defaults to true (applies to both modes).
func FieldAppliesToMode(resourceType, fieldName, mode string) bool {
	resourceFields, ok := ModeFields[resourceType]
	if !ok {
		// Resource not found in mode map, assume field applies to all modes
		return true
	}

	fieldMode, ok := resourceFields[fieldName]
	if !ok {
		// Field not found in mode map, assume it applies to all modes
		return true
	}

	switch fieldMode {
	case FieldModeBoth:
		return true
	case FieldModeDatacenter:
		return mode == "datacenter"
	case FieldModeCampus:
		return mode == "campus"
	default:
		return true
	}
}

// pendingModeFields holds the endpoints the reviewed registry does not represent
// yet; see status.md. Entries move out as the registry grows, and ModeFields is
// the union of this table and the generated one.
var pendingModeFields = map[string]map[string]FieldMode{
	"devicesettings": {
		"cli_commands":                                FieldModeBoth,
		"commit_to_flash_interval":                    FieldModeBoth,
		"cut_through_switching":                       FieldModeBoth,
		"device_aaa_profile":                          FieldModeBoth,
		"device_aaa_profile_ref_type_":                FieldModeBoth,
		"disable_tcp_udp_learned_packet_acceleration": FieldModeBoth,
		"dns_servers":                                 FieldModeBoth,
		"dns_servers.enabled":                         FieldModeBoth,
		"dns_servers.index":                           FieldModeBoth,
		"dns_servers.server":                          FieldModeBoth,
		"enable":                                      FieldModeBoth,
		"external_battery_power_available":            FieldModeBoth,
		"external_power_available":                    FieldModeBoth,
		"login_banner":                                FieldModeBoth,
		"mode":                                        FieldModeBoth,
		"name":                                        FieldModeBoth,
		"ntp_servers":                                 FieldModeBoth,
		"ntp_servers.enabled":                         FieldModeBoth,
		"ntp_servers.index":                           FieldModeBoth,
		"ntp_servers.server":                          FieldModeBoth,
		"ntp_vrf":                                     FieldModeBoth,
		"ntp_vrf_tenant":                              FieldModeBoth,
		"ntp_vrf_tenant_ref_type_":                    FieldModeBoth,
		"object_properties":                           FieldModeBoth,
		"packet_queue":                                FieldModeBoth,
		"packet_queue_ref_type_":                      FieldModeBoth,
		"rocev2":                                      FieldModeBoth,
		"security_audit_interval":                     FieldModeBoth,
		"syslog_servers":                              FieldModeBoth,
		"syslog_servers.enabled":                      FieldModeBoth,
		"syslog_servers.index":                        FieldModeBoth,
		"syslog_servers.port":                         FieldModeBoth,
		"syslog_servers.scheme":                       FieldModeBoth,
		"syslog_servers.server":                       FieldModeBoth,
		"usage_threshold":                             FieldModeBoth,
		"hold_timer":                                  FieldModeCampus,
		"mac_aging_timer_override":                    FieldModeCampus,
		"spanning_tree_priority":                      FieldModeCampus,
	},
	"fabrics": {
		"aggressive_reporting":                            FieldModeBoth,
		"allow_all_underlay_connections":                  FieldModeBoth,
		"base_bgp_as_number":                              FieldModeBoth,
		"controller_gateway":                              FieldModeBoth,
		"controller_ip_base":                              FieldModeBoth,
		"domain_for_fabric":                               FieldModeBoth,
		"domain_for_fabric_ref_type_":                     FieldModeBoth,
		"dscp_to_p_bit_map":                               FieldModeBoth,
		"duplicate_address_detection_max_number_of_moves": FieldModeBoth,
		"duplicate_address_detection_time":                FieldModeBoth,
		"enable":                                          FieldModeBoth,
		"enable_dscp":                                     FieldModeBoth,
		"evpn_mac_holdtime":                               FieldModeBoth,
		"evpn_multihoming_startup_delay":                  FieldModeBoth,
		"fabric_type":                                     FieldModeBoth,
		"force_spanning_tree_on_fabric_ports":             FieldModeBoth,
		"gpu_architecture":                                FieldModeBoth,
		"hgx_password":                                    FieldModeBoth,
		"hgx_password_encrypted":                          FieldModeBoth,
		"hgx_username":                                    FieldModeBoth,
		"link_state_timeout_value":                        FieldModeBoth,
		"max_pods":                                        FieldModeBoth,
		"max_sus":                                         FieldModeBoth,
		"max_switches":                                    FieldModeBoth,
		"multi_tenant":                                    FieldModeBoth,
		"name":                                            FieldModeBoth,
		"object_properties":                               FieldModeBoth,
		"object_properties.system_graphs":                 FieldModeBoth,
		"object_properties.system_graphs.index":           FieldModeBoth,
		"paired_ip_subnet":                                FieldModeBoth,
		"pause_validation_alarms":                         FieldModeBoth,
		"plane_count":                                     FieldModeBoth,
		"port_admin_polling_interval":                     FieldModeBoth,
		"port_status_polling_interval":                    FieldModeBoth,
		"read_only_mode":                                  FieldModeBoth,
		"region_name":                                     FieldModeBoth,
		"revision":                                        FieldModeBoth,
		"route_aggregation":                               FieldModeBoth,
		"route_aggregators":                               FieldModeBoth,
		"route_aggregators.index":                         FieldModeBoth,
		"route_aggregators.route_aggregation_num_enable":      FieldModeBoth,
		"route_aggregators.route_aggregation_num_ip_and_mask": FieldModeBoth,
		"router_id_base_prefix":                               FieldModeBoth,
		"server_management":                                   FieldModeBoth,
		"service_for_fabric":                                  FieldModeBoth,
		"service_for_fabric_ref_type_":                        FieldModeBoth,
		"set_leaf_router_id_on_bgp":                           FieldModeBoth,
		"spanning_tree_type":                                  FieldModeBoth,
		"starting_octet":                                      FieldModeBoth,
		"su_size":                                             FieldModeBoth,
		"su_support":                                          FieldModeBoth,
		"switch_gateway":                                      FieldModeBoth,
		"switch_ip_base":                                      FieldModeBoth,
		"switch_password":                                     FieldModeBoth,
		"switch_password_encrypted":                           FieldModeBoth,
		"switch_username":                                     FieldModeBoth,
		"vtep_id_base_prefix":                                 FieldModeBoth,
		"anycast_mac_address":                                 FieldModeDatacenter,
		"anycast_mac_address_auto_assigned_":                  FieldModeDatacenter,
		"bgp_hold_down_timer":                                 FieldModeDatacenter,
		"bgp_keepalive_timer":                                 FieldModeDatacenter,
		"leaf_bgp_advertisement_interval":                     FieldModeDatacenter,
		"leaf_bgp_connect_timer":                              FieldModeDatacenter,
		"leaf_bgp_hold_down_timer":                            FieldModeDatacenter,
		"leaf_bgp_keep_alive_timer":                           FieldModeDatacenter,
		"mac_address_aging_time":                              FieldModeDatacenter,
		"mlag_delay_restore_timer":                            FieldModeDatacenter,
		"spine_as_number":                                     FieldModeDatacenter,
		"spine_bgp_advertisement_interval":                    FieldModeDatacenter,
		"spine_bgp_connect_timer":                             FieldModeDatacenter,
		"enable_dhcp_snooping":                                FieldModeCampus,
		"ip_source_guard":                                     FieldModeCampus,
	},
	"gateways": {
		"advertisement_interval":            FieldModeDatacenter,
		"allowas_in_origin":                 FieldModeDatacenter,
		"anycast_ip_mask":                   FieldModeDatacenter,
		"bfd_detect_multiplier":             FieldModeDatacenter,
		"bfd_multihop":                      FieldModeDatacenter,
		"bfd_receive_interval":              FieldModeDatacenter,
		"bfd_transmission_interval":         FieldModeDatacenter,
		"bgp_instance_as_number":            FieldModeDatacenter,
		"connect_timer":                     FieldModeDatacenter,
		"default_originate":                 FieldModeDatacenter,
		"dynamic_bgp_limits":                FieldModeDatacenter,
		"dynamic_bgp_subnet":                FieldModeDatacenter,
		"ebgp_multihop":                     FieldModeDatacenter,
		"egress_vlan":                       FieldModeDatacenter,
		"enable":                            FieldModeDatacenter,
		"enable_bfd":                        FieldModeDatacenter,
		"export_route_map":                  FieldModeDatacenter,
		"export_route_map_ref_type_":        FieldModeDatacenter,
		"fabric":                            FieldModeDatacenter,
		"fabric_interconnect":               FieldModeDatacenter,
		"fabric_ref_type_":                  FieldModeDatacenter,
		"gateway_mode":                      FieldModeDatacenter,
		"helper_hop_ip_address":             FieldModeDatacenter,
		"hold_timer":                        FieldModeDatacenter,
		"import_route_map":                  FieldModeDatacenter,
		"import_route_map_ref_type_":        FieldModeDatacenter,
		"keepalive_timer":                   FieldModeDatacenter,
		"local_as_no_prepend":               FieldModeDatacenter,
		"local_as_number":                   FieldModeDatacenter,
		"max_local_as_occurrences":          FieldModeDatacenter,
		"md5_password":                      FieldModeDatacenter,
		"md5_password_encrypted":            FieldModeDatacenter,
		"name":                              FieldModeDatacenter,
		"neighbor_as_number":                FieldModeDatacenter,
		"neighbor_ip_address":               FieldModeDatacenter,
		"next_hop_self":                     FieldModeDatacenter,
		"remove_private_as":                 FieldModeDatacenter,
		"replace_as":                        FieldModeDatacenter,
		"source_ip_address":                 FieldModeDatacenter,
		"static_routes":                     FieldModeDatacenter,
		"static_routes.ad_value":            FieldModeDatacenter,
		"static_routes.enable":              FieldModeDatacenter,
		"static_routes.index":               FieldModeDatacenter,
		"static_routes.ipv4_route_prefix":   FieldModeDatacenter,
		"static_routes.next_hop_ip_address": FieldModeDatacenter,
		"switch_encrypted_md5_password":     FieldModeDatacenter,
		"tenant":                            FieldModeDatacenter,
		"tenant_ref_type_":                  FieldModeDatacenter,
		"type":                              FieldModeDatacenter,
	},
	"sfpbreakouts": {
		"breakout":             FieldModeBoth,
		"breakout.breakout":    FieldModeBoth,
		"breakout.enable":      FieldModeBoth,
		"breakout.index":       FieldModeBoth,
		"breakout.part_number": FieldModeBoth,
		"breakout.vendor":      FieldModeBoth,
		"name":                 FieldModeBoth,
		"object_properties":    FieldModeBoth,
	},
}

// ModeFields maps API field paths to the modes they apply to. It is derived from
// the reviewed spec registry, with the pending endpoints above merged in.
var ModeFields = mergeModeFields()

func mergeModeFields() map[string]map[string]FieldMode {
	merged := make(map[string]map[string]FieldMode, len(generatedModeFields)+len(pendingModeFields))
	for resource, fields := range generatedModeFields {
		merged[resource] = fields
	}
	for resource, fields := range pendingModeFields {
		merged[resource] = fields
	}
	return merged
}

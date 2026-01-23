package utils

// ResourceJSONKeys maps internal resource type names to their JSON response keys
// This is the single source of truth for both importer and bulkops modules
var ResourceJSONKeys = map[string]string{
	"tenant":                   "tenant",
	"gateway":                  "gateway",
	"gateway_profile":          "gateway_profile",
	"eth_port_profile":         "eth_port_profile_",
	"lag":                      "lag",
	"sflow_collector":          "sflow_collector",
	"diagnostics_profile":      "diagnostics_profile",
	"diagnostics_port_profile": "diagnostics_port_profile",
	"service":                  "service",
	"eth_port_settings":        "eth_port_settings",
	"bundle":                   "endpoint_bundle",
	"badge":                    "badge",
	"authenticated_eth_port":   "authenticated_eth_port",
	"device_voice_settings":    "device_voice_settings",
	"packet_broker":            "pb_egress_profile",
	"packet_queue":             "packet_queue",
	"service_port_profile":     "service_port_profile",
	"voice_port_profile":       "voice_port_profiles",
	"switchpoint":              "switchpoint",
	"device_controller":        "device_controller",
	"as_path_access_list":      "as_path_access_list",
	"community_list":           "community_list",
	"device_settings":          "eth_device_profiles",
	"extended_community_list":  "extended_community_list",
	"ipv4_list":                "ipv4_list_filter",
	"ipv4_prefix_list":         "ipv4_prefix_list",
	"ipv6_list":                "ipv6_list_filter",
	"ipv6_prefix_list":         "ipv6_prefix_list",
	"route_map_clause":         "route_map_clause",
	"route_map":                "route_map",
	"sfp_breakout":             "sfp_breakouts",
	"site":                     "site",
	"pod":                      "pod",
	"spine_plane":              "spine_plane",
	"pb_routing_acl":           "pb_routing_acl",
	"pb_routing":               "pb_routing",
	"port_acl":                 "port_acl",
	"grouping_rule":            "grouping_rules",
	"threshold_group":          "threshold_group",
	"threshold":                "threshold",
	"acls_ipv4":                "ipv4_filter",
	"acls_ipv6":                "ipv6_filter",
}

// GetResourceJSONKey returns the JSON key for a given resource type
// Returns an empty string if the resource type is not found
func GetResourceJSONKey(resourceType string) string {
	if key, ok := ResourceJSONKeys[resourceType]; ok {
		return key
	}
	return ""
}

// GetACLJSONKey returns the appropriate JSON key for ACL resources based on IP version
// ipVersion should be "4" for IPv4 or "6" for IPv6
func GetACLJSONKey(ipVersion string) string {
	if ipVersion == "4" {
		return "ipv4_filter"
	}
	return "ipv6_filter"
}

// ImporterResourceMapping maps importer resource names (plural) to their API functions and JSON keys
// This provides compatibility with the importer's naming conventions
var ImporterResourceMapping = map[string]string{
	"tenants":                 "tenant",
	"gateways":                "gateway",
	"gatewayprofiles":         "gateway_profile",
	"ethportprofiles":         "eth_port_profile_",
	"lags":                    "lag",
	"sflowcollectors":         "sflow_collector",
	"diagnosticsprofiles":     "diagnostics_profile",
	"diagnosticsportprofiles": "diagnostics_port_profile",
	"services":                "service",
	"ethportsettings":         "eth_port_settings",
	"bundles":                 "endpoint_bundle",
	"badges":                  "badge",
	"authenticatedethports":   "authenticated_eth_port",
	"devicevoicesettings":     "device_voice_settings",
	"packetbroker":            "pb_egress_profile",
	"packetqueues":            "packet_queue",
	"serviceportprofiles":     "service_port_profile",
	"voiceportprofiles":       "voice_port_profiles",
	"switchpoints":            "switchpoint",
	"devicecontrollers":       "device_controller",
	"aspathaccesslists":       "as_path_access_list",
	"communitylists":          "community_list",
	"devicesettings":          "eth_device_profiles",
	"extendedcommunitylists":  "extended_community_list",
	"ipv4lists":               "ipv4_list_filter",
	"ipv4prefixlists":         "ipv4_prefix_list",
	"ipv6lists":               "ipv6_list_filter",
	"ipv6prefixlists":         "ipv6_prefix_list",
	"routemapclauses":         "route_map_clause",
	"routemaps":               "route_map",
	"sfpbreakouts":            "sfp_breakouts",
	"sites":                   "site",
	"pods":                    "pod",
	"spineplanes":             "spine_plane",
	"policybasedroutingacl":   "pb_routing_acl",
	"policybasedrouting":      "pb_routing",
	"portacls":                "port_acl",
	"groupingrules":           "grouping_rules",
	"thresholdgroups":         "threshold_group",
	"thresholds":              "threshold",
	"aclsipv4":                "ipv4_filter",
	"aclsipv6":                "ipv6_filter",
}

// GetImporterJSONKey returns the JSON key for a given importer resource name (plural form)
func GetImporterJSONKey(importerResourceName string) string {
	if key, ok := ImporterResourceMapping[importerResourceName]; ok {
		return key
	}
	return ""
}

package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
	Mode   string
}

type NestedBlockIterationStyle struct {
	PrintIndexFirst     bool
	SkipIndexInMainLoop bool
	IterateAllAsMap     bool
}

type ResourceConfig struct {
	ResourceType                 string
	StageName                    string
	HeaderNameLineFormat         string
	HeaderDependsOnLineFormat    string
	ObjectPropsHandler           func(objProps map[string]interface{}, builder *strings.Builder, config ResourceConfig)
	NestedBlockFields            map[string]bool
	ObjectPropsNestedBlockFields map[string]bool
	FieldMappings                map[string]string
	AdditionalTopLevelSkipKeys   []string
	EmptyObjectPropsAsSingleLine bool
	NestedBlockStyles            map[string]NestedBlockIterationStyle
}

type ImporterFunc func(context.Context, *openapi.APIClient) (*http.Response, error)

var nameSplitRE = regexp.MustCompile(`(\d+|\D+)`)

// importerRegistry maps resource names to their API caller function and JSON key
var importerRegistry = map[string]struct {
	apiCaller ImporterFunc
	jsonKey   string
}{
	"tenants": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.TenantsAPI.TenantsGet(ctx).Execute()
	}, jsonKey: "tenant"},
	"gateways": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.GatewaysAPI.GatewaysGet(ctx).Execute()
	}, jsonKey: "gateway"},
	"gatewayprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.GatewayProfilesAPI.GatewayprofilesGet(ctx).Execute()
	}, jsonKey: "gateway_profile"},
	"ethportprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
	}, jsonKey: "eth_port_profile_"},
	"lags": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.LAGsAPI.LagsGet(ctx).Execute()
	}, jsonKey: "lag"},
	"sflowcollectors": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.SFlowCollectorsAPI.SflowcollectorsGet(ctx).Execute()
	}, jsonKey: "sflow_collector"},
	"diagnosticsprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx).Execute()
	}, jsonKey: "diagnostics_profile"},
	"diagnosticsportprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx).Execute()
	}, jsonKey: "diagnostics_port_profile"},
	"services": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ServicesAPI.ServicesGet(ctx).Execute()
	}, jsonKey: "service"},
	"ethportsettings": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
	}, jsonKey: "eth_port_settings"},
	"bundles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.BundlesAPI.BundlesGet(ctx).Execute()
	}, jsonKey: "endpoint_bundle"},
	"badges": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.BadgesAPI.BadgesGet(ctx).Execute()
	}, jsonKey: "badge"},
	"authenticatedethports": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx).Execute()
	}, jsonKey: "authenticated_eth_port"},
	"devicevoicesettings": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx).Execute()
	}, jsonKey: "device_voice_settings"},
	"packetbroker": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PacketBrokerAPI.PacketbrokerGet(ctx).Execute()
	}, jsonKey: "pb_egress_profile"},
	"packetqueues": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PacketQueuesAPI.PacketqueuesGet(ctx).Execute()
	}, jsonKey: "packet_queue"},
	"serviceportprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ServicePortProfilesAPI.ServiceportprofilesGet(ctx).Execute()
	}, jsonKey: "service_port_profile"},
	"voiceportprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.VoicePortProfilesAPI.VoiceportprofilesGet(ctx).Execute()
	}, jsonKey: "voice_port_profiles"},
	"switchpoints": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
	}, jsonKey: "switchpoint"},
	"devicecontrollers": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.DeviceControllersAPI.DevicecontrollersGet(ctx).Execute()
	}, jsonKey: "device_controller"},
	"aspathaccesslists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ASPathAccessListsAPI.AspathaccesslistsGet(ctx).Execute()
	}, jsonKey: "as_path_access_list"},
	"communitylists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.CommunityListsAPI.CommunitylistsGet(ctx).Execute()
	}, jsonKey: "community_list"},
	"devicesettings": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.DeviceSettingsAPI.DevicesettingsGet(ctx).Execute()
	}, jsonKey: "eth_device_profiles"},
	"extendedcommunitylists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx).Execute()
	}, jsonKey: "extended_community_list"},
	"ipv4lists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.IPv4ListFiltersAPI.Ipv4listsGet(ctx).Execute()
	}, jsonKey: "ipv4_list_filter"},
	"ipv4prefixlists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx).Execute()
	}, jsonKey: "ipv4_prefix_list"},
	"ipv6lists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.IPv6ListFiltersAPI.Ipv6listsGet(ctx).Execute()
	}, jsonKey: "ipv6_list_filter"},
	"ipv6prefixlists": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx).Execute()
	}, jsonKey: "ipv6_prefix_list"},
	"routemapclauses": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.RouteMapClausesAPI.RoutemapclausesGet(ctx).Execute()
	}, jsonKey: "route_map_clause"},
	"routemaps": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.RouteMapsAPI.RoutemapsGet(ctx).Execute()
	}, jsonKey: "route_map"},
	"sfpbreakouts": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
	}, jsonKey: "sfp_breakouts"},
	"sites": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.SitesAPI.SitesGet(ctx).Execute()
	}, jsonKey: "site"},
	"pods": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PodsAPI.PodsGet(ctx).Execute()
	}, jsonKey: "pod"},
	"spineplanes": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.SpinePlanesAPI.SpineplanesGet(ctx).Execute()
	}, jsonKey: "spine_plane"},
	"policybasedroutingacl": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PBRoutingACLAPI.PolicybasedroutingaclGet(ctx).Execute()
	}, jsonKey: "pb_routing_acl"},
	"policybasedrouting": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PBRoutingAPI.PolicybasedroutingGet(ctx).Execute()
	}, jsonKey: "pb_routing"},
	"portacls": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.PortACLsAPI.PortaclsGet(ctx).Execute()
	}, jsonKey: "port_acl"},
	"groupingrules": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.GroupingRulesAPI.GroupingrulesGet(ctx).Execute()
	}, jsonKey: "grouping_rules"},
	"thresholdgroups": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ThresholdGroupsAPI.ThresholdgroupsGet(ctx).Execute()
	}, jsonKey: "threshold_group"},
	"thresholds": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ThresholdsAPI.ThresholdsGet(ctx).Execute()
	}, jsonKey: "threshold"},
}

// terraformTypeToResourceKey maps Terraform resource types to resourceConfigs keys
var terraformTypeToResourceKey = map[string]string{
	"verity_tenant":                   "tenant",
	"verity_gateway":                  "gateway",
	"verity_gateway_profile":          "gateway_profile",
	"verity_eth_port_profile":         "eth_port_profile",
	"verity_lag":                      "lag",
	"verity_sflow_collector":          "sflow_collector",
	"verity_diagnostics_profile":      "diagnostics_profile",
	"verity_diagnostics_port_profile": "diagnostics_port_profile",
	"verity_pb_routing_acl":           "pb_routing_acl",
	"verity_pb_routing":               "pb_routing",
	"verity_service":                  "service",
	"verity_eth_port_settings":        "eth_port_settings",
	"verity_bundle":                   "bundle",
	"verity_acl_v4":                   "acl_v4",
	"verity_acl_v6":                   "acl_v6",
	"verity_badge":                    "badge",
	"verity_authenticated_eth_port":   "authenticated_eth_port",
	"verity_device_controller":        "device_controller",
	"verity_device_voice_settings":    "device_voice_settings",
	"verity_packet_broker":            "packet_broker",
	"verity_packet_queue":             "packet_queue",
	"verity_service_port_profile":     "service_port_profile",
	"verity_voice_port_profile":       "voice_port_profile",
	"verity_spine_plane":              "spine_plane",
	"verity_switchpoint":              "switchpoint",
	"verity_as_path_access_list":      "as_path_access_list",
	"verity_community_list":           "community_list",
	"verity_device_settings":          "device_settings",
	"verity_extended_community_list":  "extended_community_list",
	"verity_ipv4_list":                "ipv4_list",
	"verity_ipv4_prefix_list":         "ipv4_prefix_list",
	"verity_ipv6_list":                "ipv6_list",
	"verity_ipv6_prefix_list":         "ipv6_prefix_list",
	"verity_route_map_clause":         "route_map_clause",
	"verity_route_map":                "route_map",
	"verity_sfp_breakout":             "sfp_breakout",
	"verity_site":                     "site",
	"verity_pod":                      "pod",
	"verity_port_acl":                 "port_acl",
	"verity_grouping_rule":            "grouping_rule",
	"verity_threshold_group":          "threshold_group",
	"verity_threshold":                "threshold",
}

// resourceConfigs is a registry of all resource configurations for generating Terraform code
var resourceConfigs = map[string]ResourceConfig{
	"tenant": {
		ResourceType:              "tenant",
		StageName:                 "tenant_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"route_tenants": true},
	},
	"gateway": {
		ResourceType:              "gateway",
		StageName:                 "gateway_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"static_routes": true},
	},
	"gateway_profile": {
		ResourceType:               "gateway_profile",
		StageName:                  "gateway_profile_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"external_gateways": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	},
	"eth_port_profile": {
		ResourceType:               "eth_port_profile",
		StageName:                  "eth_port_profile_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"services": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	},
	"lag": {
		ResourceType:                 "lag",
		StageName:                    "lag_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"sflow_collector": {
		ResourceType:                 "sflow_collector",
		StageName:                    "sflow_collector_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"diagnostics_profile": {
		ResourceType:                 "diagnostics_profile",
		StageName:                    "diagnostics_profile_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"diagnostics_port_profile": {
		ResourceType:                 "diagnostics_port_profile",
		StageName:                    "diagnostics_port_profile_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"pb_routing_acl": {
		ResourceType:              "pb_routing_acl",
		StageName:                 "pb_routing_acl_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"ipv4_permit": true, "ipv4_deny": true, "ipv6_permit": true, "ipv6_deny": true},
	},
	"pb_routing": {
		ResourceType:              "pb_routing",
		StageName:                 "pb_routing_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"policy": true},
	},
	"service": {
		ResourceType:              "service",
		StageName:                 "service_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"eth_port_settings": {
		ResourceType:              "eth_port_settings",
		StageName:                 "eth_port_settings_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lldp_med": true},
	},
	"bundle": {
		ResourceType:               "bundle",
		StageName:                  "bundle_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"eth_port_paths": true, "user_services": true, "rg_services": true, "voice_port_profile_paths": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
		NestedBlockStyles: map[string]NestedBlockIterationStyle{
			"eth_port_paths":           {IterateAllAsMap: true},
			"user_services":            {IterateAllAsMap: true},
			"rg_services":              {IterateAllAsMap: true},
			"voice_port_profile_paths": {IterateAllAsMap: true},
		},
	},
	"acl_v4": {
		ResourceType:              "acl_v4",
		StageName:                 "acl_v4_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"acl_v6": {
		ResourceType:              "acl_v6",
		StageName:                 "acl_v6_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"badge": {
		ResourceType:               "badge",
		StageName:                  "badge_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		AdditionalTopLevelSkipKeys: []string{},
	},
	"authenticated_eth_port": {
		ResourceType:              "authenticated_eth_port",
		StageName:                 "authenticated_eth_port_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"eth_ports": true, "object_properties": true},
	},
	"device_controller": {
		ResourceType:                 "device_controller",
		StageName:                    "device_controller_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"device_voice_settings": {
		ResourceType:                 "device_voice_settings",
		StageName:                    "device_voice_setting_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
		NestedBlockFields:            map[string]bool{"codecs": true},
		FieldMappings:                map[string]string{"Codecs": "codecs"},
	},
	"packet_broker": {
		ResourceType:              "packet_broker",
		StageName:                 "packet_broker_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"ipv4_permit": true, "ipv4_deny": true, "ipv6_permit": true, "ipv6_deny": true},
	},
	"packet_queue": {
		ResourceType:              "packet_queue",
		StageName:                 "packet_queue_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"pbit": true, "queue": true},
	},
	"service_port_profile": {
		ResourceType:              "service_port_profile",
		StageName:                 "service_port_profile_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"services": true},
	},
	"voice_port_profile": {
		ResourceType:                 "voice_port_profile",
		StageName:                    "voice_port_profile_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: false,
	},
	"spine_plane": {
		ResourceType:              "spine_plane",
		StageName:                 "spine_plane_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"switchpoint": {
		ResourceType:                 "switchpoint",
		StageName:                    "switchpoint_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		NestedBlockFields:            map[string]bool{"badges": true, "children": true, "traffic_mirrors": true, "eths": true},
		ObjectPropsNestedBlockFields: map[string]bool{"eths": true},
	},
	"as_path_access_list": {
		ResourceType:              "as_path_access_list",
		StageName:                 "as_path_access_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lists": true},
	},
	"community_list": {
		ResourceType:              "community_list",
		StageName:                 "community_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lists": true},
	},
	"device_settings": {
		ResourceType:              "device_settings",
		StageName:                 "device_settings_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"extended_community_list": {
		ResourceType:              "extended_community_list",
		StageName:                 "extended_community_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lists": true},
	},
	"ipv4_list": {
		ResourceType:              "ipv4_list",
		StageName:                 "ipv4_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"ipv4_prefix_list": {
		ResourceType:              "ipv4_prefix_list",
		StageName:                 "ipv4_prefix_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lists": true},
	},
	"ipv6_list": {
		ResourceType:              "ipv6_list",
		StageName:                 "ipv6_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"ipv6_prefix_list": {
		ResourceType:              "ipv6_prefix_list",
		StageName:                 "ipv6_prefix_list_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lists": true},
	},
	"route_map_clause": {
		ResourceType:              "route_map_clause",
		StageName:                 "route_map_clause_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"route_map": {
		ResourceType:              "route_map",
		StageName:                 "route_map_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"route_map_clauses": true},
	},
	"sfp_breakout": {
		ResourceType:              "sfp_breakout",
		StageName:                 "sfp_breakout_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"breakout": true},
	},
	"site": {
		ResourceType:                 "site",
		StageName:                    "site_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		NestedBlockFields:            map[string]bool{"islands": true, "pairs": true, "system_graphs": true},
		ObjectPropsNestedBlockFields: map[string]bool{"system_graphs": true},
	},
	"pod": {
		ResourceType:              "pod",
		StageName:                 "pod_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"port_acl": {
		ResourceType:              "port_acl",
		StageName:                 "port_acl_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"ipv4_permit": true, "ipv4_deny": true, "ipv6_permit": true, "ipv6_deny": true},
	},
	"grouping_rule": {
		ResourceType:              "grouping_rule",
		StageName:                 "grouping_rule_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"rules": true},
	},
	"threshold_group": {
		ResourceType:              "threshold_group",
		StageName:                 "threshold_group_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"targets": true, "thresholds": true},
	},
	"threshold": {
		ResourceType:              "threshold",
		StageName:                 "threshold_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"rules": true},
	},
}

func getNaturalSortParts(s string) []interface{} {
	matches := nameSplitRE.FindAllString(s, -1)
	parts := make([]interface{}, len(matches))
	for i, match := range matches {
		if num, err := strconv.Atoi(match); err == nil {
			parts[i] = num
		} else {
			parts[i] = match
		}
	}
	return parts
}

func NewImporter(client *openapi.APIClient, mode string) *Importer {
	return &Importer{
		client: client,
		ctx:    context.Background(),
		Mode:   mode,
	}
}

func (i *Importer) getAPIVersion() string {
	defaultVersion := "6.4"

	versionResp, err := i.client.VersionAPI.VersionGet(i.ctx).Execute()
	if err != nil {
		// API 6.4 doesn't support /version endpoint (returns 404)
		tflog.Info(i.ctx, "Version endpoint not available, using default version", map[string]interface{}{"default_version": defaultVersion})
		return defaultVersion
	}
	defer versionResp.Body.Close()

	// API 6.5+ should return a valid version response
	var versionData struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(versionResp.Body).Decode(&versionData); err != nil || versionData.Version == "" {
		tflog.Warn(i.ctx, "Failed to parse version response, using default", map[string]interface{}{"error": err, "default_version": defaultVersion})
		return defaultVersion
	}

	return versionData.Version
}

// ImportAll fetches all resources and saves them as Terraform configuration files
func (i *Importer) ImportAll(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	tflog.Info(i.ctx, "Starting importer with mode", map[string]interface{}{
		"mode": i.Mode,
	})

	apiVersionString := i.getAPIVersion()
	tflog.Info(i.ctx, "Using API Version for import", map[string]interface{}{"version": apiVersionString})

	stagesTF, err := i.generateStagesTF()
	if err != nil {
		tflog.Error(i.ctx, "Failed to generate stages TF", map[string]interface{}{"error": err})
		return fmt.Errorf("failed to generate stages: %w", err)
	}

	stagesFile := filepath.Join(outputDir, "stages.tf")
	if err := os.WriteFile(stagesFile, []byte(stagesTF), 0644); err != nil {
		tflog.Error(i.ctx, "Failed to write stages TF config", map[string]interface{}{"error": err, "file": stagesFile})
		return fmt.Errorf("failed to write stages terraform config: %w", err)
	}

	allResourceTasks := []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
	}{
		{name: "tenants", terraformResourceType: "verity_tenant", importer: func() (interface{}, error) { return i.importResource("tenants") }},
		{name: "gateways", terraformResourceType: "verity_gateway", importer: func() (interface{}, error) { return i.importResource("gateways") }},
		{name: "gatewayprofiles", terraformResourceType: "verity_gateway_profile", importer: func() (interface{}, error) { return i.importResource("gatewayprofiles") }},
		{name: "ethportprofiles", terraformResourceType: "verity_eth_port_profile", importer: func() (interface{}, error) { return i.importResource("ethportprofiles") }},
		{name: "lags", terraformResourceType: "verity_lag", importer: func() (interface{}, error) { return i.importResource("lags") }},
		{name: "sflowcollectors", terraformResourceType: "verity_sflow_collector", importer: func() (interface{}, error) { return i.importResource("sflowcollectors") }},
		{name: "diagnosticsprofiles", terraformResourceType: "verity_diagnostics_profile", importer: func() (interface{}, error) { return i.importResource("diagnosticsprofiles") }},
		{name: "diagnosticsportprofiles", terraformResourceType: "verity_diagnostics_port_profile", importer: func() (interface{}, error) { return i.importResource("diagnosticsportprofiles") }},
		{name: "policybasedroutingacl", terraformResourceType: "verity_pb_routing_acl", importer: func() (interface{}, error) { return i.importResource("policybasedroutingacl") }},
		{name: "policybasedrouting", terraformResourceType: "verity_pb_routing", importer: func() (interface{}, error) { return i.importResource("policybasedrouting") }},
		{name: "services", terraformResourceType: "verity_service", importer: func() (interface{}, error) { return i.importResource("services") }},
		{name: "ethportsettings", terraformResourceType: "verity_eth_port_settings", importer: func() (interface{}, error) { return i.importResource("ethportsettings") }},
		{name: "bundles", terraformResourceType: "verity_bundle", importer: func() (interface{}, error) { return i.importResource("bundles") }},
		{name: "acls_ipv4", terraformResourceType: "verity_acl_v4", importer: i.importACLsIPv4},
		{name: "acls_ipv6", terraformResourceType: "verity_acl_v6", importer: i.importACLsIPv6},
		{name: "badges", terraformResourceType: "verity_badge", importer: func() (interface{}, error) { return i.importResource("badges") }},
		{name: "authenticatedethports", terraformResourceType: "verity_authenticated_eth_port", importer: func() (interface{}, error) { return i.importResource("authenticatedethports") }},
		{name: "devicecontrollers", terraformResourceType: "verity_device_controller", importer: func() (interface{}, error) { return i.importResource("devicecontrollers") }},
		{name: "devicevoicesettings", terraformResourceType: "verity_device_voice_settings", importer: func() (interface{}, error) { return i.importResource("devicevoicesettings") }},
		{name: "packetbroker", terraformResourceType: "verity_packet_broker", importer: func() (interface{}, error) { return i.importResource("packetbroker") }},
		{name: "packetqueues", terraformResourceType: "verity_packet_queue", importer: func() (interface{}, error) { return i.importResource("packetqueues") }},
		{name: "serviceportprofiles", terraformResourceType: "verity_service_port_profile", importer: func() (interface{}, error) { return i.importResource("serviceportprofiles") }},
		{name: "voiceportprofiles", terraformResourceType: "verity_voice_port_profile", importer: func() (interface{}, error) { return i.importResource("voiceportprofiles") }},
		{name: "spineplanes", terraformResourceType: "verity_spine_plane", importer: func() (interface{}, error) { return i.importResource("spineplanes") }},
		{name: "switchpoints", terraformResourceType: "verity_switchpoint", importer: func() (interface{}, error) { return i.importResource("switchpoints") }},
		{name: "aspathaccesslists", terraformResourceType: "verity_as_path_access_list", importer: func() (interface{}, error) { return i.importResource("aspathaccesslists") }},
		{name: "communitylists", terraformResourceType: "verity_community_list", importer: func() (interface{}, error) { return i.importResource("communitylists") }},
		{name: "devicesettings", terraformResourceType: "verity_device_settings", importer: func() (interface{}, error) { return i.importResource("devicesettings") }},
		{name: "extendedcommunitylists", terraformResourceType: "verity_extended_community_list", importer: func() (interface{}, error) { return i.importResource("extendedcommunitylists") }},
		{name: "ipv4lists", terraformResourceType: "verity_ipv4_list", importer: func() (interface{}, error) { return i.importResource("ipv4lists") }},
		{name: "ipv4prefixlists", terraformResourceType: "verity_ipv4_prefix_list", importer: func() (interface{}, error) { return i.importResource("ipv4prefixlists") }},
		{name: "ipv6lists", terraformResourceType: "verity_ipv6_list", importer: func() (interface{}, error) { return i.importResource("ipv6lists") }},
		{name: "ipv6prefixlists", terraformResourceType: "verity_ipv6_prefix_list", importer: func() (interface{}, error) { return i.importResource("ipv6prefixlists") }},
		{name: "routemapclauses", terraformResourceType: "verity_route_map_clause", importer: func() (interface{}, error) { return i.importResource("routemapclauses") }},
		{name: "routemaps", terraformResourceType: "verity_route_map", importer: func() (interface{}, error) { return i.importResource("routemaps") }},
		{name: "sfpbreakouts", terraformResourceType: "verity_sfp_breakout", importer: func() (interface{}, error) { return i.importResource("sfpbreakouts") }},
		{name: "sites", terraformResourceType: "verity_site", importer: func() (interface{}, error) { return i.importResource("sites") }},
		{name: "pods", terraformResourceType: "verity_pod", importer: func() (interface{}, error) { return i.importResource("pods") }},
		{name: "portacls", terraformResourceType: "verity_port_acl", importer: func() (interface{}, error) { return i.importResource("portacls") }},
		{name: "groupingrules", terraformResourceType: "verity_grouping_rule", importer: func() (interface{}, error) { return i.importResource("groupingrules") }},
		{name: "thresholdgroups", terraformResourceType: "verity_threshold_group", importer: func() (interface{}, error) { return i.importResource("thresholdgroups") }},
		{name: "thresholds", terraformResourceType: "verity_threshold", importer: func() (interface{}, error) { return i.importResource("thresholds") }},
	}

	// Filter tasks based on mode and API version compatibility
	var resourceTasks []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
	}

	for _, task := range allResourceTasks {
		if utils.IsResourceCompatibleWithMode(task.terraformResourceType, i.Mode) {
			resourceTasks = append(resourceTasks, task)
		} else {
			tflog.Info(i.ctx, "Skipping resource due to mode incompatibility", map[string]interface{}{
				"resource_name":           task.name,
				"terraform_resource_type": task.terraformResourceType,
				"mode":                    i.Mode,
			})
		}
	}

	for _, task := range resourceTasks {
		if !utils.IsResourceCompatibleWithVersion(task.terraformResourceType, apiVersionString) {
			tflog.Info(i.ctx, "Skipping resource due to API version incompatibility", map[string]interface{}{
				"resource_name":           task.name,
				"terraform_resource_type": task.terraformResourceType,
				"api_version":             apiVersionString,
			})
			continue
		}

		tflog.Info(i.ctx, "Importing resource compatible with API version", map[string]interface{}{
			"resource_name":           task.name,
			"terraform_resource_type": task.terraformResourceType,
			"api_version":             apiVersionString,
		})

		data, err := task.importer()
		if err != nil {
			tflog.Error(i.ctx, "Failed to import resource", map[string]interface{}{"resource_name": task.name, "error": err})
			return fmt.Errorf("failed to import %s: %w", task.name, err)
		}

		if data == nil {
			tflog.Info(i.ctx, "No data returned by importer, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}
		if m, ok := data.(map[string]map[string]interface{}); ok && len(m) == 0 {
			tflog.Info(i.ctx, "No data found for resource, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}

		// Get the resource config key from the terraform type
		resourceKey, ok := terraformTypeToResourceKey[task.terraformResourceType]
		if !ok {
			tflog.Error(i.ctx, "No resource config found for terraform type", map[string]interface{}{
				"resource_name":  task.name,
				"terraform_type": task.terraformResourceType,
			})
			return fmt.Errorf("no resource config found for %s", task.terraformResourceType)
		}

		tfConfig, err := i.generateResourceTFByName(resourceKey, data)
		if err != nil {
			tflog.Error(i.ctx, "Failed to generate Terraform config", map[string]interface{}{"resource_name": task.name, "error": err})
			return fmt.Errorf("failed to generate terraform config for %s: %w", task.name, err)
		}

		if strings.TrimSpace(tfConfig) == "" {
			tflog.Info(i.ctx, "Generated TF config is empty, skipping file write", map[string]interface{}{"resource_name": task.name})
			continue
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.tf", task.name))
		if err := os.WriteFile(outputFile, []byte(tfConfig), 0644); err != nil {
			tflog.Error(i.ctx, "Failed to write TF config to file", map[string]interface{}{"resource_name": task.name, "file": outputFile, "error": err})
			return fmt.Errorf("failed to write %s terraform config: %w", task.name, err)
		}
		tflog.Info(i.ctx, "Successfully wrote TF config for resource", map[string]interface{}{"resource_name": task.name, "file": outputFile})
	}

	return nil
}

func (i *Importer) importResource(resourceName string) (interface{}, error) {
	config, ok := importerRegistry[resourceName]
	if !ok {
		return nil, fmt.Errorf("no importer configuration found for %s", resourceName)
	}

	resp, err := config.apiCaller(i.ctx, i.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %v", resourceName, err)
	}
	defer resp.Body.Close()

	result := make(map[string]map[string]map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %v", resourceName, err)
	}

	data, ok := result[config.jsonKey]
	if !ok {
		// Return empty map if the key doesn't exist
		return make(map[string]map[string]interface{}), nil
	}

	return data, nil
}

func (i *Importer) generateResourceTFByName(resourceKey string, data interface{}) (string, error) {
	cfg, ok := resourceConfigs[resourceKey]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %s", resourceKey)
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateResourceTF(data interface{}, config ResourceConfig) (string, error) {
	resourcesMap, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data format for resource type %s", config.ResourceType)
	}

	var resourceNames []string
	for name := range resourcesMap {
		resourceNames = append(resourceNames, name)
	}

	sort.SliceStable(resourceNames, func(i, j int) bool {
		s1 := resourceNames[i]
		s2 := resourceNames[j]

		parts1 := getNaturalSortParts(s1)
		parts2 := getNaturalSortParts(s2)

		len1 := len(parts1)
		len2 := len(parts2)
		minLen := len1
		if len2 < minLen {
			minLen = len2
		}

		for k := 0; k < minLen; k++ {
			p1 := parts1[k]
			p2 := parts2[k]

			p1Int, p1IsInt := p1.(int)
			p2Int, p2IsInt := p2.(int)

			if p1IsInt && p2IsInt {
				if p1Int != p2Int {
					return p1Int < p2Int
				}
			} else if !p1IsInt && !p2IsInt {
				p1Str := p1.(string)
				p2Str := p2.(string)
				if p1Str != p2Str {
					return p1Str < p2Str
				}
			} else {
				return p1IsInt
			}
		}
		return len1 < len2
	})

	var tfConfig strings.Builder

	for _, name := range resourceNames {
		resource := resourcesMap[name]
		sanitizedName := utils.SanitizeResourceName(name)

		tfConfig.WriteString(fmt.Sprintf("\nresource \"verity_%s\" \"%s\" {\n", config.ResourceType, sanitizedName))
		tfConfig.WriteString(fmt.Sprintf(config.HeaderNameLineFormat, name))
		tfConfig.WriteString(fmt.Sprintf(config.HeaderDependsOnLineFormat, config.StageName))

		// Skip object_properties section entirely if specified
		skipObjectProperties := false
		for _, key := range config.AdditionalTopLevelSkipKeys {
			if key == "object_properties" {
				skipObjectProperties = true
				break
			}
		}

		if !skipObjectProperties {
			// Check if object_properties exists in the API response
			objPropsRaw, objectPropertiesExists := resource["object_properties"]

			// Only include object_properties if it's actually present in the API response
			if objectPropertiesExists {
				tfConfig.WriteString("	object_properties")
				objProps, _ := objPropsRaw.(map[string]interface{})

				// Check if object_properties is empty (exists but has no fields)
				isEmptyObjectProps := len(objProps) == 0

				if isEmptyObjectProps && config.EmptyObjectPropsAsSingleLine {
					tfConfig.WriteString(" {}\n")
				} else if isEmptyObjectProps {
					// Case: object_properties exists but is empty
					tfConfig.WriteString(" {\n")
					tfConfig.WriteString("	}\n")
				} else {
					// Case: object_properties exists and has content
					tfConfig.WriteString(" {\n")
					var objPropsContentBuilder strings.Builder
					if config.ObjectPropsHandler != nil {
						config.ObjectPropsHandler(objProps, &objPropsContentBuilder, config)
					}
					tfConfig.WriteString(objPropsContentBuilder.String())
					tfConfig.WriteString("	}\n")
				}
			}
			// Case: object_properties doesn't exist in API response - don't include it at all
		}

		skipKeysSet := map[string]bool{
			"name":              true,
			"object_properties": true,
		}
		for _, key := range config.AdditionalTopLevelSkipKeys {
			skipKeysSet[key] = true
		}

		var topLevelKeys []string
		for key := range resource {
			if skipKeysSet[key] {
				continue
			}
			if isAutoAssignedField(resource, key) {
				continue
			}
			topLevelKeys = append(topLevelKeys, key)
		}
		sort.Strings(topLevelKeys)

		for _, key := range topLevelKeys {
			value := resource[key]

			tfFieldName := key
			if config.FieldMappings != nil {
				if mappedName, exists := config.FieldMappings[key]; exists {
					tfFieldName = mappedName
				}
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", tfFieldName, v))
			case float64:
				// Check if it's a whole number
				if v == float64(int(v)) {
					tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", tfFieldName, int(v)))
				} else {
					tfConfig.WriteString(fmt.Sprintf("	%s = %g\n", tfFieldName, v))
				}
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", tfFieldName, formatValue(v)))
			case []interface{}:
				if _, isNestedBlock := config.NestedBlockFields[tfFieldName]; isNestedBlock {
					style, hasStyle := config.NestedBlockStyles[tfFieldName]
					if !hasStyle {
						style = NestedBlockIterationStyle{PrintIndexFirst: true, SkipIndexInMainLoop: true, IterateAllAsMap: false}
					}

					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", tfFieldName))

							if style.IterateAllAsMap {
								for itemKey, itemValue := range itemMap {
									tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
								}
							} else {
								printedIndex := false
								if style.PrintIndexFirst {
									if indexVal, idxExists := itemMap["index"]; idxExists {
										if indexFloat, isFloat := indexVal.(float64); isFloat {
											tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(indexFloat)))
											printedIndex = true
										}
									}
								}

								var nestedItemKeys []string
								for itemKey := range itemMap {
									if style.SkipIndexInMainLoop && itemKey == "index" && printedIndex {
										continue
									}
									if itemKey == "index" && style.SkipIndexInMainLoop && !printedIndex {
										if style.SkipIndexInMainLoop {
											continue
										}
									}
									nestedItemKeys = append(nestedItemKeys, itemKey)
								}
								sort.Strings(nestedItemKeys)

								for _, itemKey := range nestedItemKeys {
									itemValue := itemMap[itemKey]
									tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
								}
							}
							tfConfig.WriteString("	}\n")
						}
					}
				} else {
					tfConfig.WriteString(fmt.Sprintf("	%s = [\n", tfFieldName))
					for _, item := range v {
						if str, ok := item.(string); ok {
							tfConfig.WriteString(fmt.Sprintf("		%s,\n", formatValue(str)))
						}
					}
					tfConfig.WriteString("	]\n")
				}
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", tfFieldName))
			}
		}
		tfConfig.WriteString("}\n\n")
	}
	return tfConfig.String(), nil
}

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	apiVersionString := i.getAPIVersion()
	tflog.Info(i.ctx, "Generating stages for mode and version", map[string]interface{}{
		"mode":    i.Mode,
		"version": apiVersionString,
	})

	// Define stage orderings for each mode
	type StageDefinition struct {
		StageName      string
		ResourceType   string
		DependsOnStage string // empty string means it's the first stage
	}

	var stageOrder []StageDefinition

	if i.Mode == "campus" {
		// CAMPUS mode staging order:
		// 1. IPv4 Lists, 2. IPv6 Lists, 3. ACLs v4, 4. ACLs v6, 5. PB Routing ACL,
		// 6. PB Routing, 7. Port ACLs, 8. Services, 9. Eth Port Profiles, 10. SFlow Collectors,
		// 11. Packet Queues, 12. Service Port Profiles, 13. Diagnostics Port Profiles,
		// 14. Device Voice Settings, 15. Authenticated Eth Ports, 16. Diagnostics Profiles,
		// 17. Eth Port Settings, 18. Voice Port Profiles, 19. Device Settings, 20. Lags,
		// 21. Bundles, 22. Badges, 23. Switchpoints, 24. Thresholds, 25. Grouping Rules,
		// 26. Threshold Groups, 27. Sites, 28. Device Controllers
		stageOrder = []StageDefinition{
			{"ipv4_list_stage", "verity_ipv4_list", ""},
			{"ipv6_list_stage", "verity_ipv6_list", "ipv4_list_stage"},
			{"acl_v4_stage", "verity_acl_v4", "ipv6_list_stage"},
			{"acl_v6_stage", "verity_acl_v6", "acl_v4_stage"},
			{"pb_routing_acl_stage", "verity_pb_routing_acl", "acl_v6_stage"},
			{"pb_routing_stage", "verity_pb_routing", "pb_routing_acl_stage"},
			{"port_acl_stage", "verity_port_acl", "pb_routing_stage"},
			{"service_stage", "verity_service", "port_acl_stage"},
			{"eth_port_profile_stage", "verity_eth_port_profile", "service_stage"},
			{"sflow_collector_stage", "verity_sflow_collector", "eth_port_profile_stage"},
			{"packet_queue_stage", "verity_packet_queue", "sflow_collector_stage"},
			{"service_port_profile_stage", "verity_service_port_profile", "packet_queue_stage"},
			{"diagnostics_port_profile_stage", "verity_diagnostics_port_profile", "service_port_profile_stage"},
			{"device_voice_setting_stage", "verity_device_voice_settings", "diagnostics_port_profile_stage"},
			{"authenticated_eth_port_stage", "verity_authenticated_eth_port", "device_voice_setting_stage"},
			{"diagnostics_profile_stage", "verity_diagnostics_profile", "authenticated_eth_port_stage"},
			{"eth_port_settings_stage", "verity_eth_port_settings", "diagnostics_profile_stage"},
			{"voice_port_profile_stage", "verity_voice_port_profile", "eth_port_settings_stage"},
			{"device_settings_stage", "verity_device_settings", "voice_port_profile_stage"},
			{"lag_stage", "verity_lag", "device_settings_stage"},
			{"bundle_stage", "verity_bundle", "lag_stage"},
			{"badge_stage", "verity_badge", "bundle_stage"},
			{"switchpoint_stage", "verity_switchpoint", "badge_stage"},
			{"threshold_stage", "verity_threshold", "switchpoint_stage"},
			{"grouping_rule_stage", "verity_grouping_rule", "threshold_stage"},
			{"threshold_group_stage", "verity_threshold_group", "grouping_rule_stage"},
			{"site_stage", "verity_site", "threshold_group_stage"},
			{"device_controller_stage", "verity_device_controller", "site_stage"},
		}
	} else {
		// DATACENTER mode staging order:
		// 1. SFP Breakouts, 2. IPv6 Prefix Lists, 3. Community Lists, 4. IPv4 Prefix Lists,
		// 5. Extended Community Lists, 6. AS Path Access Lists, 7. Route Map Clauses,
		// 8. ACLs v6, 9. ACLs v4, 10. Route Maps, 11. PB Routing ACL, 12. Tenants,
		// 13. PB Routing, 14. IPv4 Lists, 15. IPv6 Lists, 16. Services, 17. Port ACLs,
		// 18. Packet Broker, 19. Eth Port Profiles, 20. Packet Queues, 21. SFlow Collectors,
		// 22. Gateways, 23. Lags, 24. Eth Port Settings, 25. Diagnostics Profiles,
		// 26. Gateway Profiles, 27. Diagnostics Port Profiles, 28. Bundles, 29. Pods,
		// 30. Badges, 31. Spine Planes, 32. Switchpoints, 33. Device Settings,
		// 34. Thresholds, 35. Grouping Rules, 36. Threshold Groups, 37. Sites, 38. Device Controllers
		stageOrder = []StageDefinition{
			{"sfp_breakout_stage", "verity_sfp_breakout", ""},
			{"ipv6_prefix_list_stage", "verity_ipv6_prefix_list", "sfp_breakout_stage"},
			{"community_list_stage", "verity_community_list", "ipv6_prefix_list_stage"},
			{"ipv4_prefix_list_stage", "verity_ipv4_prefix_list", "community_list_stage"},
			{"extended_community_list_stage", "verity_extended_community_list", "ipv4_prefix_list_stage"},
			{"as_path_access_list_stage", "verity_as_path_access_list", "extended_community_list_stage"},
			{"route_map_clause_stage", "verity_route_map_clause", "as_path_access_list_stage"},
			{"acl_v6_stage", "verity_acl_v6", "route_map_clause_stage"},
			{"acl_v4_stage", "verity_acl_v4", "acl_v6_stage"},
			{"route_map_stage", "verity_route_map", "acl_v4_stage"},
			{"pb_routing_acl_stage", "verity_pb_routing_acl", "route_map_stage"},
			{"tenant_stage", "verity_tenant", "pb_routing_acl_stage"},
			{"pb_routing_stage", "verity_pb_routing", "tenant_stage"},
			{"ipv4_list_stage", "verity_ipv4_list", "pb_routing_stage"},
			{"ipv6_list_stage", "verity_ipv6_list", "ipv4_list_stage"},
			{"service_stage", "verity_service", "ipv6_list_stage"},
			{"port_acl_stage", "verity_port_acl", "service_stage"},
			{"packet_broker_stage", "verity_packet_broker", "port_acl_stage"},
			{"eth_port_profile_stage", "verity_eth_port_profile", "packet_broker_stage"},
			{"packet_queue_stage", "verity_packet_queue", "eth_port_profile_stage"},
			{"sflow_collector_stage", "verity_sflow_collector", "packet_queue_stage"},
			{"gateway_stage", "verity_gateway", "sflow_collector_stage"},
			{"lag_stage", "verity_lag", "gateway_stage"},
			{"eth_port_settings_stage", "verity_eth_port_settings", "lag_stage"},
			{"diagnostics_profile_stage", "verity_diagnostics_profile", "eth_port_settings_stage"},
			{"gateway_profile_stage", "verity_gateway_profile", "diagnostics_profile_stage"},
			{"diagnostics_port_profile_stage", "verity_diagnostics_port_profile", "gateway_profile_stage"},
			{"bundle_stage", "verity_bundle", "diagnostics_port_profile_stage"},
			{"pod_stage", "verity_pod", "bundle_stage"},
			{"badge_stage", "verity_badge", "pod_stage"},
			{"spine_plane_stage", "verity_spine_plane", "badge_stage"},
			{"switchpoint_stage", "verity_switchpoint", "spine_plane_stage"},
			{"device_settings_stage", "verity_device_settings", "switchpoint_stage"},
			{"threshold_stage", "verity_threshold", "device_settings_stage"},
			{"grouping_rule_stage", "verity_grouping_rule", "threshold_stage"},
			{"threshold_group_stage", "verity_threshold_group", "grouping_rule_stage"},
			{"site_stage", "verity_site", "threshold_group_stage"},
			{"device_controller_stage", "verity_device_controller", "site_stage"},
		}
	}

	// Filter stages based on resource compatibility with mode and API version
	var compatibleStages []StageDefinition
	var lastCompatibleStage string

	for _, stage := range stageOrder {
		if utils.IsResourceCompatible(stage.ResourceType, i.Mode, apiVersionString) {
			if stage.DependsOnStage != "" && lastCompatibleStage != "" && stage.DependsOnStage != lastCompatibleStage {
				stage.DependsOnStage = lastCompatibleStage
			}
			compatibleStages = append(compatibleStages, stage)
			lastCompatibleStage = stage.StageName
		} else {
			tflog.Debug(i.ctx, "Excluding stage for incompatible resource", map[string]interface{}{
				"stage_name":    stage.StageName,
				"resource_type": stage.ResourceType,
				"mode":          i.Mode,
				"api_version":   apiVersionString,
			})
		}
	}

	modeComment := strings.ToUpper(i.Mode)
	tfConfig.WriteString(fmt.Sprintf("\n# These resources establish ordering for bulk operations in %s mode\n", modeComment))

	for _, stage := range compatibleStages {
		tfConfig.WriteString(fmt.Sprintf("resource \"verity_operation_stage\" \"%s\" {\n", stage.StageName))

		if stage.DependsOnStage != "" {
			tfConfig.WriteString(fmt.Sprintf("  depends_on = [verity_operation_stage.%s]\n", stage.DependsOnStage))
		}

		tfConfig.WriteString("  lifecycle {\n")
		tfConfig.WriteString("    create_before_destroy = true\n")
		tfConfig.WriteString("  }\n")
		tfConfig.WriteString("}\n\n")
	}

	tflog.Info(i.ctx, "Generated stages", map[string]interface{}{
		"mode":              i.Mode,
		"api_version":       apiVersionString,
		"total_stages":      len(stageOrder),
		"compatible_stages": len(compatibleStages),
	})

	return tfConfig.String(), nil
}

func (i *Importer) importACLsIPv4() (interface{}, error) {
	return i.importACLs("4")
}

func (i *Importer) importACLsIPv6() (interface{}, error) {
	return i.importACLs("6")
}

func (i *Importer) importACLs(ipVersion string) (map[string]map[string]interface{}, error) {
	resp, err := i.client.ACLsAPI.AclsGet(i.ctx).IpVersion(ipVersion).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s ACLs: %v", ipVersion, err)
	}
	defer resp.Body.Close()

	if ipVersion == "4" {
		var result struct {
			IpFilter map[string]map[string]interface{} `json:"ipv4_filter"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode IPv4 ACLs response: %v", err)
		}

		if result.IpFilter == nil {
			return make(map[string]map[string]interface{}), nil
		}

		return result.IpFilter, nil
	} else {
		var result struct {
			IpFilter map[string]map[string]interface{} `json:"ipv6_filter"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode IPv6 ACLs response: %v", err)
		}

		if result.IpFilter == nil {
			return make(map[string]map[string]interface{}), nil
		}

		return result.IpFilter, nil
	}
}

// universalObjectPropsHandler dynamically processes all fields present in the object_properties
// section of the API response and converts them to Terraform configuration format.
// - If objProps is nil or empty map: generates no content
// - If objProps has fields: generates TF config for all fields (including nested structures)
// - Fields specified in ObjectPropsNestedBlockFields are rendered as blocks instead of attributes
func universalObjectPropsHandler(objProps map[string]interface{}, builder *strings.Builder, config ResourceConfig) {
	if len(objProps) > 0 {
		var keys []string
		for key := range objProps {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := objProps[key]

			if config.ObjectPropsNestedBlockFields != nil && config.ObjectPropsNestedBlockFields[key] {
				// Render as nested blocks
				if valueArray, ok := value.([]interface{}); ok {
					for _, item := range valueArray {
						builder.WriteString(fmt.Sprintf("		%s {\n", key))
						if itemMap, ok := item.(map[string]interface{}); ok {
							var itemKeys []string
							for itemKey := range itemMap {
								itemKeys = append(itemKeys, itemKey)
							}
							sort.Strings(itemKeys)

							for _, itemKey := range itemKeys {
								itemValue := itemMap[itemKey]
								builder.WriteString(fmt.Sprintf("			%s = %s\n", itemKey, formatObjectPropsValue(itemValue, "		")))
							}
						}
						builder.WriteString("		}\n")
					}
				}
			} else {
				// Render as attribute assignment
				builder.WriteString(fmt.Sprintf("		%s = %s\n", key, formatObjectPropsValue(value, "	")))
			}
		}
	}
	// If objProps is nil or empty, generate no content
}

func formatObjectPropsValue(value interface{}, indent string) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		// Check if it's a whole number
		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%g", v)
	case nil:
		return "null"
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}

		var result strings.Builder
		result.WriteString("[\n")
		for i, item := range v {
			result.WriteString(indent + "		")
			if itemMap, ok := item.(map[string]interface{}); ok {
				// Handle array of objects
				result.WriteString("{\n")
				var keys []string
				for key := range itemMap {
					keys = append(keys, key)
				}
				sort.Strings(keys)

				for _, key := range keys {
					itemValue := itemMap[key]
					result.WriteString(fmt.Sprintf("%s			%s = %s\n", indent, key, formatObjectPropsValue(itemValue, indent+"		")))
				}
				result.WriteString(indent + "		}")
			} else {
				// Handle array of primitives
				result.WriteString(formatObjectPropsValue(item, indent+"		"))
			}

			if i < len(v)-1 {
				result.WriteString(",")
			}
			result.WriteString("\n")
		}
		result.WriteString(indent + "	]")
		return result.String()
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}"
		}

		var result strings.Builder
		result.WriteString("{\n")
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for i, key := range keys {
			objValue := v[key]
			result.WriteString(fmt.Sprintf("%s		%s = %s", indent, key, formatObjectPropsValue(objValue, indent+"	")))
			if i < len(keys)-1 {
				result.WriteString(",")
			}
			result.WriteString("\n")
		}
		result.WriteString(indent + "	}")
		return result.String()
	default:
		return "null"
	}
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		// Check if it's a whole number
		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%g", v)
	case nil:
		return "null"
	default:
		return "null"
	}
}

func isAutoAssignedField(resource map[string]interface{}, fieldName string) bool {
	autoAssignedFieldName := fieldName + "_auto_assigned_"

	if autoAssignedValue, ok := resource[autoAssignedFieldName].(bool); ok && autoAssignedValue {
		return true
	}

	return false
}

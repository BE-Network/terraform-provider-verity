package lifecycle

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/tests/unit/mock"
)

type fieldInfo struct {
	Name     string
	Type     string
	Required bool
}

type blockInfo struct {
	Name   string
	Fields []fieldInfo
}

type resourceSchemaInfo struct {
	Attributes []fieldInfo
	Blocks     []blockInfo
}

type ResourceCoverageEntry struct {
	TerraformType       string
	Factory             func() resource.Resource
	APIPath             string
	WrapperKey          string
	Mode                string
	ResourceName        string
	Overrides           map[string]string
	SkipCreate          bool
	RequiredQueryParams map[string]string
}

// Each entry specifies the resource metadata needed to auto-generate HCL
// from schema introspection and verify PUT field coverage.
var allResourceTests = []ResourceCoverageEntry{
	// Resources available in both modes
	{
		TerraformType: "verity_badge",
		Factory:       provider.NewVerityBadgeResource,
		APIPath:       "/api/badges",
		WrapperKey:    "badge",
		Mode:          "datacenter",
		ResourceName:  "cov_badge",
	},
	{
		TerraformType: "verity_bundle",
		Factory:       provider.NewVerityBundleResource,
		APIPath:       "/api/bundles",
		WrapperKey:    "endpoint_bundle",
		Mode:          "datacenter",
		ResourceName:  "cov_bundle",
	},
	{
		TerraformType: "verity_service",
		Factory:       provider.NewVerityServiceResource,
		APIPath:       "/api/services",
		WrapperKey:    "service",
		Mode:          "datacenter",
		ResourceName:  "cov_service",
	},
	{
		TerraformType: "verity_site",
		Factory:       provider.NewVeritySiteResource,
		APIPath:       "/api/sites",
		WrapperKey:    "site",
		Mode:          "datacenter",
		ResourceName:  "cov_site",
		SkipCreate:    true, // Site is update-only
	},
	{
		TerraformType: "verity_eth_port_profile",
		Factory:       provider.NewVerityEthPortProfileResource,
		APIPath:       "/api/ethportprofiles",
		WrapperKey:    "eth_port_profile_",
		Mode:          "datacenter",
		ResourceName:  "cov_ethpp",
	},
	{
		TerraformType: "verity_eth_port_settings",
		Factory:       provider.NewVerityEthPortSettingsResource,
		APIPath:       "/api/ethportsettings",
		WrapperKey:    "eth_port_settings",
		Mode:          "datacenter",
		ResourceName:  "cov_ethps",
	},
	{
		TerraformType: "verity_lag",
		Factory:       provider.NewVerityLagResource,
		APIPath:       "/api/lags",
		WrapperKey:    "lag",
		Mode:          "datacenter",
		ResourceName:  "cov_lag",
	},
	{
		TerraformType:       "verity_acl_v4",
		Factory:             provider.NewVerityACLV4Resource,
		APIPath:             "/api/acls",
		WrapperKey:          "ip_filter",
		Mode:                "datacenter",
		ResourceName:        "cov_aclv4",
		RequiredQueryParams: map[string]string{"ip_version": "4"},
	},
	{
		TerraformType:       "verity_acl_v6",
		Factory:             provider.NewVerityACLV6Resource,
		APIPath:             "/api/acls",
		WrapperKey:          "ip_filter",
		Mode:                "datacenter",
		ResourceName:        "cov_aclv6",
		RequiredQueryParams: map[string]string{"ip_version": "6"},
	},
	{
		TerraformType: "verity_sflow_collector",
		Factory:       provider.NewVeritySflowCollectorResource,
		APIPath:       "/api/sflowcollectors",
		WrapperKey:    "sflow_collector",
		Mode:          "datacenter",
		ResourceName:  "cov_sflow",
	},
	{
		TerraformType: "verity_switchpoint",
		Factory:       provider.NewVeritySwitchpointResource,
		APIPath:       "/api/switchpoints",
		WrapperKey:    "switchpoint",
		Mode:          "datacenter",
		ResourceName:  "cov_sp",
	},
	{
		TerraformType: "verity_device_controller",
		Factory:       provider.NewVerityDeviceControllerResource,
		APIPath:       "/api/devicecontrollers",
		WrapperKey:    "device_controller",
		Mode:          "datacenter",
		ResourceName:  "cov_dc",
	},
	{
		TerraformType: "verity_device_settings",
		Factory:       provider.NewVerityDeviceSettingsResource,
		APIPath:       "/api/devicesettings",
		WrapperKey:    "eth_device_profiles",
		Mode:          "datacenter",
		ResourceName:  "cov_ds",
	},
	{
		TerraformType: "verity_packet_queue",
		Factory:       provider.NewVerityPacketQueueResource,
		APIPath:       "/api/packetqueues",
		WrapperKey:    "packet_queue",
		Mode:          "datacenter",
		ResourceName:  "cov_pq",
	},
	{
		TerraformType: "verity_diagnostics_profile",
		Factory:       provider.NewVerityDiagnosticsProfileResource,
		APIPath:       "/api/diagnosticsprofiles",
		WrapperKey:    "diagnostics_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_diagp",
	},
	{
		TerraformType: "verity_diagnostics_port_profile",
		Factory:       provider.NewVerityDiagnosticsPortProfileResource,
		APIPath:       "/api/diagnosticsportprofiles",
		WrapperKey:    "diagnostics_port_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_diagpp",
	},
	{
		TerraformType: "verity_ipv4_list",
		Factory:       provider.NewVerityIpv4ListResource,
		APIPath:       "/api/ipv4lists",
		WrapperKey:    "ipv4_list_filter",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv4l",
	},
	{
		TerraformType: "verity_ipv6_list",
		Factory:       provider.NewVerityIpv6ListResource,
		APIPath:       "/api/ipv6lists",
		WrapperKey:    "ipv6_list_filter",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv6l",
	},
	{
		TerraformType: "verity_port_acl",
		Factory:       provider.NewVerityPortAclResource,
		APIPath:       "/api/portacls",
		WrapperKey:    "port_acl",
		Mode:          "datacenter",
		ResourceName:  "cov_pacl",
	},
	{
		TerraformType: "verity_pb_routing",
		Factory:       provider.NewVerityPBRoutingResource,
		APIPath:       "/api/policybasedrouting",
		WrapperKey:    "pb_routing",
		Mode:          "datacenter",
		ResourceName:  "cov_pbr",
	},
	{
		TerraformType: "verity_pb_routing_acl",
		Factory:       provider.NewVerityPBRoutingACLResource,
		APIPath:       "/api/policybasedroutingacl",
		WrapperKey:    "pb_routing_acl",
		Mode:          "datacenter",
		ResourceName:  "cov_pbra",
	},
	{
		TerraformType: "verity_grouping_rule",
		Factory:       provider.NewVerityGroupingRuleResource,
		APIPath:       "/api/groupingrules",
		WrapperKey:    "grouping_rules",
		Mode:          "datacenter",
		ResourceName:  "cov_gr",
	},
	{
		TerraformType: "verity_threshold_group",
		Factory:       provider.NewVerityThresholdGroupResource,
		APIPath:       "/api/thresholdgroups",
		WrapperKey:    "threshold_group",
		Mode:          "datacenter",
		ResourceName:  "cov_tg",
	},
	{
		TerraformType: "verity_threshold",
		Factory:       provider.NewVerityThresholdResource,
		APIPath:       "/api/thresholds",
		WrapperKey:    "threshold",
		Mode:          "datacenter",
		ResourceName:  "cov_th",
	},

	// Datacenter-only resources
	{
		TerraformType: "verity_tenant",
		Factory:       provider.NewVerityTenantResource,
		APIPath:       "/api/tenants",
		WrapperKey:    "tenant",
		Mode:          "datacenter",
		ResourceName:  "cov_tenant",
		Overrides: map[string]string{
			"vrf_name": `"TestVrf"`,
		},
	},
	{
		TerraformType: "verity_gateway",
		Factory:       provider.NewVerityGatewayResource,
		APIPath:       "/api/gateways",
		WrapperKey:    "gateway",
		Mode:          "datacenter",
		ResourceName:  "cov_gw",
	},
	{
		TerraformType: "verity_gateway_profile",
		Factory:       provider.NewVerityGatewayProfileResource,
		APIPath:       "/api/gatewayprofiles",
		WrapperKey:    "gateway_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_gwp",
	},
	{
		TerraformType: "verity_as_path_access_list",
		Factory:       provider.NewVerityAsPathAccessListResource,
		APIPath:       "/api/aspathaccesslists",
		WrapperKey:    "as_path_access_list",
		Mode:          "datacenter",
		ResourceName:  "cov_apal",
	},
	{
		TerraformType: "verity_community_list",
		Factory:       provider.NewVerityCommunityListResource,
		APIPath:       "/api/communitylists",
		WrapperKey:    "community_list",
		Mode:          "datacenter",
		ResourceName:  "cov_cl",
	},
	{
		TerraformType: "verity_extended_community_list",
		Factory:       provider.NewVerityExtendedCommunityListResource,
		APIPath:       "/api/extendedcommunitylists",
		WrapperKey:    "extended_community_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ecl",
	},
	{
		TerraformType: "verity_ipv4_prefix_list",
		Factory:       provider.NewVerityIpv4PrefixListResource,
		APIPath:       "/api/ipv4prefixlists",
		WrapperKey:    "ipv4_prefix_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv4pl",
	},
	{
		TerraformType: "verity_ipv6_prefix_list",
		Factory:       provider.NewVerityIpv6PrefixListResource,
		APIPath:       "/api/ipv6prefixlists",
		WrapperKey:    "ipv6_prefix_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv6pl",
	},
	{
		TerraformType: "verity_route_map_clause",
		Factory:       provider.NewVerityRouteMapClauseResource,
		APIPath:       "/api/routemapclauses",
		WrapperKey:    "route_map_clause",
		Mode:          "datacenter",
		ResourceName:  "cov_rmc",
	},
	{
		TerraformType: "verity_route_map",
		Factory:       provider.NewVerityRouteMapResource,
		APIPath:       "/api/routemaps",
		WrapperKey:    "route_map",
		Mode:          "datacenter",
		ResourceName:  "cov_rm",
	},
	{
		TerraformType: "verity_sfp_breakout",
		Factory:       provider.NewVeritySfpBreakoutResource,
		APIPath:       "/api/sfpbreakouts",
		WrapperKey:    "sfp_breakouts",
		Mode:          "datacenter",
		ResourceName:  "cov_sfpb",
		SkipCreate:    true,
	},
	{
		TerraformType: "verity_pod",
		Factory:       provider.NewVerityPodResource,
		APIPath:       "/api/pods",
		WrapperKey:    "pod",
		Mode:          "datacenter",
		ResourceName:  "cov_pod",
	},
	{
		TerraformType: "verity_spine_plane",
		Factory:       provider.NewVeritySpinePlaneResource,
		APIPath:       "/api/spineplanes",
		WrapperKey:    "spine_plane",
		Mode:          "datacenter",
		ResourceName:  "cov_spinep",
	},
	{
		TerraformType: "verity_packet_broker",
		Factory:       provider.NewVerityPacketBrokerResource,
		APIPath:       "/api/packetbroker",
		WrapperKey:    "pb_egress_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_pb",
	},

	// Campus-only resources
	{
		TerraformType: "verity_authenticated_eth_port",
		Factory:       provider.NewVerityAuthenticatedEthPortResource,
		APIPath:       "/api/authenticatedethports",
		WrapperKey:    "authenticated_eth_port",
		Mode:          "campus",
		ResourceName:  "cov_aep",
	},
	{
		TerraformType: "verity_device_voice_settings",
		Factory:       provider.NewVerityDeviceVoiceSettingsResource,
		APIPath:       "/api/devicevoicesettings",
		WrapperKey:    "device_voice_settings",
		Mode:          "campus",
		ResourceName:  "cov_dvs",
	},
	{
		TerraformType: "verity_service_port_profile",
		Factory:       provider.NewVerityServicePortProfileResource,
		APIPath:       "/api/serviceportprofiles",
		WrapperKey:    "service_port_profile",
		Mode:          "campus",
		ResourceName:  "cov_spp",
	},
	{
		TerraformType: "verity_voice_port_profile",
		Factory:       provider.NewVerityVoicePortProfileResource,
		APIPath:       "/api/voiceportprofiles",
		WrapperKey:    "voice_port_profiles",
		Mode:          "campus",
		ResourceName:  "cov_vpp",
	},
}

func attrFieldType(attr fwschema.Attribute) string {
	switch attr.(type) {
	case fwschema.StringAttribute:
		return "string"
	case fwschema.BoolAttribute:
		return "bool"
	case fwschema.Int64Attribute:
		return "int64"
	case fwschema.NumberAttribute:
		return "number"
	default:
		return ""
	}
}

func inspectSchema(factory func() resource.Resource) resourceSchemaInfo {
	res := factory()
	var resp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	var rs resourceSchemaInfo
	for name, attr := range resp.Schema.Attributes {
		fi := fieldInfo{Name: name, Type: attrFieldType(attr)}
		if sa, ok := attr.(fwschema.StringAttribute); ok {
			fi.Required = sa.Required
		}
		rs.Attributes = append(rs.Attributes, fi)
	}
	sort.Slice(rs.Attributes, func(i, j int) bool {
		return rs.Attributes[i].Name < rs.Attributes[j].Name
	})

	for name, block := range resp.Schema.Blocks {
		if lb, ok := block.(fwschema.ListNestedBlock); ok {
			bi := blockInfo{Name: name}
			for attrName, attr := range lb.NestedObject.Attributes {
				bi.Fields = append(bi.Fields, fieldInfo{Name: attrName, Type: attrFieldType(attr)})
			}
			sort.Slice(bi.Fields, func(i, j int) bool {
				return bi.Fields[i].Name < bi.Fields[j].Name
			})
			rs.Blocks = append(rs.Blocks, bi)
		}
	}
	sort.Slice(rs.Blocks, func(i, j int) bool {
		return rs.Blocks[i].Name < rs.Blocks[j].Name
	})

	return rs
}

func defaultHCLValue(fi fieldInfo) string {
	if fi.Name == "index" {
		return "1"
	}
	switch fi.Type {
	case "string":
		return `""`
	case "bool":
		if strings.HasSuffix(fi.Name, "_auto_assigned_") {
			return "false"
		}
		return "true"
	case "int64":
		return "42"
	case "number":
		return "1.5"
	default:
		return `""`
	}
}

func generateCoverageHCL(rs resourceSchemaInfo, tfType, resourceName, mode, modeFieldsKey string, overrides map[string]string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("resource %q %q {\n", tfType, "test"))

	for _, fi := range rs.Attributes {
		if !utils.FieldAppliesToMode(modeFieldsKey, fi.Name, mode) {
			continue
		}
		val := defaultHCLValue(fi)
		if fi.Name == "name" {
			val = fmt.Sprintf("%q", resourceName)
		}
		if override, ok := overrides[fi.Name]; ok {
			val = override
		}
		b.WriteString(fmt.Sprintf("  %s = %s\n", fi.Name, val))
	}

	for _, block := range rs.Blocks {
		if !utils.FieldAppliesToMode(modeFieldsKey, block.Name, mode) {
			continue
		}
		// Skip blocks with no fields
		if len(block.Fields) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\n  %s {\n", block.Name))
		for _, fi := range block.Fields {
			nestedKey := block.Name + "." + fi.Name
			if !utils.FieldAppliesToMode(modeFieldsKey, nestedKey, mode) {
				continue
			}
			val := defaultHCLValue(fi)
			if override, ok := overrides[nestedKey]; ok {
				val = override
			}
			b.WriteString(fmt.Sprintf("    %s = %s\n", fi.Name, val))
		}
		b.WriteString("  }\n")
	}

	b.WriteString("}\n")
	return b.String()
}

func (e ResourceCoverageEntry) modeFieldsKey() string {
	return strings.TrimPrefix(e.APIPath, "/api/")
}

func TestFieldCoverage_SchemaDiscovery(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			rs := inspectSchema(tc.Factory)
			if len(rs.Attributes) == 0 {
				t.Errorf("no attributes discovered for %s", tc.TerraformType)
			}

			hasName := false
			for _, fi := range rs.Attributes {
				if fi.Name == "name" {
					hasName = true
					if !fi.Required {
						t.Errorf("%s: name attribute should be Required", tc.TerraformType)
					}
					break
				}
			}
			if !hasName {
				t.Errorf("%s: missing required 'name' attribute", tc.TerraformType)
			}

			t.Logf("%s: %d attributes, %d blocks", tc.TerraformType, len(rs.Attributes), len(rs.Blocks))
			for _, block := range rs.Blocks {
				t.Logf("  block %s: %d fields", block.Name, len(block.Fields))
			}
		})
	}
}

func TestFieldCoverage_PutContainsAllFields(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping create test", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			ms.SetTestLogger(t)
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()
			hcl := generateCoverageHCL(rs, tc.TerraformType, tc.ResourceName, tc.Mode, modeKey, tc.Overrides)

			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl
			t.Logf("Generated HCL:\n%s", hcl)

			fwresource.UnitTest(t, fwresource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []fwresource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							return verifyPutFieldCoverage(t, body, tc, rs)
						},
					},
				},
			})
		})
	}
}

func verifyPutFieldCoverage(t *testing.T, body map[string]interface{}, tc ResourceCoverageEntry, rs resourceSchemaInfo) error {
	t.Helper()

	wrapper, ok := body[tc.WrapperKey].(map[string]interface{})
	if !ok {
		return fmt.Errorf("wrapper key %q not found in PUT body (keys: %v)", tc.WrapperKey, mapKeys(body))
	}
	res, ok := wrapper[tc.ResourceName].(map[string]interface{})
	if !ok {
		return fmt.Errorf("resource %q not found under %q (keys: %v)", tc.ResourceName, tc.WrapperKey, mapKeys(wrapper))
	}

	modeKey := tc.modeFieldsKey()

	// Check top-level attributes
	for _, fi := range rs.Attributes {
		if !utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
			continue
		}
		if _, exists := res[fi.Name]; !exists {
			t.Errorf("[%s] attribute %q absent from PUT body", tc.TerraformType, fi.Name)
		}
	}

	// Check nested block attributes
	for _, block := range rs.Blocks {
		if !utils.FieldAppliesToMode(modeKey, block.Name, tc.Mode) {
			continue
		}
		// Skip blocks with no fields defined in schema
		if len(block.Fields) == 0 {
			continue
		}

		// Block may appear as a JSON array (most blocks) or a JSON object (object_properties).
		var item map[string]interface{}
		if items, ok := res[block.Name].([]interface{}); ok {
			if len(items) == 0 {
				t.Errorf("[%s] block %q is an empty array in PUT body", tc.TerraformType, block.Name)
				continue
			}
			item, ok = items[0].(map[string]interface{})
			if !ok {
				t.Errorf("[%s] block %q[0] is not an object", tc.TerraformType, block.Name)
				continue
			}
		} else if obj, ok := res[block.Name].(map[string]interface{}); ok {
			item = obj
		} else {
			t.Errorf("[%s] block %q absent from PUT body", tc.TerraformType, block.Name)
			continue
		}
		for _, fi := range block.Fields {
			nestedKey := block.Name + "." + fi.Name
			if !utils.FieldAppliesToMode(modeKey, nestedKey, tc.Mode) {
				continue
			}
			if _, exists := item[fi.Name]; !exists {
				t.Errorf("[%s] nested field %s.%s absent from PUT body", tc.TerraformType, block.Name, fi.Name)
			}
		}
	}

	return nil
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	Blocks []blockInfo
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

func init() {
	byType := map[string]func() resource.Resource{}
	for _, factory := range provider.New("test")().Resources(context.Background()) {
		var metadata resource.MetadataResponse
		factory().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		byType[metadata.TypeName] = factory
	}
	for i := range allResourceTests {
		factory, registered := byType[allResourceTests[i].TerraformType]
		if !registered {
			panic("coverage entry for unregistered resource " + allResourceTests[i].TerraformType)
		}
		allResourceTests[i].Factory = factory
	}
}

var allResourceTests = []ResourceCoverageEntry{

	{
		TerraformType: "verity_badge",
		APIPath:       "/api/badges",
		WrapperKey:    "badge",
		Mode:          "datacenter",
		ResourceName:  "cov_badge",
	},
	{
		TerraformType: "verity_bundle",
		APIPath:       "/api/bundles",
		WrapperKey:    "endpoint_bundle",
		Mode:          "datacenter",
		ResourceName:  "cov_bundle",
	},
	{
		TerraformType: "verity_service",
		APIPath:       "/api/services",
		WrapperKey:    "service",
		Mode:          "datacenter",
		ResourceName:  "cov_service",
	},
	{
		TerraformType: "verity_aaa_profile",
		APIPath:       "/api/deviceaaaprofiles",
		WrapperKey:    "device_aaa_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_aaa",
	},
	{
		TerraformType: "verity_ldap_profile",
		APIPath:       "/api/ldapprofiles",
		WrapperKey:    "ldap_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_ldap",
	},
	{
		TerraformType: "verity_fabric",
		APIPath:       "/api/fabrics",
		WrapperKey:    "fabric",
		Mode:          "datacenter",
		ResourceName:  "cov_fabric",
	},
	{
		TerraformType: "verity_plane",
		APIPath:       "/api/planes",
		WrapperKey:    "plane",
		Mode:          "datacenter",
		ResourceName:  "cov_plane",
	},
	{
		TerraformType: "verity_rack",
		APIPath:       "/api/racks",
		WrapperKey:    "rack",
		Mode:          "datacenter",
		ResourceName:  "cov_rack",
	},
	{
		TerraformType: "verity_eth_port_profile",
		APIPath:       "/api/ethportprofiles",
		WrapperKey:    "eth_port_profile_",
		Mode:          "datacenter",
		ResourceName:  "cov_ethpp",
	},
	{
		TerraformType: "verity_eth_port_settings",
		APIPath:       "/api/ethportsettings",
		WrapperKey:    "eth_port_settings",
		Mode:          "datacenter",
		ResourceName:  "cov_ethps",
	},
	{
		TerraformType: "verity_lag",
		APIPath:       "/api/lags",
		WrapperKey:    "lag",
		Mode:          "datacenter",
		ResourceName:  "cov_lag",
	},
	{
		TerraformType:       "verity_acl_v4",
		APIPath:             "/api/acls",
		WrapperKey:          "ip_filter",
		Mode:                "datacenter",
		ResourceName:        "cov_aclv4",
		RequiredQueryParams: map[string]string{"ip_version": "4"},
	},
	{
		TerraformType:       "verity_acl_v6",
		APIPath:             "/api/acls",
		WrapperKey:          "ip_filter",
		Mode:                "datacenter",
		ResourceName:        "cov_aclv6",
		RequiredQueryParams: map[string]string{"ip_version": "6"},
	},
	{
		TerraformType: "verity_sflow_collector",
		APIPath:       "/api/sflowcollectors",
		WrapperKey:    "sflow_collector",
		Mode:          "datacenter",
		ResourceName:  "cov_sflow",
	},
	{
		TerraformType: "verity_switchpoint",
		APIPath:       "/api/switchpoints",
		WrapperKey:    "switchpoint",
		Mode:          "datacenter",
		ResourceName:  "cov_sp",
	},
	{
		TerraformType: "verity_device_settings",
		APIPath:       "/api/devicesettings",
		WrapperKey:    "eth_device_profiles",
		Mode:          "datacenter",
		ResourceName:  "cov_ds",
	},
	{
		TerraformType: "verity_packet_queue",
		APIPath:       "/api/packetqueues",
		WrapperKey:    "packet_queue",
		Mode:          "datacenter",
		ResourceName:  "cov_pq",
	},
	{
		TerraformType: "verity_diagnostics_profile",
		APIPath:       "/api/diagnosticsprofiles",
		WrapperKey:    "diagnostics_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_diagp",
	},
	{
		TerraformType: "verity_diagnostics_port_profile",
		APIPath:       "/api/diagnosticsportprofiles",
		WrapperKey:    "diagnostics_port_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_diagpp",
	},
	{
		TerraformType: "verity_ipv4_list",
		APIPath:       "/api/ipv4lists",
		WrapperKey:    "ipv4_list_filter",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv4l",
	},
	{
		TerraformType: "verity_ipv6_list",
		APIPath:       "/api/ipv6lists",
		WrapperKey:    "ipv6_list_filter",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv6l",
	},
	{
		TerraformType: "verity_port_acl",
		APIPath:       "/api/portacls",
		WrapperKey:    "port_acl",
		Mode:          "datacenter",
		ResourceName:  "cov_pacl",
	},
	{
		TerraformType: "verity_pb_routing",
		APIPath:       "/api/policybasedrouting",
		WrapperKey:    "pb_routing",
		Mode:          "datacenter",
		ResourceName:  "cov_pbr",
	},
	{
		TerraformType: "verity_pb_routing_acl",
		APIPath:       "/api/policybasedroutingacl",
		WrapperKey:    "pb_routing_acl",
		Mode:          "datacenter",
		ResourceName:  "cov_pbra",
	},
	{
		TerraformType: "verity_grouping_rule",
		APIPath:       "/api/groupingrules",
		WrapperKey:    "grouping_rules",
		Mode:          "datacenter",
		ResourceName:  "cov_gr",
	},
	{
		TerraformType: "verity_threshold_group",
		APIPath:       "/api/thresholdgroups",
		WrapperKey:    "threshold_group",
		Mode:          "datacenter",
		ResourceName:  "cov_tg",
	},
	{
		TerraformType: "verity_threshold",
		APIPath:       "/api/thresholds",
		WrapperKey:    "threshold",
		Mode:          "datacenter",
		ResourceName:  "cov_th",
	},

	{
		TerraformType: "verity_tenant",
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
		APIPath:       "/api/gateways",
		WrapperKey:    "gateway",
		Mode:          "datacenter",
		ResourceName:  "cov_gw",
	},
	{
		TerraformType: "verity_gateway_profile",
		APIPath:       "/api/gatewayprofiles",
		WrapperKey:    "gateway_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_gwp",
	},
	{
		TerraformType: "verity_as_path_access_list",
		APIPath:       "/api/aspathaccesslists",
		WrapperKey:    "as_path_access_list",
		Mode:          "datacenter",
		ResourceName:  "cov_apal",
	},
	{
		TerraformType: "verity_community_list",
		APIPath:       "/api/communitylists",
		WrapperKey:    "community_list",
		Mode:          "datacenter",
		ResourceName:  "cov_cl",
	},
	{
		TerraformType: "verity_extended_community_list",
		APIPath:       "/api/extendedcommunitylists",
		WrapperKey:    "extended_community_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ecl",
	},
	{
		TerraformType: "verity_ipv4_prefix_list",
		APIPath:       "/api/ipv4prefixlists",
		WrapperKey:    "ipv4_prefix_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv4pl",
	},
	{
		TerraformType: "verity_ipv6_prefix_list",
		APIPath:       "/api/ipv6prefixlists",
		WrapperKey:    "ipv6_prefix_list",
		Mode:          "datacenter",
		ResourceName:  "cov_ipv6pl",
	},
	{
		TerraformType: "verity_route_map_clause",
		APIPath:       "/api/routemapclauses",
		WrapperKey:    "route_map_clause",
		Mode:          "datacenter",
		ResourceName:  "cov_rmc",
	},
	{
		TerraformType: "verity_route_map",
		APIPath:       "/api/routemaps",
		WrapperKey:    "route_map",
		Mode:          "datacenter",
		ResourceName:  "cov_rm",
	},
	{
		TerraformType: "verity_sfp_breakout",
		APIPath:       "/api/sfpbreakouts",
		WrapperKey:    "sfp_breakouts",
		Mode:          "datacenter",
		ResourceName:  "cov_sfpb",
		SkipCreate:    true,
	},
	{
		TerraformType: "verity_pod",
		APIPath:       "/api/pods",
		WrapperKey:    "pod",
		Mode:          "datacenter",
		ResourceName:  "cov_pod",
	},
	{
		TerraformType: "verity_spine_plane",
		APIPath:       "/api/spineplanes",
		WrapperKey:    "spine_plane",
		Mode:          "datacenter",
		ResourceName:  "cov_spinep",
	},
	{
		TerraformType: "verity_packet_broker",
		APIPath:       "/api/packetbroker",
		WrapperKey:    "pb_egress_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_pb",
	},

	{
		TerraformType: "verity_authenticated_eth_port",
		APIPath:       "/api/authenticatedethports",
		WrapperKey:    "authenticated_eth_port",
		Mode:          "campus",
		ResourceName:  "cov_aep",
	},
	{
		TerraformType: "verity_device_voice_settings",
		APIPath:       "/api/devicevoicesettings",
		WrapperKey:    "device_voice_settings",
		Mode:          "campus",
		ResourceName:  "cov_dvs",
	},
	{
		TerraformType: "verity_service_port_profile",
		APIPath:       "/api/serviceportprofiles",
		WrapperKey:    "service_port_profile",
		Mode:          "campus",
		ResourceName:  "cov_spp",
	},
	{
		TerraformType: "verity_voice_port_profile",
		APIPath:       "/api/voiceportprofiles",
		WrapperKey:    "voice_port_profiles",
		Mode:          "campus",
		ResourceName:  "cov_vpp",
	},

	{
		TerraformType: "verity_mac_filter",
		APIPath:       "/api/macfilters",
		WrapperKey:    "mac_filter",
		Mode:          "campus",
		ResourceName:  "cov_mac_filter",
	},
	{
		TerraformType: "verity_pair",
		APIPath:       "/api/pairs",
		WrapperKey:    "switch_pair",
		Mode:          "datacenter",
		ResourceName:  "cov_pair",
	},
	{
		TerraformType: "verity_ssp_group",
		APIPath:       "/api/sspgroups",
		WrapperKey:    "superspine_group",
		Mode:          "datacenter",
		ResourceName:  "cov_ssp_group",
	},
	{
		TerraformType: "verity_su",
		APIPath:       "/api/sus",
		WrapperKey:    "su",
		Mode:          "datacenter",
		ResourceName:  "cov_su",
	},
	{
		TerraformType: "verity_tacacs_profile",
		APIPath:       "/api/tacacsprofiles",
		WrapperKey:    "tacacs_profile",
		Mode:          "datacenter",
		ResourceName:  "cov_tacacs_profile",
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

func inspectBlock(name string, lb fwschema.ListNestedBlock) blockInfo {
	bi := blockInfo{Name: name}
	for attrName, attr := range lb.NestedObject.Attributes {
		bi.Fields = append(bi.Fields, fieldInfo{Name: attrName, Type: attrFieldType(attr)})
	}
	sort.Slice(bi.Fields, func(i, j int) bool { return bi.Fields[i].Name < bi.Fields[j].Name })
	for nestedName, nested := range lb.NestedObject.Blocks {
		if nestedList, ok := nested.(fwschema.ListNestedBlock); ok {
			bi.Blocks = append(bi.Blocks, inspectBlock(nestedName, nestedList))
		}
	}
	sort.Slice(bi.Blocks, func(i, j int) bool { return bi.Blocks[i].Name < bi.Blocks[j].Name })
	return bi
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
			rs.Blocks = append(rs.Blocks, inspectBlock(name, lb))
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

		if len(block.Fields) == 0 && len(block.Blocks) == 0 {
			continue
		}
		b.WriteString("\n")
		writeCoverageBlock(&b, block, block.Name, "  ", mode, modeFieldsKey, overrides)
	}

	b.WriteString("}\n")
	return b.String()
}

func writeCoverageBlock(b *strings.Builder, block blockInfo, path, indent, mode, modeFieldsKey string, overrides map[string]string) {
	fmt.Fprintf(b, "%s%s {\n", indent, block.Name)
	for _, fi := range block.Fields {
		nestedKey := path + "." + fi.Name
		if !utils.FieldAppliesToMode(modeFieldsKey, nestedKey, mode) {
			continue
		}
		val := defaultHCLValue(fi)
		if override, ok := overrides[nestedKey]; ok {
			val = override
		}
		fmt.Fprintf(b, "%s  %s = %s\n", indent, fi.Name, val)
	}
	for _, nested := range block.Blocks {
		nestedPath := path + "." + nested.Name
		if !utils.FieldAppliesToMode(modeFieldsKey, nestedPath, mode) {
			continue
		}
		if len(nested.Fields) == 0 && len(nested.Blocks) == 0 {
			continue
		}
		writeCoverageBlock(b, nested, nestedPath, indent+"  ", mode, modeFieldsKey, overrides)
	}
	fmt.Fprintf(b, "%s}\n", indent)
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

	for _, fi := range rs.Attributes {
		if !utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
			continue
		}
		if _, exists := res[fi.Name]; !exists {
			t.Errorf("[%s] attribute %q absent from PUT body", tc.TerraformType, fi.Name)
		}
	}

	for _, block := range rs.Blocks {
		verifyBlockFieldCoverage(t, res, block, block.Name, tc, modeKey)
	}

	return nil
}

func verifyBlockFieldCoverage(t *testing.T, parent map[string]interface{}, block blockInfo, path string, tc ResourceCoverageEntry, modeKey string) {
	t.Helper()
	if !utils.FieldAppliesToMode(modeKey, path, tc.Mode) {
		return
	}
	if len(block.Fields) == 0 && len(block.Blocks) == 0 {
		return
	}

	var item map[string]interface{}
	switch value := parent[block.Name].(type) {
	case []interface{}:
		if len(value) == 0 {
			t.Errorf("[%s] block %q is an empty array in PUT body", tc.TerraformType, path)
			return
		}
		object, ok := value[0].(map[string]interface{})
		if !ok {
			t.Errorf("[%s] block %q[0] is not an object", tc.TerraformType, path)
			return
		}
		item = object
	case map[string]interface{}:
		item = value
	default:
		t.Errorf("[%s] block %q absent from PUT body", tc.TerraformType, path)
		return
	}

	for _, fi := range block.Fields {
		nestedKey := path + "." + fi.Name
		if !utils.FieldAppliesToMode(modeKey, nestedKey, tc.Mode) {
			continue
		}
		if _, exists := item[fi.Name]; !exists {
			t.Errorf("[%s] nested field %s absent from PUT body", tc.TerraformType, nestedKey)
		}
	}
	for _, nested := range block.Blocks {
		verifyBlockFieldCoverage(t, item, nested, path+"."+nested.Name, tc, modeKey)
	}
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestCoverageTableMatchesRegistry(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "generated_registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	var artifact struct {
		Resources []struct {
			TerraformType string `json:"terraform_type"`
			API           struct {
				EndpointPath      string `json:"endpoint_path"`
				RequestWrapperKey string `json:"request_wrapper_key"`
			} `json:"api"`
		} `json:"resources"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatal(err)
	}

	byType := make(map[string]ResourceCoverageEntry, len(allResourceTests))
	for _, entry := range allResourceTests {
		if _, duplicate := byType[entry.TerraformType]; duplicate {
			t.Errorf("%s appears twice in the coverage table", entry.TerraformType)
		}
		byType[entry.TerraformType] = entry
	}
	if len(byType) != len(artifact.Resources) {
		t.Errorf("coverage table has %d resources, the registry has %d", len(byType), len(artifact.Resources))
	}
	for _, resourceSpec := range artifact.Resources {
		entry, exists := byType[resourceSpec.TerraformType]
		if !exists {
			t.Errorf("no coverage entry for %s", resourceSpec.TerraformType)
			continue
		}
		if want := "/api" + resourceSpec.API.EndpointPath; entry.APIPath != want {
			t.Errorf("%s coverage path = %q, registry endpoint is %q", resourceSpec.TerraformType, entry.APIPath, want)
		}
		if entry.WrapperKey != resourceSpec.API.RequestWrapperKey {
			t.Errorf("%s wrapper key = %q, registry says %q", resourceSpec.TerraformType, entry.WrapperKey, resourceSpec.API.RequestWrapperKey)
		}
		delete(byType, resourceSpec.TerraformType)
	}
	for name := range byType {
		t.Errorf("coverage entry %q has no registry resource", name)
	}
}

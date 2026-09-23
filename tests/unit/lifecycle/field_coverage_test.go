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
	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
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

type coverageTestData struct {
	ResourceName string
	Overrides    map[string]string
}

var coverageTestDataByType = map[string]coverageTestData{
	"verity_aaa_profile":              {ResourceName: "cov_aaa"},
	"verity_acl_v4":                   {ResourceName: "cov_aclv4"},
	"verity_acl_v6":                   {ResourceName: "cov_aclv6"},
	"verity_as_path_access_list":      {ResourceName: "cov_apal"},
	"verity_authenticated_eth_port":   {ResourceName: "cov_aep"},
	"verity_badge":                    {ResourceName: "cov_badge"},
	"verity_bundle":                   {ResourceName: "cov_bundle"},
	"verity_community_list":           {ResourceName: "cov_cl"},
	"verity_device_settings":          {ResourceName: "cov_ds"},
	"verity_device_voice_settings":    {ResourceName: "cov_dvs"},
	"verity_diagnostics_port_profile": {ResourceName: "cov_diagpp"},
	"verity_diagnostics_profile":      {ResourceName: "cov_diagp"},
	"verity_eth_port_profile":         {ResourceName: "cov_ethpp"},
	"verity_eth_port_settings":        {ResourceName: "cov_ethps"},
	"verity_extended_community_list":  {ResourceName: "cov_ecl"},
	"verity_fabric":                   {ResourceName: "cov_fabric"},
	"verity_gateway":                  {ResourceName: "cov_gw"},
	"verity_gateway_profile":          {ResourceName: "cov_gwp"},
	"verity_grouping_rule":            {ResourceName: "cov_gr"},
	"verity_ipv4_list":                {ResourceName: "cov_ipv4l"},
	"verity_ipv4_prefix_list":         {ResourceName: "cov_ipv4pl"},
	"verity_ipv6_list":                {ResourceName: "cov_ipv6l"},
	"verity_ipv6_prefix_list":         {ResourceName: "cov_ipv6pl"},
	"verity_lag":                      {ResourceName: "cov_lag"},
	"verity_ldap_profile":             {ResourceName: "cov_ldap"},
	"verity_mac_filter":               {ResourceName: "cov_mac_filter"},
	"verity_packet_broker":            {ResourceName: "cov_pb"},
	"verity_packet_queue":             {ResourceName: "cov_pq"},
	"verity_pair":                     {ResourceName: "cov_pair"},
	"verity_pb_routing":               {ResourceName: "cov_pbr"},
	"verity_pb_routing_acl":           {ResourceName: "cov_pbra"},
	"verity_plane":                    {ResourceName: "cov_plane"},
	"verity_pod":                      {ResourceName: "cov_pod"},
	"verity_port_acl":                 {ResourceName: "cov_pacl"},
	"verity_rack":                     {ResourceName: "cov_rack"},
	"verity_route_map":                {ResourceName: "cov_rm"},
	"verity_route_map_clause":         {ResourceName: "cov_rmc"},
	"verity_service":                  {ResourceName: "cov_service"},
	"verity_service_port_profile":     {ResourceName: "cov_spp"},
	"verity_sflow_collector":          {ResourceName: "cov_sflow"},
	"verity_sfp_breakout":             {ResourceName: "cov_sfpb"},
	"verity_spine_plane":              {ResourceName: "cov_spinep"},
	"verity_ssp_group":                {ResourceName: "cov_ssp_group"},
	"verity_su":                       {ResourceName: "cov_su"},
	"verity_switchpoint":              {ResourceName: "cov_sp"},
	"verity_tacacs_profile":           {ResourceName: "cov_tacacs_profile"},
	"verity_tenant": {
		ResourceName: "cov_tenant",
		Overrides:    map[string]string{"vrf_name": `"TestVrf"`},
	},
	"verity_threshold":          {ResourceName: "cov_th"},
	"verity_threshold_group":    {ResourceName: "cov_tg"},
	"verity_voice_port_profile": {ResourceName: "cov_vpp"},
}

var allResourceTests = mustCoverageEntries()

func mustCoverageEntries() []ResourceCoverageEntry {
	entries, err := coverageEntries()
	if err != nil {
		panic("build resource coverage entries: " + err.Error())
	}
	return entries
}

func coverageEntries() ([]ResourceCoverageEntry, error) {
	resources, err := registry.Load()
	if err != nil {
		return nil, err
	}
	return coverageEntriesFor(resources, coverageTestDataByType)
}

func coverageEntriesFor(resources spec.Registry, testDataByType map[string]coverageTestData) ([]ResourceCoverageEntry, error) {
	entries := make([]ResourceCoverageEntry, 0, len(resources))
	registryTypes := make(map[string]bool, len(resources))
	for _, resourceSpec := range resources {
		registryTypes[resourceSpec.TerraformType] = true
		testData := testDataByType[resourceSpec.TerraformType]
		resourceName := testData.ResourceName
		if resourceName == "" {
			resourceName = "cov_" + strings.TrimPrefix(resourceSpec.TerraformType, "verity_")
		}
		mode, err := coverageMode(resourceSpec.Modes)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", resourceSpec.TerraformType, err)
		}
		entries = append(entries, ResourceCoverageEntry{
			TerraformType:       resourceSpec.TerraformType,
			APIPath:             "/api" + resourceSpec.API.EndpointPath,
			WrapperKey:          resourceSpec.API.RequestWrapperKey,
			Mode:                mode,
			ResourceName:        resourceName,
			Overrides:           testData.Overrides,
			SkipCreate:          !resourceSpec.Operations.Create,
			RequiredQueryParams: resourceSpec.API.FixedHeaders,
		})
	}
	for terraformType := range testDataByType {
		if !registryTypes[terraformType] {
			return nil, fmt.Errorf("test data exists for registry-absent resource %s", terraformType)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].TerraformType < entries[j].TerraformType })
	return entries, nil
}

func coverageMode(modes []spec.Mode) (string, error) {
	for _, mode := range modes {
		if mode == spec.ModeDatacenter {
			return string(mode), nil
		}
	}
	for _, mode := range modes {
		if mode == spec.ModeCampus {
			return string(mode), nil
		}
	}
	return "", fmt.Errorf("resource supports no coverage mode")
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
							put := puts[len(puts)-1]
							if err := verifyPutFieldCoverage(t, put.Body, tc, rs); err != nil {
								return err
							}
							return verifyRequiredQueryParams(put.QueryParams, tc)
						},
					},
				},
			})
		})
	}
}

func verifyRequiredQueryParams(params map[string][]string, tc ResourceCoverageEntry) error {
	for name, expected := range tc.RequiredQueryParams {
		values := params[name]
		if len(values) == 0 || values[0] != expected {
			return fmt.Errorf("%s query parameter %q = %v, want [%s]", tc.TerraformType, name, values, expected)
		}
	}
	return nil
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

func TestCoverageEntriesReadTheRegistry(t *testing.T) {
	if _, err := coverageEntries(); err != nil {
		t.Fatal(err)
	}
}

func TestCoverageEntriesDefaultTestDataFromTheRegistry(t *testing.T) {
	entries, err := coverageEntriesFor(spec.Registry{{
		TerraformType: "verity_new_resource",
		Modes:         []spec.Mode{spec.ModeCampus, spec.ModeDatacenter},
		API: spec.APIResourceSpec{
			EndpointPath:      "/newresources",
			RequestWrapperKey: "new_resource",
			FixedHeaders:      map[string]string{"scope": "global"},
		},
		Operations: spec.OperationSpec{Create: false},
	}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entry count = %d, want 1", len(entries))
	}
	entry := entries[0]
	if entry.ResourceName != "cov_new_resource" {
		t.Errorf("resource name = %q, want cov_new_resource", entry.ResourceName)
	}
	if entry.APIPath != "/api/newresources" {
		t.Errorf("API path = %q, want /api/newresources", entry.APIPath)
	}
	if entry.WrapperKey != "new_resource" {
		t.Errorf("wrapper key = %q, want new_resource", entry.WrapperKey)
	}
	if entry.Mode != "datacenter" {
		t.Errorf("mode = %q, want datacenter", entry.Mode)
	}
	if !entry.SkipCreate {
		t.Error("SkipCreate = false, want true")
	}
	if got := entry.RequiredQueryParams["scope"]; got != "global" || len(entry.RequiredQueryParams) != 1 {
		t.Errorf("required query parameters = %v, want map[scope:global]", entry.RequiredQueryParams)
	}
	if err := verifyRequiredQueryParams(map[string][]string{"scope": {"global"}}, entry); err != nil {
		t.Error(err)
	}
}

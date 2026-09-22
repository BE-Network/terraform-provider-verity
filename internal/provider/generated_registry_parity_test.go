package provider

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

type generatedRegistryArtifact struct {
	Resources spec.Registry `json:"resources"`
}

func TestGeneratedSpecsMatchLegacySchemas(t *testing.T) {
	registry := generatedRegistry(t)

	constructors := legacyConstructorsByType(t)

	cacheKeys := map[string]string{
		"verity_acl_v4":                   "acls_ipv4",
		"verity_acl_v6":                   "acls_ipv6",
		"verity_aaa_profile":              "device_aaa_profiles",
		"verity_as_path_access_list":      "as_path_access_lists",
		"verity_authenticated_eth_port":   "authenticated_eth_ports",
		"verity_badge":                    "badges",
		"verity_bundle":                   "bundles",
		"verity_community_list":           "community_lists",
		"verity_device_voice_settings":    "device_voice_settings",
		"verity_diagnostics_port_profile": "diagnostics_port_profiles",
		"verity_diagnostics_profile":      "diagnostics_profiles",
		"verity_eth_port_profile":         "eth_port_profiles",
		"verity_eth_port_settings":        "eth_port_settings",
		"verity_extended_community_list":  "extended_community_lists",
		"verity_gateway_profile":          "gateway_profiles",
		"verity_grouping_rule":            "grouping_rules",
		"verity_device_settings":          "device_settings",
		"verity_fabric":                   "fabrics",
		"verity_gateway":                  "gateways",
		"verity_sfp_breakout":             "sfp_breakouts",
		"verity_ipv4_list":                "ipv4_lists",
		"verity_ipv4_prefix_list":         "ipv4_prefix_lists",
		"verity_ipv6_list":                "ipv6_lists",
		"verity_ipv6_prefix_list":         "ipv6_prefix_lists",
		"verity_lag":                      "lags",
		"verity_ldap_profile":             "ldap_profiles",
		"verity_mac_filter":               "mac_filters",
		"verity_packet_broker":            "packet_brokers",
		"verity_packet_queue":             "packet_queues",
		"verity_pair":                     "pairs",
		"verity_pb_routing":               "pb_routing",
		"verity_pb_routing_acl":           "pb_routing_acl",
		"verity_plane":                    "planes",
		"verity_pod":                      "pods",
		"verity_port_acl":                 "port_acls",
		"verity_rack":                     "racks",
		"verity_route_map":                "route_maps",
		"verity_route_map_clause":         "route_map_clauses",
		"verity_service":                  "services",
		"verity_service_port_profile":     "service_port_profiles",
		"verity_sflow_collector":          "sflow_collectors",
		"verity_spine_plane":              "spine_planes",
		"verity_ssp_group":                "ssp_groups",
		"verity_su":                       "sus",
		"verity_switchpoint":              "switchpoints",
		"verity_tacacs_profile":           "tacacs_profiles",
		"verity_tenant":                   "tenants",
		"verity_threshold":                "thresholds",
		"verity_threshold_group":          "threshold_groups",
		"verity_voice_port_profile":       "voice_port_profiles",
	}
	if len(cacheKeys) != len(registry) {
		t.Fatalf("legacy parity mappings cover %d cache keys for %d generated resources", len(cacheKeys), len(registry))
	}
	for _, resourceSpec := range registry {
		constructor, exists := constructors[resourceSpec.TerraformType]
		if !exists {
			t.Fatalf("legacy schema constructor is missing for %q", resourceSpec.TerraformType)
		}
		if resourceSpec.API.CacheKey != cacheKeys[resourceSpec.TerraformType] {
			t.Fatalf("generated cache key = %q, legacy cache key = %q for %s", resourceSpec.API.CacheKey, cacheKeys[resourceSpec.TerraformType], resourceSpec.TerraformType)
		}
		assertGeneratedModes(t, resourceSpec)
		assertGeneratedSchemaMatchesLegacy(t, resourceSpec, constructor())
	}
}

func legacyConstructorsByType(t *testing.T) map[string]func() resource.Resource {
	t.Helper()
	ctx := context.Background()
	byType := make(map[string]func() resource.Resource)
	for _, constructor := range getAllResources() {
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		byType[metadata.TypeName] = constructor
	}
	return byType
}

func assertGeneratedModes(t *testing.T, resourceSpec spec.ResourceSpec) {
	t.Helper()
	legacyMode, exists := utils.ResourceCompatibility[resourceSpec.TerraformType]
	if !exists {
		t.Fatalf("legacy compatibility map is missing %q", resourceSpec.TerraformType)
	}
	expected := map[utils.ResourceMode][]spec.Mode{
		utils.ResourceModeDatacenter: {spec.ModeDatacenter},
		utils.ResourceModeCampus:     {spec.ModeCampus},
		utils.ResourceModeBoth:       {spec.ModeCampus, spec.ModeDatacenter},
	}[legacyMode]
	if len(expected) != len(resourceSpec.Modes) {
		t.Fatalf("legacy mode %q maps to %v, generated modes = %v for %s", legacyMode, expected, resourceSpec.Modes, resourceSpec.TerraformType)
	}
	for index, mode := range expected {
		if resourceSpec.Modes[index] != mode {
			t.Fatalf("legacy mode %q maps to %v, generated modes = %v for %s", legacyMode, expected, resourceSpec.Modes, resourceSpec.TerraformType)
		}
	}
}

func assertGeneratedSchemaMatchesLegacy(t *testing.T, resourceSpec spec.ResourceSpec, legacy resource.Resource) {
	t.Helper()
	ctx := context.Background()
	var metadata resource.MetadataResponse
	legacy.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
	if metadata.TypeName != resourceSpec.TerraformType {
		t.Fatalf("legacy Terraform type = %q, generated type = %q", metadata.TypeName, resourceSpec.TerraformType)
	}
	var response resource.SchemaResponse
	legacy.Schema(ctx, resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("legacy schema diagnostics: %s", response.Diagnostics)
	}
	assertGeneratedFieldsMatchSchema(t, resourceSpec.Fields, response.Schema.Attributes, response.Schema.Blocks, resourceSpec.TerraformType)
}

func assertGeneratedFieldsMatchSchema(t *testing.T, fields []spec.FieldSpec, attributes map[string]schema.Attribute, blocks map[string]schema.Block, path string) {
	t.Helper()

	managed := make([]spec.FieldSpec, 0, len(fields))
	for _, field := range fields {
		if !field.Unmanaged {
			managed = append(managed, field)
		}
	}
	if len(attributes)+len(blocks) != len(managed) {
		t.Fatalf("%s has %d legacy attributes and %d blocks for %d managed generated fields", path, len(attributes), len(blocks), len(managed))
	}
	for _, field := range managed {

		if field.Kind != spec.FieldKindList && field.Kind != spec.FieldKindObject {
			attribute, exists := attributes[field.TerraformName]
			if !exists {
				t.Fatalf("%s is missing generated field %q", path, field.TerraformName)
			}
			assertGeneratedAttribute(t, field, attribute)
			continue
		}
		if field.Kind == spec.FieldKindObject && field.Collection.Strategy != spec.CollectionSingleton {
			t.Fatalf("generated object field %q must use the singleton collection strategy, got %q", field.TerraformName, field.Collection.Strategy)
		}
		block, exists := blocks[field.TerraformName]
		if !exists {
			t.Fatalf("%s is missing generated block %q", path, field.TerraformName)
		}
		list, ok := block.(schema.ListNestedBlock)
		if !ok || list.Description != field.Description {
			t.Fatalf("legacy block %q does not match generated metadata", field.TerraformName)
		}
		assertGeneratedFieldsMatchSchema(t, field.Fields, list.NestedObject.Attributes, list.NestedObject.Blocks, path+"."+field.TerraformName)
	}
}

func generatedRegistry(t *testing.T) spec.Registry {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "specs", "generated_registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	var artifact generatedRegistryArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatal(err)
	}
	if err := artifact.Resources.Validate(); err != nil {
		t.Fatalf("generated registry validation failed: %v", err)
	}
	return artifact.Resources
}

func assertGeneratedAttribute(t *testing.T, field spec.FieldSpec, attribute schema.Attribute) {
	t.Helper()
	switch value := attribute.(type) {
	case schema.StringAttribute:
		assertGeneratedStringField(t, field, value)
	case schema.BoolAttribute:
		assertGeneratedBoolField(t, field, value)
	case schema.Int64Attribute:
		if field.Kind != spec.FieldKindInt64 || field.Description != value.Description || field.Sensitive != value.Sensitive {
			t.Fatalf("legacy int64 field %q does not match generated metadata", field.TerraformName)
		}
		assertGeneratedAccess(t, field, value.Required, value.Optional, value.Computed)
		assertGeneratedReplace(t, field, countRequiresReplace(value.PlanModifiers))
	case schema.NumberAttribute:
		if field.Kind != spec.FieldKindNumber || field.Description != value.Description || field.Sensitive != value.Sensitive {
			t.Fatalf("legacy number field %q does not match generated metadata", field.TerraformName)
		}
		assertGeneratedAccess(t, field, value.Required, value.Optional, value.Computed)
		assertGeneratedReplace(t, field, countRequiresReplace(value.PlanModifiers))
	default:
		t.Fatalf("legacy field %q has unsupported attribute type %T", field.TerraformName, attribute)
	}
}

func assertGeneratedStringField(t *testing.T, field spec.FieldSpec, attribute schema.StringAttribute) {
	t.Helper()
	if field.Kind != spec.FieldKindString || field.Description != attribute.Description || field.Sensitive != attribute.Sensitive {
		t.Fatalf("legacy string field %q does not match generated metadata", field.TerraformName)
	}
	assertGeneratedAccess(t, field, attribute.Required, attribute.Optional, attribute.Computed)
	assertGeneratedReplace(t, field, countRequiresReplace(attribute.PlanModifiers))
}

func assertGeneratedBoolField(t *testing.T, field spec.FieldSpec, attribute schema.BoolAttribute) {
	t.Helper()
	if field.Kind != spec.FieldKindBool || field.Description != attribute.Description || field.Sensitive != attribute.Sensitive {
		t.Fatalf("legacy bool field %q does not match generated metadata", field.TerraformName)
	}
	assertGeneratedAccess(t, field, attribute.Required, attribute.Optional, attribute.Computed)
	assertGeneratedReplace(t, field, countRequiresReplace(attribute.PlanModifiers))
}

func assertGeneratedAccess(t *testing.T, field spec.FieldSpec, required, optional, computed bool) {
	t.Helper()
	if required != (field.Access == spec.AccessRequired) || optional != (field.Access == spec.AccessOptional || field.Access == spec.AccessOptionalComputed) || computed != (field.Access == spec.AccessComputed || field.Access == spec.AccessOptionalComputed) {
		t.Fatalf("legacy field %q access does not match generated access %q", field.TerraformName, field.Access)
	}
}

var requiresReplaceTypes = map[reflect.Type]bool{
	reflect.TypeOf(stringplanmodifier.RequiresReplace()): true,
	reflect.TypeOf(boolplanmodifier.RequiresReplace()):   true,
	reflect.TypeOf(int64planmodifier.RequiresReplace()):  true,
	reflect.TypeOf(numberplanmodifier.RequiresReplace()): true,
}

func countRequiresReplace[T any](modifiers []T) int {
	count := 0
	for _, modifier := range modifiers {
		if requiresReplaceTypes[reflect.TypeOf(modifier)] {
			count++
		}
	}
	return count
}

func assertGeneratedReplace(t *testing.T, field spec.FieldSpec, requiresReplaceCount int) {
	t.Helper()
	if field.Replace != (requiresReplaceCount > 0) {
		t.Fatalf("legacy field %q replace = %v but the schema has %d RequiresReplace plan modifiers", field.TerraformName, field.Replace, requiresReplaceCount)
	}
}

func TestCountRequiresReplaceIdentifiesTheModifier(t *testing.T) {
	onlyUseState := []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
	if got := countRequiresReplace(onlyUseState); got != 0 {
		t.Fatalf("countRequiresReplace(UseStateForUnknown) = %d, want 0; a bare count would report 1", got)
	}
	replaceAndUseState := []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
		stringplanmodifier.UseStateForUnknown(),
	}
	if got := countRequiresReplace(replaceAndUseState); got != 1 {
		t.Fatalf("countRequiresReplace(RequiresReplace, UseStateForUnknown) = %d, want 1; a bare count would report 2", got)
	}
	if got := countRequiresReplace([]planmodifier.Bool{boolplanmodifier.RequiresReplace()}); got != 1 {
		t.Fatalf("countRequiresReplace(bool RequiresReplace) = %d, want 1", got)
	}
	if got := countRequiresReplace([]planmodifier.Int64(nil)); got != 0 {
		t.Fatalf("countRequiresReplace(nil) = %d, want 0", got)
	}
}

func TestGeneratedPoliciesMatchLegacyBehavior(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "legacy_field_policies.json"))
	if err != nil {
		t.Fatal(err)
	}
	var evidence map[string]map[string]map[string]string
	if err := json.Unmarshal(raw, &evidence); err != nil {
		t.Fatal(err)
	}
	registry := generatedRegistry(t)
	if len(evidence) != len(registry) {
		t.Errorf("evidence covers %d resources, the registry has %d", len(evidence), len(registry))
	}
	checked, excluded := 0, 0
	var walk func(fields []spec.FieldSpec, evidence map[string]map[string]string, prefix, resourceType, identity, collectionIdentity string)
	walk = func(fields []spec.FieldSpec, evidence map[string]map[string]string, prefix, resourceType, identity, collectionIdentity string) {
		paired := pairedFields(fields)
		for _, field := range fields {
			if field.Unmanaged {
				continue
			}
			path := prefix + field.TerraformName

			if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
				excluded++
				nestedIdentity := ""
				if field.Collection != nil {
					nestedIdentity = field.Collection.IdentityField
				}
				walk(field.Fields, evidence, path+".", resourceType, identity, nestedIdentity)
				continue
			}

			if collectionIdentity != "" && field.TerraformName == collectionIdentity {
				excluded++
				continue
			}
			policies, recorded := evidence[path]
			if !recorded {

				if paired[field.TerraformName] {
					t.Errorf("no legacy policy evidence for paired field %s.%s", resourceType, path)
					continue
				}
				t.Errorf("no legacy policy evidence for %s.%s", resourceType, path)
				continue
			}
			got := map[string]string{
				"update_clear":     string(field.UpdateClear),
				"create_null":      string(field.CreateNull),
				"response_absence": string(field.ResponseAbsence),
				"unknown_plan":     string(field.UnknownPlan),
			}
			for policy, want := range policies {
				if policy == "driver" {
					continue
				}

				if path == identity && policy != "response_absence" {
					continue
				}
				checked++
				if got[policy] != want {
					t.Errorf("%s.%s %s = %q, legacy behavior is %q",
						resourceType, path, policy, got[policy], want)
				}
			}
		}
	}
	for _, resourceSpec := range registry {
		fields, exists := evidence[resourceSpec.TerraformType]
		if !exists {
			t.Errorf("no legacy policy evidence for %s", resourceSpec.TerraformType)
			continue
		}
		walk(resourceSpec.Fields, fields, "", resourceSpec.TerraformType, resourceSpec.IdentityPath, "")
	}
	t.Logf("%d policy assertions checked, %d fields excluded by category", checked, excluded)
}

func pairedFields(fields []spec.FieldSpec) map[string]bool {
	paired := make(map[string]bool)
	for _, field := range fields {
		if field.Reference != nil {
			paired[field.TerraformName] = true
			paired[field.Reference.TypeField] = true
		}
		if field.AutoAssignment != nil {
			paired[field.TerraformName] = true
			paired[field.AutoAssignment.FlagField] = true
		}
	}
	return paired
}

func TestGeneratedPairsAreModelled(t *testing.T) {
	references, assignments := 0, 0
	var walk func(fields []spec.FieldSpec, path string)
	walk = func(fields []spec.FieldSpec, path string) {
		byName := make(map[string]spec.FieldSpec, len(fields))
		for _, field := range fields {
			byName[field.TerraformName] = field
		}
		for _, field := range fields {
			if base, found := strings.CutSuffix(field.TerraformName, "_ref_type_"); found {
				partner, exists := byName[base]
				if !exists {
					t.Errorf("%s.%s has no base field", path, field.TerraformName)
				} else if partner.Reference == nil {
					t.Errorf("%s.%s is a reference companion but %s records no reference", path, field.TerraformName, base)
				} else {
					if partner.Reference.TypeField != field.TerraformName {
						t.Errorf("%s.%s references type field %q, want %q", path, base, partner.Reference.TypeField, field.TerraformName)
					}
					if len(partner.Reference.AllowedTypes) == 0 {
						t.Errorf("%s.%s records no allowed reference types", path, base)
					}
					references++
				}
			}
			if base, found := strings.CutSuffix(field.TerraformName, "_auto_assigned_"); found {
				partner, exists := byName[base]
				if !exists {
					t.Errorf("%s.%s has no base field", path, field.TerraformName)
				} else if partner.AutoAssignment == nil {
					t.Errorf("%s.%s is an auto-assignment flag but %s records no auto assignment", path, field.TerraformName, base)
				} else {
					if partner.AutoAssignment.FlagField != field.TerraformName {
						t.Errorf("%s.%s names flag %q, want %q", path, base, partner.AutoAssignment.FlagField, field.TerraformName)
					}
					assignments++
				}
			}
			walk(field.Fields, path+"."+field.TerraformName)
		}
	}
	for _, resourceSpec := range generatedRegistry(t) {
		walk(resourceSpec.Fields, resourceSpec.TerraformType)
	}
	if references == 0 || assignments == 0 {
		t.Fatalf("found %d reference and %d auto-assignment pairs; the walk looks broken", references, assignments)
	}
	t.Logf("%d reference pairs and %d auto-assignment pairs modelled", references, assignments)
}

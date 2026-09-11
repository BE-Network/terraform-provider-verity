package provider

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
	// Constructors come from the provider's own registration list so the mapping
	// cannot drift: every generated resource must correspond to a shipped one.
	constructors := legacyConstructorsByType(t)
	// Cache keys pin the generated alias to the value each legacy resource uses
	// today. ACL derives its key from the ip_version it was constructed with.
	cacheKeys := map[string]string{
		"verity_acl_v4":                   NewVerityACLV4Resource().(*verityACLUnifiedResource).getCacheKey(),
		"verity_acl_v6":                   NewVerityACLV6Resource().(*verityACLUnifiedResource).getCacheKey(),
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

// legacyConstructorsByType indexes every resource the provider registers by its
// Terraform type name.
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

// assertGeneratedModes pins generated resource modes to the legacy compatibility
// map. Extraction derives modes from datacenter/campus OpenAPI presence, so this
// is the only check that ties that derivation to the behavior the provider ships.
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
	// Unmanaged fields record API shape the provider deliberately does not
	// surface, so they are excluded from the count the legacy schema must match.
	managed := make([]spec.FieldSpec, 0, len(resourceSpec.Fields))
	for _, field := range resourceSpec.Fields {
		if !field.Unmanaged {
			managed = append(managed, field)
		}
	}
	if len(response.Schema.Attributes)+len(response.Schema.Blocks) != len(managed) {
		t.Fatalf("legacy top-level schema has %d attributes and %d blocks for %d managed generated fields", len(response.Schema.Attributes), len(response.Schema.Blocks), len(managed))
	}
	for _, field := range managed {
		// OpenAPI models a singleton like object_properties as an object, while
		// every legacy resource models it as a single-entry ListNestedBlock. Both
		// generated kinds therefore resolve to a legacy block, and the collection
		// strategy carries the cardinality difference.
		if field.Kind != spec.FieldKindList && field.Kind != spec.FieldKindObject {
			attribute, exists := response.Schema.Attributes[field.TerraformName]
			if !exists {
				t.Fatalf("legacy schema is missing generated field %q", field.TerraformName)
			}
			assertGeneratedAttribute(t, field, attribute)
			continue
		}
		if field.Kind == spec.FieldKindObject && field.Collection.Strategy != spec.CollectionSingleton {
			t.Fatalf("generated object field %q must use the singleton collection strategy, got %q", field.TerraformName, field.Collection.Strategy)
		}
		block, exists := response.Schema.Blocks[field.TerraformName]
		if !exists {
			t.Fatalf("legacy schema is missing generated block %q", field.TerraformName)
		}
		list, ok := block.(schema.ListNestedBlock)
		if !ok || list.Description != field.Description || len(list.NestedObject.Attributes) != len(field.Fields) {
			t.Fatalf("legacy block %q does not match generated metadata", field.TerraformName)
		}
		for _, nested := range field.Fields {
			attribute, exists := list.NestedObject.Attributes[nested.TerraformName]
			if !exists {
				t.Fatalf("legacy block %q is missing generated field %q", field.TerraformName, nested.TerraformName)
			}
			assertGeneratedAttribute(t, nested, attribute)
		}
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

// requiresReplaceTypes holds the concrete type of each kind's RequiresReplace
// plan modifier. Counting modifiers instead would equate "has one modifier" with
// "replaces on change", so an unrelated modifier would read as replacement and a
// genuine one paired with another would read as none.
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

// countRequiresReplace must identify the modifier rather than count modifiers.
// Both cases below are ones a bare count gets wrong: an unrelated modifier is not
// replacement, and replacement paired with another modifier still is.
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

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func TestCanonicalJSON(t *testing.T) {
	got, err := canonicalJSON([]byte(`{"z":1.20,"a":"<value>","nested":{"b":false,"a":null}}`))
	if err != nil {
		t.Fatalf("canonicalJSON() error = %v", err)
	}
	want := "{\n  \"a\": \"<value>\",\n  \"nested\": {\n    \"a\": null,\n    \"b\": false\n  },\n  \"z\": 1.20\n}\n"
	if string(got) != want {
		t.Fatalf("canonicalJSON() = %q, want %q", got, want)
	}
}

func TestNormalizeAndVerify(t *testing.T) {
	tempDir := t.TempDir()
	datacenter := filepath.Join(tempDir, "raw-datacenter.json")
	campus := filepath.Join(tempDir, "raw-campus.json")
	if err := os.WriteFile(datacenter, []byte(`{"z":true,"a":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(campus, []byte(`{"b":null,"a":"campus"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := filepath.Join(tempDir, "canonical")
	opts := normalizeOptions{Version: "6.6", Datacenter: datacenter, Campus: campus, OutputDir: outputDir, SourceExportDate: "2026-08-11", Provenance: "test fixture"}
	if err := normalize(opts); err != nil {
		t.Fatalf("normalize() error = %v", err)
	}
	if err := verify(outputDir); err != nil {
		t.Fatalf("verify() error = %v", err)
	}

	first, err := os.ReadFile(filepath.Join(outputDir, "datacenter.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := normalize(opts); err != nil {
		t.Fatalf("second normalize() error = %v", err)
	}
	second, err := os.ReadFile(filepath.Join(outputDir, "datacenter.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("normalization output changed for identical input")
	}

	if err := os.WriteFile(filepath.Join(outputDir, "campus.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verify(outputDir); err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("verify() after mutation error = %v, want checksum mismatch", err)
	}
}

func TestExtractCoverageMergesModeSpecificOpenAPIWithoutGuessingPolicies(t *testing.T) {
	tempDir := t.TempDir()
	datacenter := filepath.Join(tempDir, "datacenter-raw.json")
	campus := filepath.Join(tempDir, "campus-raw.json")
	datacenterDocument := `{"openapi":"3.0.0","paths":{"/widgets":{"put":{"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"widget":{"type":"object","properties":{"widget_name":{"type":"object","properties":{"name":{"type":"string"},"enabled":{"type":"boolean"}}}}}}}}}}}}}}`
	campusDocument := `{"openapi":"3.0.0","paths":{"/widgets":{"patch":{"requestBody":{"content":{"application/json":{"schema":{"type":"object","properties":{"widget":{"type":"object","properties":{"widget_name":{"type":"object","properties":{"name":{"type":"string"},"campus_only":{"type":"integer","nullable":true}}}}}}}}}}}}}}`
	if err := os.WriteFile(datacenter, []byte(datacenterDocument), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(campus, []byte(campusDocument), 0o644); err != nil {
		t.Fatal(err)
	}
	inputs := filepath.Join(tempDir, "inputs")
	if err := normalize(normalizeOptions{Version: "6.6", Datacenter: datacenter, Campus: campus, OutputDir: inputs, SourceExportDate: "2026-09-03", Provenance: "test"}); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(tempDir, "coverage.json")
	if err := extract(extractOptions{InputDir: inputs, Output: output}); err != nil {
		t.Fatal(err)
	}
	if err := extract(extractOptions{InputDir: inputs, Output: output, Check: true}); err != nil {
		t.Fatal(err)
	}
	var report coverageManifest
	raw, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Resources) != 1 {
		t.Fatalf("resources = %#v", report.Resources)
	}
	resource := report.Resources[0]
	if resource.RequestWrapperKey != "widget" || strings.Join(resource.Modes, ",") != "campus,datacenter" {
		t.Fatalf("resource = %#v", resource)
	}
	if len(resource.Fields) != 3 || resource.Fields[0].APIName != "campus_only" || !resource.Fields[0].Nullable {
		t.Fatalf("fields = %#v", resource.Fields)
	}
	if !strings.Contains(strings.Join(resource.Fields[0].Review, " "), "lifecycle") {
		t.Fatalf("field policy review was not reported: %#v", resource.Fields[0])
	}
}

func TestExtractionCapturesGetResponseAndDeleteParameterShapes(t *testing.T) {
	resource := &coverageResource{}
	pathItem := map[string]any{
		"get": map[string]any{
			"responses": map[string]any{
				"200": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"properties": map[string]any{"widgets": map[string]any{"type": "object"}},
							},
						},
					},
				},
			},
		},
		"delete": map[string]any{
			"parameters": []any{map[string]any{
				"in": "query", "name": "widget_name", "required": true,
				"schema": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			}},
		},
	}
	operations := resourceOperations(pathItem)
	if strings.Join(operations, ",") != "get,delete" {
		t.Fatalf("operations = %#v", operations)
	}
	mergeResponseShape(resource, "datacenter", pathItem["get"])
	mergeDeleteParameters(resource, "datacenter", pathItem["delete"])
	if resource.ResponseCollectionKey != "widgets" {
		t.Fatalf("response key = %q", resource.ResponseCollectionKey)
	}
	if len(resource.DeleteParameters) != 1 {
		t.Fatalf("delete parameters = %#v", resource.DeleteParameters)
	}
	parameter := resource.DeleteParameters[0]
	if parameter.Name != "widget_name" || !parameter.Required || parameter.Kind != "array" || parameter.ItemKind != "string" {
		t.Fatalf("delete parameter = %#v", parameter)
	}
}

func TestGeneratedIPv4RegistryIsValid(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	output := filepath.Join(t.TempDir(), "registry.json")
	if err := generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: filepath.Join(repositoryRoot, "specs", "overrides.yaml"),
		Output:    output,
	}); err != nil {
		t.Fatalf("generateRegistry() error = %v", err)
	}
	if err := generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: filepath.Join(repositoryRoot, "specs", "overrides.yaml"),
		Output:    output,
		Check:     true,
	}); err != nil {
		t.Fatalf("generateRegistry(check) error = %v", err)
	}
	raw, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var artifact registryArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatal(err)
	}
	if artifact.APIVersion != "6.6" || len(artifact.Resources) == 0 || len(artifact.UnrepresentedResources) == 0 {
		t.Fatalf("artifact coverage = %#v", artifact)
	}
	if err := artifact.Resources.Validate(); err != nil {
		t.Fatalf("generated registry validation failed: %v", err)
	}
}

func TestRegistryRejectsStaleOverrideFieldType(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(repositoryRoot, "specs", "overrides.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	stale := strings.Replace(string(raw), "api_kind: boolean", "api_kind: string", 1)
	overrides := filepath.Join(t.TempDir(), "overrides.yaml")
	if err := os.WriteFile(overrides, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	err = generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: overrides,
		Output:    filepath.Join(t.TempDir(), "registry.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "does not match extracted OpenAPI type") {
		t.Fatalf("generateRegistry() error = %v, want stale type rejection", err)
	}
}

func TestOverridesRejectUnknownFields(t *testing.T) {
	overrides := filepath.Join(t.TempDir(), "overrides.yaml")
	if err := os.WriteFile(overrides, []byte("format_version: 1\napi_version: '6.6'\nunknown: true\nresources: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readOverrides(overrides); err == nil || !strings.Contains(err.Error(), "field unknown not found") {
		t.Fatalf("readOverrides() error = %v, want unknown field rejection", err)
	}
}

func TestCoverageFieldFromSchemaPreservesNestedArrayFields(t *testing.T) {
	field := coverageFieldFromSchema("members", map[string]any{
		"type": "array",
		"items": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
	}, "datacenter")
	if field.Kind != "array" || field.ItemKind != "object" || len(field.Fields) != 1 || field.Fields[0].APIName != "name" || field.Fields[0].Kind != "string" {
		t.Fatalf("nested coverage field = %#v", field)
	}
}

func TestMergeFieldOverridesSupportsNestedCollections(t *testing.T) {
	versions := versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"}
	overrideField := func(apiName, apiKind, terraformName string) fieldOverride {
		return fieldOverride{
			APIName: apiName, APIKind: apiKind, TerraformName: terraformName, Description: terraformName,
			Modes: []string{"datacenter"}, Versions: versions, Access: "required",
			ResponseAbsence: "error", CreateNull: "reject", UpdateClear: "reject", UnknownPlan: "reject", StateOwnership: "configuration",
		}
	}
	child := overrideField("name", "string", "name")
	parent := overrideField("members", "array", "members")
	parent.Access = "optional"
	parent.APIItemKind = "object"
	parent.ResponseAbsence, parent.CreateNull, parent.UpdateClear, parent.UnknownPlan = "terraform_null", "omit", "api_null", "reject"
	parent.Collection = &collectionOverride{Strategy: "replace", Ordering: "ordered", IdentityField: "name"}
	parent.Fields = []fieldOverride{child}
	fields, err := mergeFieldOverrides([]coverageField{{
		APIName: "members", Kind: "array", ItemKind: "object", Description: "Members", Modes: []string{"datacenter"},
		Fields: []coverageField{{APIName: "name", Kind: "string", Description: "Name", Modes: []string{"datacenter"}}},
	}}, []fieldOverride{parent}, spec.VersionRange{MinInclusive: spec.APIVersion{Major: 6, Minor: 6}, MaxExclusive: spec.APIVersion{Major: 6, Minor: 7}}, spec.APIVersion{Major: 6, Minor: 6}, nil, "", nil, versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"})
	if err != nil {
		t.Fatalf("mergeFieldOverrides() error = %v", err)
	}
	resource := spec.ResourceSpec{
		TerraformType: "verity_members", Description: "Members", Modes: []spec.Mode{spec.ModeDatacenter},
		Versions:     spec.VersionRange{MinInclusive: spec.APIVersion{Major: 6, Minor: 6}, MaxExclusive: spec.APIVersion{Major: 6, Minor: 7}},
		IdentityPath: "members", API: spec.APIResourceSpec{EndpointPath: "/members", BulkKey: "members", RequestWrapperKey: "members", ResponseCollectionKey: "members", DeleteParameter: "member_name", CacheKey: "members"},
		Operations: spec.OperationSpec{Read: true, Create: true, Delete: true}, Fields: fields,
	}
	if err := resource.Validate(); err != nil {
		t.Fatalf("nested generated resource validation failed: %v", err)
	}
	if len(fields[0].Fields) != 1 || fields[0].Collection == nil || fields[0].Collection.IdentityField != "name" {
		t.Fatalf("nested generated fields = %#v", fields)
	}
	parent.Versions = versionRangeOverride{MinInclusive: "6.4", MaxExclusive: "6.7"}
	parent.Fields[0].Versions = versionRangeOverride{MinInclusive: "6.4", MaxExclusive: "6.5"}
	_, err = mergeFieldOverrides([]coverageField{{
		APIName: "members", Kind: "array", ItemKind: "object", Description: "Members", Modes: []string{"datacenter"},
		Fields: []coverageField{{APIName: "name", Kind: "string", Description: "Name", Modes: []string{"datacenter"}}},
	}}, []fieldOverride{parent}, spec.VersionRange{MinInclusive: spec.APIVersion{Major: 6, Minor: 4}, MaxExclusive: spec.APIVersion{Major: 6, Minor: 7}}, spec.APIVersion{Major: 6, Minor: 6}, nil, "", nil, versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"})
	if err == nil || !strings.Contains(err.Error(), "do not include selected API version") {
		t.Fatalf("mergeFieldOverrides() error = %v, want selected API version rejection", err)
	}
}

func TestMergeFieldOverridesSupportsScalarLists(t *testing.T) {
	versions := versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"}
	override := fieldOverride{
		APIName: "names", APIKind: "array", APIItemKind: "string", TerraformName: "names", Description: "Names",
		Modes: []string{"datacenter"}, Versions: versions, Access: "optional",
		ResponseAbsence: "terraform_null", CreateNull: "omit", UpdateClear: "api_null", UnknownPlan: "reject", StateOwnership: "configuration",
		Collection: &collectionOverride{Strategy: "replace", Ordering: "ordered"},
	}
	fields, err := mergeFieldOverrides([]coverageField{{APIName: "names", Kind: "array", ItemKind: "string", Description: "Names", Modes: []string{"datacenter"}}}, []fieldOverride{override}, spec.VersionRange{MinInclusive: spec.APIVersion{Major: 6, Minor: 6}, MaxExclusive: spec.APIVersion{Major: 6, Minor: 7}}, spec.APIVersion{Major: 6, Minor: 6}, nil, "", nil, versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"})
	if err != nil {
		t.Fatalf("mergeFieldOverrides() error = %v", err)
	}
	resource := spec.ResourceSpec{
		TerraformType: "verity_names", Description: "Names", Modes: []spec.Mode{spec.ModeDatacenter},
		Versions:     spec.VersionRange{MinInclusive: spec.APIVersion{Major: 6, Minor: 6}, MaxExclusive: spec.APIVersion{Major: 6, Minor: 7}},
		IdentityPath: "names", API: spec.APIResourceSpec{EndpointPath: "/names", BulkKey: "names", RequestWrapperKey: "names", ResponseCollectionKey: "names", DeleteParameter: "name", CacheKey: "names"},
		Operations: spec.OperationSpec{Read: true, Create: true, Delete: true}, Fields: fields,
	}
	if err := resource.Validate(); err != nil {
		t.Fatalf("scalar-list generated resource validation failed: %v", err)
	}
	if fields[0].ElementKind != spec.FieldKindString || fields[0].Collection.IdentityField != "" {
		t.Fatalf("scalar-list generated field = %#v", fields[0])
	}
}

// The merge's core safety property is that a reviewed override must account for
// every extracted field, so a field added to the API cannot enter the registry
// without review.
func TestRegistryRejectsUnreviewedExtractedField(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(repositoryRoot, "specs", "overrides.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(raw), "\n")
	start := -1
	for index, line := range lines {
		if strings.TrimSpace(line) == "- api_name: ipv4_list" {
			start = index
			break
		}
	}
	if start < 0 {
		t.Fatal("overrides no longer contain the ipv4_list field override")
	}
	end := start + 1
	for end < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[end]), "- api_name:") && !strings.HasPrefix(lines[end], "  - path:") {
		end++
	}
	trimmed := strings.Join(append(append([]string{}, lines[:start]...), lines[end:]...), "\n")
	overrides := filepath.Join(t.TempDir(), "overrides.yaml")
	if err := os.WriteFile(overrides, []byte(trimmed), 0o644); err != nil {
		t.Fatal(err)
	}
	err = generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: overrides,
		Output:    filepath.Join(t.TempDir(), "registry.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "every extracted field needs an explicit override") {
		t.Fatalf("generateRegistry() error = %v, want unreviewed field rejection", err)
	}
}

func TestRegistryRejectsModeMismatch(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(repositoryRoot, "specs", "overrides.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	// verity_ipv4_list is datacenter-only in both OpenAPI documents, so claiming
	// campus support must fail rather than silently widen the resource.
	widened := strings.Replace(string(raw), "modes: [datacenter]", "modes: [campus, datacenter]", 1)
	if widened == string(raw) {
		t.Fatal("overrides no longer contain a datacenter-only resource to widen")
	}
	overrides := filepath.Join(t.TempDir(), "overrides.yaml")
	if err := os.WriteFile(overrides, []byte(widened), 0o644); err != nil {
		t.Fatal(err)
	}
	err = generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: overrides,
		Output:    filepath.Join(t.TempDir(), "registry.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "modes do not match extracted OpenAPI modes") {
		t.Fatalf("generateRegistry() error = %v, want mode mismatch rejection", err)
	}
}

func TestValidateDeleteParameterRejectsOptionalAndUnsupported(t *testing.T) {
	parameters := []coverageQueryParameter{
		{Name: "changeset_name", Kind: "string", Required: false},
		{Name: "ipv4_list_filter_name", Kind: "array", Required: true},
	}
	deletes := spec.OperationSpec{Read: true, Create: true, Delete: true}
	if err := validateDeleteParameter(parameters, "ipv4_list_filter_name", deletes, nil); err != nil {
		t.Fatalf("validateDeleteParameter() error = %v, want the required parameter accepted", err)
	}
	err := validateDeleteParameter(parameters, "changeset_name", deletes, nil)
	if err == nil || !strings.Contains(err.Error(), "optional query parameter") {
		t.Fatalf("validateDeleteParameter() error = %v, want optional parameter rejection", err)
	}
	err = validateDeleteParameter(parameters, "missing_name", deletes, nil)
	if err == nil || !strings.Contains(err.Error(), "is not an extracted delete parameter") {
		t.Fatalf("validateDeleteParameter() error = %v, want unknown parameter rejection", err)
	}
	err = validateDeleteParameter(nil, "", deletes, nil)
	if err == nil || !strings.Contains(err.Error(), "delete_parameter is required") {
		t.Fatalf("validateDeleteParameter() error = %v, want missing parameter rejection", err)
	}
	noDelete := spec.OperationSpec{Read: true, Update: true}
	if err := validateDeleteParameter(nil, "", noDelete, nil); err != nil {
		t.Fatalf("validateDeleteParameter() error = %v, want empty parameter accepted without delete", err)
	}
	err = validateDeleteParameter(parameters, "ipv4_list_filter_name", noDelete, nil)
	if err == nil || !strings.Contains(err.Error(), "does not support delete") {
		t.Fatalf("validateDeleteParameter() error = %v, want unsupported-delete rejection", err)
	}
}

func TestRegistryRejectsUnknownProfile(t *testing.T) {
	repositoryRoot := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(repositoryRoot, "specs", "overrides.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(raw), "profile: identity", "profile: not_a_profile", 1)
	if broken == string(raw) {
		t.Fatal("overrides no longer reference the identity profile")
	}
	overrides := filepath.Join(t.TempDir(), "overrides.yaml")
	if err := os.WriteFile(overrides, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	err = generateRegistry(registryOptions{
		InputDir:  filepath.Join(repositoryRoot, "specs", "openapi", "6.6"),
		Overrides: overrides,
		Output:    filepath.Join(t.TempDir(), "registry.json"),
	})
	if err == nil || !strings.Contains(err.Error(), "unknown profile") {
		t.Fatalf("generateRegistry() error = %v, want unknown profile rejection", err)
	}
}

// A profile supplies policies but must never mask an explicit decision, and a
// field that states a policy itself has to win over the profile it references.
func TestResolveFieldOverridePrefersExplicitValues(t *testing.T) {
	profiles := map[string]policyProfile{"server_managed": {
		Access: "optional_computed", ResponseAbsence: "terraform_null", CreateNull: "omit",
		UpdateClear: "api_null", UnknownPlan: "preserve_state", StateOwnership: "configuration_or_server",
	}}
	inheritModes := []string{"campus", "datacenter"}
	inheritVersions := versionRangeOverride{MinInclusive: "6.6", MaxExclusive: "6.7"}

	resolved, err := resolveFieldOverride(fieldOverride{APIName: "enable", Profile: "server_managed"}, profiles, "", inheritModes, inheritVersions)
	if err != nil {
		t.Fatalf("resolveFieldOverride() error = %v", err)
	}
	if resolved.TerraformName != "enable" || resolved.UpdateClear != "api_null" ||
		len(resolved.Modes) != 2 || resolved.Versions != inheritVersions {
		t.Fatalf("inherited field = %#v", resolved)
	}

	explicit := fieldOverride{APIName: "enable", Profile: "server_managed", UpdateClear: "reject", Modes: []string{"campus"}}
	resolved, err = resolveFieldOverride(explicit, profiles, "", inheritModes, inheritVersions)
	if err != nil {
		t.Fatalf("resolveFieldOverride() error = %v", err)
	}
	if resolved.UpdateClear != "reject" || len(resolved.Modes) != 1 || resolved.Modes[0] != "campus" {
		t.Fatalf("explicit field = %#v", resolved)
	}

	if _, err := resolveFieldOverride(fieldOverride{APIName: "x", Profile: "server_managed", Unmanaged: true}, profiles, "", inheritModes, inheritVersions); err == nil {
		t.Fatal("resolveFieldOverride() accepted an unmanaged field with a profile")
	}
}

// The default profile must apply only where a field says nothing about policy.
// A field that names a different profile, or states a policy outright, keeps it.
func TestResolveFieldOverrideAppliesDefaultProfile(t *testing.T) {
	profiles := map[string]policyProfile{
		"server_managed": {Access: "optional_computed", UpdateClear: "api_null"},
		"identity":       {Access: "required", Replace: true, UpdateClear: "reject"},
	}
	resolved, err := resolveFieldOverride(fieldOverride{APIName: "enable"}, profiles, "server_managed", nil, versionRangeOverride{})
	if err != nil {
		t.Fatalf("resolveFieldOverride() error = %v", err)
	}
	if resolved.Profile != "server_managed" || resolved.Access != "optional_computed" {
		t.Fatalf("defaulted field = %#v", resolved)
	}
	resolved, err = resolveFieldOverride(fieldOverride{APIName: "name", Profile: "identity"}, profiles, "server_managed", nil, versionRangeOverride{})
	if err != nil {
		t.Fatalf("resolveFieldOverride() error = %v", err)
	}
	if resolved.Access != "required" || !resolved.Replace {
		t.Fatalf("named profile lost to the default: %#v", resolved)
	}
	// An unmanaged field records API shape only and must not pick up the default.
	resolved, err = resolveFieldOverride(fieldOverride{APIName: "object_properties", Unmanaged: true}, profiles, "server_managed", nil, versionRangeOverride{})
	if err != nil {
		t.Fatalf("resolveFieldOverride() error = %v", err)
	}
	if resolved.Profile != "" || resolved.Access != "" {
		t.Fatalf("unmanaged field picked up the default profile: %#v", resolved)
	}
}

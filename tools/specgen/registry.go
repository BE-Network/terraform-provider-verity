package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"terraform-provider-verity/internal/spec"
)

const registryFormatVersion = 1

type registryOptions struct {
	InputDir  string
	Overrides string
	Output    string
	Check     bool
}

type overridesDocument struct {
	FormatVersion int                `yaml:"format_version"`
	APIVersion    string             `yaml:"api_version"`
	Defaults      defaultsOverride   `yaml:"defaults"`
	Resources     []resourceOverride `yaml:"resources"`
}

// defaultsOverride carries the values shared by most of the registry. Repeating
// them per field made the reviewed file overwhelmingly boilerplate, which hid
// the entries that actually differ; naming them here keeps every decision
// explicit while leaving only genuine deviations in each resource.
type defaultsOverride struct {
	Versions versionRangeOverride     `yaml:"versions"`
	Profile  string                   `yaml:"profile"`
	Profiles map[string]policyProfile `yaml:"profiles"`
}

// policyProfile is a reviewed bundle of lifecycle policies referenced by name.
type policyProfile struct {
	Access          string `yaml:"access"`
	Replace         bool   `yaml:"replace"`
	ResponseAbsence string `yaml:"response_absence"`
	CreateNull      string `yaml:"create_null"`
	UpdateClear     string `yaml:"update_clear"`
	UnknownPlan     string `yaml:"unknown_plan"`
	StateOwnership  string `yaml:"state_ownership"`
}

type resourceOverride struct {
	Path          string               `yaml:"path"`
	TerraformType string               `yaml:"terraform_type"`
	Description   string               `yaml:"description"`
	Modes         []string             `yaml:"modes"`
	Versions      versionRangeOverride `yaml:"versions"`
	IdentityPath  string               `yaml:"identity_path"`
	SchemaVersion int64                `yaml:"schema_version"`
	API           apiOverride          `yaml:"api"`
	Fields        []fieldOverride      `yaml:"fields"`
	Dependencies  spec.DependencySpec  `yaml:"dependencies"`
	Hooks         []string             `yaml:"hooks"`
}

type apiOverride struct {
	BulkKey               string            `yaml:"bulk_key"`
	CacheKey              string            `yaml:"cache_key"`
	ResponseCollectionKey string            `yaml:"response_collection_key"`
	DeleteParameter       string            `yaml:"delete_parameter"`
	FixedHeaders          map[string]string `yaml:"fixed_headers"`
}

type fieldOverride struct {
	APIName         string                  `yaml:"api_name"`
	APIKind         string                  `yaml:"api_kind"`
	Profile         string                  `yaml:"profile"`
	Unmanaged       bool                    `yaml:"unmanaged"`
	APIItemKind     string                  `yaml:"api_item_kind"`
	TerraformName   string                  `yaml:"terraform_name"`
	Description     string                  `yaml:"description"`
	Modes           []string                `yaml:"modes"`
	Versions        versionRangeOverride    `yaml:"versions"`
	Access          string                  `yaml:"access"`
	Sensitive       bool                    `yaml:"sensitive"`
	Replace         bool                    `yaml:"replace"`
	ResponseAbsence string                  `yaml:"response_absence"`
	CreateNull      string                  `yaml:"create_null"`
	UpdateClear     string                  `yaml:"update_clear"`
	UnknownPlan     string                  `yaml:"unknown_plan"`
	StateOwnership  string                  `yaml:"state_ownership"`
	Collection      *collectionOverride     `yaml:"collection"`
	Fields          []fieldOverride         `yaml:"fields"`
	Reference       *referenceOverride      `yaml:"reference"`
	AutoAssignment  *autoAssignmentOverride `yaml:"auto_assignment"`
}

type versionRangeOverride struct {
	MinInclusive string `yaml:"min_inclusive"`
	MaxExclusive string `yaml:"max_exclusive"`
}

type collectionOverride struct {
	Strategy      string `yaml:"strategy"`
	Ordering      string `yaml:"ordering"`
	IdentityField string `yaml:"identity_field"`
}

type referenceOverride struct {
	TypeField    string   `yaml:"type_field"`
	AllowedTypes []string `yaml:"allowed_types"`
}

type autoAssignmentOverride struct {
	FlagField string `yaml:"flag_field"`
}

type registryArtifact struct {
	FormatVersion          int           `json:"format_version"`
	GeneratedBy            string        `json:"generated_by"`
	APIVersion             string        `json:"api_version"`
	Resources              spec.Registry `json:"resources"`
	OverriddenResources    []string      `json:"overridden_resources"`
	UnrepresentedResources []string      `json:"unrepresented_resources"`
}

func generateRegistry(opts registryOptions) error {
	if opts.InputDir == "" || opts.Overrides == "" || opts.Output == "" {
		return errors.New("--input-dir, --overrides, and --output are required")
	}
	if err := verify(opts.InputDir); err != nil {
		return fmt.Errorf("verify canonical inputs: %w", err)
	}
	report, err := extractCoverage(opts.InputDir)
	if err != nil {
		return err
	}
	overrides, err := readOverrides(opts.Overrides)
	if err != nil {
		return err
	}
	if overrides.FormatVersion != 1 {
		return fmt.Errorf("overrides format_version must be 1, got %d", overrides.FormatVersion)
	}
	if overrides.APIVersion != report.APIVersion {
		return fmt.Errorf("overrides api_version %q does not match input API version %q", overrides.APIVersion, report.APIVersion)
	}
	version, err := parseAPIVersion(report.APIVersion)
	if err != nil {
		return fmt.Errorf("parse input API version: %w", err)
	}
	coverageByPath := make(map[string]coverageResource, len(report.Resources))
	for _, resource := range report.Resources {
		coverageByPath[resource.Path] = resource
	}
	// A single endpoint can back more than one Terraform resource when a required
	// discriminator selects between them, as /acls does with ip_version. Uniqueness
	// therefore belongs to the Terraform type; paths may legitimately repeat.
	overridden := make(map[string]bool, len(overrides.Resources))
	byType := make(map[string]bool, len(overrides.Resources))
	resources := make(spec.Registry, 0, len(overrides.Resources))
	for _, override := range overrides.Resources {
		if byType[override.TerraformType] {
			return fmt.Errorf("overrides contain duplicate terraform_type %q", override.TerraformType)
		}
		coverage, ok := coverageByPath[override.Path]
		if !ok {
			return fmt.Errorf("override target %q is not present in extracted OpenAPI", override.Path)
		}
		resource, err := mergeResourceOverride(coverage, override, version, overrides.Defaults)
		if err != nil {
			return fmt.Errorf("override %q: %w", override.Path, err)
		}
		resources = append(resources, resource)
		overridden[override.Path] = true
		byType[override.TerraformType] = true
	}
	if err := resources.Validate(); err != nil {
		return fmt.Errorf("validate generated registry: %w", err)
	}
	sort.Slice(resources, func(i, j int) bool { return resources[i].TerraformType < resources[j].TerraformType })
	artifact := registryArtifact{
		FormatVersion: registryFormatVersion,
		GeneratedBy:   "tools/specgen registry",
		APIVersion:    report.APIVersion,
		Resources:     resources,
	}
	for _, resource := range report.Resources {
		if overridden[resource.Path] {
			artifact.OverriddenResources = append(artifact.OverriddenResources, resource.Path)
		} else {
			artifact.UnrepresentedResources = append(artifact.UnrepresentedResources, resource.Path)
		}
	}
	sort.Strings(artifact.OverriddenResources)
	sort.Strings(artifact.UnrepresentedResources)
	encoded, err := marshalCanonical(artifact)
	if err != nil {
		return fmt.Errorf("encode generated registry: %w", err)
	}
	if opts.Check {
		existing, err := os.ReadFile(opts.Output)
		if err != nil {
			return fmt.Errorf("read generated registry for check: %w", err)
		}
		if !bytes.Equal(existing, encoded) {
			return fmt.Errorf("generated registry differs: rerun specgen registry --input-dir %s --overrides %s --output %s", opts.InputDir, opts.Overrides, opts.Output)
		}
		return nil
	}
	return writeFile(opts.Output, encoded)
}

func readOverrides(path string) (overridesDocument, error) {
	file, err := os.Open(path)
	if err != nil {
		return overridesDocument{}, fmt.Errorf("read overrides: %w", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	var overrides overridesDocument
	if err := decoder.Decode(&overrides); err != nil {
		return overridesDocument{}, fmt.Errorf("decode overrides: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return overridesDocument{}, errors.New("overrides must contain exactly one YAML document")
		}
		return overridesDocument{}, fmt.Errorf("decode overrides: %w", err)
	}
	return overrides, nil
}

func mergeResourceOverride(coverage coverageResource, override resourceOverride, version spec.APIVersion, defaults defaultsOverride) (spec.ResourceSpec, error) {
	if override.TerraformType == "" || override.Description == "" || override.IdentityPath == "" {
		return spec.ResourceSpec{}, errors.New("terraform_type, description, and identity_path are required")
	}
	if override.API.BulkKey == "" || override.API.CacheKey == "" || override.API.ResponseCollectionKey == "" {
		return spec.ResourceSpec{}, errors.New("api bulk_key, cache_key, and response_collection_key are required")
	}
	if coverage.RequestWrapperKey == "" {
		return spec.ResourceSpec{}, errors.New("OpenAPI request_wrapper_key is unresolved")
	}
	if coverage.ResponseCollectionKey != "" && coverage.ResponseCollectionKey != override.API.ResponseCollectionKey {
		return spec.ResourceSpec{}, fmt.Errorf("response_collection_key %q does not match OpenAPI value %q", override.API.ResponseCollectionKey, coverage.ResponseCollectionKey)
	}
	operations := operationsFromCoverage(coverage.Operations)
	if err := validateFixedHeaders(coverage.DeleteParameters, override.API.FixedHeaders); err != nil {
		return spec.ResourceSpec{}, err
	}
	if err := validateDeleteParameter(coverage.DeleteParameters, override.API.DeleteParameter, operations, override.API.FixedHeaders); err != nil {
		return spec.ResourceSpec{}, err
	}
	modes, err := reviewedModes(override.Modes, coverage.Modes, "resource")
	if err != nil {
		return spec.ResourceSpec{}, err
	}
	// A resource inherits the reviewed default range unless it states its own.
	declaredVersions := override.Versions
	if declaredVersions.MinInclusive == "" && declaredVersions.MaxExclusive == "" {
		declaredVersions = defaults.Versions
	}
	versions, err := parseVersionRange(declaredVersions)
	if err != nil {
		return spec.ResourceSpec{}, fmt.Errorf("versions: %w", err)
	}
	if !versions.Contains(version) {
		return spec.ResourceSpec{}, fmt.Errorf("versions do not include selected API version %s", version)
	}
	fields, err := mergeFieldOverrides(coverage.Fields, override.Fields, versions, version, defaults.Profiles, defaults.Profile, override.Modes, declaredVersions)
	if err != nil {
		return spec.ResourceSpec{}, err
	}
	return spec.ResourceSpec{
		TerraformType: override.TerraformType,
		Description:   override.Description,
		Modes:         modes,
		Versions:      versions,
		IdentityPath:  override.IdentityPath,
		SchemaVersion: override.SchemaVersion,
		API: spec.APIResourceSpec{
			EndpointPath:          coverage.Path,
			BulkKey:               override.API.BulkKey,
			RequestWrapperKey:     coverage.RequestWrapperKey,
			ResponseCollectionKey: override.API.ResponseCollectionKey,
			DeleteParameter:       override.API.DeleteParameter,
			CacheKey:              override.API.CacheKey,
			FixedHeaders:          override.API.FixedHeaders,
		},
		Operations:   operations,
		Fields:       fields,
		Dependencies: override.Dependencies,
		Hooks:        override.Hooks,
	}, nil
}

func mergeFieldOverrides(coverage []coverageField, overrides []fieldOverride, parentVersions spec.VersionRange, selectedVersion spec.APIVersion, profiles map[string]policyProfile, defaultProfile string, inheritModes []string, inheritVersions versionRangeOverride) ([]spec.FieldSpec, error) {
	coverageByName := make(map[string]coverageField, len(coverage))
	for _, field := range coverage {
		coverageByName[field.APIName] = field
	}
	overrideByName := make(map[string]fieldOverride, len(overrides))
	for _, override := range overrides {
		if override.APIName == "" {
			return nil, errors.New("field api_name is required")
		}
		if _, exists := overrideByName[override.APIName]; exists {
			return nil, fmt.Errorf("duplicate field override %q", override.APIName)
		}
		if _, exists := coverageByName[override.APIName]; !exists {
			return nil, fmt.Errorf("field override %q is not present in extracted OpenAPI", override.APIName)
		}
		overrideByName[override.APIName] = override
	}
	if len(overrideByName) != len(coverageByName) {
		missing := make([]string, 0)
		for name := range coverageByName {
			if _, exists := overrideByName[name]; !exists {
				missing = append(missing, name)
			}
		}
		sort.Strings(missing)
		return nil, fmt.Errorf("every extracted field needs an explicit override; missing %s", strings.Join(missing, ", "))
	}
	fields := make([]spec.FieldSpec, 0, len(overrides))
	for _, source := range coverage {
		override, err := resolveFieldOverride(overrideByName[source.APIName], profiles, defaultProfile, inheritModes, inheritVersions)
		if err != nil {
			return nil, err
		}
		kind, err := fieldKindFromOpenAPI(source.Kind)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", source.APIName, err)
		}
		modesForField, err := reviewedModes(override.Modes, source.Modes, fmt.Sprintf("field %q", source.APIName))
		if err != nil {
			return nil, err
		}
		versionsForField, err := parseVersionRange(override.Versions)
		if err != nil {
			return nil, fmt.Errorf("field %q versions: %w", source.APIName, err)
		}
		if !versionRangeSubset(versionsForField, parentVersions) {
			return nil, fmt.Errorf("field %q versions must be contained by resource versions", source.APIName)
		}
		if !versionsForField.Contains(selectedVersion) {
			return nil, fmt.Errorf("field %q versions do not include selected API version %s", source.APIName, selectedVersion)
		}
		if override.Unmanaged {
			if err := rejectManagedOverrideFields(override); err != nil {
				return nil, err
			}
			if len(source.Fields) != 0 {
				return nil, fmt.Errorf("field %q has nested API fields and cannot be unmanaged", source.APIName)
			}
			fields = append(fields, spec.FieldSpec{
				APIName: source.APIName, Kind: kind, ElementKind: "", Nullable: source.Nullable,
				Modes: modesForField, Versions: versionsForField, Unmanaged: true,
			})
			continue
		}
		// A description is user-facing documentation. Most match the API text, so
		// the override states only the wording the provider deliberately changes.
		description := override.Description
		if description == "" {
			description = source.Description
		}
		if override.APIKind == "" || override.TerraformName == "" || description == "" {
			return nil, fmt.Errorf("field %q requires api_kind, terraform_name, and a description in either the override or OpenAPI", source.APIName)
		}
		if override.APIKind != source.Kind {
			return nil, fmt.Errorf("field %q api_kind %q does not match extracted OpenAPI type %q", source.APIName, override.APIKind, source.Kind)
		}
		elementKind, err := fieldElementKind(source, override)
		if err != nil {
			return nil, err
		}
		modes := modesForField
		versions := versionsForField
		fields = append(fields, spec.FieldSpec{
			TerraformName:   override.TerraformName,
			APIName:         source.APIName,
			Kind:            kind,
			ElementKind:     elementKind,
			Access:          spec.Access(override.Access),
			Description:     description,
			Nullable:        source.Nullable,
			Sensitive:       override.Sensitive,
			Replace:         override.Replace,
			Modes:           modes,
			Versions:        versions,
			ResponseAbsence: spec.ResponseAbsencePolicy(override.ResponseAbsence),
			CreateNull:      spec.CreateNullPolicy(override.CreateNull),
			UpdateClear:     spec.UpdateClearPolicy(override.UpdateClear),
			UnknownPlan:     spec.UnknownPlanPolicy(override.UnknownPlan),
			StateOwnership:  spec.StateOwnershipPolicy(override.StateOwnership),
			Collection:      collectionSpec(override.Collection),
			Reference:       referenceSpec(override.Reference),
			AutoAssignment:  autoAssignmentSpec(override.AutoAssignment),
		})
		if len(source.Fields) != 0 || len(override.Fields) != 0 {
			nested, err := mergeFieldOverrides(source.Fields, override.Fields, versions, selectedVersion, profiles, defaultProfile, override.Modes, override.Versions)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", source.APIName, err)
			}
			fields[len(fields)-1].Fields = nested
		}
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].TerraformName < fields[j].TerraformName })
	return fields, nil
}

func fieldElementKind(source coverageField, override fieldOverride) (spec.FieldKind, error) {
	if source.Kind != "array" {
		if override.APIItemKind != "" {
			return "", fmt.Errorf("field %q api_item_kind is only valid for arrays", source.APIName)
		}
		return "", nil
	}
	if override.APIItemKind == "" {
		return "", fmt.Errorf("field %q requires api_item_kind", source.APIName)
	}
	if override.APIItemKind != source.ItemKind {
		return "", fmt.Errorf("field %q api_item_kind %q does not match extracted OpenAPI item type %q", source.APIName, override.APIItemKind, source.ItemKind)
	}
	kind, err := fieldKindFromOpenAPI(source.ItemKind)
	if err != nil {
		return "", fmt.Errorf("field %q item type: %w", source.APIName, err)
	}
	if kind == spec.FieldKindList {
		return "", fmt.Errorf("field %q nested array items are unsupported", source.APIName)
	}
	return kind, nil
}

// rejectManagedOverrideFields keeps an unmanaged declaration minimal, so marking
// a field unmanaged cannot quietly carry Terraform behavior alongside it.
// resolveFieldOverride applies the reviewed defaults to one field: a named
// profile supplies lifecycle policies, the Terraform name follows the API name,
// and modes and versions inherit from the resource. Anything stated explicitly
// on the field wins, so a deviation is always visible in the file.
func resolveFieldOverride(override fieldOverride, profiles map[string]policyProfile, defaultProfile string, modes []string, versions versionRangeOverride) (fieldOverride, error) {
	resolved := override
	if resolved.Profile == "" && !resolved.Unmanaged && resolved.Access == "" {
		resolved.Profile = defaultProfile
	}
	if resolved.Profile != "" {
		profile, exists := profiles[resolved.Profile]
		if !exists {
			return fieldOverride{}, fmt.Errorf("field %q references unknown profile %q", override.APIName, resolved.Profile)
		}
		if resolved.Unmanaged {
			return fieldOverride{}, fmt.Errorf("field %q is unmanaged and cannot reference a profile", override.APIName)
		}
		if resolved.Access == "" {
			resolved.Access = profile.Access
		}
		if !resolved.Replace {
			resolved.Replace = profile.Replace
		}
		if resolved.ResponseAbsence == "" {
			resolved.ResponseAbsence = profile.ResponseAbsence
		}
		if resolved.CreateNull == "" {
			resolved.CreateNull = profile.CreateNull
		}
		if resolved.UpdateClear == "" {
			resolved.UpdateClear = profile.UpdateClear
		}
		if resolved.UnknownPlan == "" {
			resolved.UnknownPlan = profile.UnknownPlan
		}
		if resolved.StateOwnership == "" {
			resolved.StateOwnership = profile.StateOwnership
		}
	}
	if !resolved.Unmanaged && resolved.TerraformName == "" {
		resolved.TerraformName = resolved.APIName
	}
	if len(resolved.Modes) == 0 {
		resolved.Modes = modes
	}
	if resolved.Versions.MinInclusive == "" && resolved.Versions.MaxExclusive == "" {
		resolved.Versions = versions
	}
	return resolved, nil
}

func rejectManagedOverrideFields(override fieldOverride) error {
	if override.TerraformName != "" || override.Access != "" || override.Description != "" {
		return fmt.Errorf("field %q is unmanaged and cannot declare terraform_name, access, or description", override.APIName)
	}
	if override.ResponseAbsence != "" || override.CreateNull != "" || override.UpdateClear != "" ||
		override.UnknownPlan != "" || override.StateOwnership != "" {
		return fmt.Errorf("field %q is unmanaged and cannot declare lifecycle policies", override.APIName)
	}
	if override.Collection != nil || len(override.Fields) != 0 || override.Reference != nil || override.AutoAssignment != nil {
		return fmt.Errorf("field %q is unmanaged and cannot declare nested structure", override.APIName)
	}
	if override.APIItemKind != "" || override.Sensitive || override.Replace {
		return fmt.Errorf("field %q is unmanaged and cannot declare item kind, sensitivity, or replacement", override.APIName)
	}
	return nil
}

func collectionSpec(value *collectionOverride) *spec.CollectionSpec {
	if value == nil {
		return nil
	}
	return &spec.CollectionSpec{Strategy: spec.CollectionStrategy(value.Strategy), Ordering: spec.CollectionOrdering(value.Ordering), IdentityField: value.IdentityField}
}

func referenceSpec(value *referenceOverride) *spec.ReferenceSpec {
	if value == nil {
		return nil
	}
	return &spec.ReferenceSpec{TypeField: value.TypeField, AllowedTypes: value.AllowedTypes}
}

func autoAssignmentSpec(value *autoAssignmentOverride) *spec.AutoAssignmentSpec {
	if value == nil {
		return nil
	}
	return &spec.AutoAssignmentSpec{FlagField: value.FlagField}
}

// validateDeleteParameter ties delete_parameter to the endpoint's DELETE support
// and to a required query parameter. Optional parameters such as the API-wide
// changeset_name identify a transaction rather than the objects being deleted,
// so accepting one by name alone would silently produce a resource that deletes
// the wrong thing.
// validateFixedHeaders checks that every declared discriminator is a required
// parameter the API actually exposes, so a resource cannot pin a header the
// endpoint ignores.
func validateFixedHeaders(parameters []coverageQueryParameter, headers map[string]string) error {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if headers[name] == "" {
			return fmt.Errorf("fixed header %q requires a value", name)
		}
		found := false
		for _, parameter := range parameters {
			if parameter.Name != name {
				continue
			}
			if !parameter.Required {
				return fmt.Errorf("fixed header %q is an optional parameter and cannot discriminate resources", name)
			}
			found = true
		}
		if !found {
			return fmt.Errorf("fixed header %q is not an extracted parameter", name)
		}
	}
	return nil
}

func validateDeleteParameter(parameters []coverageQueryParameter, name string, operations spec.OperationSpec, headers map[string]string) error {
	if !operations.Delete {
		if name != "" {
			return fmt.Errorf("delete_parameter %q is set but the endpoint does not support delete", name)
		}
		return nil
	}
	if name == "" {
		return errors.New("delete_parameter is required because the endpoint supports delete")
	}
	// A discriminator is required but selects which resources a call addresses
	// rather than naming them, so it can never be the delete identity.
	if _, pinned := headers[name]; pinned {
		return fmt.Errorf("delete_parameter %q is a fixed header and discriminates resources rather than identifying them", name)
	}
	required := make([]string, 0, len(parameters))
	for _, parameter := range parameters {
		if parameter.Required {
			required = append(required, parameter.Name)
		}
		if parameter.Name == name {
			if !parameter.Required {
				return fmt.Errorf("delete_parameter %q is an optional query parameter and cannot identify deleted objects", name)
			}
			return nil
		}
	}
	sort.Strings(required)
	return fmt.Errorf("delete_parameter %q is not an extracted delete parameter; required parameters are %s", name, strings.Join(required, ", "))
}

func modesFromCoverage(values []string) []spec.Mode {
	modes := make([]spec.Mode, 0, len(values))
	for _, value := range values {
		modes = append(modes, spec.Mode(value))
	}
	sort.Slice(modes, func(i, j int) bool { return modes[i] < modes[j] })
	return modes
}

func reviewedModes(expected, observed []string, subject string) ([]spec.Mode, error) {
	if len(expected) == 0 {
		return nil, fmt.Errorf("%s modes are required", subject)
	}
	expectedModes := modesFromCoverage(expected)
	observedModes := modesFromCoverage(observed)
	if len(expectedModes) != len(observedModes) {
		return nil, fmt.Errorf("%s modes do not match extracted OpenAPI modes", subject)
	}
	for index := range expectedModes {
		if expectedModes[index] != observedModes[index] {
			return nil, fmt.Errorf("%s modes do not match extracted OpenAPI modes", subject)
		}
	}
	return observedModes, nil
}

func operationsFromCoverage(values []string) spec.OperationSpec {
	operations := spec.OperationSpec{}
	for _, value := range values {
		switch value {
		case "put":
			operations.Create = true
		case "get":
			operations.Read = true
		case "patch":
			operations.Update = true
		case "delete":
			operations.Delete = true
		}
	}
	return operations
}

func fieldKindFromOpenAPI(kind string) (spec.FieldKind, error) {
	switch kind {
	case "string":
		return spec.FieldKindString, nil
	case "boolean":
		return spec.FieldKindBool, nil
	case "integer":
		return spec.FieldKindInt64, nil
	case "number":
		return spec.FieldKindNumber, nil
	case "object":
		return spec.FieldKindObject, nil
	case "array":
		return spec.FieldKindList, nil
	default:
		return "", fmt.Errorf("unsupported OpenAPI type %q", kind)
	}
}

func parseAPIVersion(value string) (spec.APIVersion, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return spec.APIVersion{}, fmt.Errorf("expected major.minor, got %q", value)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return spec.APIVersion{}, fmt.Errorf("invalid major version %q", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return spec.APIVersion{}, fmt.Errorf("invalid minor version %q", parts[1])
	}
	return spec.APIVersion{Major: major, Minor: minor}, nil
}

func parseVersionRange(value versionRangeOverride) (spec.VersionRange, error) {
	if value.MinInclusive == "" || value.MaxExclusive == "" {
		return spec.VersionRange{}, errors.New("min_inclusive and max_exclusive are required")
	}
	minInclusive, err := parseAPIVersion(value.MinInclusive)
	if err != nil {
		return spec.VersionRange{}, fmt.Errorf("parse min_inclusive: %w", err)
	}
	maxExclusive, err := parseAPIVersion(value.MaxExclusive)
	if err != nil {
		return spec.VersionRange{}, fmt.Errorf("parse max_exclusive: %w", err)
	}
	if minInclusive.Compare(maxExclusive) >= 0 {
		return spec.VersionRange{}, errors.New("expected a non-empty inclusive-minimum/exclusive-maximum range")
	}
	return spec.VersionRange{MinInclusive: minInclusive, MaxExclusive: maxExclusive}, nil
}

func versionRangeSubset(value, container spec.VersionRange) bool {
	return value.MinInclusive.Compare(container.MinInclusive) >= 0 && value.MaxExclusive.Compare(container.MaxExclusive) <= 0
}

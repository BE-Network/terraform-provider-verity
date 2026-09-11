package spec

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var canonicalDecimal = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?$`)

func (r ResourceSpec) Validate() error {
	if !strings.HasPrefix(r.TerraformType, "verity_") {
		return fmt.Errorf("terraform type must start with verity_: %q", r.TerraformType)
	}
	if err := validateModes(r.Modes); err != nil {
		return fmt.Errorf("resource %s: %w", r.TerraformType, err)
	}
	if err := validateRange(r.Versions); err != nil {
		return fmt.Errorf("resource %s: versions: %w", r.TerraformType, err)
	}
	if r.IdentityPath == "" || strings.Contains(r.IdentityPath, ".") {
		return fmt.Errorf("resource %s: identity path must name one top-level field", r.TerraformType)
	}
	if r.SchemaVersion < 0 {
		return fmt.Errorf("resource %s: schema version cannot be negative", r.TerraformType)
	}
	if !r.Operations.Read || (!r.Operations.Create && !r.Operations.Update) {
		return fmt.Errorf("resource %s: operations must include read and at least one mutation", r.TerraformType)
	}
	if err := r.API.validate(r.Operations); err != nil {
		return fmt.Errorf("resource %s: api: %w", r.TerraformType, err)
	}
	if err := validateFields(r.Fields, r.TerraformType, r.Versions, r.Modes); err != nil {
		return err
	}
	for _, field := range r.Fields {
		if field.TerraformName == r.IdentityPath {
			return nil
		}
	}
	return fmt.Errorf("resource %s: identity path %q does not name a top-level field", r.TerraformType, r.IdentityPath)
}

func (a APIResourceSpec) validate(operations OperationSpec) error {
	for label, value := range map[string]string{
		"endpoint path": a.EndpointPath, "bulk key": a.BulkKey, "request wrapper key": a.RequestWrapperKey,
		"response collection key": a.ResponseCollectionKey, "cache key": a.CacheKey,
	} {
		if value == "" {
			return fmt.Errorf("%s is required", label)
		}
	}
	// A delete parameter identifies the objects a DELETE removes, so it is
	// required exactly when the endpoint supports delete and must be absent
	// otherwise rather than carrying an unused reviewed value.
	if operations.Delete && a.DeleteParameter == "" {
		return fmt.Errorf("delete parameter is required when the endpoint supports delete")
	}
	if !operations.Delete && a.DeleteParameter != "" {
		return fmt.Errorf("delete parameter %q is set but the endpoint does not support delete", a.DeleteParameter)
	}
	if !strings.HasPrefix(a.EndpointPath, "/") {
		return fmt.Errorf("endpoint path must start with '/'")
	}
	return nil
}

func validateFields(fields []FieldSpec, path string, parentRange VersionRange, parentModes []Mode) error {
	if len(fields) == 0 {
		return fmt.Errorf("%s: fields are required", path)
	}
	byName := make(map[string]FieldSpec, len(fields))
	for _, field := range fields {
		if field.Unmanaged {
			if err := validateUnmanagedField(field, path+"."+field.APIName, parentRange, parentModes); err != nil {
				return err
			}
			continue
		}
		fieldPath := path + "." + field.TerraformName
		if field.TerraformName == "" || field.APIName == "" {
			return fmt.Errorf("%s: Terraform and API names are required", fieldPath)
		}
		if _, exists := byName[field.TerraformName]; exists {
			return fmt.Errorf("%s: duplicate field", fieldPath)
		}
		byName[field.TerraformName] = field
		if err := validateField(field, fieldPath, parentRange, parentModes); err != nil {
			return err
		}
	}
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		fieldPath := path + "." + field.TerraformName
		if field.Reference != nil {
			companion, exists := byName[field.Reference.TypeField]
			if !exists {
				return fmt.Errorf("%s: reference type field %q does not exist in this scope", fieldPath, field.Reference.TypeField)
			}
			if !modesSubset(field.Modes, companion.Modes) {
				return fmt.Errorf("%s: reference type field %q is unavailable in one or more field modes", fieldPath, field.Reference.TypeField)
			}
			if !rangeSubset(field.Versions, companion.Versions) {
				return fmt.Errorf("%s: reference type field %q is unavailable in one or more field API versions", fieldPath, field.Reference.TypeField)
			}
			if len(field.Reference.AllowedTypes) == 0 {
				return fmt.Errorf("%s: reference allowed types are required", fieldPath)
			}
		}
		if field.AutoAssignment != nil {
			flag, exists := byName[field.AutoAssignment.FlagField]
			if !exists || flag.Kind != FieldKindBool {
				return fmt.Errorf("%s: auto-assignment flag %q must be a bool field in this scope", fieldPath, field.AutoAssignment.FlagField)
			}
			if !modesSubset(field.Modes, flag.Modes) {
				return fmt.Errorf("%s: auto-assignment flag %q is unavailable in one or more field modes", fieldPath, field.AutoAssignment.FlagField)
			}
			if !rangeSubset(field.Versions, flag.Versions) {
				return fmt.Errorf("%s: auto-assignment flag %q is unavailable in one or more field API versions", fieldPath, field.AutoAssignment.FlagField)
			}
		}
	}
	return nil
}

// validateUnmanagedField keeps an unsurfaced API field honest: it records the
// API shape and nothing else, so it cannot smuggle in Terraform behavior.
func validateUnmanagedField(field FieldSpec, path string, parentRange VersionRange, parentModes []Mode) error {
	if field.APIName == "" {
		return fmt.Errorf("%s: API name is required", path)
	}
	if field.TerraformName != "" {
		return fmt.Errorf("%s: unmanaged fields cannot declare a Terraform name", path)
	}
	if field.Access != "" || field.ResponseAbsence != "" || field.CreateNull != "" ||
		field.UpdateClear != "" || field.UnknownPlan != "" || field.StateOwnership != "" {
		return fmt.Errorf("%s: unmanaged fields cannot declare access or lifecycle policies", path)
	}
	if len(field.Fields) != 0 || field.Collection != nil || field.Reference != nil ||
		field.AutoAssignment != nil || field.Default != nil || len(field.Validators) != 0 {
		return fmt.Errorf("%s: unmanaged fields cannot declare nested structure or behavior", path)
	}
	// An unmanaged field still records API shape and applicability, so its kind,
	// modes, and version range are validated exactly as a managed field's are.
	// Only Terraform behavior is absent.
	if !validKind(field.Kind) {
		return fmt.Errorf("%s: field kind is required", path)
	}
	if err := validateModes(field.Modes); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if !modesSubset(field.Modes, parentModes) {
		return fmt.Errorf("%s: field modes must be a subset of parent modes", path)
	}
	if err := validateRange(field.Versions); err != nil {
		return fmt.Errorf("%s: versions: %w", path, err)
	}
	if !rangeSubset(field.Versions, parentRange) {
		return fmt.Errorf("%s: version range must be contained by its parent", path)
	}
	return nil
}

func validateField(field FieldSpec, path string, parentRange VersionRange, parentModes []Mode) error {
	if err := validateModes(field.Modes); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if err := validateRange(field.Versions); err != nil {
		return fmt.Errorf("%s: versions: %w", path, err)
	}
	if !modesSubset(field.Modes, parentModes) {
		return fmt.Errorf("%s: field modes must be a subset of parent modes", path)
	}
	if !rangeSubset(field.Versions, parentRange) {
		return fmt.Errorf("%s: version range must be contained by its parent", path)
	}
	if !validKind(field.Kind) || !validAccess(field.Access) {
		return fmt.Errorf("%s: field kind and access are required", path)
	}
	if !validPolicies(field) {
		return fmt.Errorf("%s: all lifecycle policies are required", path)
	}
	if err := validateOwnership(field, path); err != nil {
		return err
	}
	if field.Default != nil {
		if err := validateLiteral(*field.Default, field.Kind); err != nil {
			return fmt.Errorf("%s: default: %w", path, err)
		}
	}
	if err := validateValidators(field.Validators, field.Kind, path); err != nil {
		return err
	}
	if field.Kind == FieldKindObject {
		if field.Collection == nil {
			return fmt.Errorf("%s: object/list fields require an explicit collection policy", path)
		}
		if field.ElementKind != "" {
			return fmt.Errorf("%s: object fields cannot declare an element kind", path)
		}
		if err := validateCollection(*field.Collection, field.Fields, field.ElementKind, path); err != nil {
			return err
		}
		// The API declares some objects with no properties, and the resources that
		// expose them ship an empty block. Representing that faithfully needs a
		// singleton object with no fields; every other object still requires them.
		if len(field.Fields) == 0 {
			if field.Collection.Strategy != CollectionSingleton {
				return fmt.Errorf("%s: only a singleton object may have no fields", path)
			}
		} else if err := validateFields(field.Fields, path, field.Versions, field.Modes); err != nil {
			return err
		}
	} else if field.Kind == FieldKindList {
		if field.Collection == nil || !validKind(field.ElementKind) || field.ElementKind == FieldKindList {
			return fmt.Errorf("%s: list fields require an explicit collection policy and element kind", path)
		}
		if err := validateCollection(*field.Collection, field.Fields, field.ElementKind, path); err != nil {
			return err
		}
		if field.ElementKind == FieldKindObject {
			if err := validateFields(field.Fields, path, field.Versions, field.Modes); err != nil {
				return err
			}
		} else if len(field.Fields) != 0 {
			return fmt.Errorf("%s: scalar list elements cannot have nested fields", path)
		}
	} else if field.ElementKind != "" || len(field.Fields) != 0 || field.Collection != nil {
		return fmt.Errorf("%s: scalar fields cannot have nested fields or a collection policy", path)
	}
	return nil
}

func rangeSubset(rangeToCheck, container VersionRange) bool {
	return rangeToCheck.MinInclusive.Compare(container.MinInclusive) >= 0 && rangeToCheck.MaxExclusive.Compare(container.MaxExclusive) <= 0
}

func modesSubset(modes, parentModes []Mode) bool {
	allowed := make(map[Mode]bool, len(parentModes))
	for _, mode := range parentModes {
		allowed[mode] = true
	}
	for _, mode := range modes {
		if !allowed[mode] {
			return false
		}
	}
	return true
}

func validateRange(r VersionRange) error {
	if r.MinInclusive.Major < 0 || r.MinInclusive.Minor < 0 || r.MaxExclusive.Major < 0 || r.MaxExclusive.Minor < 0 || r.MinInclusive.Compare(r.MaxExclusive) >= 0 {
		return fmt.Errorf("expected a non-empty inclusive-minimum/exclusive-maximum range")
	}
	return nil
}

func validateModes(modes []Mode) error {
	if len(modes) == 0 {
		return fmt.Errorf("modes are required")
	}
	seen := map[Mode]bool{}
	for _, mode := range modes {
		if mode != ModeDatacenter && mode != ModeCampus {
			return fmt.Errorf("unknown mode %q", mode)
		}
		if seen[mode] {
			return fmt.Errorf("duplicate mode %q", mode)
		}
		seen[mode] = true
	}
	return nil
}

func validKind(kind FieldKind) bool {
	return kind == FieldKindString || kind == FieldKindBool || kind == FieldKindInt64 || kind == FieldKindNumber || kind == FieldKindObject || kind == FieldKindList
}
func validAccess(access Access) bool {
	return access == AccessRequired || access == AccessOptional || access == AccessComputed || access == AccessOptionalComputed
}
func validPolicies(f FieldSpec) bool {
	responseAbsenceValid := f.ResponseAbsence == ResponseAbsenceTerraformNull || f.ResponseAbsence == ResponseAbsenceDefault || f.ResponseAbsence == ResponseAbsencePreserve || f.ResponseAbsence == ResponseAbsenceError
	return responseAbsenceValid &&
		(f.CreateNull == CreateNullOmit || f.CreateNull == CreateNullAPINull || f.CreateNull == CreateNullDefault || f.CreateNull == CreateNullReject) &&
		(f.UpdateClear == UpdateClearAPINull || f.UpdateClear == UpdateClearEmptyString || f.UpdateClear == UpdateClearZero || f.UpdateClear == UpdateClearFalse || f.UpdateClear == UpdateClearDefault || f.UpdateClear == UpdateClearOmitUnmanaged || f.UpdateClear == UpdateClearReject) &&
		(f.UnknownPlan == UnknownPlanOmitAndRead || f.UnknownPlan == UnknownPlanPreserve || f.UnknownPlan == UnknownPlanReject) &&
		(f.StateOwnership == StateConfiguration || f.StateOwnership == StateServer || f.StateOwnership == StateConfigurationOrServer)
}

func validateOwnership(field FieldSpec, path string) error {
	expected := map[Access]StateOwnershipPolicy{AccessRequired: StateConfiguration, AccessOptional: StateConfiguration, AccessComputed: StateServer, AccessOptionalComputed: StateConfigurationOrServer}
	if field.StateOwnership != expected[field.Access] {
		return fmt.Errorf("%s: access %q requires state ownership %q", path, field.Access, expected[field.Access])
	}
	if field.Access == AccessComputed && (field.CreateNull != CreateNullOmit || field.UpdateClear != UpdateClearOmitUnmanaged) {
		return fmt.Errorf("%s: computed-only fields must omit create nulls and unmanaged update clears", path)
	}
	return nil
}

func validateCollection(collection CollectionSpec, fields []FieldSpec, elementKind FieldKind, path string) error {
	if collection.Strategy != CollectionSingleton && collection.Strategy != CollectionIndexedPatch && collection.Strategy != CollectionIndexedServerAssigned && collection.Strategy != CollectionReplace && collection.Strategy != CollectionComputedSubset {
		return fmt.Errorf("%s: collection strategy is required", path)
	}
	if collection.Ordering != CollectionOrdered && collection.Ordering != CollectionUnordered {
		return fmt.Errorf("%s: collection ordering is required", path)
	}
	if collection.Strategy != CollectionSingleton && collection.IdentityField == "" {
		if elementKind != "" && elementKind != FieldKindObject {
			return nil
		}
		return fmt.Errorf("%s: collection strategy %q requires an identity field", path, collection.Strategy)
	}
	if elementKind != "" && elementKind != FieldKindObject && collection.IdentityField != "" {
		return fmt.Errorf("%s: scalar list elements cannot declare an identity field", path)
	}
	if collection.IdentityField != "" {
		for _, field := range fields {
			if field.TerraformName == collection.IdentityField {
				return nil
			}
		}
		return fmt.Errorf("%s: collection identity field %q does not name a nested field", path, collection.IdentityField)
	}
	return nil
}

func validateLiteral(literal LiteralSpec, kind FieldKind) error {
	if literal.Kind == LiteralNull {
		if literal.Value != "" {
			return fmt.Errorf("null must have an empty value")
		}
		return nil
	}
	switch kind {
	case FieldKindString:
		if literal.Kind != LiteralString {
			return fmt.Errorf("string field requires a string literal")
		}
	case FieldKindBool:
		if literal.Kind != LiteralBool || (literal.Value != "true" && literal.Value != "false") {
			return fmt.Errorf("bool literal must be true or false")
		}
	case FieldKindInt64:
		if literal.Kind != LiteralInt64 {
			return fmt.Errorf("int64 field requires an int64 literal")
		}
		value, err := strconv.ParseInt(literal.Value, 10, 64)
		if err != nil || strconv.FormatInt(value, 10) != literal.Value {
			return fmt.Errorf("int64 must be canonical base-10")
		}
	case FieldKindNumber:
		if literal.Kind != LiteralDecimal || !canonicalDecimal.MatchString(literal.Value) {
			return fmt.Errorf("number field requires a canonical decimal literal")
		}
	default:
		return fmt.Errorf("object/list defaults are not supported")
	}
	return nil
}

func validateValidators(validators []ValidatorSpec, kind FieldKind, path string) error {
	seen := make(map[ValidatorKind]bool, len(validators))
	for _, validator := range validators {
		if seen[validator.Kind] {
			return fmt.Errorf("%s: duplicate validator %q", path, validator.Kind)
		}
		seen[validator.Kind] = true
		if err := validateValidator(validator, kind); err != nil {
			return fmt.Errorf("%s: validator %q: %w", path, validator.Kind, err)
		}
	}
	return nil
}

func validateValidator(validator ValidatorSpec, kind FieldKind) error {
	hasBounds := validator.Min != nil || validator.Max != nil
	hasNumberBounds := validator.MinValue != nil || validator.MaxValue != nil
	switch validator.Kind {
	case ValidatorOneOf:
		if len(validator.Values) == 0 || hasBounds || hasNumberBounds || validator.Pattern != "" {
			return fmt.Errorf("one_of requires values only")
		}
		for _, value := range validator.Values {
			if value.Kind == LiteralNull {
				return fmt.Errorf("one_of values cannot be null")
			}
			if err := validateLiteral(value, kind); err != nil {
				return err
			}
		}
	case ValidatorStringLength:
		if kind != FieldKindString || !hasBounds || hasNumberBounds || len(validator.Values) != 0 || validator.Pattern != "" {
			return fmt.Errorf("string_length requires string field and min and/or max only")
		}
		if (validator.Min != nil && *validator.Min < 0) || (validator.Max != nil && *validator.Max < 0) || (validator.Min != nil && validator.Max != nil && *validator.Min > *validator.Max) {
			return fmt.Errorf("string length bounds must be non-negative and ordered")
		}
	case ValidatorInt64Range:
		if kind != FieldKindInt64 || !hasBounds || hasNumberBounds || len(validator.Values) != 0 || validator.Pattern != "" {
			return fmt.Errorf("int64_range requires int64 field and min and/or max only")
		}
		if validator.Min != nil && validator.Max != nil && *validator.Min > *validator.Max {
			return fmt.Errorf("int64 range bounds must be ordered")
		}
	case ValidatorNumberRange:
		if kind != FieldKindNumber || !hasNumberBounds || hasBounds || len(validator.Values) != 0 || validator.Pattern != "" {
			return fmt.Errorf("number_range requires number field and decimal min_value and/or max_value only")
		}
		if err := validateDecimalBound("min_value", validator.MinValue); err != nil {
			return err
		}
		if err := validateDecimalBound("max_value", validator.MaxValue); err != nil {
			return err
		}
		if validator.MinValue != nil && validator.MaxValue != nil {
			min, minOK := new(big.Rat).SetString(validator.MinValue.Value)
			max, maxOK := new(big.Rat).SetString(validator.MaxValue.Value)
			if !minOK || !maxOK {
				return fmt.Errorf("number bounds must be valid decimals")
			}
			if min.Cmp(max) > 0 {
				return fmt.Errorf("number range bounds must be ordered")
			}
		}
	case ValidatorStringRegex:
		if kind != FieldKindString || validator.Pattern == "" || hasBounds || hasNumberBounds || len(validator.Values) != 0 {
			return fmt.Errorf("string_regex requires string field and pattern only")
		}
		if _, err := regexp.Compile(validator.Pattern); err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	default:
		return fmt.Errorf("unsupported validator kind")
	}
	return nil
}

func validateDecimalBound(name string, bound *LiteralSpec) error {
	if bound == nil {
		return nil
	}
	if bound.Kind != LiteralDecimal {
		return fmt.Errorf("%s must be a decimal literal", name)
	}
	return validateLiteral(*bound, FieldKindNumber)
}

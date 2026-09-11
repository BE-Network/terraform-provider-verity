// Mode field mappings.
//
// The bulk of this data is now derived from the reviewed spec registry and lives
// in generated_mode_metadata.go, produced by `specgen metadata`. What remains
// here is the hand-maintained residue for endpoints the registry does not
// represent yet, plus the lookup helpers.

package utils

type FieldMode string

const (
	FieldModeBoth       FieldMode = "both"
	FieldModeDatacenter FieldMode = "datacenter"
	FieldModeCampus     FieldMode = "campus"
)

// FieldAppliesToMode checks if a field applies to the given mode.
// Returns true if the field should be populated for the given mode.
// If the field is not found in ModeFields, it defaults to true (applies to both modes).
func FieldAppliesToMode(resourceType, fieldName, mode string) bool {
	resourceFields, ok := ModeFields[resourceType]
	if !ok {
		// Resource not found in mode map, assume field applies to all modes
		return true
	}

	fieldMode, ok := resourceFields[fieldName]
	if !ok {
		// Field not found in mode map, assume it applies to all modes
		return true
	}

	switch fieldMode {
	case FieldModeBoth:
		return true
	case FieldModeDatacenter:
		return mode == "datacenter"
	case FieldModeCampus:
		return mode == "campus"
	default:
		return true
	}
}

// pendingModeFields holds the endpoints the reviewed registry does not represent
// yet; see status.md. Entries move out as the registry grows, and ModeFields is
// the union of this table and the generated one.
// pendingModeFields is empty: every endpoint the provider exposes is now
// represented in the reviewed registry, so all mode-field data is generated.
var pendingModeFields = map[string]map[string]FieldMode{}

// ModeFields maps API field paths to the modes they apply to. It is derived from
// the reviewed spec registry, with the pending endpoints above merged in.
var ModeFields = mergeModeFields()

func mergeModeFields() map[string]map[string]FieldMode {
	merged := make(map[string]map[string]FieldMode, len(generatedModeFields)+len(pendingModeFields))
	for resource, fields := range generatedModeFields {
		merged[resource] = fields
	}
	for resource, fields := range pendingModeFields {
		merged[resource] = fields
	}
	return merged
}

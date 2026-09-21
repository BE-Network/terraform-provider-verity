package utils

type FieldMode string

const (
	FieldModeBoth       FieldMode = "both"
	FieldModeDatacenter FieldMode = "datacenter"
	FieldModeCampus     FieldMode = "campus"
)

func FieldAppliesToMode(resourceType, fieldName, mode string) bool {
	resourceFields, ok := ModeFields[resourceType]
	if !ok {

		return true
	}

	fieldMode, ok := resourceFields[fieldName]
	if !ok {

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

var pendingModeFields = map[string]map[string]FieldMode{}

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

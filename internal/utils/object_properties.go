package utils

import (
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ObjectPropertiesField struct {
	Name     string
	TFValue  interface{} // Can be types.String, types.Bool, types.Int64
	APIValue interface{} // Pointer to API field (e.g., **string, **bool, **int32)
}

// SetObjectPropertiesFields sets API fields from TF values for object_properties
func SetObjectPropertiesFields(fields []ObjectPropertiesField) {
	for _, field := range fields {
		switch tfVal := field.TFValue.(type) {
		case types.String:
			if apiPtr, ok := field.APIValue.(**string); ok {
				if !tfVal.IsNull() {
					*apiPtr = openapi.PtrString(tfVal.ValueString())
				} else {
					*apiPtr = nil
				}
			}
		case types.Bool:
			if apiPtr, ok := field.APIValue.(**bool); ok {
				if !tfVal.IsNull() {
					*apiPtr = openapi.PtrBool(tfVal.ValueBool())
				} else {
					*apiPtr = nil
				}
			}
		case types.Int64:
			// Try regular **int32 pointer first
			if apiPtr, ok := field.APIValue.(**int32); ok {
				if !tfVal.IsNull() {
					val := int32(tfVal.ValueInt64())
					*apiPtr = openapi.PtrInt32(val)
				} else {
					*apiPtr = nil
				}
			} else if apiNullablePtr, ok := field.APIValue.(*openapi.NullableInt32); ok {
				// Handle NullableInt32 type
				if !tfVal.IsNull() {
					val := int32(tfVal.ValueInt64())
					*apiNullablePtr = *openapi.NewNullableInt32(&val)
				} else {
					*apiNullablePtr = *openapi.NewNullableInt32(nil)
				}
			}
		}
	}
}

// CompareObjectPropertiesFields checks if any object_properties fields have changed
// Returns true if any field differs between plan and state
func CompareObjectPropertiesFields(fields []ObjectPropertiesFieldComparison) bool {
	for _, field := range fields {
		switch planVal := field.PlanValue.(type) {
		case types.String:
			if stateVal, ok := field.StateValue.(types.String); ok {
				if !planVal.Equal(stateVal) {
					return true
				}
			}
		case types.Bool:
			if stateVal, ok := field.StateValue.(types.Bool); ok {
				if !planVal.Equal(stateVal) {
					return true
				}
			}
		case types.Int64:
			if stateVal, ok := field.StateValue.(types.Int64); ok {
				if !planVal.Equal(stateVal) {
					return true
				}
			}
		}
	}
	return false
}

type ObjectPropertiesFieldComparison struct {
	Name       string
	PlanValue  interface{}
	StateValue interface{}
}

type ObjectPropertiesFieldWithComparison struct {
	Name       string
	PlanValue  interface{}
	StateValue interface{}
	APIValue   interface{}
}

// CompareAndSetObjectPropertiesFields sets only the fields that have changed between plan and state
func CompareAndSetObjectPropertiesFields(fields []ObjectPropertiesFieldWithComparison, hasChanges *bool) {
	for _, field := range fields {
		// Check if field has changed
		switch planVal := field.PlanValue.(type) {
		case types.String:
			if stateVal, ok := field.StateValue.(types.String); ok {
				if !planVal.Equal(stateVal) {
					if apiPtr, ok := field.APIValue.(**string); ok {
						if !planVal.IsNull() {
							*apiPtr = openapi.PtrString(planVal.ValueString())
						} else {
							*apiPtr = nil
						}
						*hasChanges = true
					}
				}
			}
		case types.Bool:
			if stateVal, ok := field.StateValue.(types.Bool); ok {
				if !planVal.Equal(stateVal) {
					if apiPtr, ok := field.APIValue.(**bool); ok {
						if !planVal.IsNull() {
							*apiPtr = openapi.PtrBool(planVal.ValueBool())
						} else {
							*apiPtr = nil
						}
						*hasChanges = true
					}
				}
			}
		case types.Int64:
			if stateVal, ok := field.StateValue.(types.Int64); ok {
				if !planVal.Equal(stateVal) {
					// Try regular **int32 pointer first
					if apiPtr, ok := field.APIValue.(**int32); ok {
						if !planVal.IsNull() {
							val := int32(planVal.ValueInt64())
							*apiPtr = openapi.PtrInt32(val)
						} else {
							*apiPtr = nil
						}
						*hasChanges = true
					} else if apiNullablePtr, ok := field.APIValue.(*openapi.NullableInt32); ok {
						// Handle NullableInt32 type
						if !planVal.IsNull() {
							val := int32(planVal.ValueInt64())
							*apiNullablePtr = *openapi.NewNullableInt32(&val)
						} else {
							*apiNullablePtr = *openapi.NewNullableInt32(nil)
						}
						*hasChanges = true
					}
				}
			}
		}
	}
}

// MapObjectPropertiesFromAPI reads object_properties from API map and populates TF state
// The mapper function is called with each field from the API map to allow custom mapping
type ObjectPropertiesMapper func(fieldName string, apiValue interface{}) interface{}

func MapObjectPropertiesFieldsFromAPI(objPropsMap map[string]interface{}, fieldNames []string) map[string]types.String {
	result := make(map[string]types.String)
	for _, fieldName := range fieldNames {
		result[fieldName] = MapStringFromAPI(objPropsMap[fieldName])
	}
	return result
}

func MapObjectPropertiesBoolFieldsFromAPI(objPropsMap map[string]interface{}, fieldNames []string) map[string]types.Bool {
	result := make(map[string]types.Bool)
	for _, fieldName := range fieldNames {
		result[fieldName] = MapBoolFromAPI(objPropsMap[fieldName])
	}
	return result
}

func MapObjectPropertiesInt64FieldsFromAPI(objPropsMap map[string]interface{}, fieldNames []string) map[string]types.Int64 {
	result := make(map[string]types.Int64)
	for _, fieldName := range fieldNames {
		result[fieldName] = MapInt64FromAPI(objPropsMap[fieldName])
	}
	return result
}

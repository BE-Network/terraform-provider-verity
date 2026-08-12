package utils

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type postOperationFallbackContextKey struct{}

type postOperationFallback struct {
	plan         interface{}
	resourceType string
	mode         string
}

// WithPostOperationFallback marks a Create/Update-initiated Read. Normal Read
// operations never carry this value and therefore always reflect the live API.
func WithPostOperationFallback(ctx context.Context, plan interface{}, resourceType, mode string) context.Context {
	return context.WithValue(ctx, postOperationFallbackContextKey{}, postOperationFallback{
		plan: plan, resourceType: resourceType, mode: mode,
	})
}

// ApplyPostOperationFallback fills only absent top-level scalar API fields from a
// known Create/Update plan. It copies the map so the bulk response cache remains a
// faithful record of the API response. An explicit API null is present and wins.
func ApplyPostOperationFallback(ctx context.Context, apiData map[string]interface{}) map[string]interface{} {
	fallback, ok := ctx.Value(postOperationFallbackContextKey{}).(postOperationFallback)
	if !ok {
		return apiData
	}
	return MergeMissingPlanScalars(apiData, fallback.plan, fallback.resourceType, fallback.mode)
}

// SetPostOperationFallbackState stores normalized plan state for a
// Create/Update-initiated Read that has no usable API object. It returns false
// for ordinary Reads, preserving their live-drift behavior.
func SetPostOperationFallbackState(ctx context.Context, state *tfsdk.State) (bool, diag.Diagnostics) {
	fallback, ok := ctx.Value(postOperationFallbackContextKey{}).(postOperationFallback)
	if !ok {
		return false, nil
	}
	return true, state.Set(ctx, nullifyUnknownPlanScalars(fallback.plan))
}

// MergeMissingPlanScalars is the non-context form used by direct cached-response
// Create/Update paths. It intentionally supports only top-level scalar attributes.
func MergeMissingPlanScalars(apiData map[string]interface{}, plan interface{}, resourceType, mode string) map[string]interface{} {
	merged := make(map[string]interface{}, len(apiData))
	for key, value := range apiData {
		merged[key] = value
	}

	value := reflect.ValueOf(plan)
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return merged
		}
		value = value.Elem()
	}
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return merged
	}

	modelType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := modelType.Field(i)
		key := field.Tag.Get("tfsdk")
		if key == "" || key == "-" || !FieldAppliesToMode(resourceType, key, mode) {
			continue
		}
		if _, present := merged[key]; present {
			continue
		}
		if scalar, known := terraformScalarValue(value.Field(i).Interface()); known {
			merged[key] = scalar
		}
	}

	return merged
}

func terraformScalarValue(value interface{}) (interface{}, bool) {
	switch typed := value.(type) {
	case types.Bool:
		if typed.IsNull() || typed.IsUnknown() {
			return nil, false
		}
		return typed.ValueBool(), true
	case types.String:
		if typed.IsNull() || typed.IsUnknown() {
			return nil, false
		}
		return typed.ValueString(), true
	case types.Int64:
		if typed.IsNull() || typed.IsUnknown() {
			return nil, false
		}
		return typed.ValueInt64(), true
	case types.Number:
		if typed.IsNull() || typed.IsUnknown() {
			return nil, false
		}
		return typed.ValueBigFloat().Text('g', -1), true
	default:
		return nil, false
	}
}

// nullifyUnknownPlanScalars makes a planned value safe as final state when a
// post-operation Read has no API result. Terraform cannot accept unknown values
// after apply; the top-level scalar types covered by the fallback become null.
func nullifyUnknownPlanScalars(plan interface{}) interface{} {
	value := reflect.ValueOf(plan)
	if !value.IsValid() {
		return plan
	}
	clonedValue := reflect.New(value.Type()).Elem()
	clonedValue.Set(value)
	if clonedValue.Kind() != reflect.Struct {
		return clonedValue.Interface()
	}
	for i := 0; i < clonedValue.NumField(); i++ {
		field := clonedValue.Field(i)
		switch typed := field.Interface().(type) {
		case types.Bool:
			if typed.IsUnknown() {
				field.Set(reflect.ValueOf(types.BoolNull()))
			}
		case types.String:
			if typed.IsUnknown() {
				field.Set(reflect.ValueOf(types.StringNull()))
			}
		case types.Int64:
			if typed.IsUnknown() {
				field.Set(reflect.ValueOf(types.Int64Null()))
			}
		case types.Number:
			if typed.IsUnknown() {
				field.Set(reflect.ValueOf(types.NumberNull()))
			}
		}
	}
	return clonedValue.Interface()
}

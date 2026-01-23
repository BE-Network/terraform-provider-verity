package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ============================================================================
// ModifyPlan helpers for mode-aware field nullification
// ============================================================================

// ModeFieldNullifier provides a way to nullify fields that don't apply to the current mode
type ModeFieldNullifier struct {
	Ctx          context.Context
	ResourceType string
	Mode         string
	Plan         interface {
		SetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
}

// NullifyStrings sets string fields to null if they don't apply to the current mode
func (n *ModeFieldNullifier) NullifyStrings(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.StringNull())
		}
	}
}

// NullifyBools sets bool fields to null if they don't apply to the current mode
func (n *ModeFieldNullifier) NullifyBools(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.BoolNull())
		}
	}
}

// NullifyInt64s sets int64 fields to null if they don't apply to the current mode
func (n *ModeFieldNullifier) NullifyInt64s(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.Int64Null())
		}
	}
}

// NullifyNumbers sets Number fields to null if they don't apply to the current mode
func (n *ModeFieldNullifier) NullifyNumbers(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.NumberNull())
		}
	}
}

// ============================================================================
// Nullable Field Handling for ModifyPlan
// ============================================================================

// NullableInt64Field represents a nullable Int64 field for explicit null detection
type NullableInt64Field struct {
	AttrName  string
	ConfigVal types.Int64
	StateVal  types.Int64
}

// NullableNumberField represents a nullable Number field for explicit null detection
type NullableNumberField struct {
	AttrName  string
	ConfigVal types.Number
	StateVal  types.Number
}

// NullableFieldsConfig holds the configuration for handling nullable fields in ModifyPlan
type NullableFieldsConfig struct {
	Ctx  context.Context
	Plan interface {
		SetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
	ConfiguredAttrs *ConfiguredAttributes
	Int64Fields     []NullableInt64Field
	NumberFields    []NullableNumberField
}

// HandleNullableFields processes nullable Int64 and Number fields for explicit null detection.
// For Optional+Computed fields, Terraform copies state to plan when config is null.
// This function detects explicit null in HCL and forces plan to null.
func HandleNullableFields(cfg NullableFieldsConfig) {
	// Handle nullable Int64 fields
	for _, field := range cfg.Int64Fields {
		// If explicitly configured as null in .tf AND state has a value -> force null
		if cfg.ConfiguredAttrs.IsConfigured(field.AttrName) && field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			cfg.Plan.SetAttribute(cfg.Ctx, path.Root(field.AttrName), types.Int64Null())
		}
	}

	// Handle nullable Number fields
	for _, field := range cfg.NumberFields {
		// If explicitly configured as null in .tf AND state has a value -> force null
		if cfg.ConfiguredAttrs.IsConfigured(field.AttrName) && field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			cfg.Plan.SetAttribute(cfg.Ctx, path.Root(field.AttrName), types.NumberNull())
		}
	}
}

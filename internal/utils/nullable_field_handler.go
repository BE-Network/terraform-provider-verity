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

// NullifyNestedBlocks sets nested block fields to null if the block doesn't apply to the current mode.
func (n *ModeFieldNullifier) NullifyNestedBlocks(blocks ...string) {
	for _, block := range blocks {
		if !FieldAppliesToMode(n.ResourceType, block, n.Mode) {
			// Set the nested block to an empty list when mode doesn't apply
			n.Plan.SetAttribute(n.Ctx, path.Root(block), types.ListNull(types.ObjectType{}))
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

// ============================================================================
// Nullable Nested Field Handling for ModifyPlan
// ============================================================================

// NullableNestedInt64Field represents a nullable Int64 field within an indexed nested block
type NullableNestedInt64Field struct {
	BlockIndex int64
	AttrName   string
	ConfigVal  types.Int64
	StateVal   types.Int64
}

// NullableNestedNumberField represents a nullable Number field within an indexed nested block
type NullableNestedNumberField struct {
	BlockIndex int64
	AttrName   string
	ConfigVal  types.Number
	StateVal   types.Number
}

// NullableNestedFieldsConfig holds the configuration for handling nullable nested fields in ModifyPlan
type NullableNestedFieldsConfig struct {
	Ctx  context.Context
	Plan interface {
		SetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
	ConfiguredAttrs *ConfiguredAttributes
	BlockType       string
	BlockListPath   string
	BlockListIndex  int
	Int64Fields     []NullableNestedInt64Field
	NumberFields    []NullableNestedNumberField
}

// HandleNullableNestedFields processes nullable fields within indexed nested blocks.
// For Optional+Computed fields in nested blocks, Terraform copies state to plan when config is null.
// This function detects explicit null in HCL and forces plan to null.
func HandleNullableNestedFields(cfg NullableNestedFieldsConfig) {
	// Handle nullable Int64 fields in nested blocks
	for _, field := range cfg.Int64Fields {
		// If explicitly configured as null in .tf AND state has a value -> force null
		if cfg.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(cfg.BlockType, field.BlockIndex, field.AttrName) &&
			field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			attrPath := path.Root(cfg.BlockListPath).AtListIndex(cfg.BlockListIndex).AtName(field.AttrName)
			cfg.Plan.SetAttribute(cfg.Ctx, attrPath, types.Int64Null())
		}
	}

	// Handle nullable Number fields in nested blocks
	for _, field := range cfg.NumberFields {
		// If explicitly configured as null in .tf AND state has a value -> force null
		if cfg.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(cfg.BlockType, field.BlockIndex, field.AttrName) &&
			field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			attrPath := path.Root(cfg.BlockListPath).AtListIndex(cfg.BlockListIndex).AtName(field.AttrName)
			cfg.Plan.SetAttribute(cfg.Ctx, attrPath, types.NumberNull())
		}
	}
}

// ============================================================================
// Indexed Block Item Helpers for Create/Update
// ============================================================================

// IndexedBlockItem is an interface for nested block items that have an Index field
type IndexedBlockItem interface {
	GetIndex() types.Int64
}

// BuildIndexedConfigMap builds a map from index value to config item for quick lookup.
func BuildIndexedConfigMap[T IndexedBlockItem](configItems []T) map[int64]T {
	result := make(map[int64]T)
	for _, item := range configItems {
		idx := item.GetIndex()
		if !idx.IsNull() && !idx.IsUnknown() {
			result[idx.ValueInt64()] = item
		}
	}
	return result
}

// IndexedBlockNullableFieldConfig provides config for setting nullable fields in indexed blocks
type IndexedBlockNullableFieldConfig struct {
	BlockType       string
	BlockIndex      int64
	ConfiguredAttrs *ConfiguredAttributes
}

// IsFieldConfigured checks if a field is configured in this indexed block
func (c *IndexedBlockNullableFieldConfig) IsFieldConfigured(attrName string) bool {
	return c.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(c.BlockType, c.BlockIndex, attrName)
}

// ============================================================================
// Object Properties Nullable Field Helpers
// ============================================================================

// ObjectPropertiesNullableFieldConfig provides config for checking nullable fields in object_properties blocks
type ObjectPropertiesNullableFieldConfig struct {
	ConfiguredAttrs *ConfiguredAttributes
}

// IsFieldConfigured checks if a field is configured in object_properties
func (c *ObjectPropertiesNullableFieldConfig) IsFieldConfigured(attrName string) bool {
	return c.ConfiguredAttrs.IsBlockAttributeConfigured("object_properties." + attrName)
}

// GetObjectPropertiesConfig returns the config item (from configItems if present, otherwise planItem)
// and the ObjectPropertiesNullableFieldConfig for checking which fields are explicitly configured.
// This is the equivalent of GetIndexedBlockConfig but for non-indexed object_properties blocks.
func GetObjectPropertiesConfig[T any](
	planItem T,
	configItems []T,
	configuredAttrs *ConfiguredAttributes,
) (T, *ObjectPropertiesNullableFieldConfig) {
	configItem := planItem
	if len(configItems) > 0 {
		configItem = configItems[0]
	}

	cfg := &ObjectPropertiesNullableFieldConfig{
		ConfiguredAttrs: configuredAttrs,
	}

	return configItem, cfg
}

// GetIndexedBlockConfig returns the config item (from configMap if present, otherwise planItem)
// and the IndexedBlockNullableFieldConfig for checking which fields are explicitly configured.
func GetIndexedBlockConfig[T IndexedBlockItem](
	planItem T,
	configMap map[int64]T,
	blockType string,
	configuredAttrs *ConfiguredAttributes,
) (T, *IndexedBlockNullableFieldConfig) {
	itemIndex := planItem.GetIndex().ValueInt64()

	// Get config item from map, fallback to plan item
	configItem := planItem
	if cfgItem, ok := configMap[itemIndex]; ok {
		configItem = cfgItem
	}

	cfg := &IndexedBlockNullableFieldConfig{
		BlockType:       blockType,
		BlockIndex:      itemIndex,
		ConfiguredAttrs: configuredAttrs,
	}

	return configItem, cfg
}

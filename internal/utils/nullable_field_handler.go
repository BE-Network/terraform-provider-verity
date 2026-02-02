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

// SubBlockFieldConfig defines fields within a nested block that is inside another nested block.
type SubBlockFieldConfig struct {
	SubBlockName string
	ItemCounts   []int
	StringFields []string
	BoolFields   []string
	Int64Fields  []string
	NumberFields []string
}

// NestedBlockFieldConfig defines the fields within a nested block that need mode-aware nullification.
type NestedBlockFieldConfig struct {
	BlockName    string
	ItemCount    int
	StringFields []string
	BoolFields   []string
	Int64Fields  []string
	NumberFields []string
	SubBlocks    []SubBlockFieldConfig
}

// NullifyNestedBlockFields nullifies individual fields within a nested block based on mode.
func (n *ModeFieldNullifier) NullifyNestedBlockFields(config NestedBlockFieldConfig) {
	// First check if the block itself applies to the mode
	if !FieldAppliesToMode(n.ResourceType, config.BlockName, n.Mode) {
		// Block doesn't apply, nullify the entire block
		n.Plan.SetAttribute(n.Ctx, path.Root(config.BlockName), types.ListNull(types.ObjectType{}))
		return
	}

	// Block applies to the mode, but we need to check individual fields
	// Iterate through each item in the block
	for i := 0; i < config.ItemCount; i++ {
		basePath := path.Root(config.BlockName).AtListIndex(i)

		// Nullify string fields that don't apply
		for _, field := range config.StringFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.StringNull())
			}
		}

		// Nullify bool fields that don't apply
		for _, field := range config.BoolFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.BoolNull())
			}
		}

		// Nullify int64 fields that don't apply
		for _, field := range config.Int64Fields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.Int64Null())
			}
		}

		// Nullify number fields that don't apply
		for _, field := range config.NumberFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.NumberNull())
			}
		}

		// Handle sub-blocks (deeply nested blocks)
		for _, subBlock := range config.SubBlocks {
			subBlockPath := config.BlockName + "." + subBlock.SubBlockName

			// Check if sub-block applies to mode
			if !FieldAppliesToMode(n.ResourceType, subBlockPath, n.Mode) {
				// Sub-block doesn't apply, nullify the entire sub-block
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(subBlock.SubBlockName), types.ListNull(types.ObjectType{}))
				continue
			}

			// Get item count for this parent index
			subBlockItemCount := 0
			if i < len(subBlock.ItemCounts) {
				subBlockItemCount = subBlock.ItemCounts[i]
			}

			// Iterate through each item in the sub-block
			for j := 0; j < subBlockItemCount; j++ {
				subBasePath := basePath.AtName(subBlock.SubBlockName).AtListIndex(j)

				// Nullify string fields that don't apply
				for _, field := range subBlock.StringFields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.StringNull())
					}
				}

				// Nullify bool fields that don't apply
				for _, field := range subBlock.BoolFields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.BoolNull())
					}
				}

				// Nullify int64 fields that don't apply
				for _, field := range subBlock.Int64Fields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.Int64Null())
					}
				}

				// Nullify number fields that don't apply
				for _, field := range subBlock.NumberFields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.NumberNull())
					}
				}
			}
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

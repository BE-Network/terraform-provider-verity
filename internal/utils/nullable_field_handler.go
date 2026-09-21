package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ModeFieldNullifier struct {
	Ctx          context.Context
	ResourceType string
	Mode         string
	Plan         interface {
		SetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
}

func (n *ModeFieldNullifier) NullifyStrings(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.StringNull())
		}
	}
}

func (n *ModeFieldNullifier) NullifyBools(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.BoolNull())
		}
	}
}

func (n *ModeFieldNullifier) NullifyInt64s(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.Int64Null())
		}
	}
}

func (n *ModeFieldNullifier) NullifyNumbers(fields ...string) {
	for _, field := range fields {
		if !FieldAppliesToMode(n.ResourceType, field, n.Mode) {
			n.Plan.SetAttribute(n.Ctx, path.Root(field), types.NumberNull())
		}
	}
}

type SubBlockFieldConfig struct {
	SubBlockName string
	ItemCounts   []int
	StringFields []string
	BoolFields   []string
	Int64Fields  []string
	NumberFields []string
}

type NestedBlockFieldConfig struct {
	BlockName    string
	ItemCount    int
	StringFields []string
	BoolFields   []string
	Int64Fields  []string
	NumberFields []string
	SubBlocks    []SubBlockFieldConfig
}

func (n *ModeFieldNullifier) NullifyNestedBlockFields(config NestedBlockFieldConfig) {

	if !FieldAppliesToMode(n.ResourceType, config.BlockName, n.Mode) {

		n.Plan.SetAttribute(n.Ctx, path.Root(config.BlockName), types.ListNull(types.ObjectType{}))
		return
	}

	for i := 0; i < config.ItemCount; i++ {
		basePath := path.Root(config.BlockName).AtListIndex(i)

		for _, field := range config.StringFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.StringNull())
			}
		}

		for _, field := range config.BoolFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.BoolNull())
			}
		}

		for _, field := range config.Int64Fields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.Int64Null())
			}
		}

		for _, field := range config.NumberFields {
			fieldPath := config.BlockName + "." + field
			if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
				n.Plan.SetAttribute(n.Ctx, basePath.AtName(field), types.NumberNull())
			}
		}

		for _, subBlock := range config.SubBlocks {
			subBlockPath := config.BlockName + "." + subBlock.SubBlockName

			if !FieldAppliesToMode(n.ResourceType, subBlockPath, n.Mode) {

				n.Plan.SetAttribute(n.Ctx, basePath.AtName(subBlock.SubBlockName), types.ListNull(types.ObjectType{}))
				continue
			}

			subBlockItemCount := 0
			if i < len(subBlock.ItemCounts) {
				subBlockItemCount = subBlock.ItemCounts[i]
			}

			for j := 0; j < subBlockItemCount; j++ {
				subBasePath := basePath.AtName(subBlock.SubBlockName).AtListIndex(j)

				for _, field := range subBlock.StringFields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.StringNull())
					}
				}

				for _, field := range subBlock.BoolFields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.BoolNull())
					}
				}

				for _, field := range subBlock.Int64Fields {
					fieldPath := subBlockPath + "." + field
					if !FieldAppliesToMode(n.ResourceType, fieldPath, n.Mode) {
						n.Plan.SetAttribute(n.Ctx, subBasePath.AtName(field), types.Int64Null())
					}
				}

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

type NullableInt64Field struct {
	AttrName  string
	ConfigVal types.Int64
	StateVal  types.Int64
}

type NullableNumberField struct {
	AttrName  string
	ConfigVal types.Number
	StateVal  types.Number
}

type NullableFieldsConfig struct {
	Ctx  context.Context
	Plan interface {
		SetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}
	ConfiguredAttrs *ConfiguredAttributes
	Int64Fields     []NullableInt64Field
	NumberFields    []NullableNumberField
}

func HandleNullableFields(cfg NullableFieldsConfig) {

	for _, field := range cfg.Int64Fields {

		if cfg.ConfiguredAttrs.IsConfigured(field.AttrName) && field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			cfg.Plan.SetAttribute(cfg.Ctx, path.Root(field.AttrName), types.Int64Null())
		}
	}

	for _, field := range cfg.NumberFields {

		if cfg.ConfiguredAttrs.IsConfigured(field.AttrName) && field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			cfg.Plan.SetAttribute(cfg.Ctx, path.Root(field.AttrName), types.NumberNull())
		}
	}
}

type NullableNestedInt64Field struct {
	BlockIndex int64
	AttrName   string
	ConfigVal  types.Int64
	StateVal   types.Int64
}

type NullableNestedNumberField struct {
	BlockIndex int64
	AttrName   string
	ConfigVal  types.Number
	StateVal   types.Number
}

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

func HandleNullableNestedFields(cfg NullableNestedFieldsConfig) {

	for _, field := range cfg.Int64Fields {

		if cfg.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(cfg.BlockType, field.BlockIndex, field.AttrName) &&
			field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			attrPath := path.Root(cfg.BlockListPath).AtListIndex(cfg.BlockListIndex).AtName(field.AttrName)
			cfg.Plan.SetAttribute(cfg.Ctx, attrPath, types.Int64Null())
		}
	}

	for _, field := range cfg.NumberFields {

		if cfg.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(cfg.BlockType, field.BlockIndex, field.AttrName) &&
			field.ConfigVal.IsNull() && !field.StateVal.IsNull() {
			attrPath := path.Root(cfg.BlockListPath).AtListIndex(cfg.BlockListIndex).AtName(field.AttrName)
			cfg.Plan.SetAttribute(cfg.Ctx, attrPath, types.NumberNull())
		}
	}
}

type IndexedBlockItem interface {
	GetIndex() types.Int64
}

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

type IndexedBlockNullableFieldConfig struct {
	BlockType       string
	BlockIndex      int64
	ConfiguredAttrs *ConfiguredAttributes
}

func (c *IndexedBlockNullableFieldConfig) IsFieldConfigured(attrName string) bool {
	return c.ConfiguredAttrs.IsIndexedBlockAttributeConfigured(c.BlockType, c.BlockIndex, attrName)
}

type ObjectPropertiesNullableFieldConfig struct {
	ConfiguredAttrs *ConfiguredAttributes
}

func (c *ObjectPropertiesNullableFieldConfig) IsFieldConfigured(attrName string) bool {
	return c.ConfiguredAttrs.IsBlockAttributeConfigured("object_properties." + attrName)
}

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

func GetIndexedBlockConfig[T IndexedBlockItem](
	planItem T,
	configMap map[int64]T,
	blockType string,
	configuredAttrs *ConfiguredAttributes,
) (T, *IndexedBlockNullableFieldConfig) {
	itemIndex := planItem.GetIndex().ValueInt64()

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

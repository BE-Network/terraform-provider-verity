package importer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
	Mode   string

	supported   map[string]*SchemaFields
	unsupported map[string]map[string]bool
}

type ResourceConfig struct {
	ResourceType                 string
	StageName                    string
	NestedBlockFields            map[string]bool
	ObjectPropsNestedBlockFields map[string]bool
	FieldMappings                map[string]string
	SkipTopLevelKeys             map[string]bool
}

var apiFieldRenames = map[string]map[string]string{
	"verity_device_voice_settings": {"Codecs": "codecs"},
}

var rootIndexSkipped = map[string]bool{
	"verity_gateway_profile":  true,
	"verity_eth_port_profile": true,
	"verity_bundle":           true,
}

func (i *Importer) resourceConfig(terraformType string) (ResourceConfig, error) {
	resource, err := registry.Lookup(terraformType)
	if err != nil {
		return ResourceConfig{}, err
	}
	stageName, err := stageNameFor(i.Mode, terraformType)
	if err != nil {
		return ResourceConfig{}, err
	}
	config := ResourceConfig{
		ResourceType:                 strings.TrimPrefix(terraformType, "verity_"),
		StageName:                    stageName,
		NestedBlockFields:            map[string]bool{},
		ObjectPropsNestedBlockFields: map[string]bool{},
		FieldMappings:                map[string]string{},
		SkipTopLevelKeys:             map[string]bool{"name": true},
	}
	if rootIndexSkipped[terraformType] {
		config.SkipTopLevelKeys["index"] = true
	}
	for apiName, terraformName := range apiFieldRenames[terraformType] {
		config.FieldMappings[apiName] = terraformName
	}
	for _, field := range resource.Fields {
		if field.Unmanaged {
			continue
		}
		if field.APIName != field.TerraformName {
			config.FieldMappings[field.APIName] = field.TerraformName
		}
		switch {
		case field.Kind == spec.FieldKindList && len(field.Fields) > 0:
			config.NestedBlockFields[field.TerraformName] = true
		case field.Kind == spec.FieldKindObject:
			for _, member := range field.Fields {
				if !member.Unmanaged && member.Kind == spec.FieldKindList && len(member.Fields) > 0 {
					config.ObjectPropsNestedBlockFields[member.TerraformName] = true
				}
			}
		}
	}
	return config, nil
}

func stageNameFor(mode, terraformType string) (string, error) {
	stages, err := stageOrder(mode)
	if err != nil {
		return "", err
	}
	for _, stage := range stages {
		if stage.ResourceType == terraformType {
			return stage.StageName, nil
		}
	}
	return "", fmt.Errorf("%s has no import stage for mode %q", terraformType, mode)
}

var nameSplitRE = regexp.MustCompile(`(\d+|\D+)`)

func getNaturalSortParts(s string) []interface{} {
	matches := nameSplitRE.FindAllString(s, -1)
	parts := make([]interface{}, len(matches))
	for i, match := range matches {
		if num, err := strconv.Atoi(match); err == nil {
			parts[i] = num
		} else {
			parts[i] = match
		}
	}
	return parts
}

func NewImporter(client *openapi.APIClient, mode string) *Importer {
	return &Importer{
		client: client,
		ctx:    context.Background(),
		Mode:   mode,
	}
}

func (i *Importer) ImportAll(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	tflog.Info(i.ctx, "Starting importer with mode", map[string]interface{}{
		"mode":        i.Mode,
		"api_version": utils.GetSupportedAPIVersionString(),
	})

	tasks, err := i.importTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		tflog.Info(i.ctx, "Importing resource", map[string]interface{}{
			"resource_name":           task.name,
			"terraform_resource_type": task.terraformResourceType,
		})

		data, err := i.fetchResource(task.resource)
		if err != nil {
			tflog.Error(i.ctx, "Failed to import resource", map[string]interface{}{"resource_name": task.name, "error": err})
			return fmt.Errorf("failed to import %s: %w", task.name, err)
		}
		if len(data) == 0 {
			tflog.Info(i.ctx, "No data found for resource, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}

		config, err := i.resourceConfig(task.terraformResourceType)
		if err != nil {
			tflog.Error(i.ctx, "No registry entry for terraform type", map[string]interface{}{
				"resource_name":  task.name,
				"terraform_type": task.terraformResourceType,
				"error":          err,
			})
			return fmt.Errorf("no registry entry for %s: %w", task.terraformResourceType, err)
		}

		i.PruneUnsupported(task.terraformResourceType, data)

		tfConfig, err := i.generateResourceTF(data, config)
		if err != nil {
			tflog.Error(i.ctx, "Failed to generate Terraform config", map[string]interface{}{"resource_name": task.name, "error": err})
			return fmt.Errorf("failed to generate terraform config for %s: %w", task.name, err)
		}

		if strings.TrimSpace(tfConfig) == "" {
			tflog.Info(i.ctx, "Generated TF config is empty, skipping file write", map[string]interface{}{"resource_name": task.name})
			continue
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.tf", task.name))
		if err := os.WriteFile(outputFile, []byte(tfConfig), 0644); err != nil {
			tflog.Error(i.ctx, "Failed to write TF config to file", map[string]interface{}{"resource_name": task.name, "file": outputFile, "error": err})
			return fmt.Errorf("failed to write %s terraform config: %w", task.name, err)
		}
		tflog.Info(i.ctx, "Successfully wrote TF config for resource", map[string]interface{}{"resource_name": task.name, "file": outputFile})
	}

	stagesTF, err := i.generateStagesTF()
	if err != nil {
		tflog.Error(i.ctx, "Failed to generate stages TF", map[string]interface{}{"error": err})
		return fmt.Errorf("failed to generate stages: %w", err)
	}

	stagesFile := filepath.Join(outputDir, "stages.tf")
	if err := os.WriteFile(stagesFile, []byte(stagesTF), 0644); err != nil {
		tflog.Error(i.ctx, "Failed to write stages TF config", map[string]interface{}{"error": err, "file": stagesFile})
		return fmt.Errorf("failed to write stages terraform config: %w", err)
	}

	return nil
}

type importTask struct {
	name                  string
	terraformResourceType string
	resource              spec.ResourceSpec
}

var importFileNames = map[string]string{
	"verity_acl_v4": "acls_ipv4",
	"verity_acl_v6": "acls_ipv6",
}

func (i *Importer) importTasks() ([]importTask, error) {
	tasks := make([]importTask, 0, 50)
	resourceTypes, err := ResourceTypeOrder(i.Mode)
	if err != nil {
		return nil, err
	}
	for _, terraformType := range resourceTypes {
		resource, err := registry.Lookup(terraformType)
		if err != nil {
			return nil, fmt.Errorf("no registry entry for %s: %w", terraformType, err)
		}
		name, named := importFileNames[terraformType]
		if !named {
			name = strings.Trim(resource.API.EndpointPath, "/")
		}
		tasks = append(tasks, importTask{name: name, terraformResourceType: terraformType, resource: resource})
	}
	return tasks, nil
}

func (i *Importer) fetchResource(resource spec.ResourceSpec) (map[string]map[string]interface{}, error) {
	return i.fetch(resource.TerraformType, resource.API.EndpointPath, resource.API.FixedHeaders, resource.API.ResponseCollectionKey)
}

func (i *Importer) fetch(label, endpointPath string, fixedHeaders map[string]string, collectionKey string) (map[string]map[string]interface{}, error) {
	collection, err := transport.FetchCollection(i.ctx, i.client, label, endpointPath, fixedHeaders, collectionKey)
	if err != nil {
		return nil, err
	}
	objects := make(map[string]map[string]interface{}, len(collection))
	for name, raw := range collection {
		object, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s object %q is not an object", label, name)
		}
		objects[name] = object
	}
	return objects, nil
}

func (i *Importer) generateResourceTF(resourcesMap map[string]map[string]interface{}, config ResourceConfig) (string, error) {

	var resourceNames []string
	for name := range resourcesMap {
		resourceNames = append(resourceNames, name)
	}

	sort.SliceStable(resourceNames, func(i, j int) bool {
		s1 := resourceNames[i]
		s2 := resourceNames[j]

		parts1 := getNaturalSortParts(s1)
		parts2 := getNaturalSortParts(s2)

		len1 := len(parts1)
		len2 := len(parts2)
		minLen := len1
		if len2 < minLen {
			minLen = len2
		}

		for k := 0; k < minLen; k++ {
			p1 := parts1[k]
			p2 := parts2[k]

			p1Int, p1IsInt := p1.(int)
			p2Int, p2IsInt := p2.(int)

			if p1IsInt && p2IsInt {
				if p1Int != p2Int {
					return p1Int < p2Int
				}
			} else if !p1IsInt && !p2IsInt {
				p1Str := p1.(string)
				p2Str := p2.(string)
				if p1Str != p2Str {
					return p1Str < p2Str
				}
			} else {
				return p1IsInt
			}
		}
		return len1 < len2
	})

	var tfConfig strings.Builder

	for _, name := range resourceNames {
		resource := resourcesMap[name]
		sanitizedName := utils.SanitizeResourceName(name)

		tfConfig.WriteString(fmt.Sprintf("\nresource \"verity_%s\" \"%s\" {\n", config.ResourceType, sanitizedName))
		tfConfig.WriteString(fmt.Sprintf("    name = %q\n", name))
		tfConfig.WriteString(fmt.Sprintf("    depends_on = [verity_operation_stage.%s]\n", config.StageName))

		skipObjectProperties := config.SkipTopLevelKeys["object_properties"]

		if !skipObjectProperties {

			objPropsRaw, objectPropertiesExists := resource["object_properties"]

			if objectPropertiesExists {
				tfConfig.WriteString("	object_properties")
				objProps, _ := objPropsRaw.(map[string]interface{})

				isEmptyObjectProps := len(objProps) == 0

				if isEmptyObjectProps {
					tfConfig.WriteString(" {}\n")
				} else {

					tfConfig.WriteString(" {\n")
					var objPropsContentBuilder strings.Builder
					universalObjectPropsHandler(objProps, &objPropsContentBuilder, config)
					tfConfig.WriteString(objPropsContentBuilder.String())
					tfConfig.WriteString("	}\n")
				}
			}

		}

		skipKeysSet := map[string]bool{"object_properties": true}
		for key := range config.SkipTopLevelKeys {
			skipKeysSet[key] = true
		}

		var topLevelKeys []string
		for key := range resource {
			if skipKeysSet[key] {
				continue
			}
			if isAutoAssignedField(resource, key) {
				continue
			}
			topLevelKeys = append(topLevelKeys, key)
		}
		sort.Strings(topLevelKeys)

		for _, key := range topLevelKeys {
			value := resource[key]

			tfFieldName := key
			if config.FieldMappings != nil {
				if mappedName, exists := config.FieldMappings[key]; exists {
					tfFieldName = mappedName
				}
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", tfFieldName, v))
			case float64:

				if v == float64(int(v)) {
					tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", tfFieldName, int(v)))
				} else {
					tfConfig.WriteString(fmt.Sprintf("	%s = %g\n", tfFieldName, v))
				}
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", tfFieldName, formatValue(v)))
			case []interface{}:
				if _, isNestedBlock := config.NestedBlockFields[tfFieldName]; isNestedBlock {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", tfFieldName))

							printedIndex := false
							if indexVal, idxExists := itemMap["index"]; idxExists {
								if indexFloat, isFloat := indexVal.(float64); isFloat {
									tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(indexFloat)))
									printedIndex = true
								}
							}

							var nestedItemKeys []string
							for itemKey := range itemMap {
								if itemKey == "index" && printedIndex {
									continue
								}
								nestedItemKeys = append(nestedItemKeys, itemKey)
							}
							sort.Strings(nestedItemKeys)

							for _, itemKey := range nestedItemKeys {
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemMap[itemKey])))
							}
							tfConfig.WriteString("	}\n")
						}
					}
				} else {
					tfConfig.WriteString(fmt.Sprintf("	%s = [\n", tfFieldName))
					for _, item := range v {
						if str, ok := item.(string); ok {
							tfConfig.WriteString(fmt.Sprintf("		%s,\n", formatValue(str)))
						}
					}
					tfConfig.WriteString("	]\n")
				}
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", tfFieldName))
			}
		}
		tfConfig.WriteString("}\n\n")
	}
	return tfConfig.String(), nil
}

type stageDefinition struct {
	StageName      string
	ResourceType   string
	DependsOnStage string
	order          int
}

func stageOrder(mode string) ([]stageDefinition, error) {
	resources, err := registry.Load()
	if err != nil {
		return nil, fmt.Errorf("load import stages: %w", err)
	}
	stages := make([]stageDefinition, 0, len(resources))
	for _, resource := range resources {
		stage, declared := resource.ImportStages[spec.Mode(mode)]
		if !declared {
			continue
		}
		stages = append(stages, stageDefinition{StageName: stage.Name, ResourceType: resource.TerraformType, order: stage.Order})
	}
	sort.Slice(stages, func(i, j int) bool { return stages[i].order < stages[j].order })
	for index := range stages {
		if index > 0 {
			stages[index].DependsOnStage = stages[index-1].StageName
		}
	}
	return stages, nil
}

func ResourceTypeOrder(mode string) ([]string, error) {
	stages, err := stageOrder(mode)
	if err != nil {
		return nil, err
	}
	order := make([]string, 0, len(stages))
	for _, stage := range stages {
		if utils.IsResourceCompatibleWithMode(stage.ResourceType, mode) {
			order = append(order, stage.ResourceType)
		}
	}
	return order, nil
}

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	tflog.Info(i.ctx, "Generating stages for mode", map[string]interface{}{
		"mode": i.Mode,
	})

	stages, err := stageOrder(i.Mode)
	if err != nil {
		return "", err
	}

	var compatibleStages []stageDefinition
	var lastCompatibleStage string

	for _, stage := range stages {
		if utils.IsResourceCompatibleWithMode(stage.ResourceType, i.Mode) {
			if stage.DependsOnStage != "" && lastCompatibleStage != "" && stage.DependsOnStage != lastCompatibleStage {
				stage.DependsOnStage = lastCompatibleStage
			}
			compatibleStages = append(compatibleStages, stage)
			lastCompatibleStage = stage.StageName
		} else {
			tflog.Debug(i.ctx, "Excluding stage for incompatible resource", map[string]interface{}{
				"stage_name":    stage.StageName,
				"resource_type": stage.ResourceType,
				"mode":          i.Mode,
			})
		}
	}

	modeComment := strings.ToUpper(i.Mode)
	tfConfig.WriteString(fmt.Sprintf("\n# These resources establish ordering for bulk operations in %s mode\n", modeComment))

	for _, stage := range compatibleStages {
		tfConfig.WriteString(fmt.Sprintf("resource \"verity_operation_stage\" \"%s\" {\n", stage.StageName))

		if stage.DependsOnStage != "" {
			tfConfig.WriteString(fmt.Sprintf("  depends_on = [verity_operation_stage.%s]\n", stage.DependsOnStage))
		}

		tfConfig.WriteString("  lifecycle {\n")
		tfConfig.WriteString("    create_before_destroy = true\n")
		tfConfig.WriteString("  }\n")
		tfConfig.WriteString("}\n\n")
	}

	tflog.Info(i.ctx, "Generated stages", map[string]interface{}{
		"mode":              i.Mode,
		"total_stages":      len(stages),
		"compatible_stages": len(compatibleStages),
	})

	return tfConfig.String(), nil
}

func universalObjectPropsHandler(objProps map[string]interface{}, builder *strings.Builder, config ResourceConfig) {
	if len(objProps) > 0 {
		var keys []string
		for key := range objProps {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := objProps[key]

			if config.ObjectPropsNestedBlockFields != nil && config.ObjectPropsNestedBlockFields[key] {

				if valueArray, ok := value.([]interface{}); ok {
					for _, item := range valueArray {
						builder.WriteString(fmt.Sprintf("		%s {\n", key))
						if itemMap, ok := item.(map[string]interface{}); ok {
							var itemKeys []string
							for itemKey := range itemMap {
								itemKeys = append(itemKeys, itemKey)
							}
							sort.Strings(itemKeys)

							for _, itemKey := range itemKeys {
								itemValue := itemMap[itemKey]
								builder.WriteString(fmt.Sprintf("			%s = %s\n", itemKey, formatObjectPropsValue(itemValue, "		")))
							}
						}
						builder.WriteString("		}\n")
					}
				}
			} else {

				builder.WriteString(fmt.Sprintf("		%s = %s\n", key, formatObjectPropsValue(value, "	")))
			}
		}
	}

}

func formatObjectPropsValue(value interface{}, indent string) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:

		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%g", v)
	case nil:
		return "null"
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}

		var result strings.Builder
		result.WriteString("[\n")
		for i, item := range v {
			result.WriteString(indent + "		")
			if itemMap, ok := item.(map[string]interface{}); ok {

				result.WriteString("{\n")
				var keys []string
				for key := range itemMap {
					keys = append(keys, key)
				}
				sort.Strings(keys)

				for _, key := range keys {
					itemValue := itemMap[key]
					result.WriteString(fmt.Sprintf("%s			%s = %s\n", indent, key, formatObjectPropsValue(itemValue, indent+"		")))
				}
				result.WriteString(indent + "		}")
			} else {

				result.WriteString(formatObjectPropsValue(item, indent+"		"))
			}

			if i < len(v)-1 {
				result.WriteString(",")
			}
			result.WriteString("\n")
		}
		result.WriteString(indent + "	]")
		return result.String()
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}"
		}

		var result strings.Builder
		result.WriteString("{\n")
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for i, key := range keys {
			objValue := v[key]
			result.WriteString(fmt.Sprintf("%s		%s = %s", indent, key, formatObjectPropsValue(objValue, indent+"	")))
			if i < len(keys)-1 {
				result.WriteString(",")
			}
			result.WriteString("\n")
		}
		result.WriteString(indent + "	}")
		return result.String()
	default:
		return "null"
	}
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:

		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%g", v)
	case nil:
		return "null"
	default:
		return "null"
	}
}

func isAutoAssignedField(resource map[string]interface{}, fieldName string) bool {
	autoAssignedFieldName := fieldName + "_auto_assigned_"

	if autoAssignedValue, ok := resource[autoAssignedFieldName].(bool); ok && autoAssignedValue {
		return true
	}

	return false
}

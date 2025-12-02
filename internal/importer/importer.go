package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
	Mode   string
}

type NestedBlockIterationStyle struct {
	PrintIndexFirst     bool
	SkipIndexInMainLoop bool
	IterateAllAsMap     bool
}

type ResourceConfig struct {
	ResourceType                 string
	StageName                    string
	HeaderNameLineFormat         string
	HeaderDependsOnLineFormat    string
	ObjectPropsHandler           func(objProps map[string]interface{}, builder *strings.Builder, config ResourceConfig)
	NestedBlockFields            map[string]bool
	ObjectPropsNestedBlockFields map[string]bool
	FieldMappings                map[string]string
	AdditionalTopLevelSkipKeys   []string
	EmptyObjectPropsAsSingleLine bool
	NestedBlockStyles            map[string]NestedBlockIterationStyle
}

type ImporterFunc func(context.Context, *openapi.APIClient) (*http.Response, error)

var nameSplitRE = regexp.MustCompile(`(\d+|\D+)`)

// importerRegistry maps resource names to their API caller function and JSON key
var importerRegistry = map[string]struct {
	apiCaller ImporterFunc
	jsonKey   string
}{
	"gateways": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.GatewaysAPI.GatewaysGet(ctx).Execute()
	}, jsonKey: "gateway"},
	"gatewayprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.GatewayProfilesAPI.GatewayprofilesGet(ctx).Execute()
	}, jsonKey: "gateway_profile"},
	"ethportprofiles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
	}, jsonKey: "eth_port_profile_"},
	"ethportsettings": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
	}, jsonKey: "eth_port_settings"},
	"lags": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.LAGsAPI.LagsGet(ctx).Execute()
	}, jsonKey: "lag"},
	"services": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.ServicesAPI.ServicesGet(ctx).Execute()
	}, jsonKey: "service"},
	"tenants": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.TenantsAPI.TenantsGet(ctx).Execute()
	}, jsonKey: "tenant"},
	"bundles": {apiCaller: func(ctx context.Context, client *openapi.APIClient) (*http.Response, error) {
		return client.BundlesAPI.BundlesGet(ctx).Execute()
	}, jsonKey: "endpoint_bundle"},
}

// terraformTypeToResourceKey maps Terraform resource types to resourceConfigs keys
var terraformTypeToResourceKey = map[string]string{
	"verity_tenant":            "tenant",
	"verity_gateway":           "gateway",
	"verity_gateway_profile":   "gateway_profile",
	"verity_eth_port_profile":  "eth_port_profile",
	"verity_lag":               "lag",
	"verity_service":           "service",
	"verity_eth_port_settings": "eth_port_settings",
	"verity_bundle":            "bundle",
}

// resourceConfigs is a registry of all resource configurations for generating Terraform code
var resourceConfigs = map[string]ResourceConfig{
	"tenant": {
		ResourceType:              "tenant",
		StageName:                 "tenant_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"route_tenants": true},
	},
	"gateway": {
		ResourceType:              "gateway",
		StageName:                 "gateway_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"static_routes": true},
	},
	"gateway_profile": {
		ResourceType:              "gateway_profile",
		StageName:                 "gateway_profile_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"external_gateways": true},
	},
	"eth_port_profile": {
		ResourceType:              "eth_port_profile",
		StageName:                 "eth_port_profile_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"services": true},
	},
	"eth_port_settings": {
		ResourceType:              "eth_port_settings",
		StageName:                 "eth_port_settings_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"lag": {
		ResourceType:                 "lag",
		StageName:                    "lag_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	},
	"service": {
		ResourceType:              "service",
		StageName:                 "service_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	},
	"bundle": {
		ResourceType:              "bundle",
		StageName:                 "bundle_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"eth_port_paths": true, "user_services": true},
		NestedBlockStyles: map[string]NestedBlockIterationStyle{
			"eth_port_paths": {IterateAllAsMap: true},
			"user_services":  {IterateAllAsMap: true},
		},
	},
}

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

func (i *Importer) getAPIVersion() string {
	defaultVersion := "6.4"

	versionResp, err := i.client.VersionAPI.VersionGet(i.ctx).Execute()
	if err != nil {
		// API 6.4 doesn't support /version endpoint (returns 404)
		tflog.Info(i.ctx, "Version endpoint not available, using default version", map[string]interface{}{"default_version": defaultVersion})
		return defaultVersion
	}
	defer versionResp.Body.Close()

	// API 6.5+ should return a valid version response
	var versionData struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(versionResp.Body).Decode(&versionData); err != nil || versionData.Version == "" {
		tflog.Warn(i.ctx, "Failed to parse version response, using default", map[string]interface{}{"error": err, "default_version": defaultVersion})
		return defaultVersion
	}

	return versionData.Version
}

// ImportAll fetches all resources and saves them as Terraform configuration files
func (i *Importer) ImportAll(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	tflog.Info(i.ctx, "Starting importer with mode", map[string]interface{}{
		"mode": i.Mode,
	})

	apiVersionString := i.getAPIVersion()
	tflog.Info(i.ctx, "Using API Version for import", map[string]interface{}{"version": apiVersionString})

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

	allResourceTasks := []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
	}{
		{name: "tenants", terraformResourceType: "verity_tenant", importer: func() (interface{}, error) { return i.importResource("tenants") }},
		{name: "gateways", terraformResourceType: "verity_gateway", importer: func() (interface{}, error) { return i.importResource("gateways") }},
		{name: "gatewayprofiles", terraformResourceType: "verity_gateway_profile", importer: func() (interface{}, error) { return i.importResource("gatewayprofiles") }},
		{name: "ethportprofiles", terraformResourceType: "verity_eth_port_profile", importer: func() (interface{}, error) { return i.importResource("ethportprofiles") }},
		{name: "ethportsettings", terraformResourceType: "verity_eth_port_settings", importer: func() (interface{}, error) { return i.importResource("ethportsettings") }},
		{name: "lags", terraformResourceType: "verity_lag", importer: func() (interface{}, error) { return i.importResource("lags") }},
		{name: "services", terraformResourceType: "verity_service", importer: func() (interface{}, error) { return i.importResource("services") }},
		{name: "bundles", terraformResourceType: "verity_bundle", importer: func() (interface{}, error) { return i.importResource("bundles") }},
	}

	// Filter tasks based on mode and API version compatibility
	var resourceTasks []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
	}

	for _, task := range allResourceTasks {
		if utils.IsResourceCompatibleWithMode(task.terraformResourceType, i.Mode) {
			resourceTasks = append(resourceTasks, task)
		} else {
			tflog.Info(i.ctx, "Skipping resource due to mode incompatibility", map[string]interface{}{
				"resource_name":           task.name,
				"terraform_resource_type": task.terraformResourceType,
				"mode":                    i.Mode,
			})
		}
	}

	for _, task := range resourceTasks {
		if !utils.IsResourceCompatibleWithVersion(task.terraformResourceType, apiVersionString) {
			tflog.Info(i.ctx, "Skipping resource due to API version incompatibility", map[string]interface{}{
				"resource_name":           task.name,
				"terraform_resource_type": task.terraformResourceType,
				"api_version":             apiVersionString,
			})
			continue
		}

		tflog.Info(i.ctx, "Importing resource compatible with API version", map[string]interface{}{
			"resource_name":           task.name,
			"terraform_resource_type": task.terraformResourceType,
			"api_version":             apiVersionString,
		})

		data, err := task.importer()
		if err != nil {
			tflog.Error(i.ctx, "Failed to import resource", map[string]interface{}{"resource_name": task.name, "error": err})
			return fmt.Errorf("failed to import %s: %w", task.name, err)
		}

		if data == nil {
			tflog.Info(i.ctx, "No data returned by importer, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}
		if m, ok := data.(map[string]map[string]interface{}); ok && len(m) == 0 {
			tflog.Info(i.ctx, "No data found for resource, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}

		// Get the resource config key from the terraform type
		resourceKey, ok := terraformTypeToResourceKey[task.terraformResourceType]
		if !ok {
			tflog.Error(i.ctx, "No resource config found for terraform type", map[string]interface{}{
				"resource_name":  task.name,
				"terraform_type": task.terraformResourceType,
			})
			return fmt.Errorf("no resource config found for %s", task.terraformResourceType)
		}

		tfConfig, err := i.generateResourceTFByName(resourceKey, data)
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

	return nil
}

func (i *Importer) importResource(resourceName string) (interface{}, error) {
	config, ok := importerRegistry[resourceName]
	if !ok {
		return nil, fmt.Errorf("no importer configuration found for %s", resourceName)
	}

	resp, err := config.apiCaller(i.ctx, i.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %v", resourceName, err)
	}
	defer resp.Body.Close()

	result := make(map[string]map[string]map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %v", resourceName, err)
	}

	data, ok := result[config.jsonKey]
	if !ok {
		// Return empty map if the key doesn't exist
		return make(map[string]map[string]interface{}), nil
	}

	return data, nil
}

func (i *Importer) generateResourceTFByName(resourceKey string, data interface{}) (string, error) {
	cfg, ok := resourceConfigs[resourceKey]
	if !ok {
		return "", fmt.Errorf("unknown resource type: %s", resourceKey)
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateResourceTF(data interface{}, config ResourceConfig) (string, error) {
	resourcesMap, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data format for resource type %s", config.ResourceType)
	}

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
		tfConfig.WriteString(fmt.Sprintf(config.HeaderNameLineFormat, name))
		tfConfig.WriteString(fmt.Sprintf(config.HeaderDependsOnLineFormat, config.StageName))

		// Skip object_properties section entirely if specified
		skipObjectProperties := false
		for _, key := range config.AdditionalTopLevelSkipKeys {
			if key == "object_properties" {
				skipObjectProperties = true
				break
			}
		}

		if !skipObjectProperties {
			// Check if object_properties exists in the API response
			objPropsRaw, objectPropertiesExists := resource["object_properties"]

			// Only include object_properties if it's actually present in the API response
			if objectPropertiesExists {
				tfConfig.WriteString("	object_properties")
				objProps, _ := objPropsRaw.(map[string]interface{})

				// Check if object_properties is empty (exists but has no fields)
				isEmptyObjectProps := len(objProps) == 0

				if isEmptyObjectProps && config.EmptyObjectPropsAsSingleLine {
					tfConfig.WriteString(" {}\n")
				} else if isEmptyObjectProps {
					// Case: object_properties exists but is empty
					tfConfig.WriteString(" {\n")
					tfConfig.WriteString("	}\n")
				} else {
					// Case: object_properties exists and has content
					tfConfig.WriteString(" {\n")
					var objPropsContentBuilder strings.Builder
					if config.ObjectPropsHandler != nil {
						config.ObjectPropsHandler(objProps, &objPropsContentBuilder, config)
					}
					tfConfig.WriteString(objPropsContentBuilder.String())
					tfConfig.WriteString("	}\n")
				}
			}
			// Case: object_properties doesn't exist in API response - don't include it at all
		}

		skipKeysSet := map[string]bool{
			"name":              true,
			"object_properties": true,
		}
		for _, key := range config.AdditionalTopLevelSkipKeys {
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
				// Check if it's a whole number
				if v == float64(int(v)) {
					tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", tfFieldName, int(v)))
				} else {
					tfConfig.WriteString(fmt.Sprintf("	%s = %g\n", tfFieldName, v))
				}
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", tfFieldName, formatValue(v)))
			case []interface{}:
				if _, isNestedBlock := config.NestedBlockFields[tfFieldName]; isNestedBlock {
					style, hasStyle := config.NestedBlockStyles[tfFieldName]
					if !hasStyle {
						style = NestedBlockIterationStyle{PrintIndexFirst: true, SkipIndexInMainLoop: true, IterateAllAsMap: false}
					}

					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", tfFieldName))

							if style.IterateAllAsMap {
								for itemKey, itemValue := range itemMap {
									tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
								}
							} else {
								printedIndex := false
								if style.PrintIndexFirst {
									if indexVal, idxExists := itemMap["index"]; idxExists {
										if indexFloat, isFloat := indexVal.(float64); isFloat {
											tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(indexFloat)))
											printedIndex = true
										}
									}
								}

								var nestedItemKeys []string
								for itemKey := range itemMap {
									if style.SkipIndexInMainLoop && itemKey == "index" && printedIndex {
										continue
									}
									if itemKey == "index" && style.SkipIndexInMainLoop && !printedIndex {
										if style.SkipIndexInMainLoop {
											continue
										}
									}
									nestedItemKeys = append(nestedItemKeys, itemKey)
								}
								sort.Strings(nestedItemKeys)

								for _, itemKey := range nestedItemKeys {
									itemValue := itemMap[itemKey]
									tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
								}
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

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	apiVersionString := i.getAPIVersion()
	tflog.Info(i.ctx, "Generating stages for mode and version", map[string]interface{}{
		"mode":    i.Mode,
		"version": apiVersionString,
	})

	// Define stage orderings for each mode
	type StageDefinition struct {
		StageName      string
		ResourceType   string
		DependsOnStage string // empty string means it's the first stage
	}

	// Order for 6.4 API DATACENTER:
	// 1. Tenants, 2. Gateways, 3. Gateway Profiles, 4. Services,
	// 5. Eth Port Profiles, 6. Eth Port Settings, 7. Lags, 8. Bundles
	stageOrder := []StageDefinition{
		{"tenant_stage", "verity_tenant", ""},
		{"gateway_stage", "verity_gateway", "tenant_stage"},
		{"gateway_profile_stage", "verity_gateway_profile", "gateway_stage"},
		{"service_stage", "verity_service", "gateway_profile_stage"},
		{"eth_port_profile_stage", "verity_eth_port_profile", "service_stage"},
		{"eth_port_settings_stage", "verity_eth_port_settings", "eth_port_profile_stage"},
		{"lag_stage", "verity_lag", "eth_port_settings_stage"},
		{"bundle_stage", "verity_bundle", "lag_stage"},
	}

	// Filter stages based on resource compatibility with mode and API version
	var compatibleStages []StageDefinition
	var lastCompatibleStage string

	for _, stage := range stageOrder {
		if utils.IsResourceCompatible(stage.ResourceType, i.Mode, apiVersionString) {
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
				"api_version":   apiVersionString,
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
		"api_version":       apiVersionString,
		"total_stages":      len(stageOrder),
		"compatible_stages": len(compatibleStages),
	})

	return tfConfig.String(), nil
}

// universalObjectPropsHandler dynamically processes all fields present in the object_properties
// section of the API response and converts them to Terraform configuration format.
// - If objProps is nil or empty map: generates no content
// - If objProps has fields: generates TF config for all fields (including nested structures)
// - Fields specified in ObjectPropsNestedBlockFields are rendered as blocks instead of attributes
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
				// Render as nested blocks
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
				// Render as attribute assignment
				builder.WriteString(fmt.Sprintf("		%s = %s\n", key, formatObjectPropsValue(value, "	")))
			}
		}
	}
	// If objProps is nil or empty, generate no content
}

func formatObjectPropsValue(value interface{}, indent string) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		// Check if it's a whole number
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
				// Handle array of objects
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
				// Handle array of primitives
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
		// Check if it's a whole number
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

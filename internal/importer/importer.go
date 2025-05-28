package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
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
	ObjectPropsHandler           func(objProps map[string]interface{}, builder *strings.Builder)
	NestedBlockFields            map[string]bool
	AdditionalTopLevelSkipKeys   []string
	EmptyObjectPropsAsSingleLine bool
	NestedBlockStyles            map[string]NestedBlockIterationStyle
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

func NewImporter(client *openapi.APIClient) *Importer {
	return &Importer{
		client: client,
		ctx:    context.Background(),
	}
}

// ImportAll fetches all resources and saves them as Terraform configuration files
func (i *Importer) ImportAll(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	stagesTF, err := i.generateStagesTF()
	if err != nil {
		return fmt.Errorf("failed to generate stages: %v", err)
	}

	stagesFile := filepath.Join(outputDir, "stages.tf")
	if err := os.WriteFile(stagesFile, []byte(stagesTF), 0644); err != nil {
		return fmt.Errorf("failed to write stages terraform config: %v", err)
	}

	resources := []struct {
		name        string
		importer    func() (interface{}, error)
		tfGenerator func(interface{}) (string, error)
	}{
		{
			name:        "tenants",
			importer:    i.importTenants,
			tfGenerator: i.generateTenantsTF,
		},
		{
			name:        "gateways",
			importer:    i.importGateways,
			tfGenerator: i.generateGatewaysTF,
		},
		{
			name:        "gatewayprofiles",
			importer:    i.importGatewayProfiles,
			tfGenerator: i.generateGatewayProfilesTF,
		},
		{
			name:        "ethportprofiles",
			importer:    i.importEthPortProfiles,
			tfGenerator: i.generateEthPortProfilesTF,
		},
		{
			name:        "lags",
			importer:    i.importLags,
			tfGenerator: i.generateLagsTF,
		},
		{
			name:        "services",
			importer:    i.importServices,
			tfGenerator: i.generateServicesTF,
		},
		{
			name:        "ethportsettings",
			importer:    i.importEthPortSettings,
			tfGenerator: i.generateEthPortSettingsTF,
		},
		{
			name:        "bundles",
			importer:    i.importBundles,
			tfGenerator: i.generateBundlesTF,
		},
	}

	for _, resource := range resources {
		data, err := resource.importer()
		if err != nil {
			return fmt.Errorf("failed to import %s: %v", resource.name, err)
		}

		tfConfig, err := resource.tfGenerator(data)
		if err != nil {
			return fmt.Errorf("failed to generate terraform config for %s: %v", resource.name, err)
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.tf", resource.name))
		if err := os.WriteFile(outputFile, []byte(tfConfig), 0644); err != nil {
			return fmt.Errorf("failed to write %s terraform config: %v", resource.name, err)
		}
	}

	return nil
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

		tfConfig.WriteString("	object_properties")
		var objPropsContentBuilder strings.Builder
		objProps, _ := resource["object_properties"].(map[string]interface{})

		if config.ObjectPropsHandler != nil {
			config.ObjectPropsHandler(objProps, &objPropsContentBuilder)
		}
		objPropsContent := objPropsContentBuilder.String()

		if objPropsContent == "" && config.EmptyObjectPropsAsSingleLine {
			tfConfig.WriteString(" {}\n")
		} else {
			tfConfig.WriteString(" {\n")
			tfConfig.WriteString(objPropsContent)
			tfConfig.WriteString("	}\n")
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
			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if _, isNestedBlock := config.NestedBlockFields[key]; isNestedBlock {
					style, hasStyle := config.NestedBlockStyles[key]
					if !hasStyle {
						style = NestedBlockIterationStyle{PrintIndexFirst: true, SkipIndexInMainLoop: true, IterateAllAsMap: false}
					}

					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))

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
					tfConfig.WriteString(fmt.Sprintf("	%s = [\n", key))
					for _, item := range v {
						if str, ok := item.(string); ok {
							tfConfig.WriteString(fmt.Sprintf("		%s,\n", formatValue(str)))
						}
					}
					tfConfig.WriteString("	]\n")
				}
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", key))
			}
		}
		tfConfig.WriteString("}\n\n")
	}
	return tfConfig.String(), nil
}

func (i *Importer) generateGatewayProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "gateway_profile",
		StageName:                 "gateway_profile_stage",
		HeaderNameLineFormat:      "\tname = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler: func(objProps map[string]interface{}, builder *strings.Builder) {
			if objProps != nil {
				if group, ok := objProps["group"].(string); ok && group != "" {
					builder.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
				} else {
					builder.WriteString("		group = \"\"\n")
				}
			} else {
				builder.WriteString("		group = \"\"\n")
			}
		},
		NestedBlockFields:          map[string]bool{"external_gateways": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateEthPortProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "eth_port_profile",
		StageName:                 "eth_port_profile_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler: func(objProps map[string]interface{}, builder *strings.Builder) {
			if objProps != nil {
				if group, ok := objProps["group"]; ok && group != nil {
					builder.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
				} else {
					builder.WriteString("		group = null\n")
				}
				if portMonitoring, ok := objProps["port_monitoring"]; ok && portMonitoring != nil {
					builder.WriteString(fmt.Sprintf("		port_monitoring = %s\n", formatValue(portMonitoring)))
				} else {
					builder.WriteString("		port_monitoring = null\n")
				}
			} else {
				builder.WriteString("		group = null\n")
				builder.WriteString("		port_monitoring = null\n")
			}
		},
		NestedBlockFields:          map[string]bool{"services": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateBundlesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "bundle",
		StageName:                 "bundle_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler: func(objProps map[string]interface{}, builder *strings.Builder) {
			if objProps != nil {
				if isForSwitch, ok := objProps["is_for_switch"].(bool); ok {
					builder.WriteString(fmt.Sprintf("		is_for_switch = %t\n", isForSwitch))
				} else {
					builder.WriteString("		is_for_switch = false\n")
				}
			} else {
				builder.WriteString("		is_for_switch = false\n")
			}
		},
		NestedBlockFields:          map[string]bool{"eth_port_paths": true, "user_services": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
		NestedBlockStyles: map[string]NestedBlockIterationStyle{
			"eth_port_paths": {IterateAllAsMap: true},
			"user_services":  {IterateAllAsMap: true},
		},
	}
	return i.generateResourceTF(data, cfg)
}

func commonObjectPropsStringGroupHandler(objProps map[string]interface{}, builder *strings.Builder) {
	if objProps != nil {
		if groupVal, keyExists := objProps["group"]; keyExists {
			if groupStr, ok := groupVal.(string); ok {
				builder.WriteString(fmt.Sprintf("		group = %s\n", formatValue(groupStr)))
			} else {
				builder.WriteString("		group = null\n")
			}
		} else {
			builder.WriteString("		group = null\n")
		}
	} else {
		builder.WriteString("		group = null\n")
	}
}

func (i *Importer) generateTenantsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "tenant",
		StageName:                 "tenant_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        commonObjectPropsStringGroupHandler,
		NestedBlockFields:         map[string]bool{"route_tenants": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateLagsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "lag",
		StageName:                    "lag_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           func(objProps map[string]interface{}, builder *strings.Builder) {},
		EmptyObjectPropsAsSingleLine: true,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateServicesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "service",
		StageName:                 "service_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        commonObjectPropsStringGroupHandler,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateEthPortSettingsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "eth_port_settings",
		StageName:                 "eth_port_settings_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        commonObjectPropsStringGroupHandler,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateGatewaysTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "gateway",
		StageName:                 "gateway_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "\tdepends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        commonObjectPropsStringGroupHandler,
		NestedBlockFields:         map[string]bool{"static_routes": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	tfConfig.WriteString(`
# These resources establish ordering for bulk operations
resource "verity_operation_stage" "tenant_stage" {
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "gateway_stage" {
  depends_on = [verity_operation_stage.tenant_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "gateway_profile_stage" {
  depends_on = [verity_operation_stage.gateway_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "service_stage" {
  depends_on = [verity_operation_stage.gateway_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_profile_stage" {
  depends_on = [verity_operation_stage.service_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "eth_port_settings_stage" {
  depends_on = [verity_operation_stage.eth_port_profile_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "lag_stage" {
  depends_on = [verity_operation_stage.eth_port_settings_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "bundle_stage" {
  depends_on = [verity_operation_stage.lag_stage]
  lifecycle {
    create_before_destroy = true
  }
}
`)

	return tfConfig.String(), nil
}

func (i *Importer) importEthPortProfiles() (interface{}, error) {
	resp, err := i.client.EthPortProfilesAPI.EthportprofilesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get ethport profiles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		EthPort map[string]map[string]interface{} `json:"eth_port_profile_"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ethport profiles response: %v", err)
	}

	return result.EthPort, nil
}

func (i *Importer) importTenants() (interface{}, error) {
	resp, err := i.client.TenantsAPI.TenantsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Tenant map[string]map[string]interface{} `json:"tenant"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode tenants response: %v", err)
	}

	return result.Tenant, nil
}

func (i *Importer) importGatewayProfiles() (interface{}, error) {
	resp, err := i.client.GatewayProfilesAPI.GatewayprofilesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway profiles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		GatewayProfile map[string]map[string]interface{} `json:"gateway_profile"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode gateway profiles response: %v", err)
	}

	return result.GatewayProfile, nil
}

func (i *Importer) importLags() (interface{}, error) {
	resp, err := i.client.LAGsAPI.LagsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get lags: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		LaggGroup map[string]map[string]interface{} `json:"lag"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode lags response: %v", err)
	}

	return result.LaggGroup, nil
}

func (i *Importer) importServices() (interface{}, error) {
	resp, err := i.client.ServicesAPI.ServicesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Service map[string]map[string]interface{} `json:"service"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode services response: %v", err)
	}

	return result.Service, nil
}

func (i *Importer) importEthPortSettings() (interface{}, error) {
	resp, err := i.client.EthPortSettingsAPI.EthportsettingsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get ethport settings: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		EthPortSetting map[string]map[string]interface{} `json:"eth_port_settings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode eth port settings response: %v", err)
	}

	return result.EthPortSetting, nil
}

func (i *Importer) importBundles() (interface{}, error) {
	resp, err := i.client.BundlesAPI.BundlesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get bundles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Bundle map[string]map[string]interface{} `json:"endpoint_bundle"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode bundles response: %v", err)
	}

	return result.Bundle, nil
}

func (i *Importer) importGateways() (interface{}, error) {
	resp, err := i.client.GatewaysAPI.GatewaysGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateways: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Gateway map[string]map[string]interface{} `json:"gateway"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode gateways response: %v", err)
	}

	return result.Gateway, nil
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		return fmt.Sprintf("%d", int(v))
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

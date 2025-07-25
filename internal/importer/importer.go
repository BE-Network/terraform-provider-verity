package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ImportAll fetches all resources and saves them as Terraform configuration files
func (i *Importer) ImportAll(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	tflog.Info(i.ctx, "Starting importer with mode", map[string]interface{}{
		"mode": i.Mode,
	})

	var apiVersionString string
	defaultVersion := "6.4"

	versionResp, err := i.client.VersionAPI.VersionGet(i.ctx).Execute()
	if err != nil {
		tflog.Error(i.ctx, "Failed to execute API version GET request, defaulting version", map[string]interface{}{"error": err, "default_version": defaultVersion})
		apiVersionString = defaultVersion
	} else {
		defer versionResp.Body.Close()
		versionBodyBytes, err := io.ReadAll(versionResp.Body)
		if err != nil {
			tflog.Error(i.ctx, "Failed to read API version response body, defaulting version", map[string]interface{}{"error": err, "default_version": defaultVersion})
			apiVersionString = defaultVersion
		} else if versionResp.StatusCode != http.StatusOK {
			tflog.Error(i.ctx, "API version GET request failed with non-OK status, defaulting version", map[string]interface{}{
				"status_code":     versionResp.StatusCode,
				"body":            string(versionBodyBytes),
				"default_version": defaultVersion,
			})
			apiVersionString = defaultVersion
		} else {
			var directVersion string
			if err := json.Unmarshal(versionBodyBytes, &directVersion); err == nil && directVersion != "" {
				tflog.Info(i.ctx, "Successfully parsed API version as direct string", map[string]interface{}{
					"version": directVersion,
				})
				apiVersionString = directVersion
			} else {
				var versionData struct {
					Version string `json:"version"`
				}
				if err := json.Unmarshal(versionBodyBytes, &versionData); err != nil {
					tflog.Error(i.ctx, "Failed to decode API version JSON response, defaulting version", map[string]interface{}{
						"error":           err,
						"body":            string(versionBodyBytes),
						"default_version": defaultVersion,
					})
					apiVersionString = defaultVersion
				} else if versionData.Version == "" {
					tflog.Warn(i.ctx, "API version string is empty in response, defaulting version", map[string]interface{}{"body": string(versionBodyBytes), "default_version": defaultVersion})
					apiVersionString = defaultVersion
				} else {
					apiVersionString = versionData.Version
				}
			}
		}
	}

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
		tfGenerator           func(interface{}) (string, error)
	}{
		{name: "tenants", terraformResourceType: "verity_tenant", importer: i.importTenants, tfGenerator: i.generateTenantsTF},
		{name: "gateways", terraformResourceType: "verity_gateway", importer: i.importGateways, tfGenerator: i.generateGatewaysTF},
		{name: "gatewayprofiles", terraformResourceType: "verity_gateway_profile", importer: i.importGatewayProfiles, tfGenerator: i.generateGatewayProfilesTF},
		{name: "ethportprofiles", terraformResourceType: "verity_eth_port_profile", importer: i.importEthPortProfiles, tfGenerator: i.generateEthPortProfilesTF},
		{name: "lags", terraformResourceType: "verity_lag", importer: i.importLags, tfGenerator: i.generateLagsTF},
		{name: "services", terraformResourceType: "verity_service", importer: i.importServices, tfGenerator: i.generateServicesTF},
		{name: "ethportsettings", terraformResourceType: "verity_eth_port_settings", importer: i.importEthPortSettings, tfGenerator: i.generateEthPortSettingsTF},
		{name: "bundles", terraformResourceType: "verity_bundle", importer: i.importBundles, tfGenerator: i.generateBundlesTF},
		{name: "acls_ipv4", terraformResourceType: "verity_acl_v4", importer: i.importACLsIPv4, tfGenerator: i.generateACLsIPv4TF},
		{name: "acls_ipv6", terraformResourceType: "verity_acl_v6", importer: i.importACLsIPv6, tfGenerator: i.generateACLsIPv6TF},
		{name: "badges", terraformResourceType: "verity_badge", importer: i.importBadges, tfGenerator: i.generateBadgesTF},
		{name: "authenticatedethports", terraformResourceType: "verity_authenticated_eth_port", importer: i.importAuthenticatedEthPorts, tfGenerator: i.generateAuthenticatedEthPortsTF},
		{name: "devicecontrollers", terraformResourceType: "verity_device_controller", importer: i.importDeviceControllers, tfGenerator: i.generateDeviceControllersTF},
		{name: "devicevoicesettings", terraformResourceType: "verity_device_voice_settings", importer: i.importDeviceVoiceSettings, tfGenerator: i.generateDeviceVoiceSettingsTF},
		{name: "packetbroker", terraformResourceType: "verity_packet_broker", importer: i.importPacketBroker, tfGenerator: i.generatePacketBrokerTF},
		{name: "packetqueues", terraformResourceType: "verity_packet_queue", importer: i.importPacketQueues, tfGenerator: i.generatePacketQueuesTF},
		{name: "serviceportprofiles", terraformResourceType: "verity_service_port_profile", importer: i.importServicePortProfiles, tfGenerator: i.generateServicePortProfilesTF},
		{name: "voiceportprofiles", terraformResourceType: "verity_voice_port_profile", importer: i.importVoicePortProfiles, tfGenerator: i.generateVoicePortProfilesTF},
		{name: "switchpoints", terraformResourceType: "verity_switchpoint", importer: i.importSwitchpoints, tfGenerator: i.generateSwitchpointsTF},
	}

	// Filter tasks based on mode and API version compatibility
	var resourceTasks []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
		tfGenerator           func(interface{}) (string, error)
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

		tfConfig, err := task.tfGenerator(data)
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
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", tfFieldName, int(v)))
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

func (i *Importer) generateGatewayProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:               "gateway_profile",
		StageName:                  "gateway_profile_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"external_gateways": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateEthPortProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:               "eth_port_profile",
		StageName:                  "eth_port_profile_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"services": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateBundlesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:               "bundle",
		StageName:                  "bundle_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		NestedBlockFields:          map[string]bool{"eth_port_paths": true, "user_services": true, "rg_services": true, "voice_port_profile_paths": true},
		AdditionalTopLevelSkipKeys: []string{"index"},
		NestedBlockStyles: map[string]NestedBlockIterationStyle{
			"eth_port_paths":           {IterateAllAsMap: true},
			"user_services":            {IterateAllAsMap: true},
			"rg_services":              {IterateAllAsMap: true},
			"voice_port_profile_paths": {IterateAllAsMap: true},
		},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateTenantsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "tenant",
		StageName:                 "tenant_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"route_tenants": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateLagsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "lag",
		StageName:                    "lag_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateServicesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "service",
		StageName:                 "service_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateEthPortSettingsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "eth_port_settings",
		StageName:                 "eth_port_settings_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"lldp_med": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateGatewaysTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "gateway",
		StageName:                 "gateway_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"static_routes": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateBadgesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:               "badge",
		StageName:                  "badge_stage",
		HeaderNameLineFormat:       "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:  "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:         universalObjectPropsHandler,
		AdditionalTopLevelSkipKeys: []string{},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateAuthenticatedEthPortsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "authenticated_eth_port",
		StageName:                 "authenticated_eth_port_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"eth_ports": true, "object_properties": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateDeviceVoiceSettingsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "device_voice_settings",
		StageName:                    "device_voice_setting_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
		NestedBlockFields:            map[string]bool{"codecs": true},
		FieldMappings:                map[string]string{"Codecs": "codecs"},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generatePacketBrokerTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "packet_broker",
		StageName:                 "packet_broker_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"ipv4_permit": true, "ipv4_deny": true, "ipv6_permit": true, "ipv6_deny": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generatePacketQueuesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "packet_queue",
		StageName:                 "packet_queue_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"pbit": true, "queue": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateServicePortProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "service_port_profile",
		StageName:                 "service_port_profile_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
		NestedBlockFields:         map[string]bool{"services": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateVoicePortProfilesTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "voice_port_profile",
		StageName:                    "voice_port_profile_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: false,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateSwitchpointsTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "switchpoint",
		StageName:                    "switchpoint_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		NestedBlockFields:            map[string]bool{"badges": true, "children": true, "traffic_mirrors": true, "eths": true},
		ObjectPropsNestedBlockFields: map[string]bool{"eths": true},
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateDeviceControllersTF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:                 "device_controller",
		StageName:                    "device_controller_stage",
		HeaderNameLineFormat:         "    name = \"%s\"\n",
		HeaderDependsOnLineFormat:    "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:           universalObjectPropsHandler,
		EmptyObjectPropsAsSingleLine: true,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateACLsIPv4TF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "acl_v4",
		StageName:                 "acl_v4_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateACLsIPv6TF(data interface{}) (string, error) {
	cfg := ResourceConfig{
		ResourceType:              "acl_v6",
		StageName:                 "acl_v6_stage",
		HeaderNameLineFormat:      "    name = \"%s\"\n",
		HeaderDependsOnLineFormat: "    depends_on = [verity_operation_stage.%s]\n",
		ObjectPropsHandler:        universalObjectPropsHandler,
	}
	return i.generateResourceTF(data, cfg)
}

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	var apiVersionString string
	defaultVersion := "6.4" // Default to 6.4

	versionResp, err := i.client.VersionAPI.VersionGet(i.ctx).Execute()
	if err != nil {
		tflog.Warn(i.ctx, "Failed to get API version for stages, using default", map[string]interface{}{
			"error":   err,
			"default": defaultVersion,
		})
		apiVersionString = defaultVersion
	} else {
		defer versionResp.Body.Close()
		versionBodyBytes, err := io.ReadAll(versionResp.Body)
		if err != nil {
			apiVersionString = defaultVersion
		} else {

			var directVersion string
			if err := json.Unmarshal(versionBodyBytes, &directVersion); err == nil && directVersion != "" {
				apiVersionString = directVersion
			} else {
				var versionData struct {
					Version string `json:"version"`
				}
				if err := json.Unmarshal(versionBodyBytes, &versionData); err != nil {
					apiVersionString = defaultVersion
				} else if versionData.Version == "" {
					apiVersionString = defaultVersion
				} else {
					apiVersionString = versionData.Version
				}
			}
		}
	}

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

	var stageOrder []StageDefinition

	if i.Mode == "campus" {
		// CAMPUS mode staging order:
		// 1. Services, 2. Eth Port Profiles, 3. Authenticated Eth-Ports, 4. Device Voice Settings,
		// 5. Packet Queues, 6. Service Port Profiles, 7. Voice-Port Profiles, 8. Eth Port Settings,
		// 9. Lags, 10. Bundles, 11. Badges, 12. Switchpoints, 13. Device controllers
		stageOrder = []StageDefinition{
			{"service_stage", "verity_service", ""},
			{"eth_port_profile_stage", "verity_eth_port_profile", "service_stage"},
			{"authenticated_eth_port_stage", "verity_authenticated_eth_port", "eth_port_profile_stage"},
			{"device_voice_setting_stage", "verity_device_voice_settings", "authenticated_eth_port_stage"},
			{"packet_queue_stage", "verity_packet_queue", "device_voice_setting_stage"},
			{"service_port_profile_stage", "verity_service_port_profile", "packet_queue_stage"},
			{"voice_port_profile_stage", "verity_voice_port_profile", "service_port_profile_stage"},
			{"eth_port_settings_stage", "verity_eth_port_settings", "voice_port_profile_stage"},
			{"lag_stage", "verity_lag", "eth_port_settings_stage"},
			{"bundle_stage", "verity_bundle", "lag_stage"},
			{"badge_stage", "verity_badge", "bundle_stage"},
			{"switchpoint_stage", "verity_switchpoint", "badge_stage"},
			{"device_controller_stage", "verity_device_controller", "switchpoint_stage"},
		}
	} else {
		// DATACENTER mode staging order:
		// 1. Tenants, 2. Gateways, 3. Gateway Profiles, 4. Services, 5. Packet Queues,
		// 6. Eth Port Profiles, 7. Eth Port Settings, 8. Lags, 9. Bundles,
		// 10. IPv4 ACLs, 11. IPv6 ACLs, 12. PacketBroker, 13. Badges, 14. Switchpoints, 15. Device controllers
		stageOrder = []StageDefinition{
			{"tenant_stage", "verity_tenant", ""},
			{"gateway_stage", "verity_gateway", "tenant_stage"},
			{"gateway_profile_stage", "verity_gateway_profile", "gateway_stage"},
			{"service_stage", "verity_service", "gateway_profile_stage"},
			{"packet_queue_stage", "verity_packet_queue", "service_stage"},
			{"eth_port_profile_stage", "verity_eth_port_profile", "packet_queue_stage"},
			{"eth_port_settings_stage", "verity_eth_port_settings", "eth_port_profile_stage"},
			{"lag_stage", "verity_lag", "eth_port_settings_stage"},
			{"bundle_stage", "verity_bundle", "lag_stage"},
			{"acl_v4_stage", "verity_acl_v4", "bundle_stage"},
			{"acl_v6_stage", "verity_acl_v6", "acl_v4_stage"},
			{"packet_broker_stage", "verity_packet_broker", "acl_v6_stage"},
			{"badge_stage", "verity_badge", "packet_broker_stage"},
			{"switchpoint_stage", "verity_switchpoint", "badge_stage"},
			{"device_controller_stage", "verity_device_controller", "switchpoint_stage"},
		}
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

func (i *Importer) importBadges() (interface{}, error) {
	resp, err := i.client.BadgesAPI.BadgesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get badges: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Badge map[string]map[string]interface{} `json:"badge"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode badges response: %v", err)
	}

	return result.Badge, nil
}

func (i *Importer) importAuthenticatedEthPorts() (interface{}, error) {
	resp, err := i.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated eth ports: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		AuthenticatedEthPort map[string]map[string]interface{} `json:"authenticated_eth_port"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode authenticated eth ports response: %v", err)
	}

	return result.AuthenticatedEthPort, nil
}

func (i *Importer) importDeviceVoiceSettings() (interface{}, error) {
	resp, err := i.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get device voice settings: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		DeviceVoiceSettings map[string]map[string]interface{} `json:"device_voice_settings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode device voice settings response: %v", err)
	}

	return result.DeviceVoiceSettings, nil
}

func (i *Importer) importPacketBroker() (interface{}, error) {
	resp, err := i.client.PacketBrokerAPI.PacketbrokerGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get packet broker: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		PacketBroker map[string]map[string]interface{} `json:"pb_egress_profile"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode packet broker response: %v", err)
	}

	return result.PacketBroker, nil
}

func (i *Importer) importPacketQueues() (interface{}, error) {
	resp, err := i.client.PacketQueuesAPI.PacketqueuesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get packet queues: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		PacketQueue map[string]map[string]interface{} `json:"packet_queue"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode packet queues response: %v", err)
	}

	return result.PacketQueue, nil
}

func (i *Importer) importServicePortProfiles() (interface{}, error) {
	resp, err := i.client.ServicePortProfilesAPI.ServiceportprofilesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get service port profiles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ServicePortProfile map[string]map[string]interface{} `json:"service_port_profile"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode service port profiles response: %v", err)
	}

	return result.ServicePortProfile, nil
}

func (i *Importer) importVoicePortProfiles() (interface{}, error) {
	resp, err := i.client.VoicePortProfilesAPI.VoiceportprofilesGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get voice port profiles: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		VoicePortProfile map[string]map[string]interface{} `json:"voice_port_profiles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode voice port profiles response: %v", err)
	}

	return result.VoicePortProfile, nil
}

func (i *Importer) importSwitchpoints() (interface{}, error) {
	resp, err := i.client.SwitchpointsAPI.SwitchpointsGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get switchpoints: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Switchpoint map[string]map[string]interface{} `json:"switchpoint"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode switchpoints response: %v", err)
	}

	return result.Switchpoint, nil
}

func (i *Importer) importDeviceControllers() (interface{}, error) {
	resp, err := i.client.DeviceControllersAPI.DevicecontrollersGet(i.ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get device controllers: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		DeviceController map[string]map[string]interface{} `json:"device_controller"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode device controllers response: %v", err)
	}

	return result.DeviceController, nil
}

func (i *Importer) importACLsIPv4() (interface{}, error) {
	return i.importACLs("4")
}

func (i *Importer) importACLsIPv6() (interface{}, error) {
	return i.importACLs("6")
}

func (i *Importer) importACLs(ipVersion string) (map[string]map[string]interface{}, error) {
	resp, err := i.client.ACLsAPI.AclsGet(i.ctx).IpVersion(ipVersion).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s ACLs: %v", ipVersion, err)
	}
	defer resp.Body.Close()

	if ipVersion == "4" {
		var result struct {
			IpFilter map[string]map[string]interface{} `json:"ipv4_filter"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode IPv4 ACLs response: %v", err)
		}

		if result.IpFilter == nil {
			return make(map[string]map[string]interface{}), nil
		}

		return result.IpFilter, nil
	} else {
		var result struct {
			IpFilter map[string]map[string]interface{} `json:"ipv6_filter"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode IPv6 ACLs response: %v", err)
		}

		if result.IpFilter == nil {
			return make(map[string]map[string]interface{}), nil
		}

		return result.IpFilter, nil
	}
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
		return fmt.Sprintf("%d", int(v))
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

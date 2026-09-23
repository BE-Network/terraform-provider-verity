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
	config := ResourceConfig{
		ResourceType:                 strings.TrimPrefix(terraformType, "verity_"),
		StageName:                    stageNameFor(i.Mode, terraformType),
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

func stageNameFor(mode, terraformType string) string {
	for _, stage := range stageOrder(mode) {
		if stage.ResourceType == terraformType {
			return stage.StageName
		}
	}
	return ""
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

	allResourceTasks := []struct {
		name                  string
		terraformResourceType string
		importer              func() (interface{}, error)
	}{
		{name: "tenants", terraformResourceType: "verity_tenant", importer: func() (interface{}, error) { return i.importResource("tenants") }},
		{name: "gateways", terraformResourceType: "verity_gateway", importer: func() (interface{}, error) { return i.importResource("gateways") }},
		{name: "gatewayprofiles", terraformResourceType: "verity_gateway_profile", importer: func() (interface{}, error) { return i.importResource("gatewayprofiles") }},
		{name: "deviceaaaprofiles", terraformResourceType: "verity_aaa_profile", importer: func() (interface{}, error) { return i.importResource("deviceaaaprofiles") }},
		{name: "ldapprofiles", terraformResourceType: "verity_ldap_profile", importer: func() (interface{}, error) { return i.importResource("ldapprofiles") }},
		{name: "ethportprofiles", terraformResourceType: "verity_eth_port_profile", importer: func() (interface{}, error) { return i.importResource("ethportprofiles") }},
		{name: "lags", terraformResourceType: "verity_lag", importer: func() (interface{}, error) { return i.importResource("lags") }},
		{name: "sflowcollectors", terraformResourceType: "verity_sflow_collector", importer: func() (interface{}, error) { return i.importResource("sflowcollectors") }},
		{name: "diagnosticsprofiles", terraformResourceType: "verity_diagnostics_profile", importer: func() (interface{}, error) { return i.importResource("diagnosticsprofiles") }},
		{name: "diagnosticsportprofiles", terraformResourceType: "verity_diagnostics_port_profile", importer: func() (interface{}, error) { return i.importResource("diagnosticsportprofiles") }},
		{name: "policybasedroutingacl", terraformResourceType: "verity_pb_routing_acl", importer: func() (interface{}, error) { return i.importResource("policybasedroutingacl") }},
		{name: "policybasedrouting", terraformResourceType: "verity_pb_routing", importer: func() (interface{}, error) { return i.importResource("policybasedrouting") }},
		{name: "services", terraformResourceType: "verity_service", importer: func() (interface{}, error) { return i.importResource("services") }},
		{name: "ethportsettings", terraformResourceType: "verity_eth_port_settings", importer: func() (interface{}, error) { return i.importResource("ethportsettings") }},
		{name: "bundles", terraformResourceType: "verity_bundle", importer: func() (interface{}, error) { return i.importResource("bundles") }},
		{name: "acls_ipv4", terraformResourceType: "verity_acl_v4", importer: i.importACLsIPv4},
		{name: "acls_ipv6", terraformResourceType: "verity_acl_v6", importer: i.importACLsIPv6},
		{name: "badges", terraformResourceType: "verity_badge", importer: func() (interface{}, error) { return i.importResource("badges") }},
		{name: "authenticatedethports", terraformResourceType: "verity_authenticated_eth_port", importer: func() (interface{}, error) { return i.importResource("authenticatedethports") }},
		{name: "devicevoicesettings", terraformResourceType: "verity_device_voice_settings", importer: func() (interface{}, error) { return i.importResource("devicevoicesettings") }},
		{name: "packetbroker", terraformResourceType: "verity_packet_broker", importer: func() (interface{}, error) { return i.importResource("packetbroker") }},
		{name: "packetqueues", terraformResourceType: "verity_packet_queue", importer: func() (interface{}, error) { return i.importResource("packetqueues") }},
		{name: "tacacsprofiles", terraformResourceType: "verity_tacacs_profile", importer: func() (interface{}, error) { return i.importResource("tacacsprofiles") }},
		{name: "serviceportprofiles", terraformResourceType: "verity_service_port_profile", importer: func() (interface{}, error) { return i.importResource("serviceportprofiles") }},
		{name: "voiceportprofiles", terraformResourceType: "verity_voice_port_profile", importer: func() (interface{}, error) { return i.importResource("voiceportprofiles") }},
		{name: "spineplanes", terraformResourceType: "verity_spine_plane", importer: func() (interface{}, error) { return i.importResource("spineplanes") }},
		{name: "switchpoints", terraformResourceType: "verity_switchpoint", importer: func() (interface{}, error) { return i.importResource("switchpoints") }},
		{name: "aspathaccesslists", terraformResourceType: "verity_as_path_access_list", importer: func() (interface{}, error) { return i.importResource("aspathaccesslists") }},
		{name: "communitylists", terraformResourceType: "verity_community_list", importer: func() (interface{}, error) { return i.importResource("communitylists") }},
		{name: "macfilters", terraformResourceType: "verity_mac_filter", importer: func() (interface{}, error) { return i.importResource("macfilters") }},
		{name: "devicesettings", terraformResourceType: "verity_device_settings", importer: func() (interface{}, error) { return i.importResource("devicesettings") }},
		{name: "extendedcommunitylists", terraformResourceType: "verity_extended_community_list", importer: func() (interface{}, error) { return i.importResource("extendedcommunitylists") }},
		{name: "ipv4lists", terraformResourceType: "verity_ipv4_list", importer: func() (interface{}, error) { return i.importResource("ipv4lists") }},
		{name: "ipv4prefixlists", terraformResourceType: "verity_ipv4_prefix_list", importer: func() (interface{}, error) { return i.importResource("ipv4prefixlists") }},
		{name: "ipv6lists", terraformResourceType: "verity_ipv6_list", importer: func() (interface{}, error) { return i.importResource("ipv6lists") }},
		{name: "ipv6prefixlists", terraformResourceType: "verity_ipv6_prefix_list", importer: func() (interface{}, error) { return i.importResource("ipv6prefixlists") }},
		{name: "routemapclauses", terraformResourceType: "verity_route_map_clause", importer: func() (interface{}, error) { return i.importResource("routemapclauses") }},
		{name: "routemaps", terraformResourceType: "verity_route_map", importer: func() (interface{}, error) { return i.importResource("routemaps") }},
		{name: "sfpbreakouts", terraformResourceType: "verity_sfp_breakout", importer: func() (interface{}, error) { return i.importResource("sfpbreakouts") }},
		{name: "fabrics", terraformResourceType: "verity_fabric", importer: func() (interface{}, error) { return i.importResource("fabrics") }},
		{name: "planes", terraformResourceType: "verity_plane", importer: func() (interface{}, error) { return i.importResource("planes") }},
		{name: "racks", terraformResourceType: "verity_rack", importer: func() (interface{}, error) { return i.importResource("racks") }},
		{name: "pairs", terraformResourceType: "verity_pair", importer: func() (interface{}, error) { return i.importResource("pairs") }},
		{name: "pods", terraformResourceType: "verity_pod", importer: func() (interface{}, error) { return i.importResource("pods") }},
		{name: "sspgroups", terraformResourceType: "verity_ssp_group", importer: func() (interface{}, error) { return i.importResource("sspgroups") }},
		{name: "sus", terraformResourceType: "verity_su", importer: func() (interface{}, error) { return i.importResource("sus") }},
		{name: "portacls", terraformResourceType: "verity_port_acl", importer: func() (interface{}, error) { return i.importResource("portacls") }},
		{name: "groupingrules", terraformResourceType: "verity_grouping_rule", importer: func() (interface{}, error) { return i.importResource("groupingrules") }},
		{name: "thresholdgroups", terraformResourceType: "verity_threshold_group", importer: func() (interface{}, error) { return i.importResource("thresholdgroups") }},
		{name: "thresholds", terraformResourceType: "verity_threshold", importer: func() (interface{}, error) { return i.importResource("thresholds") }},
	}

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
		tflog.Info(i.ctx, "Importing resource", map[string]interface{}{
			"resource_name":           task.name,
			"terraform_resource_type": task.terraformResourceType,
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

		config, err := i.resourceConfig(task.terraformResourceType)
		if err != nil {
			tflog.Error(i.ctx, "No registry entry for terraform type", map[string]interface{}{
				"resource_name":  task.name,
				"terraform_type": task.terraformResourceType,
				"error":          err,
			})
			return fmt.Errorf("no registry entry for %s: %w", task.terraformResourceType, err)
		}

		if objects, ok := data.(map[string]map[string]interface{}); ok {
			i.PruneUnsupported(task.terraformResourceType, objects)
		}

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

func (i *Importer) importResource(resourceName string) (interface{}, error) {
	collectionKey := utils.ResponseCollectionKeyForEndpoint(resourceName)
	if collectionKey == "" {
		return nil, fmt.Errorf("the resource registry describes no endpoint %q", resourceName)
	}
	return i.fetch(resourceName, "/"+resourceName, nil, collectionKey)
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
}

func stageOrder(mode string) []stageDefinition {
	if mode == "campus" {
		return []stageDefinition{
			{"sfp_breakout_stage", "verity_sfp_breakout", ""},
			{"acl_v6_stage", "verity_acl_v6", "sfp_breakout_stage"},
			{"acl_v4_stage", "verity_acl_v4", "acl_v6_stage"},
			{"mac_filter_stage", "verity_mac_filter", "acl_v4_stage"},
			{"service_stage", "verity_service", "mac_filter_stage"},
			{"port_acl_stage", "verity_port_acl", "service_stage"},
			{"tacacs_profile_stage", "verity_tacacs_profile", "port_acl_stage"},
			{"ldap_profile_stage", "verity_ldap_profile", "tacacs_profile_stage"},
			{"sflow_collector_stage", "verity_sflow_collector", "ldap_profile_stage"},
			{"eth_port_profile_stage", "verity_eth_port_profile", "sflow_collector_stage"},
			{"packet_queue_stage", "verity_packet_queue", "eth_port_profile_stage"},
			{"device_aaa_profile_stage", "verity_aaa_profile", "packet_queue_stage"},
			{"fabric_stage", "verity_fabric", "device_aaa_profile_stage"},
			{"service_port_profile_stage", "verity_service_port_profile", "fabric_stage"},
			{"diagnostics_profile_stage", "verity_diagnostics_profile", "service_port_profile_stage"},
			{"authenticated_eth_port_stage", "verity_authenticated_eth_port", "diagnostics_profile_stage"},
			{"device_settings_stage", "verity_device_settings", "authenticated_eth_port_stage"},
			{"voice_port_profile_stage", "verity_voice_port_profile", "device_settings_stage"},
			{"lag_stage", "verity_lag", "voice_port_profile_stage"},
			{"device_voice_setting_stage", "verity_device_voice_settings", "lag_stage"},
			{"eth_port_settings_stage", "verity_eth_port_settings", "device_voice_setting_stage"},
			{"diagnostics_port_profile_stage", "verity_diagnostics_port_profile", "eth_port_settings_stage"},
			{"bundle_stage", "verity_bundle", "diagnostics_port_profile_stage"},
			{"badge_stage", "verity_badge", "bundle_stage"},
			{"grouping_rule_stage", "verity_grouping_rule", "badge_stage"},
			{"switchpoint_stage", "verity_switchpoint", "grouping_rule_stage"},
			{"threshold_stage", "verity_threshold", "switchpoint_stage"},
			{"threshold_group_stage", "verity_threshold_group", "threshold_stage"},
			{"pair_stage", "verity_pair", "threshold_group_stage"},
		}
	}
	return []stageDefinition{
		{"sfp_breakout_stage", "verity_sfp_breakout", ""},
		{"community_list_stage", "verity_community_list", "sfp_breakout_stage"},
		{"as_path_access_list_stage", "verity_as_path_access_list", "community_list_stage"},
		{"ipv6_prefix_list_stage", "verity_ipv6_prefix_list", "as_path_access_list_stage"},
		{"ipv4_prefix_list_stage", "verity_ipv4_prefix_list", "ipv6_prefix_list_stage"},
		{"extended_community_list_stage", "verity_extended_community_list", "ipv4_prefix_list_stage"},
		{"acl_v6_stage", "verity_acl_v6", "extended_community_list_stage"},
		{"acl_v4_stage", "verity_acl_v4", "acl_v6_stage"},
		{"route_map_clause_stage", "verity_route_map_clause", "acl_v4_stage"},
		{"pb_routing_acl_stage", "verity_pb_routing_acl", "route_map_clause_stage"},
		{"route_map_stage", "verity_route_map", "pb_routing_acl_stage"},
		{"pb_routing_stage", "verity_pb_routing", "route_map_stage"},
		{"tenant_stage", "verity_tenant", "pb_routing_stage"},
		{"service_stage", "verity_service", "tenant_stage"},
		{"fabric_stage", "verity_fabric", "service_stage"},
		{"tacacs_profile_stage", "verity_tacacs_profile", "fabric_stage"},
		{"ldap_profile_stage", "verity_ldap_profile", "tacacs_profile_stage"},
		{"port_acl_stage", "verity_port_acl", "ldap_profile_stage"},
		{"ipv6_list_stage", "verity_ipv6_list", "port_acl_stage"},
		{"ipv4_list_stage", "verity_ipv4_list", "ipv6_list_stage"},
		{"pod_stage", "verity_pod", "ipv4_list_stage"},
		{"packet_queue_stage", "verity_packet_queue", "pod_stage"},
		{"device_aaa_profile_stage", "verity_aaa_profile", "packet_queue_stage"},
		{"eth_port_profile_stage", "verity_eth_port_profile", "device_aaa_profile_stage"},
		{"packet_broker_stage", "verity_packet_broker", "eth_port_profile_stage"},
		{"sflow_collector_stage", "verity_sflow_collector", "packet_broker_stage"},
		{"gateway_stage", "verity_gateway", "sflow_collector_stage"},
		{"su_stage", "verity_su", "gateway_stage"},
		{"diagnostics_port_profile_stage", "verity_diagnostics_port_profile", "su_stage"},
		{"device_settings_stage", "verity_device_settings", "diagnostics_port_profile_stage"},
		{"lag_stage", "verity_lag", "device_settings_stage"},
		{"diagnostics_profile_stage", "verity_diagnostics_profile", "lag_stage"},
		{"gateway_profile_stage", "verity_gateway_profile", "diagnostics_profile_stage"},
		{"eth_port_settings_stage", "verity_eth_port_settings", "gateway_profile_stage"},
		{"badge_stage", "verity_badge", "eth_port_settings_stage"},
		{"plane_stage", "verity_plane", "badge_stage"},
		{"spine_plane_stage", "verity_spine_plane", "plane_stage"},
		{"rack_stage", "verity_rack", "spine_plane_stage"},
		{"bundle_stage", "verity_bundle", "rack_stage"},
		{"ssp_group_stage", "verity_ssp_group", "bundle_stage"},
		{"grouping_rule_stage", "verity_grouping_rule", "ssp_group_stage"},
		{"switchpoint_stage", "verity_switchpoint", "grouping_rule_stage"},
		{"threshold_stage", "verity_threshold", "switchpoint_stage"},
		{"threshold_group_stage", "verity_threshold_group", "threshold_stage"},
		{"pair_stage", "verity_pair", "threshold_group_stage"},
	}
}

func ResourceTypeOrder(mode string) []string {
	stages := stageOrder(mode)
	order := make([]string, 0, len(stages))
	for _, stage := range stages {
		if utils.IsResourceCompatibleWithMode(stage.ResourceType, mode) {
			order = append(order, stage.ResourceType)
		}
	}
	return order
}

func (i *Importer) generateStagesTF() (string, error) {
	var tfConfig strings.Builder

	tflog.Info(i.ctx, "Generating stages for mode", map[string]interface{}{
		"mode": i.Mode,
	})

	stages := stageOrder(i.Mode)

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

func (i *Importer) importACLsIPv4() (interface{}, error) {
	return i.importACLs("4")
}

func (i *Importer) importACLsIPv6() (interface{}, error) {
	return i.importACLs("6")
}

func (i *Importer) importACLs(ipVersion string) (map[string]map[string]interface{}, error) {
	return i.fetch("IPv"+ipVersion+" ACLs", "/acls", map[string]string{"ip_version": ipVersion},
		utils.ResponseCollectionKeyForType("verity_acl_v"+ipVersion))
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

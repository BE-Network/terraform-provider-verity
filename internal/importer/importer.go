package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
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
		{
			name:        "gateways",
			importer:    i.importGateways,
			tfGenerator: i.generateGatewaysTF,
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

func (i *Importer) generateGatewayProfilesTF(data interface{}) (string, error) {
	profiles, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid gateway profiles data format")
	}

	var tfConfig strings.Builder
	for name, profile := range profiles {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_gateway_profile" "%s" {
	name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := profile["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"].(string); ok && group != "" {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = \"\"\n")
			}
		} else {
			tfConfig.WriteString("		group = \"\"\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range profile {
			if key == "name" || key == "index" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if key == "external_gateways" {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))
							if index, ok := itemMap["index"].(float64); ok {
								tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(index)))
							}
							for itemKey, itemValue := range itemMap {
								if itemKey == "index" {
									continue
								}
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
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
func (i *Importer) generateEthPortProfilesTF(data interface{}) (string, error) {
	profiles, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid eth port profiles data format")
	}

	var tfConfig strings.Builder
	for name, profile := range profiles {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_eth_port_profile" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := profile["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"]; ok && group != nil {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = null\n")
			}

			if portMonitoring, ok := objProps["port_monitoring"]; ok && portMonitoring != nil {
				tfConfig.WriteString(fmt.Sprintf("		port_monitoring = %s\n", formatValue(portMonitoring)))
			} else {
				tfConfig.WriteString("		port_monitoring = null\n")
			}
		} else {
			tfConfig.WriteString("		group = null\n")
			tfConfig.WriteString("		port_monitoring = null\n")
		}

		tfConfig.WriteString("	}\n")

		for key, value := range profile {
			if key == "name" || key == "index" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if key == "services" {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))
							if index, ok := itemMap["index"].(float64); ok {
								tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(index)))
							}
							for itemKey, itemValue := range itemMap {
								if itemKey == "index" {
									continue
								}
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
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

func (i *Importer) generateBundlesTF(data interface{}) (string, error) {
	bundles, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid bundles data format")
	}

	var tfConfig strings.Builder
	for name, bundle := range bundles {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_bundle" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := bundle["object_properties"].(map[string]interface{}); ok {
			if isForSwitch, ok := objProps["is_for_switch"].(bool); ok {
				tfConfig.WriteString(fmt.Sprintf("		is_for_switch = %t\n", isForSwitch))
			} else {
				tfConfig.WriteString("		is_for_switch = false\n")
			}
		} else {
			tfConfig.WriteString("		is_for_switch = false\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range bundle {
			if key == "name" || key == "index" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if key == "eth_port_paths" || key == "user_services" {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))
							for itemKey, itemValue := range itemMap {
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
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

func (i *Importer) generateTenantsTF(data interface{}) (string, error) {
	tenants, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid tenants data format")
	}

	var tfConfig strings.Builder
	for name, tenant := range tenants {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_tenant" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := tenant["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"].(string); ok {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = null\n")
			}
		} else {
			tfConfig.WriteString("		group = null\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range tenant {
			if key == "name" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if key == "route_tenants" {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))
							if index, ok := itemMap["index"].(float64); ok {
								tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(index)))
							}
							for itemKey, itemValue := range itemMap {
								if itemKey == "index" {
									continue
								}
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
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

func (i *Importer) generateLagsTF(data interface{}) (string, error) {
	lags, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid lags data format")
	}

	var tfConfig strings.Builder
	for name, lag := range lags {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_lag" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {}\n")

		for key, value := range lag {
			if key == "name" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", key))
			}
		}

		tfConfig.WriteString("}\n\n")
	}

	return tfConfig.String(), nil
}

func (i *Importer) generateServicesTF(data interface{}) (string, error) {
	services, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid services data format")
	}

	var tfConfig strings.Builder
	for name, service := range services {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_service" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := service["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"].(string); ok {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = null\n")
			}
		} else {
			tfConfig.WriteString("		group = null\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range service {
			if key == "name" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", key))
			}
		}

		tfConfig.WriteString("}\n\n")
	}

	return tfConfig.String(), nil
}

func (i *Importer) generateEthPortSettingsTF(data interface{}) (string, error) {
	settings, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid eth port settings data format")
	}

	var tfConfig strings.Builder
	for name, setting := range settings {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_eth_port_settings" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := setting["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"].(string); ok {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = null\n")
			}
		} else {
			tfConfig.WriteString("		group = null\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range setting {
			if key == "name" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case nil:
				tfConfig.WriteString(fmt.Sprintf("	%s = null\n", key))
			}
		}

		tfConfig.WriteString("}\n\n")
	}

	return tfConfig.String(), nil
}

func (i *Importer) generateGatewaysTF(data interface{}) (string, error) {
	gateways, ok := data.(map[string]map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid gateways data format")
	}

	var tfConfig strings.Builder
	for name, gateway := range gateways {
		sanitizedName := utils.SanitizeResourceName(name)
		tfConfig.WriteString(fmt.Sprintf(`
resource "verity_gateway" "%s" {
    name = "%s"
`, sanitizedName, name))

		tfConfig.WriteString("	object_properties {\n")

		if objProps, ok := gateway["object_properties"].(map[string]interface{}); ok {
			if group, ok := objProps["group"].(string); ok {
				tfConfig.WriteString(fmt.Sprintf("		group = %s\n", formatValue(group)))
			} else {
				tfConfig.WriteString("		group = null\n")
			}
		} else {
			tfConfig.WriteString("		group = null\n")
		}
		tfConfig.WriteString("	}\n")

		for key, value := range gateway {
			if key == "name" || key == "object_properties" {
				continue
			}

			switch v := value.(type) {
			case bool:
				tfConfig.WriteString(fmt.Sprintf("	%s = %t\n", key, v))
			case float64:
				tfConfig.WriteString(fmt.Sprintf("	%s = %d\n", key, int(v)))
			case string:
				tfConfig.WriteString(fmt.Sprintf("	%s = %s\n", key, formatValue(v)))
			case []interface{}:
				if key == "static_routes" {
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							tfConfig.WriteString(fmt.Sprintf("	%s {\n", key))
							if index, ok := itemMap["index"].(float64); ok {
								tfConfig.WriteString(fmt.Sprintf("		index = %d\n", int(index)))
							}
							for itemKey, itemValue := range itemMap {
								if itemKey == "index" {
									continue
								}
								tfConfig.WriteString(fmt.Sprintf("		%s = %s\n", itemKey, formatValue(itemValue)))
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

package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type OperationMode string

const (
	ModeDatacenter OperationMode = "datacenter"
	ModeCampus     OperationMode = "campus"
)

type ResourceMode string

const (
	ResourceModeDatacenter ResourceMode = "datacenter"
	ResourceModeCampus     ResourceMode = "campus"
	ResourceModeBoth       ResourceMode = "both"
)

type ApiVersionSupport struct {
	Major int
	Minor int
}

var ResourceCompatibility = map[string]ResourceMode{
	// Datacenter-only resources
	"verity_tenant":                  ResourceModeDatacenter,
	"verity_gateway":                 ResourceModeDatacenter,
	"verity_gateway_profile":         ResourceModeDatacenter,
	"verity_packet_broker":           ResourceModeDatacenter,
	"verity_as_path_access_list":     ResourceModeDatacenter,
	"verity_community_list":          ResourceModeDatacenter,
	"verity_extended_community_list": ResourceModeDatacenter,
	"verity_ipv4_prefix_list":        ResourceModeDatacenter,
	"verity_ipv6_prefix_list":        ResourceModeDatacenter,
	"verity_route_map_clause":        ResourceModeDatacenter,
	"verity_route_map":               ResourceModeDatacenter,
	"verity_sfp_breakout":            ResourceModeDatacenter,
	"verity_pod":                     ResourceModeDatacenter,
	"verity_spine_plane":             ResourceModeDatacenter,

	// Campus-only resources
	"verity_authenticated_eth_port": ResourceModeCampus,
	"verity_device_voice_settings":  ResourceModeCampus,
	"verity_service_port_profile":   ResourceModeCampus,
	"verity_voice_port_profile":     ResourceModeCampus,

	// Both mode resources
	"verity_service":                  ResourceModeBoth,
	"verity_eth_port_profile":         ResourceModeBoth,
	"verity_eth_port_settings":        ResourceModeBoth,
	"verity_device_settings":          ResourceModeBoth,
	"verity_lag":                      ResourceModeBoth,
	"verity_bundle":                   ResourceModeBoth,
	"verity_acl_v4":                   ResourceModeBoth,
	"verity_acl_v6":                   ResourceModeBoth,
	"verity_ipv4_list":                ResourceModeBoth,
	"verity_ipv6_list":                ResourceModeBoth,
	"verity_port_acl":                 ResourceModeBoth,
	"verity_badge":                    ResourceModeBoth,
	"verity_switchpoint":              ResourceModeBoth,
	"verity_device_controller":        ResourceModeBoth,
	"verity_site":                     ResourceModeBoth,
	"verity_packet_queue":             ResourceModeBoth,
	"verity_sflow_collector":          ResourceModeBoth,
	"verity_diagnostics_profile":      ResourceModeBoth,
	"verity_diagnostics_port_profile": ResourceModeBoth,
	"verity_pb_routing":               ResourceModeBoth,
	"verity_pb_routing_acl":           ResourceModeBoth,
	"verity_grouping_rule":            ResourceModeBoth,
	"verity_threshold_group":          ResourceModeBoth,
	"verity_threshold":                ResourceModeBoth,
}

var ResourceVersionCompatibility = map[string]ApiVersionSupport{
	"verity_tenant":                   {Major: 6, Minor: 4},
	"verity_gateway":                  {Major: 6, Minor: 4},
	"verity_service":                  {Major: 6, Minor: 4},
	"verity_eth_port_profile":         {Major: 6, Minor: 4},
	"verity_eth_port_settings":        {Major: 6, Minor: 4},
	"verity_bundle":                   {Major: 6, Minor: 4},
	"verity_lag":                      {Major: 6, Minor: 4},
	"verity_gateway_profile":          {Major: 6, Minor: 4},
	"verity_badge":                    {Major: 6, Minor: 5},
	"verity_switchpoint":              {Major: 6, Minor: 5},
	"verity_acl_v4":                   {Major: 6, Minor: 5},
	"verity_acl_v6":                   {Major: 6, Minor: 5},
	"verity_packet_broker":            {Major: 6, Minor: 5},
	"verity_authenticated_eth_port":   {Major: 6, Minor: 5},
	"verity_device_voice_settings":    {Major: 6, Minor: 5},
	"verity_packet_queue":             {Major: 6, Minor: 5},
	"verity_service_port_profile":     {Major: 6, Minor: 5},
	"verity_voice_port_profile":       {Major: 6, Minor: 5},
	"verity_device_controller":        {Major: 6, Minor: 5},
	"verity_as_path_access_list":      {Major: 6, Minor: 5},
	"verity_community_list":           {Major: 6, Minor: 5},
	"verity_device_settings":          {Major: 6, Minor: 5},
	"verity_extended_community_list":  {Major: 6, Minor: 5},
	"verity_ipv4_list":                {Major: 6, Minor: 5},
	"verity_ipv4_prefix_list":         {Major: 6, Minor: 5},
	"verity_ipv6_list":                {Major: 6, Minor: 5},
	"verity_ipv6_prefix_list":         {Major: 6, Minor: 5},
	"verity_route_map_clause":         {Major: 6, Minor: 5},
	"verity_route_map":                {Major: 6, Minor: 5},
	"verity_sfp_breakout":             {Major: 6, Minor: 5},
	"verity_site":                     {Major: 6, Minor: 5},
	"verity_pod":                      {Major: 6, Minor: 5},
	"verity_port_acl":                 {Major: 6, Minor: 5},
	"verity_sflow_collector":          {Major: 6, Minor: 5},
	"verity_diagnostics_profile":      {Major: 6, Minor: 5},
	"verity_diagnostics_port_profile": {Major: 6, Minor: 5},
	"verity_pb_routing":               {Major: 6, Minor: 5},
	"verity_pb_routing_acl":           {Major: 6, Minor: 5},
	"verity_spine_plane":              {Major: 6, Minor: 5},
	"verity_grouping_rule":            {Major: 6, Minor: 5},
	"verity_threshold_group":          {Major: 6, Minor: 5},
	"verity_threshold":                {Major: 6, Minor: 5},
}

func FilterResourcesByMode(
	ctx context.Context,
	resources []func() resource.Resource,
	currentMode string,
	apiVersion string,
) []func() resource.Resource {
	// If currentMode is empty, default to "datacenter"
	if currentMode == "" {
		currentMode = string(ModeDatacenter)
	}

	tflog.Info(ctx, "Filtering resources", map[string]interface{}{
		"mode":        currentMode,
		"api_version": apiVersion,
	})

	compatibleResources := make([]func() resource.Resource, 0, len(resources))
	knownTypes := make(map[string]bool)
	unknownTypeCount := 0
	resourceTypeToConstructor := make(map[string]func() resource.Resource, len(resources))

	for _, constructorFn := range resources {
		instance := constructorFn()

		for resourceType := range ResourceCompatibility {
			baseName := strings.TrimPrefix(resourceType, "verity_")

			constructorName := fmt.Sprintf("%T", instance)
			if strings.Contains(strings.ToLower(constructorName), strings.ToLower(baseName)) {
				resourceTypeToConstructor[resourceType] = constructorFn
				knownTypes[resourceType] = true
				break
			}
		}

		if _, found := knownTypes[fmt.Sprintf("%T", instance)]; !found {
			unknownTypeCount++
		}
	}

	tflog.Debug(ctx, "Resource mapping results", map[string]interface{}{
		"known_types":        len(knownTypes),
		"unknown_type_count": unknownTypeCount,
	})

	for resourceType, constructorFn := range resourceTypeToConstructor {
		if IsResourceCompatible(resourceType, currentMode, apiVersion) {
			tflog.Debug(ctx, "Resource is compatible and included", map[string]interface{}{
				"resource_type": resourceType,
				"mode":          currentMode,
				"api_version":   apiVersion,
			})
			compatibleResources = append(compatibleResources, constructorFn)
		} else {
			tflog.Debug(ctx, "Resource is NOT compatible and excluded", map[string]interface{}{
				"resource_type": resourceType,
				"mode":          currentMode,
				"api_version":   apiVersion,
			})
		}
	}

	parsedMajor, parsedMinor, err := ParseApiVersion(apiVersion)
	if err == nil && currentMode == string(ModeDatacenter) &&
		(parsedMajor > 6 || (parsedMajor == 6 && parsedMinor >= 4)) {
		tflog.Info(ctx, "Using datacenter mode with API 6.4+, returning all resources", map[string]interface{}{
			"total_resources": len(resources),
		})
		return resources
	}

	tflog.Info(ctx, "Resource filtering complete", map[string]interface{}{
		"total_resources":      len(resources),
		"compatible_resources": len(compatibleResources),
	})

	return compatibleResources
}

func ParseApiVersion(version string) (int, int, error) {
	parts := strings.Split(version, ".")

	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid version format: %s", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	return major, minor, nil
}

func IsResourceCompatibleWithMode(resourceType string, mode string) bool {
	compatMode, exists := ResourceCompatibility[resourceType]
	if !exists {
		return false
	}

	return compatMode == ResourceModeBoth ||
		(compatMode == ResourceModeDatacenter && mode == string(ModeDatacenter)) ||
		(compatMode == ResourceModeCampus && mode == string(ModeCampus))
}

func IsResourceCompatibleWithVersion(resourceType string, apiVersion string) bool {
	currentMajor, currentMinor, err := ParseApiVersion(apiVersion)
	if err != nil {
		return false
	}

	requiredVersion, exists := ResourceVersionCompatibility[resourceType]
	if !exists {
		return false
	}

	if currentMajor > requiredVersion.Major {
		return true
	}

	if currentMajor == requiredVersion.Major && currentMinor >= requiredVersion.Minor {
		return true
	}

	return false
}

func IsResourceCompatible(resourceType string, mode string, apiVersion string) bool {
	modeCompatible := IsResourceCompatibleWithMode(resourceType, mode)
	versionCompatible := IsResourceCompatibleWithVersion(resourceType, apiVersion)

	tflog.Debug(context.Background(), "Resource compatibility check", map[string]interface{}{
		"resource_type":      resourceType,
		"mode":               mode,
		"api_version":        apiVersion,
		"mode_compatible":    modeCompatible,
		"version_compatible": versionCompatible,
	})

	return modeCompatible && versionCompatible
}

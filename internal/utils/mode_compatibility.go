package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SupportedAPIVersion defines the API version this provider is built for
const (
	SupportedAPIMajor = 6
	SupportedAPIMinor = 5
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

var ResourceCompatibility = map[string]ResourceMode{
	// Datacenter-only resources
	"verity_tenant":                  ResourceModeDatacenter,
	"verity_gateway":                 ResourceModeDatacenter,
	"verity_gateway_profile":         ResourceModeDatacenter,
	"verity_pod":                     ResourceModeDatacenter,
	"verity_spine_plane":             ResourceModeDatacenter,
	"verity_route_map_clause":        ResourceModeDatacenter,
	"verity_route_map":               ResourceModeDatacenter,
	"verity_as_path_access_list":     ResourceModeDatacenter,
	"verity_community_list":          ResourceModeDatacenter,
	"verity_extended_community_list": ResourceModeDatacenter,
	"verity_ipv4_prefix_list":        ResourceModeDatacenter,
	"verity_ipv6_prefix_list":        ResourceModeDatacenter,
	"verity_packet_broker":           ResourceModeDatacenter,
	"verity_sfp_breakout":            ResourceModeDatacenter,

	// Campus-only resources
	"verity_authenticated_eth_port": ResourceModeCampus,
	"verity_device_voice_settings":  ResourceModeCampus,
	"verity_service_port_profile":   ResourceModeCampus,
	"verity_voice_port_profile":     ResourceModeCampus,

	// Resources available on both datacenter and campus systems
	"verity_acl_v4":                   ResourceModeBoth,
	"verity_acl_v6":                   ResourceModeBoth,
	"verity_badge":                    ResourceModeBoth,
	"verity_bundle":                   ResourceModeBoth,
	"verity_device_controller":        ResourceModeBoth,
	"verity_device_settings":          ResourceModeBoth,
	"verity_diagnostics_port_profile": ResourceModeBoth,
	"verity_diagnostics_profile":      ResourceModeBoth,
	"verity_eth_port_profile":         ResourceModeBoth,
	"verity_eth_port_settings":        ResourceModeBoth,
	"verity_ipv4_list":                ResourceModeBoth,
	"verity_ipv6_list":                ResourceModeBoth,
	"verity_lag":                      ResourceModeBoth,
	"verity_packet_queue":             ResourceModeBoth,
	"verity_pb_routing":               ResourceModeBoth,
	"verity_pb_routing_acl":           ResourceModeBoth,
	"verity_port_acl":                 ResourceModeBoth,
	"verity_service":                  ResourceModeBoth,
	"verity_sflow_collector":          ResourceModeBoth,
	"verity_site":                     ResourceModeBoth,
	"verity_switchpoint":              ResourceModeBoth,
	"verity_threshold_group":          ResourceModeBoth,
	"verity_threshold":                ResourceModeBoth,
	"verity_grouping_rule":            ResourceModeBoth,
}

// ValidateAPIVersion checks if the API version matches the supported version.
func ValidateAPIVersion(apiVersion string) error {
	major, minor, err := ParseApiVersion(apiVersion)
	if err != nil {
		return fmt.Errorf("failed to parse API version '%s': %w. This Terraform provider requires API version %d.%d",
			apiVersion, err, SupportedAPIMajor, SupportedAPIMinor)
	}

	if major != SupportedAPIMajor || minor != SupportedAPIMinor {
		return fmt.Errorf("API version mismatch: server is running API version %d.%d, but this Terraform provider is built for API version %d.%d. Please use a Terraform provider version that matches your API version",
			major, minor, SupportedAPIMajor, SupportedAPIMinor)
	}

	return nil
}

func GetSupportedAPIVersionString() string {
	return fmt.Sprintf("%d.%d", SupportedAPIMajor, SupportedAPIMinor)
}

// FilterResourcesByMode filters resources based on the operation mode (datacenter/campus).
func FilterResourcesByMode(
	ctx context.Context,
	resources []func() resource.Resource,
	currentMode string,
	apiVersion string,
) []func() resource.Resource {
	tflog.Info(ctx, "Filtering resources by mode", map[string]interface{}{
		"mode":        currentMode,
		"api_version": apiVersion,
	})

	compatibleResources := make([]func() resource.Resource, 0, len(resources))
	resourceTypeToConstructor := make(map[string]func() resource.Resource, len(resources))

	for _, constructorFn := range resources {
		instance := constructorFn()

		for resourceType := range ResourceCompatibility {
			baseName := strings.TrimPrefix(resourceType, "verity_")
			constructorName := fmt.Sprintf("%T", instance)
			if strings.Contains(strings.ToLower(constructorName), strings.ToLower(baseName)) {
				resourceTypeToConstructor[resourceType] = constructorFn
				break
			}
		}
	}

	for resourceType, constructorFn := range resourceTypeToConstructor {
		if IsResourceCompatibleWithMode(resourceType, currentMode) {
			tflog.Debug(ctx, "Resource is compatible with mode", map[string]interface{}{
				"resource_type": resourceType,
				"mode":          currentMode,
			})
			compatibleResources = append(compatibleResources, constructorFn)
		} else {
			tflog.Debug(ctx, "Resource is NOT compatible with mode", map[string]interface{}{
				"resource_type": resourceType,
				"mode":          currentMode,
			})
		}
	}

	tflog.Info(ctx, "Resource filtering complete", map[string]interface{}{
		"total_resources":      len(resources),
		"compatible_resources": len(compatibleResources),
	})

	return compatibleResources
}

// ParseApiVersion extracts major and minor version numbers from a version string.
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

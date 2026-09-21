package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	SupportedAPIMajor = 6
	SupportedAPIMinor = 6
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

var pendingResourceCompatibility = map[string]ResourceMode{

	"verity_operation_stage": ResourceModeBoth,
}

var ResourceCompatibility = mergeResourceCompatibility()

func mergeResourceCompatibility() map[string]ResourceMode {
	merged := make(map[string]ResourceMode, len(generatedResourceCompatibility)+len(pendingResourceCompatibility))
	for name, mode := range generatedResourceCompatibility {
		merged[name] = mode
	}
	for name, mode := range pendingResourceCompatibility {
		merged[name] = mode
	}
	return merged
}

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
		metadataResponse := resource.MetadataResponse{}
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadataResponse)

		if _, isModeControlled := ResourceCompatibility[metadataResponse.TypeName]; isModeControlled {
			resourceTypeToConstructor[metadataResponse.TypeName] = constructorFn
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

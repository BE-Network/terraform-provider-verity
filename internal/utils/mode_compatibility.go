package utils

import (
	"fmt"
	"strconv"
	"strings"
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

var nonAPIResourceCompatibility = map[string]ResourceMode{
	"verity_operation_stage": ResourceModeBoth,
}

func ResourceModeFor(resourceType string) (ResourceMode, bool) {
	if mode, found := nonAPIResourceCompatibility[resourceType]; found {
		return mode, true
	}
	resource, found := registryResource(resourceType)
	if !found {
		return "", false
	}
	modes := specModeSet(resource.Modes)
	switch {
	case modes[string(ModeDatacenter)] && modes[string(ModeCampus)]:
		return ResourceModeBoth, true
	case modes[string(ModeDatacenter)]:
		return ResourceModeDatacenter, true
	case modes[string(ModeCampus)]:
		return ResourceModeCampus, true
	default:
		return "", false
	}
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
	compatMode, exists := ResourceModeFor(resourceType)
	if !exists {
		return false
	}

	return compatMode == ResourceModeBoth ||
		(compatMode == ResourceModeDatacenter && mode == string(ModeDatacenter)) ||
		(compatMode == ResourceModeCampus && mode == string(ModeCampus))
}

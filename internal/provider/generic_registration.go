package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/transport"
)

const LegacyResourcesEnvVar = "VERITY_LEGACY_RESOURCES"

func genericAdapter(terraformType string) (genericresource.TransportAdapter, bool) {
	adapter, found := transport.GeneratedAdapters[terraformType]
	return adapter, found
}

func legacySelection() map[string]bool {
	raw := strings.TrimSpace(os.Getenv(LegacyResourcesEnvVar))
	if raw == "" {
		return nil
	}
	selected := make(map[string]bool)
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		selected[entry] = true
	}
	return selected
}

func genericConstructor(terraformType string, legacy map[string]bool) (func() resource.Resource, error) {
	if legacy[terraformType] || legacy["all"] {
		return nil, nil
	}
	adapter, hasAdapter := genericAdapter(terraformType)
	if !hasAdapter {
		return nil, fmt.Errorf("%s has no generated transport adapter", terraformType)
	}
	resourceSpec, err := registry.Lookup(terraformType)
	if err != nil {
		return nil, err
	}
	factory, err := genericresource.New(resourceSpec, adapter, bindGenericRuntime)
	if err != nil {
		return nil, fmt.Errorf("%s cannot be served generically: %w", terraformType, err)
	}
	return factory, nil
}

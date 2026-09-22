package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/transport"
)

func genericConstructor(terraformType string) (func() resource.Resource, error) {
	adapter, hasAdapter := transport.GeneratedAdapters[terraformType]
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

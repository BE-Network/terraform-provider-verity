package provider

import (
	"fmt"
	"sort"

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

func registryResourceOrder() []string {
	resources, err := registry.Load()
	if err != nil {
		panic("provider: " + err.Error())
	}
	order := make([]string, 0, len(resources))
	for _, resource := range resources {
		order = append(order, resource.TerraformType)
	}
	sort.Strings(order)
	return order
}

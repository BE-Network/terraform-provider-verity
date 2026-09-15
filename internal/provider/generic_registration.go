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

// GenericResourcesEnvVar names the resources the generic engine serves.
//
// Phase 2 migrates one resource, and a migration is only safe if it can be
// turned off in the field: the switch is what lets the generic engine be run
// against a real 6.6 deployment beside the implementation it replaces, without
// committing every user to it. Set it to a comma-separated list of Terraform
// types, or to "all" for every resource the engine can serve.
//
// Unset, the provider registers the handwritten resources exactly as before, so
// the default path through this file is the one that changes nothing.
const GenericResourcesEnvVar = "VERITY_GENERIC_RESOURCES"

// genericAdapters holds the transport adapter for each resource the engine can
// serve. The bulk manager type-asserts the value it is handed, so an adapter is
// what lets a canonical object reach it; a resource with no adapter here cannot
// be served generically however complete its spec is.
var genericAdapters = map[string]genericresource.TransportAdapter{
	"verity_ipv4_list": transport.IPv4ListAdapter{},
}

// genericSelection reads the switch. An empty set means every resource stays
// handwritten.
func genericSelection() map[string]bool {
	raw := strings.TrimSpace(os.Getenv(GenericResourcesEnvVar))
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

// genericConstructor returns the generic factory for a resource when the switch
// selects it and the engine can serve it.
//
// A selection the engine cannot honor is an error rather than a silent fallback
// to the handwritten resource: someone who asked for the generic engine and
// quietly got the old one would be testing the wrong thing and would not know.
func genericConstructor(terraformType string, selected map[string]bool) (func() resource.Resource, error) {
	if !selected[terraformType] && !selected["all"] {
		return nil, nil
	}
	adapter, hasAdapter := genericAdapters[terraformType]
	if !hasAdapter {
		if selected["all"] {
			// "all" is a convenience, not a claim that every resource is ready.
			return nil, nil
		}
		return nil, fmt.Errorf("%s=%s selects %s, which has no transport adapter yet",
			GenericResourcesEnvVar, os.Getenv(GenericResourcesEnvVar), terraformType)
	}
	resourceSpec, err := registry.Lookup(terraformType)
	if err != nil {
		return nil, err
	}
	factory, err := genericresource.New(resourceSpec, adapter, bindGenericRuntime)
	if err != nil {
		if selected["all"] {
			return nil, nil
		}
		return nil, fmt.Errorf("%s cannot be served generically: %w", terraformType, err)
	}
	return factory, nil
}

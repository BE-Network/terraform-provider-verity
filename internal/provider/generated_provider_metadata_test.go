package provider

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-verity/internal/spec"
)

// TestGeneratedResourceKeysCoverEveryResource checks that the generated key table
// serves every resource the provider registers. The resource implementations now
// take their endpoint name and cache key from it rather than declaring their own
// constants, so a missing entry would leave a resource addressing an empty
// endpoint rather than failing to compile.
func TestGeneratedResourceKeysCoverEveryResource(t *testing.T) {
	registry := generatedRegistry(t)
	if len(generatedResourceKeys) != len(registry) {
		t.Fatalf("key table has %d entries, the registry has %d", len(generatedResourceKeys), len(registry))
	}
	for _, resourceSpec := range registry {
		keys, exists := generatedResourceKeys[resourceSpec.TerraformType]
		if !exists {
			t.Errorf("no generated keys for %s", resourceSpec.TerraformType)
			continue
		}
		if keys.Endpoint == "" || keys.CacheKey == "" || keys.ResponseCollectionKey == "" {
			t.Errorf("%s has an incomplete key entry: %+v", resourceSpec.TerraformType, keys)
		}
		if want := strings.TrimPrefix(resourceSpec.API.EndpointPath, "/"); keys.Endpoint != want {
			t.Errorf("%s endpoint = %q, registry says %q", resourceSpec.TerraformType, keys.Endpoint, want)
		}
		if keys.CacheKey != resourceSpec.API.CacheKey {
			t.Errorf("%s cache key = %q, registry says %q", resourceSpec.TerraformType, keys.CacheKey, resourceSpec.API.CacheKey)
		}
	}
}

// TestRegistrationFollowsTheRegistry checks that registration is driven by the
// registry rather than by a handwritten list. Every registry resource needs a
// constructor, every constructor must belong to a registry resource, and each one
// must report the Terraform type it is registered under. A drift in any direction
// would either drop a resource from the provider or register one the registry does
// not describe.
func TestRegistrationFollowsTheRegistry(t *testing.T) {
	registry := generatedRegistry(t)
	if len(generatedResourceOrder) != len(registry) {
		t.Fatalf("registration order lists %d resources, the registry has %d", len(generatedResourceOrder), len(registry))
	}
	if len(resourceConstructors) != len(registry) {
		t.Fatalf("constructor map has %d entries, the registry has %d", len(resourceConstructors), len(registry))
	}

	inRegistry := make(map[string]bool, len(registry))
	for _, resourceSpec := range registry {
		inRegistry[resourceSpec.TerraformType] = true
	}
	ctx := context.Background()
	for terraformType, constructor := range resourceConstructors {
		if !inRegistry[terraformType] {
			t.Errorf("constructor registered for %q, which the registry does not describe", terraformType)
			continue
		}
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		if metadata.TypeName != terraformType {
			t.Errorf("constructor keyed %q reports type %q", terraformType, metadata.TypeName)
		}
	}

	// The order the provider registers in must be canonical, so the list is
	// reproducible from the registry rather than from edit history.
	if !sort.StringsAreSorted(generatedResourceOrder) {
		t.Error("registration order is not in canonical Terraform-name order")
	}

	// Every registered resource must appear exactly once, including the non-API
	// ones the plan keeps bespoke.
	seen := map[string]int{}
	for _, constructor := range getAllResources() {
		var metadata resource.MetadataResponse
		constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		seen[metadata.TypeName]++
	}
	for name, count := range seen {
		if count != 1 {
			t.Errorf("%s is registered %d times", name, count)
		}
	}
	if len(seen) != len(registry)+len(nonAPIResources) {
		t.Errorf("provider registers %d resources, want %d registry resources plus %d non-API", len(seen), len(registry), len(nonAPIResources))
	}
}

// TestModeRestrictedFieldsReachThePlanNullifier checks that every field the
// registry marks as narrower than its resource is handed to the ModifyPlan
// nullifier.
//
// A field that does not apply to the running mode is set to null in the plan so
// Terraform does not show "known after apply" for something the API will never
// return. The nullifier decides per field by asking FieldAppliesToMode, but the
// list of fields it is given is handwritten in each resource, so a mode-restricted
// field omitted from that list silently keeps showing as unknown. The registry
// knows which fields those are, so the omission is checkable.
func TestModeRestrictedFieldsReachThePlanNullifier(t *testing.T) {
	sources, err := filepath.Glob(filepath.Join("resource_verity_*.go"))
	if err != nil {
		t.Fatal(err)
	}
	listedByType := make(map[string]map[string]bool, len(sources))
	modifyPlan := regexp.MustCompile(`(?s)func \(r \*\w+\) ModifyPlan\(.*?\n\}`)
	typeName := regexp.MustCompile(`resp\.TypeName = req\.ProviderTypeName \+ "([a-z0-9_]+)"`)
	quoted := regexp.MustCompile(`"([a-z0-9_]+)"`)
	for _, source := range sources {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		names := typeName.FindAllStringSubmatch(string(raw), -1)
		body := modifyPlan.Find(raw)
		if len(names) == 0 || body == nil {
			continue
		}
		listed := make(map[string]bool)
		for _, match := range quoted.FindAllStringSubmatch(string(body), -1) {
			listed[match[1]] = true
		}
		for _, name := range names {
			// ACL builds its type name from the ip_version it was constructed with.
			if name[1] == "_acl_v" {
				listedByType["verity_acl_v4"], listedByType["verity_acl_v6"] = listed, listed
				continue
			}
			listedByType["verity"+name[1]] = listed
		}
	}

	checked := 0
	for _, resourceSpec := range generatedRegistry(t) {
		listed, exists := listedByType[resourceSpec.TerraformType]
		if !exists {
			continue
		}
		resourceModes := make(map[spec.Mode]bool, len(resourceSpec.Modes))
		for _, mode := range resourceSpec.Modes {
			resourceModes[mode] = true
		}
		for _, field := range resourceSpec.Fields {
			if field.Unmanaged || field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
				continue
			}
			if len(field.Modes) == len(resourceModes) {
				continue // applies wherever the resource does, so nothing to nullify
			}
			checked++
			if !listed[field.TerraformName] {
				t.Errorf("%s.%s applies to %v but the resource supports %v, and ModifyPlan never nullifies it",
					resourceSpec.TerraformType, field.TerraformName, field.Modes, resourceSpec.Modes)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no mode-restricted fields were checked; the scan looks broken")
	}
	t.Logf("%d mode-restricted fields confirmed reachable by the plan nullifier", checked)
}

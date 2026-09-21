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

	if !sort.StringsAreSorted(generatedResourceOrder) {
		t.Error("registration order is not in canonical Terraform-name order")
	}

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
				continue
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

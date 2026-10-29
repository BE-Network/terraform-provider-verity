package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type registeredResource struct {
	TypeName  string
	StateType string
}

func TestResourcesAreModeIndependentAndStable(t *testing.T) {
	provider := &verityProvider{}
	beforeConfiguration := registeredResources(t, provider.Resources(context.Background()))

	for _, mode := range []string{"datacenter", "campus"} {
		ctx := context.WithValue(context.Background(), "providerData", &providerContext{mode: mode})
		afterConfiguration := registeredResources(t, provider.Resources(ctx))
		if !reflect.DeepEqual(beforeConfiguration, afterConfiguration) {
			t.Fatalf("resource registration changed for mode %q:\n before: %#v\n after:  %#v", mode, beforeConfiguration, afterConfiguration)
		}
	}
}

func registeredResources(t *testing.T, constructors []func() resource.Resource) []registeredResource {
	t.Helper()
	ctx := context.Background()
	result := make([]registeredResource, 0, len(constructors))
	for _, constructor := range constructors {
		instance := constructor()
		var metadataResponse resource.MetadataResponse
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadataResponse)
		var schemaResponse resource.SchemaResponse
		instance.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
		result = append(result, registeredResource{
			TypeName:  metadataResponse.TypeName,
			StateType: schemaResponse.Schema.Type().TerraformType(ctx).String(),
		})
	}
	return result
}

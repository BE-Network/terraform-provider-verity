package lifecycle

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"sort"
	"terraform-provider-verity/internal/genericresource"

	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/tests/unit/mock"
)

func TestGenericResourcesMatchLegacyGoldenFixtures(t *testing.T) {
	for _, terraformType := range genericallyServable(t) {
		t.Run(terraformType, func(t *testing.T) {
			assertGoldenParity(t, terraformType)
		})
	}
}

func genericallyServable(t *testing.T) []string {
	t.Helper()
	names := make([]string, 0, len(transport.GeneratedAdapters))
	for terraformType := range transport.GeneratedAdapters {
		names = append(names, terraformType)
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatal("no generated adapters, so this test proved nothing")
	}
	return names
}

func assertGoldenParity(t *testing.T, terraformType string) {
	entry := coverageEntry(t, terraformType)
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	rs := inspectSchema(entry.Factory)
	modeKey := entry.modeFieldsKey()
	createConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"enable": "true"}))
	updateConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"enable": "false"}))

	address := entry.TerraformType + ".test"

	if entry.SkipCreate {

		fwresource.UnitTest(t, fwresource.TestCase{
			ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
			Steps: []fwresource.TestStep{{
				PreConfig:     func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
				Config:        createConfig,
				ResourceName:  address,
				ImportState:   true,
				ImportStateId: importIDFor(entry),
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("%s: imported %d instances, want 1", entry.TerraformType, len(states))
					}
					return compareGoldenAttributes(t, entry.TerraformType, states[0].Attributes)
				},
			}},
		})
		return
	}

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
				Config:    createConfig,
				Check: func(s *terraform.State) error {
					puts := ms.GetRequestsByMethodAndPath("PUT", entry.APIPath)
					if len(puts) == 0 {
						return fmt.Errorf("the generic engine sent no PUT to %s", entry.APIPath)
					}
					if err := compareGolden(t, entry.TerraformType, "put", puts[len(puts)-1].Body); err != nil {
						return err
					}
					return compareGoldenState(t, entry.TerraformType, address, s)
				},
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
				Config:    updateConfig,
				Check: func(*terraform.State) error {
					patches := ms.GetRequestsByMethodAndPath("PATCH", entry.APIPath)
					if len(patches) == 0 {
						return fmt.Errorf("the generic engine sent no PATCH to %s", entry.APIPath)
					}
					return compareGolden(t, entry.TerraformType, "patch", patches[len(patches)-1].Body)
				},
			},
		},
	})
}

func TestGenericIPv4ListImportsAndDeletes(t *testing.T) {
	entry := coverageEntry(t, "verity_ipv4_list")
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	rs := inspectSchema(entry.Factory)
	config := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, entry.modeFieldsKey(),
			mergeOverrides(entry.Overrides, map[string]string{"enable": "true"}))

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
			},
			{

				PreConfig:         func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:            config,
				ResourceName:      entry.TerraformType + ".test",
				ImportState:       true,
				ImportStateId:     entry.ResourceName,
				ImportStateVerify: true,

				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})

	deletes := ms.GetRequestsByMethodAndPath("DELETE", entry.APIPath)
	if len(deletes) == 0 {
		t.Fatalf("the generic engine sent no DELETE to %s", entry.APIPath)
	}
}

func assertServedGenerically(t *testing.T, terraformType string) {
	t.Helper()

	var served resource.Resource
	for _, factory := range provider.New("test")().Resources(context.Background()) {
		candidate := factory()
		var metadata resource.MetadataResponse
		candidate.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		if metadata.TypeName == terraformType {
			served = candidate
			break
		}
	}
	if served == nil {
		t.Fatalf("%s is not registered at all", terraformType)
	}
	if _, generic := served.(*genericresource.Resource); !generic {
		t.Fatalf("%s is served by %T, not the generic engine", terraformType, served)
	}
}

func coverageEntry(t *testing.T, terraformType string) ResourceCoverageEntry {
	t.Helper()
	for _, candidate := range allResourceTests {
		if candidate.TerraformType == terraformType {
			return candidate
		}
	}
	t.Fatalf("%s has no coverage entry", terraformType)
	return ResourceCoverageEntry{}
}

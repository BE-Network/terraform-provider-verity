package lifecycle

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/tests/unit/mock"
)

// The Phase 2 pilot is judged against the same fixtures as the resource it
// replaces.
//
// Every other check compares the generic engine to a description of the legacy
// implementation. This compares it to the legacy implementation's recorded
// output: the golden PUT, PATCH, and post-apply state captured from the
// handwritten verity_ipv4_list. The generic engine drives the same
// configuration against the same mock responses, and the bytes have to match
// what is already committed.
//
// That is the whole parity claim, and it is why this test reuses
// tests/unit/lifecycle/testdata/golden/verity_ipv4_list rather than recording a
// second baseline: a fixture of its own would only prove the engine agrees with
// itself.
//
// It cannot run in parallel, because the switch that selects the engine is an
// environment variable and t.Setenv forbids it.
func TestGenericIPv4ListMatchesLegacyGoldenFixtures(t *testing.T) {
	t.Setenv(provider.GenericResourcesEnvVar, "verity_ipv4_list")

	entry := coverageEntry(t, "verity_ipv4_list")
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	// The same configuration the golden fixtures were captured from: a full
	// create, then one field flipped so the update carries exactly one change.
	rs := inspectSchema(entry.Factory)
	modeKey := entry.modeFieldsKey()
	createConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"enable": "true"}))
	updateConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"enable": "false"}))

	address := entry.TerraformType + ".test"

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

// The engine must also delete and import, which the golden fixtures do not
// cover: the harness destroys at the end of every case, so a delete that failed
// would surface as a case error rather than a fixture diff, and import has no
// fixture at all for a resource that is normally created.
func TestGenericIPv4ListImportsAndDeletes(t *testing.T) {
	t.Setenv(provider.GenericResourcesEnvVar, "verity_ipv4_list")

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
				// Import reads the object back by name through the generic decoder
				// and must produce the state the apply already holds.
				PreConfig:         func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:            config,
				ResourceName:      entry.TerraformType + ".test",
				ImportState:       true,
				ImportStateId:     entry.ResourceName,
				ImportStateVerify: true,
				// These resources carry no "id"; the name is the identifier, which
				// is what ImportStatePassthroughID writes.
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})

	// The harness destroys what it created, so reaching here with no error means
	// the generic delete ran. Assert the request rather than infer it.
	deletes := ms.GetRequestsByMethodAndPath("DELETE", entry.APIPath)
	if len(deletes) == 0 {
		t.Fatalf("the generic engine sent no DELETE to %s", entry.APIPath)
	}
}

// assertServedGenerically fails if the switch did not actually take effect.
//
// Without this every assertion above would still pass while the handwritten
// resource served the request, and the test would claim parity for an engine it
// never exercised. It is the control the rest of the file depends on.
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
		t.Fatalf("%s is served by %T, not the generic engine; the %s switch did not take effect",
			terraformType, served, provider.GenericResourcesEnvVar)
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

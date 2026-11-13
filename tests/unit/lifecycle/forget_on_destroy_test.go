package lifecycle

import (
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestDestroyingSfpBreakoutRemovesItFromStateOnly(t *testing.T) {
	entry := coverageEntry(t, "verity_sfp_breakout")
	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}
	config := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(t, inspectSchema(entry.Factory), entry.TerraformType, importIDFor(entry), entry.Mode, entry.modeFieldsKey(), entry.Overrides)

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{{
			PreConfig:          func() { mock.WriteTFConfig(t, ms.URL(), config) },
			Config:             config,
			ResourceName:       entry.TerraformType + ".test",
			ImportState:        true,
			ImportStateId:      importIDFor(entry),
			ImportStatePersist: true,
		}},
	})

	if deletes := ms.GetRequestsByMethodAndPath("DELETE", entry.APIPath); len(deletes) != 0 {
		t.Fatalf("got %d DELETE requests for %s, want none", len(deletes), entry.APIPath)
	}
}

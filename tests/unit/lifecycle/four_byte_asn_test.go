package lifecycle

import (
	"strconv"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestFourByteASNRoundTrips(t *testing.T) {
	const asn = "4201200024"
	entry := coverageEntry(t, "verity_gateway")
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}

	rs := inspectSchema(entry.Factory)
	overrides := mergeOverrides(entry.Overrides, map[string]string{"neighbor_as_number": asn})
	config := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(rs, entry.TerraformType, entry.ResourceName, entry.Mode, entry.modeFieldsKey(), overrides)

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{{
			PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
			Config:    config,
			Check:     fwresource.TestCheckResourceAttr("verity_gateway.test", "neighbor_as_number", asn),
		}},
	})

	puts := ms.GetRequestsByMethodAndPath("PUT", entry.APIPath)
	if len(puts) == 0 {
		t.Fatal("no PUT captured")
	}
	object, err := extractResourceFromBody(puts[len(puts)-1].Body, entry.WrapperKey, entry.ResourceName)
	if err != nil {
		t.Fatal(err)
	}
	sent, ok := object["neighbor_as_number"].(float64)
	if !ok || strconv.FormatInt(int64(sent), 10) != asn {
		t.Fatalf("PUT neighbor_as_number = %v, want %s", object["neighbor_as_number"], asn)
	}
}

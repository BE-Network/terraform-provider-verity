package lifecycle

import (
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestUnconfiguredReferenceSurvivesAnUnrelatedUpdate(t *testing.T) {
	entry := coverageEntry(t, "verity_service")
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}
	ms.SetPostPutEnrichment(entry.APIPath, entry.WrapperKey, "diffsvc", map[string]interface{}{
		"policy_based_routing":           "pbr-ui",
		"policy_based_routing_ref_type_": "pb_routing",
	})

	service := func(vlan string) string {
		return `
resource "verity_service" "test" {
  name = "diffsvc"
  enable = true
  vlan = ` + vlan + `
  vni_auto_assigned_ = true
  mtu = 1500
  tenant = ""
  tenant_ref_type_ = ""
}
`
	}
	provider := mock.ProviderConfig(ms.URL(), entry.Mode)
	create, update := provider+service("10"), provider+service("20")

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), create) },
				Config:    create,
				Check: fwresource.ComposeAggregateTestCheckFunc(
					fwresource.TestCheckResourceAttr("verity_service.test", "policy_based_routing", "pbr-ui"),
				),
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), update) },
				Config:    update,
				Check: fwresource.ComposeAggregateTestCheckFunc(
					fwresource.TestCheckResourceAttr("verity_service.test", "vlan", "20"),
					fwresource.TestCheckResourceAttr("verity_service.test", "policy_based_routing", "pbr-ui"),
					fwresource.TestCheckResourceAttr("verity_service.test", "policy_based_routing_ref_type_", "pb_routing"),
				),
			},
		},
	})

	patches := ms.GetRequestsByMethodAndPath("PATCH", entry.APIPath)
	if len(patches) != 1 {
		t.Fatalf("got %d PATCH requests, want exactly the vlan update", len(patches))
	}
	if got, want := canonical(t, patches[0].Body), `{"service":{"diffsvc":{"vlan":20}}}`; got != want {
		t.Fatalf("PATCH = %s\nwant %s\nan unconfigured reference pair must not be sent during an unrelated update", got, want)
	}
}

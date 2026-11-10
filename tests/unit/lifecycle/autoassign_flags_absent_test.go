package lifecycle

import (
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestSwitchpointRefreshKeepsFlagsMissingFromGET(t *testing.T) {
	entry := coverageEntry(t, "verity_switchpoint")
	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}

	config := mock.ProviderConfig(ms.URL(), entry.Mode) + withPairs(switchpointBase(t), func(p switchpointPair) string {
		return "  " + p.flag() + " = true\n"
	})
	var flags []string
	for _, pair := range switchpointPairs {
		flags = append(flags, pair.flag())
	}

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
			},
			{
				PreConfig: func() {
					mock.WriteTFConfig(t, ms.URL(), config)
					ms.SetPostGetOmissions(entry.APIPath, entry.WrapperKey, "diffsp", flags, 10)
				},
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

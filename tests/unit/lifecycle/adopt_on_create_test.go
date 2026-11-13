package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func sfpBreakoutConfig(name, vendor string) string {
	return fmt.Sprintf(`
resource "verity_sfp_breakout" "test" {
  name = %q

  breakout {
    index       = 1
    breakout    = "1x100G"
    enable      = false
    part_number = ""
    vendor      = %q
  }
}
`, name, vendor)
}

func newSfpBreakoutServer(t *testing.T) (*mock.MockServer, string) {
	t.Helper()
	ms := mock.NewMockServer("datacenter")
	t.Cleanup(ms.Close)
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir("datacenter")); err != nil {
		t.Fatal(err)
	}
	return ms, mock.ProviderConfig(ms.URL(), "datacenter")
}

func requireOnlyGetAndPatch(t *testing.T, ms *mock.MockServer) {
	t.Helper()
	for _, method := range []string{"PUT", "DELETE", "POST"} {
		if requests := ms.GetRequestsByMethodAndPath(method, "/api/sfpbreakouts"); len(requests) != 0 {
			t.Fatalf("got %d %s requests for /api/sfpbreakouts, want none", len(requests), method)
		}
	}
}

func importThenForgetSteps(ms *mock.MockServer, t *testing.T, provider, config string) []fwresource.TestStep {
	return []fwresource.TestStep{
		{
			PreConfig:          func() { mock.WriteTFConfig(t, ms.URL(), provider+config) },
			Config:             provider + config,
			ResourceName:       "verity_sfp_breakout.test",
			ImportState:        true,
			ImportStateId:      "SFP Breakouts",
			ImportStatePersist: true,
		},
		{
			PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider) },
			Config:    provider,
		},
	}
}

func TestSfpBreakoutIsAdoptedAfterDestroyWithoutChanges(t *testing.T) {
	ms, provider := newSfpBreakoutServer(t)
	config := sfpBreakoutConfig("SFP Breakouts", "aaa")
	steps := append(importThenForgetSteps(ms, t, provider, config),
		fwresource.TestStep{
			PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider+config) },
			Config:    provider + config,
			Check:     fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "name", "SFP Breakouts"),
		},
		fwresource.TestStep{Config: provider + config, PlanOnly: true},
	)

	fwresource.UnitTest(t, fwresource.TestCase{ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(), Steps: steps})

	if patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/sfpbreakouts"); len(patches) != 0 {
		t.Fatalf("got %d PATCH requests, want none: the configuration matches the system", len(patches))
	}
	requireOnlyGetAndPatch(t, ms)
}

func TestSfpBreakoutAdoptionAppliesChangedSettings(t *testing.T) {
	ms, provider := newSfpBreakoutServer(t)
	original, changed := sfpBreakoutConfig("SFP Breakouts", "aaa"), sfpBreakoutConfig("SFP Breakouts", "bbb")
	steps := append(importThenForgetSteps(ms, t, provider, original),
		fwresource.TestStep{
			PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider+changed) },
			Config:    provider + changed,
			Check: fwresource.ComposeTestCheckFunc(
				fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "name", "SFP Breakouts"),
				fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "breakout.0.vendor", "bbb"),
			),
		},
		fwresource.TestStep{Config: provider + changed, PlanOnly: true},
	)

	fwresource.UnitTest(t, fwresource.TestCase{ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(), Steps: steps})

	patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/sfpbreakouts")
	if len(patches) != 1 || !strings.Contains(string(patches[0].RawBody), `"bbb"`) {
		t.Fatalf("got %d PATCH requests, want one with the changed vendor", len(patches))
	}
	requireOnlyGetAndPatch(t, ms)
}

func TestSfpBreakoutThatDoesNotExistIsNotCreated(t *testing.T) {
	ms, provider := newSfpBreakoutServer(t)
	config := provider + sfpBreakoutConfig("Missing Breakouts", "aaa")

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{{
			PreConfig:   func() { mock.WriteTFConfig(t, ms.URL(), config) },
			Config:      config,
			ExpectError: regexp.MustCompile(`doesn't exist on the system`),
		}},
	})

	if patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/sfpbreakouts"); len(patches) != 0 {
		t.Fatalf("got %d PATCH requests, want none", len(patches))
	}
	requireOnlyGetAndPatch(t, ms)
}

func TestSfpBreakoutAdoptionRecoversFromAFailedPatch(t *testing.T) {
	ms, provider := newSfpBreakoutServer(t)
	original, changed := sfpBreakoutConfig("SFP Breakouts", "aaa"), sfpBreakoutConfig("SFP Breakouts", "bbb")
	steps := append(importThenForgetSteps(ms, t, provider, original),
		fwresource.TestStep{
			PreConfig: func() {
				mock.WriteTFConfig(t, ms.URL(), provider+changed)
				ms.FailRequests("PATCH", "/api/sfpbreakouts", 1, "patch rejected by the test")
			},
			Config:      provider + changed,
			ExpectError: regexp.MustCompile(`patch rejected by the test`),
		},
		fwresource.TestStep{
			Config: provider + changed,
			Check:  fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "breakout.0.vendor", "bbb"),
		},
		fwresource.TestStep{Config: provider + changed, PlanOnly: true},
	)

	fwresource.UnitTest(t, fwresource.TestCase{ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(), Steps: steps})

	if patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/sfpbreakouts"); len(patches) != 2 {
		t.Fatalf("got %d PATCH requests, want the failed one and the retry", len(patches))
	}
	requireOnlyGetAndPatch(t, ms)
}

const twoBreakoutsInAPIOrder = `{"sfp_breakouts": {"SFP Breakouts": {"enable": false, "breakout": [
  {"enable": false, "vendor": "aaa", "part_number": "", "breakout": "1x100G", "index": 1},
  {"enable": false, "vendor": "bbb", "part_number": "", "breakout": "1x100G", "index": 2}
]}}}`

func reversedBreakoutConfig(secondVendor string) string {
	return fmt.Sprintf(`
resource "verity_sfp_breakout" "test" {
  name = "SFP Breakouts"

  breakout {
    index       = 2
    breakout    = "1x100G"
    enable      = false
    part_number = ""
    vendor      = %q
  }

  breakout {
    index       = 1
    breakout    = "1x100G"
    enable      = false
    part_number = ""
    vendor      = "aaa"
  }
}
`, secondVendor)
}

func TestSfpBreakoutAdoptionKeepsTheConfiguredEntryOrder(t *testing.T) {
	for _, tc := range []struct {
		name    string
		vendor  string
		patches int
	}{
		{"unchanged", "bbb", 0},
		{"changed entry", "zzz", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ms, provider := newSfpBreakoutServer(t)
			ms.SetGetResponse("/api/sfpbreakouts", []byte(twoBreakoutsInAPIOrder))
			config := provider + reversedBreakoutConfig(tc.vendor)

			fwresource.UnitTest(t, fwresource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []fwresource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: fwresource.ComposeTestCheckFunc(
							fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "breakout.0.index", "2"),
							fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "breakout.0.vendor", tc.vendor),
							fwresource.TestCheckResourceAttr("verity_sfp_breakout.test", "breakout.1.index", "1"),
						),
					},
					{Config: config, PlanOnly: true},
				},
			})

			patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/sfpbreakouts")
			if len(patches) != tc.patches {
				t.Fatalf("got %d PATCH requests, want %d", len(patches), tc.patches)
			}
			if tc.patches == 1 && (!strings.Contains(string(patches[0].RawBody), `"zzz"`) || strings.Contains(string(patches[0].RawBody), `"aaa"`)) {
				t.Fatalf("PATCH = %s, want only the changed entry", patches[0].RawBody)
			}
			requireOnlyGetAndPatch(t, ms)
		})
	}
}

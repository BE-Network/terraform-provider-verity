package lifecycle

import (
	"fmt"
	"testing"
)

func TestGenericMatchesLegacyOnSingletonUpdates(t *testing.T) {
	const badgeFixed = `  color = "red"
  enable = true
  number = 42
`
	badge := func(body string) string {
		return fmt.Sprintf("resource \"verity_badge\" \"test\" {\n  name = \"diffbadge\"\n%s}\n", badgeFixed+body)
	}

	const lagFixed = `  color = ""
  crc_failure_threshold = 42
  enable = true
  eth_port_profile = ""
  eth_port_profile_ref_type_ = ""
  fallback = true
  fast_rate = true
  is_peer_link = true
  lacp = true
  peer_link_vlan = 42
  uplink = true
`
	lag := func(body string) string {
		return fmt.Sprintf("resource \"verity_lag\" \"test\" {\n  name = \"difflag\"\n%s}\n", lagFixed+body)
	}

	scenarios := []struct {
		name           string
		terraformType  string
		create, update string

		outcome lifecycleOutcome
	}{
		{
			name:          "member changes",
			terraformType: "verity_badge",
			create:        badge("  object_properties {\n    notes = \"first\"\n  }\n"),
			update:        badge("  object_properties {\n    notes = \"second\"\n  }\n"),
		},
		{
			name:          "member is cleared to an empty string",
			terraformType: "verity_badge",
			create:        badge("  object_properties {\n    notes = \"first\"\n  }\n"),
			update:        badge("  object_properties {\n    notes = \"\"\n  }\n"),
		},
		{
			name:          "member is removed from the block",
			terraformType: "verity_badge",
			create:        badge("  object_properties {\n    notes = \"first\"\n  }\n"),
			update:        badge("  object_properties {\n  }\n"),
		},
		{

			name:          "block is removed",
			terraformType: "verity_badge",
			create:        badge("  object_properties {\n    notes = \"first\"\n  }\n"),
			update:        badge(""),
			outcome:       lifecycleOutcome{driftAfterApply: true},
		},
		{

			name:          "block is added",
			terraformType: "verity_badge",
			create:        badge(""),
			update:        badge("  object_properties {\n    notes = \"added\"\n  }\n"),
			outcome:       lifecycleOutcome{driftAfterApply: true},
		},
		{
			name:          "block is unchanged while another field changes",
			terraformType: "verity_badge",
			create:        badge("  object_properties {\n    notes = \"kept\"\n  }\n"),
			update: fmt.Sprintf("resource \"verity_badge\" \"test\" {\n  name = \"diffbadge\"\n%s}\n",
				"  color = \"blue\"\n  enable = true\n  number = 42\n  object_properties {\n    notes = \"kept\"\n  }\n"),
		},
		{

			name:          "reference pair inside the block changes its value",
			terraformType: "verity_lag",
			create:        lag("  object_properties {\n    fabric = \"fabric-a\"\n    fabric_ref_type_ = \"fabric\"\n  }\n"),
			update:        lag("  object_properties {\n    fabric = \"fabric-b\"\n    fabric_ref_type_ = \"fabric\"\n  }\n"),
		},
		{
			name:          "reference pair inside the block is cleared",
			terraformType: "verity_lag",
			create:        lag("  object_properties {\n    fabric = \"fabric-a\"\n    fabric_ref_type_ = \"fabric\"\n  }\n"),
			update:        lag("  object_properties {\n    fabric = \"\"\n    fabric_ref_type_ = \"\"\n  }\n"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, scenario.terraformType, false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, scenario.terraformType, true, scenario.create, scenario.update, scenario.outcome)

			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func TestSingletonEmptyBlockOnCreateOmitsUnknownMembers(t *testing.T) {
	create := fmt.Sprintf("resource \"verity_badge\" \"test\" {\n  name = \"diffbadge\"\n%s}\n",
		"  color = \"red\"\n  enable = true\n  number = 42\n  object_properties {\n  }\n")

	legacy := captureLifecycle(t, "verity_badge", false, create, create)
	generic := captureLifecycle(t, "verity_badge", true, create, create)

	const (
		wantLegacy  = `{"badge":{"diffbadge":{"color":"red","enable":true,"name":"diffbadge","number":42,"object_properties":{"notes":""}}}}`
		wantGeneric = `{"badge":{"diffbadge":{"color":"red","enable":true,"name":"diffbadge","number":42,"object_properties":{}}}}`
	)
	if got := canonical(t, legacy["PUT"]); got != wantLegacy {
		t.Errorf("legacy PUT = %s, want %s: if the handwritten helper now skips unknowns, the two agree and this test should assert equality", got, wantLegacy)
	}
	if got := canonical(t, generic["PUT"]); got != wantGeneric {
		t.Errorf("generic PUT = %s, want %s: an unknown member must follow unknown_plan", got, wantGeneric)
	}
}

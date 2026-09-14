package lifecycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/tests/unit/mock"
)

// Golden wire fixtures record the exact JSON each resource sends when a fully
// populated configuration is applied.
//
// Every other check in this repository compares the registry against the legacy
// *source*: schema shape, modes, aliases, lifecycle policies. That proves the
// registry describes the implementation, but it cannot prove a replacement
// implementation puts the same bytes on the wire. The migration plan asks for
// behavior parity, and parity needs a recorded baseline to compare against.
//
// These fixtures are that baseline. They are captured from the handwritten
// resources as they ship today, so when a resource moves to the generic engine
// the same test re-runs against the new implementation and any difference in the
// request body is a diff in a committed file rather than a silent behavior
// change.
//
// Regenerate deliberately, never to make a red test pass:
//
//	UPDATE_GOLDEN=1 go test ./tests/unit/lifecycle/ -run GoldenWireFixtures
//
// A diff here means the request changed. Read it before accepting it.
func TestGoldenWireFixtures(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			ms.SetTestLogger(t)
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			// Create with everything populated, then flip one field. The first
			// request records what a full create sends; the second records which
			// fields an update carries, which is where omission and explicit-null
			// behavior actually shows.
			createOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "true"})
			updateOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "false"})
			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) +
				generateCoverageHCL(rs, tc.TerraformType, tc.ResourceName, tc.Mode, modeKey, createOverrides)
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) +
				generateCoverageHCL(rs, tc.TerraformType, tc.ResourceName, tc.Mode, modeKey, updateOverrides)

			address := tc.TerraformType + ".test" // generateCoverageHCL labels every resource "test"

			var steps []fwresource.TestStep
			if tc.SkipCreate {
				// This resource refuses both create and delete: it represents
				// existing hardware that can only be read and updated. It is
				// therefore brought into state by import, and the state is not
				// persisted, because the harness destroys at the end of a case and
				// the destroy would fail. That pins its read path. Its update path
				// cannot be exercised here for the same reason, so it has no PUT or
				// PATCH fixture; see docs/exception_inventory.md.
				steps = append(steps, fwresource.TestStep{
					PreConfig:     func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
					Config:        createConfig,
					ResourceName:  address,
					ImportState:   true,
					ImportStateId: importIDFor(tc),
					// An import step reports through ImportStateCheck; Check is only
					// called for apply steps.
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						if len(states) != 1 {
							return fmt.Errorf("%s: imported %d instances, want 1", tc.TerraformType, len(states))
						}
						return compareGoldenAttributes(t, tc.TerraformType, states[0].Attributes)
					},
				})
			} else {
				steps = append(steps, fwresource.TestStep{
					PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
					Config:    createConfig,
					Check: func(s *terraform.State) error {
						puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
						if len(puts) == 0 {
							return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
						}
						if err := compareGolden(t, tc.TerraformType, "put", puts[len(puts)-1].Body); err != nil {
							return err
						}
						// The state after apply is the read path's output: the
						// provider fetched the object back and decoded it. Pinning
						// it covers the direction the request fixtures do not.
						return compareGoldenState(t, tc.TerraformType, address, s)
					},
				})
			}
			// Only resources that expose enable can be updated this way; the rest
			// record a create fixture alone rather than a contrived edit.
			if hasEnableAttribute(rs) {
				steps = append(steps, fwresource.TestStep{
					PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
					Config:    updateConfig,
					Check: func(s *terraform.State) error {
						patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
						if len(patches) == 0 {
							return fmt.Errorf("no PATCH request captured for %s", tc.APIPath)
						}
						return compareGolden(t, tc.TerraformType, "patch", patches[len(patches)-1].Body)
					},
				})
			}

			fwresource.UnitTest(t, fwresource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps:                    steps,
			})
		})
	}
}

// importIDFor names an object the mock actually serves. An update-only resource
// is imported rather than created, so the identifier has to match a real entry in
// the response corpus instead of the name the harness would have created.
func importIDFor(tc ResourceCoverageEntry) string {
	if tc.TerraformType == "verity_sfp_breakout" {
		return "SFP Breakouts"
	}
	return tc.ResourceName
}

func hasEnableAttribute(rs resourceSchemaInfo) bool {
	for _, fi := range rs.Attributes {
		if fi.Name == "enable" {
			return true
		}
	}
	return false
}

// compareGoldenState records the Terraform state an apply produced. Together with
// the request fixtures this pins both directions: what the provider sends, and
// what it makes of what the API sends back. A generic decoder that dropped a
// field, or read it as the wrong type, would show up here rather than passing
// silently because nothing asserted on the decoded value.
func compareGoldenState(t *testing.T, terraformType, address string, state *terraform.State) error {
	t.Helper()
	root := state.RootModule()
	if root == nil {
		return fmt.Errorf("%s: no root module in state", terraformType)
	}
	resourceState, exists := root.Resources[address]
	if !exists || resourceState.Primary == nil {
		return fmt.Errorf("%s: no state recorded at %s", terraformType, address)
	}
	return compareGoldenAttributes(t, terraformType, resourceState.Primary.Attributes)
}

// compareGoldenAttributes records a decoded instance's attributes, dropping the
// identifiers Terraform keeps for its own bookkeeping rather than reading from
// the API.
func compareGoldenAttributes(t *testing.T, terraformType string, raw map[string]string) error {
	t.Helper()
	attributes := make(map[string]string, len(raw))
	for name, value := range raw {
		if name == "id" || name == "%" {
			continue
		}
		attributes[name] = value
	}
	return compareGoldenValue(t, terraformType, "state", attributes)
}

// compareGolden writes the fixture when regenerating and otherwise compares
// against it. The body is re-encoded canonically so a fixture diff reflects a
// change in what was sent, not in key ordering.
func compareGolden(t *testing.T, terraformType, operation string, body map[string]interface{}) error {
	t.Helper()
	return compareGoldenValue(t, terraformType, operation, body)
}

func compareGoldenValue(t *testing.T, terraformType, operation string, value any) error {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s %s body: %w", terraformType, operation, err)
	}
	encoded = append(encoded, '\n')

	path := filepath.Join("testdata", "golden", terraformType, operation+".json")
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		return os.WriteFile(path, encoded, 0o644)
	}

	want, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("no golden fixture for %s %s: %w; regenerate with UPDATE_GOLDEN=1", terraformType, operation, err)
	}
	if string(want) != string(encoded) {
		return fmt.Errorf("%s %s request differs from its golden fixture.\n--- recorded\n%s\n--- sent now\n%s",
			terraformType, operation, want, encoded)
	}
	return nil
}

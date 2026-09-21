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

			createOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "true"})
			updateOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "false"})
			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) +
				generateCoverageHCL(rs, tc.TerraformType, tc.ResourceName, tc.Mode, modeKey, createOverrides)
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) +
				generateCoverageHCL(rs, tc.TerraformType, tc.ResourceName, tc.Mode, modeKey, updateOverrides)

			address := tc.TerraformType + ".test"

			var steps []fwresource.TestStep
			if tc.SkipCreate {

				steps = append(steps, fwresource.TestStep{
					PreConfig:     func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
					Config:        createConfig,
					ResourceName:  address,
					ImportState:   true,
					ImportStateId: importIDFor(tc),

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

						return compareGoldenState(t, tc.TerraformType, address, s)
					},
				})
			}

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

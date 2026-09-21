package lifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestImporterLeavesOutUnsupportedArguments(t *testing.T) {
	ms := mock.NewMockServer("datacenter")
	defer ms.Close()
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir("datacenter")); err != nil {
		t.Fatal(err)
	}
	tenants := map[string]map[string]map[string]interface{}{"tenant": {"newer_tenant": {
		"name":                    "newer_tenant",
		"enable":                  true,
		"maximum_ebgp_paths":      nil,
		"maximum_ebgp_paths_mode": "automated",
	}}}
	body, err := json.Marshal(tenants)
	if err != nil {
		t.Fatal(err)
	}
	ms.SetGetResponse("/api/tenants", body)
	ms.SetGetResponse("/api/tenants:tenants", body)

	outputDir := t.TempDir()
	config := mock.ProviderConfig(ms.URL(), "datacenter") + `
data "verity_state_importer" "test" {
  output_dir = "` + filepath.ToSlash(outputDir) + `"
}
`
	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps:                    []fwresource.TestStep{{Config: config}},
	})

	generated, err := os.ReadFile(filepath.Join(outputDir, "tenants.tf"))
	if err != nil {
		t.Fatal(err)
	}
	tf := string(generated)
	if !strings.Contains(tf, `name = "newer_tenant"`) || !strings.Contains(tf, "enable = true") {
		t.Fatalf("tenants.tf lost the tenant or a supported argument:\n%s", tf)
	}
	for _, unsupported := range []string{"maximum_ebgp_paths", "maximum_ebgp_paths_mode"} {
		if strings.Contains(tf, unsupported) {
			t.Errorf("tenants.tf still writes %s:\n%s", unsupported, tf)
		}
	}

	listed, err := os.ReadFile(filepath.Join(outputDir, "unsupported_arguments.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(listed), "verity_tenant: maximum_ebgp_paths, maximum_ebgp_paths_mode") {
		t.Errorf("unsupported_arguments.txt does not name the tenant's arguments:\n%s", listed)
	}
}

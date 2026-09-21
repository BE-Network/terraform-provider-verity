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

// A Verity system newer than the provider's API version returns arguments the
// schema does not have, and one of them in a generated file makes Terraform
// reject the whole import. The importer leaves them out; this runs the importer
// data source against a system whose tenants carry two such arguments and
// requires the generated tenants.tf to hold the tenant without them, and
// unsupported_arguments.txt, which the import scripts print at the end, to name
// them.
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

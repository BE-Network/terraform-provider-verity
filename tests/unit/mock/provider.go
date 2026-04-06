package mock

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/internal/utils"
)

var testDirs sync.Map

// ProtoV6ProviderFactories returns provider factories for use with terraform-plugin-testing.
// These point to the real Verity provider — the mock server intercepts HTTP calls.
func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"verity": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

// ProviderConfig returns the HCL provider configuration block pointing at the mock server.
func ProviderConfig(serverURL, mode string) string {
	return fmt.Sprintf(`
provider "verity" {
  uri      = %q
  username = "test"
  password = "test"
  mode     = %q
}
`, serverURL, mode)
}

// StageHCL returns an operation stage resource block.
// Resources require depends_on a stage to trigger bulk op execution.
func StageHCL(name string) string {
	return fmt.Sprintf(`
resource "verity_operation_stage" %q {
}
`, name)
}

// WriteTFConfig writes the given HCL config to a temp directory so that
// ParseResourceConfiguredAttributes can find .tf files during test execution.
// The serverURL is the mock server's URL, used as a key to isolate multi-step tests.
// On first call it creates a temp dir and registers it in the workdir registry.
// Subsequent calls for the same serverURL update the file in place.
// Use this in PreConfig callbacks for multi-step tests.
func WriteTFConfig(t *testing.T, serverURL, config string) {
	t.Helper()
	dir, loaded := testDirs.LoadOrStore(serverURL, "")
	if !loaded || dir.(string) == "" {
		d := t.TempDir()
		testDirs.Store(serverURL, d)
		utils.RegisterWorkDir(serverURL, d)
		t.Cleanup(func() {
			utils.UnregisterWorkDir(serverURL)
			testDirs.Delete(serverURL)
		})
		dir = d
	}
	if err := os.WriteFile(filepath.Join(dir.(string), "test.tf"), []byte(config), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
}

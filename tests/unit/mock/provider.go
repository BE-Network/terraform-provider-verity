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

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"verity": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

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

func StageHCL(name string) string {
	return fmt.Sprintf(`
resource "verity_operation_stage" %q {
}
`, name)
}

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
	utils.InvalidateConfigIndex(dir.(string))
}

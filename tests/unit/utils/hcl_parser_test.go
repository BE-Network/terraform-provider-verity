package utils_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	providerutils "terraform-provider-verity/internal/utils"
)

func writeTF(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

// The index must be rebuilt when a .tf file changes on disk.
func TestConfigIndexPicksUpFileChanges(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	writeTF(t, dir, "a.tf", `
resource "verity_service" "svc" {
  name = "svc"
  vlan = 100
}
`)

	got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "svc")
	if !got.IsConfigured("vlan") || got.IsConfigured("tenant") {
		t.Fatalf("first parse: got %+v", got.Attributes)
	}

	// Rewrite with a different attribute set; the size/mtime fingerprint alone
	// must invalidate the index.
	writeTF(t, dir, "a.tf", `
resource "verity_service" "svc" {
  name   = "svc"
  tenant = "t1"
}
`)

	got = providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "svc")
	if got.IsConfigured("vlan") || !got.IsConfigured("tenant") {
		t.Fatalf("after rewrite: got %+v", got.Attributes)
	}

	// A newly added file changes the file set, so the fingerprint alone must
	// invalidate the index without an explicit call.
	writeTF(t, dir, "b.tf", `
resource "verity_service" "other" {
  name = "other"
  mtu  = 9000
}
`)

	got = providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "other")
	if !got.IsConfigured("mtu") {
		t.Fatalf("after adding file: got %+v", got.Attributes)
	}
}

// A same-size in-place rewrite can land inside one timestamp tick, which the
// fingerprint cannot see. InvalidateConfigIndex is the escape hatch for callers
// that rewrite .tf files themselves, such as the unit-test mock provider.
func TestInvalidateConfigIndexForcesReparse(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	writeTF(t, dir, "a.tf", "resource \"verity_service\" \"svc\" {\n  name = \"svc\"\n  vlan = 100\n}\n")
	if got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "svc"); !got.IsConfigured("vlan") {
		t.Fatalf("first parse: %+v", got.Attributes)
	}

	// Same byte length, different attribute.
	writeTF(t, dir, "a.tf", "resource \"verity_service\" \"svc\" {\n  name = \"svc\"\n  mtu1 = 100\n}\n")
	providerutils.InvalidateConfigIndex(dir)

	got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "svc")
	if got.IsConfigured("vlan") || !got.IsConfigured("mtu1") {
		t.Fatalf("after invalidation: %+v", got.Attributes)
	}
}

// Matching priority: the "name" attribute wins; the block label (raw or
// sanitized) is only a fallback when "name" is absent.
func TestConfigIndexMatching(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	writeTF(t, dir, "a.tf", `
resource "verity_service" "label_only" {
  vlan = 10
}

resource "verity_service" "wrong_label" {
  name = "By Name"
  mtu  = 1500
}

resource "verity_gateway" "_Odd_Name_" {
  asn = 65000
}
`)

	// Matched by the name attribute, not by the label.
	if got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "By Name"); !got.IsConfigured("mtu") {
		t.Errorf("name-attribute match failed: %+v", got.Attributes)
	}
	// The label must not match when a differing name attribute is present.
	if got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "wrong_label"); len(got.Attributes) != 0 {
		t.Errorf("label should not match a block with a differing name attribute: %+v", got.Attributes)
	}
	// Fallback to the raw label when no name attribute exists.
	if got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "label_only"); !got.IsConfigured("vlan") {
		t.Errorf("label fallback failed: %+v", got.Attributes)
	}
	// Fallback to the sanitized label for API names that are not valid labels.
	if got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_gateway", "(Odd Name)"); !got.IsConfigured("asn") {
		t.Errorf("sanitized label fallback failed: %+v", got.Attributes)
	}
	// Unknown resources yield an empty, non-nil result.
	got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "missing")
	if got == nil || len(got.Attributes) != 0 || got.IsConfigured("vlan") {
		t.Errorf("unknown resource: %+v", got)
	}
}

// Results are merged across files, first match per file.
func TestConfigIndexMergesAcrossFiles(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	writeTF(t, dir, "a.tf", `
resource "verity_service" "svc" {
  name = "svc"
  vlan = 100
}
`)
	writeTF(t, dir, "b.tf", `
resource "verity_service" "svc" {
  name = "svc"
  mtu  = 9000
}
`)

	got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", "svc")
	if !got.IsConfigured("vlan") || !got.IsConfigured("mtu") {
		t.Fatalf("merge across files: %+v", got.Attributes)
	}
}

// Nested and indexed blocks survive the index round-trip.
func TestConfigIndexNestedAndIndexedBlocks(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	writeTF(t, dir, "a.tf", `
resource "verity_eth_port_profile" "epp" {
  name = "epp"

  services {
    index                 = 1
    row_num_enable        = true
    row_num_external_vlan = 200
  }

  services {
    index          = 2
    row_num_enable = false
  }

  object_properties {
    group = "g1"
  }
}
`)

	got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_eth_port_profile", "epp")

	if !got.IsBlockConfigured("services") || !got.IsBlockConfigured("object_properties") {
		t.Errorf("blocks: %+v", got.Blocks)
	}
	if !got.IsBlockAttributeConfigured("services.row_num_external_vlan") {
		t.Errorf("block attributes: %+v", got.BlockAttributes)
	}
	if !got.IsIndexedBlockAttributeConfigured("services", 1, "row_num_external_vlan") {
		t.Errorf("index 1 should have row_num_external_vlan: %+v", got.IndexedBlockAttributes)
	}
	if got.IsIndexedBlockAttributeConfigured("services", 2, "row_num_external_vlan") {
		t.Errorf("index 2 should not have row_num_external_vlan: %+v", got.IndexedBlockAttributes)
	}
}

// Concurrent lookups mirror Terraform's parallel graph walk; run under -race.
func TestConfigIndexConcurrentLookups(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	providerutils.ClearConfigIndexCache()

	for f := 0; f < 4; f++ {
		body := ""
		for i := 0; i < 50; i++ {
			body += fmt.Sprintf("resource \"verity_service\" \"svc_%d_%d\" {\n  name = \"svc_%d_%d\"\n  vlan = %d\n}\n", f, i, f, i, i)
		}
		writeTF(t, dir, fmt.Sprintf("f%d.tf", f), body)
	}

	var wg sync.WaitGroup
	for w := 0; w < 16; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				name := fmt.Sprintf("svc_%d_%d", w%4, i)
				got := providerutils.ParseResourceConfiguredAttributes(ctx, dir, "verity_service", name)
				if !got.IsConfigured("vlan") {
					t.Errorf("worker %d: %s missing vlan", w, name)
					return
				}
			}
		}(w)
	}
	wg.Wait()
}

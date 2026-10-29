package utils

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func writeIndexTF(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "a.tf"), []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestWorkDirIndexIsReused(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	ClearConfigIndexCache()

	writeIndexTF(t, dir, "resource \"verity_service\" \"svc\" {\n  name = \"svc\"\n  vlan = 100\n}\n")

	first := getWorkDirIndex(ctx, dir)
	if first == nil {
		t.Fatal("no index built")
	}
	for i := 0; i < 10; i++ {
		if got := getWorkDirIndex(ctx, dir); got != first {
			t.Fatalf("lookup %d rebuilt the index instead of reusing it", i)
		}
	}
	if len(configIndexCache) != 1 {
		t.Errorf("expected one cached directory, got %d", len(configIndexCache))
	}

	writeIndexTF(t, dir, "resource \"verity_service\" \"svc\" {\n  name = \"svc\"\n  tenant = \"t1\"\n}\n")
	second := getWorkDirIndex(ctx, dir)
	if second == first {
		t.Fatal("index was reused after the file changed")
	}
	if second.fingerprint == first.fingerprint {
		t.Error("fingerprint did not change with the file")
	}
}

func TestWorkDirIndexIsPerDirectory(t *testing.T) {
	ctx := context.Background()
	ClearConfigIndexCache()

	dirA, dirB := t.TempDir(), t.TempDir()
	writeIndexTF(t, dirA, "resource \"verity_service\" \"a\" {\n  name = \"a\"\n  vlan = 1\n}\n")
	writeIndexTF(t, dirB, "resource \"verity_service\" \"b\" {\n  name = \"b\"\n  mtu = 9000\n}\n")

	idxA, idxB := getWorkDirIndex(ctx, dirA), getWorkDirIndex(ctx, dirB)
	if idxA == idxB {
		t.Fatal("two directories shared one index")
	}
	if len(configIndexCache) != 2 {
		t.Fatalf("expected two cached directories, got %d", len(configIndexCache))
	}

	InvalidateConfigIndex(dirA)
	if _, ok := configIndexCache[dirA]; ok {
		t.Error("InvalidateConfigIndex did not drop the entry")
	}
	if _, ok := configIndexCache[dirB]; !ok {
		t.Error("InvalidateConfigIndex dropped an unrelated directory")
	}
}

func TestWorkDirIndexMissingDirectory(t *testing.T) {
	ctx := context.Background()
	ClearConfigIndexCache()

	got := ParseResourceConfiguredAttributes(ctx, filepath.Join(t.TempDir(), "absent"), "verity_service", "svc")
	if got == nil || len(got.Attributes) != 0 {
		t.Fatalf("expected an empty result, got %+v", got)
	}
}

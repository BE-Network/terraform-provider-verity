package bulkops_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/transport"
)

// The generic engine's update request for an update-only resource must match the
// one the handwritten resource sends.
//
// TestUpdateOnlyResourceGoldenPatch recorded verity_sfp_breakout's PATCH from a
// value built the way the handwritten resource builds it. The acceptance harness
// cannot hold this resource across an update, so the generic path is checked at
// the same layer: the canonical object the codec would produce for the same entry
// goes through the generated adapter and the bulk manager, and the request has to
// equal the committed fixture. It never regenerates the fixture, so it cannot
// agree with itself.
func TestGenericUpdateOnlyAdapterMatchesGoldenPatch(t *testing.T) {
	t.Parallel()

	adapter, found := transport.GeneratedAdapters["verity_sfp_breakout"]
	if !found {
		t.Fatal("verity_sfp_breakout has no generated adapter")
	}

	var (
		mu       sync.Mutex
		captured []byte
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
			}
			mu.Lock()
			captured = body
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	// The entry TestUpdateOnlyResourceGoldenPatch builds by hand, as the codec
	// represents it.
	object := transport.WireObject{
		"breakout": transport.List([]transport.WireValue{
			transport.Object(transport.WireObject{
				"index":    transport.Int64(1),
				"breakout": transport.String("1x100G"),
				"enable":   transport.Bool(true),
			}),
		}),
	}
	value, err := adapter.ResourceValue(object)
	if err != nil {
		t.Fatalf("ResourceValue: %v", err)
	}

	manager := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	manager.AddPatch(ctx, "sfp_breakout", "SFP Breakouts", value)
	if diags, _ := manager.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	mu.Lock()
	body := captured
	mu.Unlock()
	if len(body) == 0 {
		t.Fatal("no PATCH request captured for sfp_breakout")
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("decode captured PATCH body: %v", err)
	}
	encoded, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	encoded = append(encoded, '\n')

	want, err := os.ReadFile(filepath.Join("..", "lifecycle", "testdata", "golden", "verity_sfp_breakout", "patch.json"))
	if err != nil {
		t.Fatalf("read golden patch: %v", err)
	}
	if string(want) != string(encoded) {
		t.Fatalf("generic sfp_breakout patch differs from the handwritten fixture.\n--- recorded\n%s\n--- sent now\n%s", want, encoded)
	}
}

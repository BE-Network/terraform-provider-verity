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
	"terraform-provider-verity/openapi"
)

func TestUpdateOnlyResourceGoldenPatch(t *testing.T) {
	t.Parallel()

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

	manager := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")

	breakout := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{}
	breakout.SetIndex(1)
	breakout.SetBreakout("1x100G")
	breakout.SetEnable(true)
	value := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue{}
	value.SetBreakout([]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{breakout})

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

	path := filepath.Join("..", "lifecycle", "testdata", "golden", "verity_sfp_breakout", "patch.json")
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, encoded, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no golden fixture for the sfp_breakout update path: %v; regenerate with UPDATE_GOLDEN=1", err)
	}
	if string(want) != string(encoded) {
		t.Fatalf("sfp_breakout patch request differs from its golden fixture.\n--- recorded\n%s\n--- sent now\n%s", want, encoded)
	}
}

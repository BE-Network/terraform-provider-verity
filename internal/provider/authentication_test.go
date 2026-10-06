package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"terraform-provider-verity/internal/auth"
	"terraform-provider-verity/openapi"
)

func TestAuthenticatePreservesNamedCredentialsAcrossSDKRegeneration(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost || r.URL.Path != "/api/auth" {
			t.Errorf("unexpected authentication route: %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Auth struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"auth"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Auth.Username != "test-user" || body.Auth.Password != "test-password" {
			t.Error("SDK regeneration changed the transmitted username or password")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"test-session"}`))
	}))
	defer server.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	config := openapi.NewConfiguration()
	config.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	config.HTTPClient = &http.Client{Jar: jar}
	provCtx := &providerContext{
		client: openapi.NewAPIClient(config), config: config,
		tokenManager: auth.NewTokenManager(jar),
	}
	provCtx.credentials.username = "test-user"
	provCtx.credentials.password = "test-password"
	for range 2 {
		if err := authenticate(context.Background(), provCtx); err != nil {
			t.Fatal(err)
		}
	}
	if requests != 1 {
		t.Errorf("authentication requests = %d, want 1 with token reuse", requests)
	}
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	for _, cookie := range jar.Cookies(u) {
		if cookie.Name == "nbi_auth" && cookie.Value == "test-session" {
			return
		}
	}
	t.Fatal("authentication did not preserve the session cookie")
}

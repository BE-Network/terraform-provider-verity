package transport

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"terraform-provider-verity/openapi"
)

type fetchRoundTripFunc func(*http.Request) (*http.Response, error)

func (f fetchRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchDebugLoggingRedactsSessionCookieAndPreservesCollection(t *testing.T) {
	var logs bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(previous) })
	response := &http.Response{
		StatusCode: http.StatusOK, Status: "200 OK", Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
		Header: http.Header{"Content-Type": {"application/json"}, "Set-Cookie": {"nbi_auth=fetch-session-secret"}},
		Body:   io.NopCloser(strings.NewReader(`{"objects":{"example":{"name":"example"}}}`)),
	}
	config := openapi.NewConfiguration()
	config.Debug = true
	config.Servers = openapi.ServerConfigurations{{URL: "http://verity.invalid/api"}}
	config.HTTPClient = &http.Client{Transport: fetchRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return response, nil
	})}
	objects, err := FetchCollection(context.Background(), openapi.NewAPIClient(config), "example", "/objects", nil, "objects")
	if err != nil {
		t.Fatal(err)
	}
	if objects["example"].(map[string]interface{})["name"] != "example" {
		t.Fatal("logging changed the decoded collection")
	}
	if response.Header.Get("Set-Cookie") != "nbi_auth=fetch-session-secret" {
		t.Fatal("logging changed the real session cookie")
	}
	if strings.Contains(logs.String(), "fetch-session-secret") {
		t.Fatal("generic GET debug logs exposed a session credential")
	}
	if !strings.Contains(logs.String(), `"name":"example"`) || !strings.Contains(logs.String(), "[REDACTED]") {
		t.Fatal("generic GET debug output lost its response body or redaction marker")
	}
}

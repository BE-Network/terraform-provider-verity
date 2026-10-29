package sdk_test

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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(previous) })
	return &output
}

func TestAuthDebugLoggingRedactsCredentialsAndPreservesTraffic(t *testing.T) {
	for _, debug := range []bool{true, false} {
		for _, basePath := range []string{"/api", "/api/auth/session"} {
			t.Run(basePath+map[bool]string{true: "/debug", false: "/quiet"}[debug], func(t *testing.T) {
				logs := captureLogs(t)
				const username = "sdk-test-user-secret"
				const password = "sdk-test-password-secret"
				const token = "sdk-test-response-token-secret"
				const authorization = "Bearer sdk-test-authorization-secret"
				const cookie = "nbi_auth=sdk-test-cookie-secret"
				const responseBody = `{"token":"` + token + `"}`
				config := openapi.NewConfiguration()
				config.Debug = debug
				config.Servers = openapi.ServerConfigurations{{URL: "http://verity.invalid" + basePath}}
				config.DefaultHeader["Authorization"] = authorization
				config.DefaultHeader["Cookie"] = cookie
				called := false
				config.HTTPClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					called = true
					body, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Contains(body, []byte(username)) || !bytes.Contains(body, []byte(password)) {
						t.Fatal("logging changed the authentication request body")
					}
					if request.Header.Get("Authorization") != authorization || request.Header.Get("Cookie") != cookie {
						t.Fatal("logging changed the real request credentials")
					}
					return &http.Response{
						StatusCode: http.StatusOK, Status: "200 OK", Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
						Header: http.Header{"Content-Type": {"application/json"}, "Set-Cookie": {"nbi_auth=" + token}},
						Body:   io.NopCloser(strings.NewReader(responseBody)), Request: request,
					}, nil
				})}
				client := openapi.NewAPIClient(config)
				response, err := client.AuthorizationAPI.AuthPost(context.Background()).AuthPostRequest(openapi.AuthPostRequest{
					Auth: &openapi.AuthPostRequestAuth{Username: username, Password: password},
				}).Execute()
				if err != nil {
					t.Fatal(err)
				}
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				if err != nil || string(body) != responseBody || response.Header.Get("Set-Cookie") != "nbi_auth="+token || !called {
					t.Fatal("logging changed the authentication response")
				}
				for _, secret := range []string{username, password, token, authorization, cookie} {
					if strings.Contains(logs.String(), secret) {
						t.Fatal("authentication credentials appeared in SDK debug logs")
					}
				}
				if debug {
					for _, marker := range []string{"POST ", "200 OK", "[REDACTED AUTH REQUEST BODY]", "[REDACTED AUTH RESPONSE BODY]"} {
						if !strings.Contains(logs.String(), marker) {
							t.Errorf("debug logs missing %q", marker)
						}
					}
				} else if logs.Len() != 0 {
					t.Fatal("SDK logged traffic with debug disabled")
				}
			})
		}
	}
}

func TestNonAuthDebugLoggingPreservesResponseAndRedactsSessionHeaders(t *testing.T) {
	logs := captureLogs(t)
	const body = `{"version":"6.6"}`
	config := openapi.NewConfiguration()
	config.Debug = true
	config.Servers = openapi.ServerConfigurations{{URL: "http://verity.invalid/api"}}
	config.DefaultHeader["Cookie"] = "nbi_auth=request-session-secret"
	config.HTTPClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Cookie") != "nbi_auth=request-session-secret" {
			t.Fatal("session header changed on the real request")
		}
		return &http.Response{
			StatusCode: http.StatusOK, Status: "200 OK", Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
			Header: http.Header{"Content-Type": {"application/json"}, "Set-Cookie": {"nbi_auth=response-session-secret"}},
			Body:   io.NopCloser(strings.NewReader(body)), Request: request,
		}, nil
	})}
	response, err := openapi.NewAPIClient(config).VersionAPI.VersionGet(context.Background()).Execute()
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	got, err := io.ReadAll(response.Body)
	if err != nil || string(got) != body {
		t.Fatal("debug logging consumed the response body")
	}
	if !strings.Contains(logs.String(), body) || !strings.Contains(logs.String(), "GET ") {
		t.Fatal("non-auth traffic was not logged")
	}
	for _, secret := range []string{"request-session-secret", "response-session-secret"} {
		if strings.Contains(logs.String(), secret) {
			t.Fatal("session credentials appeared in SDK debug logs")
		}
	}
}

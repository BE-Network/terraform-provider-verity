package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"terraform-provider-verity/openapi"
)

func strictFetch(status int, body string) (map[string]interface{}, error) {
	config := openapi.NewConfiguration()
	config.Servers = openapi.ServerConfigurations{{URL: "http://verity.invalid/api"}}
	config.HTTPClient = &http.Client{Transport: fetchRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status, Status: http.StatusText(status),
			Header: http.Header{"Content-Type": {"application/json"}},
			Body:   io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	return FetchCollectionStrict(context.Background(), openapi.NewAPIClient(config), "sfp breakouts", "/sfpbreakouts", nil, "sfp_breakouts")
}

func TestFetchCollectionStrict(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		body      string
		objects   int
		malformed bool
		failure   bool
	}{
		{name: "collection", status: 200, body: `{"sfp_breakouts": {"SFP Breakouts": {"enable": false}}}`, objects: 1},
		{name: "empty collection means absence", status: 200, body: `{"sfp_breakouts": {}}`},
		{name: "empty object", status: 200, body: `{}`, malformed: true},
		{name: "unrelated object", status: 200, body: `{"status": "ok"}`, malformed: true},
		{name: "null collection", status: 200, body: `{"sfp_breakouts": null}`, malformed: true},
		{name: "string collection", status: 200, body: `{"sfp_breakouts": "none"}`, malformed: true},
		{name: "array collection", status: 200, body: `{"sfp_breakouts": []}`, malformed: true},
		{name: "HTTP error", status: 500, body: `{"error": "down"}`, failure: true},
		{name: "undecodable body", status: 200, body: `not json`, failure: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			objects, err := strictFetch(tc.status, tc.body)
			switch {
			case tc.malformed:
				if !errors.Is(err, ErrMalformedCollection) {
					t.Fatalf("err = %v, want ErrMalformedCollection", err)
				}
			case tc.failure:
				if err == nil || errors.Is(err, ErrMalformedCollection) {
					t.Fatalf("err = %v, want a request or decode error", err)
				}
			default:
				if err != nil || len(objects) != tc.objects {
					t.Fatalf("objects = %v, err = %v, want %d objects", objects, err, tc.objects)
				}
			}
		})
	}
}

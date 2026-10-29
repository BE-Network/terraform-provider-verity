package utils

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestFetchResourceWithRetry(t *testing.T) {
	fetchError := errors.New("fetch failed")
	want := map[string]interface{}{"resource": "value"}
	tests := []struct {
		name          string
		failures      int
		wrongType     bool
		wantAttempts  int
		wantResult    map[string]interface{}
		wantError     error
		wantErrorText string
	}{
		{name: "first attempt succeeds", wantAttempts: 1, wantResult: want},
		{name: "retry succeeds", failures: 1, wantAttempts: 2, wantResult: want},
		{name: "last attempt succeeds", failures: 2, wantAttempts: 3, wantResult: want},
		{name: "all attempts fail", failures: 3, wantAttempts: 3, wantError: fetchError},
		{name: "conversion fails", wrongType: true, wantAttempts: 3, wantErrorText: "failed to cast result to expected type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			got, err := FetchResourceWithRetry(context.Background(), nil, "cache", "resource",
				func() (map[string]interface{}, error) {
					if attempts <= tt.failures {
						return nil, fetchError
					}
					return want, nil
				},
				func(_ context.Context, _ interface{}, _ string, fetch func() (interface{}, error), _ ...bool) (interface{}, error) {
					attempts++
					if tt.wrongType {
						return "unexpected cached value", nil
					}
					return fetch()
				},
			)
			if attempts != tt.wantAttempts {
				t.Fatalf("attempts = %d, want %d", attempts, tt.wantAttempts)
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Fatalf("result = %#v, want %#v", got, tt.wantResult)
			}
			if tt.wantErrorText != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrorText) {
					t.Fatalf("error = %v, want containing %q", err, tt.wantErrorText)
				}
			} else if !errors.Is(err, tt.wantError) {
				t.Fatalf("error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

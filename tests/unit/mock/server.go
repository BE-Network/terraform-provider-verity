package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
)

type CapturedRequest struct {
	Method      string
	Path        string
	QueryParams map[string][]string
	Body        map[string]interface{}
	RawBody     []byte
}

type MockServer struct {
	Server            *httptest.Server
	mu                sync.Mutex
	requests          []CapturedRequest
	getResponses      map[string][]byte
	resourceState     map[string]map[string]map[string]interface{}
	postPutEnrichment map[string]map[string]map[string]map[string]interface{}
	testLogger        testing.TB
	versionResponse   []byte
}

func NewMockServer(mode string) *MockServer {
	ms := &MockServer{
		getResponses:  make(map[string][]byte),
		resourceState: make(map[string]map[string]map[string]interface{}),
	}

	datacenter := mode == "datacenter"
	ms.versionResponse = []byte(fmt.Sprintf(
		`{"version":"6.5","datacenter":%t}`, datacenter,
	))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.handleRequest(w, r)
	})

	ms.Server = httptest.NewServer(handler)
	return ms
}

func (ms *MockServer) LoadResponsesFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read response directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		apiPath := "/api/" + name

		if name == "acls_ipv4" || name == "acls_ipv6" {
			apiPath = "/api/acls"
		}
		ms.getResponses[apiPath+":"+name] = data
		ms.getResponses[apiPath] = data
	}
	return nil
}

func (ms *MockServer) SetGetResponse(path string, data []byte) {
	ms.getResponses[path] = data
}

func (ms *MockServer) SetPostPutEnrichment(path, wrapperKey, resourceName string, fields map[string]interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if ms.postPutEnrichment == nil {
		ms.postPutEnrichment = make(map[string]map[string]map[string]map[string]interface{})
	}
	if ms.postPutEnrichment[path] == nil {
		ms.postPutEnrichment[path] = make(map[string]map[string]map[string]interface{})
	}
	if ms.postPutEnrichment[path][wrapperKey] == nil {
		ms.postPutEnrichment[path][wrapperKey] = make(map[string]map[string]interface{})
	}
	ms.postPutEnrichment[path][wrapperKey][resourceName] = fields
}

func (ms *MockServer) GetRequests() []CapturedRequest {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	result := make([]CapturedRequest, len(ms.requests))
	copy(result, ms.requests)
	return result
}

func (ms *MockServer) GetRequestsByMethod(method string) []CapturedRequest {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	var result []CapturedRequest
	for _, req := range ms.requests {
		if req.Method == method {
			result = append(result, req)
		}
	}
	return result
}

func (ms *MockServer) GetRequestsByMethodAndPath(method, path string) []CapturedRequest {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	var result []CapturedRequest
	for _, req := range ms.requests {
		if req.Method == method && req.Path == path {
			result = append(result, req)
		}
	}
	return result
}

func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.requests = nil
}

func (ms *MockServer) ResetAll() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.requests = nil
	ms.resourceState = make(map[string]map[string]map[string]interface{})
}

func (ms *MockServer) SetTestLogger(t testing.TB) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.testLogger = t
}

func (ms *MockServer) Close() {
	ms.Server.Close()
}

func (ms *MockServer) URL() string {
	return ms.Server.URL
}

func (ms *MockServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/auth" && r.Method == http.MethodPost {
		io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"mock-test-token"}`))
		return
	}

	if r.URL.Path == "/api/version" && r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ms.versionResponse)
		return
	}

	captured := CapturedRequest{
		Method:      r.Method,
		Path:        r.URL.Path,
		QueryParams: r.URL.Query(),
	}

	if r.Method == http.MethodPut || r.Method == http.MethodPatch {
		body, err := io.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			captured.RawBody = body
			var decoded map[string]interface{}
			if json.Unmarshal(body, &decoded) == nil {
				captured.Body = decoded
			}
		}
	}

	ms.mu.Lock()
	ms.requests = append(ms.requests, captured)

	if ms.testLogger != nil {
		switch r.Method {
		case http.MethodPut, http.MethodPatch:
			if captured.Body != nil {
				formatted, _ := json.MarshalIndent(captured.Body, "  ", "  ")
				ms.testLogger.Logf("[MOCK] %s %s\n  %s", r.Method, r.URL.Path, formatted)
			}
		case http.MethodDelete:
			ms.testLogger.Logf("[MOCK] DELETE %s?%s", r.URL.Path, r.URL.RawQuery)
		}
	}

	queryParams := r.URL.Query()
	if r.Method == http.MethodPut && captured.Body != nil {
		ms.applyPutState(r.URL.Path, captured.Body, queryParams)
	} else if r.Method == http.MethodPatch && captured.Body != nil {
		ms.applyPatchState(r.URL.Path, captured.Body, queryParams)
	} else if r.Method == http.MethodDelete {
		ms.applyDeleteState(r.URL.Path, queryParams)
	}
	ms.mu.Unlock()

	if r.Method == http.MethodGet {
		ms.mu.Lock()
		response := ms.buildGetResponse(r.URL.Path, r.URL.Query())
		ms.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if captured.RawBody != nil {
		w.Write(captured.RawBody)
	} else {
		w.Write([]byte(`{"status":"success"}`))
	}
}

func (ms *MockServer) applyPutState(path string, body map[string]interface{}, queryParams map[string][]string) {
	for wrapperKey, wrapper := range body {
		resources, ok := wrapper.(map[string]interface{})
		if !ok {
			continue
		}
		if ms.resourceState[path] == nil {
			ms.resourceState[path] = make(map[string]map[string]interface{})
		}
		if ms.resourceState[path][wrapperKey] == nil {
			ms.resourceState[path][wrapperKey] = make(map[string]interface{})
		}
		for resourceName, data := range resources {
			ms.resourceState[path][wrapperKey][resourceName] = deepCopy(data)
		}

		// ACL special case
		if path == "/api/acls" && wrapperKey == "ip_filter" {
			versionKey := aclVersionKey(queryParams)
			if ms.resourceState[path][versionKey] == nil {
				ms.resourceState[path][versionKey] = make(map[string]interface{})
			}
			for resourceName, data := range resources {
				ms.resourceState[path][versionKey][resourceName] = deepCopy(data)
			}
		}
	}

	// Merge post-PUT enrichment fields (simulates auto-assigned values)
	if ms.postPutEnrichment != nil {
		if pathEnrich, ok := ms.postPutEnrichment[path]; ok {
			for wk, resources := range pathEnrich {
				if ms.resourceState[path] != nil && ms.resourceState[path][wk] != nil {
					for rn, enrichFields := range resources {
						if existing, ok := ms.resourceState[path][wk][rn].(map[string]interface{}); ok {
							for k, v := range enrichFields {
								existing[k] = v
							}
						}
					}
				}
			}
		}
	}
}

func (ms *MockServer) applyPatchState(path string, body map[string]interface{}, queryParams map[string][]string) {
	for wrapperKey, wrapper := range body {
		resources, ok := wrapper.(map[string]interface{})
		if !ok {
			continue
		}

		// For ACLs, also merge into the version-specific key so GET returns updated state
		keysToMerge := []string{wrapperKey}
		if path == "/api/acls" && wrapperKey == "ip_filter" {
			keysToMerge = append(keysToMerge, aclVersionKey(queryParams))
		}

		for resourceName, patchData := range resources {
			patchFields, ok := patchData.(map[string]interface{})
			if !ok {
				continue
			}
			for _, key := range keysToMerge {
				// Get existing state or create new
				if ms.resourceState[path] == nil || ms.resourceState[path][key] == nil {
					// No existing state to patch — treat as full set
					ms.applyPutState(path, body, queryParams)
					return
				}
				existing, ok := ms.resourceState[path][key][resourceName].(map[string]interface{})
				if !ok {
					existing = make(map[string]interface{})
				}
				// Merge patch fields into existing
				for k, v := range patchFields {
					if patchArr, ok := v.([]interface{}); ok && isIndexedArray(patchArr) {
						existingArr, _ := existing[k].([]interface{})
						existing[k] = mergeIndexedArray(existingArr, patchArr)
						continue
					}
					existing[k] = v
				}
				ms.resourceState[path][key][resourceName] = existing
			}
		}
	}
}

func isIndexedArray(arr []interface{}) bool {
	if len(arr) == 0 {
		return false
	}
	if m, ok := arr[0].(map[string]interface{}); ok {
		_, hasIndex := m["index"]
		return hasIndex
	}
	return false
}

func mergeIndexedArray(existing, patch []interface{}) []interface{} {
	indexMap := make(map[float64]map[string]interface{})
	for _, item := range existing {
		if m, ok := item.(map[string]interface{}); ok {
			if idx, ok := m["index"].(float64); ok {
				indexMap[idx] = m
			}
		}
	}

	for _, item := range patch {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		idx, ok := m["index"].(float64)
		if !ok {
			continue
		}

		if len(m) == 1 {
			// Delete: only "index" field present
			delete(indexMap, idx)
		} else if existingItem, exists := indexMap[idx]; exists {
			// Update: merge fields into existing item
			for k, v := range m {
				existingItem[k] = v
			}
		} else {
			// Add: new item
			indexMap[idx] = m
		}
	}

	indices := make([]float64, 0, len(indexMap))
	for idx := range indexMap {
		indices = append(indices, idx)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })

	result := make([]interface{}, 0, len(indices))
	for _, idx := range indices {
		result = append(result, indexMap[idx])
	}
	return result
}

func (ms *MockServer) applyDeleteState(path string, queryParams map[string][]string) {
	pathState, exists := ms.resourceState[path]
	if !exists {
		return
	}

	for _, names := range queryParams {
		for _, name := range names {
			for wrapperKey := range pathState {
				delete(pathState[wrapperKey], name)
			}
		}
	}
}

func (ms *MockServer) buildGetResponse(path string, queryParams map[string][]string) []byte {
	var baseResponse map[string]interface{}

	if path == "/api/acls" {
		if versions, ok := queryParams["ip_version"]; ok && len(versions) > 0 {
			var key string
			switch versions[0] {
			case "4":
				key = path + ":acls_ipv4"
			case "6":
				key = path + ":acls_ipv6"
			}
			if key != "" {
				if data, ok := ms.getResponses[key]; ok {
					json.Unmarshal(data, &baseResponse)
				}
			}
		}
	}

	if baseResponse == nil {
		if data, ok := ms.getResponses[path]; ok {
			json.Unmarshal(data, &baseResponse)
		}
	}
	if baseResponse == nil {
		baseResponse = make(map[string]interface{})
	}

	// Merge live state on top (field-level deep merge preserves canned data for fields not in PUT)
	if pathState, exists := ms.resourceState[path]; exists {
		for wrapperKey, resources := range pathState {
			existing, ok := baseResponse[wrapperKey].(map[string]interface{})
			if !ok {
				existing = make(map[string]interface{})
			}
			for name, data := range resources {
				// Deep merge: merge fields from live state into canned resource data
				if existingResource, ok := existing[name].(map[string]interface{}); ok {
					if newFields, ok := data.(map[string]interface{}); ok {
						for k, v := range newFields {
							existingResource[k] = v
						}
						existing[name] = existingResource
					} else {
						existing[name] = data
					}
				} else {
					existing[name] = data
				}
			}
			baseResponse[wrapperKey] = existing
		}
	}

	result, err := json.Marshal(baseResponse)
	if err != nil {
		return []byte(`{}`)
	}
	return result
}

func TestdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

func aclVersionKey(queryParams map[string][]string) string {
	if versions, ok := queryParams["ip_version"]; ok && len(versions) > 0 && versions[0] == "6" {
		return "ipv6_filter"
	}
	return "ipv4_filter"
}

func (ms *MockServer) GetOrderedPaths(method string) []string {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	seen := make(map[string]bool)
	var result []string
	for _, req := range ms.requests {
		if req.Method == method && !seen[req.Path] {
			seen[req.Path] = true
			result = append(result, req.Path)
		}
	}
	return result
}

func deepCopy(v interface{}) interface{} {
	data, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return v
	}
	return result
}

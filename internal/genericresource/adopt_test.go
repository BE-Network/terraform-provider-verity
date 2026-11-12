package genericresource

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

const adoptName = "SFP Breakouts"

type fetchResult struct {
	collection map[string]interface{}
	err        error
}

type adoptRuntime struct {
	mu        sync.Mutex
	manager   *bulkops.Manager
	fetches   []fetchResult
	fetched   int
	authErr   error
	requests  []string
	bodies    []string
	patchFail bool
}

func (f *adoptRuntime) Mode() string { return "datacenter" }
func (f *adoptRuntime) EnsureAuthenticated(context.Context) error {
	return f.authErr
}
func (*adoptRuntime) ClearCache(context.Context, string) {}
func (f *adoptRuntime) BulkManager() *bulkops.Manager    { return f.manager }
func (f *adoptRuntime) NotifyOperationAdded() {
	go f.manager.ExecuteBulk(context.Background(), "sfp_breakout", "PATCH")
}
func (f *adoptRuntime) FetchCollection(context.Context, spec.ResourceSpec, string) (map[string]interface{}, error) {
	return nil, errors.New("adoption must use the strict fetch")
}
func (f *adoptRuntime) FetchCollectionStrict(context.Context, spec.ResourceSpec) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.fetches) == 0 {
		return nil, errors.New("unexpected fetch")
	}
	next := f.fetches[0]
	if len(f.fetches) > 1 {
		f.fetches = f.fetches[1:]
	}
	f.fetched++
	return next.collection, next.err
}
func (*adoptRuntime) ConfiguredAttributes(context.Context, string, string) *utils.ConfiguredAttributes {
	return nil
}

func (f *adoptRuntime) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	f.mu.Lock()
	f.requests = append(f.requests, r.Method+" "+r.URL.Path)
	if r.Method == http.MethodPatch {
		f.bodies = append(f.bodies, string(body))
	}
	fail := f.patchFail && r.Method == http.MethodPatch
	f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if fail {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "patch rejected by the test"}`))
		return
	}
	if r.Method == http.MethodGet {
		w.Write([]byte(`{"sfp_breakouts": {}}`))
		return
	}
	w.Write([]byte(`{}`))
}

func (f *adoptRuntime) seen() ([]string, []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.requests...), append([]string(nil), f.bodies...)
}

func breakoutObject(entries ...map[string]interface{}) map[string]interface{} {
	list := make([]interface{}, 0, len(entries))
	for _, entry := range entries {
		list = append(list, entry)
	}
	return map[string]interface{}{"enable": false, "breakout": list}
}

func sfpEntry(index int, vendor string) map[string]interface{} {
	return map[string]interface{}{"enable": false, "vendor": vendor, "part_number": "", "breakout": "1x100G", "index": index}
}

func collectionWith(object map[string]interface{}) fetchResult {
	return fetchResult{collection: map[string]interface{}{adoptName: object}}
}

func setupAdopt(t *testing.T, fetches ...fetchResult) (*Resource, *adoptRuntime) {
	t.Helper()
	resourceSpec, err := registry.Lookup("verity_sfp_breakout")
	if err != nil {
		t.Fatal(err)
	}
	runtime := &adoptRuntime{fetches: fetches}
	server := httptest.NewServer(http.HandlerFunc(runtime.handle))
	t.Cleanup(server.Close)
	config := openapi.NewConfiguration()
	config.Servers = openapi.ServerConfigurations{{URL: server.URL}}
	runtime.manager = bulkops.NewManager(openapi.NewAPIClient(config), nil, nil, "datacenter")
	oldDelay, oldBackoff := bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff
	bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff = 0, 0
	t.Cleanup(func() {
		bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff = oldDelay, oldBackoff
	})
	r := &Resource{spec: resourceSpec, runtime: runtime, adapter: transport.GeneratedAdapters["verity_sfp_breakout"]}
	compiled, err := CompileSchema(resourceSpec)
	if err != nil {
		t.Fatal(err)
	}
	r.schema = compiled
	return r, runtime
}

func planValues(t *testing.T, r *Resource, planned map[string]interface{}) map[string]attr.Value {
	t.Helper()
	values, err := stateFromAPI(r.spec.Fields, planned, "datacenter", nil)
	if err != nil {
		t.Fatal(err)
	}
	values["name"] = types.StringValue(adoptName)
	return values
}

func createWithPlan(t *testing.T, r *Resource, planned map[string]interface{}) resource.CreateResponse {
	t.Helper()
	return createWithValues(t, r, planValues(t, r, planned))
}

func createWithValues(t *testing.T, r *Resource, values map[string]attr.Value) resource.CreateResponse {
	t.Helper()
	ctx := context.Background()
	state := tfsdk.State{Schema: r.schema}
	if ds := r.setState(ctx, &state, values); ds.HasError() {
		t.Fatal(ds)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: r.schema, Raw: tftypes.NewValue(r.schema.Type().TerraformType(ctx), nil)}}
	r.Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: state.Raw, Schema: r.schema},
		Config: tfsdk.Config{Raw: state.Raw, Schema: r.schema},
	}, &resp)
	return resp
}

func requireNoPutOrDelete(t *testing.T, requests []string) {
	t.Helper()
	for _, request := range requests {
		if strings.HasPrefix(request, "PUT ") || strings.HasPrefix(request, "DELETE ") {
			t.Fatalf("requests = %v: adoption must never PUT or DELETE", requests)
		}
	}
}

func patchCount(requests []string) int {
	count := 0
	for _, request := range requests {
		if strings.HasPrefix(request, "PATCH ") {
			count++
		}
	}
	return count
}

func requireAdoptedName(t *testing.T, resp resource.CreateResponse) {
	t.Helper()
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	var name types.String
	if ds := resp.State.GetAttribute(context.Background(), path.Root("name"), &name); ds.HasError() {
		t.Fatal(ds)
	}
	if name.ValueString() != adoptName {
		t.Fatalf("name = %v, want the collection key %q", name, adoptName)
	}
}

func requireFailedWithoutState(t *testing.T, resp resource.CreateResponse, want string) {
	t.Helper()
	if !resp.Diagnostics.HasError() || !strings.Contains(fmt.Sprint(resp.Diagnostics), want) {
		t.Fatalf("diagnostics = %v, want an error containing %q", resp.Diagnostics, want)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("a failed adoption must not publish any state")
	}
}

func TestAdoptUnchangedSendsNothing(t *testing.T) {
	current := breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "bbb"))
	r, runtime := setupAdopt(t, collectionWith(current))
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "bbb")))

	requireAdoptedName(t, resp)
	requests, _ := runtime.seen()
	requireNoPutOrDelete(t, requests)
	if patchCount(requests) != 0 || runtime.fetched != 1 {
		t.Fatalf("requests = %v, fetches = %d, want one read and no write", requests, runtime.fetched)
	}
}

func TestAdoptPatchesOnlyTheChangedEntry(t *testing.T) {
	r, runtime := setupAdopt(t,
		collectionWith(breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "bbb"))),
		collectionWith(breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "zzz"))))
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "zzz")))

	requireAdoptedName(t, resp)
	requests, bodies := runtime.seen()
	requireNoPutOrDelete(t, requests)
	if patchCount(requests) != 1 {
		t.Fatalf("requests = %v, want one PATCH", requests)
	}
	if !strings.Contains(bodies[0], `"zzz"`) || strings.Contains(bodies[0], `"aaa"`) {
		t.Fatalf("PATCH = %s, want only the changed entry", bodies[0])
	}
	var breakout types.List
	resp.State.GetAttribute(context.Background(), path.Root("breakout"), &breakout)
	if !strings.Contains(breakout.String(), `"zzz"`) {
		t.Fatalf("state breakout = %s, want the confirmed value", breakout)
	}
}

func TestAdoptIgnoresUnknownPlanValues(t *testing.T) {
	current := breakoutObject(sfpEntry(1, "aaa"))
	r, runtime := setupAdopt(t, collectionWith(current))
	values := planValues(t, r, breakoutObject(sfpEntry(1, "aaa")))
	properties, ok := values["object_properties"].(types.List)
	if !ok {
		t.Fatalf("object_properties = %T", values["object_properties"])
	}
	values["object_properties"] = types.ListUnknown(properties.ElementType(context.Background()))
	resp := createWithValues(t, r, values)

	requireAdoptedName(t, resp)
	if requests, _ := runtime.seen(); patchCount(requests) != 0 {
		t.Fatalf("requests = %v, want no PATCH for an unknown planned value", requests)
	}
}

func TestAdoptMissingObjectIsNotCreated(t *testing.T) {
	r, runtime := setupAdopt(t, fetchResult{collection: map[string]interface{}{}})
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa")))

	requireFailedWithoutState(t, resp, "doesn't exist on the system")
	if requests, _ := runtime.seen(); len(requests) != 0 {
		t.Fatalf("requests = %v, want none", requests)
	}
}

func TestAdoptReadErrorsAreNotAbsence(t *testing.T) {
	for _, tc := range []struct {
		name  string
		fetch fetchResult
		want  string
	}{
		{"transport error", fetchResult{err: errors.New("connection refused")}, "Failed to Read"},
		{"missing collection key", fetchResult{err: fmt.Errorf("%w: no collection", transport.ErrMalformedCollection)}, "Invalid Response"},
		{"entry is not an object", fetchResult{collection: map[string]interface{}{adoptName: "text"}}, "Invalid Response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, runtime := setupAdopt(t, tc.fetch)
			resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa")))

			requireFailedWithoutState(t, resp, tc.want)
			if strings.Contains(fmt.Sprint(resp.Diagnostics), "doesn't exist") {
				t.Fatalf("diagnostics = %v: a read error must not look like absence", resp.Diagnostics)
			}
			if requests, _ := runtime.seen(); len(requests) != 0 {
				t.Fatalf("requests = %v, want none", requests)
			}
		})
	}
}

func TestAdoptPatchFailurePublishesNoState(t *testing.T) {
	r, runtime := setupAdopt(t, collectionWith(breakoutObject(sfpEntry(1, "aaa"))))
	runtime.patchFail = true
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "zzz")))

	requireFailedWithoutState(t, resp, "patch rejected by the test")
	requests, _ := runtime.seen()
	requireNoPutOrDelete(t, requests)
}

func TestAdoptReadBackMustConfirm(t *testing.T) {
	old := collectionWith(breakoutObject(sfpEntry(1, "aaa")))
	for _, tc := range []struct {
		name     string
		readBack fetchResult
		want     string
	}{
		{"read-back fails", fetchResult{err: errors.New("connection reset")}, "reading it back failed"},
		{"read-back keeps the old value", old, "still shows different values for [breakout]"},
		{"read-back omits the object", fetchResult{collection: map[string]interface{}{}}, "missing from the API response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, runtime := setupAdopt(t, old, tc.readBack)
			resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "zzz")))

			requireFailedWithoutState(t, resp, tc.want)
			if want := 2 + bulkops.PostOperationVerificationRetries; runtime.fetched != want {
				t.Fatalf("fetches = %d, want the initial read plus %d read-backs", runtime.fetched, want-1)
			}
			requests, _ := runtime.seen()
			requireNoPutOrDelete(t, requests)

			r.runtime.(*adoptRuntime).fetches = []fetchResult{collectionWith(breakoutObject(sfpEntry(1, "zzz")))}
			before := patchCount(requests)
			retry := createWithPlan(t, r, breakoutObject(sfpEntry(1, "zzz")))
			requireAdoptedName(t, retry)
			if after, _ := runtime.seen(); patchCount(after) != before {
				t.Fatalf("requests = %v: the next apply must not PATCH settings that already match", after)
			}
		})
	}
}

func TestAdoptReadBackRetriesAStaleResponse(t *testing.T) {
	r, _ := setupAdopt(t,
		collectionWith(breakoutObject(sfpEntry(1, "aaa"))),
		collectionWith(breakoutObject(sfpEntry(1, "aaa"))),
		collectionWith(breakoutObject(sfpEntry(1, "zzz"))))
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "zzz")))
	requireAdoptedName(t, resp)
}

func TestAdoptChecksAuthenticationFirst(t *testing.T) {
	r, runtime := setupAdopt(t, collectionWith(breakoutObject(sfpEntry(1, "aaa"))))
	runtime.authErr = errors.New("bad credentials")
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa")))

	requireFailedWithoutState(t, resp, "Failed to Authenticate")
	if runtime.fetched != 0 {
		t.Fatalf("fetches = %d, want none before authentication", runtime.fetched)
	}
}

func TestAdoptValidatesIdentityFirst(t *testing.T) {
	r, runtime := setupAdopt(t, collectionWith(breakoutObject(sfpEntry(1, "aaa"))))
	ctx := context.Background()
	values, err := stateFromAPI(r.spec.Fields, breakoutObject(sfpEntry(1, "aaa")), "datacenter", nil)
	if err != nil {
		t.Fatal(err)
	}
	values["name"] = types.StringValue("")
	state := tfsdk.State{Schema: r.schema}
	if ds := r.setState(ctx, &state, values); ds.HasError() {
		t.Fatal(ds)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: r.schema, Raw: tftypes.NewValue(r.schema.Type().TerraformType(ctx), nil)}}
	r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Raw: state.Raw, Schema: r.schema}, Config: tfsdk.Config{Raw: state.Raw, Schema: r.schema}}, &resp)

	if !resp.Diagnostics.HasError() || runtime.fetched != 0 {
		t.Fatalf("diagnostics = %v, fetches = %d: an invalid identity must fail before any request", resp.Diagnostics, runtime.fetched)
	}
}

func TestCreateWithoutAdoptPolicyStillFails(t *testing.T) {
	r, runtime := setupAdopt(t, collectionWith(breakoutObject(sfpEntry(1, "aaa"))))
	r.spec.CreateBehavior = spec.CreateBehaviorAPI
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(1, "aaa")))

	requireFailedWithoutState(t, resp, "Create Not Supported")
	if runtime.fetched != 0 {
		t.Fatalf("fetches = %d, want none without the policy", runtime.fetched)
	}
}

func breakoutOrder(t *testing.T, resp resource.CreateResponse) string {
	t.Helper()
	var breakout types.List
	if ds := resp.State.GetAttribute(context.Background(), path.Root("breakout"), &breakout); ds.HasError() {
		t.Fatal(ds)
	}
	var order []string
	for _, element := range breakout.Elements() {
		object := element.(types.Object).Attributes()
		order = append(order, fmt.Sprintf("%v:%v", object["index"], strings.Trim(object["vendor"].String(), `"`)))
	}
	return strings.Join(order, ",")
}

func TestAdoptKeepsThePlannedEntryOrderWhenUnchanged(t *testing.T) {
	r, runtime := setupAdopt(t, collectionWith(breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "bbb"))))
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(2, "bbb"), sfpEntry(1, "aaa")))

	requireAdoptedName(t, resp)
	if requests, _ := runtime.seen(); patchCount(requests) != 0 {
		t.Fatalf("requests = %v, want no PATCH when only the order differs", requests)
	}
	if got := breakoutOrder(t, resp); got != "2:bbb,1:aaa" {
		t.Fatalf("state breakout = %s, want the planned order 2:bbb,1:aaa", got)
	}
}

func TestAdoptKeepsThePlannedEntryOrderAfterAPatch(t *testing.T) {
	r, runtime := setupAdopt(t,
		collectionWith(breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "bbb"))),
		collectionWith(breakoutObject(sfpEntry(1, "aaa"), sfpEntry(2, "zzz"))))
	resp := createWithPlan(t, r, breakoutObject(sfpEntry(2, "zzz"), sfpEntry(1, "aaa")))

	requireAdoptedName(t, resp)
	_, bodies := runtime.seen()
	if len(bodies) != 1 || !strings.Contains(bodies[0], `{"index":2,"vendor":"zzz"}`) || strings.Contains(bodies[0], `"aaa"`) {
		t.Fatalf("PATCH bodies = %v, want only entry 2's vendor", bodies)
	}
	if got := breakoutOrder(t, resp); got != "2:zzz,1:aaa" {
		t.Fatalf("state breakout = %s, want the planned order with the confirmed value", got)
	}
}

func TestDecodedIndexedListFollowsThePriorOrder(t *testing.T) {
	r, _ := setupAdopt(t)
	prior := planValues(t, r, breakoutObject(sfpEntry(2, "bbb"), sfpEntry(1, "aaa")))
	api := breakoutObject(sfpEntry(1, "aaa"), sfpEntry(3, "ccc"), sfpEntry(2, "bbb"))

	decoded, err := stateFromAPI(r.spec.Fields, api, "datacenter", prior)
	if err != nil {
		t.Fatal(err)
	}
	decoded["name"] = types.StringValue(adoptName)
	state := tfsdk.State{Schema: r.schema}
	if ds := r.setState(context.Background(), &state, decoded); ds.HasError() {
		t.Fatal(ds)
	}
	if got := breakoutOrder(t, resource.CreateResponse{State: state}); got != "2:bbb,1:aaa,3:ccc" {
		t.Fatalf("decoded order = %s, want the prior order first and new entries after, in API order", got)
	}
	unordered, err := stateFromAPI(r.spec.Fields, api, "datacenter", nil)
	if err != nil {
		t.Fatal(err)
	}
	unordered["name"] = types.StringValue(adoptName)
	if ds := r.setState(context.Background(), &state, unordered); ds.HasError() {
		t.Fatal(ds)
	}
	if got := breakoutOrder(t, resource.CreateResponse{State: state}); got != "1:aaa,3:ccc,2:bbb" {
		t.Fatalf("decoded order without a prior = %s, want the API order", got)
	}
}

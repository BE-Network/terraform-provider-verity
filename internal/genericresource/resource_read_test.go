package genericresource

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

type readRuntime struct {
	mode                   string
	collection             map[string]interface{}
	manager                *bulkops.Manager
	fetches, notifications int
	fetchedSpec            spec.ResourceSpec
}

func (f *readRuntime) Mode() string                            { return f.mode }
func (*readRuntime) EnsureAuthenticated(context.Context) error { return nil }
func (*readRuntime) ClearCache(context.Context, string)        {}
func (f *readRuntime) BulkManager() *bulkops.Manager           { return f.manager }
func (f *readRuntime) NotifyOperationAdded()                   { f.notifications++ }
func (f *readRuntime) FetchCollection(_ context.Context, s spec.ResourceSpec, _ string) (map[string]interface{}, error) {
	f.fetches++
	f.fetchedSpec = s
	return f.collection, nil
}
func (*readRuntime) ConfiguredAttributes(context.Context, string, string) *utils.ConfiguredAttributes {
	return nil
}

func readTestResource(t *testing.T, runtime *readRuntime) (*Resource, tfsdk.State) {
	t.Helper()
	fields := scalarFields()
	for i := range fields {
		fields[i].Modes = []spec.Mode{spec.ModeDatacenter, spec.ModeCampus}
	}
	r := &Resource{spec: spec.ResourceSpec{TerraformType: "verity_ipv4_list", IdentityPath: "name", Fields: fields, Operations: spec.OperationSpec{Create: true, Read: true, Update: true, Delete: true}, API: spec.APIResourceSpec{BulkKey: "ipv4_list"}}, runtime: runtime}
	compiled, err := CompileSchema(r.spec)
	if err != nil {
		t.Fatal(err)
	}
	r.schema = compiled
	state := tfsdk.State{Schema: compiled}
	diags := r.setState(context.Background(), &state, map[string]attr.Value{"name": types.StringValue("Gateway / Łódź"), "enable": types.BoolValue(false), "ipv4_list": types.StringNull()})
	if diags.HasError() {
		t.Fatal(diags)
	}
	return r, state
}

func TestReadCollectionIdentity(t *testing.T) {
	for _, tc := range []struct {
		name         string
		objects      map[string]interface{}
		absent, fail bool
	}{
		{name: "no inner name", objects: map[string]interface{}{"Gateway / Łódź": map[string]interface{}{"enable": true}}},
		{name: "inner name does not identify an object", objects: map[string]interface{}{"other-key": map[string]interface{}{"name": "Gateway / Łódź", "enable": true}}, absent: true},
		{name: "conflicting body", objects: map[string]interface{}{"Gateway / Łódź": map[string]interface{}{"name": "other", "enable": true}, "other": map[string]interface{}{"name": "Gateway / Łódź", "enable": false}}},
		{name: "malformed body name", objects: map[string]interface{}{"Gateway / Łódź": map[string]interface{}{"name": 42, "enable": true}}},
		{name: "absent", objects: map[string]interface{}{}, absent: true},
		{name: "invalid exact object", objects: map[string]interface{}{"Gateway / Łódź": 42}, fail: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := &readRuntime{mode: "datacenter", collection: tc.objects}
			r, state := readTestResource(t, f)
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			if resp.Diagnostics.HasError() != tc.fail {
				t.Fatal(resp.Diagnostics)
			}
			if resp.State.Raw.IsNull() != tc.absent {
				t.Fatalf("null state=%v", resp.State.Raw.IsNull())
			}
			if !tc.absent && !tc.fail {
				var name types.String
				var enabled types.Bool
				resp.State.GetAttribute(context.Background(), path.Root("name"), &name)
				resp.State.GetAttribute(context.Background(), path.Root("enable"), &enabled)
				if name.ValueString() != "Gateway / Łódź" || !enabled.ValueBool() {
					t.Fatalf("name=%v enabled=%v", name, enabled)
				}
			}
		})
	}
}

func TestReadCampusParameters(t *testing.T) {
	f := &readRuntime{mode: "campus", collection: map[string]interface{}{"Gateway / Łódź": map[string]interface{}{"enable": true}}}
	r, state := readTestResource(t, f)
	r.spec.API.FixedHeaders = map[string]string{"ip_version": "4"}
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() || resp.State.Raw.IsNull() || f.fetchedSpec.API.FixedHeaders["ip_version"] != "4" {
		t.Fatalf("response=%v fetched=%v", resp.Diagnostics, f.fetchedSpec)
	}
}

func TestCachedReadAndSettlePreserveIdentity(t *testing.T) {
	for _, bodyName := range []interface{}{nil, "conflicting", 42} {
		t.Run("body", func(t *testing.T) {

			object := map[string]interface{}{"enable": true}
			if bodyName != nil {
				object["name"] = bodyName
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.Method == http.MethodGet {
					json.NewEncoder(w).Encode(map[string]interface{}{"ipv4_list_filter": map[string]interface{}{"Gateway / Łódź": object}})
				} else {
					w.Write([]byte(`{}`))
				}
			}))
			defer server.Close()
			config := openapi.NewConfiguration()
			config.Servers = openapi.ServerConfigurations{{URL: server.URL}}
			manager := bulkops.NewManager(openapi.NewAPIClient(config), nil, nil, "datacenter")
			oldDelay, oldBackoff := bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff
			bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff = 0, 0
			defer func() {
				bulkops.ResponseProcessorDelay, bulkops.PostOperationVerificationBackoff = oldDelay, oldBackoff
			}()
			value, err := transport.GeneratedAdapters["verity_ipv4_list"].ResourceValue(transport.WireObject{"enable": transport.Bool(true)})
			if err != nil {
				t.Fatal(err)
			}
			manager.AddPatch(context.Background(), "ipv4_list", "Gateway / Łódź", value)
			if ds := manager.ExecuteBulk(context.Background(), "ipv4_list", "PATCH"); ds.HasError() {
				t.Fatal(ds)
			}
			f := &readRuntime{mode: "datacenter", manager: manager}
			r, state := readTestResource(t, f)
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			if resp.Diagnostics.HasError() || f.fetches != 0 {
				t.Fatal(resp.Diagnostics)
			}
			var name types.String
			resp.State.GetAttribute(context.Background(), path.Root("name"), &name)
			if name.ValueString() != "Gateway / Łódź" {
				t.Fatal(name)
			}
			plan, ds := readScalars(context.Background(), state, r.spec.Fields)
			if ds.HasError() {
				t.Fatal(ds)
			}
			var diagnostics diag.Diagnostics
			r.settleAfterWrite(context.Background(), "Gateway / Łódź", plan, &state, &diagnostics)
			if diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			state.GetAttribute(context.Background(), path.Root("name"), &name)
			if name.ValueString() != "Gateway / Łódź" {
				t.Fatal(name)
			}
			if bodyName == nil {
				if _, present := object["name"]; present {
					t.Fatal("cached body mutated")
				}
			}
		})
	}
}

func TestInvalidIdentityNeverLooksUpOrEnqueues(t *testing.T) {
	for _, identity := range []types.String{types.StringNull(), types.StringUnknown(), types.StringValue("")} {
		for _, op := range []string{"read", "create", "update-plan", "update-state", "delete"} {
			t.Run(op, func(t *testing.T) {
				f := &readRuntime{mode: "datacenter"}
				r, state := readTestResource(t, f)
				invalid := state
				ds := invalid.SetAttribute(context.Background(), path.Root("name"), identity)
				if ds.HasError() {
					t.Fatal(ds)
				}
				var errors bool
				switch op {
				case "read":
					resp := resource.ReadResponse{State: invalid}
					r.Read(context.Background(), resource.ReadRequest{State: invalid}, &resp)
					errors = resp.Diagnostics.HasError()
					if resp.State.Raw.IsNull() {
						t.Fatal("removed invalid state")
					}
				case "delete":
					resp := resource.DeleteResponse{State: invalid}
					r.Delete(context.Background(), resource.DeleteRequest{State: invalid}, &resp)
					errors = resp.Diagnostics.HasError()
				case "create":
					resp := resource.CreateResponse{State: state}
					r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: invalid.Raw, Schema: invalid.Schema}}, &resp)
					errors = resp.Diagnostics.HasError()
				default:
					plan, before := invalid, state
					if op == "update-state" {
						plan, before = state, invalid
					}
					resp := resource.UpdateResponse{State: before}
					r.Update(context.Background(), resource.UpdateRequest{Plan: tfsdk.Plan{Raw: plan.Raw, Schema: plan.Schema}, State: before}, &resp)
					errors = resp.Diagnostics.HasError()
				}
				if !errors || f.fetches != 0 || f.notifications != 0 {
					t.Fatalf("errors=%v fetches=%d notifications=%d", errors, f.fetches, f.notifications)
				}
			})
		}
	}
}

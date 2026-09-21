package genericresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
	_ resource.ResourceWithModifyPlan  = &Resource{}
)

// Resource is the generic lifecycle engine: one implementation that serves any
// resource the registry describes, for as much of the registry as it can yet
// reach.
//
// Supported decides what it can serve, and a spec it rejects is refused at
// construction rather than served partially. Everything this type does is
// decided by the spec it holds — the schema, what a request carries, what a
// response means, and which fields the plan nullifies — so migrating a resource
// is a matter of the registry describing it, not of code being written for it.
type Resource struct {
	spec     spec.ResourceSpec
	adapter  TransportAdapter
	schema   schema.Schema
	bind     func(providerData interface{}) (Runtime, error)
	runtime  Runtime
	typeName string
}

// New returns a factory for one registry resource. bind turns the provider's own
// data into a Runtime, which is how the engine stays independent of the package
// that registers it.
func New(resourceSpec spec.ResourceSpec, adapter TransportAdapter, bind func(providerData interface{}) (Runtime, error)) (func() resource.Resource, error) {
	// CompileSchema applies Supported, so an unsupported spec never becomes a
	// resource the provider registers.
	compiled, err := CompileSchema(resourceSpec)
	if err != nil {
		return nil, err
	}
	if adapter == nil {
		return nil, fmt.Errorf("%s: no transport adapter", resourceSpec.TerraformType)
	}
	if resourceSpec.IdentityPath == "" {
		return nil, fmt.Errorf("%s: no identity path", resourceSpec.TerraformType)
	}
	if _, found := identityField(resourceSpec); !found {
		return nil, fmt.Errorf("%s: identity %q names no field", resourceSpec.TerraformType, resourceSpec.IdentityPath)
	}
	return func() resource.Resource {
		return &Resource{
			spec:     resourceSpec,
			adapter:  adapter,
			schema:   compiled,
			bind:     bind,
			typeName: strings.TrimPrefix(resourceSpec.TerraformType, "verity_"),
		}
	}, nil
}

func identityField(resourceSpec spec.ResourceSpec) (spec.FieldSpec, bool) {
	for _, field := range resourceSpec.Fields {
		if field.TerraformName == resourceSpec.IdentityPath {
			return field, true
		}
	}
	return spec.FieldSpec{}, false
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	runtime, err := r.bind(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}
	r.runtime = runtime
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(r.spec.IdentityPath), req, resp)
}

// ModifyPlan nullifies fields that do not apply to the running mode.
//
// Without it Terraform shows "known after apply" for a Computed attribute the
// API will never return in this mode, which is noise on every plan. The legacy
// resources each carry a handwritten list of which fields those are; here the
// list is the spec's own Modes, so it cannot drift from what the registry says.
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// The resource is being destroyed; there is no plan to modify.
		return
	}
	if r.runtime == nil {
		return
	}
	mode := r.runtime.Mode()
	for _, field := range r.spec.Fields {
		if field.Unmanaged {
			continue
		}
		if !appliesToMode(field, mode) {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(field.TerraformName), nullFor(field))...)
			continue
		}
		if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
			r.nullifyOutOfModeMembers(ctx, field, mode, req, resp)
		}
	}

	if req.State.Raw.IsNull() {
		// Creating: there is no prior value for a cleared field to differ from,
		// but a value the server will assign is still unknown until it does.
		r.planAutoAssignment(ctx, req, resp)
		return
	}
	r.planExplicitNulls(ctx, req, resp)
	r.planAutoAssignment(ctx, req, resp)
}

// nullifyOutOfModeMembers applies the mode rule inside a block: a
// member the running mode does not expose is nulled in every entry the plan
// holds, so it does not show as "known after apply" either. The handwritten
// resources nullify every entry, not only the first, and so does this.
func (r *Resource) nullifyOutOfModeMembers(ctx context.Context, field spec.FieldSpec, mode string, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var block types.List
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(field.TerraformName), &block)...)
	if resp.Diagnostics.HasError() || block.IsNull() || block.IsUnknown() {
		return
	}
	r.nullifyOutOfModeEntries(ctx, path.Root(field.TerraformName), block, field, mode, resp)
}

// nullifyOutOfModeEntries nulls the out-of-mode members of every entry of one
// block, and descends into a list the block contains — verity_fabric's
// object_properties.system_graphs — so a member nested two levels down follows
// the same rule.
func (r *Resource) nullifyOutOfModeEntries(ctx context.Context, blockPath path.Path, block types.List, field spec.FieldSpec, mode string, resp *resource.ModifyPlanResponse) {
	for index, element := range block.Elements() {
		entry := blockPath.AtListIndex(index)
		object, ok := element.(types.Object)
		if !ok || object.IsNull() || object.IsUnknown() {
			continue
		}
		for _, member := range field.Fields {
			if member.Unmanaged {
				continue
			}
			if !appliesToMode(member, mode) {
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, entry.AtName(member.TerraformName), nullFor(member))...)
				continue
			}
			if member.Kind != spec.FieldKindList {
				continue
			}
			nested, ok := object.Attributes()[member.TerraformName].(types.List)
			if !ok || nested.IsNull() || nested.IsUnknown() {
				continue
			}
			r.nullifyOutOfModeEntries(ctx, entry.AtName(member.TerraformName), nested, member, mode, resp)
		}
	}
}

// planExplicitNulls makes a cleared nullable field visible as a change.
//
// Terraform copies state into the plan for an Optional and Computed attribute
// whose configuration is null, so writing `x = null` plans as no change at all
// and Update is never called. The handwritten resources detect the explicit null
// by parsing the .tf files and force the planned value to null; without the same
// step the generic engine silently ignores a clear.
//
// Only a nullable field needs this. Every other kind is cleared by a zero value
// the configuration can state outright, which plans as an ordinary change.
func (r *Resource) planExplicitNulls(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	nullable := make([]spec.FieldSpec, 0, len(r.spec.Fields))
	var lists []spec.FieldSpec
	for _, field := range r.spec.Fields {
		if field.Unmanaged {
			continue
		}
		if field.Nullable {
			nullable = append(nullable, field)
		}
		if hasNullableMember(field) {
			lists = append(lists, field)
		}
	}
	if len(nullable) == 0 && len(lists) == 0 {
		return
	}

	config, diags := readScalars(ctx, req.Config, r.spec.Fields)
	resp.Diagnostics.Append(diags...)
	state, stateDiags := readScalars(ctx, req.State, r.spec.Fields)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan map[string]attr.Value
	plan, diags = readScalars(ctx, req.Plan, r.spec.Fields)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	configured := r.runtime.ConfiguredAttributes(ctx, r.spec.TerraformType, r.identityOf(plan))

	for _, field := range nullable {
		if !configured.IsConfigured(field.TerraformName) {
			continue
		}
		configValue, held := config[field.TerraformName]
		if !held || configValue == nil || !configValue.IsNull() {
			continue
		}
		previous, hadPrevious := state[field.TerraformName]
		if !hadPrevious || previous == nil || previous.IsNull() {
			continue
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(field.TerraformName), nullOf(field.Kind))...)
	}

	for _, field := range lists {
		if field.Kind == spec.FieldKindObject {
			r.planSingletonExplicitNulls(ctx, field, config, state, configured, resp)
			continue
		}
		r.planEntryExplicitNulls(ctx, field, config, state, configured, resp)
	}
}

// planSingletonExplicitNulls applies the explicit-null rule to a singleton's
// members: when the configuration and state both hold the block, a nullable
// member written as null where state holds a value is planned as null in the
// block's one entry. It is what verity_switchpoint's handwritten plan does for
// object_properties.number_of_multipoints.
func (r *Resource) planSingletonExplicitNulls(ctx context.Context, field spec.FieldSpec, config, state map[string]attr.Value, configured *utils.ConfiguredAttributes, resp *resource.ModifyPlanResponse) {
	configMembers, configPresent, err := singletonMembers(field, config[field.TerraformName])
	if err != nil || !configPresent {
		return
	}
	stateMembers, statePresent, err := singletonMembers(field, state[field.TerraformName])
	if err != nil || !statePresent {
		return
	}
	for _, member := range field.Fields {
		if member.Unmanaged || !member.Nullable {
			continue
		}
		if !configured.IsBlockAttributeConfigured(field.TerraformName + "." + member.TerraformName) {
			continue
		}
		written, held := configMembers[member.TerraformName]
		if !held || written == nil || !written.IsNull() {
			continue
		}
		prior, had := stateMembers[member.TerraformName]
		if !had || prior == nil || prior.IsNull() {
			continue
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx,
			path.Root(field.TerraformName).AtListIndex(0).AtName(member.TerraformName), nullOf(member.Kind))...)
	}
}

// planEntryExplicitNulls applies the same rule inside a list. Each configuration
// entry is matched to the state entry with its index, and a nullable member
// written as null where state holds a value is planned as null at that entry's
// position. The handwritten resources read the index the same way: an unknown or
// null one reads as zero.
func (r *Resource) planEntryExplicitNulls(ctx context.Context, field spec.FieldSpec, config, state map[string]attr.Value, configured *utils.ConfiguredAttributes, resp *resource.ModifyPlanResponse) {
	configEntries, present, err := listEntries(field, config[field.TerraformName])
	if err != nil || !present {
		return
	}
	stateEntries, _, err := listEntries(field, state[field.TerraformName])
	if err != nil {
		return
	}
	for position, configEntry := range configEntries {
		var index int64
		if value, ok := configEntry[field.Collection.IdentityField].(types.Int64); ok {
			index = value.ValueInt64()
		}
		var before map[string]attr.Value
		for _, candidate := range stateEntries {
			if value, ok := candidate[field.Collection.IdentityField].(types.Int64); ok && value.ValueInt64() == index {
				before = candidate
				break
			}
		}
		if before == nil {
			continue
		}
		for _, member := range field.Fields {
			if member.Unmanaged || !member.Nullable {
				continue
			}
			if !configured.IsIndexedBlockAttributeConfigured(field.TerraformName, index, member.TerraformName) {
				continue
			}
			written, held := configEntry[member.TerraformName]
			if !held || written == nil || !written.IsNull() {
				continue
			}
			prior, had := before[member.TerraformName]
			if !had || prior == nil || prior.IsNull() {
				continue
			}
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx,
				path.Root(field.TerraformName).AtListIndex(position).AtName(member.TerraformName), nullOf(member.Kind))...)
		}
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan, diagnostics := readScalars(ctx, req.Plan, r.spec.Fields)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !r.ready(&resp.Diagnostics) {
		return
	}
	if !r.spec.Operations.Create {
		// Some objects exist only as hardware or configuration the API manages;
		// they are imported and updated, never created.
		resp.Diagnostics.AddError(
			"Create Not Supported",
			fmt.Sprintf("%s cannot be created through this API; import it instead", r.spec.TerraformType),
		)
		return
	}
	if err := r.runtime.EnsureAuthenticated(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to Authenticate", fmt.Sprintf("Error authenticating with API: %s", err))
		return
	}

	name := r.identityOf(plan)
	object, err := buildCreate(r.spec.Fields, plan, r.nullableSource(ctx, req.Config, name, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err))
		return
	}
	value, err := r.adapter.ResourceValue(object)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Build Request", fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err))
		return
	}

	if !bulkops.ExecuteResourceOperationWithOptions(ctx, r.runtime.BulkManager(), r.runtime.NotifyOperationAdded,
		"create", r.spec.API.BulkKey, name, value, &resp.Diagnostics, r.operationOptions()) {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("%s %s creation operation completed successfully", r.spec.TerraformType, name))
	r.runtime.ClearCache(ctx, r.spec.API.CacheKey)

	r.settleAfterWrite(ctx, name, plan, &resp.State, &resp.Diagnostics)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan, diagnostics := readScalars(ctx, req.Plan, r.spec.Fields)
	resp.Diagnostics.Append(diagnostics...)
	state, stateDiagnostics := readScalars(ctx, req.State, r.spec.Fields)
	resp.Diagnostics.Append(stateDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !r.ready(&resp.Diagnostics) {
		return
	}
	if err := r.runtime.EnsureAuthenticated(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to Authenticate", fmt.Sprintf("Error authenticating with API: %s", err))
		return
	}

	name := r.identityOf(plan)
	object, changed, err := buildUpdate(r.spec.Fields, plan, state, r.nullableSource(ctx, req.Config, name, &resp.Diagnostics), &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err))
		return
	}
	// A reference pair reports a refused combination through diagnostics rather
	// than an error, in the words the handwritten helpers use.
	if resp.Diagnostics.HasError() {
		return
	}
	if !changed {
		// Nothing to send. The plan is already what the object is, so it becomes
		// state unchanged rather than provoking an empty request.
		resp.Diagnostics.Append(r.setState(ctx, &resp.State, plan)...)
		return
	}
	value, err := r.adapter.ResourceValue(object)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Build Request", fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err))
		return
	}

	if !bulkops.ExecuteResourceOperationWithOptions(ctx, r.runtime.BulkManager(), r.runtime.NotifyOperationAdded,
		"update", r.spec.API.BulkKey, name, value, &resp.Diagnostics, r.operationOptions()) {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("%s %s update operation completed successfully", r.spec.TerraformType, name))
	r.runtime.ClearCache(ctx, r.spec.API.CacheKey)

	r.settleAfterWrite(ctx, name, plan, &resp.State, &resp.Diagnostics)
}

// nullableSource reads the configuration and the .tf files, but only when the
// resource actually has a nullable field. Parsing is not free, and a resource
// with none has nothing to learn from it.
func (r *Resource) nullableSource(ctx context.Context, config tfsdk.Config, name string, diagnostics *diag.Diagnostics) nullableSource {
	// The configuration is needed for a nullable field, and for an
	// auto-assignment flag, which is sent only when the configuration states it.
	if !r.needsConfiguration() {
		return nullableSource{}
	}
	values, diags := readScalars(ctx, config, r.spec.Fields)
	diagnostics.Append(diags...)
	attributes := r.runtime.ConfiguredAttributes(ctx, r.spec.TerraformType, name)
	return nullableSource{
		config:     values,
		configured: attributes.IsConfigured,
		indexed:    attributes.IsIndexedBlockAttributeConfigured,
		block:      attributes.IsBlockAttributeConfigured,
	}
}

// needsConfiguration reports whether any field's request depends on the
// configuration itself: a nullable field, a nullable member of a list entry, or
// an auto-assignment flag, which is sent only when the configuration states it.
func (r *Resource) needsConfiguration() bool {
	for _, field := range r.spec.Fields {
		if field.Unmanaged {
			continue
		}
		if field.Nullable || field.AutoAssignment != nil || hasNullableMember(field) {
			return true
		}
	}
	return false
}

func hasNullableMember(field spec.FieldSpec) bool {
	if field.Kind != spec.FieldKindList && field.Kind != spec.FieldKindObject {
		return false
	}
	for _, member := range field.Fields {
		if member.Nullable && !member.Unmanaged {
			return true
		}
	}
	return false
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state, diagnostics := readScalars(ctx, req.State, r.spec.Fields)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !r.ready(&resp.Diagnostics) {
		return
	}
	if !r.spec.Operations.Delete {
		resp.Diagnostics.AddError(
			"Delete Not Supported",
			fmt.Sprintf("%s cannot be deleted through this API", r.spec.TerraformType),
		)
		return
	}
	if err := r.runtime.EnsureAuthenticated(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to Authenticate", fmt.Sprintf("Error authenticating with API: %s", err))
		return
	}

	name := r.identityOf(state)
	if !bulkops.ExecuteResourceOperationWithOptions(ctx, r.runtime.BulkManager(), r.runtime.NotifyOperationAdded,
		"delete", r.spec.API.BulkKey, name, nil, &resp.Diagnostics, r.operationOptions()) {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("%s %s deletion operation completed successfully", r.spec.TerraformType, name))
	r.runtime.ClearCache(ctx, r.spec.API.CacheKey)
	resp.State.RemoveResource(ctx)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state, diagnostics := readScalars(ctx, req.State, r.spec.Fields)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !r.ready(&resp.Diagnostics) {
		return
	}
	if err := r.runtime.EnsureAuthenticated(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to Authenticate", fmt.Sprintf("Error authenticating with API: %s", err))
		return
	}

	name := r.identityOf(state)
	manager := r.runtime.BulkManager()

	// A write this provider just made is authoritative, and re-reading it would
	// race the server's own settling. Both branches below are why: the response
	// the operation returned, then the knowledge that one is still in flight.
	if manager != nil {
		if data, exists := manager.GetResourceResponse(r.spec.API.BulkKey, name); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached %s data for %s from recent operation", r.spec.TerraformType, name))
			r.applyResponse(ctx, applyFallback(ctx, data), state, &resp.State, &resp.Diagnostics)
			return
		}
		if manager.HasPendingOrRecentOperations(r.spec.API.BulkKey) {
			tflog.Info(ctx, fmt.Sprintf("Skipping %s %s verification - trusting recent successful API operation", r.spec.TerraformType, name))
			if handled, diags := setFallbackState(ctx, &resp.State); handled {
				resp.Diagnostics.Append(diags...)
			}
			return
		}
	}

	collection, err := r.runtime.FetchCollection(ctx, r.spec, name)
	if err != nil {
		resp.Diagnostics.Append(utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read %s %s", r.spec.TerraformType, name))...)
		return
	}

	data, actualAPIName, exists := utils.FindResourceByAPIName(collection, name, func(candidate interface{}) (string, bool) {
		object, ok := candidate.(map[string]interface{})
		if !ok {
			return "", false
		}
		identity, ok := object[r.spec.IdentityPath].(string)
		return identity, ok
	})
	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("%s with name '%s' not found in API response", r.spec.TerraformType, name))
		resp.State.RemoveResource(ctx)
		return
	}
	object, ok := data.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Response",
			fmt.Sprintf("%s data is not in the expected format for %s", r.spec.TerraformType, name),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Found %s '%s' under API key '%s'", r.spec.TerraformType, name, actualAPIName))

	r.applyResponse(ctx, applyFallback(ctx, object), state, &resp.State, &resp.Diagnostics)
}

// settleAfterWrite turns a completed write into state. The operation's own
// response is preferred, since it is what the server stored; otherwise the
// resource is read back, and if even that yields nothing the plan stands in so
// the apply does not end with no state at all.
func (r *Resource) settleAfterWrite(ctx context.Context, name string, plan map[string]attr.Value, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	// State has to exist before it can be refined, and at this point the write
	// has happened: losing the identity here would orphan the object.
	minimal := map[string]attr.Value{}
	for _, field := range r.spec.Fields {
		if field.Unmanaged {
			continue
		}
		minimal[field.TerraformName] = nullFor(field)
	}
	minimal[r.spec.IdentityPath] = types.StringValue(name)
	diagnostics.Append(r.setState(ctx, state, minimal)...)
	if diagnostics.HasError() {
		return
	}

	if manager := r.runtime.BulkManager(); manager != nil {
		if data, exists := manager.GetResourceResponse(r.spec.API.BulkKey, name); exists {
			merged := r.mergePlanScalars(data, plan)
			values, err := stateFromAPI(r.spec.Fields, merged, r.runtime.Mode(), plan)
			if err != nil {
				diagnostics.AddError("Invalid Response", fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err))
				return
			}
			diagnostics.Append(r.setState(ctx, state, values)...)
			return
		}
	}

	readResp := resource.ReadResponse{State: *state, Diagnostics: *diagnostics}
	r.Read(r.withFallback(ctx, plan), resource.ReadRequest{State: *state}, &readResp)
	if readResp.State.Raw.IsNull() {
		// The read found nothing to write, so the plan is all that is left to
		// record. Without this the apply ends having made a change it did not save.
		readResp.Diagnostics.Append(r.setState(ctx, &readResp.State, nullifyUnknown(r.spec.Fields, plan))...)
	}
	*state = readResp.State
	*diagnostics = readResp.Diagnostics
}

// applyResponse decodes an API object into state, keeping the prior state for
// the fields whose absence policy asks for it.
func (r *Resource) applyResponse(ctx context.Context, data map[string]interface{}, prior map[string]attr.Value, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	values, err := stateFromAPI(r.spec.Fields, data, r.runtime.Mode(), prior)
	if err != nil {
		diagnostics.AddError("Invalid Response", fmt.Sprintf("%s: %s", r.spec.TerraformType, err))
		return
	}
	diagnostics.Append(r.setState(ctx, state, values)...)
}

// mergePlanScalars fills fields the response left out from the plan that was
// just applied. The server echoes what it stored, and a field it omits is one it
// accepted without comment, so the planned value is the better answer than null.
func (r *Resource) mergePlanScalars(data map[string]interface{}, plan map[string]attr.Value) map[string]interface{} {
	merged := make(map[string]interface{}, len(data))
	for key, value := range data {
		merged[key] = value
	}
	mode := r.runtime.Mode()
	for _, field := range r.spec.Fields {
		if field.Unmanaged {
			continue
		}
		if !appliesToMode(field, mode) {
			continue
		}
		if _, present := merged[field.APIName]; present {
			continue
		}
		planned, held := plan[field.TerraformName]
		if !held || planned == nil || planned.IsNull() || planned.IsUnknown() {
			continue
		}
		if scalar, ok := goValue(field.Kind, planned); ok {
			merged[field.APIName] = scalar
		}
	}
	return merged
}

func (r *Resource) setState(ctx context.Context, state *tfsdk.State, values map[string]attr.Value) diag.Diagnostics {
	object, diagnostics := types.ObjectValue(attributeTypes(r.spec.Fields), values)
	if diagnostics.HasError() {
		return diagnostics
	}
	return state.Set(ctx, object)
}

func (r *Resource) identityOf(values map[string]attr.Value) string {
	identity, held := values[r.spec.IdentityPath]
	if !held {
		return ""
	}
	if text, ok := identity.(types.String); ok && !text.IsNull() && !text.IsUnknown() {
		return text.ValueString()
	}
	return ""
}

// operationOptions carries the fixed parameters that select this resource on a
// shared endpoint, which the bulk manager adds to every write and uses to keep
// resources that share an endpoint in separate batches. The handwritten ACL
// resources pass the same map.
func (r *Resource) operationOptions() *bulkops.ResourceOperationOptions {
	if len(r.spec.API.FixedHeaders) == 0 {
		return nil
	}
	return &bulkops.ResourceOperationOptions{HeaderParams: r.spec.API.FixedHeaders}
}

func (r *Resource) ready(diagnostics *diag.Diagnostics) bool {
	if r.runtime == nil {
		diagnostics.AddError(
			"Provider Not Configured",
			fmt.Sprintf("%s was used before the provider finished configuring", r.spec.TerraformType),
		)
		return false
	}
	return true
}

// goValue converts a known Framework value back to the Go value an API response
// would have carried, so a planned value can stand in for an absent member.
func goValue(kind spec.FieldKind, value attr.Value) (interface{}, bool) {
	switch kind {
	case spec.FieldKindString:
		typed, ok := value.(types.String)
		return typed.ValueString(), ok
	case spec.FieldKindBool:
		typed, ok := value.(types.Bool)
		return typed.ValueBool(), ok
	case spec.FieldKindInt64:
		typed, ok := value.(types.Int64)
		return typed.ValueInt64(), ok
	case spec.FieldKindNumber:
		typed, ok := value.(types.Number)
		if !ok {
			return nil, false
		}
		result, _ := typed.ValueBigFloat().Float64()
		return result, true
	default:
		return nil, false
	}
}

// nullifyUnknown replaces unknowns with nulls so a plan can be written to state.
// An unknown cannot be stored, and the field is Computed, so null is what the
// next read will fill in.
func nullifyUnknown(fields []spec.FieldSpec, values map[string]attr.Value) map[string]attr.Value {
	settled := make(map[string]attr.Value, len(values))
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		value, held := values[field.TerraformName]
		if !held || value == nil || value.IsUnknown() {
			settled[field.TerraformName] = nullFor(field)
			continue
		}
		if field.Kind == spec.FieldKindObject {
			settled[field.TerraformName] = settleSingleton(field, value)
			continue
		}
		if field.Kind == spec.FieldKindList {
			settled[field.TerraformName] = settleList(field, value)
			continue
		}
		settled[field.TerraformName] = value
	}
	return settled
}

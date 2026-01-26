package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityACLUnifiedResource{}
	_ resource.ResourceWithConfigure   = &verityACLUnifiedResource{}
	_ resource.ResourceWithImportState = &verityACLUnifiedResource{}
	_ resource.ResourceWithModifyPlan  = &verityACLUnifiedResource{}
)

const aclResourceType = "acls"

func NewVerityACLV4Resource() resource.Resource {
	return &verityACLUnifiedResource{ipVersion: "4"}
}

func NewVerityACLV6Resource() resource.Resource {
	return &verityACLUnifiedResource{ipVersion: "6"}
}

type verityACLUnifiedResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
	ipVersion            string // "4" for IPv4, "6" for IPv6
}

type verityACLUnifiedResourceModel struct {
	Name                    types.String                            `tfsdk:"name"`
	Enable                  types.Bool                              `tfsdk:"enable"`
	Protocol                types.String                            `tfsdk:"protocol"`
	Bidirectional           types.Bool                              `tfsdk:"bidirectional"`
	SourceIP                types.String                            `tfsdk:"source_ip"`
	SourcePortOperator      types.String                            `tfsdk:"source_port_operator"`
	SourcePort1             types.Int64                             `tfsdk:"source_port_1"`
	SourcePort2             types.Int64                             `tfsdk:"source_port_2"`
	DestinationIP           types.String                            `tfsdk:"destination_ip"`
	DestinationPortOperator types.String                            `tfsdk:"destination_port_operator"`
	DestinationPort1        types.Int64                             `tfsdk:"destination_port_1"`
	DestinationPort2        types.Int64                             `tfsdk:"destination_port_2"`
	ObjectProperties        []verityACLUnifiedObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityACLUnifiedObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityACLUnifiedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl_v" + r.ipVersion
}

func (r *verityACLUnifiedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *verityACLUnifiedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	ipVersionDesc := "IPv4"
	ipVersionName := "IPv4"
	if r.ipVersion == "6" {
		ipVersionDesc = "IPv6"
		ipVersionName = "IPv6"
	}

	resp.Schema = schema.Schema{
		Description: fmt.Sprintf("Manages a Verity %s IP Filter (ACL)", ipVersionName),
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object.",
				Optional:    true,
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets.  Value IP will match all IP protocols.",
				Optional:    true,
				Computed:    true,
			},
			"bidirectional": schema.BoolAttribute{
				Description: "If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet.",
				Optional:    true,
				Computed:    true,
			},
			"source_ip": schema.StringAttribute{
				Description: fmt.Sprintf("This field matches the source IP address of an %s packet", ipVersionDesc),
				Optional:    true,
				Computed:    true,
			},
			"source_port_operator": schema.StringAttribute{
				Description: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:    true,
				Computed:    true,
			},
			"source_port_1": schema.Int64Attribute{
				Description: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:    true,
				Computed:    true,
			},
			"source_port_2": schema.Int64Attribute{
				Description: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:    true,
				Computed:    true,
			},
			"destination_ip": schema.StringAttribute{
				Description: fmt.Sprintf("This field matches the destination IP address of an %s packet.", ipVersionDesc),
				Optional:    true,
				Computed:    true,
			},
			"destination_port_operator": schema.StringAttribute{
				Description: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:    true,
				Computed:    true,
			},
			"destination_port_1": schema.Int64Attribute{
				Description: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:    true,
				Computed:    true,
			},
			"destination_port_2": schema.Int64Attribute{
				Description: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Additional properties for this object",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User notes",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityACLUnifiedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityACLUnifiedResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityACLUnifiedResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	aclProps := &openapi.AclsPutRequestIpFilterValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Protocol", APIField: &aclProps.Protocol, TFValue: plan.Protocol},
		{FieldName: "SourceIp", APIField: &aclProps.SourceIp, TFValue: plan.SourceIP},
		{FieldName: "SourcePortOperator", APIField: &aclProps.SourcePortOperator, TFValue: plan.SourcePortOperator},
		{FieldName: "DestinationIp", APIField: &aclProps.DestinationIp, TFValue: plan.DestinationIP},
		{FieldName: "DestinationPortOperator", APIField: &aclProps.DestinationPortOperator, TFValue: plan.DestinationPortOperator},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &aclProps.Enable, TFValue: plan.Enable},
		{FieldName: "Bidirectional", APIField: &aclProps.Bidirectional, TFValue: plan.Bidirectional},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_acl_v"+r.ipVersion, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "SourcePort1", APIField: &aclProps.SourcePort1, TFValue: config.SourcePort1, IsConfigured: configuredAttrs.IsConfigured("source_port_1")},
		{FieldName: "SourcePort2", APIField: &aclProps.SourcePort2, TFValue: config.SourcePort2, IsConfigured: configuredAttrs.IsConfigured("source_port_2")},
		{FieldName: "DestinationPort1", APIField: &aclProps.DestinationPort1, TFValue: config.DestinationPort1, IsConfigured: configuredAttrs.IsConfigured("destination_port_1")},
		{FieldName: "DestinationPort2", APIField: &aclProps.DestinationPort2, TFValue: config.DestinationPort2, IsConfigured: configuredAttrs.IsConfigured("destination_port_2")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		aclProps.ObjectProperties = &objProps
	}

	options := &bulkops.ResourceOperationOptions{
		HeaderParams: map[string]string{"ip_version": r.ipVersion},
	}

	success := bulkops.ExecuteResourceOperationWithOptions(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "acl", name, *aclProps, &resp.Diagnostics, options)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s creation operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())

	var minState verityACLUnifiedResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aclData, exists := bulkMgr.GetResourceResponse("acl_v"+r.ipVersion, name); exists {
			state := populateACLState(ctx, minState, aclData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityACLUnifiedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityACLUnifiedResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	aclName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if aclData, exists := r.bulkOpsMgr.GetResourceResponse("acl_v"+r.ipVersion, aclName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached IPv%s ACL data for %s from recent operation", r.ipVersion, aclName))
			state = populateACLState(ctx, state, aclData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv%s ACL %s verification – trusting recent successful API operation", r.ipVersion, aclName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv%s ACLs for verification of %s", r.ipVersion, aclName))

	type ACLsResponse struct {
		ACLs map[string]interface{}
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, r.getCacheKey(), fmt.Sprintf("IPv%s ACLs", r.ipVersion),
		func() (ACLsResponse, error) {
			tflog.Debug(ctx, fmt.Sprintf("Making API call to fetch IPv%s ACLs", r.ipVersion))

			req := r.client.ACLsAPI.AclsGet(ctx).IpVersion(r.ipVersion)
			apiResp, err := req.Execute()
			if err != nil {
				return ACLsResponse{}, fmt.Errorf("error reading IPv%s ACL: %v", r.ipVersion, err)
			}
			defer apiResp.Body.Close()

			var rawResponse map[string]interface{}
			if err := json.NewDecoder(apiResp.Body).Decode(&rawResponse); err != nil {
				return ACLsResponse{}, fmt.Errorf("failed to decode IPv%s ACLs response: %v", r.ipVersion, err)
			}

			// Extract the correct field based on IP version
			var filterKey string
			if r.ipVersion == "6" {
				filterKey = "ipv6_filter"
			} else {
				filterKey = "ipv4_filter"
			}

			if ipFilter, ok := rawResponse[filterKey].(map[string]interface{}); ok {
				tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv%s ACLs", len(ipFilter), r.ipVersion))
				return ACLsResponse{ACLs: ipFilter}, nil
			}

			return ACLsResponse{ACLs: make(map[string]interface{})}, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv%s ACL %s", r.ipVersion, aclName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv%s ACL with name: %s", r.ipVersion, aclName))

	// ACLs use the map key as the name
	aclData, exists := utils.FindResourceByKey(result.ACLs, aclName)
	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv%s ACL with name '%s' not found in API response", r.ipVersion, aclName))
		resp.State.RemoveResource(ctx)
		return
	}

	aclDataMap, ok := aclData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid ACL Data",
			fmt.Sprintf("IPv%s ACL data is not in expected format for %s", r.ipVersion, aclName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found IPv%s ACL '%s'", r.ipVersion, aclName))

	state = populateACLState(ctx, state, aclDataMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityACLUnifiedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityACLUnifiedResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	aclProps := openapi.AclsPutRequestIpFilterValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { aclProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Protocol, state.Protocol, func(v *string) { aclProps.Protocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SourceIP, state.SourceIP, func(v *string) { aclProps.SourceIp = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SourcePortOperator, state.SourcePortOperator, func(v *string) { aclProps.SourcePortOperator = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DestinationIP, state.DestinationIP, func(v *string) { aclProps.DestinationIp = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DestinationPortOperator, state.DestinationPortOperator, func(v *string) { aclProps.DestinationPortOperator = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { aclProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Bidirectional, state.Bidirectional, func(v *bool) { aclProps.Bidirectional = v }, &hasChanges)

	// Handle nullable int64 field changes (ports)
	utils.CompareAndSetNullableInt64Field(plan.SourcePort1, state.SourcePort1, func(v *openapi.NullableInt32) { aclProps.SourcePort1 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.SourcePort2, state.SourcePort2, func(v *openapi.NullableInt32) { aclProps.SourcePort2 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DestinationPort1, state.DestinationPort1, func(v *openapi.NullableInt32) { aclProps.DestinationPort1 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DestinationPort2, state.DestinationPort2, func(v *openapi.NullableInt32) { aclProps.DestinationPort2 = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			aclProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	options := &bulkops.ResourceOperationOptions{
		HeaderParams: map[string]string{"ip_version": r.ipVersion},
	}
	success := bulkops.ExecuteResourceOperationWithOptions(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "acl", name, aclProps, &resp.Diagnostics, options)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s update operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())

	var minState verityACLUnifiedResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aclData, exists := bulkMgr.GetResourceResponse("acl_v"+r.ipVersion, name); exists {
			newState := populateACLState(ctx, minState, aclData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityACLUnifiedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityACLUnifiedResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	options := &bulkops.ResourceOperationOptions{
		HeaderParams: map[string]string{"ip_version": r.ipVersion},
	}
	success := bulkops.ExecuteResourceOperationWithOptions(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "acl", name, nil, &resp.Diagnostics, options)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s deletion operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())
	resp.State.RemoveResource(ctx)
}

func (r *verityACLUnifiedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityACLUnifiedResource) getCacheKey() string {
	return "acls_ipv" + r.ipVersion
}

func populateACLState(ctx context.Context, state verityACLUnifiedResourceModel, data map[string]interface{}, mode string) verityACLUnifiedResourceModel {
	const resourceType = aclResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.Protocol = utils.MapStringWithMode(data, "protocol", resourceType, mode)
	state.SourceIP = utils.MapStringWithMode(data, "source_ip", resourceType, mode)
	state.SourcePortOperator = utils.MapStringWithMode(data, "source_port_operator", resourceType, mode)
	state.DestinationIP = utils.MapStringWithMode(data, "destination_ip", resourceType, mode)
	state.DestinationPortOperator = utils.MapStringWithMode(data, "destination_port_operator", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Bidirectional = utils.MapBoolWithMode(data, "bidirectional", resourceType, mode)

	// Int fields
	state.SourcePort1 = utils.MapInt64WithMode(data, "source_port_1", resourceType, mode)
	state.SourcePort2 = utils.MapInt64WithMode(data, "source_port_2", resourceType, mode)
	state.DestinationPort1 = utils.MapInt64WithMode(data, "destination_port_1", resourceType, mode)
	state.DestinationPort2 = utils.MapInt64WithMode(data, "destination_port_2", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityACLUnifiedObjectPropertiesModel{
				Notes: utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
			}
			state.ObjectProperties = []verityACLUnifiedObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityACLUnifiedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityACLUnifiedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = aclResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"protocol",
		"source_ip", "source_port_operator",
		"destination_ip", "destination_port_operator",
	)

	nullifier.NullifyBools(
		"enable", "bidirectional",
	)

	nullifier.NullifyInt64s(
		"source_port_1", "source_port_2",
		"destination_port_1", "destination_port_2",
	)

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityACLUnifiedResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityACLUnifiedResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable Int64 fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_acl_v"+r.ipVersion, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "source_port_1", ConfigVal: config.SourcePort1, StateVal: state.SourcePort1},
			{AttrName: "source_port_2", ConfigVal: config.SourcePort2, StateVal: state.SourcePort2},
			{AttrName: "destination_port_1", ConfigVal: config.DestinationPort1, StateVal: state.DestinationPort1},
			{AttrName: "destination_port_2", ConfigVal: config.DestinationPort2, StateVal: state.DestinationPort2},
		},
	})
}

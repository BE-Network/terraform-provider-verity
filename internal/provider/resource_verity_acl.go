package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityACLUnifiedResource{}
	_ resource.ResourceWithConfigure   = &verityACLUnifiedResource{}
	_ resource.ResourceWithImportState = &verityACLUnifiedResource{}
)

func NewVerityACLV4Resource() resource.Resource {
	return &verityACLUnifiedResource{ipVersion: "4"}
}

func NewVerityACLV6Resource() resource.Resource {
	return &verityACLUnifiedResource{ipVersion: "6"}
}

type verityACLUnifiedResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
			},
			"protocol": schema.StringAttribute{
				Description: "Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets.  Value IP will match all IP protocols.",
				Optional:    true,
			},
			"bidirectional": schema.BoolAttribute{
				Description: "If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet.",
				Optional:    true,
			},
			"source_ip": schema.StringAttribute{
				Description: fmt.Sprintf("This field matches the source IP address of an %s packet", ipVersionDesc),
				Optional:    true,
			},
			"source_port_operator": schema.StringAttribute{
				Description: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:    true,
			},
			"source_port_1": schema.Int64Attribute{
				Description: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:    true,
			},
			"source_port_2": schema.Int64Attribute{
				Description: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:    true,
			},
			"destination_ip": schema.StringAttribute{
				Description: fmt.Sprintf("This field matches the destination IP address of an %s packet.", ipVersionDesc),
				Optional:    true,
			},
			"destination_port_operator": schema.StringAttribute{
				Description: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:    true,
			},
			"destination_port_1": schema.Int64Attribute{
				Description: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:    true,
			},
			"destination_port_2": schema.Int64Attribute{
				Description: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:    true,
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

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "SourcePort1", APIField: &aclProps.SourcePort1, TFValue: plan.SourcePort1},
		{FieldName: "SourcePort2", APIField: &aclProps.SourcePort2, TFValue: plan.SourcePort2},
		{FieldName: "DestinationPort1", APIField: &aclProps.DestinationPort1, TFValue: plan.DestinationPort1},
		{FieldName: "DestinationPort2", APIField: &aclProps.DestinationPort2, TFValue: plan.DestinationPort2},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		objProps := &openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !plan.ObjectProperties[0].Notes.IsNull() {
			objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
		} else {
			objProps.Notes = nil
		}
		aclProps.ObjectProperties = objProps
	}

	// Special handling for dual resource types
	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddAclPut(ctx, name, *aclProps, r.ipVersion)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv%s ACL creation operation %s to complete", r.ipVersion, operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to create IPv%s ACL %s", r.ipVersion, name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s creation operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv%s ACL %s verification – trusting recent successful API operation", r.ipVersion, aclName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv%s ACLs for verification of %s", r.ipVersion, aclName))

	var err error
	maxRetries := 3
	var aclsMap map[string]interface{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		aclsData, fetchErr := getCachedResponse(ctx, r.provCtx, r.getCacheKey(), func() (interface{}, error) {
			tflog.Debug(ctx, fmt.Sprintf("Making API call to fetch IPv%s ACLs", r.ipVersion))

			req := r.client.ACLsAPI.AclsGet(ctx).IpVersion(r.ipVersion)
			apiResp, err := req.Execute()

			if err != nil {
				return nil, fmt.Errorf("error reading IPv%s ACL: %v", r.ipVersion, err)
			}
			defer apiResp.Body.Close()

			var rawResponse map[string]interface{}
			if err := json.NewDecoder(apiResp.Body).Decode(&rawResponse); err != nil {
				return nil, fmt.Errorf("failed to decode IPv%s ACLs response: %v", r.ipVersion, err)
			}

			// Extract the correct field based on IP version
			var filterKey string
			if r.ipVersion == "6" {
				filterKey = "ipv6_filter"
			} else {
				filterKey = "ipv4_filter"
			}

			if ipFilter, ok := rawResponse[filterKey].(map[string]interface{}); ok {
				tflog.Debug(ctx, fmt.Sprintf("Successfully fetched IPv%s ACLs from API", r.ipVersion))
				return ipFilter, nil
			}

			return make(map[string]interface{}), nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch IPv%s ACLs on attempt %d, retrying in %v", r.ipVersion, attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		aclsMap = aclsData.(map[string]interface{})
		break
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv%s ACL %s", r.ipVersion, aclName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv%s ACL with name: %s", r.ipVersion, aclName))

	// ACLs use the map key as the name
	aclData, exists := aclsMap[aclName]
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

	state.Name = types.StringValue(aclName)

	// Handle object properties
	if objProps, ok := aclDataMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityACLUnifiedObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objProps["notes"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"protocol":                  &state.Protocol,
		"source_ip":                 &state.SourceIP,
		"source_port_operator":      &state.SourcePortOperator,
		"destination_ip":            &state.DestinationIP,
		"destination_port_operator": &state.DestinationPortOperator,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(aclDataMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":        &state.Enable,
		"bidirectional": &state.Bidirectional,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(aclDataMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"source_port_1":      &state.SourcePort1,
		"source_port_2":      &state.SourcePort2,
		"destination_port_1": &state.DestinationPort1,
		"destination_port_2": &state.DestinationPort2,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(aclDataMap[apiKey])
	}

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
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			aclProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	// Special handling for dual resource types
	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddAclPatch(ctx, name, aclProps, r.ipVersion)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv%s ACL update operation %s to complete", r.ipVersion, operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv%s ACL %s", r.ipVersion, name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s update operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

	// Special handling for dual resource types
	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddAclDelete(ctx, name, r.ipVersion)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv%s ACL deletion operation %s to complete", r.ipVersion, operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv%s ACL %s", r.ipVersion, name))...,
		)
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

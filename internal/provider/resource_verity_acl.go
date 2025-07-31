package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"encoding/json"
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
		MarkdownDescription: fmt.Sprintf("Manages a Verity %s IP Filter (ACL)", ipVersionName),
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Object Name. Must be unique.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable object.",
				Optional:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets.  Value IP will match all IP protocols.",
				Optional:            true,
			},
			"bidirectional": schema.BoolAttribute{
				MarkdownDescription: "If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet.",
				Optional:            true,
			},
			"source_ip": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("This field matches the source IP address of an %s packet", ipVersionDesc),
				Optional:            true,
			},
			"source_port_operator": schema.StringAttribute{
				MarkdownDescription: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:            true,
			},
			"source_port_1": schema.Int64Attribute{
				MarkdownDescription: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:            true,
			},
			"source_port_2": schema.Int64Attribute{
				MarkdownDescription: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:            true,
			},
			"destination_ip": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("This field matches the destination IP address of an %s packet.", ipVersionDesc),
				Optional:            true,
			},
			"destination_port_operator": schema.StringAttribute{
				MarkdownDescription: "This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.",
				Optional:            true,
			},
			"destination_port_1": schema.Int64Attribute{
				MarkdownDescription: "This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.",
				Optional:            true,
			},
			"destination_port_2": schema.Int64Attribute{
				MarkdownDescription: "This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				MarkdownDescription: "Additional properties for this object",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							MarkdownDescription: "User notes",
							Optional:            true,
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
	aclProps := openapi.NewAclsPutRequestIpFilterValue()

	if !plan.Enable.IsNull() {
		aclProps.SetEnable(plan.Enable.ValueBool())
	}
	if !plan.Protocol.IsNull() {
		aclProps.SetProtocol(plan.Protocol.ValueString())
	}
	if !plan.Bidirectional.IsNull() {
		aclProps.SetBidirectional(plan.Bidirectional.ValueBool())
	}
	if !plan.SourceIP.IsNull() {
		aclProps.SetSourceIp(plan.SourceIP.ValueString())
	}
	if !plan.SourcePortOperator.IsNull() {
		aclProps.SetSourcePortOperator(plan.SourcePortOperator.ValueString())
	}
	if !plan.SourcePort1.IsNull() {
		aclProps.SetSourcePort1(int32(plan.SourcePort1.ValueInt64()))
	}
	if !plan.SourcePort2.IsNull() {
		aclProps.SetSourcePort2(int32(plan.SourcePort2.ValueInt64()))
	}
	if !plan.DestinationIP.IsNull() {
		aclProps.SetDestinationIp(plan.DestinationIP.ValueString())
	}
	if !plan.DestinationPortOperator.IsNull() {
		aclProps.SetDestinationPortOperator(plan.DestinationPortOperator.ValueString())
	}
	if !plan.DestinationPort1.IsNull() {
		aclProps.SetDestinationPort1(int32(plan.DestinationPort1.ValueInt64()))
	}
	if !plan.DestinationPort2.IsNull() {
		aclProps.SetDestinationPort2(int32(plan.DestinationPort2.ValueInt64()))
	}

	if len(plan.ObjectProperties) > 0 && !plan.ObjectProperties[0].Notes.IsNull() {
		objProps := openapi.NewAclsPutRequestIpFilterValueObjectProperties()
		objProps.SetNotes(plan.ObjectProperties[0].Notes.ValueString())
		aclProps.SetObjectProperties(*objProps)
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddAclPut(ctx, name, *aclProps, r.ipVersion)
	r.provCtx.NotifyOperationAdded()

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
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

	tflog.Debug(ctx, fmt.Sprintf("Reading IPv%s ACL resource", r.ipVersion))

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	aclName := state.Name.ValueString()

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv%s ACL %s verification - trusting recent successful API operation", r.ipVersion, aclName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent IPv%s ACL operations found, performing normal verification for %s", r.ipVersion, aclName))

	var err error
	maxRetries := 3
	var aclsMap map[string]interface{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		aclsData, fetchErr := getCachedResponse(ctx, provCtx, r.getCacheKey(), func() (interface{}, error) {
			tflog.Debug(ctx, fmt.Sprintf("Making API call to fetch IPv%s ACLs", r.ipVersion))

			req := provCtx.client.ACLsAPI.AclsGet(ctx).IpVersion(r.ipVersion)
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
	var aclData map[string]interface{}
	exists := false

	if data, ok := aclsMap[aclName]; ok {
		if aclProps, ok := data.(map[string]interface{}); ok {
			aclData = aclProps
			exists = true
			tflog.Debug(ctx, fmt.Sprintf("Found IPv%s ACL with name: %s", r.ipVersion, aclName))
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv%s ACL with name '%s' not found in API response", r.ipVersion, aclName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(aclName)

	if val, ok := aclData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else {
		state.Enable = types.BoolNull()
	}

	if val, ok := aclData["protocol"].(string); ok {
		state.Protocol = types.StringValue(val)
	} else {
		state.Protocol = types.StringNull()
	}

	if val, ok := aclData["bidirectional"].(bool); ok {
		state.Bidirectional = types.BoolValue(val)
	} else {
		state.Bidirectional = types.BoolNull()
	}

	if val, ok := aclData["source_ip"].(string); ok {
		state.SourceIP = types.StringValue(val)
	} else {
		state.SourceIP = types.StringNull()
	}

	if val, ok := aclData["source_port_operator"].(string); ok {
		state.SourcePortOperator = types.StringValue(val)
	} else {
		state.SourcePortOperator = types.StringNull()
	}

	if val, ok := aclData["source_port_1"]; ok {
		switch v := val.(type) {
		case float64:
			state.SourcePort1 = types.Int64Value(int64(v))
		case int:
			state.SourcePort1 = types.Int64Value(int64(v))
		case int32:
			state.SourcePort1 = types.Int64Value(int64(v))
		default:
			state.SourcePort1 = types.Int64Null()
		}
	} else {
		state.SourcePort1 = types.Int64Null()
	}

	if val, ok := aclData["source_port_2"]; ok {
		switch v := val.(type) {
		case float64:
			state.SourcePort2 = types.Int64Value(int64(v))
		case int:
			state.SourcePort2 = types.Int64Value(int64(v))
		case int32:
			state.SourcePort2 = types.Int64Value(int64(v))
		default:
			state.SourcePort2 = types.Int64Null()
		}
	} else {
		state.SourcePort2 = types.Int64Null()
	}

	if val, ok := aclData["destination_ip"].(string); ok {
		state.DestinationIP = types.StringValue(val)
	} else {
		state.DestinationIP = types.StringNull()
	}

	if val, ok := aclData["destination_port_operator"].(string); ok {
		state.DestinationPortOperator = types.StringValue(val)
	} else {
		state.DestinationPortOperator = types.StringNull()
	}

	if val, ok := aclData["destination_port_1"]; ok {
		switch v := val.(type) {
		case float64:
			state.DestinationPort1 = types.Int64Value(int64(v))
		case int:
			state.DestinationPort1 = types.Int64Value(int64(v))
		case int32:
			state.DestinationPort1 = types.Int64Value(int64(v))
		default:
			state.DestinationPort1 = types.Int64Null()
		}
	} else {
		state.DestinationPort1 = types.Int64Null()
	}

	if val, ok := aclData["destination_port_2"]; ok {
		switch v := val.(type) {
		case float64:
			state.DestinationPort2 = types.Int64Value(int64(v))
		case int:
			state.DestinationPort2 = types.Int64Value(int64(v))
		case int32:
			state.DestinationPort2 = types.Int64Value(int64(v))
		default:
			state.DestinationPort2 = types.Int64Null()
		}
	} else {
		state.DestinationPort2 = types.Int64Null()
	}

	if objProps, ok := aclData["object_properties"].(map[string]interface{}); ok {
		objPropsModel := verityACLUnifiedObjectPropertiesModel{}
		if notes, ok := objProps["notes"].(string); ok {
			objPropsModel.Notes = types.StringValue(notes)
		} else {
			objPropsModel.Notes = types.StringNull()
		}
		state.ObjectProperties = []verityACLUnifiedObjectPropertiesModel{objPropsModel}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityACLUnifiedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityACLUnifiedResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()
	aclProps := &openapi.AclsPutRequestIpFilterValue{}
	hasChanges := false

	objPropsChanged := false
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		if !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objPropsChanged = true
		}
	} else if len(plan.ObjectProperties) != len(state.ObjectProperties) {
		objPropsChanged = true
	}

	if objPropsChanged {
		if len(plan.ObjectProperties) > 0 {
			objProps := openapi.NewAclsPutRequestIpFilterValueObjectProperties()
			objProps.SetNotes(plan.ObjectProperties[0].Notes.ValueString())
			aclProps.SetObjectProperties(*objProps)
		}
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		aclProps.SetEnable(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Protocol.Equal(state.Protocol) {
		aclProps.SetProtocol(plan.Protocol.ValueString())
		hasChanges = true
	}

	if !plan.Bidirectional.Equal(state.Bidirectional) {
		aclProps.SetBidirectional(plan.Bidirectional.ValueBool())
		hasChanges = true
	}

	if !plan.SourceIP.Equal(state.SourceIP) {
		aclProps.SetSourceIp(plan.SourceIP.ValueString())
		hasChanges = true
	}

	if !plan.SourcePortOperator.Equal(state.SourcePortOperator) {
		aclProps.SetSourcePortOperator(plan.SourcePortOperator.ValueString())
		hasChanges = true
	}

	if !plan.SourcePort1.Equal(state.SourcePort1) {
		if !plan.SourcePort1.IsNull() {
			aclProps.SetSourcePort1(int32(plan.SourcePort1.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.SourcePort2.Equal(state.SourcePort2) {
		if !plan.SourcePort2.IsNull() {
			aclProps.SetSourcePort2(int32(plan.SourcePort2.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.DestinationIP.Equal(state.DestinationIP) {
		aclProps.SetDestinationIp(plan.DestinationIP.ValueString())
		hasChanges = true
	}

	if !plan.DestinationPortOperator.Equal(state.DestinationPortOperator) {
		aclProps.SetDestinationPortOperator(plan.DestinationPortOperator.ValueString())
		hasChanges = true
	}

	if !plan.DestinationPort1.Equal(state.DestinationPort1) {
		if !plan.DestinationPort1.IsNull() {
			aclProps.SetDestinationPort1(int32(plan.DestinationPort1.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.DestinationPort2.Equal(state.DestinationPort2) {
		if !plan.DestinationPort2.IsNull() {
			aclProps.SetDestinationPort2(int32(plan.DestinationPort2.ValueInt64()))
		}
		hasChanges = true
	}

	if !hasChanges {
		tflog.Info(ctx, fmt.Sprintf("No changes detected for IPv%s ACL %s, skipping update", r.ipVersion, name))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddAclPatch(ctx, name, *aclProps, r.ipVersion)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv%s ACL update operation %s to complete", r.ipVersion, operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv%s ACL %s", r.ipVersion, name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s update operation completed successfully", r.ipVersion, name))
	clearCache(ctx, r.provCtx, r.getCacheKey())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	aclName := state.Name.ValueString()
	bulkOpsMgr := r.provCtx.bulkOpsMgr

	operationID := bulkOpsMgr.AddAclDelete(ctx, aclName, r.ipVersion)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv%s ACL deletion operation %s to complete", r.ipVersion, operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv%s ACL %s", r.ipVersion, aclName))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv%s ACL %s deletion operation completed successfully", r.ipVersion, aclName))
	clearCache(ctx, r.provCtx, r.getCacheKey())
	resp.State.RemoveResource(ctx)
}

func (r *verityACLUnifiedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityACLUnifiedResource) getCacheKey() string {
	return "acls_ipv" + r.ipVersion
}

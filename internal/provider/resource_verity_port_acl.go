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
	_ resource.Resource                = &verityPortAclResource{}
	_ resource.ResourceWithConfigure   = &verityPortAclResource{}
	_ resource.ResourceWithImportState = &verityPortAclResource{}
)

func NewVerityPortAclResource() resource.Resource {
	return &verityPortAclResource{}
}

type verityPortAclResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityPortAclResourceModel struct {
	Name       types.String                   `tfsdk:"name"`
	Enable     types.Bool                     `tfsdk:"enable"`
	Ipv4Permit []verityPortAclIpv4PermitModel `tfsdk:"ipv4_permit"`
	Ipv4Deny   []verityPortAclIpv4DenyModel   `tfsdk:"ipv4_deny"`
	Ipv6Permit []verityPortAclIpv6PermitModel `tfsdk:"ipv6_permit"`
	Ipv6Deny   []verityPortAclIpv6DenyModel   `tfsdk:"ipv6_deny"`
}

type verityPortAclIpv4PermitModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

type verityPortAclIpv4DenyModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

type verityPortAclIpv6PermitModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

type verityPortAclIpv6DenyModel struct {
	Enable        types.Bool   `tfsdk:"enable"`
	Filter        types.String `tfsdk:"filter"`
	FilterRefType types.String `tfsdk:"filter_ref_type_"`
	Index         types.Int64  `tfsdk:"index"`
}

func (r *verityPortAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_acl"
}

func (r *verityPortAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPortAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Port ACL",
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
		},
		Blocks: map[string]schema.Block{
			"ipv4_permit": schema.ListNestedBlock{
				Description: "List of IPv4 permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv4_deny": schema.ListNestedBlock{
				Description: "List of IPv4 deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv6_permit": schema.ListNestedBlock{
				Description: "List of IPv6 permit filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"ipv6_deny": schema.ListNestedBlock{
				Description: "List of IPv6 deny filters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"filter": schema.StringAttribute{
							Description: "Filter",
							Optional:    true,
						},
						"filter_ref_type_": schema.StringAttribute{
							Description: "Object type for filter field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPortAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPortAclResourceModel
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
	portAclProps := &openapi.PortaclsPutRequestPortAclValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		portAclProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	// Handle IPv4 Permit
	if len(plan.Ipv4Permit) > 0 {
		ipv4PermitSlice := make([]openapi.PortaclsPutRequestPortAclValueIpv4PermitInner, len(plan.Ipv4Permit))
		for i, permit := range plan.Ipv4Permit {
			permitItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{}
			if !permit.Enable.IsNull() {
				permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			}
			if !permit.Filter.IsNull() {
				permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
			}
			if !permit.FilterRefType.IsNull() {
				permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
			}
			if !permit.Index.IsNull() {
				permitItem.Index = openapi.PtrInt32(int32(permit.Index.ValueInt64()))
			}
			ipv4PermitSlice[i] = permitItem
		}
		portAclProps.Ipv4Permit = ipv4PermitSlice
	}

	// Handle IPv4 Deny
	if len(plan.Ipv4Deny) > 0 {
		ipv4DenySlice := make([]openapi.PortaclsPutRequestPortAclValueIpv4PermitInner, len(plan.Ipv4Deny))
		for i, deny := range plan.Ipv4Deny {
			denyItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{}
			if !deny.Enable.IsNull() {
				denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			}
			if !deny.Filter.IsNull() {
				denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
			}
			if !deny.FilterRefType.IsNull() {
				denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
			}
			ipv4DenySlice[i] = denyItem
		}
		portAclProps.Ipv4Deny = ipv4DenySlice
	}

	// Handle IPv6 Permit
	if len(plan.Ipv6Permit) > 0 {
		ipv6PermitSlice := make([]openapi.PortaclsPutRequestPortAclValueIpv6PermitInner, len(plan.Ipv6Permit))
		for i, permit := range plan.Ipv6Permit {
			permitItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{}
			if !permit.Enable.IsNull() {
				permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			}
			if !permit.Filter.IsNull() {
				permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
			}
			if !permit.FilterRefType.IsNull() {
				permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
			}
			if !permit.Index.IsNull() {
				permitItem.Index = openapi.PtrInt32(int32(permit.Index.ValueInt64()))
			}
			ipv6PermitSlice[i] = permitItem
		}
		portAclProps.Ipv6Permit = ipv6PermitSlice
	}

	// Handle IPv6 Deny
	if len(plan.Ipv6Deny) > 0 {
		ipv6DenySlice := make([]openapi.PortaclsPutRequestPortAclValueIpv6PermitInner, len(plan.Ipv6Deny))
		for i, deny := range plan.Ipv6Deny {
			denyItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{}
			if !deny.Enable.IsNull() {
				denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			}
			if !deny.Filter.IsNull() {
				denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
			}
			if !deny.FilterRefType.IsNull() {
				denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
			}
			if !deny.Index.IsNull() {
				denyItem.Index = openapi.PtrInt32(int32(deny.Index.ValueInt64()))
			}
			ipv6DenySlice[i] = denyItem
		}
		portAclProps.Ipv6Deny = ipv6DenySlice
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "port_acl", name, *portAclProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for port ACL creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Port ACL %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Port ACL %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityPortAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPortAclResourceModel
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("port_acl") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Port ACL %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Port ACLs for verification of %s", name))

	type PortAclResponse struct {
		PortAcl map[string]map[string]interface{} `json:"port_acl"`
	}

	var result PortAclResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		portAclsData, fetchErr := getCachedResponse(ctx, r.provCtx, "port_acls", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch port ACLs")
			req := r.client.PortACLsAPI.PortaclsGet(ctx)
			resp, err := req.Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading port ACLs: %v", err)
			}
			defer resp.Body.Close()

			var res PortAclResponse
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode port ACL response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d port ACLs", len(res.PortAcl)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch port ACLs on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = portAclsData.(PortAclResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Port ACL %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for port ACL with ID: %s", name))
	var portAclData map[string]interface{}
	exists := false

	if data, ok := result.PortAcl[name]; ok {
		portAclData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found port ACL directly by ID: %s", name))
	} else {
		for apiName, p := range result.PortAcl {
			if nameVal, ok := p["name"].(string); ok && nameVal == name {
				portAclData = p
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found port ACL with name '%s' under API key '%s'", nameVal, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Port ACL with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", portAclData["name"]))

	if enable, ok := portAclData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	// Handle IPv4 Permit
	if ipv4PermitData, ok := portAclData["ipv4_permit"].([]interface{}); ok && len(ipv4PermitData) > 0 {
		var ipv4Permits []verityPortAclIpv4PermitModel
		for _, item := range ipv4PermitData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				permit := verityPortAclIpv4PermitModel{}

				if enable, ok := itemMap["enable"].(bool); ok {
					permit.Enable = types.BoolValue(enable)
				} else {
					permit.Enable = types.BoolNull()
				}
				if filter, ok := itemMap["filter"].(string); ok {
					permit.Filter = types.StringValue(filter)
				} else {
					permit.Filter = types.StringNull()
				}
				if filterRefType, ok := itemMap["filter_ref_type_"].(string); ok {
					permit.FilterRefType = types.StringValue(filterRefType)
				} else {
					permit.FilterRefType = types.StringNull()
				}
				if index, exists := itemMap["index"]; exists && index != nil {
					switch v := index.(type) {
					case int:
						permit.Index = types.Int64Value(int64(v))
					case float64:
						permit.Index = types.Int64Value(int64(v))
					default:
						permit.Index = types.Int64Null()
					}
				} else {
					permit.Index = types.Int64Null()
				}
				ipv4Permits = append(ipv4Permits, permit)
			}
		}
		state.Ipv4Permit = ipv4Permits
	} else {
		state.Ipv4Permit = nil
	}

	// Handle IPv4 Deny
	if ipv4DenyData, ok := portAclData["ipv4_deny"].([]interface{}); ok && len(ipv4DenyData) > 0 {
		var ipv4Denies []verityPortAclIpv4DenyModel
		for _, item := range ipv4DenyData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				deny := verityPortAclIpv4DenyModel{}

				if enable, ok := itemMap["enable"].(bool); ok {
					deny.Enable = types.BoolValue(enable)
				} else {
					deny.Enable = types.BoolNull()
				}
				if filter, ok := itemMap["filter"].(string); ok {
					deny.Filter = types.StringValue(filter)
				} else {
					deny.Filter = types.StringNull()
				}
				if filterRefType, ok := itemMap["filter_ref_type_"].(string); ok {
					deny.FilterRefType = types.StringValue(filterRefType)
				} else {
					deny.FilterRefType = types.StringNull()
				}
				if index, exists := itemMap["index"]; exists && index != nil {
					switch v := index.(type) {
					case int:
						deny.Index = types.Int64Value(int64(v))
					case float64:
						deny.Index = types.Int64Value(int64(v))
					default:
						deny.Index = types.Int64Null()
					}
				} else {
					deny.Index = types.Int64Null()
				}
				ipv4Denies = append(ipv4Denies, deny)
			}
		}
		state.Ipv4Deny = ipv4Denies
	} else {
		state.Ipv4Deny = nil
	}

	// Handle IPv6 Permit
	if ipv6PermitData, ok := portAclData["ipv6_permit"].([]interface{}); ok && len(ipv6PermitData) > 0 {
		var ipv6Permits []verityPortAclIpv6PermitModel
		for _, item := range ipv6PermitData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				permit := verityPortAclIpv6PermitModel{}

				if enable, ok := itemMap["enable"].(bool); ok {
					permit.Enable = types.BoolValue(enable)
				} else {
					permit.Enable = types.BoolNull()
				}
				if filter, ok := itemMap["filter"].(string); ok {
					permit.Filter = types.StringValue(filter)
				} else {
					permit.Filter = types.StringNull()
				}
				if filterRefType, ok := itemMap["filter_ref_type_"].(string); ok {
					permit.FilterRefType = types.StringValue(filterRefType)
				} else {
					permit.FilterRefType = types.StringNull()
				}
				if index, exists := itemMap["index"]; exists && index != nil {
					switch v := index.(type) {
					case int:
						permit.Index = types.Int64Value(int64(v))
					case float64:
						permit.Index = types.Int64Value(int64(v))
					default:
						permit.Index = types.Int64Null()
					}
				} else {
					permit.Index = types.Int64Null()
				}
				ipv6Permits = append(ipv6Permits, permit)
			}
		}
		state.Ipv6Permit = ipv6Permits
	} else {
		state.Ipv6Permit = nil
	}

	// Handle IPv6 Deny
	if ipv6DenyData, ok := portAclData["ipv6_deny"].([]interface{}); ok && len(ipv6DenyData) > 0 {
		var ipv6Denies []verityPortAclIpv6DenyModel
		for _, item := range ipv6DenyData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				deny := verityPortAclIpv6DenyModel{}

				if enable, ok := itemMap["enable"].(bool); ok {
					deny.Enable = types.BoolValue(enable)
				} else {
					deny.Enable = types.BoolNull()
				}
				if filter, ok := itemMap["filter"].(string); ok {
					deny.Filter = types.StringValue(filter)
				} else {
					deny.Filter = types.StringNull()
				}
				if filterRefType, ok := itemMap["filter_ref_type_"].(string); ok {
					deny.FilterRefType = types.StringValue(filterRefType)
				} else {
					deny.FilterRefType = types.StringNull()
				}
				if index, exists := itemMap["index"]; exists && index != nil {
					switch v := index.(type) {
					case int:
						deny.Index = types.Int64Value(int64(v))
					case float64:
						deny.Index = types.Int64Value(int64(v))
					default:
						deny.Index = types.Int64Null()
					}
				} else {
					deny.Index = types.Int64Null()
				}
				ipv6Denies = append(ipv6Denies, deny)
			}
		}
		state.Ipv6Deny = ipv6Denies
	} else {
		state.Ipv6Deny = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPortAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPortAclResourceModel

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
	portAclProps := openapi.PortaclsPutRequestPortAclValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		portAclProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		portAclProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	// Handle IPv4 Permit
	oldIpv4PermitsByIndex := make(map[int64]verityPortAclIpv4PermitModel)
	for _, permit := range state.Ipv4Permit {
		if !permit.Index.IsNull() {
			oldIpv4PermitsByIndex[permit.Index.ValueInt64()] = permit
		}
	}

	var changedIpv4Permits []openapi.PortaclsPutRequestPortAclValueIpv4PermitInner
	ipv4PermitsChanged := false

	for _, permit := range plan.Ipv4Permit {
		if permit.Index.IsNull() {
			continue
		}

		index := permit.Index.ValueInt64()
		oldPermit, exists := oldIpv4PermitsByIndex[index]

		if !exists {
			// new permit, include all fields
			permitItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				permit.Filter, permit.FilterRefType,
				"filter", "filter_ref_type_",
				true, false) {
				return
			}

			if !permit.Enable.IsNull() {
				permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			} else {
				permitItem.Enable = openapi.PtrBool(false)
			}

			if !permit.Filter.IsNull() {
				permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
			} else {
				permitItem.Filter = openapi.PtrString("")
			}

			if !permit.FilterRefType.IsNull() {
				permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
			} else {
				permitItem.FilterRefType = openapi.PtrString("")
			}

			changedIpv4Permits = append(changedIpv4Permits, permitItem)
			ipv4PermitsChanged = true
			continue
		}

		// existing permit, check which fields changed
		permitItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !permit.Enable.Equal(oldPermit.Enable) {
			permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			fieldChanged = true
		}

		filterChanged := !permit.Filter.Equal(oldPermit.Filter)
		filterRefTypeChanged := !permit.FilterRefType.Equal(oldPermit.FilterRefType)

		if filterChanged || filterRefTypeChanged {
			// Validate using "one ref type supported" rules
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				permit.Filter, permit.FilterRefType,
				"filter", "filter_ref_type_",
				filterChanged, filterRefTypeChanged) {
				return
			}

			if filterChanged && !filterRefTypeChanged {
				// Just send the base field
				if !permit.Filter.IsNull() {
					permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
				} else {
					permitItem.Filter = openapi.PtrString("")
				}
				fieldChanged = true
			} else if filterRefTypeChanged {
				// Send both fields
				if !permit.Filter.IsNull() {
					permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
				} else {
					permitItem.Filter = openapi.PtrString("")
				}

				if !permit.FilterRefType.IsNull() {
					permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
				} else {
					permitItem.FilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		if fieldChanged {
			changedIpv4Permits = append(changedIpv4Permits, permitItem)
			ipv4PermitsChanged = true
		}
	}

	for idx := range oldIpv4PermitsByIndex {
		found := false
		for _, permit := range plan.Ipv4Permit {
			if !permit.Index.IsNull() && permit.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// permit removed - include only the index for deletion
			deletedPermit := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedIpv4Permits = append(changedIpv4Permits, deletedPermit)
			ipv4PermitsChanged = true
		}
	}

	if ipv4PermitsChanged && len(changedIpv4Permits) > 0 {
		portAclProps.Ipv4Permit = changedIpv4Permits
		hasChanges = true
	}

	// Handle IPv4 Deny
	oldIpv4DeniesByIndex := make(map[int64]verityPortAclIpv4DenyModel)
	for _, deny := range state.Ipv4Deny {
		if !deny.Index.IsNull() {
			oldIpv4DeniesByIndex[deny.Index.ValueInt64()] = deny
		}
	}

	var changedIpv4Denies []openapi.PortaclsPutRequestPortAclValueIpv4PermitInner
	ipv4DeniesChanged := false

	for _, deny := range plan.Ipv4Deny {
		if deny.Index.IsNull() {
			continue
		}

		index := deny.Index.ValueInt64()
		oldDeny, exists := oldIpv4DeniesByIndex[index]

		if !exists {
			// new deny, include all fields
			denyItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			// Handle ref_type validation
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				deny.Filter, deny.FilterRefType,
				"filter", "filter_ref_type_",
				true, false) {
				return
			}

			if !deny.Enable.IsNull() {
				denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			} else {
				denyItem.Enable = openapi.PtrBool(false)
			}

			if !deny.Filter.IsNull() {
				denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
			} else {
				denyItem.Filter = openapi.PtrString("")
			}

			if !deny.FilterRefType.IsNull() {
				denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
			} else {
				denyItem.FilterRefType = openapi.PtrString("")
			}

			changedIpv4Denies = append(changedIpv4Denies, denyItem)
			ipv4DeniesChanged = true
			continue
		}

		// existing deny, check which fields changed
		denyItem := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !deny.Enable.Equal(oldDeny.Enable) {
			denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			fieldChanged = true
		}

		filterChanged := !deny.Filter.Equal(oldDeny.Filter)
		filterRefTypeChanged := !deny.FilterRefType.Equal(oldDeny.FilterRefType)

		if filterChanged || filterRefTypeChanged {
			// Validate using "one ref type supported" rules
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				deny.Filter, deny.FilterRefType,
				"filter", "filter_ref_type_",
				filterChanged, filterRefTypeChanged) {
				return
			}

			if filterChanged && !filterRefTypeChanged {
				// Just send the base field
				if !deny.Filter.IsNull() {
					denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
				} else {
					denyItem.Filter = openapi.PtrString("")
				}
				fieldChanged = true
			} else if filterRefTypeChanged {
				// Send both fields
				if !deny.Filter.IsNull() {
					denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
				} else {
					denyItem.Filter = openapi.PtrString("")
				}

				if !deny.FilterRefType.IsNull() {
					denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
				} else {
					denyItem.FilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		if fieldChanged {
			changedIpv4Denies = append(changedIpv4Denies, denyItem)
			ipv4DeniesChanged = true
		}
	}

	for idx := range oldIpv4DeniesByIndex {
		found := false
		for _, deny := range plan.Ipv4Deny {
			if !deny.Index.IsNull() && deny.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// deny removed - include only the index for deletion
			deletedDeny := openapi.PortaclsPutRequestPortAclValueIpv4PermitInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedIpv4Denies = append(changedIpv4Denies, deletedDeny)
			ipv4DeniesChanged = true
		}
	}

	if ipv4DeniesChanged && len(changedIpv4Denies) > 0 {
		portAclProps.Ipv4Deny = changedIpv4Denies
		hasChanges = true
	}

	// Handle IPv6 Permit
	oldIpv6PermitsByIndex := make(map[int64]verityPortAclIpv6PermitModel)
	for _, permit := range state.Ipv6Permit {
		if !permit.Index.IsNull() {
			oldIpv6PermitsByIndex[permit.Index.ValueInt64()] = permit
		}
	}

	var changedIpv6Permits []openapi.PortaclsPutRequestPortAclValueIpv6PermitInner
	ipv6PermitsChanged := false

	for _, permit := range plan.Ipv6Permit {
		if permit.Index.IsNull() {
			continue
		}

		index := permit.Index.ValueInt64()
		oldPermit, exists := oldIpv6PermitsByIndex[index]

		if !exists {
			// new permit, include all fields
			permitItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			// Handle ref_type validation
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				permit.Filter, permit.FilterRefType,
				"filter", "filter_ref_type_",
				true, false) {
				return
			}

			if !permit.Enable.IsNull() {
				permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			} else {
				permitItem.Enable = openapi.PtrBool(false)
			}

			if !permit.Filter.IsNull() {
				permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
			} else {
				permitItem.Filter = openapi.PtrString("")
			}

			if !permit.FilterRefType.IsNull() {
				permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
			} else {
				permitItem.FilterRefType = openapi.PtrString("")
			}

			changedIpv6Permits = append(changedIpv6Permits, permitItem)
			ipv6PermitsChanged = true
			continue
		}

		// existing permit, check which fields changed
		permitItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !permit.Enable.Equal(oldPermit.Enable) {
			permitItem.Enable = openapi.PtrBool(permit.Enable.ValueBool())
			fieldChanged = true
		}

		filterChanged := !permit.Filter.Equal(oldPermit.Filter)
		filterRefTypeChanged := !permit.FilterRefType.Equal(oldPermit.FilterRefType)

		if filterChanged || filterRefTypeChanged {
			// Validate using "one ref type supported" rules
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				permit.Filter, permit.FilterRefType,
				"filter", "filter_ref_type_",
				filterChanged, filterRefTypeChanged) {
				return
			}

			if filterChanged && !filterRefTypeChanged {
				// Just send the base field
				if !permit.Filter.IsNull() {
					permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
				} else {
					permitItem.Filter = openapi.PtrString("")
				}
				fieldChanged = true
			} else if filterRefTypeChanged {
				// Send both fields
				if !permit.Filter.IsNull() {
					permitItem.Filter = openapi.PtrString(permit.Filter.ValueString())
				} else {
					permitItem.Filter = openapi.PtrString("")
				}

				if !permit.FilterRefType.IsNull() {
					permitItem.FilterRefType = openapi.PtrString(permit.FilterRefType.ValueString())
				} else {
					permitItem.FilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		if fieldChanged {
			changedIpv6Permits = append(changedIpv6Permits, permitItem)
			ipv6PermitsChanged = true
		}
	}

	for idx := range oldIpv6PermitsByIndex {
		found := false
		for _, permit := range plan.Ipv6Permit {
			if !permit.Index.IsNull() && permit.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// permit removed - include only the index for deletion
			deletedPermit := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedIpv6Permits = append(changedIpv6Permits, deletedPermit)
			ipv6PermitsChanged = true
		}
	}

	if ipv6PermitsChanged && len(changedIpv6Permits) > 0 {
		portAclProps.Ipv6Permit = changedIpv6Permits
		hasChanges = true
	}

	// Handle IPv6 Deny
	oldIpv6DeniesByIndex := make(map[int64]verityPortAclIpv6DenyModel)
	for _, deny := range state.Ipv6Deny {
		if !deny.Index.IsNull() {
			oldIpv6DeniesByIndex[deny.Index.ValueInt64()] = deny
		}
	}

	var changedIpv6Denies []openapi.PortaclsPutRequestPortAclValueIpv6PermitInner
	ipv6DeniesChanged := false

	for _, deny := range plan.Ipv6Deny {
		if deny.Index.IsNull() {
			continue
		}

		index := deny.Index.ValueInt64()
		oldDeny, exists := oldIpv6DeniesByIndex[index]

		if !exists {
			// new deny, include all fields
			denyItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				deny.Filter, deny.FilterRefType,
				"filter", "filter_ref_type_",
				true, false) {
				return
			}

			if !deny.Enable.IsNull() {
				denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			} else {
				denyItem.Enable = openapi.PtrBool(false)
			}

			if !deny.Filter.IsNull() {
				denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
			} else {
				denyItem.Filter = openapi.PtrString("")
			}

			if !deny.FilterRefType.IsNull() {
				denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
			} else {
				denyItem.FilterRefType = openapi.PtrString("")
			}

			changedIpv6Denies = append(changedIpv6Denies, denyItem)
			ipv6DeniesChanged = true
			continue
		}

		// existing deny, check which fields changed
		denyItem := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !deny.Enable.Equal(oldDeny.Enable) {
			denyItem.Enable = openapi.PtrBool(deny.Enable.ValueBool())
			fieldChanged = true
		}

		filterChanged := !deny.Filter.Equal(oldDeny.Filter)
		filterRefTypeChanged := !deny.FilterRefType.Equal(oldDeny.FilterRefType)

		if filterChanged || filterRefTypeChanged {
			// Validate using "one ref type supported" rules
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				deny.Filter, deny.FilterRefType,
				"filter", "filter_ref_type_",
				filterChanged, filterRefTypeChanged) {
				return
			}

			if filterChanged && !filterRefTypeChanged {
				// Just send the base field
				if !deny.Filter.IsNull() {
					denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
				} else {
					denyItem.Filter = openapi.PtrString("")
				}
				fieldChanged = true
			} else if filterRefTypeChanged {
				// Send both fields
				if !deny.Filter.IsNull() {
					denyItem.Filter = openapi.PtrString(deny.Filter.ValueString())
				} else {
					denyItem.Filter = openapi.PtrString("")
				}

				if !deny.FilterRefType.IsNull() {
					denyItem.FilterRefType = openapi.PtrString(deny.FilterRefType.ValueString())
				} else {
					denyItem.FilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		if fieldChanged {
			changedIpv6Denies = append(changedIpv6Denies, denyItem)
			ipv6DeniesChanged = true
		}
	}

	for idx := range oldIpv6DeniesByIndex {
		found := false
		for _, deny := range plan.Ipv6Deny {
			if !deny.Index.IsNull() && deny.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// deny removed - include only the index for deletion
			deletedDeny := openapi.PortaclsPutRequestPortAclValueIpv6PermitInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedIpv6Denies = append(changedIpv6Denies, deletedDeny)
			ipv6DeniesChanged = true
		}
	}

	if ipv6DeniesChanged && len(changedIpv6Denies) > 0 {
		portAclProps.Ipv6Deny = changedIpv6Denies
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "port_acl", name, portAclProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for port ACL update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Port ACL %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Port ACL %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityPortAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPortAclResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "port_acl", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for port ACL deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Port ACL %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Port ACL %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "port_acls")
	resp.State.RemoveResource(ctx)
}

func (r *verityPortAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

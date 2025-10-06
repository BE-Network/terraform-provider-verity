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

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityGatewayProfileResource{}
	_ resource.ResourceWithConfigure   = &verityGatewayProfileResource{}
	_ resource.ResourceWithImportState = &verityGatewayProfileResource{}
)

func NewVerityGatewayProfileResource() resource.Resource {
	return &verityGatewayProfileResource{}
}

type verityGatewayProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityGatewayProfileResourceModel struct {
	Name             types.String                                `tfsdk:"name"`
	Enable           types.Bool                                  `tfsdk:"enable"`
	ObjectProperties []verityGatewayProfileObjectPropertiesModel `tfsdk:"object_properties"`
	ExternalGateways []verityGatewayProfileExternalGatewaysModel `tfsdk:"external_gateways"`
}

type verityGatewayProfileObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

type verityGatewayProfileExternalGatewaysModel struct {
	Enable         types.Bool   `tfsdk:"enable"`
	Gateway        types.String `tfsdk:"gateway"`
	GatewayRefType types.String `tfsdk:"gateway_ref_type_"`
	SourceIpMask   types.String `tfsdk:"source_ip_mask"`
	PeerGw         types.Bool   `tfsdk:"peer_gw"`
	Index          types.Int64  `tfsdk:"index"`
}

func (eg verityGatewayProfileExternalGatewaysModel) GetIndex() types.Int64 {
	return eg.Index
}

func (r *verityGatewayProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_profile"
}

func (r *verityGatewayProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityGatewayProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Gateway Profile",
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
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the gateway profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
					},
				},
			},
			"external_gateways": schema.ListNestedBlock{
				Description: "List of external gateway configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "BGP Gateway referenced for port profile",
							Optional:    true,
						},
						"gateway_ref_type_": schema.StringAttribute{
							Description: "Object type for gateway field",
							Optional:    true,
						},
						"source_ip_mask": schema.StringAttribute{
							Description: "Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session",
							Optional:    true,
						},
						"peer_gw": schema.BoolAttribute{
							Description: "Setting for paired switches only. Flag indicating that this gateway is a peer gateway.",
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

func (r *verityGatewayProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityGatewayProfileResourceModel
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
	profileProps := &openapi.GatewayprofilesPutRequestGatewayProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &profileProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.GatewayprofilesPutRequestGatewayProfileValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		profileProps.ObjectProperties = &objProps
	}

	// Handle external gateways
	if len(plan.ExternalGateways) > 0 {
		gateways := make([]openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner, len(plan.ExternalGateways))
		for i, gateway := range plan.ExternalGateways {
			gwItem := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{}
			if !gateway.Enable.IsNull() {
				gwItem.Enable = openapi.PtrBool(gateway.Enable.ValueBool())
			}
			if !gateway.Gateway.IsNull() {
				gwItem.Gateway = openapi.PtrString(gateway.Gateway.ValueString())
			}
			if !gateway.GatewayRefType.IsNull() {
				gwItem.GatewayRefType = openapi.PtrString(gateway.GatewayRefType.ValueString())
			}
			if !gateway.SourceIpMask.IsNull() {
				gwItem.SourceIpMask = openapi.PtrString(gateway.SourceIpMask.ValueString())
			}
			if !gateway.PeerGw.IsNull() {
				gwItem.PeerGw = openapi.PtrBool(gateway.PeerGw.ValueBool())
			}
			if !gateway.Index.IsNull() {
				gwItem.Index = openapi.PtrInt32(int32(gateway.Index.ValueInt64()))
			}
			gateways[i] = gwItem
		}
		profileProps.ExternalGateways = gateways
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "gateway_profile", name, *profileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateway_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityGatewayProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityGatewayProfileResourceModel
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

	profileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("gateway_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping gateway profile %s verification – trusting recent successful API operation", profileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching gateway profiles for verification of %s", profileName))

	type GatewayProfileResponse struct {
		GatewayProfile map[string]interface{} `json:"gateway_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "gateway_profiles", profileName,
		func() (GatewayProfileResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch gateway profiles")
			respAPI, err := r.client.GatewayProfilesAPI.GatewayprofilesGet(ctx).Execute()
			if err != nil {
				return GatewayProfileResponse{}, fmt.Errorf("error reading gateway profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res GatewayProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return GatewayProfileResponse{}, fmt.Errorf("failed to decode gateway profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d gateway profiles", len(res.GatewayProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Gateway Profile %s", profileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for gateway profile with name: %s", profileName))

	profileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.GatewayProfile,
		profileName,
		func(data interface{}) (string, bool) {
			if profile, ok := data.(map[string]interface{}); ok {
				if name, ok := profile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Gateway Profile with name '%s' not found in API response", profileName))
		resp.State.RemoveResource(ctx)
		return
	}

	profileMap, ok := profileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Gateway Profile Data",
			fmt.Sprintf("Gateway Profile data is not in expected format for %s", profileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found gateway profile '%s' under API key '%s'", profileName, actualAPIName))

	state.Name = utils.MapStringFromAPI(profileMap["name"])

	// Handle object properties
	if objProps, ok := profileMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityGatewayProfileObjectPropertiesModel{
			{Group: utils.MapStringFromAPI(objProps["group"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(profileMap[apiKey])
	}

	// Handle external gateways
	if ext, ok := profileMap["external_gateways"].([]interface{}); ok && len(ext) > 0 {
		var egList []verityGatewayProfileExternalGatewaysModel

		for _, item := range ext {
			gateway, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			egModel := verityGatewayProfileExternalGatewaysModel{
				Enable:         utils.MapBoolFromAPI(gateway["enable"]),
				Gateway:        utils.MapStringFromAPI(gateway["gateway"]),
				GatewayRefType: utils.MapStringFromAPI(gateway["gateway_ref_type_"]),
				SourceIpMask:   utils.MapStringFromAPI(gateway["source_ip_mask"]),
				PeerGw:         utils.MapBoolFromAPI(gateway["peer_gw"]),
				Index:          utils.MapInt64FromAPI(gateway["index"]),
			}

			egList = append(egList, egModel)
		}

		state.ExternalGateways = egList
	} else {
		state.ExternalGateways = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityGatewayProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityGatewayProfileResourceModel

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
	profileProps := openapi.GatewayprofilesPutRequestGatewayProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { profileProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { profileProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.GatewayprofilesPutRequestGatewayProfileValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			profileProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle external gateways
	externalGatewaysHandler := utils.IndexedItemHandler[verityGatewayProfileExternalGatewaysModel, openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner]{
		CreateNew: func(planItem verityGatewayProfileExternalGatewaysModel) openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner {
			gateway := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			if !planItem.Enable.IsNull() {
				gateway.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
			} else {
				gateway.Enable = openapi.PtrBool(false)
			}

			if !planItem.Gateway.IsNull() {
				gateway.Gateway = openapi.PtrString(planItem.Gateway.ValueString())
			} else {
				gateway.Gateway = openapi.PtrString("")
			}

			if !planItem.GatewayRefType.IsNull() {
				gateway.GatewayRefType = openapi.PtrString(planItem.GatewayRefType.ValueString())
			} else {
				gateway.GatewayRefType = openapi.PtrString("")
			}

			if !planItem.SourceIpMask.IsNull() {
				gateway.SourceIpMask = openapi.PtrString(planItem.SourceIpMask.ValueString())
			} else {
				gateway.SourceIpMask = openapi.PtrString("")
			}

			if !planItem.PeerGw.IsNull() {
				gateway.PeerGw = openapi.PtrBool(planItem.PeerGw.ValueBool())
			} else {
				gateway.PeerGw = openapi.PtrBool(false)
			}

			return gateway
		},
		UpdateExisting: func(planItem verityGatewayProfileExternalGatewaysModel, stateItem verityGatewayProfileExternalGatewaysModel) (openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner, bool) {
			gateway := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			fieldChanged := false

			if !planItem.Enable.Equal(stateItem.Enable) {
				gateway.Enable = openapi.PtrBool(planItem.Enable.ValueBool())
				fieldChanged = true
			}

			if !planItem.Gateway.Equal(stateItem.Gateway) {
				if !planItem.Gateway.IsNull() {
					gateway.Gateway = openapi.PtrString(planItem.Gateway.ValueString())
				} else {
					gateway.Gateway = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.GatewayRefType.Equal(stateItem.GatewayRefType) {
				if !planItem.GatewayRefType.IsNull() {
					gateway.GatewayRefType = openapi.PtrString(planItem.GatewayRefType.ValueString())
				} else {
					gateway.GatewayRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.SourceIpMask.Equal(stateItem.SourceIpMask) {
				if !planItem.SourceIpMask.IsNull() {
					gateway.SourceIpMask = openapi.PtrString(planItem.SourceIpMask.ValueString())
				} else {
					gateway.SourceIpMask = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.PeerGw.Equal(stateItem.PeerGw) {
				gateway.PeerGw = openapi.PtrBool(planItem.PeerGw.ValueBool())
				fieldChanged = true
			}

			return gateway, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner {
			return openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedExternalGateways, externalGatewaysChanged := utils.ProcessIndexedArrayUpdates(plan.ExternalGateways, state.ExternalGateways, externalGatewaysHandler)
	if externalGatewaysChanged {
		profileProps.ExternalGateways = changedExternalGateways
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "gateway_profile", name, profileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateway_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityGatewayProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityGatewayProfileResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "gateway_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateway_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityGatewayProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

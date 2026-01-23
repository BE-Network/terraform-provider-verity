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
	_ resource.Resource                = &verityGatewayProfileResource{}
	_ resource.ResourceWithConfigure   = &verityGatewayProfileResource{}
	_ resource.ResourceWithImportState = &verityGatewayProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityGatewayProfileResource{}
)

const gatewayProfileResourceType = "gatewayprofiles"

func NewVerityGatewayProfileResource() resource.Resource {
	return &verityGatewayProfileResource{}
}

type verityGatewayProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
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
							Computed:    true,
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
							Computed:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "BGP Gateway referenced for port profile",
							Optional:    true,
							Computed:    true,
						},
						"gateway_ref_type_": schema.StringAttribute{
							Description: "Object type for gateway field",
							Optional:    true,
							Computed:    true,
						},
						"source_ip_mask": schema.StringAttribute{
							Description: "Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session",
							Optional:    true,
							Computed:    true,
						},
						"peer_gw": schema.BoolAttribute{
							Description: "Setting for paired switches only. Flag indicating that this gateway is a peer gateway.",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
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
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		profileProps.ObjectProperties = &objProps
	}

	// Handle external gateways
	if len(plan.ExternalGateways) > 0 {
		gateways := make([]openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner, len(plan.ExternalGateways))
		for i, item := range plan.ExternalGateways {
			gwItem := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &gwItem.Enable, TFValue: item.Enable},
				{FieldName: "PeerGw", APIField: &gwItem.PeerGw, TFValue: item.PeerGw},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Gateway", APIField: &gwItem.Gateway, TFValue: item.Gateway},
				{FieldName: "GatewayRefType", APIField: &gwItem.GatewayRefType, TFValue: item.GatewayRefType},
				{FieldName: "SourceIpMask", APIField: &gwItem.SourceIpMask, TFValue: item.SourceIpMask},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &gwItem.Index, TFValue: item.Index},
			})
			gateways[i] = gwItem
		}
		profileProps.ExternalGateways = gateways
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "gateway_profile", name, *profileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateway_profiles")

	var minState verityGatewayProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if gatewayProfileData, exists := bulkMgr.GetResourceResponse("gateway_profile", name); exists {
			state := populateGatewayProfileState(ctx, minState, gatewayProfileData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if gatewayProfileData, exists := r.bulkOpsMgr.GetResourceResponse("gateway_profile", profileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached gateway profile data for %s from recent operation", profileName))
			state = populateGatewayProfileState(ctx, state, gatewayProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateGatewayProfileState(ctx, state, profileMap, r.provCtx.mode)
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
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
		}, &objPropsChanged)

		if objPropsChanged {
			profileProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle external gateways
	externalGatewaysHandler := utils.IndexedItemHandler[verityGatewayProfileExternalGatewaysModel, openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner]{
		CreateNew: func(planItem verityGatewayProfileExternalGatewaysModel) openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner {
			gateway := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &gateway.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &gateway.Enable, TFValue: planItem.Enable},
				{FieldName: "PeerGw", APIField: &gateway.PeerGw, TFValue: planItem.PeerGw},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Gateway", APIField: &gateway.Gateway, TFValue: planItem.Gateway},
				{FieldName: "GatewayRefType", APIField: &gateway.GatewayRefType, TFValue: planItem.GatewayRefType},
				{FieldName: "SourceIpMask", APIField: &gateway.SourceIpMask, TFValue: planItem.SourceIpMask},
			})

			return gateway
		},
		UpdateExisting: func(planItem verityGatewayProfileExternalGatewaysModel, stateItem verityGatewayProfileExternalGatewaysModel) (openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner, bool) {
			gateway := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &gateway.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { gateway.Enable = v }, &fieldChanged)

			// Handle gateway and gateway_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.Gateway, stateItem.Gateway, planItem.GatewayRefType, stateItem.GatewayRefType,
				func(v *string) { gateway.Gateway = v },
				func(v *string) { gateway.GatewayRefType = v },
				"gateway", "gateway_ref_type_",
				&fieldChanged, &resp.Diagnostics,
			) {
				return gateway, false
			}

			// Handle non-ref-type string fields
			utils.CompareAndSetStringField(planItem.SourceIpMask, stateItem.SourceIpMask, func(v *string) { gateway.SourceIpMask = v }, &fieldChanged)

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.PeerGw, stateItem.PeerGw, func(v *bool) { gateway.PeerGw = v }, &fieldChanged)

			return gateway, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner {
			gateway := openapi.GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &gateway.Index, TFValue: types.Int64Value(index)},
			})
			return gateway
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "gateway_profile", name, profileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateway_profiles")

	var minState verityGatewayProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if gatewayProfileData, exists := bulkMgr.GetResourceResponse("gateway_profile", name); exists {
			newState := populateGatewayProfileState(ctx, minState, gatewayProfileData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "gateway_profile", name, nil, &resp.Diagnostics)
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

// populateGatewayProfileState populates the state from API response data with mode-aware field mapping
func populateGatewayProfileState(ctx context.Context, state verityGatewayProfileResourceModel, data map[string]interface{}, mode string) verityGatewayProfileResourceModel {
	const resourceType = gatewayProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle object properties
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			state.ObjectProperties = []verityGatewayProfileObjectPropertiesModel{
				{Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode)},
			}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle external gateways
	if utils.FieldAppliesToMode(resourceType, "external_gateways", mode) {
		if ext, ok := data["external_gateways"].([]interface{}); ok && len(ext) > 0 {
			var egList []verityGatewayProfileExternalGatewaysModel

			for _, item := range ext {
				gateway, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				egModel := verityGatewayProfileExternalGatewaysModel{
					Enable:         utils.MapBoolWithModeNested(gateway, "enable", resourceType, "external_gateways.enable", mode),
					Gateway:        utils.MapStringWithModeNested(gateway, "gateway", resourceType, "external_gateways.gateway", mode),
					GatewayRefType: utils.MapStringWithModeNested(gateway, "gateway_ref_type_", resourceType, "external_gateways.gateway_ref_type_", mode),
					SourceIpMask:   utils.MapStringWithModeNested(gateway, "source_ip_mask", resourceType, "external_gateways.source_ip_mask", mode),
					PeerGw:         utils.MapBoolWithModeNested(gateway, "peer_gw", resourceType, "external_gateways.peer_gw", mode),
					Index:          utils.MapInt64WithModeNested(gateway, "index", resourceType, "external_gateways.index", mode),
				}

				egList = append(egList, egModel)
			}

			state.ExternalGateways = egList
		} else {
			state.ExternalGateways = nil
		}
	} else {
		state.ExternalGateways = nil
	}

	return state
}

func (r *verityGatewayProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityGatewayProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = gatewayProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)
}

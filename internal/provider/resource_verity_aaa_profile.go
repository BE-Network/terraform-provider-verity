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
	_ resource.Resource                = &verityAaaProfileResource{}
	_ resource.ResourceWithConfigure   = &verityAaaProfileResource{}
	_ resource.ResourceWithImportState = &verityAaaProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityAaaProfileResource{}
)

const aaaProfileResourceType = "deviceaaaprofiles"
const aaaProfileTerraformType = "verity_aaa_profile"

func NewVerityAaaProfileResource() resource.Resource {
	return &verityAaaProfileResource{}
}

type verityAaaProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityAaaProfileResourceModel struct {
	Name                 types.String                        `tfsdk:"name"`
	Enable               types.Bool                          `tfsdk:"enable"`
	FailThrough          types.Bool                          `tfsdk:"fail-through"`
	TacacsProfile        types.String                        `tfsdk:"tacacs_profile"`
	TacacsProfileRefType types.String                        `tfsdk:"tacacs_profile_ref_type_"`
	LoginDefault         []verityAaaProfileLoginDefaultModel `tfsdk:"login_default"`
}

type verityAaaProfileLoginDefaultModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	LoginMethod types.String `tfsdk:"login_method"`
	Index       types.Int64  `tfsdk:"index"`
}

func (m verityAaaProfileLoginDefaultModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityAaaProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_profile"
}

func (r *verityAaaProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityAaaProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Device AAA Profile.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Template Name. Must be unique within type.",
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
			"fail-through": schema.BoolAttribute{
				Description: "When enabled, authentication continues to access each server in the method list if an authentication request fails on one server",
				Optional:    true,
				Computed:    true,
			},
			"tacacs_profile": schema.StringAttribute{
				Description: "TACACS+ profile for authentication",
				Optional:    true,
				Computed:    true,
			},
			"tacacs_profile_ref_type_": schema.StringAttribute{
				Description: "Object type for tacacs_profile field",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"login_default": schema.ListNestedBlock{
				Description: "Authentication method list for remote access",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable this login method",
							Optional:    true,
							Computed:    true,
						},
						"login_method": schema.StringAttribute{
							Description: "Authentication method for remote access (SSH, etc.)",
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

func (r *verityAaaProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityAaaProfileResourceModel
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
	aaaProfileProps := &openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "TacacsProfile", APIField: &aaaProfileProps.TacacsProfile, TFValue: plan.TacacsProfile},
		{FieldName: "TacacsProfileRefType", APIField: &aaaProfileProps.TacacsProfileRefType, TFValue: plan.TacacsProfileRefType},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &aaaProfileProps.Enable, TFValue: plan.Enable},
		{FieldName: "FailThrough", APIField: &aaaProfileProps.FailThrough, TFValue: plan.FailThrough},
	})

	if len(plan.LoginDefault) > 0 {
		loginDefault := make([]openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner, len(plan.LoginDefault))

		for i, item := range plan.LoginDefault {
			loginMethod := openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &loginMethod.Enabled, TFValue: item.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "LoginMethod", APIField: &loginMethod.LoginMethod, TFValue: item.LoginMethod},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &loginMethod.Index, TFValue: item.Index},
			})

			loginDefault[i] = loginMethod
		}

		aaaProfileProps.LoginDefault = loginDefault
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_aaa_profile", name, *aaaProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AAA Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_aaa_profiles")

	var minState verityAaaProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aaaProfileData, exists := bulkMgr.GetResourceResponse("device_aaa_profile", name); exists {
			state := populateAaaProfileState(ctx, minState, aaaProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityAaaProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityAaaProfileResourceModel
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

	aaaProfileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if aaaProfileData, exists := r.bulkOpsMgr.GetResourceResponse("device_aaa_profile", aaaProfileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached AAA profile data for %s from recent operation", aaaProfileName))
			state = populateAaaProfileState(ctx, state, aaaProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("device_aaa_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping AAA profile %s verification – trusting recent successful API operation", aaaProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching AAA profiles for verification of %s", aaaProfileName))

	type AaaProfilesResponse struct {
		DeviceAaaProfile map[string]interface{} `json:"device_aaa_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "device_aaa_profiles", aaaProfileName,
		func() (AaaProfilesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch AAA profiles")
			respAPI, err := r.client.DeviceAAAProfilesAPI.DeviceaaaprofilesGet(ctx).Execute()
			if err != nil {
				return AaaProfilesResponse{}, fmt.Errorf("error reading AAA profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res AaaProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return AaaProfilesResponse{}, fmt.Errorf("failed to decode AAA profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d AAA profiles", len(res.DeviceAaaProfile)))
			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read AAA Profile %s", aaaProfileName))...,
		)
		return
	}

	aaaProfileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.DeviceAaaProfile,
		aaaProfileName,
		func(data interface{}) (string, bool) {
			if aaaProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := aaaProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)
	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("AAA Profile with name '%s' not found in API response", aaaProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	aaaProfileMap, ok := aaaProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid AAA Profile Data",
			fmt.Sprintf("AAA Profile data is not in expected format for %s", aaaProfileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found AAA profile '%s' under API key '%s'", aaaProfileName, actualAPIName))

	state = populateAaaProfileState(ctx, state, aaaProfileMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityAaaProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityAaaProfileResourceModel

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
	aaaProfileProps := openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValue{}
	hasChanges := false

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { aaaProfileProps.Name = v }, &hasChanges)

	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { aaaProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FailThrough, state.FailThrough, func(v *bool) { aaaProfileProps.FailThrough = v }, &hasChanges)

	if !utils.HandleOneRefTypeSupported(
		plan.TacacsProfile, state.TacacsProfile, plan.TacacsProfileRefType, state.TacacsProfileRefType,
		func(v *string) { aaaProfileProps.TacacsProfile = v },
		func(v *string) { aaaProfileProps.TacacsProfileRefType = v },
		"tacacs_profile", "tacacs_profile_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	loginDefaultHandler := utils.IndexedItemHandler[verityAaaProfileLoginDefaultModel, openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner]{
		CreateNew: func(planItem verityAaaProfileLoginDefaultModel) openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner {
			loginMethod := openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &loginMethod.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &loginMethod.Enabled, TFValue: planItem.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "LoginMethod", APIField: &loginMethod.LoginMethod, TFValue: planItem.LoginMethod},
			})

			return loginMethod
		},
		UpdateExisting: func(planItem verityAaaProfileLoginDefaultModel, stateItem verityAaaProfileLoginDefaultModel) (openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner, bool) {
			loginMethod := openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner{}
			fieldChanged := false

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &loginMethod.Index, TFValue: planItem.Index},
			})

			utils.CompareAndSetBoolField(planItem.Enabled, stateItem.Enabled, func(v *bool) { loginMethod.Enabled = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.LoginMethod, stateItem.LoginMethod, func(v *string) { loginMethod.LoginMethod = v }, &fieldChanged)

			return loginMethod, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner {
			loginMethod := openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &loginMethod.Index, TFValue: types.Int64Value(index)},
			})
			return loginMethod
		},
	}

	changedLoginDefault, loginDefaultChanged := utils.ProcessIndexedArrayUpdates(plan.LoginDefault, state.LoginDefault, loginDefaultHandler)
	if loginDefaultChanged {
		aaaProfileProps.LoginDefault = changedLoginDefault
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_aaa_profile", name, aaaProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AAA Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_aaa_profiles")

	var minState verityAaaProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aaaProfileData, exists := bulkMgr.GetResourceResponse("device_aaa_profile", name); exists {
			newState := populateAaaProfileState(ctx, minState, aaaProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityAaaProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityAaaProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_aaa_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("AAA Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_aaa_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityAaaProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateAaaProfileState(ctx context.Context, state verityAaaProfileResourceModel, data map[string]interface{}, mode string) verityAaaProfileResourceModel {
	const resourceType = aaaProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.FailThrough = utils.MapBoolWithMode(data, "fail-through", resourceType, mode)
	state.TacacsProfile = utils.MapStringWithMode(data, "tacacs_profile", resourceType, mode)
	state.TacacsProfileRefType = utils.MapStringWithMode(data, "tacacs_profile_ref_type_", resourceType, mode)

	if utils.FieldAppliesToMode(resourceType, "login_default", mode) {
		if loginDefaultData, ok := data["login_default"].([]interface{}); ok && len(loginDefaultData) > 0 {
			var loginDefault []verityAaaProfileLoginDefaultModel

			for _, item := range loginDefaultData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				loginDefault = append(loginDefault, verityAaaProfileLoginDefaultModel{
					Enabled:     utils.MapBoolWithModeNested(itemMap, "enabled", resourceType, "login_default.enabled", mode),
					LoginMethod: utils.MapStringWithModeNested(itemMap, "login_method", resourceType, "login_default.login_method", mode),
					Index:       utils.MapInt64WithModeNested(itemMap, "index", resourceType, "login_default.index", mode),
				})
			}

			state.LoginDefault = loginDefault
		} else {
			state.LoginDefault = nil
		}
	} else {
		state.LoginDefault = nil
	}

	return state
}

func (r *verityAaaProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityAaaProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = aaaProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"tacacs_profile", "tacacs_profile_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "fail-through",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "login_default",
		ItemCount:    len(plan.LoginDefault),
		StringFields: []string{"login_method"},
		BoolFields:   []string{"enabled"},
		Int64Fields:  []string{"index"},
	})
}

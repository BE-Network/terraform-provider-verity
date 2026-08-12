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
	_ resource.Resource                = &verityTacacsProfileResource{}
	_ resource.ResourceWithConfigure   = &verityTacacsProfileResource{}
	_ resource.ResourceWithImportState = &verityTacacsProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityTacacsProfileResource{}
)

const tacacsProfileResourceType = "tacacsprofiles"
const tacacsProfileTerraformType = "verity_tacacs_profile"

func NewVerityTacacsProfileResource() resource.Resource {
	return &verityTacacsProfileResource{}
}

type verityTacacsProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityTacacsProfileResourceModel struct {
	Name          types.String                           `tfsdk:"name"`
	Enable        types.Bool                             `tfsdk:"enable"`
	TacacsServers []verityTacacsProfileTacacsServerModel `tfsdk:"tacacs_servers"`
}

type verityTacacsProfileTacacsServerModel struct {
	Enabled   types.Bool   `tfsdk:"enabled"`
	Server    types.String `tfsdk:"server"`
	AuthType  types.String `tfsdk:"auth_type"`
	Port      types.String `tfsdk:"port"`
	Timeout   types.Int64  `tfsdk:"timeout"`
	Secret    types.String `tfsdk:"secret"`
	EncSecret types.String `tfsdk:"enc_secret"`
	Index     types.Int64  `tfsdk:"index"`
}

func (m verityTacacsProfileTacacsServerModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityTacacsProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tacacs_profile"
}

func (r *verityTacacsProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityTacacsProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity TACACS Profile.",
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
		},
		Blocks: map[string]schema.Block{
			"tacacs_servers": schema.ListNestedBlock{
				Description: "List of TACACS+ servers",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable TACACS+ server",
							Optional:    true,
							Computed:    true,
						},
						"server": schema.StringAttribute{
							Description: "IPv4, IPv6, or DNS name for TACACS+ server",
							Optional:    true,
							Computed:    true,
						},
						"auth_type": schema.StringAttribute{
							Description: "TACACS+ authentication type",
							Optional:    true,
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: "TACACS+ server port",
							Optional:    true,
							Computed:    true,
						},
						"timeout": schema.Int64Attribute{
							Description: "TACACS+ server timeout in seconds",
							Optional:    true,
							Computed:    true,
						},
						"secret": schema.StringAttribute{
							Description: "TACACS+ shared secret",
							Optional:    true,
							Computed:    true,
						},
						"enc_secret": schema.StringAttribute{
							Description: "TACACS+ shared secret (encrypted)",
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

func (r *verityTacacsProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityTacacsProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityTacacsProfileResourceModel
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
	tacacsProfileProps := &openapi.TacacsprofilesPutRequestTacacsProfileValue{
		Name: openapi.PtrString(name),
	}

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &tacacsProfileProps.Enable, TFValue: plan.Enable},
	})

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, tacacsProfileTerraformType, name)

	if len(plan.TacacsServers) > 0 {
		tacacsServersConfigMap := utils.BuildIndexedConfigMap(config.TacacsServers)
		tacacsServers := make([]openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner, len(plan.TacacsServers))

		for i, item := range plan.TacacsServers {
			server := openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &server.Enabled, TFValue: item.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Server", APIField: &server.Server, TFValue: item.Server},
				{FieldName: "AuthType", APIField: &server.AuthType, TFValue: item.AuthType},
				{FieldName: "Port", APIField: &server.Port, TFValue: item.Port},
				{FieldName: "Secret", APIField: &server.Secret, TFValue: item.Secret},
				{FieldName: "EncSecret", APIField: &server.EncSecret, TFValue: item.EncSecret},
			})

			configItem, cfg := utils.GetIndexedBlockConfig(item, tacacsServersConfigMap, "tacacs_servers", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "Timeout", APIField: &server.Timeout, TFValue: configItem.Timeout, IsConfigured: cfg.IsFieldConfigured("timeout")},
			})

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: item.Index},
			})

			tacacsServers[i] = server
		}

		tacacsProfileProps.TacacsServers = tacacsServers
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "tacacs_profile", name, *tacacsProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("TACACS Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "tacacs_profiles")

	var minState verityTacacsProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if tacacsProfileData, exists := bulkMgr.GetResourceResponse("tacacs_profile", name); exists {
			state := populateTacacsProfileState(ctx, minState, utils.MergeMissingPlanScalars(tacacsProfileData, plan, tacacsProfileResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, tacacsProfileResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityTacacsProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityTacacsProfileResourceModel
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

	tacacsProfileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil {
		if tacacsProfileData, exists := r.bulkOpsMgr.GetResourceResponse("tacacs_profile", tacacsProfileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached TACACS profile data for %s from recent operation", tacacsProfileName))
			state = populateTacacsProfileState(ctx, state, utils.ApplyPostOperationFallback(ctx, tacacsProfileData), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("tacacs_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping TACACS profile %s verification - trusting recent successful API operation", tacacsProfileName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	type TacacsProfilesResponse struct {
		TacacsProfile map[string]interface{} `json:"tacacs_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "tacacs_profiles", tacacsProfileName,
		func() (TacacsProfilesResponse, error) {
			respAPI, err := r.client.TACACSProfilesAPI.TacacsprofilesGet(ctx).Execute()
			if err != nil {
				return TacacsProfilesResponse{}, fmt.Errorf("error reading TACACS profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res TacacsProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return TacacsProfilesResponse{}, fmt.Errorf("failed to decode TACACS profiles response: %v", err)
			}

			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read TACACS Profile %s", tacacsProfileName))...,
		)
		return
	}

	tacacsProfileData, _, exists := utils.FindResourceByAPIName(
		result.TacacsProfile,
		tacacsProfileName,
		func(data interface{}) (string, bool) {
			if tacacsProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := tacacsProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	tacacsProfileMap, ok := tacacsProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid TACACS Profile Data",
			fmt.Sprintf("TACACS Profile data is not in expected format for %s", tacacsProfileName),
		)
		return
	}

	state = populateTacacsProfileState(ctx, state, utils.ApplyPostOperationFallback(ctx, tacacsProfileMap), r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityTacacsProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityTacacsProfileResourceModel

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

	var config verityTacacsProfileResourceModel
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
	tacacsProfileProps := openapi.TacacsprofilesPutRequestTacacsProfileValue{}
	hasChanges := false

	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { tacacsProfileProps.Name = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { tacacsProfileProps.Enable = v }, &hasChanges)

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, tacacsProfileTerraformType, name)

	tacacsServersConfigMap := utils.BuildIndexedConfigMap(config.TacacsServers)

	tacacsServersHandler := utils.IndexedItemHandler[verityTacacsProfileTacacsServerModel, openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner]{
		CreateNew: func(planItem verityTacacsProfileTacacsServerModel) openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner {
			server := openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &server.Enabled, TFValue: planItem.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Server", APIField: &server.Server, TFValue: planItem.Server},
				{FieldName: "AuthType", APIField: &server.AuthType, TFValue: planItem.AuthType},
				{FieldName: "Port", APIField: &server.Port, TFValue: planItem.Port},
				{FieldName: "Secret", APIField: &server.Secret, TFValue: planItem.Secret},
				{FieldName: "EncSecret", APIField: &server.EncSecret, TFValue: planItem.EncSecret},
			})

			configItem, cfg := utils.GetIndexedBlockConfig(planItem, tacacsServersConfigMap, "tacacs_servers", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "Timeout", APIField: &server.Timeout, TFValue: configItem.Timeout, IsConfigured: cfg.IsFieldConfigured("timeout")},
			})

			return server
		},
		UpdateExisting: func(planItem verityTacacsProfileTacacsServerModel, stateItem verityTacacsProfileTacacsServerModel) (openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner, bool) {
			server := openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner{}
			fieldChanged := false

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: planItem.Index},
			})

			utils.CompareAndSetBoolField(planItem.Enabled, stateItem.Enabled, func(v *bool) { server.Enabled = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Server, stateItem.Server, func(v *string) { server.Server = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.AuthType, stateItem.AuthType, func(v *string) { server.AuthType = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Port, stateItem.Port, func(v *string) { server.Port = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.Secret, stateItem.Secret, func(v *string) { server.Secret = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.EncSecret, stateItem.EncSecret, func(v *string) { server.EncSecret = v }, &fieldChanged)

			configItem, cfg := utils.GetIndexedBlockConfig(planItem, tacacsServersConfigMap, "tacacs_servers", configuredAttrs)
			utils.CompareAndSetNullableInt64Field(configItem.Timeout, stateItem.Timeout, cfg.IsFieldConfigured("timeout"), func(v *openapi.NullableInt64) { server.Timeout = *v }, &fieldChanged)

			return server, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner {
			server := openapi.TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: types.Int64Value(index)},
			})
			return server
		},
	}

	changedTacacsServers, tacacsServersChanged := utils.ProcessIndexedArrayUpdates(plan.TacacsServers, state.TacacsServers, tacacsServersHandler)
	if tacacsServersChanged {
		hasChanges = true
		tacacsProfileProps.TacacsServers = changedTacacsServers
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "tacacs_profile", name, tacacsProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("TACACS Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "tacacs_profiles")

	var minState verityTacacsProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if tacacsProfileData, exists := bulkMgr.GetResourceResponse("tacacs_profile", name); exists {
			newState := populateTacacsProfileState(ctx, minState, utils.MergeMissingPlanScalars(tacacsProfileData, plan, tacacsProfileResourceType, r.provCtx.mode), r.provCtx.mode)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, tacacsProfileResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityTacacsProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityTacacsProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "tacacs_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("TACACS Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "tacacs_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityTacacsProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateTacacsProfileState(ctx context.Context, state verityTacacsProfileResourceModel, data map[string]interface{}, mode string) verityTacacsProfileResourceModel {
	const resourceType = tacacsProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	if utils.FieldAppliesToMode(resourceType, "tacacs_servers", mode) {
		if tacacsServersData, ok := data["tacacs_servers"].([]interface{}); ok && len(tacacsServersData) > 0 {
			var tacacsServers []verityTacacsProfileTacacsServerModel

			for _, item := range tacacsServersData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				tacacsServers = append(tacacsServers, verityTacacsProfileTacacsServerModel{
					Enabled:   utils.MapBoolWithModeNested(itemMap, "enabled", resourceType, "tacacs_servers.enabled", mode),
					Server:    utils.MapStringWithModeNested(itemMap, "server", resourceType, "tacacs_servers.server", mode),
					AuthType:  utils.MapStringWithModeNested(itemMap, "auth_type", resourceType, "tacacs_servers.auth_type", mode),
					Port:      utils.MapStringWithModeNested(itemMap, "port", resourceType, "tacacs_servers.port", mode),
					Timeout:   utils.MapInt64WithModeNested(itemMap, "timeout", resourceType, "tacacs_servers.timeout", mode),
					Secret:    utils.MapStringWithModeNested(itemMap, "secret", resourceType, "tacacs_servers.secret", mode),
					EncSecret: utils.MapStringWithModeNested(itemMap, "enc_secret", resourceType, "tacacs_servers.enc_secret", mode),
					Index:     utils.MapInt64WithModeNested(itemMap, "index", resourceType, "tacacs_servers.index", mode),
				})
			}

			state.TacacsServers = tacacsServers
		} else {
			state.TacacsServers = nil
		}
	} else {
		state.TacacsServers = nil
	}

	return state
}

func (r *verityTacacsProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityTacacsProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const resourceType = tacacsProfileResourceType
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

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "tacacs_servers",
		ItemCount:    len(plan.TacacsServers),
		StringFields: []string{"server", "auth_type", "port", "secret", "enc_secret"},
		BoolFields:   []string{"enabled"},
		Int64Fields:  []string{"timeout", "index"},
	})

	if req.State.Raw.IsNull() {
		return
	}

	var state verityTacacsProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityTacacsProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, tacacsProfileTerraformType, name)

	for i, configServer := range config.TacacsServers {
		serverIndex := configServer.Index.ValueInt64()
		var stateServer *verityTacacsProfileTacacsServerModel
		for j := range state.TacacsServers {
			if state.TacacsServers[j].Index.ValueInt64() == serverIndex {
				stateServer = &state.TacacsServers[j]
				break
			}
		}

		if stateServer != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "tacacs_servers",
				BlockListPath:   "tacacs_servers",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: serverIndex, AttrName: "timeout", ConfigVal: configServer.Timeout, StateVal: stateServer.Timeout},
				},
			})
		}
	}
}

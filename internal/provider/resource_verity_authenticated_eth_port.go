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
	_ resource.Resource                = &verityAuthenticatedEthPortResource{}
	_ resource.ResourceWithConfigure   = &verityAuthenticatedEthPortResource{}
	_ resource.ResourceWithImportState = &verityAuthenticatedEthPortResource{}
	_ resource.ResourceWithModifyPlan  = &verityAuthenticatedEthPortResource{}
)

const authenticatedEthPortResourceType = "authenticatedethports"

func NewVerityAuthenticatedEthPortResource() resource.Resource {
	return &verityAuthenticatedEthPortResource{}
}

type verityAuthenticatedEthPortResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityAuthenticatedEthPortResourceModel struct {
	Name                        types.String                                      `tfsdk:"name"`
	Enable                      types.Bool                                        `tfsdk:"enable"`
	ConnectionMode              types.String                                      `tfsdk:"connection_mode"`
	ReauthorizationPeriodSec    types.Int64                                       `tfsdk:"reauthorization_period_sec"`
	AllowMacBasedAuthentication types.Bool                                        `tfsdk:"allow_mac_based_authentication"`
	MacAuthenticationHoldoffSec types.Int64                                       `tfsdk:"mac_authentication_holdoff_sec"`
	TrustedPort                 types.Bool                                        `tfsdk:"trusted_port"`
	EthPorts                    []verityAuthenticatedEthPortEthPortModel          `tfsdk:"eth_ports"`
	ObjectProperties            []verityAuthenticatedEthPortObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityAuthenticatedEthPortEthPortModel struct {
	EthPortProfileNumEnable          types.Bool   `tfsdk:"eth_port_profile_num_enable"`
	EthPortProfileNumEthPort         types.String `tfsdk:"eth_port_profile_num_eth_port"`
	EthPortProfileNumEthPortRefType  types.String `tfsdk:"eth_port_profile_num_eth_port_ref_type_"`
	EthPortProfileNumWalledGardenSet types.Bool   `tfsdk:"eth_port_profile_num_walled_garden_set"`
	EthPortProfileNumRadiusFilterId  types.String `tfsdk:"eth_port_profile_num_radius_filter_id"`
	Index                            types.Int64  `tfsdk:"index"`
}

func (ep verityAuthenticatedEthPortEthPortModel) GetIndex() types.Int64 {
	return ep.Index
}

type verityAuthenticatedEthPortObjectPropertiesModel struct {
	Group          types.String `tfsdk:"group"`
	PortMonitoring types.String `tfsdk:"port_monitoring"`
}

func (r *verityAuthenticatedEthPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticated_eth_port"
}

func (r *verityAuthenticatedEthPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityAuthenticatedEthPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Authenticated Eth-Port",
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
			"connection_mode": schema.StringAttribute{
				Description: "Choose connection mode for Authenticated Eth-Port",
				Optional:    true,
				Computed:    true,
			},
			"reauthorization_period_sec": schema.Int64Attribute{
				Description: "Amount of time in seconds before 802.1X requires reauthorization of an active session. \"0\" disables reauthorization (not recommended)",
				Optional:    true,
				Computed:    true,
			},
			"allow_mac_based_authentication": schema.BoolAttribute{
				Description: "Enables 802.1x to capture the connected MAC address and send it to the Radius Server instead of requesting credentials. Useful for printers and similar devices",
				Optional:    true,
				Computed:    true,
			},
			"mac_authentication_holdoff_sec": schema.Int64Attribute{
				Description: "Amount of time in seconds 802.1X authentication is allowed to run before MAC-based authentication has begun",
				Optional:    true,
				Computed:    true,
			},
			"trusted_port": schema.BoolAttribute{
				Description: "Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"eth_ports": schema.ListNestedBlock{
				Description: "Ethernet port configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"eth_port_profile_num_enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_profile_num_eth_port": schema.StringAttribute{
							Description: "Choose an Eth Port Profile",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_profile_num_eth_port_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_profile_num_eth_port field",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_profile_num_walled_garden_set": schema.BoolAttribute{
							Description: "Flag indicating this Eth Port Profile is the Walled Garden",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_profile_num_radius_filter_id": schema.StringAttribute{
							Description: "The value of filter-id in the RADIUS response which will evoke this Eth Port Profile",
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
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the authenticated eth-port",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityAuthenticatedEthPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityAuthenticatedEthPortResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityAuthenticatedEthPortResourceModel
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
	aepProps := &openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "ConnectionMode", APIField: &aepProps.ConnectionMode, TFValue: plan.ConnectionMode},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &aepProps.Enable, TFValue: plan.Enable},
		{FieldName: "AllowMacBasedAuthentication", APIField: &aepProps.AllowMacBasedAuthentication, TFValue: plan.AllowMacBasedAuthentication},
		{FieldName: "TrustedPort", APIField: &aepProps.TrustedPort, TFValue: plan.TrustedPort},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_authenticated_eth_port", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "ReauthorizationPeriodSec", APIField: &aepProps.ReauthorizationPeriodSec, TFValue: config.ReauthorizationPeriodSec, IsConfigured: configuredAttrs.IsConfigured("reauthorization_period_sec")},
		{FieldName: "MacAuthenticationHoldoffSec", APIField: &aepProps.MacAuthenticationHoldoffSec, TFValue: config.MacAuthenticationHoldoffSec, IsConfigured: configuredAttrs.IsConfigured("mac_authentication_holdoff_sec")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", TFValue: op.PortMonitoring, APIValue: &objProps.PortMonitoring},
		})
		aepProps.ObjectProperties = &objProps
	}

	// Handle eth ports
	if len(plan.EthPorts) > 0 {
		ethPorts := make([]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner, len(plan.EthPorts))
		for i, item := range plan.EthPorts {
			ethPortItem := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "EthPortProfileNumEnable", APIField: &ethPortItem.EthPortProfileNumEnable, TFValue: item.EthPortProfileNumEnable},
				{FieldName: "EthPortProfileNumWalledGardenSet", APIField: &ethPortItem.EthPortProfileNumWalledGardenSet, TFValue: item.EthPortProfileNumWalledGardenSet},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "EthPortProfileNumEthPort", APIField: &ethPortItem.EthPortProfileNumEthPort, TFValue: item.EthPortProfileNumEthPort},
				{FieldName: "EthPortProfileNumEthPortRefType", APIField: &ethPortItem.EthPortProfileNumEthPortRefType, TFValue: item.EthPortProfileNumEthPortRefType},
				{FieldName: "EthPortProfileNumRadiusFilterId", APIField: &ethPortItem.EthPortProfileNumRadiusFilterId, TFValue: item.EthPortProfileNumRadiusFilterId},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethPortItem.Index, TFValue: item.Index},
			})
			ethPorts[i] = ethPortItem
		}
		aepProps.EthPorts = ethPorts
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "authenticated_eth_port", name, *aepProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")

	var minState verityAuthenticatedEthPortResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aepData, exists := bulkMgr.GetResourceResponse("authenticated_eth_port", name); exists {
			state := populateAuthenticatedEthPortState(ctx, minState, aepData, r.provCtx.mode)
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

func (r *verityAuthenticatedEthPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityAuthenticatedEthPortResourceModel
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

	aepName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if authenticatedEthPortData, exists := r.bulkOpsMgr.GetResourceResponse("authenticated_eth_port", aepName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached authenticated eth port data for %s from recent operation", aepName))
			state = populateAuthenticatedEthPortState(ctx, state, authenticatedEthPortData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("authenticated_eth_port") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Authenticated Eth-Port %s verification – trusting recent successful API operation", aepName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Authenticated Eth-Ports for verification of %s", aepName))

	type AuthenticatedEthPortResponse struct {
		AuthenticatedEthPort map[string]interface{} `json:"authenticated_eth_port"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "authenticated_eth_ports", aepName,
		func() (AuthenticatedEthPortResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Authenticated Eth-Ports")
			respAPI, err := r.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx).Execute()
			if err != nil {
				return AuthenticatedEthPortResponse{}, fmt.Errorf("error reading Authenticated Eth-Ports: %v", err)
			}
			defer respAPI.Body.Close()

			var res AuthenticatedEthPortResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return AuthenticatedEthPortResponse{}, fmt.Errorf("failed to decode Authenticated Eth-Ports response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Authenticated Eth-Ports", len(res.AuthenticatedEthPort)))
			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Authenticated Eth-Port %s", aepName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Authenticated Eth-Port with name: %s", aepName))

	aepData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.AuthenticatedEthPort,
		aepName,
		func(data interface{}) (string, bool) {
			if aethPort, ok := data.(map[string]interface{}); ok {
				if name, ok := aethPort["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Authenticated Eth-Port with name '%s' not found in API response", aepName))
		resp.State.RemoveResource(ctx)
		return
	}

	aepMap, ok := aepData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Authenticated Eth-Port Data",
			fmt.Sprintf("Authenticated Eth-Port data is not in expected format for %s", aepName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Authenticated Eth-Port '%s' under API key '%s'", aepName, actualAPIName))

	state = populateAuthenticatedEthPortState(ctx, state, aepMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityAuthenticatedEthPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityAuthenticatedEthPortResourceModel

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
	aepProps := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { aepProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ConnectionMode, state.ConnectionMode, func(v *string) { aepProps.ConnectionMode = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { aepProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AllowMacBasedAuthentication, state.AllowMacBasedAuthentication, func(v *bool) { aepProps.AllowMacBasedAuthentication = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TrustedPort, state.TrustedPort, func(v *bool) { aepProps.TrustedPort = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.ReauthorizationPeriodSec, state.ReauthorizationPeriodSec, func(v *openapi.NullableInt32) { aepProps.ReauthorizationPeriodSec = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MacAuthenticationHoldoffSec, state.MacAuthenticationHoldoffSec, func(v *openapi.NullableInt32) { aepProps.MacAuthenticationHoldoffSec = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", PlanValue: op.PortMonitoring, StateValue: st.PortMonitoring, APIValue: &objProps.PortMonitoring},
		}, &objPropsChanged)

		if objPropsChanged {
			aepProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle eth ports
	ethPortsHandler := utils.IndexedItemHandler[verityAuthenticatedEthPortEthPortModel, openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner]{
		CreateNew: func(planItem verityAuthenticatedEthPortEthPortModel) openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner {
			item := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "EthPortProfileNumEnable", APIField: &item.EthPortProfileNumEnable, TFValue: planItem.EthPortProfileNumEnable},
				{FieldName: "EthPortProfileNumWalledGardenSet", APIField: &item.EthPortProfileNumWalledGardenSet, TFValue: planItem.EthPortProfileNumWalledGardenSet},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "EthPortProfileNumEthPort", APIField: &item.EthPortProfileNumEthPort, TFValue: planItem.EthPortProfileNumEthPort},
				{FieldName: "EthPortProfileNumEthPortRefType", APIField: &item.EthPortProfileNumEthPortRefType, TFValue: planItem.EthPortProfileNumEthPortRefType},
				{FieldName: "EthPortProfileNumRadiusFilterId", APIField: &item.EthPortProfileNumRadiusFilterId, TFValue: planItem.EthPortProfileNumRadiusFilterId},
			})

			return item
		},
		UpdateExisting: func(planItem verityAuthenticatedEthPortEthPortModel, stateItem verityAuthenticatedEthPortEthPortModel) (openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner, bool) {
			item := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean field changes
			utils.CompareAndSetBoolField(planItem.EthPortProfileNumEnable, stateItem.EthPortProfileNumEnable, func(v *bool) { item.EthPortProfileNumEnable = v }, &fieldChanged)

			// Handle eth_port_profile_num_eth_port and eth_port_profile_num_eth_port_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.EthPortProfileNumEthPort, stateItem.EthPortProfileNumEthPort, planItem.EthPortProfileNumEthPortRefType, stateItem.EthPortProfileNumEthPortRefType,
				func(v *string) { item.EthPortProfileNumEthPort = v },
				func(v *string) { item.EthPortProfileNumEthPortRefType = v },
				"eth_port_profile_num_eth_port", "eth_port_profile_num_eth_port_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return item, false
			}

			utils.CompareAndSetBoolField(planItem.EthPortProfileNumWalledGardenSet, stateItem.EthPortProfileNumWalledGardenSet, func(v *bool) { item.EthPortProfileNumWalledGardenSet = v }, &fieldChanged)

			// Handle string field changes
			utils.CompareAndSetStringField(planItem.EthPortProfileNumRadiusFilterId, stateItem.EthPortProfileNumRadiusFilterId, func(v *string) { item.EthPortProfileNumRadiusFilterId = v }, &fieldChanged)

			return item, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner {
			item := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedEthPorts, ethPortsChanged := utils.ProcessIndexedArrayUpdates(plan.EthPorts, state.EthPorts, ethPortsHandler)
	if ethPortsChanged {
		aepProps.EthPorts = changedEthPorts
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "authenticated_eth_port", name, aepProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")

	var minState verityAuthenticatedEthPortResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if aepData, exists := bulkMgr.GetResourceResponse("authenticated_eth_port", name); exists {
			newState := populateAuthenticatedEthPortState(ctx, minState, aepData, r.provCtx.mode)
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

func (r *verityAuthenticatedEthPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityAuthenticatedEthPortResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "authenticated_eth_port", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")
	resp.State.RemoveResource(ctx)
}

func (r *verityAuthenticatedEthPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateAuthenticatedEthPortState(ctx context.Context, state verityAuthenticatedEthPortResourceModel, data map[string]interface{}, mode string) verityAuthenticatedEthPortResourceModel {
	const resourceType = authenticatedEthPortResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.AllowMacBasedAuthentication = utils.MapBoolWithMode(data, "allow_mac_based_authentication", resourceType, mode)
	state.TrustedPort = utils.MapBoolWithMode(data, "trusted_port", resourceType, mode)

	// String fields
	state.ConnectionMode = utils.MapStringWithMode(data, "connection_mode", resourceType, mode)

	// Int64 fields
	state.ReauthorizationPeriodSec = utils.MapInt64WithMode(data, "reauthorization_period_sec", resourceType, mode)
	state.MacAuthenticationHoldoffSec = utils.MapInt64WithMode(data, "mac_authentication_holdoff_sec", resourceType, mode)

	// Handle eth_ports array
	if utils.FieldAppliesToMode(resourceType, "eth_ports", mode) {
		if ethPortsData, ok := data["eth_ports"].([]interface{}); ok && len(ethPortsData) > 0 {
			var ethPorts []verityAuthenticatedEthPortEthPortModel
			for _, e := range ethPortsData {
				ethPort, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				ethPortModel := verityAuthenticatedEthPortEthPortModel{
					EthPortProfileNumEnable:          utils.MapBoolWithModeNested(ethPort, "eth_port_profile_num_enable", resourceType, "eth_ports.eth_port_profile_num_enable", mode),
					EthPortProfileNumEthPort:         utils.MapStringWithModeNested(ethPort, "eth_port_profile_num_eth_port", resourceType, "eth_ports.eth_port_profile_num_eth_port", mode),
					EthPortProfileNumEthPortRefType:  utils.MapStringWithModeNested(ethPort, "eth_port_profile_num_eth_port_ref_type_", resourceType, "eth_ports.eth_port_profile_num_eth_port_ref_type_", mode),
					EthPortProfileNumWalledGardenSet: utils.MapBoolWithModeNested(ethPort, "eth_port_profile_num_walled_garden_set", resourceType, "eth_ports.eth_port_profile_num_walled_garden_set", mode),
					EthPortProfileNumRadiusFilterId:  utils.MapStringWithModeNested(ethPort, "eth_port_profile_num_radius_filter_id", resourceType, "eth_ports.eth_port_profile_num_radius_filter_id", mode),
					Index:                            utils.MapInt64WithModeNested(ethPort, "index", resourceType, "eth_ports.index", mode),
				}
				ethPorts = append(ethPorts, ethPortModel)
			}
			state.EthPorts = ethPorts
		} else {
			state.EthPorts = nil
		}
	} else {
		state.EthPorts = nil
	}

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityAuthenticatedEthPortObjectPropertiesModel{
				Group:          utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
				PortMonitoring: utils.MapStringWithModeNested(objProps, "port_monitoring", resourceType, "object_properties.port_monitoring", mode),
			}
			state.ObjectProperties = []verityAuthenticatedEthPortObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityAuthenticatedEthPortResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityAuthenticatedEthPortResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = authenticatedEthPortResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"connection_mode",
	)

	nullifier.NullifyBools(
		"enable", "allow_mac_based_authentication", "trusted_port",
	)

	nullifier.NullifyInt64s(
		"reauthorization_period_sec", "mac_authentication_holdoff_sec",
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
	var state verityAuthenticatedEthPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityAuthenticatedEthPortResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_authenticated_eth_port", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "reauthorization_period_sec", ConfigVal: config.ReauthorizationPeriodSec, StateVal: state.ReauthorizationPeriodSec},
			{AttrName: "mac_authentication_holdoff_sec", ConfigVal: config.MacAuthenticationHoldoffSec, StateVal: state.MacAuthenticationHoldoffSec},
		},
	})
}

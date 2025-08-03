package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	_ resource.Resource                = &verityAuthenticatedEthPortResource{}
	_ resource.ResourceWithConfigure   = &verityAuthenticatedEthPortResource{}
	_ resource.ResourceWithImportState = &verityAuthenticatedEthPortResource{}
)

func NewVerityAuthenticatedEthPortResource() resource.Resource {
	return &verityAuthenticatedEthPortResource{}
}

type verityAuthenticatedEthPortResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
			},
			"connection_mode": schema.StringAttribute{
				Description: "Choose connection mode for Authenticated Eth-Port",
				Optional:    true,
			},
			"reauthorization_period_sec": schema.Int64Attribute{
				Description: "Amount of time in seconds before 802.1X requires reauthorization of an active session. \"0\" disables reauthorization (not recommended)",
				Optional:    true,
			},
			"allow_mac_based_authentication": schema.BoolAttribute{
				Description: "Enables 802.1x to capture the connected MAC address and send it to the Radius Server instead of requesting credentials. Useful for printers and similar devices",
				Optional:    true,
			},
			"mac_authentication_holdoff_sec": schema.Int64Attribute{
				Description: "Amount of time in seconds 802.1X authentication is allowed to run before MAC-based authentication has begun",
				Optional:    true,
			},
			"trusted_port": schema.BoolAttribute{
				Description: "Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.",
				Optional:    true,
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
						},
						"eth_port_profile_num_eth_port": schema.StringAttribute{
							Description: "Choose an Eth Port Profile",
							Optional:    true,
						},
						"eth_port_profile_num_eth_port_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_profile_num_eth_port field",
							Optional:    true,
						},
						"eth_port_profile_num_walled_garden_set": schema.BoolAttribute{
							Description: "Flag indicating this Eth Port Profile is the Walled Garden",
							Optional:    true,
						},
						"eth_port_profile_num_radius_filter_id": schema.StringAttribute{
							Description: "The value of filter-id in the RADIUS response which will evoke this Eth Port Profile",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
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
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
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

	if !plan.Enable.IsNull() {
		aepProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.ConnectionMode.IsNull() {
		aepProps.ConnectionMode = openapi.PtrString(plan.ConnectionMode.ValueString())
	}
	if !plan.ReauthorizationPeriodSec.IsNull() {
		aepProps.ReauthorizationPeriodSec = openapi.PtrInt32(int32(plan.ReauthorizationPeriodSec.ValueInt64()))
	}
	if !plan.AllowMacBasedAuthentication.IsNull() {
		aepProps.AllowMacBasedAuthentication = openapi.PtrBool(plan.AllowMacBasedAuthentication.ValueBool())
	}
	if !plan.MacAuthenticationHoldoffSec.IsNull() {
		aepProps.MacAuthenticationHoldoffSec = openapi.PtrInt32(int32(plan.MacAuthenticationHoldoffSec.ValueInt64()))
	}
	if !plan.TrustedPort.IsNull() {
		aepProps.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
	}

	if len(plan.EthPorts) > 0 {
		ethPorts := make([]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner, len(plan.EthPorts))
		for i, ethPort := range plan.EthPorts {
			ethPortItem := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{}
			if !ethPort.EthPortProfileNumEnable.IsNull() {
				ethPortItem.EthPortProfileNumEnable = openapi.PtrBool(ethPort.EthPortProfileNumEnable.ValueBool())
			}
			if !ethPort.EthPortProfileNumEthPort.IsNull() {
				ethPortItem.EthPortProfileNumEthPort = openapi.PtrString(ethPort.EthPortProfileNumEthPort.ValueString())
			}
			if !ethPort.EthPortProfileNumEthPortRefType.IsNull() {
				ethPortItem.EthPortProfileNumEthPortRefType = openapi.PtrString(ethPort.EthPortProfileNumEthPortRefType.ValueString())
			}
			if !ethPort.EthPortProfileNumWalledGardenSet.IsNull() {
				ethPortItem.EthPortProfileNumWalledGardenSet = openapi.PtrBool(ethPort.EthPortProfileNumWalledGardenSet.ValueBool())
			}
			if !ethPort.EthPortProfileNumRadiusFilterId.IsNull() {
				ethPortItem.EthPortProfileNumRadiusFilterId = openapi.PtrString(ethPort.EthPortProfileNumRadiusFilterId.ValueString())
			}
			if !ethPort.Index.IsNull() {
				ethPortItem.Index = openapi.PtrInt32(int32(ethPort.Index.ValueInt64()))
			}
			ethPorts[i] = ethPortItem
		}
		aepProps.EthPorts = ethPorts
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		if !op.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
		}
		aepProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "authenticated_eth_port", name, *aepProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for authenticated eth-port creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Authenticated Eth-Port %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("authenticated_eth_port") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Authenticated Eth-Port %s verification – trusting recent successful API operation", aepName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Authenticated Eth-Ports for verification of %s", aepName))

	type AuthenticatedEthPortResponse struct {
		AuthenticatedEthPort map[string]interface{} `json:"authenticated_eth_port"`
	}

	var result AuthenticatedEthPortResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		aepData, fetchErr := getCachedResponse(ctx, r.provCtx, "authenticated_eth_ports", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Authenticated Eth-Ports")
			respAPI, err := r.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Authenticated Eth-Ports: %v", err)
			}
			defer respAPI.Body.Close()

			var res AuthenticatedEthPortResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Authenticated Eth-Ports response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Authenticated Eth-Ports", len(res.AuthenticatedEthPort)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Authenticated Eth-Ports on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = aepData.(AuthenticatedEthPortResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Authenticated Eth-Port %s", aepName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Authenticated Eth-Port with ID: %s", aepName))
	var aepData map[string]interface{}
	exists := false

	if data, ok := result.AuthenticatedEthPort[aepName].(map[string]interface{}); ok {
		aepData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Authenticated Eth-Port directly by ID: %s", aepName))
	} else {
		for apiName, a := range result.AuthenticatedEthPort {
			aethPort, ok := a.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := aethPort["name"].(string); ok && name == aepName {
				aepData = aethPort
				aepName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Authenticated Eth-Port with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Authenticated Eth-Port with ID '%s' not found in API response", aepName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", aepData["name"]))

	if enable, ok := aepData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if allowMac, ok := aepData["allow_mac_based_authentication"].(bool); ok {
		state.AllowMacBasedAuthentication = types.BoolValue(allowMac)
	} else {
		state.AllowMacBasedAuthentication = types.BoolNull()
	}

	if trusted, ok := aepData["trusted_port"].(bool); ok {
		state.TrustedPort = types.BoolValue(trusted)
	} else {
		state.TrustedPort = types.BoolNull()
	}

	if connMode, ok := aepData["connection_mode"].(string); ok {
		state.ConnectionMode = types.StringValue(connMode)
	} else {
		state.ConnectionMode = types.StringNull()
	}

	if reauth, ok := aepData["reauthorization_period_sec"]; ok && reauth != nil {
		switch v := reauth.(type) {
		case int:
			state.ReauthorizationPeriodSec = types.Int64Value(int64(v))
		case int32:
			state.ReauthorizationPeriodSec = types.Int64Value(int64(v))
		case int64:
			state.ReauthorizationPeriodSec = types.Int64Value(v)
		case float64:
			state.ReauthorizationPeriodSec = types.Int64Value(int64(v))
		case string:
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				state.ReauthorizationPeriodSec = types.Int64Value(intVal)
			} else {
				state.ReauthorizationPeriodSec = types.Int64Null()
			}
		default:
			state.ReauthorizationPeriodSec = types.Int64Null()
		}
	} else {
		state.ReauthorizationPeriodSec = types.Int64Null()
	}

	if macHoldoff, ok := aepData["mac_authentication_holdoff_sec"]; ok && macHoldoff != nil {
		switch v := macHoldoff.(type) {
		case int:
			state.MacAuthenticationHoldoffSec = types.Int64Value(int64(v))
		case int32:
			state.MacAuthenticationHoldoffSec = types.Int64Value(int64(v))
		case int64:
			state.MacAuthenticationHoldoffSec = types.Int64Value(v)
		case float64:
			state.MacAuthenticationHoldoffSec = types.Int64Value(int64(v))
		case string:
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				state.MacAuthenticationHoldoffSec = types.Int64Value(intVal)
			} else {
				state.MacAuthenticationHoldoffSec = types.Int64Null()
			}
		default:
			state.MacAuthenticationHoldoffSec = types.Int64Null()
		}
	} else {
		state.MacAuthenticationHoldoffSec = types.Int64Null()
	}

	if ethPortsArray, ok := aepData["eth_ports"].([]interface{}); ok && len(ethPortsArray) > 0 {
		var ethPorts []verityAuthenticatedEthPortEthPortModel
		for _, e := range ethPortsArray {
			ethPort, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			ethPortModel := verityAuthenticatedEthPortEthPortModel{}

			if enable, ok := ethPort["eth_port_profile_num_enable"].(bool); ok {
				ethPortModel.EthPortProfileNumEnable = types.BoolValue(enable)
			} else {
				ethPortModel.EthPortProfileNumEnable = types.BoolNull()
			}

			if ethPortProfile, ok := ethPort["eth_port_profile_num_eth_port"].(string); ok {
				ethPortModel.EthPortProfileNumEthPort = types.StringValue(ethPortProfile)
			} else {
				ethPortModel.EthPortProfileNumEthPort = types.StringNull()
			}

			if refType, ok := ethPort["eth_port_profile_num_eth_port_ref_type_"].(string); ok {
				ethPortModel.EthPortProfileNumEthPortRefType = types.StringValue(refType)
			} else {
				ethPortModel.EthPortProfileNumEthPortRefType = types.StringNull()
			}

			if walledGarden, ok := ethPort["eth_port_profile_num_walled_garden_set"].(bool); ok {
				ethPortModel.EthPortProfileNumWalledGardenSet = types.BoolValue(walledGarden)
			} else {
				ethPortModel.EthPortProfileNumWalledGardenSet = types.BoolNull()
			}

			if radiusFilter, ok := ethPort["eth_port_profile_num_radius_filter_id"].(string); ok {
				ethPortModel.EthPortProfileNumRadiusFilterId = types.StringValue(radiusFilter)
			} else {
				ethPortModel.EthPortProfileNumRadiusFilterId = types.StringNull()
			}

			if index, ok := ethPort["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					ethPortModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					ethPortModel.Index = types.Int64Value(int64(intVal))
				} else {
					ethPortModel.Index = types.Int64Null()
				}
			} else {
				ethPortModel.Index = types.Int64Null()
			}

			ethPorts = append(ethPorts, ethPortModel)
		}
		state.EthPorts = ethPorts
	} else {
		state.EthPorts = nil
	}

	if objProps, ok := aepData["object_properties"].(map[string]interface{}); ok {
		op := verityAuthenticatedEthPortObjectPropertiesModel{}
		if group, ok := objProps["group"].(string); ok {
			op.Group = types.StringValue(group)
		} else {
			op.Group = types.StringNull()
		}
		if portMon, ok := objProps["port_monitoring"].(string); ok {
			op.PortMonitoring = types.StringValue(portMon)
		} else {
			op.PortMonitoring = types.StringNull()
		}
		state.ObjectProperties = []verityAuthenticatedEthPortObjectPropertiesModel{op}
	} else {
		state.ObjectProperties = nil
	}

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

	if !plan.Name.Equal(state.Name) {
		aepProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		aepProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.ConnectionMode.Equal(state.ConnectionMode) {
		aepProps.ConnectionMode = openapi.PtrString(plan.ConnectionMode.ValueString())
		hasChanges = true
	}
	if !plan.ReauthorizationPeriodSec.Equal(state.ReauthorizationPeriodSec) {
		aepProps.ReauthorizationPeriodSec = openapi.PtrInt32(int32(plan.ReauthorizationPeriodSec.ValueInt64()))
		hasChanges = true
	}
	if !plan.AllowMacBasedAuthentication.Equal(state.AllowMacBasedAuthentication) {
		aepProps.AllowMacBasedAuthentication = openapi.PtrBool(plan.AllowMacBasedAuthentication.ValueBool())
		hasChanges = true
	}
	if !plan.MacAuthenticationHoldoffSec.Equal(state.MacAuthenticationHoldoffSec) {
		aepProps.MacAuthenticationHoldoffSec = openapi.PtrInt32(int32(plan.MacAuthenticationHoldoffSec.ValueInt64()))
		hasChanges = true
	}
	if !plan.TrustedPort.Equal(state.TrustedPort) {
		aepProps.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
		hasChanges = true
	}

	oldEthPortsByIndex := make(map[int64]verityAuthenticatedEthPortEthPortModel)
	for _, ethPort := range state.EthPorts {
		if !ethPort.Index.IsNull() {
			idx := ethPort.Index.ValueInt64()
			oldEthPortsByIndex[idx] = ethPort
		}
	}

	var changedEthPorts []openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner
	ethPortsChanged := false

	for _, planEthPort := range plan.EthPorts {
		if planEthPort.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planEthPort.Index.ValueInt64()
		stateEthPort, exists := oldEthPortsByIndex[idx]

		if !exists {
			// CREATE: new eth port, include all fields
			newEthPort := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planEthPort.EthPortProfileNumEnable.IsNull() {
				newEthPort.EthPortProfileNumEnable = openapi.PtrBool(planEthPort.EthPortProfileNumEnable.ValueBool())
			} else {
				newEthPort.EthPortProfileNumEnable = openapi.PtrBool(false)
			}

			if !planEthPort.EthPortProfileNumEthPort.IsNull() && planEthPort.EthPortProfileNumEthPort.ValueString() != "" {
				newEthPort.EthPortProfileNumEthPort = openapi.PtrString(planEthPort.EthPortProfileNumEthPort.ValueString())
			} else {
				newEthPort.EthPortProfileNumEthPort = openapi.PtrString("")
			}

			if !planEthPort.EthPortProfileNumEthPortRefType.IsNull() && planEthPort.EthPortProfileNumEthPortRefType.ValueString() != "" {
				newEthPort.EthPortProfileNumEthPortRefType = openapi.PtrString(planEthPort.EthPortProfileNumEthPortRefType.ValueString())
			} else {
				newEthPort.EthPortProfileNumEthPortRefType = openapi.PtrString("")
			}

			if !planEthPort.EthPortProfileNumWalledGardenSet.IsNull() {
				newEthPort.EthPortProfileNumWalledGardenSet = openapi.PtrBool(planEthPort.EthPortProfileNumWalledGardenSet.ValueBool())
			} else {
				newEthPort.EthPortProfileNumWalledGardenSet = openapi.PtrBool(false)
			}

			if !planEthPort.EthPortProfileNumRadiusFilterId.IsNull() && planEthPort.EthPortProfileNumRadiusFilterId.ValueString() != "" {
				newEthPort.EthPortProfileNumRadiusFilterId = openapi.PtrString(planEthPort.EthPortProfileNumRadiusFilterId.ValueString())
			} else {
				newEthPort.EthPortProfileNumRadiusFilterId = openapi.PtrString("")
			}

			changedEthPorts = append(changedEthPorts, newEthPort)
			ethPortsChanged = true
			continue
		}

		// UPDATE: existing eth port, check which fields changed
		updateEthPort := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planEthPort.EthPortProfileNumEnable.Equal(stateEthPort.EthPortProfileNumEnable) {
			updateEthPort.EthPortProfileNumEnable = openapi.PtrBool(planEthPort.EthPortProfileNumEnable.ValueBool())
			fieldChanged = true
		}

		if !planEthPort.EthPortProfileNumEthPort.Equal(stateEthPort.EthPortProfileNumEthPort) {
			if !planEthPort.EthPortProfileNumEthPort.IsNull() && planEthPort.EthPortProfileNumEthPort.ValueString() != "" {
				updateEthPort.EthPortProfileNumEthPort = openapi.PtrString(planEthPort.EthPortProfileNumEthPort.ValueString())
			} else {
				updateEthPort.EthPortProfileNumEthPort = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planEthPort.EthPortProfileNumEthPortRefType.Equal(stateEthPort.EthPortProfileNumEthPortRefType) {
			if !planEthPort.EthPortProfileNumEthPortRefType.IsNull() && planEthPort.EthPortProfileNumEthPortRefType.ValueString() != "" {
				updateEthPort.EthPortProfileNumEthPortRefType = openapi.PtrString(planEthPort.EthPortProfileNumEthPortRefType.ValueString())
			} else {
				updateEthPort.EthPortProfileNumEthPortRefType = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planEthPort.EthPortProfileNumWalledGardenSet.Equal(stateEthPort.EthPortProfileNumWalledGardenSet) {
			updateEthPort.EthPortProfileNumWalledGardenSet = openapi.PtrBool(planEthPort.EthPortProfileNumWalledGardenSet.ValueBool())
			fieldChanged = true
		}

		if !planEthPort.EthPortProfileNumRadiusFilterId.Equal(stateEthPort.EthPortProfileNumRadiusFilterId) {
			if !planEthPort.EthPortProfileNumRadiusFilterId.IsNull() && planEthPort.EthPortProfileNumRadiusFilterId.ValueString() != "" {
				updateEthPort.EthPortProfileNumRadiusFilterId = openapi.PtrString(planEthPort.EthPortProfileNumRadiusFilterId.ValueString())
			} else {
				updateEthPort.EthPortProfileNumRadiusFilterId = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedEthPorts = append(changedEthPorts, updateEthPort)
			ethPortsChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldEthPortsByIndex {
		found := false
		for _, planEthPort := range plan.EthPorts {
			if !planEthPort.Index.IsNull() && planEthPort.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// eth port removed - include only the index for deletion
			deletedEthPort := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedEthPorts = append(changedEthPorts, deletedEthPort)
			ethPortsChanged = true
		}
	}

	if ethPortsChanged && len(changedEthPorts) > 0 {
		aepProps.EthPorts = changedEthPorts
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !r.equalObjectProperties(plan.ObjectProperties[0], state.ObjectProperties[0]) {
			op := plan.ObjectProperties[0]
			objProps := openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			}
			if !op.PortMonitoring.IsNull() {
				objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
			}
			aepProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "authenticated_eth_port", name, aepProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Authenticated Eth-Port update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Authenticated Eth-Port %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "authenticated_eth_port", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Authenticated Eth-Port deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Authenticated Eth-Port %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Authenticated Eth-Port %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "authenticated_eth_ports")
	resp.State.RemoveResource(ctx)
}

func (r *verityAuthenticatedEthPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityAuthenticatedEthPortResource) equalObjectProperties(a, b verityAuthenticatedEthPortObjectPropertiesModel) bool {
	return a.Group.Equal(b.Group) && a.PortMonitoring.Equal(b.PortMonitoring)
}

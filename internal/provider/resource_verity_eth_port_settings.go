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
	_ resource.Resource                = &verityEthPortSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortSettingsResource{}
	_ resource.ResourceWithImportState = &verityEthPortSettingsResource{}
)

func NewVerityEthPortSettingsResource() resource.Resource {
	return &verityEthPortSettingsResource{}
}

type verityEthPortSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityEthPortSettingsResourceModel struct {
	Name                                   types.String                                 `tfsdk:"name"`
	Enable                                 types.Bool                                   `tfsdk:"enable"`
	ObjectProperties                       []verityEthPortSettingsObjectPropertiesModel `tfsdk:"object_properties"`
	AutoNegotiation                        types.Bool                                   `tfsdk:"auto_negotiation"`
	MaxBitRate                             types.String                                 `tfsdk:"max_bit_rate"`
	DuplexMode                             types.String                                 `tfsdk:"duplex_mode"`
	StpEnable                              types.Bool                                   `tfsdk:"stp_enable"`
	FastLearningMode                       types.Bool                                   `tfsdk:"fast_learning_mode"`
	BpduGuard                              types.Bool                                   `tfsdk:"bpdu_guard"`
	BpduFilter                             types.Bool                                   `tfsdk:"bpdu_filter"`
	GuardLoop                              types.Bool                                   `tfsdk:"guard_loop"`
	PoeEnable                              types.Bool                                   `tfsdk:"poe_enable"`
	Priority                               types.String                                 `tfsdk:"priority"`
	AllocatedPower                         types.String                                 `tfsdk:"allocated_power"`
	BspEnable                              types.Bool                                   `tfsdk:"bsp_enable"`
	Broadcast                              types.Bool                                   `tfsdk:"broadcast"`
	Multicast                              types.Bool                                   `tfsdk:"multicast"`
	MaxAllowedValue                        types.Int64                                  `tfsdk:"max_allowed_value"`
	MaxAllowedUnit                         types.String                                 `tfsdk:"max_allowed_unit"`
	Action                                 types.String                                 `tfsdk:"action"`
	Fec                                    types.String                                 `tfsdk:"fec"`
	SingleLink                             types.Bool                                   `tfsdk:"single_link"`
	MinimumWredThreshold                   types.Int64                                  `tfsdk:"minimum_wred_threshold"`
	MaximumWredThreshold                   types.Int64                                  `tfsdk:"maximum_wred_threshold"`
	WredDropProbability                    types.Int64                                  `tfsdk:"wred_drop_probability"`
	PriorityFlowControlWatchdogAction      types.String                                 `tfsdk:"priority_flow_control_watchdog_action"`
	PriorityFlowControlWatchdogDetectTime  types.Int64                                  `tfsdk:"priority_flow_control_watchdog_detect_time"`
	PriorityFlowControlWatchdogRestoreTime types.Int64                                  `tfsdk:"priority_flow_control_watchdog_restore_time"`
	PacketQueue                            types.String                                 `tfsdk:"packet_queue"`
	PacketQueueRefType                     types.String                                 `tfsdk:"packet_queue_ref_type_"`
	EnableWredTuning                       types.Bool                                   `tfsdk:"enable_wred_tuning"`
	EnableEcn                              types.Bool                                   `tfsdk:"enable_ecn"`
	EnableWatchdogTuning                   types.Bool                                   `tfsdk:"enable_watchdog_tuning"`
	CliCommands                            types.String                                 `tfsdk:"cli_commands"`
	DetectBridgingLoops                    types.Bool                                   `tfsdk:"detect_bridging_loops"`
	UnidirectionalLinkDetection            types.Bool                                   `tfsdk:"unidirectional_link_detection"`
	MacSecurityMode                        types.String                                 `tfsdk:"mac_security_mode"`
	MacLimit                               types.Int64                                  `tfsdk:"mac_limit"`
	SecurityViolationAction                types.String                                 `tfsdk:"security_violation_action"`
	AgingType                              types.String                                 `tfsdk:"aging_type"`
	AgingTime                              types.Int64                                  `tfsdk:"aging_time"`
	LldpEnable                             types.Bool                                   `tfsdk:"lldp_enable"`
	LldpMode                               types.String                                 `tfsdk:"lldp_mode"`
	LldpMedEnable                          types.Bool                                   `tfsdk:"lldp_med_enable"`
	LldpMed                                []verityEthPortSettingsLldpMedModel          `tfsdk:"lldp_med"`
}

type verityEthPortSettingsObjectPropertiesModel struct {
	Group     types.String `tfsdk:"group"`
	IsDefault types.Bool   `tfsdk:"isdefault"`
}

type verityEthPortSettingsLldpMedModel struct {
	LldpMedRowNumEnable                types.Bool   `tfsdk:"lldp_med_row_num_enable"`
	LldpMedRowNumAdvertisedApplication types.String `tfsdk:"lldp_med_row_num_advertised_applicatio"`
	LldpMedRowNumDscpMark              types.Int64  `tfsdk:"lldp_med_row_num_dscp_mark"`
	LldpMedRowNumPriority              types.Int64  `tfsdk:"lldp_med_row_num_priority"`
	LldpMedRowNumService               types.String `tfsdk:"lldp_med_row_num_service"`
	LldpMedRowNumServiceRefType        types.String `tfsdk:"lldp_med_row_num_service_ref_type_"`
	Index                              types.Int64  `tfsdk:"index"`
}

func (lm verityEthPortSettingsLldpMedModel) GetIndex() types.Int64 {
	return lm.Index
}

func (r *verityEthPortSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eth_port_settings"
}

func (r *verityEthPortSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityEthPortSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Ethernet port settings",
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
			"auto_negotiation": schema.BoolAttribute{
				Description: "Indicates if port speed and duplex mode should be auto negotiated",
				Optional:    true,
			},
			"max_bit_rate": schema.StringAttribute{
				Description: "Maximum Bit Rate allowed",
				Optional:    true,
			},
			"duplex_mode": schema.StringAttribute{
				Description: "Duplex Mode",
				Optional:    true,
			},
			"stp_enable": schema.BoolAttribute{
				Description: "Enable Spanning Tree on the port. Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings",
				Optional:    true,
			},
			"fast_learning_mode": schema.BoolAttribute{
				Description: "Enable Immediate Transition to Forwarding",
				Optional:    true,
			},
			"bpdu_guard": schema.BoolAttribute{
				Description: "Block port on BPDU Receive",
				Optional:    true,
			},
			"bpdu_filter": schema.BoolAttribute{
				Description: "Drop all Rx and Tx BPDUs",
				Optional:    true,
			},
			"guard_loop": schema.BoolAttribute{
				Description: "Enable Cisco Guard Loop",
				Optional:    true,
			},
			"poe_enable": schema.BoolAttribute{
				Description: "PoE Enable",
				Optional:    true,
			},
			"priority": schema.StringAttribute{
				Description: "Priority given when assigning power in a limited power situation",
				Optional:    true,
			},
			"allocated_power": schema.StringAttribute{
				Description: "Power the PoE system will attempt to allocate on this port",
				Optional:    true,
			},
			"bsp_enable": schema.BoolAttribute{
				Description: "Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU",
				Optional:    true,
			},
			"broadcast": schema.BoolAttribute{
				Description: "Broadcast",
				Optional:    true,
			},
			"multicast": schema.BoolAttribute{
				Description: "Multicast",
				Optional:    true,
			},
			"max_allowed_value": schema.Int64Attribute{
				Description: "Max Percentage of the port's bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action",
				Optional:    true,
			},
			"max_allowed_unit": schema.StringAttribute{
				Description: "Max Percentage unit for broadcast/multicast/unknown-unicast traffic",
				Optional:    true,
			},
			"action": schema.StringAttribute{
				Description: "Action taken if broadcast/multicast/unknown-unicast traffic exceeds the Max",
				Optional:    true,
			},
			"fec": schema.StringAttribute{
				Description: "FEC is Forward Error Correction which is error correction on the fiber link",
				Optional:    true,
			},
			"single_link": schema.BoolAttribute{
				Description: "Ports with this setting will be disabled when link state tracking takes effect",
				Optional:    true,
			},
			"minimum_wred_threshold": schema.Int64Attribute{
				Description: "A value between 1 to 12480(in KiloBytes)",
				Optional:    true,
			},
			"maximum_wred_threshold": schema.Int64Attribute{
				Description: "A value between 1 to 12480(in KiloBytes)",
				Optional:    true,
			},
			"wred_drop_probability": schema.Int64Attribute{
				Description: "A value between 0 to 100",
				Optional:    true,
			},
			"priority_flow_control_watchdog_action": schema.StringAttribute{
				Description: "Ports with this setting will be disabled when link state tracking takes effect",
				Optional:    true,
			},
			"priority_flow_control_watchdog_detect_time": schema.Int64Attribute{
				Description: "A value between 100 to 5000",
				Optional:    true,
			},
			"priority_flow_control_watchdog_restore_time": schema.Int64Attribute{
				Description: "A value between 100 to 60000",
				Optional:    true,
			},
			"packet_queue": schema.StringAttribute{
				Description: "Packet Queue",
				Optional:    true,
			},
			"packet_queue_ref_type_": schema.StringAttribute{
				Description: "Object type for packet_queue field",
				Optional:    true,
			},
			"enable_wred_tuning": schema.BoolAttribute{
				Description: "Enables custom tuning of WRED values. Uncheck to use Switch default values.",
				Optional:    true,
			},
			"enable_ecn": schema.BoolAttribute{
				Description: "Enables Explicit Congestion Notification for WRED.",
				Optional:    true,
			},
			"enable_watchdog_tuning": schema.BoolAttribute{
				Description: "Enables custom tuning of Watchdog values. Uncheck to use Switch default values.",
				Optional:    true,
			},
			"cli_commands": schema.StringAttribute{
				Description: "CLI Commands",
				Optional:    true,
			},
			"detect_bridging_loops": schema.BoolAttribute{
				Description: "Enable Detection of Bridging Loops",
				Optional:    true,
			},
			"unidirectional_link_detection": schema.BoolAttribute{
				Description: "Enable Detection of Unidirectional Link",
				Optional:    true,
			},
			"mac_security_mode": schema.StringAttribute{
				Description: "MAC security mode",
				Optional:    true,
			},
			"mac_limit": schema.Int64Attribute{
				Description: "Between 1-1000",
				Optional:    true,
			},
			"security_violation_action": schema.StringAttribute{
				Description: "Security violation action",
				Optional:    true,
			},
			"aging_type": schema.StringAttribute{
				Description: "Limit MAC authentication based on inactivity or on absolute time. See Also Aging Time",
				Optional:    true,
			},
			"aging_time": schema.Int64Attribute{
				Description: "In minutes, how long the client will stay authenticated. See Also Aging Type",
				Optional:    true,
			},
			"lldp_enable": schema.BoolAttribute{
				Description: "LLDP enable",
				Optional:    true,
			},
			"lldp_mode": schema.StringAttribute{
				Description: "LLDP mode. Enables LLDP Rx and/or LLDP Tx",
				Optional:    true,
			},
			"lldp_med_enable": schema.BoolAttribute{
				Description: "LLDP med enable",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the eth port settings",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"isdefault": schema.BoolAttribute{
							Description: "Default object.",
							Optional:    true,
						},
					},
				},
			},
			"lldp_med": schema.ListNestedBlock{
				Description: "LLDP MED configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"lldp_med_row_num_enable": schema.BoolAttribute{
							Description: "Per LLDP Med row enable",
							Optional:    true,
						},
						"lldp_med_row_num_advertised_applicatio": schema.StringAttribute{
							Description: "Advertised application",
							Optional:    true,
						},
						"lldp_med_row_num_dscp_mark": schema.Int64Attribute{
							Description: "LLDP DSCP Mark",
							Optional:    true,
						},
						"lldp_med_row_num_priority": schema.Int64Attribute{
							Description: "LLDP Priority",
							Optional:    true,
						},
						"lldp_med_row_num_service": schema.StringAttribute{
							Description: "LLDP Service",
							Optional:    true,
						},
						"lldp_med_row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for lldp_med_row_num_service field",
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

func (r *verityEthPortSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityEthPortSettingsResourceModel
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
	ethPortSettingsProps := &openapi.EthportsettingsPutRequestEthPortSettingsValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "MaxBitRate", APIField: &ethPortSettingsProps.MaxBitRate, TFValue: plan.MaxBitRate},
		{FieldName: "DuplexMode", APIField: &ethPortSettingsProps.DuplexMode, TFValue: plan.DuplexMode},
		{FieldName: "Priority", APIField: &ethPortSettingsProps.Priority, TFValue: plan.Priority},
		{FieldName: "AllocatedPower", APIField: &ethPortSettingsProps.AllocatedPower, TFValue: plan.AllocatedPower},
		{FieldName: "MaxAllowedUnit", APIField: &ethPortSettingsProps.MaxAllowedUnit, TFValue: plan.MaxAllowedUnit},
		{FieldName: "Action", APIField: &ethPortSettingsProps.Action, TFValue: plan.Action},
		{FieldName: "Fec", APIField: &ethPortSettingsProps.Fec, TFValue: plan.Fec},
		{FieldName: "PriorityFlowControlWatchdogAction", APIField: &ethPortSettingsProps.PriorityFlowControlWatchdogAction, TFValue: plan.PriorityFlowControlWatchdogAction},
		{FieldName: "PacketQueue", APIField: &ethPortSettingsProps.PacketQueue, TFValue: plan.PacketQueue},
		{FieldName: "PacketQueueRefType", APIField: &ethPortSettingsProps.PacketQueueRefType, TFValue: plan.PacketQueueRefType},
		{FieldName: "CliCommands", APIField: &ethPortSettingsProps.CliCommands, TFValue: plan.CliCommands},
		{FieldName: "MacSecurityMode", APIField: &ethPortSettingsProps.MacSecurityMode, TFValue: plan.MacSecurityMode},
		{FieldName: "SecurityViolationAction", APIField: &ethPortSettingsProps.SecurityViolationAction, TFValue: plan.SecurityViolationAction},
		{FieldName: "AgingType", APIField: &ethPortSettingsProps.AgingType, TFValue: plan.AgingType},
		{FieldName: "LldpMode", APIField: &ethPortSettingsProps.LldpMode, TFValue: plan.LldpMode},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ethPortSettingsProps.Enable, TFValue: plan.Enable},
		{FieldName: "AutoNegotiation", APIField: &ethPortSettingsProps.AutoNegotiation, TFValue: plan.AutoNegotiation},
		{FieldName: "StpEnable", APIField: &ethPortSettingsProps.StpEnable, TFValue: plan.StpEnable},
		{FieldName: "FastLearningMode", APIField: &ethPortSettingsProps.FastLearningMode, TFValue: plan.FastLearningMode},
		{FieldName: "BpduGuard", APIField: &ethPortSettingsProps.BpduGuard, TFValue: plan.BpduGuard},
		{FieldName: "BpduFilter", APIField: &ethPortSettingsProps.BpduFilter, TFValue: plan.BpduFilter},
		{FieldName: "GuardLoop", APIField: &ethPortSettingsProps.GuardLoop, TFValue: plan.GuardLoop},
		{FieldName: "PoeEnable", APIField: &ethPortSettingsProps.PoeEnable, TFValue: plan.PoeEnable},
		{FieldName: "BspEnable", APIField: &ethPortSettingsProps.BspEnable, TFValue: plan.BspEnable},
		{FieldName: "Broadcast", APIField: &ethPortSettingsProps.Broadcast, TFValue: plan.Broadcast},
		{FieldName: "Multicast", APIField: &ethPortSettingsProps.Multicast, TFValue: plan.Multicast},
		{FieldName: "SingleLink", APIField: &ethPortSettingsProps.SingleLink, TFValue: plan.SingleLink},
		{FieldName: "EnableWredTuning", APIField: &ethPortSettingsProps.EnableWredTuning, TFValue: plan.EnableWredTuning},
		{FieldName: "EnableEcn", APIField: &ethPortSettingsProps.EnableEcn, TFValue: plan.EnableEcn},
		{FieldName: "EnableWatchdogTuning", APIField: &ethPortSettingsProps.EnableWatchdogTuning, TFValue: plan.EnableWatchdogTuning},
		{FieldName: "DetectBridgingLoops", APIField: &ethPortSettingsProps.DetectBridgingLoops, TFValue: plan.DetectBridgingLoops},
		{FieldName: "UnidirectionalLinkDetection", APIField: &ethPortSettingsProps.UnidirectionalLinkDetection, TFValue: plan.UnidirectionalLinkDetection},
		{FieldName: "LldpEnable", APIField: &ethPortSettingsProps.LldpEnable, TFValue: plan.LldpEnable},
		{FieldName: "LldpMedEnable", APIField: &ethPortSettingsProps.LldpMedEnable, TFValue: plan.LldpMedEnable},
	})

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "MaxAllowedValue", APIField: &ethPortSettingsProps.MaxAllowedValue, TFValue: plan.MaxAllowedValue},
		{FieldName: "MinimumWredThreshold", APIField: &ethPortSettingsProps.MinimumWredThreshold, TFValue: plan.MinimumWredThreshold},
		{FieldName: "MaximumWredThreshold", APIField: &ethPortSettingsProps.MaximumWredThreshold, TFValue: plan.MaximumWredThreshold},
		{FieldName: "WredDropProbability", APIField: &ethPortSettingsProps.WredDropProbability, TFValue: plan.WredDropProbability},
		{FieldName: "PriorityFlowControlWatchdogDetectTime", APIField: &ethPortSettingsProps.PriorityFlowControlWatchdogDetectTime, TFValue: plan.PriorityFlowControlWatchdogDetectTime},
		{FieldName: "PriorityFlowControlWatchdogRestoreTime", APIField: &ethPortSettingsProps.PriorityFlowControlWatchdogRestoreTime, TFValue: plan.PriorityFlowControlWatchdogRestoreTime},
		{FieldName: "MacLimit", APIField: &ethPortSettingsProps.MacLimit, TFValue: plan.MacLimit},
		{FieldName: "AgingTime", APIField: &ethPortSettingsProps.AgingTime, TFValue: plan.AgingTime},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		}
		ethPortSettingsProps.ObjectProperties = &objProps
	}

	// Handle LLDP Med
	if len(plan.LldpMed) > 0 {
		lldpMedItems := make([]openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner, len(plan.LldpMed))
		for i, lldpMedItem := range plan.LldpMed {
			item := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

			if !lldpMedItem.LldpMedRowNumEnable.IsNull() {
				item.LldpMedRowNumEnable = openapi.PtrBool(lldpMedItem.LldpMedRowNumEnable.ValueBool())
			}
			if !lldpMedItem.LldpMedRowNumAdvertisedApplication.IsNull() {
				item.LldpMedRowNumAdvertisedApplicatio = openapi.PtrString(lldpMedItem.LldpMedRowNumAdvertisedApplication.ValueString())
			}
			if !lldpMedItem.LldpMedRowNumDscpMark.IsNull() {
				item.LldpMedRowNumDscpMark = openapi.PtrInt32(int32(lldpMedItem.LldpMedRowNumDscpMark.ValueInt64()))
			}
			if !lldpMedItem.LldpMedRowNumPriority.IsNull() {
				item.LldpMedRowNumPriority = openapi.PtrInt32(int32(lldpMedItem.LldpMedRowNumPriority.ValueInt64()))
			}
			if !lldpMedItem.LldpMedRowNumService.IsNull() {
				item.LldpMedRowNumService = openapi.PtrString(lldpMedItem.LldpMedRowNumService.ValueString())
			}
			if !lldpMedItem.LldpMedRowNumServiceRefType.IsNull() {
				item.LldpMedRowNumServiceRefType = openapi.PtrString(lldpMedItem.LldpMedRowNumServiceRefType.ValueString())
			}
			if !lldpMedItem.Index.IsNull() {
				item.Index = openapi.PtrInt32(int32(lldpMedItem.Index.ValueInt64()))
			}

			lldpMedItems[i] = item
		}
		ethPortSettingsProps.LldpMed = lldpMedItems
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_settings", name, *ethPortSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityEthPortSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityEthPortSettingsResourceModel
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

	ethPortSettingsName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("eth_port_settings") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port settings %s verification – trusting recent successful API operation", ethPortSettingsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching eth port settings for verification of %s", ethPortSettingsName))

	type EthPortSettingsResponse struct {
		EthPortSettings map[string]interface{} `json:"eth_port_settings"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "eth_port_settings", ethPortSettingsName,
		func() (EthPortSettingsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch eth port settings")
			respAPI, err := r.client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
			if err != nil {
				return EthPortSettingsResponse{}, fmt.Errorf("error reading eth port settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res EthPortSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return EthPortSettingsResponse{}, fmt.Errorf("failed to decode eth port settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d eth port settings", len(res.EthPortSettings)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Eth Port Settings %s", ethPortSettingsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for eth port settings with name: %s", ethPortSettingsName))

	ethPortSettingsData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EthPortSettings,
		ethPortSettingsName,
		func(data interface{}) (string, bool) {
			if ethPortSettings, ok := data.(map[string]interface{}); ok {
				if name, ok := ethPortSettings["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Eth Port Settings with name '%s' not found in API response", ethPortSettingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	ethPortSettingsMap, ok := ethPortSettingsData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Eth Port Settings Data",
			fmt.Sprintf("Eth Port Settings data is not in expected format for %s", ethPortSettingsName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found eth port settings '%s' under API key '%s'", ethPortSettingsName, actualAPIName))

	state.Name = utils.MapStringFromAPI(ethPortSettingsMap["name"])

	// Handle object properties
	if objProps, ok := ethPortSettingsMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityEthPortSettingsObjectPropertiesModel{
			{
				Group:     utils.MapStringFromAPI(objProps["group"]),
				IsDefault: utils.MapBoolFromAPI(objProps["isdefault"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"max_bit_rate":                          &state.MaxBitRate,
		"duplex_mode":                           &state.DuplexMode,
		"priority":                              &state.Priority,
		"allocated_power":                       &state.AllocatedPower,
		"max_allowed_unit":                      &state.MaxAllowedUnit,
		"action":                                &state.Action,
		"fec":                                   &state.Fec,
		"priority_flow_control_watchdog_action": &state.PriorityFlowControlWatchdogAction,
		"packet_queue":                          &state.PacketQueue,
		"packet_queue_ref_type_":                &state.PacketQueueRefType,
		"cli_commands":                          &state.CliCommands,
		"mac_security_mode":                     &state.MacSecurityMode,
		"security_violation_action":             &state.SecurityViolationAction,
		"aging_type":                            &state.AgingType,
		"lldp_mode":                             &state.LldpMode,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(ethPortSettingsMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":                        &state.Enable,
		"auto_negotiation":              &state.AutoNegotiation,
		"stp_enable":                    &state.StpEnable,
		"fast_learning_mode":            &state.FastLearningMode,
		"bpdu_guard":                    &state.BpduGuard,
		"bpdu_filter":                   &state.BpduFilter,
		"guard_loop":                    &state.GuardLoop,
		"poe_enable":                    &state.PoeEnable,
		"bsp_enable":                    &state.BspEnable,
		"broadcast":                     &state.Broadcast,
		"multicast":                     &state.Multicast,
		"single_link":                   &state.SingleLink,
		"enable_wred_tuning":            &state.EnableWredTuning,
		"enable_ecn":                    &state.EnableEcn,
		"enable_watchdog_tuning":        &state.EnableWatchdogTuning,
		"detect_bridging_loops":         &state.DetectBridgingLoops,
		"unidirectional_link_detection": &state.UnidirectionalLinkDetection,
		"lldp_enable":                   &state.LldpEnable,
		"lldp_med_enable":               &state.LldpMedEnable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ethPortSettingsMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"max_allowed_value":                           &state.MaxAllowedValue,
		"minimum_wred_threshold":                      &state.MinimumWredThreshold,
		"maximum_wred_threshold":                      &state.MaximumWredThreshold,
		"wred_drop_probability":                       &state.WredDropProbability,
		"priority_flow_control_watchdog_detect_time":  &state.PriorityFlowControlWatchdogDetectTime,
		"priority_flow_control_watchdog_restore_time": &state.PriorityFlowControlWatchdogRestoreTime,
		"mac_limit":  &state.MacLimit,
		"aging_time": &state.AgingTime,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(ethPortSettingsMap[apiKey])
	}

	// Handle LLDP Med
	if lldpMedData, ok := ethPortSettingsMap["lldp_med"].([]interface{}); ok && len(lldpMedData) > 0 {
		var lldpMedList []verityEthPortSettingsLldpMedModel

		for _, item := range lldpMedData {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			lldpMedItem := verityEthPortSettingsLldpMedModel{
				LldpMedRowNumEnable:                utils.MapBoolFromAPI(itemMap["lldp_med_row_num_enable"]),
				LldpMedRowNumAdvertisedApplication: utils.MapStringFromAPI(itemMap["lldp_med_row_num_advertised_applicatio"]),
				LldpMedRowNumDscpMark:              utils.MapInt64FromAPI(itemMap["lldp_med_row_num_dscp_mark"]),
				LldpMedRowNumPriority:              utils.MapInt64FromAPI(itemMap["lldp_med_row_num_priority"]),
				LldpMedRowNumService:               utils.MapStringFromAPI(itemMap["lldp_med_row_num_service"]),
				LldpMedRowNumServiceRefType:        utils.MapStringFromAPI(itemMap["lldp_med_row_num_service_ref_type_"]),
				Index:                              utils.MapInt64FromAPI(itemMap["index"]),
			}

			lldpMedList = append(lldpMedList, lldpMedItem)
		}

		state.LldpMed = lldpMedList
	} else {
		state.LldpMed = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityEthPortSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan verityEthPortSettingsResourceModel
	var state verityEthPortSettingsResourceModel

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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ethPortSettingsProps := openapi.EthportsettingsPutRequestEthPortSettingsValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ethPortSettingsProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MaxBitRate, state.MaxBitRate, func(v *string) { ethPortSettingsProps.MaxBitRate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DuplexMode, state.DuplexMode, func(v *string) { ethPortSettingsProps.DuplexMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Priority, state.Priority, func(v *string) { ethPortSettingsProps.Priority = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AllocatedPower, state.AllocatedPower, func(v *string) { ethPortSettingsProps.AllocatedPower = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MaxAllowedUnit, state.MaxAllowedUnit, func(v *string) { ethPortSettingsProps.MaxAllowedUnit = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Action, state.Action, func(v *string) { ethPortSettingsProps.Action = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Fec, state.Fec, func(v *string) { ethPortSettingsProps.Fec = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PriorityFlowControlWatchdogAction, state.PriorityFlowControlWatchdogAction, func(v *string) { ethPortSettingsProps.PriorityFlowControlWatchdogAction = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CliCommands, state.CliCommands, func(v *string) { ethPortSettingsProps.CliCommands = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MacSecurityMode, state.MacSecurityMode, func(v *string) { ethPortSettingsProps.MacSecurityMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SecurityViolationAction, state.SecurityViolationAction, func(v *string) { ethPortSettingsProps.SecurityViolationAction = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AgingType, state.AgingType, func(v *string) { ethPortSettingsProps.AgingType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LldpMode, state.LldpMode, func(v *string) { ethPortSettingsProps.LldpMode = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ethPortSettingsProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AutoNegotiation, state.AutoNegotiation, func(v *bool) { ethPortSettingsProps.AutoNegotiation = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.StpEnable, state.StpEnable, func(v *bool) { ethPortSettingsProps.StpEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FastLearningMode, state.FastLearningMode, func(v *bool) { ethPortSettingsProps.FastLearningMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BpduGuard, state.BpduGuard, func(v *bool) { ethPortSettingsProps.BpduGuard = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BpduFilter, state.BpduFilter, func(v *bool) { ethPortSettingsProps.BpduFilter = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.GuardLoop, state.GuardLoop, func(v *bool) { ethPortSettingsProps.GuardLoop = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.PoeEnable, state.PoeEnable, func(v *bool) { ethPortSettingsProps.PoeEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BspEnable, state.BspEnable, func(v *bool) { ethPortSettingsProps.BspEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Broadcast, state.Broadcast, func(v *bool) { ethPortSettingsProps.Broadcast = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Multicast, state.Multicast, func(v *bool) { ethPortSettingsProps.Multicast = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.SingleLink, state.SingleLink, func(v *bool) { ethPortSettingsProps.SingleLink = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableWredTuning, state.EnableWredTuning, func(v *bool) { ethPortSettingsProps.EnableWredTuning = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableEcn, state.EnableEcn, func(v *bool) { ethPortSettingsProps.EnableEcn = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableWatchdogTuning, state.EnableWatchdogTuning, func(v *bool) { ethPortSettingsProps.EnableWatchdogTuning = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.DetectBridgingLoops, state.DetectBridgingLoops, func(v *bool) { ethPortSettingsProps.DetectBridgingLoops = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.UnidirectionalLinkDetection, state.UnidirectionalLinkDetection, func(v *bool) { ethPortSettingsProps.UnidirectionalLinkDetection = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.LldpEnable, state.LldpEnable, func(v *bool) { ethPortSettingsProps.LldpEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.LldpMedEnable, state.LldpMedEnable, func(v *bool) { ethPortSettingsProps.LldpMedEnable = v }, &hasChanges)

	// Handle int64 field changes
	utils.CompareAndSetInt64Field(plan.MaxAllowedValue, state.MaxAllowedValue, func(v *int32) { ethPortSettingsProps.MaxAllowedValue = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.MinimumWredThreshold, state.MinimumWredThreshold, func(v *int32) { ethPortSettingsProps.MinimumWredThreshold = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.MaximumWredThreshold, state.MaximumWredThreshold, func(v *int32) { ethPortSettingsProps.MaximumWredThreshold = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.WredDropProbability, state.WredDropProbability, func(v *int32) { ethPortSettingsProps.WredDropProbability = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.PriorityFlowControlWatchdogDetectTime, state.PriorityFlowControlWatchdogDetectTime, func(v *int32) { ethPortSettingsProps.PriorityFlowControlWatchdogDetectTime = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.PriorityFlowControlWatchdogRestoreTime, state.PriorityFlowControlWatchdogRestoreTime, func(v *int32) { ethPortSettingsProps.PriorityFlowControlWatchdogRestoreTime = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.MacLimit, state.MacLimit, func(v *int32) { ethPortSettingsProps.MacLimit = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.AgingTime, state.AgingTime, func(v *int32) { ethPortSettingsProps.AgingTime = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) {
			objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			if !plan.ObjectProperties[0].IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(plan.ObjectProperties[0].IsDefault.ValueBool())
			}
			ethPortSettingsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle packet_queue and packet_queue_ref_type_ using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.PacketQueue, state.PacketQueue, plan.PacketQueueRefType, state.PacketQueueRefType,
		func(v *string) { ethPortSettingsProps.PacketQueue = v },
		func(v *string) { ethPortSettingsProps.PacketQueueRefType = v },
		"packet_queue", "packet_queue_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle LLDP Med
	lldpMedHandler := utils.IndexedItemHandler[verityEthPortSettingsLldpMedModel, openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner]{
		CreateNew: func(item verityEthPortSettingsLldpMedModel) openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

			if !item.Index.IsNull() {
				lldpMedItem.Index = openapi.PtrInt32(int32(item.Index.ValueInt64()))
			}
			if !item.LldpMedRowNumEnable.IsNull() {
				lldpMedItem.LldpMedRowNumEnable = openapi.PtrBool(item.LldpMedRowNumEnable.ValueBool())
			}
			if !item.LldpMedRowNumAdvertisedApplication.IsNull() {
				lldpMedItem.LldpMedRowNumAdvertisedApplicatio = openapi.PtrString(item.LldpMedRowNumAdvertisedApplication.ValueString())
			}
			if !item.LldpMedRowNumDscpMark.IsNull() {
				lldpMedItem.LldpMedRowNumDscpMark = openapi.PtrInt32(int32(item.LldpMedRowNumDscpMark.ValueInt64()))
			}
			if !item.LldpMedRowNumPriority.IsNull() {
				lldpMedItem.LldpMedRowNumPriority = openapi.PtrInt32(int32(item.LldpMedRowNumPriority.ValueInt64()))
			}
			if !item.LldpMedRowNumService.IsNull() {
				lldpMedItem.LldpMedRowNumService = openapi.PtrString(item.LldpMedRowNumService.ValueString())
			}
			if !item.LldpMedRowNumServiceRefType.IsNull() {
				lldpMedItem.LldpMedRowNumServiceRefType = openapi.PtrString(item.LldpMedRowNumServiceRefType.ValueString())
			}

			return lldpMedItem
		},
		UpdateExisting: func(planItem, stateItem verityEthPortSettingsLldpMedModel) (openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner, bool) {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}
			hasChanges := false

			if !planItem.LldpMedRowNumEnable.Equal(stateItem.LldpMedRowNumEnable) {
				lldpMedItem.LldpMedRowNumEnable = openapi.PtrBool(planItem.LldpMedRowNumEnable.ValueBool())
				hasChanges = true
			}
			if !planItem.LldpMedRowNumAdvertisedApplication.Equal(stateItem.LldpMedRowNumAdvertisedApplication) {
				if !planItem.LldpMedRowNumAdvertisedApplication.IsNull() {
					lldpMedItem.LldpMedRowNumAdvertisedApplicatio = openapi.PtrString(planItem.LldpMedRowNumAdvertisedApplication.ValueString())
				}
				hasChanges = true
			}
			if !planItem.LldpMedRowNumDscpMark.Equal(stateItem.LldpMedRowNumDscpMark) {
				lldpMedItem.LldpMedRowNumDscpMark = openapi.PtrInt32(int32(planItem.LldpMedRowNumDscpMark.ValueInt64()))
				hasChanges = true
			}
			if !planItem.LldpMedRowNumPriority.Equal(stateItem.LldpMedRowNumPriority) {
				lldpMedItem.LldpMedRowNumPriority = openapi.PtrInt32(int32(planItem.LldpMedRowNumPriority.ValueInt64()))
				hasChanges = true
			}
			if !planItem.LldpMedRowNumService.Equal(stateItem.LldpMedRowNumService) {
				if !planItem.LldpMedRowNumService.IsNull() {
					lldpMedItem.LldpMedRowNumService = openapi.PtrString(planItem.LldpMedRowNumService.ValueString())
				}
				hasChanges = true
			}
			if !planItem.LldpMedRowNumServiceRefType.Equal(stateItem.LldpMedRowNumServiceRefType) {
				if !planItem.LldpMedRowNumServiceRefType.IsNull() {
					lldpMedItem.LldpMedRowNumServiceRefType = openapi.PtrString(planItem.LldpMedRowNumServiceRefType.ValueString())
				}
				hasChanges = true
			}

			return lldpMedItem, hasChanges
		},
		CreateDeleted: func(index int64) openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
			return openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedLldpMed, lldpMedChanged := utils.ProcessIndexedArrayUpdates(plan.LldpMed, state.LldpMed, lldpMedHandler)

	if lldpMedChanged {
		ethPortSettingsProps.LldpMed = changedLldpMed
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_settings", name, ethPortSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityEthPortSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityEthPortSettingsResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_settings", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

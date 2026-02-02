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
	_ resource.Resource                = &verityEthPortSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortSettingsResource{}
	_ resource.ResourceWithImportState = &verityEthPortSettingsResource{}
	_ resource.ResourceWithModifyPlan  = &verityEthPortSettingsResource{}
)

const ethPortSettingsResourceType = "ethportsettings"
const ethPortSettingsTerraformType = "verity_eth_port_settings"

func NewVerityEthPortSettingsResource() resource.Resource {
	return &verityEthPortSettingsResource{}
}

type verityEthPortSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
	Group types.String `tfsdk:"group"`
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
				Computed:    true,
			},
			"auto_negotiation": schema.BoolAttribute{
				Description: "Indicates if port speed and duplex mode should be auto negotiated",
				Optional:    true,
				Computed:    true,
			},
			"max_bit_rate": schema.StringAttribute{
				Description: "Maximum Bit Rate allowed",
				Optional:    true,
				Computed:    true,
			},
			"duplex_mode": schema.StringAttribute{
				Description: "Duplex Mode",
				Optional:    true,
				Computed:    true,
			},
			"stp_enable": schema.BoolAttribute{
				Description: "Enable Spanning Tree on the port. Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings",
				Optional:    true,
				Computed:    true,
			},
			"fast_learning_mode": schema.BoolAttribute{
				Description: "Enable Immediate Transition to Forwarding",
				Optional:    true,
				Computed:    true,
			},
			"bpdu_guard": schema.BoolAttribute{
				Description: "Block port on BPDU Receive",
				Optional:    true,
				Computed:    true,
			},
			"bpdu_filter": schema.BoolAttribute{
				Description: "Drop all Rx and Tx BPDUs",
				Optional:    true,
				Computed:    true,
			},
			"guard_loop": schema.BoolAttribute{
				Description: "Enable Cisco Guard Loop",
				Optional:    true,
				Computed:    true,
			},
			"poe_enable": schema.BoolAttribute{
				Description: "PoE Enable",
				Optional:    true,
				Computed:    true,
			},
			"priority": schema.StringAttribute{
				Description: "Priority given when assigning power in a limited power situation",
				Optional:    true,
				Computed:    true,
			},
			"allocated_power": schema.StringAttribute{
				Description: "Power the PoE system will attempt to allocate on this port",
				Optional:    true,
				Computed:    true,
			},
			"bsp_enable": schema.BoolAttribute{
				Description: "Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU",
				Optional:    true,
				Computed:    true,
			},
			"broadcast": schema.BoolAttribute{
				Description: "Broadcast",
				Optional:    true,
				Computed:    true,
			},
			"multicast": schema.BoolAttribute{
				Description: "Multicast",
				Optional:    true,
				Computed:    true,
			},
			"max_allowed_value": schema.Int64Attribute{
				Description: "Max Percentage of the port's bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action",
				Optional:    true,
				Computed:    true,
			},
			"max_allowed_unit": schema.StringAttribute{
				Description: "Max Percentage unit for broadcast/multicast/unknown-unicast traffic",
				Optional:    true,
				Computed:    true,
			},
			"action": schema.StringAttribute{
				Description: "Action taken if broadcast/multicast/unknown-unicast traffic exceeds the Max",
				Optional:    true,
				Computed:    true,
			},
			"fec": schema.StringAttribute{
				Description: "FEC is Forward Error Correction which is error correction on the fiber link",
				Optional:    true,
				Computed:    true,
			},
			"single_link": schema.BoolAttribute{
				Description: "Ports with this setting will be disabled when link state tracking takes effect",
				Optional:    true,
				Computed:    true,
			},
			"minimum_wred_threshold": schema.Int64Attribute{
				Description: "A value between 1 to 12480(in KiloBytes)",
				Optional:    true,
				Computed:    true,
			},
			"maximum_wred_threshold": schema.Int64Attribute{
				Description: "A value between 1 to 12480(in KiloBytes)",
				Optional:    true,
				Computed:    true,
			},
			"wred_drop_probability": schema.Int64Attribute{
				Description: "A value between 0 to 100",
				Optional:    true,
				Computed:    true,
			},
			"priority_flow_control_watchdog_action": schema.StringAttribute{
				Description: "Ports with this setting will be disabled when link state tracking takes effect",
				Optional:    true,
				Computed:    true,
			},
			"priority_flow_control_watchdog_detect_time": schema.Int64Attribute{
				Description: "A value between 100 to 5000",
				Optional:    true,
				Computed:    true,
			},
			"priority_flow_control_watchdog_restore_time": schema.Int64Attribute{
				Description: "A value between 100 to 60000",
				Optional:    true,
				Computed:    true,
			},
			"packet_queue": schema.StringAttribute{
				Description: "Packet Queue",
				Optional:    true,
				Computed:    true,
			},
			"packet_queue_ref_type_": schema.StringAttribute{
				Description: "Object type for packet_queue field",
				Optional:    true,
				Computed:    true,
			},
			"enable_wred_tuning": schema.BoolAttribute{
				Description: "Enables custom tuning of WRED values. Uncheck to use Switch default values.",
				Optional:    true,
				Computed:    true,
			},
			"enable_ecn": schema.BoolAttribute{
				Description: "Enables Explicit Congestion Notification for WRED.",
				Optional:    true,
				Computed:    true,
			},
			"enable_watchdog_tuning": schema.BoolAttribute{
				Description: "Enables custom tuning of Watchdog values. Uncheck to use Switch default values.",
				Optional:    true,
				Computed:    true,
			},
			"cli_commands": schema.StringAttribute{
				Description: "CLI Commands",
				Optional:    true,
				Computed:    true,
			},
			"detect_bridging_loops": schema.BoolAttribute{
				Description: "Enable Detection of Bridging Loops",
				Optional:    true,
				Computed:    true,
			},
			"unidirectional_link_detection": schema.BoolAttribute{
				Description: "Enable Detection of Unidirectional Link",
				Optional:    true,
				Computed:    true,
			},
			"mac_security_mode": schema.StringAttribute{
				Description: "MAC security mode",
				Optional:    true,
				Computed:    true,
			},
			"mac_limit": schema.Int64Attribute{
				Description: "Between 1-1000",
				Optional:    true,
				Computed:    true,
			},
			"security_violation_action": schema.StringAttribute{
				Description: "Security violation action",
				Optional:    true,
				Computed:    true,
			},
			"aging_type": schema.StringAttribute{
				Description: "Limit MAC authentication based on inactivity or on absolute time. See Also Aging Time",
				Optional:    true,
				Computed:    true,
			},
			"aging_time": schema.Int64Attribute{
				Description: "In minutes, how long the client will stay authenticated. See Also Aging Type",
				Optional:    true,
				Computed:    true,
			},
			"lldp_enable": schema.BoolAttribute{
				Description: "LLDP enable",
				Optional:    true,
				Computed:    true,
			},
			"lldp_mode": schema.StringAttribute{
				Description: "LLDP mode. Enables LLDP Rx and/or LLDP Tx",
				Optional:    true,
				Computed:    true,
			},
			"lldp_med_enable": schema.BoolAttribute{
				Description: "LLDP med enable",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
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
							Computed:    true,
						},
						"lldp_med_row_num_advertised_applicatio": schema.StringAttribute{
							Description: "Advertised application",
							Optional:    true,
							Computed:    true,
						},
						"lldp_med_row_num_dscp_mark": schema.Int64Attribute{
							Description: "LLDP DSCP Mark",
							Optional:    true,
							Computed:    true,
						},
						"lldp_med_row_num_priority": schema.Int64Attribute{
							Description: "LLDP Priority",
							Optional:    true,
							Computed:    true,
						},
						"lldp_med_row_num_service": schema.StringAttribute{
							Description: "LLDP Service",
							Optional:    true,
							Computed:    true,
						},
						"lldp_med_row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for lldp_med_row_num_service field",
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

func (r *verityEthPortSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityEthPortSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityEthPortSettingsResourceModel
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

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortSettingsTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "MaxAllowedValue", APIField: &ethPortSettingsProps.MaxAllowedValue, TFValue: config.MaxAllowedValue, IsConfigured: configuredAttrs.IsConfigured("max_allowed_value")},
		{FieldName: "MinimumWredThreshold", APIField: &ethPortSettingsProps.MinimumWredThreshold, TFValue: config.MinimumWredThreshold, IsConfigured: configuredAttrs.IsConfigured("minimum_wred_threshold")},
		{FieldName: "MaximumWredThreshold", APIField: &ethPortSettingsProps.MaximumWredThreshold, TFValue: config.MaximumWredThreshold, IsConfigured: configuredAttrs.IsConfigured("maximum_wred_threshold")},
		{FieldName: "WredDropProbability", APIField: &ethPortSettingsProps.WredDropProbability, TFValue: config.WredDropProbability, IsConfigured: configuredAttrs.IsConfigured("wred_drop_probability")},
		{FieldName: "PriorityFlowControlWatchdogDetectTime", APIField: &ethPortSettingsProps.PriorityFlowControlWatchdogDetectTime, TFValue: config.PriorityFlowControlWatchdogDetectTime, IsConfigured: configuredAttrs.IsConfigured("priority_flow_control_watchdog_detect_time")},
		{FieldName: "PriorityFlowControlWatchdogRestoreTime", APIField: &ethPortSettingsProps.PriorityFlowControlWatchdogRestoreTime, TFValue: config.PriorityFlowControlWatchdogRestoreTime, IsConfigured: configuredAttrs.IsConfigured("priority_flow_control_watchdog_restore_time")},
		{FieldName: "MacLimit", APIField: &ethPortSettingsProps.MacLimit, TFValue: config.MacLimit, IsConfigured: configuredAttrs.IsConfigured("mac_limit")},
		{FieldName: "AgingTime", APIField: &ethPortSettingsProps.AgingTime, TFValue: config.AgingTime, IsConfigured: configuredAttrs.IsConfigured("aging_time")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		ethPortSettingsProps.ObjectProperties = &objProps
	}

	// Handle LLDP Med
	if len(plan.LldpMed) > 0 {
		lldpMedItems := make([]openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner, len(plan.LldpMed))
		lldpMedConfigMap := utils.BuildIndexedConfigMap(config.LldpMed)
		for i, item := range plan.LldpMed {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "LldpMedRowNumEnable", APIField: &lldpMedItem.LldpMedRowNumEnable, TFValue: item.LldpMedRowNumEnable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "LldpMedRowNumAdvertisedApplication", APIField: &lldpMedItem.LldpMedRowNumAdvertisedApplicatio, TFValue: item.LldpMedRowNumAdvertisedApplication},
				{FieldName: "LldpMedRowNumService", APIField: &lldpMedItem.LldpMedRowNumService, TFValue: item.LldpMedRowNumService},
				{FieldName: "LldpMedRowNumServiceRefType", APIField: &lldpMedItem.LldpMedRowNumServiceRefType, TFValue: item.LldpMedRowNumServiceRefType},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(item, lldpMedConfigMap, "lldp_med", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "LldpMedRowNumDscpMark", APIField: &lldpMedItem.LldpMedRowNumDscpMark, TFValue: configItem.LldpMedRowNumDscpMark, IsConfigured: cfg.IsFieldConfigured("lldp_med_row_num_dscp_mark")},
				{FieldName: "LldpMedRowNumPriority", APIField: &lldpMedItem.LldpMedRowNumPriority, TFValue: configItem.LldpMedRowNumPriority, IsConfigured: cfg.IsFieldConfigured("lldp_med_row_num_priority")},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &lldpMedItem.Index, TFValue: item.Index},
			})

			lldpMedItems[i] = lldpMedItem
		}
		ethPortSettingsProps.LldpMed = lldpMedItems
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_settings", name, *ethPortSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")

	var minState verityEthPortSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ethPortSettingsData, exists := bulkMgr.GetResourceResponse("eth_port_settings", name); exists {
			state := populateEthPortSettingsState(ctx, minState, ethPortSettingsData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if ethPortSettingsData, exists := r.bulkOpsMgr.GetResourceResponse("eth_port_settings", ethPortSettingsName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached eth port settings data for %s from recent operation", ethPortSettingsName))
			state = populateEthPortSettingsState(ctx, state, ethPortSettingsData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateEthPortSettingsState(ctx, state, ethPortSettingsMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityEthPortSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityEthPortSettingsResourceModel

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

	// Get config for nullable field handling
	var config verityEthPortSettingsResourceModel
	diags = req.Config.Get(ctx, &config)
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

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortSettingsTerraformType, name)

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

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.MaxAllowedValue, state.MaxAllowedValue, configuredAttrs.IsConfigured("max_allowed_value"), func(v *openapi.NullableInt32) { ethPortSettingsProps.MaxAllowedValue = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MinimumWredThreshold, state.MinimumWredThreshold, configuredAttrs.IsConfigured("minimum_wred_threshold"), func(v *openapi.NullableInt32) { ethPortSettingsProps.MinimumWredThreshold = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MaximumWredThreshold, state.MaximumWredThreshold, configuredAttrs.IsConfigured("maximum_wred_threshold"), func(v *openapi.NullableInt32) { ethPortSettingsProps.MaximumWredThreshold = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.WredDropProbability, state.WredDropProbability, configuredAttrs.IsConfigured("wred_drop_probability"), func(v *openapi.NullableInt32) { ethPortSettingsProps.WredDropProbability = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.PriorityFlowControlWatchdogDetectTime, state.PriorityFlowControlWatchdogDetectTime, configuredAttrs.IsConfigured("priority_flow_control_watchdog_detect_time"), func(v *openapi.NullableInt32) { ethPortSettingsProps.PriorityFlowControlWatchdogDetectTime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.PriorityFlowControlWatchdogRestoreTime, state.PriorityFlowControlWatchdogRestoreTime, configuredAttrs.IsConfigured("priority_flow_control_watchdog_restore_time"), func(v *openapi.NullableInt32) { ethPortSettingsProps.PriorityFlowControlWatchdogRestoreTime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MacLimit, state.MacLimit, configuredAttrs.IsConfigured("mac_limit"), func(v *openapi.NullableInt32) { ethPortSettingsProps.MacLimit = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.AgingTime, state.AgingTime, configuredAttrs.IsConfigured("aging_time"), func(v *openapi.NullableInt32) { ethPortSettingsProps.AgingTime = *v }, &hasChanges)

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
	lldpMedConfigMap := utils.BuildIndexedConfigMap(config.LldpMed)

	lldpMedHandler := utils.IndexedItemHandler[verityEthPortSettingsLldpMedModel, openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner]{
		CreateNew: func(item verityEthPortSettingsLldpMedModel) openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &lldpMedItem.Index, TFValue: item.Index},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(item, lldpMedConfigMap, "lldp_med", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "LldpMedRowNumDscpMark", APIField: &lldpMedItem.LldpMedRowNumDscpMark, TFValue: configItem.LldpMedRowNumDscpMark, IsConfigured: cfg.IsFieldConfigured("lldp_med_row_num_dscp_mark")},
				{FieldName: "LldpMedRowNumPriority", APIField: &lldpMedItem.LldpMedRowNumPriority, TFValue: configItem.LldpMedRowNumPriority, IsConfigured: cfg.IsFieldConfigured("lldp_med_row_num_priority")},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "LldpMedRowNumEnable", APIField: &lldpMedItem.LldpMedRowNumEnable, TFValue: item.LldpMedRowNumEnable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "LldpMedRowNumAdvertisedApplication", APIField: &lldpMedItem.LldpMedRowNumAdvertisedApplicatio, TFValue: item.LldpMedRowNumAdvertisedApplication},
				{FieldName: "LldpMedRowNumService", APIField: &lldpMedItem.LldpMedRowNumService, TFValue: item.LldpMedRowNumService},
				{FieldName: "LldpMedRowNumServiceRefType", APIField: &lldpMedItem.LldpMedRowNumServiceRefType, TFValue: item.LldpMedRowNumServiceRefType},
			})

			return lldpMedItem
		},
		UpdateExisting: func(planItem, stateItem verityEthPortSettingsLldpMedModel) (openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner, bool) {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &lldpMedItem.Index, TFValue: planItem.Index},
			})

			hasChanges := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.LldpMedRowNumEnable, stateItem.LldpMedRowNumEnable, func(v *bool) { lldpMedItem.LldpMedRowNumEnable = v }, &hasChanges)

			// Handle string fields (non-ref-type)
			utils.CompareAndSetStringField(planItem.LldpMedRowNumAdvertisedApplication, stateItem.LldpMedRowNumAdvertisedApplication, func(v *string) { lldpMedItem.LldpMedRowNumAdvertisedApplicatio = v }, &hasChanges)

			// Handle nullable int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(planItem, lldpMedConfigMap, "lldp_med", configuredAttrs)
			utils.CompareAndSetNullableInt64Field(configItem.LldpMedRowNumDscpMark, stateItem.LldpMedRowNumDscpMark, cfg.IsFieldConfigured("lldp_med_row_num_dscp_mark"), func(v *openapi.NullableInt32) { lldpMedItem.LldpMedRowNumDscpMark = *v }, &hasChanges)
			utils.CompareAndSetNullableInt64Field(configItem.LldpMedRowNumPriority, stateItem.LldpMedRowNumPriority, cfg.IsFieldConfigured("lldp_med_row_num_priority"), func(v *openapi.NullableInt32) { lldpMedItem.LldpMedRowNumPriority = *v }, &hasChanges)

			// Handle lldp_med_row_num_service and lldp_med_row_num_service_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.LldpMedRowNumService, stateItem.LldpMedRowNumService, planItem.LldpMedRowNumServiceRefType, stateItem.LldpMedRowNumServiceRefType,
				func(v *string) { lldpMedItem.LldpMedRowNumService = v },
				func(v *string) { lldpMedItem.LldpMedRowNumServiceRefType = v },
				"lldp_med_row_num_service", "lldp_med_row_num_service_ref_type_",
				&hasChanges,
				&resp.Diagnostics,
			) {
				return lldpMedItem, false
			}

			return lldpMedItem, hasChanges
		},
		CreateDeleted: func(index int64) openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
			lldpMedItem := openapi.EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &lldpMedItem.Index, TFValue: types.Int64Value(index)},
			})
			return lldpMedItem
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_settings", name, ethPortSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")

	var minState verityEthPortSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ethPortSettingsData, exists := bulkMgr.GetResourceResponse("eth_port_settings", name); exists {
			newState := populateEthPortSettingsState(ctx, minState, ethPortSettingsData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_settings", name, nil, &resp.Diagnostics)
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

func populateEthPortSettingsState(ctx context.Context, state verityEthPortSettingsResourceModel, data map[string]interface{}, mode string) verityEthPortSettingsResourceModel {
	const resourceType = ethPortSettingsResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Int fields
	state.MaxAllowedValue = utils.MapInt64WithMode(data, "max_allowed_value", resourceType, mode)
	state.MinimumWredThreshold = utils.MapInt64WithMode(data, "minimum_wred_threshold", resourceType, mode)
	state.MaximumWredThreshold = utils.MapInt64WithMode(data, "maximum_wred_threshold", resourceType, mode)
	state.WredDropProbability = utils.MapInt64WithMode(data, "wred_drop_probability", resourceType, mode)
	state.PriorityFlowControlWatchdogDetectTime = utils.MapInt64WithMode(data, "priority_flow_control_watchdog_detect_time", resourceType, mode)
	state.PriorityFlowControlWatchdogRestoreTime = utils.MapInt64WithMode(data, "priority_flow_control_watchdog_restore_time", resourceType, mode)
	state.MacLimit = utils.MapInt64WithMode(data, "mac_limit", resourceType, mode)
	state.AgingTime = utils.MapInt64WithMode(data, "aging_time", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.AutoNegotiation = utils.MapBoolWithMode(data, "auto_negotiation", resourceType, mode)
	state.StpEnable = utils.MapBoolWithMode(data, "stp_enable", resourceType, mode)
	state.FastLearningMode = utils.MapBoolWithMode(data, "fast_learning_mode", resourceType, mode)
	state.BpduGuard = utils.MapBoolWithMode(data, "bpdu_guard", resourceType, mode)
	state.BpduFilter = utils.MapBoolWithMode(data, "bpdu_filter", resourceType, mode)
	state.GuardLoop = utils.MapBoolWithMode(data, "guard_loop", resourceType, mode)
	state.PoeEnable = utils.MapBoolWithMode(data, "poe_enable", resourceType, mode)
	state.BspEnable = utils.MapBoolWithMode(data, "bsp_enable", resourceType, mode)
	state.Broadcast = utils.MapBoolWithMode(data, "broadcast", resourceType, mode)
	state.Multicast = utils.MapBoolWithMode(data, "multicast", resourceType, mode)
	state.SingleLink = utils.MapBoolWithMode(data, "single_link", resourceType, mode)
	state.EnableWredTuning = utils.MapBoolWithMode(data, "enable_wred_tuning", resourceType, mode)
	state.EnableEcn = utils.MapBoolWithMode(data, "enable_ecn", resourceType, mode)
	state.EnableWatchdogTuning = utils.MapBoolWithMode(data, "enable_watchdog_tuning", resourceType, mode)
	state.DetectBridgingLoops = utils.MapBoolWithMode(data, "detect_bridging_loops", resourceType, mode)
	state.UnidirectionalLinkDetection = utils.MapBoolWithMode(data, "unidirectional_link_detection", resourceType, mode)
	state.LldpEnable = utils.MapBoolWithMode(data, "lldp_enable", resourceType, mode)
	state.LldpMedEnable = utils.MapBoolWithMode(data, "lldp_med_enable", resourceType, mode)

	// String fields
	state.MaxBitRate = utils.MapStringWithMode(data, "max_bit_rate", resourceType, mode)
	state.DuplexMode = utils.MapStringWithMode(data, "duplex_mode", resourceType, mode)
	state.Priority = utils.MapStringWithMode(data, "priority", resourceType, mode)
	state.AllocatedPower = utils.MapStringWithMode(data, "allocated_power", resourceType, mode)
	state.MaxAllowedUnit = utils.MapStringWithMode(data, "max_allowed_unit", resourceType, mode)
	state.Action = utils.MapStringWithMode(data, "action", resourceType, mode)
	state.Fec = utils.MapStringWithMode(data, "fec", resourceType, mode)
	state.PriorityFlowControlWatchdogAction = utils.MapStringWithMode(data, "priority_flow_control_watchdog_action", resourceType, mode)
	state.PacketQueue = utils.MapStringWithMode(data, "packet_queue", resourceType, mode)
	state.PacketQueueRefType = utils.MapStringWithMode(data, "packet_queue_ref_type_", resourceType, mode)
	state.CliCommands = utils.MapStringWithMode(data, "cli_commands", resourceType, mode)
	state.MacSecurityMode = utils.MapStringWithMode(data, "mac_security_mode", resourceType, mode)
	state.SecurityViolationAction = utils.MapStringWithMode(data, "security_violation_action", resourceType, mode)
	state.AgingType = utils.MapStringWithMode(data, "aging_type", resourceType, mode)
	state.LldpMode = utils.MapStringWithMode(data, "lldp_mode", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityEthPortSettingsObjectPropertiesModel{
				Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityEthPortSettingsObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle lldp_med list block
	if utils.FieldAppliesToMode(resourceType, "lldp_med", mode) {
		if lldpMedData, ok := data["lldp_med"].([]interface{}); ok && len(lldpMedData) > 0 {
			var lldpMedList []verityEthPortSettingsLldpMedModel

			for _, item := range lldpMedData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				lldpMedItem := verityEthPortSettingsLldpMedModel{
					LldpMedRowNumEnable:                utils.MapBoolWithModeNested(itemMap, "lldp_med_row_num_enable", resourceType, "lldp_med.lldp_med_row_num_enable", mode),
					LldpMedRowNumAdvertisedApplication: utils.MapStringWithModeNested(itemMap, "lldp_med_row_num_advertised_applicatio", resourceType, "lldp_med.lldp_med_row_num_advertised_applicatio", mode),
					LldpMedRowNumDscpMark:              utils.MapInt64WithModeNested(itemMap, "lldp_med_row_num_dscp_mark", resourceType, "lldp_med.lldp_med_row_num_dscp_mark", mode),
					LldpMedRowNumPriority:              utils.MapInt64WithModeNested(itemMap, "lldp_med_row_num_priority", resourceType, "lldp_med.lldp_med_row_num_priority", mode),
					LldpMedRowNumService:               utils.MapStringWithModeNested(itemMap, "lldp_med_row_num_service", resourceType, "lldp_med.lldp_med_row_num_service", mode),
					LldpMedRowNumServiceRefType:        utils.MapStringWithModeNested(itemMap, "lldp_med_row_num_service_ref_type_", resourceType, "lldp_med.lldp_med_row_num_service_ref_type_", mode),
					Index:                              utils.MapInt64WithModeNested(itemMap, "index", resourceType, "lldp_med.index", mode),
				}

				lldpMedList = append(lldpMedList, lldpMedItem)
			}

			state.LldpMed = lldpMedList
		} else {
			state.LldpMed = nil
		}
	} else {
		state.LldpMed = nil
	}

	return state
}

func (r *verityEthPortSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityEthPortSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = ethPortSettingsResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"max_bit_rate", "duplex_mode", "priority", "allocated_power",
		"max_allowed_unit", "action", "fec",
		"priority_flow_control_watchdog_action",
		"packet_queue", "packet_queue_ref_type_",
		"cli_commands", "mac_security_mode", "security_violation_action",
		"aging_type", "lldp_mode",
	)

	nullifier.NullifyBools(
		"enable", "auto_negotiation", "stp_enable", "fast_learning_mode",
		"bpdu_guard", "bpdu_filter", "guard_loop", "poe_enable",
		"bsp_enable", "broadcast", "multicast", "single_link",
		"enable_wred_tuning", "enable_ecn", "enable_watchdog_tuning",
		"detect_bridging_loops", "unidirectional_link_detection",
		"lldp_enable", "lldp_med_enable",
	)

	nullifier.NullifyInt64s(
		"max_allowed_value",
		"minimum_wred_threshold", "maximum_wred_threshold", "wred_drop_probability",
		"priority_flow_control_watchdog_detect_time", "priority_flow_control_watchdog_restore_time",
		"mac_limit", "aging_time",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"group"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "lldp_med",
		ItemCount:    len(plan.LldpMed),
		StringFields: []string{"lldp_med_row_num_advertised_applicatio", "lldp_med_row_num_service", "lldp_med_row_num_service_ref_type_"},
		BoolFields:   []string{"lldp_med_row_num_enable"},
		Int64Fields:  []string{"index", "lldp_med_row_num_dscp_mark", "lldp_med_row_num_priority"},
	})

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityEthPortSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityEthPortSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortSettingsTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "max_allowed_value", ConfigVal: config.MaxAllowedValue, StateVal: state.MaxAllowedValue},
			{AttrName: "minimum_wred_threshold", ConfigVal: config.MinimumWredThreshold, StateVal: state.MinimumWredThreshold},
			{AttrName: "maximum_wred_threshold", ConfigVal: config.MaximumWredThreshold, StateVal: state.MaximumWredThreshold},
			{AttrName: "wred_drop_probability", ConfigVal: config.WredDropProbability, StateVal: state.WredDropProbability},
			{AttrName: "priority_flow_control_watchdog_detect_time", ConfigVal: config.PriorityFlowControlWatchdogDetectTime, StateVal: state.PriorityFlowControlWatchdogDetectTime},
			{AttrName: "priority_flow_control_watchdog_restore_time", ConfigVal: config.PriorityFlowControlWatchdogRestoreTime, StateVal: state.PriorityFlowControlWatchdogRestoreTime},
			{AttrName: "mac_limit", ConfigVal: config.MacLimit, StateVal: state.MacLimit},
			{AttrName: "aging_time", ConfigVal: config.AgingTime, StateVal: state.AgingTime},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	for i, configItem := range config.LldpMed {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *verityEthPortSettingsLldpMedModel
		for j := range state.LldpMed {
			if state.LldpMed[j].Index.ValueInt64() == itemIndex {
				stateItem = &state.LldpMed[j]
				break
			}
		}

		if stateItem != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "lldp_med",
				BlockListPath:   "lldp_med",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: itemIndex, AttrName: "lldp_med_row_num_dscp_mark", ConfigVal: configItem.LldpMedRowNumDscpMark, StateVal: stateItem.LldpMedRowNumDscpMark},
					{BlockIndex: itemIndex, AttrName: "lldp_med_row_num_priority", ConfigVal: configItem.LldpMedRowNumPriority, StateVal: stateItem.LldpMedRowNumPriority},
				},
			})
		}
	}
}

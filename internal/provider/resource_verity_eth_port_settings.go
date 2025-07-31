package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	Group                   types.String `tfsdk:"group"`
	OverriddenObject        types.String `tfsdk:"overridden_object"`
	OverriddenObjectRefType types.String `tfsdk:"overridden_object_ref_type_"`
	IsDefault               types.Bool   `tfsdk:"isdefault"`
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
				Default:     booldefault.StaticBool(false),
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
						"overridden_object": schema.StringAttribute{
							Description: "Overridden object.",
							Optional:    true,
						},
						"overridden_object_ref_type_": schema.StringAttribute{
							Description: "Object type for overridden_object field",
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()

	ethPortSettingsReq := openapi.EthportsettingsPutRequestEthPortSettingsValue{}
	ethPortSettingsReq.Name = openapi.PtrString(name)
	if !plan.Enable.IsNull() {
		ethPortSettingsReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.AutoNegotiation.IsNull() {
		ethPortSettingsReq.AutoNegotiation = openapi.PtrBool(plan.AutoNegotiation.ValueBool())
	}
	if !plan.MaxBitRate.IsNull() {
		ethPortSettingsReq.MaxBitRate = openapi.PtrString(plan.MaxBitRate.ValueString())
	}
	if !plan.DuplexMode.IsNull() {
		ethPortSettingsReq.DuplexMode = openapi.PtrString(plan.DuplexMode.ValueString())
	}
	if !plan.StpEnable.IsNull() {
		ethPortSettingsReq.StpEnable = openapi.PtrBool(plan.StpEnable.ValueBool())
	}
	if !plan.FastLearningMode.IsNull() {
		ethPortSettingsReq.FastLearningMode = openapi.PtrBool(plan.FastLearningMode.ValueBool())
	}
	if !plan.BpduGuard.IsNull() {
		ethPortSettingsReq.BpduGuard = openapi.PtrBool(plan.BpduGuard.ValueBool())
	}
	if !plan.BpduFilter.IsNull() {
		ethPortSettingsReq.BpduFilter = openapi.PtrBool(plan.BpduFilter.ValueBool())
	}
	if !plan.GuardLoop.IsNull() {
		ethPortSettingsReq.GuardLoop = openapi.PtrBool(plan.GuardLoop.ValueBool())
	}
	if !plan.PoeEnable.IsNull() {
		ethPortSettingsReq.PoeEnable = openapi.PtrBool(plan.PoeEnable.ValueBool())
	}
	if !plan.Priority.IsNull() {
		ethPortSettingsReq.Priority = openapi.PtrString(plan.Priority.ValueString())
	}
	if !plan.AllocatedPower.IsNull() {
		ethPortSettingsReq.AllocatedPower = openapi.PtrString(plan.AllocatedPower.ValueString())
	}
	if !plan.BspEnable.IsNull() {
		ethPortSettingsReq.BspEnable = openapi.PtrBool(plan.BspEnable.ValueBool())
	}
	if !plan.Broadcast.IsNull() {
		ethPortSettingsReq.Broadcast = openapi.PtrBool(plan.Broadcast.ValueBool())
	}
	if !plan.Multicast.IsNull() {
		ethPortSettingsReq.Multicast = openapi.PtrBool(plan.Multicast.ValueBool())
	}
	if !plan.MaxAllowedValue.IsNull() {
		ethPortSettingsReq.MaxAllowedValue = openapi.PtrInt32(int32(plan.MaxAllowedValue.ValueInt64()))
	}
	if !plan.MaxAllowedUnit.IsNull() {
		ethPortSettingsReq.MaxAllowedUnit = openapi.PtrString(plan.MaxAllowedUnit.ValueString())
	}
	if !plan.Action.IsNull() {
		ethPortSettingsReq.Action = openapi.PtrString(plan.Action.ValueString())
	}
	if !plan.Fec.IsNull() {
		ethPortSettingsReq.Fec = openapi.PtrString(plan.Fec.ValueString())
	}
	if !plan.SingleLink.IsNull() {
		ethPortSettingsReq.SingleLink = openapi.PtrBool(plan.SingleLink.ValueBool())
	}
	if !plan.MinimumWredThreshold.IsNull() {
		ethPortSettingsReq.MinimumWredThreshold = openapi.PtrInt32(int32(plan.MinimumWredThreshold.ValueInt64()))
	}
	if !plan.MaximumWredThreshold.IsNull() {
		ethPortSettingsReq.MaximumWredThreshold = openapi.PtrInt32(int32(plan.MaximumWredThreshold.ValueInt64()))
	}
	if !plan.WredDropProbability.IsNull() {
		ethPortSettingsReq.WredDropProbability = openapi.PtrInt32(int32(plan.WredDropProbability.ValueInt64()))
	}
	if !plan.PriorityFlowControlWatchdogAction.IsNull() {
		ethPortSettingsReq.PriorityFlowControlWatchdogAction = openapi.PtrString(plan.PriorityFlowControlWatchdogAction.ValueString())
	}
	if !plan.PriorityFlowControlWatchdogDetectTime.IsNull() {
		ethPortSettingsReq.PriorityFlowControlWatchdogDetectTime = openapi.PtrInt32(int32(plan.PriorityFlowControlWatchdogDetectTime.ValueInt64()))
	}
	if !plan.PriorityFlowControlWatchdogRestoreTime.IsNull() {
		ethPortSettingsReq.PriorityFlowControlWatchdogRestoreTime = openapi.PtrInt32(int32(plan.PriorityFlowControlWatchdogRestoreTime.ValueInt64()))
	}
	if !plan.PacketQueue.IsNull() {
		ethPortSettingsReq.PacketQueue = openapi.PtrString(plan.PacketQueue.ValueString())
	}
	if !plan.PacketQueueRefType.IsNull() {
		ethPortSettingsReq.PacketQueueRefType = openapi.PtrString(plan.PacketQueueRefType.ValueString())
	}
	if !plan.EnableWredTuning.IsNull() {
		ethPortSettingsReq.EnableWredTuning = openapi.PtrBool(plan.EnableWredTuning.ValueBool())
	}
	if !plan.EnableEcn.IsNull() {
		ethPortSettingsReq.EnableEcn = openapi.PtrBool(plan.EnableEcn.ValueBool())
	}
	if !plan.EnableWatchdogTuning.IsNull() {
		ethPortSettingsReq.EnableWatchdogTuning = openapi.PtrBool(plan.EnableWatchdogTuning.ValueBool())
	}
	if !plan.CliCommands.IsNull() {
		ethPortSettingsReq.CliCommands = openapi.PtrString(plan.CliCommands.ValueString())
	}
	if !plan.DetectBridgingLoops.IsNull() {
		ethPortSettingsReq.DetectBridgingLoops = openapi.PtrBool(plan.DetectBridgingLoops.ValueBool())
	}
	if !plan.UnidirectionalLinkDetection.IsNull() {
		ethPortSettingsReq.UnidirectionalLinkDetection = openapi.PtrBool(plan.UnidirectionalLinkDetection.ValueBool())
	}
	if !plan.MacSecurityMode.IsNull() {
		ethPortSettingsReq.MacSecurityMode = openapi.PtrString(plan.MacSecurityMode.ValueString())
	}
	if !plan.MacLimit.IsNull() {
		ethPortSettingsReq.MacLimit = openapi.PtrInt32(int32(plan.MacLimit.ValueInt64()))
	}
	if !plan.SecurityViolationAction.IsNull() {
		ethPortSettingsReq.SecurityViolationAction = openapi.PtrString(plan.SecurityViolationAction.ValueString())
	}
	if !plan.AgingType.IsNull() {
		ethPortSettingsReq.AgingType = openapi.PtrString(plan.AgingType.ValueString())
	}
	if !plan.AgingTime.IsNull() {
		ethPortSettingsReq.AgingTime = openapi.PtrInt32(int32(plan.AgingTime.ValueInt64()))
	}
	if !plan.LldpEnable.IsNull() {
		ethPortSettingsReq.LldpEnable = openapi.PtrBool(plan.LldpEnable.ValueBool())
	}
	if !plan.LldpMode.IsNull() {
		ethPortSettingsReq.LldpMode = openapi.PtrString(plan.LldpMode.ValueString())
	}
	if !plan.LldpMedEnable.IsNull() {
		ethPortSettingsReq.LldpMedEnable = openapi.PtrBool(plan.LldpMedEnable.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		if !op.OverriddenObject.IsNull() {
			objProps.OverriddenObject = openapi.PtrString(op.OverriddenObject.ValueString())
		}
		if !op.OverriddenObjectRefType.IsNull() {
			objProps.OverriddenObjectRefType = openapi.PtrString(op.OverriddenObjectRefType.ValueString())
		}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		}
		ethPortSettingsReq.ObjectProperties = &objProps
	} else {
		ethPortSettingsReq.ObjectProperties = nil
	}

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

		ethPortSettingsReq.LldpMed = lldpMedItems
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPut(ctx, "eth_port_settings", name, ethPortSettingsReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings creation operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Eth Port Settings %s", name))...,
		)
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	settingsName := state.Name.ValueString()

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("eth_port_settings") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port settings %s verification - trusting recent successful API operation", settingsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent eth port settings operations found, performing normal verification for %s", settingsName))

	type EthPortSettingsResponse struct {
		EthPortSettings map[string]map[string]interface{} `json:"eth_port_settings"`
	}

	var result EthPortSettingsResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		settingsData, fetchErr := getCachedResponse(ctx, provCtx, "eth_port_settings", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Ethernet port settings")
			apiResp, err := provCtx.client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading eth port settings: %v", err)
			}
			defer apiResp.Body.Close()

			var res EthPortSettingsResponse
			if err := json.NewDecoder(apiResp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("error decoding eth port settings response: %v", err)
			}
			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Ethernet port settings from API", len(res.EthPortSettings)))
			return res, nil
		})

		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch eth port settings on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = settingsData.(EthPortSettingsResponse)
		break
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Eth Port Settings %s", settingsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Ethernet port settings with ID: %s", settingsName))
	var settings map[string]interface{}
	exists := false

	if data, ok := result.EthPortSettings[settingsName]; ok {
		settings = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings directly by ID: %s", settingsName))
	} else {
		var nameStr types.String
		diags := req.State.GetAttribute(ctx, path.Root("name"), &nameStr)
		if !diags.HasError() && !nameStr.IsNull() {
			settingsNameFromAttr := nameStr.ValueString()
			tflog.Debug(ctx, fmt.Sprintf("Settings not found by ID, trying name attribute: %s", settingsNameFromAttr))
			if data, ok := result.EthPortSettings[settingsNameFromAttr]; ok {
				settings = data
				settingsName = settingsNameFromAttr
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings by name attribute: %s", settingsNameFromAttr))
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Settings not found directly by ID '%s', searching through all settings", settingsName))
		for apiName, s := range result.EthPortSettings {
			if name, ok := s["name"].(string); ok && name == settingsName {
				settings = s
				settingsName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Ethernet port settings with ID '%s' not found in API response", settingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	if op, ok := settings["object_properties"].(map[string]interface{}); ok {
		objProp := verityEthPortSettingsObjectPropertiesModel{}

		if group, ok := op["group"].(string); ok {
			objProp.Group = types.StringValue(group)
		} else {
			objProp.Group = types.StringNull()
		}

		if overriddenObject, ok := op["overridden_object"].(string); ok {
			objProp.OverriddenObject = types.StringValue(overriddenObject)
		} else {
			objProp.OverriddenObject = types.StringNull()
		}

		if overriddenObjectRefType, ok := op["overridden_object_ref_type_"].(string); ok {
			objProp.OverriddenObjectRefType = types.StringValue(overriddenObjectRefType)
		} else {
			objProp.OverriddenObjectRefType = types.StringNull()
		}

		if isDefault, ok := op["isdefault"].(bool); ok {
			objProp.IsDefault = types.BoolValue(isDefault)
		} else {
			objProp.IsDefault = types.BoolNull()
		}

		state.ObjectProperties = []verityEthPortSettingsObjectPropertiesModel{objProp}
	} else {
		state.ObjectProperties = nil
	}

	if val, ok := settings["name"]; ok {
		state.Name = types.StringValue(fmt.Sprintf("%v", val))
	}
	if val, ok := settings["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else {
		state.Enable = types.BoolNull()
	}
	if val, ok := settings["auto_negotiation"].(bool); ok {
		state.AutoNegotiation = types.BoolValue(val)
	} else {
		state.AutoNegotiation = types.BoolNull()
	}
	if val, ok := settings["max_bit_rate"].(string); ok {
		state.MaxBitRate = types.StringValue(val)
	} else {
		state.MaxBitRate = types.StringNull()
	}
	if val, ok := settings["duplex_mode"].(string); ok {
		state.DuplexMode = types.StringValue(val)
	} else {
		state.DuplexMode = types.StringNull()
	}
	if val, ok := settings["stp_enable"].(bool); ok {
		state.StpEnable = types.BoolValue(val)
	} else {
		state.StpEnable = types.BoolNull()
	}
	if val, ok := settings["fast_learning_mode"].(bool); ok {
		state.FastLearningMode = types.BoolValue(val)
	} else {
		state.FastLearningMode = types.BoolNull()
	}
	if val, ok := settings["bpdu_guard"].(bool); ok {
		state.BpduGuard = types.BoolValue(val)
	} else {
		state.BpduGuard = types.BoolNull()
	}
	if val, ok := settings["bpdu_filter"].(bool); ok {
		state.BpduFilter = types.BoolValue(val)
	} else {
		state.BpduFilter = types.BoolNull()
	}
	if val, ok := settings["guard_loop"].(bool); ok {
		state.GuardLoop = types.BoolValue(val)
	} else {
		state.GuardLoop = types.BoolNull()
	}
	if val, ok := settings["poe_enable"].(bool); ok {
		state.PoeEnable = types.BoolValue(val)
	} else {
		state.PoeEnable = types.BoolNull()
	}
	if val, ok := settings["priority"].(string); ok {
		state.Priority = types.StringValue(val)
	} else {
		state.Priority = types.StringNull()
	}
	if val, ok := settings["allocated_power"].(string); ok {
		state.AllocatedPower = types.StringValue(val)
	} else {
		state.AllocatedPower = types.StringNull()
	}
	if val, ok := settings["bsp_enable"].(bool); ok {
		state.BspEnable = types.BoolValue(val)
	} else {
		state.BspEnable = types.BoolNull()
	}
	if val, ok := settings["broadcast"].(bool); ok {
		state.Broadcast = types.BoolValue(val)
	} else {
		state.Broadcast = types.BoolNull()
	}
	if val, ok := settings["multicast"].(bool); ok {
		state.Multicast = types.BoolValue(val)
	} else {
		state.Multicast = types.BoolNull()
	}
	if val, ok := settings["max_allowed_value"]; ok {
		switch v := val.(type) {
		case float64:
			state.MaxAllowedValue = types.Int64Value(int64(v))
		case int:
			state.MaxAllowedValue = types.Int64Value(int64(v))
		default:
			state.MaxAllowedValue = types.Int64Null()
		}
	} else {
		state.MaxAllowedValue = types.Int64Null()
	}
	if val, ok := settings["max_allowed_unit"].(string); ok {
		state.MaxAllowedUnit = types.StringValue(val)
	} else {
		state.MaxAllowedUnit = types.StringNull()
	}
	if val, ok := settings["action"].(string); ok {
		state.Action = types.StringValue(val)
	} else {
		state.Action = types.StringNull()
	}
	if val, ok := settings["fec"].(string); ok {
		state.Fec = types.StringValue(val)
	} else {
		state.Fec = types.StringNull()
	}
	if val, ok := settings["single_link"].(bool); ok {
		state.SingleLink = types.BoolValue(val)
	} else {
		state.SingleLink = types.BoolNull()
	}

	if val, ok := settings["minimum_wred_threshold"]; ok {
		switch v := val.(type) {
		case float64:
			state.MinimumWredThreshold = types.Int64Value(int64(v))
		case int:
			state.MinimumWredThreshold = types.Int64Value(int64(v))
		default:
			state.MinimumWredThreshold = types.Int64Null()
		}
	} else {
		state.MinimumWredThreshold = types.Int64Null()
	}
	if val, ok := settings["maximum_wred_threshold"]; ok {
		switch v := val.(type) {
		case float64:
			state.MaximumWredThreshold = types.Int64Value(int64(v))
		case int:
			state.MaximumWredThreshold = types.Int64Value(int64(v))
		default:
			state.MaximumWredThreshold = types.Int64Null()
		}
	} else {
		state.MaximumWredThreshold = types.Int64Null()
	}
	if val, ok := settings["wred_drop_probability"]; ok {
		switch v := val.(type) {
		case float64:
			state.WredDropProbability = types.Int64Value(int64(v))
		case int:
			state.WredDropProbability = types.Int64Value(int64(v))
		default:
			state.WredDropProbability = types.Int64Null()
		}
	} else {
		state.WredDropProbability = types.Int64Null()
	}
	if val, ok := settings["priority_flow_control_watchdog_action"].(string); ok {
		state.PriorityFlowControlWatchdogAction = types.StringValue(val)
	} else {
		state.PriorityFlowControlWatchdogAction = types.StringNull()
	}
	if val, ok := settings["priority_flow_control_watchdog_detect_time"]; ok {
		switch v := val.(type) {
		case float64:
			state.PriorityFlowControlWatchdogDetectTime = types.Int64Value(int64(v))
		case int:
			state.PriorityFlowControlWatchdogDetectTime = types.Int64Value(int64(v))
		default:
			state.PriorityFlowControlWatchdogDetectTime = types.Int64Null()
		}
	} else {
		state.PriorityFlowControlWatchdogDetectTime = types.Int64Null()
	}
	if val, ok := settings["priority_flow_control_watchdog_restore_time"]; ok {
		switch v := val.(type) {
		case float64:
			state.PriorityFlowControlWatchdogRestoreTime = types.Int64Value(int64(v))
		case int:
			state.PriorityFlowControlWatchdogRestoreTime = types.Int64Value(int64(v))
		default:
			state.PriorityFlowControlWatchdogRestoreTime = types.Int64Null()
		}
	} else {
		state.PriorityFlowControlWatchdogRestoreTime = types.Int64Null()
	}
	if val, ok := settings["packet_queue"].(string); ok {
		state.PacketQueue = types.StringValue(val)
	} else {
		state.PacketQueue = types.StringNull()
	}
	if val, ok := settings["packet_queue_ref_type_"].(string); ok {
		state.PacketQueueRefType = types.StringValue(val)
	} else {
		state.PacketQueueRefType = types.StringNull()
	}
	if val, ok := settings["enable_wred_tuning"].(bool); ok {
		state.EnableWredTuning = types.BoolValue(val)
	} else {
		state.EnableWredTuning = types.BoolNull()
	}
	if val, ok := settings["enable_ecn"].(bool); ok {
		state.EnableEcn = types.BoolValue(val)
	} else {
		state.EnableEcn = types.BoolNull()
	}
	if val, ok := settings["enable_watchdog_tuning"].(bool); ok {
		state.EnableWatchdogTuning = types.BoolValue(val)
	} else {
		state.EnableWatchdogTuning = types.BoolNull()
	}
	if val, ok := settings["cli_commands"].(string); ok {
		state.CliCommands = types.StringValue(val)
	} else {
		state.CliCommands = types.StringNull()
	}
	if val, ok := settings["detect_bridging_loops"].(bool); ok {
		state.DetectBridgingLoops = types.BoolValue(val)
	} else {
		state.DetectBridgingLoops = types.BoolNull()
	}
	if val, ok := settings["unidirectional_link_detection"].(bool); ok {
		state.UnidirectionalLinkDetection = types.BoolValue(val)
	} else {
		state.UnidirectionalLinkDetection = types.BoolNull()
	}
	if val, ok := settings["mac_security_mode"].(string); ok {
		state.MacSecurityMode = types.StringValue(val)
	} else {
		state.MacSecurityMode = types.StringNull()
	}
	if val, ok := settings["mac_limit"]; ok {
		switch v := val.(type) {
		case float64:
			state.MacLimit = types.Int64Value(int64(v))
		case int:
			state.MacLimit = types.Int64Value(int64(v))
		default:
			state.MacLimit = types.Int64Null()
		}
	} else {
		state.MacLimit = types.Int64Null()
	}
	if val, ok := settings["security_violation_action"].(string); ok {
		state.SecurityViolationAction = types.StringValue(val)
	} else {
		state.SecurityViolationAction = types.StringNull()
	}
	if val, ok := settings["aging_type"].(string); ok {
		state.AgingType = types.StringValue(val)
	} else {
		state.AgingType = types.StringNull()
	}
	if val, ok := settings["aging_time"]; ok {
		switch v := val.(type) {
		case float64:
			state.AgingTime = types.Int64Value(int64(v))
		case int:
			state.AgingTime = types.Int64Value(int64(v))
		default:
			state.AgingTime = types.Int64Null()
		}
	} else {
		state.AgingTime = types.Int64Null()
	}
	if val, ok := settings["lldp_enable"].(bool); ok {
		state.LldpEnable = types.BoolValue(val)
	} else {
		state.LldpEnable = types.BoolNull()
	}
	if val, ok := settings["lldp_mode"].(string); ok {
		state.LldpMode = types.StringValue(val)
	} else {
		state.LldpMode = types.StringNull()
	}
	if val, ok := settings["lldp_med_enable"].(bool); ok {
		state.LldpMedEnable = types.BoolValue(val)
	} else {
		state.LldpMedEnable = types.BoolNull()
	}

	if lldpMedData, ok := settings["lldp_med"].([]interface{}); ok {
		var lldpMedList []verityEthPortSettingsLldpMedModel

		for _, item := range lldpMedData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				lldpMedItem := verityEthPortSettingsLldpMedModel{}

				if enable, ok := itemMap["lldp_med_row_num_enable"].(bool); ok {
					lldpMedItem.LldpMedRowNumEnable = types.BoolValue(enable)
				} else {
					lldpMedItem.LldpMedRowNumEnable = types.BoolNull()
				}

				if app, ok := itemMap["lldp_med_row_num_advertised_applicatio"].(string); ok {
					lldpMedItem.LldpMedRowNumAdvertisedApplication = types.StringValue(app)
				} else {
					lldpMedItem.LldpMedRowNumAdvertisedApplication = types.StringNull()
				}

				if dscpMark, ok := itemMap["lldp_med_row_num_dscp_mark"]; ok {
					switch v := dscpMark.(type) {
					case float64:
						lldpMedItem.LldpMedRowNumDscpMark = types.Int64Value(int64(v))
					case int:
						lldpMedItem.LldpMedRowNumDscpMark = types.Int64Value(int64(v))
					default:
						lldpMedItem.LldpMedRowNumDscpMark = types.Int64Null()
					}
				} else {
					lldpMedItem.LldpMedRowNumDscpMark = types.Int64Null()
				}

				if priority, ok := itemMap["lldp_med_row_num_priority"]; ok {
					switch v := priority.(type) {
					case float64:
						lldpMedItem.LldpMedRowNumPriority = types.Int64Value(int64(v))
					case int:
						lldpMedItem.LldpMedRowNumPriority = types.Int64Value(int64(v))
					default:
						lldpMedItem.LldpMedRowNumPriority = types.Int64Null()
					}
				} else {
					lldpMedItem.LldpMedRowNumPriority = types.Int64Null()
				}

				if service, ok := itemMap["lldp_med_row_num_service"].(string); ok {
					lldpMedItem.LldpMedRowNumService = types.StringValue(service)
				} else {
					lldpMedItem.LldpMedRowNumService = types.StringNull()
				}

				if serviceRefType, ok := itemMap["lldp_med_row_num_service_ref_type_"].(string); ok {
					lldpMedItem.LldpMedRowNumServiceRefType = types.StringValue(serviceRefType)
				} else {
					lldpMedItem.LldpMedRowNumServiceRefType = types.StringNull()
				}

				if index, ok := itemMap["index"]; ok {
					switch v := index.(type) {
					case float64:
						lldpMedItem.Index = types.Int64Value(int64(v))
					case int:
						lldpMedItem.Index = types.Int64Value(int64(v))
					default:
						lldpMedItem.Index = types.Int64Null()
					}
				} else {
					lldpMedItem.Index = types.Int64Null()
				}

				lldpMedList = append(lldpMedList, lldpMedItem)
			}
		}

		state.LldpMed = lldpMedList
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

	ethPortSettingsReq := openapi.EthportsettingsPutRequestEthPortSettingsValue{}
	hasChanges := false

	objectPropertiesChanged := false

	if (len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) == 0) ||
		(len(plan.ObjectProperties) == 0 && len(state.ObjectProperties) > 0) {
		objectPropertiesChanged = true
	} else if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		if !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].OverriddenObject.Equal(state.ObjectProperties[0].OverriddenObject) ||
			!plan.ObjectProperties[0].OverriddenObjectRefType.Equal(state.ObjectProperties[0].OverriddenObjectRefType) ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) {
			objectPropertiesChanged = true
		}
	}

	if objectPropertiesChanged {
		if len(plan.ObjectProperties) > 0 {
			objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			if !plan.ObjectProperties[0].OverriddenObject.IsNull() {
				objProps.OverriddenObject = openapi.PtrString(plan.ObjectProperties[0].OverriddenObject.ValueString())
			}
			if !plan.ObjectProperties[0].OverriddenObjectRefType.IsNull() {
				objProps.OverriddenObjectRefType = openapi.PtrString(plan.ObjectProperties[0].OverriddenObjectRefType.ValueString())
			}
			if !plan.ObjectProperties[0].IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(plan.ObjectProperties[0].IsDefault.ValueBool())
			}
			ethPortSettingsReq.ObjectProperties = &objProps
		} else {
			ethPortSettingsReq.ObjectProperties = nil
		}
		hasChanges = true
	}

	boolFields := []struct {
		planValue  types.Bool
		stateValue types.Bool
		setter     func(val bool)
	}{
		{
			planValue:  plan.Enable,
			stateValue: state.Enable,
			setter: func(val bool) {
				ethPortSettingsReq.Enable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.AutoNegotiation,
			stateValue: state.AutoNegotiation,
			setter: func(val bool) {
				ethPortSettingsReq.AutoNegotiation = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.StpEnable,
			stateValue: state.StpEnable,
			setter: func(val bool) {
				ethPortSettingsReq.StpEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.FastLearningMode,
			stateValue: state.FastLearningMode,
			setter: func(val bool) {
				ethPortSettingsReq.FastLearningMode = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BpduGuard,
			stateValue: state.BpduGuard,
			setter: func(val bool) {
				ethPortSettingsReq.BpduGuard = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BpduFilter,
			stateValue: state.BpduFilter,
			setter: func(val bool) {
				ethPortSettingsReq.BpduFilter = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.GuardLoop,
			stateValue: state.GuardLoop,
			setter: func(val bool) {
				ethPortSettingsReq.GuardLoop = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.PoeEnable,
			stateValue: state.PoeEnable,
			setter: func(val bool) {
				ethPortSettingsReq.PoeEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BspEnable,
			stateValue: state.BspEnable,
			setter: func(val bool) {
				ethPortSettingsReq.BspEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.Broadcast,
			stateValue: state.Broadcast,
			setter: func(val bool) {
				ethPortSettingsReq.Broadcast = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.Multicast,
			stateValue: state.Multicast,
			setter: func(val bool) {
				ethPortSettingsReq.Multicast = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.SingleLink,
			stateValue: state.SingleLink,
			setter: func(val bool) {
				ethPortSettingsReq.SingleLink = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.EnableWredTuning,
			stateValue: state.EnableWredTuning,
			setter: func(val bool) {
				ethPortSettingsReq.EnableWredTuning = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.EnableEcn,
			stateValue: state.EnableEcn,
			setter: func(val bool) {
				ethPortSettingsReq.EnableEcn = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.EnableWatchdogTuning,
			stateValue: state.EnableWatchdogTuning,
			setter: func(val bool) {
				ethPortSettingsReq.EnableWatchdogTuning = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.DetectBridgingLoops,
			stateValue: state.DetectBridgingLoops,
			setter: func(val bool) {
				ethPortSettingsReq.DetectBridgingLoops = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.UnidirectionalLinkDetection,
			stateValue: state.UnidirectionalLinkDetection,
			setter: func(val bool) {
				ethPortSettingsReq.UnidirectionalLinkDetection = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.LldpEnable,
			stateValue: state.LldpEnable,
			setter: func(val bool) {
				ethPortSettingsReq.LldpEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.LldpMedEnable,
			stateValue: state.LldpMedEnable,
			setter: func(val bool) {
				ethPortSettingsReq.LldpMedEnable = openapi.PtrBool(val)
			},
		},
	}

	for _, field := range boolFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(field.planValue.ValueBool())
			hasChanges = true
		}
	}

	stringFields := []struct {
		planValue  types.String
		stateValue types.String
		setter     func(val string)
	}{
		{
			planValue:  plan.MaxBitRate,
			stateValue: state.MaxBitRate,
			setter: func(val string) {
				ethPortSettingsReq.MaxBitRate = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.DuplexMode,
			stateValue: state.DuplexMode,
			setter: func(val string) {
				ethPortSettingsReq.DuplexMode = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Priority,
			stateValue: state.Priority,
			setter: func(val string) {
				ethPortSettingsReq.Priority = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.AllocatedPower,
			stateValue: state.AllocatedPower,
			setter: func(val string) {
				ethPortSettingsReq.AllocatedPower = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.MaxAllowedUnit,
			stateValue: state.MaxAllowedUnit,
			setter: func(val string) {
				ethPortSettingsReq.MaxAllowedUnit = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Action,
			stateValue: state.Action,
			setter: func(val string) {
				ethPortSettingsReq.Action = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Fec,
			stateValue: state.Fec,
			setter: func(val string) {
				ethPortSettingsReq.Fec = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.PriorityFlowControlWatchdogAction,
			stateValue: state.PriorityFlowControlWatchdogAction,
			setter: func(val string) {
				ethPortSettingsReq.PriorityFlowControlWatchdogAction = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.PacketQueue,
			stateValue: state.PacketQueue,
			setter: func(val string) {
				ethPortSettingsReq.PacketQueue = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.PacketQueueRefType,
			stateValue: state.PacketQueueRefType,
			setter: func(val string) {
				ethPortSettingsReq.PacketQueueRefType = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.CliCommands,
			stateValue: state.CliCommands,
			setter: func(val string) {
				ethPortSettingsReq.CliCommands = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.MacSecurityMode,
			stateValue: state.MacSecurityMode,
			setter: func(val string) {
				ethPortSettingsReq.MacSecurityMode = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.SecurityViolationAction,
			stateValue: state.SecurityViolationAction,
			setter: func(val string) {
				ethPortSettingsReq.SecurityViolationAction = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.AgingType,
			stateValue: state.AgingType,
			setter: func(val string) {
				ethPortSettingsReq.AgingType = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.LldpMode,
			stateValue: state.LldpMode,
			setter: func(val string) {
				ethPortSettingsReq.LldpMode = openapi.PtrString(val)
			},
		},
	}

	for _, field := range stringFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(field.planValue.ValueString())
			hasChanges = true
		}
	}

	intFields := []struct {
		planValue  types.Int64
		stateValue types.Int64
		setter     func(val int32)
	}{
		{
			planValue:  plan.MaxAllowedValue,
			stateValue: state.MaxAllowedValue,
			setter: func(val int32) {
				ethPortSettingsReq.MaxAllowedValue = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.MinimumWredThreshold,
			stateValue: state.MinimumWredThreshold,
			setter: func(val int32) {
				ethPortSettingsReq.MinimumWredThreshold = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.MaximumWredThreshold,
			stateValue: state.MaximumWredThreshold,
			setter: func(val int32) {
				ethPortSettingsReq.MaximumWredThreshold = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.WredDropProbability,
			stateValue: state.WredDropProbability,
			setter: func(val int32) {
				ethPortSettingsReq.WredDropProbability = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.PriorityFlowControlWatchdogDetectTime,
			stateValue: state.PriorityFlowControlWatchdogDetectTime,
			setter: func(val int32) {
				ethPortSettingsReq.PriorityFlowControlWatchdogDetectTime = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.PriorityFlowControlWatchdogRestoreTime,
			stateValue: state.PriorityFlowControlWatchdogRestoreTime,
			setter: func(val int32) {
				ethPortSettingsReq.PriorityFlowControlWatchdogRestoreTime = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.MacLimit,
			stateValue: state.MacLimit,
			setter: func(val int32) {
				ethPortSettingsReq.MacLimit = openapi.PtrInt32(val)
			},
		},
		{
			planValue:  plan.AgingTime,
			stateValue: state.AgingTime,
			setter: func(val int32) {
				ethPortSettingsReq.AgingTime = openapi.PtrInt32(val)
			},
		},
	}

	for _, field := range intFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(int32(field.planValue.ValueInt64()))
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	provCtx := r.provCtx
	bulkMgr := provCtx.bulkOpsMgr
	resourceName := plan.Name.ValueString()

	operationID := bulkMgr.AddPatch(ctx, "eth_port_settings", resourceName, ethPortSettingsReq)
	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Eth Port Settings %s", resourceName))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s update operation completed successfully", resourceName))
	clearCache(ctx, provCtx, "eth_port_settings")
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := state.Name.ValueString()
	bulkMgr := r.provCtx.bulkOpsMgr
	operationID := bulkMgr.AddDelete(ctx, "eth_port_settings", name)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Eth Port Settings %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

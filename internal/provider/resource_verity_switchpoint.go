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
	_ resource.Resource                = &veritySwitchpointResource{}
	_ resource.ResourceWithConfigure   = &veritySwitchpointResource{}
	_ resource.ResourceWithImportState = &veritySwitchpointResource{}
	_ resource.ResourceWithModifyPlan  = &veritySwitchpointResource{}
)

func NewVeritySwitchpointResource() resource.Resource {
	return &veritySwitchpointResource{}
}

type veritySwitchpointResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type veritySwitchpointResourceModel struct {
	Name                             types.String                             `tfsdk:"name"`
	DeviceSerialNumber               types.String                             `tfsdk:"device_serial_number"`
	ConnectedBundle                  types.String                             `tfsdk:"connected_bundle"`
	ConnectedBundleRefType           types.String                             `tfsdk:"connected_bundle_ref_type_"`
	ReadOnlyMode                     types.Bool                               `tfsdk:"read_only_mode"`
	Locked                           types.Bool                               `tfsdk:"locked"`
	DisabledPorts                    types.String                             `tfsdk:"disabled_ports"`
	OutOfBandManagement              types.Bool                               `tfsdk:"out_of_band_management"`
	Type                             types.String                             `tfsdk:"type"`
	SuperPod                         types.String                             `tfsdk:"super_pod"`
	Pod                              types.String                             `tfsdk:"pod"`
	Rack                             types.String                             `tfsdk:"rack"`
	SwitchRouterIdIpMask             types.String                             `tfsdk:"switch_router_id_ip_mask"`
	SwitchRouterIdIpMaskAutoAssigned types.Bool                               `tfsdk:"switch_router_id_ip_mask_auto_assigned_"`
	SwitchVtepIdIpMask               types.String                             `tfsdk:"switch_vtep_id_ip_mask"`
	SwitchVtepIdIpMaskAutoAssigned   types.Bool                               `tfsdk:"switch_vtep_id_ip_mask_auto_assigned_"`
	BgpAsNumber                      types.Int64                              `tfsdk:"bgp_as_number"`
	BgpAsNumberAutoAssigned          types.Bool                               `tfsdk:"bgp_as_number_auto_assigned_"`
	Badges                           []veritySwitchpointBadgeModel            `tfsdk:"badges"`
	Children                         []veritySwitchpointChildModel            `tfsdk:"children"`
	TrafficMirrors                   []veritySwitchpointTrafficMirrorModel    `tfsdk:"traffic_mirrors"`
	Eths                             []veritySwitchpointEthModel              `tfsdk:"eths"`
	ObjectProperties                 []veritySwitchpointObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySwitchpointBadgeModel struct {
	Badge        types.String `tfsdk:"badge"`
	BadgeRefType types.String `tfsdk:"badge_ref_type_"`
	Index        types.Int64  `tfsdk:"index"`
}

type veritySwitchpointChildModel struct {
	ChildNumEndpoint        types.String `tfsdk:"child_num_endpoint"`
	ChildNumEndpointRefType types.String `tfsdk:"child_num_endpoint_ref_type_"`
	ChildNumDevice          types.String `tfsdk:"child_num_device"`
	Index                   types.Int64  `tfsdk:"index"`
}

type veritySwitchpointTrafficMirrorModel struct {
	TrafficMirrorNumEnable             types.Bool   `tfsdk:"traffic_mirror_num_enable"`
	TrafficMirrorNumSourcePort         types.String `tfsdk:"traffic_mirror_num_source_port"`
	TrafficMirrorNumSourceLagIndicator types.Bool   `tfsdk:"traffic_mirror_num_source_lag_indicator"`
	TrafficMirrorNumDestinationPort    types.String `tfsdk:"traffic_mirror_num_destination_port"`
	TrafficMirrorNumInboundTraffic     types.Bool   `tfsdk:"traffic_mirror_num_inbound_traffic"`
	TrafficMirrorNumOutboundTraffic    types.Bool   `tfsdk:"traffic_mirror_num_outbound_traffic"`
}

type veritySwitchpointEthModel struct {
	Breakout types.String `tfsdk:"breakout"`
	Index    types.Int64  `tfsdk:"index"`
}

type veritySwitchpointObjectPropertiesEthModel struct {
	EthNumIcon  types.String `tfsdk:"eth_num_icon"`
	EthNumLabel types.String `tfsdk:"eth_num_label"`
	Index       types.Int64  `tfsdk:"index"`
}

type veritySwitchpointObjectPropertiesModel struct {
	UserNotes                     types.String                                `tfsdk:"user_notes"`
	ExpectedParentEndpoint        types.String                                `tfsdk:"expected_parent_endpoint"`
	ExpectedParentEndpointRefType types.String                                `tfsdk:"expected_parent_endpoint_ref_type_"`
	NumberOfMultipoints           types.Int64                                 `tfsdk:"number_of_multipoints"`
	Aggregate                     types.Bool                                  `tfsdk:"aggregate"`
	IsHost                        types.Bool                                  `tfsdk:"is_host"`
	Eths                          []veritySwitchpointObjectPropertiesEthModel `tfsdk:"eths"`
}

func (r *veritySwitchpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_switchpoint"
}

func (r *veritySwitchpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySwitchpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Switchpoint",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_serial_number": schema.StringAttribute{
				Description: "Device Serial Number",
				Optional:    true,
			},
			"connected_bundle": schema.StringAttribute{
				Description: "Connected Bundle",
				Optional:    true,
			},
			"connected_bundle_ref_type_": schema.StringAttribute{
				Description: "Object type for connected_bundle field",
				Optional:    true,
			},
			"read_only_mode": schema.BoolAttribute{
				Description: "When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware",
				Optional:    true,
			},
			"locked": schema.BoolAttribute{
				Description: "Permission lock",
				Optional:    true,
			},
			"disabled_ports": schema.StringAttribute{
				Description: "Disabled Ports - comma separated list of ports to disable",
				Optional:    true,
			},
			"out_of_band_management": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch is managed out of band via the management port",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of Switchpoint",
				Optional:    true,
			},
			"super_pod": schema.StringAttribute{
				Description: "Super Pod - subgrouping of super spines and pods",
				Optional:    true,
			},
			"pod": schema.StringAttribute{
				Description: "Pod - subgrouping of spine and leaf switches",
				Optional:    true,
			},
			"rack": schema.StringAttribute{
				Description: "Physical Rack location of the Switch",
				Optional:    true,
			},
			"switch_router_id_ip_mask": schema.StringAttribute{
				Description: "Switch BGP Router Identifier",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in switch_router_id_ip_mask field has been automatically assigned",
				Optional:    true,
			},
			"switch_vtep_id_ip_mask": schema.StringAttribute{
				Description: "Switch VTEP Identifier",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in switch_vtep_id_ip_mask field has been automatically assigned",
				Optional:    true,
			},
			"bgp_as_number": schema.Int64Attribute{
				Description: "BGP Autonomous System Number for the site underlay",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in bgp_as_number field has been automatically assigned",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"badges": schema.ListNestedBlock{
				Description: "Badge configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"badge": schema.StringAttribute{
							Description: "Badge name",
							Optional:    true,
						},
						"badge_ref_type_": schema.StringAttribute{
							Description: "Object type for badge field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
						},
					},
				},
			},
			"children": schema.ListNestedBlock{
				Description: "Child configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"child_num_endpoint": schema.StringAttribute{
							Description: "Switchpoint associated with the Child",
							Optional:    true,
						},
						"child_num_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for child_num_endpoint field",
							Optional:    true,
						},
						"child_num_device": schema.StringAttribute{
							Description: "Device associated with the Child",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
						},
					},
				},
			},
			"traffic_mirrors": schema.ListNestedBlock{
				Description: "Traffic mirror configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"traffic_mirror_num_enable": schema.BoolAttribute{
							Description: "Enable Traffic Mirror",
							Optional:    true,
						},
						"traffic_mirror_num_source_port": schema.StringAttribute{
							Description: "Source Port for Traffic Mirror",
							Optional:    true,
						},
						"traffic_mirror_num_source_lag_indicator": schema.BoolAttribute{
							Description: "Source LAG Indicator for Traffic Mirror",
							Optional:    true,
						},
						"traffic_mirror_num_destination_port": schema.StringAttribute{
							Description: "Destination Port for Traffic Mirror",
							Optional:    true,
						},
						"traffic_mirror_num_inbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for inbound traffic",
							Optional:    true,
						},
						"traffic_mirror_num_outbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for outbound traffic",
							Optional:    true,
						},
					},
				},
			},
			"eths": schema.ListNestedBlock{
				Description: "Ethernet port configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"breakout": schema.StringAttribute{
							Description: "Breakout Port Override",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the switchpoint",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"user_notes": schema.StringAttribute{
							Description: "Notes written by User about the site",
							Optional:    true,
						},
						"expected_parent_endpoint": schema.StringAttribute{
							Description: "Expected Parent Endpoint",
							Optional:    true,
						},
						"expected_parent_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for expected_parent_endpoint field",
							Optional:    true,
						},
						"number_of_multipoints": schema.Int64Attribute{
							Description: "Number of Multipoints",
							Optional:    true,
						},
						"aggregate": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes switch aggregated with all of its sub switches",
							Optional:    true,
						},
						"is_host": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes the Host Switch",
							Optional:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"eths": schema.ListNestedBlock{
							Description: "Ethernet port properties within object_properties",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"eth_num_icon": schema.StringAttribute{
										Description: "Icon of this Eth Port",
										Optional:    true,
									},
									"eth_num_label": schema.StringAttribute{
										Description: "Label of this Eth Port",
										Optional:    true,
									},
									"index": schema.Int64Attribute{
										Description: "The index identifying the object",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *veritySwitchpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySwitchpointResourceModel
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
	spProps := openapi.SwitchpointsPutRequestSwitchpointValue{}
	spProps.Name = openapi.PtrString(name)

	if !plan.DeviceSerialNumber.IsNull() {
		spProps.DeviceSerialNumber = openapi.PtrString(plan.DeviceSerialNumber.ValueString())
	}
	if !plan.ConnectedBundle.IsNull() {
		spProps.ConnectedBundle = openapi.PtrString(plan.ConnectedBundle.ValueString())
	}
	if !plan.ConnectedBundleRefType.IsNull() {
		spProps.ConnectedBundleRefType = openapi.PtrString(plan.ConnectedBundleRefType.ValueString())
	}
	if !plan.DisabledPorts.IsNull() {
		spProps.DisabledPorts = openapi.PtrString(plan.DisabledPorts.ValueString())
	}
	if !plan.Type.IsNull() {
		spProps.Type = openapi.PtrString(plan.Type.ValueString())
	}
	if !plan.SuperPod.IsNull() {
		spProps.SuperPod = openapi.PtrString(plan.SuperPod.ValueString())
	}
	if !plan.Pod.IsNull() {
		spProps.Pod = openapi.PtrString(plan.Pod.ValueString())
	}
	if !plan.Rack.IsNull() {
		spProps.Rack = openapi.PtrString(plan.Rack.ValueString())
	}
	if !plan.SwitchRouterIdIpMask.IsNull() {
		spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
	}
	if !plan.SwitchVtepIdIpMask.IsNull() {
		spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
	}

	if !plan.ReadOnlyMode.IsNull() {
		spProps.ReadOnlyMode = openapi.PtrBool(plan.ReadOnlyMode.ValueBool())
	}
	if !plan.Locked.IsNull() {
		spProps.Locked = openapi.PtrBool(plan.Locked.ValueBool())
	}
	if !plan.OutOfBandManagement.IsNull() {
		spProps.OutOfBandManagement = openapi.PtrBool(plan.OutOfBandManagement.ValueBool())
	}
	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
		spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool())
	}
	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
		spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool())
	}
	if !plan.BgpAsNumberAutoAssigned.IsNull() {
		spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(plan.BgpAsNumberAutoAssigned.ValueBool())
	}

	if !plan.BgpAsNumber.IsNull() {
		val := int32(plan.BgpAsNumber.ValueInt64())
		spProps.BgpAsNumber = &val
	} else {
		spProps.BgpAsNumber = nil
	}

	if len(plan.Badges) > 0 {
		badges := make([]openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner, len(plan.Badges))
		for i, badge := range plan.Badges {
			badgeItem := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}
			if !badge.Badge.IsNull() {
				badgeItem.Badge = openapi.PtrString(badge.Badge.ValueString())
			}
			if !badge.BadgeRefType.IsNull() {
				badgeItem.BadgeRefType = openapi.PtrString(badge.BadgeRefType.ValueString())
			}
			if !badge.Index.IsNull() {
				badgeItem.Index = openapi.PtrInt32(int32(badge.Index.ValueInt64()))
			}
			badges[i] = badgeItem
		}
		spProps.Badges = badges
	}

	if len(plan.Children) > 0 {
		children := make([]openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner, len(plan.Children))
		for i, child := range plan.Children {
			childItem := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}
			if !child.ChildNumEndpoint.IsNull() {
				childItem.ChildNumEndpoint = openapi.PtrString(child.ChildNumEndpoint.ValueString())
			}
			if !child.ChildNumEndpointRefType.IsNull() {
				childItem.ChildNumEndpointRefType = openapi.PtrString(child.ChildNumEndpointRefType.ValueString())
			}
			if !child.ChildNumDevice.IsNull() {
				childItem.ChildNumDevice = openapi.PtrString(child.ChildNumDevice.ValueString())
			}
			if !child.Index.IsNull() {
				childItem.Index = openapi.PtrInt32(int32(child.Index.ValueInt64()))
			}
			children[i] = childItem
		}
		spProps.Children = children
	}

	if len(plan.TrafficMirrors) > 0 {
		mirrors := make([]openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, len(plan.TrafficMirrors))
		for i, mirror := range plan.TrafficMirrors {
			mirrorItem := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}
			if !mirror.TrafficMirrorNumEnable.IsNull() {
				mirrorItem.TrafficMirrorNumEnable = openapi.PtrBool(mirror.TrafficMirrorNumEnable.ValueBool())
			}
			if !mirror.TrafficMirrorNumSourcePort.IsNull() {
				mirrorItem.TrafficMirrorNumSourcePort = openapi.PtrString(mirror.TrafficMirrorNumSourcePort.ValueString())
			}
			if !mirror.TrafficMirrorNumSourceLagIndicator.IsNull() {
				mirrorItem.TrafficMirrorNumSourceLagIndicator = openapi.PtrBool(mirror.TrafficMirrorNumSourceLagIndicator.ValueBool())
			}
			if !mirror.TrafficMirrorNumDestinationPort.IsNull() {
				mirrorItem.TrafficMirrorNumDestinationPort = openapi.PtrString(mirror.TrafficMirrorNumDestinationPort.ValueString())
			}
			if !mirror.TrafficMirrorNumInboundTraffic.IsNull() {
				mirrorItem.TrafficMirrorNumInboundTraffic = openapi.PtrBool(mirror.TrafficMirrorNumInboundTraffic.ValueBool())
			}
			if !mirror.TrafficMirrorNumOutboundTraffic.IsNull() {
				mirrorItem.TrafficMirrorNumOutboundTraffic = openapi.PtrBool(mirror.TrafficMirrorNumOutboundTraffic.ValueBool())
			}
			mirrors[i] = mirrorItem
		}
		spProps.TrafficMirrors = mirrors
	}

	if len(plan.Eths) > 0 {
		eths := make([]openapi.SwitchpointsPutRequestSwitchpointValueEthsInner, len(plan.Eths))
		for i, eth := range plan.Eths {
			ethItem := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}
			if !eth.Breakout.IsNull() {
				ethItem.Breakout = openapi.PtrString(eth.Breakout.ValueString())
			}
			if !eth.Index.IsNull() {
				ethItem.Index = openapi.PtrInt32(int32(eth.Index.ValueInt64()))
			}
			eths[i] = ethItem
		}
		spProps.Eths = eths
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		if !op.UserNotes.IsNull() {
			objProps.UserNotes = openapi.PtrString(op.UserNotes.ValueString())
		}
		if !op.ExpectedParentEndpoint.IsNull() {
			objProps.ExpectedParentEndpoint = openapi.PtrString(op.ExpectedParentEndpoint.ValueString())
		}
		if !op.ExpectedParentEndpointRefType.IsNull() {
			objProps.ExpectedParentEndpointRefType = openapi.PtrString(op.ExpectedParentEndpointRefType.ValueString())
		}
		if !op.NumberOfMultipoints.IsNull() {
			val := int32(op.NumberOfMultipoints.ValueInt64())
			objProps.NumberOfMultipoints = *openapi.NewNullableInt32(&val)
		} else {
			objProps.NumberOfMultipoints = *openapi.NewNullableInt32(nil)
		}
		if !op.Aggregate.IsNull() {
			objProps.Aggregate = openapi.PtrBool(op.Aggregate.ValueBool())
		}
		if !op.IsHost.IsNull() {
			objProps.IsHost = openapi.PtrBool(op.IsHost.ValueBool())
		}

		if len(op.Eths) > 0 {
			ethsSlice := make([]openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner, len(op.Eths))
			for i, eth := range op.Eths {
				ethItem := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{}
				if !eth.EthNumIcon.IsNull() {
					ethItem.EthNumIcon = openapi.PtrString(eth.EthNumIcon.ValueString())
				}
				if !eth.EthNumLabel.IsNull() {
					ethItem.EthNumLabel = openapi.PtrString(eth.EthNumLabel.ValueString())
				}
				if !eth.Index.IsNull() {
					ethItem.Index = openapi.PtrInt32(int32(eth.Index.ValueInt64()))
				}
				ethsSlice[i] = ethItem
			}
			objProps.Eths = ethsSlice
		}

		spProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "switchpoint", name, spProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for switchpoint creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Switchpoint %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")

	var minState veritySwitchpointResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if switchpointData, exists := bulkMgr.GetResourceResponse("switchpoint", name); exists {
			// Use the cached data with plan values as fallback
			state := r.populateSwitchpointState(ctx, minState, switchpointData, &plan)
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

func (r *veritySwitchpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySwitchpointResourceModel
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

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	spName := state.Name.ValueString()

	var switchpointData map[string]interface{}
	var exists bool

	if bulkOpsMgr != nil {
		switchpointData, exists = bulkOpsMgr.GetResourceResponse("switchpoint", spName)
		if exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached switchpoint data for %s from recent operation", spName))
			state = r.populateSwitchpointState(ctx, state, switchpointData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("switchpoint") {
		tflog.Info(ctx, fmt.Sprintf("Skipping switchpoint %s verification - trusting recent successful API operation", spName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent switchpoint operations found, performing normal verification for %s", spName))

	type SwitchpointResponse struct {
		Switchpoint map[string]interface{} `json:"switchpoint"`
	}

	var result SwitchpointResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		spData, fetchErr := getCachedResponse(ctx, r.provCtx, "switchpoints", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Switchpoints")
			respAPI, err := r.client.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Switchpoints: %v", err)
			}
			defer respAPI.Body.Close()

			var res SwitchpointResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Switchpoints response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Switchpoints", len(res.Switchpoint)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Switchpoints on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = spData.(SwitchpointResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Switchpoint %s", spName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Switchpoint with ID: %s", spName))
	var spData map[string]interface{}
	foundInAPI := false

	if data, ok := result.Switchpoint[spName].(map[string]interface{}); ok {
		spData = data
		foundInAPI = true
		tflog.Debug(ctx, fmt.Sprintf("Found Switchpoint directly by ID: %s", spName))
	} else {
		for apiName, s := range result.Switchpoint {
			switchpoint, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := switchpoint["name"].(string); ok && name == spName {
				spData = switchpoint
				spName = apiName
				foundInAPI = true
				tflog.Debug(ctx, fmt.Sprintf("Found Switchpoint with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !foundInAPI {
		tflog.Debug(ctx, fmt.Sprintf("Switchpoint with ID '%s' not found in API response", spName))
		resp.State.RemoveResource(ctx)
		return
	}

	state = r.populateSwitchpointState(ctx, state, spData, nil)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySwitchpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySwitchpointResourceModel

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

	if !plan.BgpAsNumber.Equal(state.BgpAsNumber) && !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'bgp_as_number' field cannot be modified because 'bgp_as_number_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) && !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'switch_router_id_ip_mask' field cannot be modified because 'switch_router_id_ip_mask_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) && !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'switch_vtep_id_ip_mask' field cannot be modified because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true.",
		)
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
	spProps := openapi.SwitchpointsPutRequestSwitchpointValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		spProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	stringFields := map[string]struct {
		planField  types.String
		stateField types.String
		setter     func(string)
	}{
		"device_serial_number":       {plan.DeviceSerialNumber, state.DeviceSerialNumber, func(v string) { spProps.DeviceSerialNumber = openapi.PtrString(v) }},
		"connected_bundle":           {plan.ConnectedBundle, state.ConnectedBundle, func(v string) { spProps.ConnectedBundle = openapi.PtrString(v) }},
		"connected_bundle_ref_type_": {plan.ConnectedBundleRefType, state.ConnectedBundleRefType, func(v string) { spProps.ConnectedBundleRefType = openapi.PtrString(v) }},
		"disabled_ports":             {plan.DisabledPorts, state.DisabledPorts, func(v string) { spProps.DisabledPorts = openapi.PtrString(v) }},
		"type":                       {plan.Type, state.Type, func(v string) { spProps.Type = openapi.PtrString(v) }},
		"super_pod":                  {plan.SuperPod, state.SuperPod, func(v string) { spProps.SuperPod = openapi.PtrString(v) }},
		"pod":                        {plan.Pod, state.Pod, func(v string) { spProps.Pod = openapi.PtrString(v) }},
		"rack":                       {plan.Rack, state.Rack, func(v string) { spProps.Rack = openapi.PtrString(v) }},
	}

	for _, field := range stringFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField.ValueString())
			hasChanges = true
		}
	}

	if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) {
		if plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() || !plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() ||
			!plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
			spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
			hasChanges = true
		}
	}

	if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) {
		if plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() || !plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() ||
			!plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
			spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
			hasChanges = true
		}
	}

	boolFields := map[string]struct {
		planField  types.Bool
		stateField types.Bool
		setter     func(bool)
	}{
		"read_only_mode":         {plan.ReadOnlyMode, state.ReadOnlyMode, func(v bool) { spProps.ReadOnlyMode = openapi.PtrBool(v) }},
		"locked":                 {plan.Locked, state.Locked, func(v bool) { spProps.Locked = openapi.PtrBool(v) }},
		"out_of_band_management": {plan.OutOfBandManagement, state.OutOfBandManagement, func(v bool) { spProps.OutOfBandManagement = openapi.PtrBool(v) }},
		"switch_router_id_ip_mask_auto_assigned_": {plan.SwitchRouterIdIpMaskAutoAssigned, state.SwitchRouterIdIpMaskAutoAssigned, func(v bool) { spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(v) }},
		"switch_vtep_id_ip_mask_auto_assigned_":   {plan.SwitchVtepIdIpMaskAutoAssigned, state.SwitchVtepIdIpMaskAutoAssigned, func(v bool) { spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(v) }},
		"bgp_as_number_auto_assigned_":            {plan.BgpAsNumberAutoAssigned, state.BgpAsNumberAutoAssigned, func(v bool) { spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(v) }},
	}

	for _, field := range boolFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField.ValueBool())
			hasChanges = true
		}
	}

	if !plan.BgpAsNumber.Equal(state.BgpAsNumber) {
		if plan.BgpAsNumberAutoAssigned.IsNull() || !plan.BgpAsNumberAutoAssigned.ValueBool() ||
			!plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
			if !plan.BgpAsNumber.IsNull() {
				val := int32(plan.BgpAsNumber.ValueInt64())
				spProps.BgpAsNumber = &val
			} else {
				spProps.BgpAsNumber = nil
			}
			hasChanges = true
		}
	}

	oldBadgesByIndex := make(map[int64]veritySwitchpointBadgeModel)
	for _, badge := range state.Badges {
		if !badge.Index.IsNull() {
			idx := badge.Index.ValueInt64()
			oldBadgesByIndex[idx] = badge
		}
	}

	var changedBadges []openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner
	badgesChanged := false

	for _, planBadge := range plan.Badges {
		if planBadge.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planBadge.Index.ValueInt64()
		stateBadge, exists := oldBadgesByIndex[idx]

		if !exists {
			// CREATE: new badge, include all fields
			newBadge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planBadge.Badge.IsNull() && planBadge.Badge.ValueString() != "" {
				newBadge.Badge = openapi.PtrString(planBadge.Badge.ValueString())
			} else {
				newBadge.Badge = openapi.PtrString("")
			}

			if !planBadge.BadgeRefType.IsNull() && planBadge.BadgeRefType.ValueString() != "" {
				newBadge.BadgeRefType = openapi.PtrString(planBadge.BadgeRefType.ValueString())
			} else {
				newBadge.BadgeRefType = openapi.PtrString("")
			}

			changedBadges = append(changedBadges, newBadge)
			badgesChanged = true
			continue
		}

		// UPDATE: existing badge, check which fields changed
		updateBadge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planBadge.Badge.Equal(stateBadge.Badge) {
			if !planBadge.Badge.IsNull() && planBadge.Badge.ValueString() != "" {
				updateBadge.Badge = openapi.PtrString(planBadge.Badge.ValueString())
			} else {
				updateBadge.Badge = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planBadge.BadgeRefType.Equal(stateBadge.BadgeRefType) {
			if !planBadge.BadgeRefType.IsNull() && planBadge.BadgeRefType.ValueString() != "" {
				updateBadge.BadgeRefType = openapi.PtrString(planBadge.BadgeRefType.ValueString())
			} else {
				updateBadge.BadgeRefType = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedBadges = append(changedBadges, updateBadge)
			badgesChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldBadgesByIndex {
		found := false
		for _, planBadge := range plan.Badges {
			if !planBadge.Index.IsNull() && planBadge.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// badge removed - include only the index for deletion
			deletedBadge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedBadges = append(changedBadges, deletedBadge)
			badgesChanged = true
		}
	}

	if badgesChanged && len(changedBadges) > 0 {
		spProps.Badges = changedBadges
		hasChanges = true
	}

	oldChildrenByIndex := make(map[int64]veritySwitchpointChildModel)
	for _, child := range state.Children {
		if !child.Index.IsNull() {
			idx := child.Index.ValueInt64()
			oldChildrenByIndex[idx] = child
		}
	}

	var changedChildren []openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner
	childrenChanged := false

	// Process all items in plan
	for _, planChild := range plan.Children {
		if planChild.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planChild.Index.ValueInt64()
		stateChild, exists := oldChildrenByIndex[idx]

		if !exists {
			// CREATE: new child, include all fields
			newChild := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planChild.ChildNumEndpoint.IsNull() && planChild.ChildNumEndpoint.ValueString() != "" {
				newChild.ChildNumEndpoint = openapi.PtrString(planChild.ChildNumEndpoint.ValueString())
			} else {
				newChild.ChildNumEndpoint = openapi.PtrString("")
			}

			if !planChild.ChildNumEndpointRefType.IsNull() && planChild.ChildNumEndpointRefType.ValueString() != "" {
				newChild.ChildNumEndpointRefType = openapi.PtrString(planChild.ChildNumEndpointRefType.ValueString())
			} else {
				newChild.ChildNumEndpointRefType = openapi.PtrString("")
			}

			if !planChild.ChildNumDevice.IsNull() && planChild.ChildNumDevice.ValueString() != "" {
				newChild.ChildNumDevice = openapi.PtrString(planChild.ChildNumDevice.ValueString())
			} else {
				newChild.ChildNumDevice = openapi.PtrString("")
			}

			changedChildren = append(changedChildren, newChild)
			childrenChanged = true
			continue
		}

		// UPDATE: existing child, check which fields changed
		updateChild := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planChild.ChildNumEndpoint.Equal(stateChild.ChildNumEndpoint) {
			if !planChild.ChildNumEndpoint.IsNull() && planChild.ChildNumEndpoint.ValueString() != "" {
				updateChild.ChildNumEndpoint = openapi.PtrString(planChild.ChildNumEndpoint.ValueString())
			} else {
				updateChild.ChildNumEndpoint = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planChild.ChildNumEndpointRefType.Equal(stateChild.ChildNumEndpointRefType) {
			if !planChild.ChildNumEndpointRefType.IsNull() && planChild.ChildNumEndpointRefType.ValueString() != "" {
				updateChild.ChildNumEndpointRefType = openapi.PtrString(planChild.ChildNumEndpointRefType.ValueString())
			} else {
				updateChild.ChildNumEndpointRefType = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planChild.ChildNumDevice.Equal(stateChild.ChildNumDevice) {
			if !planChild.ChildNumDevice.IsNull() && planChild.ChildNumDevice.ValueString() != "" {
				updateChild.ChildNumDevice = openapi.PtrString(planChild.ChildNumDevice.ValueString())
			} else {
				updateChild.ChildNumDevice = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedChildren = append(changedChildren, updateChild)
			childrenChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldChildrenByIndex {
		found := false
		for _, planChild := range plan.Children {
			if !planChild.Index.IsNull() && planChild.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// child removed - include only the index for deletion
			deletedChild := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedChildren = append(changedChildren, deletedChild)
			childrenChanged = true
		}
	}

	if childrenChanged && len(changedChildren) > 0 {
		spProps.Children = changedChildren
		hasChanges = true
	}

	if !r.equalTrafficMirrorArrays(plan.TrafficMirrors, state.TrafficMirrors) {
		mirrors := make([]openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, len(plan.TrafficMirrors))
		for i, mirror := range plan.TrafficMirrors {
			mirrorItem := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}
			if !mirror.TrafficMirrorNumEnable.IsNull() {
				mirrorItem.TrafficMirrorNumEnable = openapi.PtrBool(mirror.TrafficMirrorNumEnable.ValueBool())
			}
			if !mirror.TrafficMirrorNumSourcePort.IsNull() {
				mirrorItem.TrafficMirrorNumSourcePort = openapi.PtrString(mirror.TrafficMirrorNumSourcePort.ValueString())
			}
			if !mirror.TrafficMirrorNumSourceLagIndicator.IsNull() {
				mirrorItem.TrafficMirrorNumSourceLagIndicator = openapi.PtrBool(mirror.TrafficMirrorNumSourceLagIndicator.ValueBool())
			}
			if !mirror.TrafficMirrorNumDestinationPort.IsNull() {
				mirrorItem.TrafficMirrorNumDestinationPort = openapi.PtrString(mirror.TrafficMirrorNumDestinationPort.ValueString())
			}
			if !mirror.TrafficMirrorNumInboundTraffic.IsNull() {
				mirrorItem.TrafficMirrorNumInboundTraffic = openapi.PtrBool(mirror.TrafficMirrorNumInboundTraffic.ValueBool())
			}
			if !mirror.TrafficMirrorNumOutboundTraffic.IsNull() {
				mirrorItem.TrafficMirrorNumOutboundTraffic = openapi.PtrBool(mirror.TrafficMirrorNumOutboundTraffic.ValueBool())
			}
			mirrors[i] = mirrorItem
		}
		spProps.TrafficMirrors = mirrors
		hasChanges = true
	}

	oldEthsByIndex := make(map[int64]veritySwitchpointEthModel)
	for _, eth := range state.Eths {
		if !eth.Index.IsNull() {
			idx := eth.Index.ValueInt64()
			oldEthsByIndex[idx] = eth
		}
	}

	var changedEths []openapi.SwitchpointsPutRequestSwitchpointValueEthsInner
	ethsChanged := false

	for _, planEth := range plan.Eths {
		if planEth.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planEth.Index.ValueInt64()
		stateEth, exists := oldEthsByIndex[idx]

		if !exists {
			// CREATE: new eth, include all fields
			newEth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planEth.Breakout.IsNull() && planEth.Breakout.ValueString() != "" {
				newEth.Breakout = openapi.PtrString(planEth.Breakout.ValueString())
			} else {
				newEth.Breakout = openapi.PtrString("")
			}

			changedEths = append(changedEths, newEth)
			ethsChanged = true
			continue
		}

		// UPDATE: existing eth, check which fields changed
		updateEth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planEth.Breakout.Equal(stateEth.Breakout) {
			if !planEth.Breakout.IsNull() && planEth.Breakout.ValueString() != "" {
				updateEth.Breakout = openapi.PtrString(planEth.Breakout.ValueString())
			} else {
				updateEth.Breakout = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedEths = append(changedEths, updateEth)
			ethsChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldEthsByIndex {
		found := false
		for _, planEth := range plan.Eths {
			if !planEth.Index.IsNull() && planEth.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// eth removed - include only the index for deletion
			deletedEth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedEths = append(changedEths, deletedEth)
			ethsChanged = true
		}
	}

	if ethsChanged && len(changedEths) > 0 {
		spProps.Eths = changedEths
		hasChanges = true
	}

	objectPropertiesChanged := false
	var objProps openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties

	if len(plan.ObjectProperties) > 0 {
		planOP := plan.ObjectProperties[0]
		var stateOP *veritySwitchpointObjectPropertiesModel
		if len(state.ObjectProperties) > 0 {
			stateOP = &state.ObjectProperties[0]
		}

		// Check if non-eths fields changed
		fieldsChanged := false
		if stateOP == nil ||
			!planOP.UserNotes.Equal(stateOP.UserNotes) ||
			!planOP.ExpectedParentEndpoint.Equal(stateOP.ExpectedParentEndpoint) ||
			!planOP.ExpectedParentEndpointRefType.Equal(stateOP.ExpectedParentEndpointRefType) ||
			!planOP.NumberOfMultipoints.Equal(stateOP.NumberOfMultipoints) ||
			!planOP.Aggregate.Equal(stateOP.Aggregate) ||
			!planOP.IsHost.Equal(stateOP.IsHost) {
			fieldsChanged = true
		}

		if fieldsChanged {
			if !planOP.UserNotes.IsNull() {
				objProps.UserNotes = openapi.PtrString(planOP.UserNotes.ValueString())
			}
			if !planOP.ExpectedParentEndpoint.IsNull() {
				objProps.ExpectedParentEndpoint = openapi.PtrString(planOP.ExpectedParentEndpoint.ValueString())
			}
			if !planOP.ExpectedParentEndpointRefType.IsNull() {
				objProps.ExpectedParentEndpointRefType = openapi.PtrString(planOP.ExpectedParentEndpointRefType.ValueString())
			}
			if !planOP.NumberOfMultipoints.IsNull() {
				val := int32(planOP.NumberOfMultipoints.ValueInt64())
				objProps.NumberOfMultipoints = *openapi.NewNullableInt32(&val)
			} else {
				objProps.NumberOfMultipoints = *openapi.NewNullableInt32(nil)
			}
			if !planOP.Aggregate.IsNull() {
				objProps.Aggregate = openapi.PtrBool(planOP.Aggregate.ValueBool())
			}
			if !planOP.IsHost.IsNull() {
				objProps.IsHost = openapi.PtrBool(planOP.IsHost.ValueBool())
			}
			objectPropertiesChanged = true
		}

		var stateEths []veritySwitchpointObjectPropertiesEthModel
		if stateOP != nil {
			stateEths = stateOP.Eths
		}

		oldEthsByIndex := make(map[int64]veritySwitchpointObjectPropertiesEthModel)
		for _, eth := range stateEths {
			if !eth.Index.IsNull() {
				idx := eth.Index.ValueInt64()
				oldEthsByIndex[idx] = eth
			}
		}

		var changedEths []openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner
		ethsChanged := false

		for _, planEth := range planOP.Eths {
			if planEth.Index.IsNull() {
				continue // Skip items without identifier
			}

			idx := planEth.Index.ValueInt64()
			stateEth, exists := oldEthsByIndex[idx]

			if !exists {
				// CREATE: new eth, include all fields
				newEth := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{
					Index: openapi.PtrInt32(int32(idx)),
				}

				if !planEth.EthNumIcon.IsNull() && planEth.EthNumIcon.ValueString() != "" {
					newEth.EthNumIcon = openapi.PtrString(planEth.EthNumIcon.ValueString())
				} else {
					newEth.EthNumIcon = openapi.PtrString("empty")
				}

				if !planEth.EthNumLabel.IsNull() && planEth.EthNumLabel.ValueString() != "" {
					newEth.EthNumLabel = openapi.PtrString(planEth.EthNumLabel.ValueString())
				} else {
					newEth.EthNumLabel = openapi.PtrString("")
				}

				changedEths = append(changedEths, newEth)
				ethsChanged = true
				continue
			}

			// UPDATE: existing eth, check which fields changed
			updateEth := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			fieldChanged := false

			if !planEth.EthNumIcon.Equal(stateEth.EthNumIcon) {
				if !planEth.EthNumIcon.IsNull() && planEth.EthNumIcon.ValueString() != "" {
					updateEth.EthNumIcon = openapi.PtrString(planEth.EthNumIcon.ValueString())
				} else {
					updateEth.EthNumIcon = openapi.PtrString("empty")
				}
				fieldChanged = true
			}

			if !planEth.EthNumLabel.Equal(stateEth.EthNumLabel) {
				if !planEth.EthNumLabel.IsNull() && planEth.EthNumLabel.ValueString() != "" {
					updateEth.EthNumLabel = openapi.PtrString(planEth.EthNumLabel.ValueString())
				} else {
					updateEth.EthNumLabel = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if fieldChanged {
				changedEths = append(changedEths, updateEth)
				ethsChanged = true
			}
		}

		// DELETE: Check for deleted items
		for stateIdx := range oldEthsByIndex {
			found := false
			for _, planEth := range planOP.Eths {
				if !planEth.Index.IsNull() && planEth.Index.ValueInt64() == stateIdx {
					found = true
					break
				}
			}

			if !found {
				// eth removed - include only the index for deletion
				deletedEth := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{
					Index: openapi.PtrInt32(int32(stateIdx)),
				}
				changedEths = append(changedEths, deletedEth)
				ethsChanged = true
			}
		}

		if ethsChanged && len(changedEths) > 0 {
			objProps.Eths = changedEths
			objectPropertiesChanged = true
		}
	}

	if objectPropertiesChanged {
		spProps.ObjectProperties = &objProps
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "switchpoint", name, spProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Switchpoint update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Switchpoint %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")

	var minState veritySwitchpointResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if switchpointData, exists := bulkMgr.GetResourceResponse("switchpoint", name); exists {
			// Use the cached data from the API response with plan values as fallback
			state := r.populateSwitchpointState(ctx, minState, switchpointData, &plan)
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

func (r *veritySwitchpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySwitchpointResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "switchpoint", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Switchpoint deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Switchpoint %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")
	resp.State.RemoveResource(ctx)
}

func (r *veritySwitchpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *veritySwitchpointResource) equalTrafficMirrorArrays(a, b []veritySwitchpointTrafficMirrorModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].TrafficMirrorNumEnable.Equal(b[i].TrafficMirrorNumEnable) ||
			!a[i].TrafficMirrorNumSourcePort.Equal(b[i].TrafficMirrorNumSourcePort) ||
			!a[i].TrafficMirrorNumSourceLagIndicator.Equal(b[i].TrafficMirrorNumSourceLagIndicator) ||
			!a[i].TrafficMirrorNumDestinationPort.Equal(b[i].TrafficMirrorNumDestinationPort) ||
			!a[i].TrafficMirrorNumInboundTraffic.Equal(b[i].TrafficMirrorNumInboundTraffic) ||
			!a[i].TrafficMirrorNumOutboundTraffic.Equal(b[i].TrafficMirrorNumOutboundTraffic) {
			return false
		}
	}
	return true
}

func (r *veritySwitchpointResource) populateSwitchpointState(
	ctx context.Context,
	state veritySwitchpointResourceModel,
	switchpointData map[string]interface{},
	plan *veritySwitchpointResourceModel,
) veritySwitchpointResourceModel {
	state.Name = types.StringValue(fmt.Sprintf("%v", switchpointData["name"]))

	// For each field, check if it's in the API response first,
	// if not: use plan value (if plan provided and not null), otherwise preserve current state value

	stringFields := map[string]*types.String{
		"device_serial_number":       &state.DeviceSerialNumber,
		"connected_bundle":           &state.ConnectedBundle,
		"connected_bundle_ref_type_": &state.ConnectedBundleRefType,
		"disabled_ports":             &state.DisabledPorts,
		"type":                       &state.Type,
		"super_pod":                  &state.SuperPod,
		"pod":                        &state.Pod,
		"rack":                       &state.Rack,
	}

	for apiKey, stateField := range stringFields {
		if val, ok := switchpointData[apiKey].(string); ok {
			*stateField = types.StringValue(val)
		} else if plan != nil {
			switch apiKey {
			case "device_serial_number":
				if !plan.DeviceSerialNumber.IsNull() {
					*stateField = plan.DeviceSerialNumber
				}
			case "connected_bundle":
				if !plan.ConnectedBundle.IsNull() {
					*stateField = plan.ConnectedBundle
				}
			case "connected_bundle_ref_type_":
				if !plan.ConnectedBundleRefType.IsNull() {
					*stateField = plan.ConnectedBundleRefType
				}
			case "disabled_ports":
				if !plan.DisabledPorts.IsNull() {
					*stateField = plan.DisabledPorts
				}
			case "type":
				if !plan.Type.IsNull() {
					*stateField = plan.Type
				}
			case "super_pod":
				if !plan.SuperPod.IsNull() {
					*stateField = plan.SuperPod
				}
			case "pod":
				if !plan.Pod.IsNull() {
					*stateField = plan.Pod
				}
			case "rack":
				if !plan.Rack.IsNull() {
					*stateField = plan.Rack
				}
			}
		}
	}

	// Auto-assigned string fields with special handling
	// For auto-assigned fields, always use API values when available
	if val, ok := switchpointData["switch_router_id_ip_mask"].(string); ok {
		state.SwitchRouterIdIpMask = types.StringValue(val)
	} else if plan != nil && !plan.SwitchRouterIdIpMask.IsNull() {
		state.SwitchRouterIdIpMask = plan.SwitchRouterIdIpMask
	}

	if val, ok := switchpointData["switch_router_id_ip_mask_auto_assigned_"].(bool); ok {
		state.SwitchRouterIdIpMaskAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
		state.SwitchRouterIdIpMaskAutoAssigned = plan.SwitchRouterIdIpMaskAutoAssigned
	}

	if val, ok := switchpointData["switch_vtep_id_ip_mask"].(string); ok {
		state.SwitchVtepIdIpMask = types.StringValue(val)
	} else if plan != nil && !plan.SwitchVtepIdIpMask.IsNull() {
		state.SwitchVtepIdIpMask = plan.SwitchVtepIdIpMask
	}

	if val, ok := switchpointData["switch_vtep_id_ip_mask_auto_assigned_"].(bool); ok {
		state.SwitchVtepIdIpMaskAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
		state.SwitchVtepIdIpMaskAutoAssigned = plan.SwitchVtepIdIpMaskAutoAssigned
	}

	if val, ok := switchpointData["read_only_mode"].(bool); ok {
		state.ReadOnlyMode = types.BoolValue(val)
	} else if plan != nil && !plan.ReadOnlyMode.IsNull() {
		state.ReadOnlyMode = plan.ReadOnlyMode
	}

	if val, ok := switchpointData["locked"].(bool); ok {
		state.Locked = types.BoolValue(val)
	} else if plan != nil && !plan.Locked.IsNull() {
		state.Locked = plan.Locked
	}

	if val, ok := switchpointData["out_of_band_management"].(bool); ok {
		state.OutOfBandManagement = types.BoolValue(val)
	} else if plan != nil && !plan.OutOfBandManagement.IsNull() {
		state.OutOfBandManagement = plan.OutOfBandManagement
	}

	if val, ok := switchpointData["bgp_as_number"]; ok {
		if val == nil {
			state.BgpAsNumber = types.Int64Null()
		} else {
			switch v := val.(type) {
			case int:
				state.BgpAsNumber = types.Int64Value(int64(v))
			case int32:
				state.BgpAsNumber = types.Int64Value(int64(v))
			case int64:
				state.BgpAsNumber = types.Int64Value(v)
			case float64:
				state.BgpAsNumber = types.Int64Value(int64(v))
			case string:
				if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
					state.BgpAsNumber = types.Int64Value(intVal)
				} else {
					if plan != nil && !plan.BgpAsNumber.IsNull() {
						state.BgpAsNumber = plan.BgpAsNumber
					}
				}
			default:
				if plan != nil && !plan.BgpAsNumber.IsNull() {
					state.BgpAsNumber = plan.BgpAsNumber
				}
			}
		}
	} else if plan != nil && !plan.BgpAsNumber.IsNull() {
		state.BgpAsNumber = plan.BgpAsNumber
	}

	if val, ok := switchpointData["bgp_as_number_auto_assigned_"].(bool); ok {
		state.BgpAsNumberAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.BgpAsNumberAutoAssigned.IsNull() {
		state.BgpAsNumberAutoAssigned = plan.BgpAsNumberAutoAssigned
	}

	if badgesArray, ok := switchpointData["badges"].([]interface{}); ok && len(badgesArray) > 0 {
		var badges []veritySwitchpointBadgeModel
		for _, b := range badgesArray {
			badge, ok := b.(map[string]interface{})
			if !ok {
				continue
			}
			badgeModel := veritySwitchpointBadgeModel{}
			if val, ok := badge["badge"].(string); ok {
				badgeModel.Badge = types.StringValue(val)
			} else {
				badgeModel.Badge = types.StringNull()
			}
			if val, ok := badge["badge_ref_type_"].(string); ok {
				badgeModel.BadgeRefType = types.StringValue(val)
			} else {
				badgeModel.BadgeRefType = types.StringNull()
			}
			if index, ok := badge["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					badgeModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					badgeModel.Index = types.Int64Value(int64(intVal))
				} else {
					badgeModel.Index = types.Int64Null()
				}
			} else {
				badgeModel.Index = types.Int64Null()
			}
			badges = append(badges, badgeModel)
		}
		state.Badges = badges
	} else if plan != nil && len(plan.Badges) > 0 {
		state.Badges = plan.Badges
	}

	if childrenArray, ok := switchpointData["children"].([]interface{}); ok && len(childrenArray) > 0 {
		var children []veritySwitchpointChildModel
		for _, c := range childrenArray {
			child, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			childModel := veritySwitchpointChildModel{}
			if val, ok := child["child_num_endpoint"].(string); ok {
				childModel.ChildNumEndpoint = types.StringValue(val)
			} else {
				childModel.ChildNumEndpoint = types.StringNull()
			}
			if val, ok := child["child_num_endpoint_ref_type_"].(string); ok {
				childModel.ChildNumEndpointRefType = types.StringValue(val)
			} else {
				childModel.ChildNumEndpointRefType = types.StringNull()
			}
			if val, ok := child["child_num_device"].(string); ok {
				childModel.ChildNumDevice = types.StringValue(val)
			} else {
				childModel.ChildNumDevice = types.StringNull()
			}
			if index, ok := child["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					childModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					childModel.Index = types.Int64Value(int64(intVal))
				} else {
					childModel.Index = types.Int64Null()
				}
			} else {
				childModel.Index = types.Int64Null()
			}
			children = append(children, childModel)
		}
		state.Children = children
	} else if plan != nil && len(plan.Children) > 0 {
		state.Children = plan.Children
	}

	if mirrorsArray, ok := switchpointData["traffic_mirrors"].([]interface{}); ok && len(mirrorsArray) > 0 {
		var mirrors []veritySwitchpointTrafficMirrorModel
		for _, m := range mirrorsArray {
			mirror, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			mirrorModel := veritySwitchpointTrafficMirrorModel{}
			if val, ok := mirror["traffic_mirror_num_enable"].(bool); ok {
				mirrorModel.TrafficMirrorNumEnable = types.BoolValue(val)
			} else {
				mirrorModel.TrafficMirrorNumEnable = types.BoolNull()
			}
			if val, ok := mirror["traffic_mirror_num_source_port"].(string); ok {
				mirrorModel.TrafficMirrorNumSourcePort = types.StringValue(val)
			} else {
				mirrorModel.TrafficMirrorNumSourcePort = types.StringNull()
			}
			if val, ok := mirror["traffic_mirror_num_source_lag_indicator"].(bool); ok {
				mirrorModel.TrafficMirrorNumSourceLagIndicator = types.BoolValue(val)
			} else {
				mirrorModel.TrafficMirrorNumSourceLagIndicator = types.BoolNull()
			}
			if val, ok := mirror["traffic_mirror_num_destination_port"].(string); ok {
				mirrorModel.TrafficMirrorNumDestinationPort = types.StringValue(val)
			} else {
				mirrorModel.TrafficMirrorNumDestinationPort = types.StringNull()
			}
			if val, ok := mirror["traffic_mirror_num_inbound_traffic"].(bool); ok {
				mirrorModel.TrafficMirrorNumInboundTraffic = types.BoolValue(val)
			} else {
				mirrorModel.TrafficMirrorNumInboundTraffic = types.BoolNull()
			}
			if val, ok := mirror["traffic_mirror_num_outbound_traffic"].(bool); ok {
				mirrorModel.TrafficMirrorNumOutboundTraffic = types.BoolValue(val)
			} else {
				mirrorModel.TrafficMirrorNumOutboundTraffic = types.BoolNull()
			}
			mirrors = append(mirrors, mirrorModel)
		}
		state.TrafficMirrors = mirrors
	} else if plan != nil && len(plan.TrafficMirrors) > 0 {
		state.TrafficMirrors = plan.TrafficMirrors
	}

	if ethsArray, ok := switchpointData["eths"].([]interface{}); ok && len(ethsArray) > 0 {
		var eths []veritySwitchpointEthModel
		for _, e := range ethsArray {
			eth, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			ethModel := veritySwitchpointEthModel{}
			if val, ok := eth["breakout"].(string); ok {
				ethModel.Breakout = types.StringValue(val)
			} else {
				ethModel.Breakout = types.StringNull()
			}
			if index, ok := eth["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					ethModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					ethModel.Index = types.Int64Value(int64(intVal))
				} else {
					ethModel.Index = types.Int64Null()
				}
			} else {
				ethModel.Index = types.Int64Null()
			}
			eths = append(eths, ethModel)
		}
		state.Eths = eths
	} else if plan != nil && len(plan.Eths) > 0 {
		state.Eths = plan.Eths
	}

	if objProps, ok := switchpointData["object_properties"].(map[string]interface{}); ok {
		op := veritySwitchpointObjectPropertiesModel{}
		if val, ok := objProps["user_notes"].(string); ok {
			op.UserNotes = types.StringValue(val)
		} else {
			op.UserNotes = types.StringNull()
		}
		if val, ok := objProps["expected_parent_endpoint"].(string); ok {
			op.ExpectedParentEndpoint = types.StringValue(val)
		} else {
			op.ExpectedParentEndpoint = types.StringNull()
		}
		if val, ok := objProps["expected_parent_endpoint_ref_type_"].(string); ok {
			op.ExpectedParentEndpointRefType = types.StringValue(val)
		} else {
			op.ExpectedParentEndpointRefType = types.StringNull()
		}
		if val, ok := objProps["number_of_multipoints"]; ok && val != nil {
			if intVal, ok := val.(float64); ok {
				op.NumberOfMultipoints = types.Int64Value(int64(intVal))
			} else if intVal, ok := val.(int); ok {
				op.NumberOfMultipoints = types.Int64Value(int64(intVal))
			} else {
				op.NumberOfMultipoints = types.Int64Null()
			}
		} else {
			op.NumberOfMultipoints = types.Int64Null()
		}
		if val, ok := objProps["aggregate"].(bool); ok {
			op.Aggregate = types.BoolValue(val)
		} else {
			op.Aggregate = types.BoolNull()
		}
		if val, ok := objProps["is_host"].(bool); ok {
			op.IsHost = types.BoolValue(val)
		} else {
			op.IsHost = types.BoolNull()
		}

		if ethsArray, ok := objProps["eths"].([]interface{}); ok && len(ethsArray) > 0 {
			var eths []veritySwitchpointObjectPropertiesEthModel
			for _, e := range ethsArray {
				eth, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				ethModel := veritySwitchpointObjectPropertiesEthModel{}
				if val, ok := eth["eth_num_icon"].(string); ok {
					ethModel.EthNumIcon = types.StringValue(val)
				} else {
					ethModel.EthNumIcon = types.StringNull()
				}
				if val, ok := eth["eth_num_label"].(string); ok {
					ethModel.EthNumLabel = types.StringValue(val)
				} else {
					ethModel.EthNumLabel = types.StringNull()
				}
				if index, ok := eth["index"]; ok && index != nil {
					if intVal, ok := index.(float64); ok {
						ethModel.Index = types.Int64Value(int64(intVal))
					} else if intVal, ok := index.(int); ok {
						ethModel.Index = types.Int64Value(int64(intVal))
					} else {
						ethModel.Index = types.Int64Null()
					}
				} else {
					ethModel.Index = types.Int64Null()
				}
				eths = append(eths, ethModel)
			}
			op.Eths = eths
		} else {
			op.Eths = nil
		}

		state.ObjectProperties = []veritySwitchpointObjectPropertiesModel{op}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
	}

	return state
}

func (r *veritySwitchpointResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip modification if we're deleting the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For new resources (where state is null)
	if req.State.Raw.IsNull() {
		if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
			resp.Diagnostics.AddWarning(
				"BGP AS Number will be assigned by the API",
				"The 'bgp_as_number' field value in your configuration will be ignored because 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
		}
		if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
			resp.Diagnostics.AddWarning(
				"Switch Router ID IP Mask will be assigned by the API",
				"The 'switch_router_id_ip_mask' field value in your configuration will be ignored because 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
		}
		if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
			resp.Diagnostics.AddWarning(
				"Switch VTEP ID IP Mask will be assigned by the API",
				"The 'switch_vtep_id_ip_mask' field value in your configuration will be ignored because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
		}
		return
	}

	var state veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only show warning and suppress diff if auto-assignment is enabled AND the field is actually changing
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() && !plan.BgpAsNumber.Equal(state.BgpAsNumber) {
		resp.Diagnostics.AddWarning(
			"Ignoring bgp_as_number changes with auto-assignment enabled",
			"The 'bgp_as_number' field changes will be ignored because 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
		)

		// Use current state value to suppress the diff
		if !state.BgpAsNumber.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), state.BgpAsNumber)
		}
	}

	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() && !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) {
		resp.Diagnostics.AddWarning(
			"Ignoring switch_router_id_ip_mask changes with auto-assignment enabled",
			"The 'switch_router_id_ip_mask' field changes will be ignored because 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
		)

		// Use current state value to suppress the diff
		if !state.SwitchRouterIdIpMask.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), state.SwitchRouterIdIpMask)
		}
	}

	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() && !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) {
		resp.Diagnostics.AddWarning(
			"Ignoring switch_vtep_id_ip_mask changes with auto-assignment enabled",
			"The 'switch_vtep_id_ip_mask' field changes will be ignored because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
		)

		// Use current state value to suppress the diff
		if !state.SwitchVtepIdIpMask.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), state.SwitchVtepIdIpMask)
		}
	}
}

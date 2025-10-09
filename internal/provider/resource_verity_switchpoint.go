package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
	Enable                           types.Bool                               `tfsdk:"enable"`
	DeviceSerialNumber               types.String                             `tfsdk:"device_serial_number"`
	ConnectedBundle                  types.String                             `tfsdk:"connected_bundle"`
	ConnectedBundleRefType           types.String                             `tfsdk:"connected_bundle_ref_type_"`
	ReadOnlyMode                     types.Bool                               `tfsdk:"read_only_mode"`
	Locked                           types.Bool                               `tfsdk:"locked"`
	DisabledPorts                    types.String                             `tfsdk:"disabled_ports"`
	OutOfBandManagement              types.Bool                               `tfsdk:"out_of_band_management"`
	Type                             types.String                             `tfsdk:"type"`
	SpinePlane                       types.String                             `tfsdk:"spine_plane"`
	SpinePlaneRefType                types.String                             `tfsdk:"spine_plane_ref_type_"`
	Pod                              types.String                             `tfsdk:"pod"`
	PodRefType                       types.String                             `tfsdk:"pod_ref_type_"`
	Rack                             types.String                             `tfsdk:"rack"`
	SwitchRouterIdIpMask             types.String                             `tfsdk:"switch_router_id_ip_mask"`
	SwitchRouterIdIpMaskAutoAssigned types.Bool                               `tfsdk:"switch_router_id_ip_mask_auto_assigned_"`
	SwitchVtepIdIpMask               types.String                             `tfsdk:"switch_vtep_id_ip_mask"`
	SwitchVtepIdIpMaskAutoAssigned   types.Bool                               `tfsdk:"switch_vtep_id_ip_mask_auto_assigned_"`
	BgpAsNumber                      types.Int64                              `tfsdk:"bgp_as_number"`
	BgpAsNumberAutoAssigned          types.Bool                               `tfsdk:"bgp_as_number_auto_assigned_"`
	IsFabric                         types.Bool                               `tfsdk:"is_fabric"`
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

func (b veritySwitchpointBadgeModel) GetIndex() types.Int64 {
	return b.Index
}

type veritySwitchpointChildModel struct {
	ChildNumEndpoint        types.String `tfsdk:"child_num_endpoint"`
	ChildNumEndpointRefType types.String `tfsdk:"child_num_endpoint_ref_type_"`
	ChildNumDevice          types.String `tfsdk:"child_num_device"`
	Index                   types.Int64  `tfsdk:"index"`
}

func (c veritySwitchpointChildModel) GetIndex() types.Int64 {
	return c.Index
}

type veritySwitchpointTrafficMirrorModel struct {
	TrafficMirrorNumEnable             types.Bool   `tfsdk:"traffic_mirror_num_enable"`
	TrafficMirrorNumSourcePort         types.String `tfsdk:"traffic_mirror_num_source_port"`
	TrafficMirrorNumSourceLagIndicator types.Bool   `tfsdk:"traffic_mirror_num_source_lag_indicator"`
	TrafficMirrorNumDestinationPort    types.String `tfsdk:"traffic_mirror_num_destination_port"`
	TrafficMirrorNumInboundTraffic     types.Bool   `tfsdk:"traffic_mirror_num_inbound_traffic"`
	TrafficMirrorNumOutboundTraffic    types.Bool   `tfsdk:"traffic_mirror_num_outbound_traffic"`
	Index                              types.Int64  `tfsdk:"index"`
}

func (tm veritySwitchpointTrafficMirrorModel) GetIndex() types.Int64 {
	return tm.Index
}

type veritySwitchpointEthModel struct {
	Breakout types.String `tfsdk:"breakout"`
	Index    types.Int64  `tfsdk:"index"`
}

func (e veritySwitchpointEthModel) GetIndex() types.Int64 {
	return e.Index
}

type veritySwitchpointObjectPropertiesEthModel struct {
	EthNumIcon  types.String `tfsdk:"eth_num_icon"`
	EthNumLabel types.String `tfsdk:"eth_num_label"`
	Index       types.Int64  `tfsdk:"index"`
}

func (ope veritySwitchpointObjectPropertiesEthModel) GetIndex() types.Int64 {
	return ope.Index
}

type veritySwitchpointObjectPropertiesModel struct {
	UserNotes                     types.String                                `tfsdk:"user_notes"`
	ExpectedParentEndpoint        types.String                                `tfsdk:"expected_parent_endpoint"`
	ExpectedParentEndpointRefType types.String                                `tfsdk:"expected_parent_endpoint_ref_type_"`
	NumberOfMultipoints           types.Int64                                 `tfsdk:"number_of_multipoints"`
	Aggregate                     types.Bool                                  `tfsdk:"aggregate"`
	IsHost                        types.Bool                                  `tfsdk:"is_host"`
	DrawAsEdgeDevice              types.Bool                                  `tfsdk:"draw_as_edge_device"`
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
			"enable": schema.BoolAttribute{
				Description: "Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.",
				Optional:    true,
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
			"spine_plane": schema.StringAttribute{
				Description: "Spine Plane - subgrouping of super spine and spine",
				Optional:    true,
			},
			"spine_plane_ref_type_": schema.StringAttribute{
				Description: "Object type for spine_plane field",
				Optional:    true,
			},
			"pod": schema.StringAttribute{
				Description: "Pod - subgrouping of spine and leaf switches",
				Optional:    true,
			},
			"pod_ref_type_": schema.StringAttribute{
				Description: "Object type for pod field",
				Optional:    true,
			},
			"rack": schema.StringAttribute{
				Description: "Physical Rack location of the Switch",
				Optional:    true,
			},
			"switch_router_id_ip_mask": schema.StringAttribute{
				Description: "Switch BGP Router Identifier. This field should not be specified when 'switch_router_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch BGP Router Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_router_id_ip_mask' field in your configuration.",
				Optional:    true,
			},
			"switch_vtep_id_ip_mask": schema.StringAttribute{
				Description: "Switch VTEP Identifier. This field should not be specified when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch VTEP Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_vtep_id_ip_mask' field in your configuration.",
				Optional:    true,
			},
			"bgp_as_number": schema.Int64Attribute{
				Description: "BGP Autonomous System Number for the site underlay. This field should not be specified when 'bgp_as_number_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the BGP AS Number should be automatically assigned by the API. When set to true, do not specify the 'bgp_as_number' field in your configuration.",
				Optional:    true,
			},
			"is_fabric": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device",
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
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
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
						"draw_as_edge_device": schema.BoolAttribute{
							Description: "Turn on to display the switch as an edge device instead of as a switch",
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

	// Validate auto-assigned field specifications
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		if !plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() {
			resp.Diagnostics.AddError(
				"BGP AS Number cannot be specified when auto-assigned",
				"The 'bgp_as_number' field cannot be specified in the configuration when 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchRouterIdIpMask.IsNull() && !plan.SwitchRouterIdIpMask.IsUnknown() && plan.SwitchRouterIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch Router ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_router_id_ip_mask' field cannot be specified in the configuration when 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchVtepIdIpMask.IsNull() && !plan.SwitchVtepIdIpMask.IsUnknown() && plan.SwitchVtepIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch VTEP ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_vtep_id_ip_mask' field cannot be specified in the configuration when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	spProps := &openapi.SwitchpointsPutRequestSwitchpointValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "DeviceSerialNumber", APIField: &spProps.DeviceSerialNumber, TFValue: plan.DeviceSerialNumber},
		{FieldName: "ConnectedBundle", APIField: &spProps.ConnectedBundle, TFValue: plan.ConnectedBundle},
		{FieldName: "ConnectedBundleRefType", APIField: &spProps.ConnectedBundleRefType, TFValue: plan.ConnectedBundleRefType},
		{FieldName: "DisabledPorts", APIField: &spProps.DisabledPorts, TFValue: plan.DisabledPorts},
		{FieldName: "Type", APIField: &spProps.Type, TFValue: plan.Type},
		{FieldName: "SpinePlane", APIField: &spProps.SpinePlane, TFValue: plan.SpinePlane},
		{FieldName: "SpinePlaneRefType", APIField: &spProps.SpinePlaneRefType, TFValue: plan.SpinePlaneRefType},
		{FieldName: "Pod", APIField: &spProps.Pod, TFValue: plan.Pod},
		{FieldName: "PodRefType", APIField: &spProps.PodRefType, TFValue: plan.PodRefType},
		{FieldName: "Rack", APIField: &spProps.Rack, TFValue: plan.Rack},
		{FieldName: "SwitchRouterIdIpMask", APIField: &spProps.SwitchRouterIdIpMask, TFValue: plan.SwitchRouterIdIpMask},
		{FieldName: "SwitchVtepIdIpMask", APIField: &spProps.SwitchVtepIdIpMask, TFValue: plan.SwitchVtepIdIpMask},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &spProps.Enable, TFValue: plan.Enable},
		{FieldName: "ReadOnlyMode", APIField: &spProps.ReadOnlyMode, TFValue: plan.ReadOnlyMode},
		{FieldName: "Locked", APIField: &spProps.Locked, TFValue: plan.Locked},
		{FieldName: "OutOfBandManagement", APIField: &spProps.OutOfBandManagement, TFValue: plan.OutOfBandManagement},
		{FieldName: "SwitchRouterIdIpMaskAutoAssigned", APIField: &spProps.SwitchRouterIdIpMaskAutoAssigned, TFValue: plan.SwitchRouterIdIpMaskAutoAssigned},
		{FieldName: "SwitchVtepIdIpMaskAutoAssigned", APIField: &spProps.SwitchVtepIdIpMaskAutoAssigned, TFValue: plan.SwitchVtepIdIpMaskAutoAssigned},
		{FieldName: "BgpAsNumberAutoAssigned", APIField: &spProps.BgpAsNumberAutoAssigned, TFValue: plan.BgpAsNumberAutoAssigned},
		{FieldName: "IsFabric", APIField: &spProps.IsFabric, TFValue: plan.IsFabric},
	})

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "BgpAsNumber", APIField: &spProps.BgpAsNumber, TFValue: plan.BgpAsNumber},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		if !op.UserNotes.IsNull() {
			objProps.UserNotes = openapi.PtrString(op.UserNotes.ValueString())
		} else {
			objProps.UserNotes = nil
		}
		if !op.ExpectedParentEndpoint.IsNull() {
			objProps.ExpectedParentEndpoint = openapi.PtrString(op.ExpectedParentEndpoint.ValueString())
		} else {
			objProps.ExpectedParentEndpoint = nil
		}
		if !op.ExpectedParentEndpointRefType.IsNull() {
			objProps.ExpectedParentEndpointRefType = openapi.PtrString(op.ExpectedParentEndpointRefType.ValueString())
		} else {
			objProps.ExpectedParentEndpointRefType = nil
		}
		if !op.NumberOfMultipoints.IsNull() {
			val := int32(op.NumberOfMultipoints.ValueInt64())
			objProps.NumberOfMultipoints = *openapi.NewNullableInt32(&val)
		} else {
			objProps.NumberOfMultipoints = *openapi.NewNullableInt32(nil)
		}
		if !op.Aggregate.IsNull() {
			objProps.Aggregate = openapi.PtrBool(op.Aggregate.ValueBool())
		} else {
			objProps.Aggregate = nil
		}
		if !op.IsHost.IsNull() {
			objProps.IsHost = openapi.PtrBool(op.IsHost.ValueBool())
		} else {
			objProps.IsHost = nil
		}
		if !op.DrawAsEdgeDevice.IsNull() {
			objProps.DrawAsEdgeDevice = openapi.PtrBool(op.DrawAsEdgeDevice.ValueBool())
		} else {
			objProps.DrawAsEdgeDevice = nil
		}

		if len(op.Eths) > 0 {
			ethsSlice := make([]openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner, len(op.Eths))
			for i, eth := range op.Eths {
				ethItem := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{}
				if !eth.Index.IsNull() {
					ethItem.Index = openapi.PtrInt32(int32(eth.Index.ValueInt64()))
				}
				if !eth.EthNumIcon.IsNull() {
					ethItem.EthNumIcon = openapi.PtrString(eth.EthNumIcon.ValueString())
				}
				if !eth.EthNumLabel.IsNull() {
					ethItem.EthNumLabel = openapi.PtrString(eth.EthNumLabel.ValueString())
				}
				ethsSlice[i] = ethItem
			}
			objProps.Eths = ethsSlice
		}

		spProps.ObjectProperties = &objProps
	}

	// Handle badges
	if len(plan.Badges) > 0 {
		badges := make([]openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner, len(plan.Badges))
		for i, badge := range plan.Badges {
			badgeItem := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Badge", APIField: &badgeItem.Badge, TFValue: badge.Badge},
				{FieldName: "BadgeRefType", APIField: &badgeItem.BadgeRefType, TFValue: badge.BadgeRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &badgeItem.Index, TFValue: badge.Index},
			})

			badges[i] = badgeItem
		}
		spProps.Badges = badges
	}

	// Handle children
	if len(plan.Children) > 0 {
		children := make([]openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner, len(plan.Children))
		for i, child := range plan.Children {
			childItem := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "ChildNumEndpoint", APIField: &childItem.ChildNumEndpoint, TFValue: child.ChildNumEndpoint},
				{FieldName: "ChildNumEndpointRefType", APIField: &childItem.ChildNumEndpointRefType, TFValue: child.ChildNumEndpointRefType},
				{FieldName: "ChildNumDevice", APIField: &childItem.ChildNumDevice, TFValue: child.ChildNumDevice},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &childItem.Index, TFValue: child.Index},
			})

			children[i] = childItem
		}
		spProps.Children = children
	}

	// Handle traffic mirrors
	if len(plan.TrafficMirrors) > 0 {
		mirrors := make([]openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, len(plan.TrafficMirrors))
		for i, mirror := range plan.TrafficMirrors {
			mirrorItem := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "TrafficMirrorNumEnable", APIField: &mirrorItem.TrafficMirrorNumEnable, TFValue: mirror.TrafficMirrorNumEnable},
				{FieldName: "TrafficMirrorNumSourceLagIndicator", APIField: &mirrorItem.TrafficMirrorNumSourceLagIndicator, TFValue: mirror.TrafficMirrorNumSourceLagIndicator},
				{FieldName: "TrafficMirrorNumInboundTraffic", APIField: &mirrorItem.TrafficMirrorNumInboundTraffic, TFValue: mirror.TrafficMirrorNumInboundTraffic},
				{FieldName: "TrafficMirrorNumOutboundTraffic", APIField: &mirrorItem.TrafficMirrorNumOutboundTraffic, TFValue: mirror.TrafficMirrorNumOutboundTraffic},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "TrafficMirrorNumSourcePort", APIField: &mirrorItem.TrafficMirrorNumSourcePort, TFValue: mirror.TrafficMirrorNumSourcePort},
				{FieldName: "TrafficMirrorNumDestinationPort", APIField: &mirrorItem.TrafficMirrorNumDestinationPort, TFValue: mirror.TrafficMirrorNumDestinationPort},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &mirrorItem.Index, TFValue: mirror.Index},
			})

			mirrors[i] = mirrorItem
		}
		spProps.TrafficMirrors = mirrors
	}

	// Handle eths
	if len(plan.Eths) > 0 {
		eths := make([]openapi.SwitchpointsPutRequestSwitchpointValueEthsInner, len(plan.Eths))
		for i, eth := range plan.Eths {
			ethItem := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Breakout", APIField: &ethItem.Breakout, TFValue: eth.Breakout},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethItem.Index, TFValue: eth.Index},
			})

			eths[i] = ethItem
		}
		spProps.Eths = eths
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "switchpoint", name, *spProps, &resp.Diagnostics)
	if !success {
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

	spName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if switchpointData, exists := r.bulkOpsMgr.GetResourceResponse("switchpoint", spName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached switchpoint data for %s from recent operation", spName))
			state = r.populateSwitchpointState(ctx, state, switchpointData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("switchpoint") {
		tflog.Info(ctx, fmt.Sprintf("Skipping switchpoint %s verification – trusting recent successful API operation", spName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching switchpoints for verification of %s", spName))

	type SwitchpointResponse struct {
		Switchpoint map[string]interface{} `json:"switchpoint"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "switchpoints", spName,
		func() (SwitchpointResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch switchpoints")
			respAPI, err := r.client.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
			if err != nil {
				return SwitchpointResponse{}, fmt.Errorf("error reading switchpoints: %v", err)
			}
			defer respAPI.Body.Close()

			var res SwitchpointResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SwitchpointResponse{}, fmt.Errorf("failed to decode switchpoints response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d switchpoints", len(res.Switchpoint)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Switchpoint %s", spName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for switchpoint with name: %s", spName))

	spData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Switchpoint,
		spName,
		func(data interface{}) (string, bool) {
			if switchpoint, ok := data.(map[string]interface{}); ok {
				if name, ok := switchpoint["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Switchpoint with name '%s' not found in API response", spName))
		resp.State.RemoveResource(ctx)
		return
	}

	switchpointMap, ok := spData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Switchpoint Data",
			fmt.Sprintf("Switchpoint data is not in expected format for %s", spName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found switchpoint '%s' under API key '%s'", spName, actualAPIName))

	state = r.populateSwitchpointState(ctx, state, switchpointMap, nil)
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

	// Validate auto-assigned fields - these checks prevent ineffective API calls
	// Only error if the auto-assigned flag is enabled AND the user is explicitly setting a value
	// AND the auto-assigned flag itself is not changing (which would be a valid operation)
	// Don't error if the field is unknown (computed during plan recalculation)
	if !plan.BgpAsNumber.Equal(state.BgpAsNumber) &&
		!plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() && // User is explicitly setting a value
		!plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() &&
		plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'bgp_as_number' field cannot be modified because 'bgp_as_number_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) &&
		!plan.SwitchRouterIdIpMask.IsNull() && !plan.SwitchRouterIdIpMask.IsUnknown() && // User is explicitly setting a value
		!plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() &&
		plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'switch_router_id_ip_mask' field cannot be modified because 'switch_router_id_ip_mask_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) &&
		!plan.SwitchVtepIdIpMask.IsNull() && !plan.SwitchVtepIdIpMask.IsUnknown() && // User is explicitly setting a value
		!plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() &&
		plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { spProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DeviceSerialNumber, state.DeviceSerialNumber, func(v *string) { spProps.DeviceSerialNumber = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DisabledPorts, state.DisabledPorts, func(v *string) { spProps.DisabledPorts = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { spProps.Type = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Rack, state.Rack, func(v *string) { spProps.Rack = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { spProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ReadOnlyMode, state.ReadOnlyMode, func(v *bool) { spProps.ReadOnlyMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Locked, state.Locked, func(v *bool) { spProps.Locked = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.OutOfBandManagement, state.OutOfBandManagement, func(v *bool) { spProps.OutOfBandManagement = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsFabric, state.IsFabric, func(v *bool) { spProps.IsFabric = v }, &hasChanges)

	// Handle ConnectedBundle and ConnectedBundleRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ConnectedBundle, state.ConnectedBundle, plan.ConnectedBundleRefType, state.ConnectedBundleRefType,
		func(val *string) { spProps.ConnectedBundle = val },
		func(val *string) { spProps.ConnectedBundleRefType = val },
		"connected_bundle", "connected_bundle_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	// Handle SpinePlane and SpinePlaneRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.SpinePlane, state.SpinePlane, plan.SpinePlaneRefType, state.SpinePlaneRefType,
		func(v *string) { spProps.SpinePlane = v },
		func(v *string) { spProps.SpinePlaneRefType = v },
		"spine_plane", "spine_plane_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle Pod and PodRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Pod, state.Pod, plan.PodRefType, state.PodRefType,
		func(v *string) { spProps.Pod = v },
		func(v *string) { spProps.PodRefType = v },
		"pod", "pod_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle BgpAsNumber and BgpAsNumberAutoAssigned changes
	bgpAsNumberChanged := !plan.BgpAsNumber.IsUnknown() && !plan.BgpAsNumber.Equal(state.BgpAsNumber)
	bgpAsNumberAutoAssignedChanged := !plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned)

	if bgpAsNumberChanged || bgpAsNumberAutoAssignedChanged {
		// Handle BgpAsNumber field changes
		if bgpAsNumberChanged {
			if !plan.BgpAsNumber.IsNull() {
				val := int32(plan.BgpAsNumber.ValueInt64())
				spProps.BgpAsNumber = &val
			} else {
				spProps.BgpAsNumber = nil
			}
		}

		// Handle BgpAsNumberAutoAssigned field changes
		if bgpAsNumberAutoAssignedChanged {
			// Only send bgp_as_number_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedBgpAsNumberAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedBgpAsNumberAutoAssigned = !config.BgpAsNumberAutoAssigned.IsNull()
				}
			}

			if userSpecifiedBgpAsNumberAutoAssigned {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(plan.BgpAsNumberAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both bgp_as_number_auto_assigned_ and bgp_as_number fields to be sent.
				if !state.BgpAsNumberAutoAssigned.IsNull() && state.BgpAsNumberAutoAssigned.ValueBool() &&
					!plan.BgpAsNumberAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include BgpAsNumber value in the request for the change to take effect
					if !plan.BgpAsNumber.IsNull() {
						val := int32(plan.BgpAsNumber.ValueInt64())
						spProps.BgpAsNumber = &val
					} else if !state.BgpAsNumber.IsNull() {
						// Use current state BgpAsNumber if plan doesn't specify one
						val := int32(state.BgpAsNumber.ValueInt64())
						spProps.BgpAsNumber = &val
					}
				}
			}
		} else if bgpAsNumberChanged {
			// BgpAsNumber changed but BgpAsNumberAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.BgpAsNumberAutoAssigned.IsNull() {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(plan.BgpAsNumberAutoAssigned.ValueBool())
			} else if !state.BgpAsNumberAutoAssigned.IsNull() {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(state.BgpAsNumberAutoAssigned.ValueBool())
			} else {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle SwitchRouterIdIpMask and SwitchRouterIdIpMaskAutoAssigned changes
	switchRouterIdIpMaskChanged := !plan.SwitchRouterIdIpMask.IsUnknown() && !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask)
	switchRouterIdIpMaskAutoAssignedChanged := !plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned)

	if switchRouterIdIpMaskChanged || switchRouterIdIpMaskAutoAssignedChanged {
		// Handle SwitchRouterIdIpMask field changes
		if switchRouterIdIpMaskChanged {
			spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
		}

		// Handle SwitchRouterIdIpMaskAutoAssigned field changes
		if switchRouterIdIpMaskAutoAssignedChanged {
			// Only send switch_router_id_ip_mask_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedSwitchRouterIdIpMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedSwitchRouterIdIpMaskAutoAssigned = !config.SwitchRouterIdIpMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedSwitchRouterIdIpMaskAutoAssigned {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both switch_router_id_ip_mask_auto_assigned_ and switch_router_id_ip_mask fields to be sent.
				if !state.SwitchRouterIdIpMaskAutoAssigned.IsNull() && state.SwitchRouterIdIpMaskAutoAssigned.ValueBool() &&
					!plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include SwitchRouterIdIpMask value in the request for the change to take effect
					if !plan.SwitchRouterIdIpMask.IsNull() {
						spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
					} else if !state.SwitchRouterIdIpMask.IsNull() {
						// Use current state SwitchRouterIdIpMask if plan doesn't specify one
						spProps.SwitchRouterIdIpMask = openapi.PtrString(state.SwitchRouterIdIpMask.ValueString())
					}
				}
			}
		} else if switchRouterIdIpMaskChanged {
			// SwitchRouterIdIpMask changed but SwitchRouterIdIpMaskAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool())
			} else if !state.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(state.SwitchRouterIdIpMaskAutoAssigned.ValueBool())
			} else {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle SwitchVtepIdIpMask and SwitchVtepIdIpMaskAutoAssigned changes
	switchVtepIdIpMaskChanged := !plan.SwitchVtepIdIpMask.IsUnknown() && !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask)
	switchVtepIdIpMaskAutoAssignedChanged := !plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned)

	if switchVtepIdIpMaskChanged || switchVtepIdIpMaskAutoAssignedChanged {
		// Handle SwitchVtepIdIpMask field changes
		if switchVtepIdIpMaskChanged {
			spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
		}

		// Handle SwitchVtepIdIpMaskAutoAssigned field changes
		if switchVtepIdIpMaskAutoAssignedChanged {
			// Only send switch_vtep_id_ip_mask_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedSwitchVtepIdIpMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedSwitchVtepIdIpMaskAutoAssigned = !config.SwitchVtepIdIpMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedSwitchVtepIdIpMaskAutoAssigned {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both switch_vtep_id_ip_mask_auto_assigned_ and switch_vtep_id_ip_mask fields to be sent.
				if !state.SwitchVtepIdIpMaskAutoAssigned.IsNull() && state.SwitchVtepIdIpMaskAutoAssigned.ValueBool() &&
					!plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include SwitchVtepIdIpMask value in the request for the change to take effect
					if !plan.SwitchVtepIdIpMask.IsNull() {
						spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
					} else if !state.SwitchVtepIdIpMask.IsNull() {
						// Use current state SwitchVtepIdIpMask if plan doesn't specify one
						spProps.SwitchVtepIdIpMask = openapi.PtrString(state.SwitchVtepIdIpMask.ValueString())
					}
				}
			}
		} else if switchVtepIdIpMaskChanged {
			// SwitchVtepIdIpMask changed but SwitchVtepIdIpMaskAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool())
			} else if !state.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(state.SwitchVtepIdIpMaskAutoAssigned.ValueBool())
			} else {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle badges
	changedBadges, badgesChanged := utils.ProcessIndexedArrayUpdates(plan.Badges, state.Badges,
		utils.IndexedItemHandler[veritySwitchpointBadgeModel, openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner]{
			CreateNew: func(planItem veritySwitchpointBadgeModel) openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner {
				badge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Badge", APIField: &badge.Badge, TFValue: planItem.Badge},
					{FieldName: "BadgeRefType", APIField: &badge.BadgeRefType, TFValue: planItem.BadgeRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &badge.Index, TFValue: planItem.Index},
				})

				return badge
			},
			UpdateExisting: func(planItem veritySwitchpointBadgeModel, stateItem veritySwitchpointBadgeModel) (openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner, bool) {
				badge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}
				fieldChanged := false

				// Handle badge and badge_ref_type_ using "One ref type supported" pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Badge, stateItem.Badge, planItem.BadgeRefType, stateItem.BadgeRefType,
					func(v *string) { badge.Badge = v },
					func(v *string) { badge.BadgeRefType = v },
					"badge", "badge_ref_type_",
					&fieldChanged,
					&resp.Diagnostics,
				) {
					return badge, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { badge.Index = v }, &fieldChanged)

				return badge, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if badgesChanged {
		spProps.Badges = changedBadges
		hasChanges = true
	}

	// Handle children
	changedChildren, childrenChanged := utils.ProcessIndexedArrayUpdates(plan.Children, state.Children,
		utils.IndexedItemHandler[veritySwitchpointChildModel, openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner]{
			CreateNew: func(planItem veritySwitchpointChildModel) openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner {
				child := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "ChildNumEndpoint", APIField: &child.ChildNumEndpoint, TFValue: planItem.ChildNumEndpoint},
					{FieldName: "ChildNumEndpointRefType", APIField: &child.ChildNumEndpointRefType, TFValue: planItem.ChildNumEndpointRefType},
					{FieldName: "ChildNumDevice", APIField: &child.ChildNumDevice, TFValue: planItem.ChildNumDevice},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &child.Index, TFValue: planItem.Index},
				})

				return child
			},
			UpdateExisting: func(planItem veritySwitchpointChildModel, stateItem veritySwitchpointChildModel) (openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner, bool) {
				child := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}
				fieldChanged := false

				// Handle child_num_endpoint and child_num_endpoint_ref_type_ using "One ref type supported" pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.ChildNumEndpoint, stateItem.ChildNumEndpoint, planItem.ChildNumEndpointRefType, stateItem.ChildNumEndpointRefType,
					func(v *string) { child.ChildNumEndpoint = v },
					func(v *string) { child.ChildNumEndpointRefType = v },
					"child_num_endpoint", "child_num_endpoint_ref_type_",
					&fieldChanged,
					&resp.Diagnostics,
				) {
					return child, false
				}

				// Handle other string field changes
				utils.CompareAndSetStringField(planItem.ChildNumDevice, stateItem.ChildNumDevice, func(v *string) { child.ChildNumDevice = v }, &fieldChanged)

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { child.Index = v }, &fieldChanged)

				return child, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if childrenChanged {
		spProps.Children = changedChildren
		hasChanges = true
	}

	// Handle traffic mirrors
	changedTrafficMirrors, trafficMirrorsChanged := utils.ProcessIndexedArrayUpdates(plan.TrafficMirrors, state.TrafficMirrors,
		utils.IndexedItemHandler[veritySwitchpointTrafficMirrorModel, openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner]{
			CreateNew: func(planItem veritySwitchpointTrafficMirrorModel) openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner {
				mirror := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "TrafficMirrorNumEnable", APIField: &mirror.TrafficMirrorNumEnable, TFValue: planItem.TrafficMirrorNumEnable},
					{FieldName: "TrafficMirrorNumSourceLagIndicator", APIField: &mirror.TrafficMirrorNumSourceLagIndicator, TFValue: planItem.TrafficMirrorNumSourceLagIndicator},
					{FieldName: "TrafficMirrorNumInboundTraffic", APIField: &mirror.TrafficMirrorNumInboundTraffic, TFValue: planItem.TrafficMirrorNumInboundTraffic},
					{FieldName: "TrafficMirrorNumOutboundTraffic", APIField: &mirror.TrafficMirrorNumOutboundTraffic, TFValue: planItem.TrafficMirrorNumOutboundTraffic},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "TrafficMirrorNumSourcePort", APIField: &mirror.TrafficMirrorNumSourcePort, TFValue: planItem.TrafficMirrorNumSourcePort},
					{FieldName: "TrafficMirrorNumDestinationPort", APIField: &mirror.TrafficMirrorNumDestinationPort, TFValue: planItem.TrafficMirrorNumDestinationPort},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &mirror.Index, TFValue: planItem.Index},
				})

				return mirror
			},
			UpdateExisting: func(planItem veritySwitchpointTrafficMirrorModel, stateItem veritySwitchpointTrafficMirrorModel) (openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, bool) {
				mirror := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumEnable, stateItem.TrafficMirrorNumEnable, func(v *bool) { mirror.TrafficMirrorNumEnable = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumSourceLagIndicator, stateItem.TrafficMirrorNumSourceLagIndicator, func(v *bool) { mirror.TrafficMirrorNumSourceLagIndicator = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumInboundTraffic, stateItem.TrafficMirrorNumInboundTraffic, func(v *bool) { mirror.TrafficMirrorNumInboundTraffic = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumOutboundTraffic, stateItem.TrafficMirrorNumOutboundTraffic, func(v *bool) { mirror.TrafficMirrorNumOutboundTraffic = v }, &fieldChanged)

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.TrafficMirrorNumSourcePort, stateItem.TrafficMirrorNumSourcePort, func(v *string) { mirror.TrafficMirrorNumSourcePort = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.TrafficMirrorNumDestinationPort, stateItem.TrafficMirrorNumDestinationPort, func(v *string) { mirror.TrafficMirrorNumDestinationPort = v }, &fieldChanged)

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { mirror.Index = v }, &fieldChanged)

				return mirror, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if trafficMirrorsChanged {
		spProps.TrafficMirrors = changedTrafficMirrors
		hasChanges = true
	}

	// Handle eths
	changedEths, ethsChanged := utils.ProcessIndexedArrayUpdates(plan.Eths, state.Eths,
		utils.IndexedItemHandler[veritySwitchpointEthModel, openapi.SwitchpointsPutRequestSwitchpointValueEthsInner]{
			CreateNew: func(planItem veritySwitchpointEthModel) openapi.SwitchpointsPutRequestSwitchpointValueEthsInner {
				eth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Breakout", APIField: &eth.Breakout, TFValue: planItem.Breakout},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &eth.Index, TFValue: planItem.Index},
				})

				return eth
			},
			UpdateExisting: func(planItem veritySwitchpointEthModel, stateItem veritySwitchpointEthModel) (openapi.SwitchpointsPutRequestSwitchpointValueEthsInner, bool) {
				eth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}
				fieldChanged := false

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.Breakout, stateItem.Breakout, func(v *string) { eth.Breakout = v }, &fieldChanged)

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { eth.Index = v }, &fieldChanged)

				return eth, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueEthsInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if ethsChanged {
		spProps.Eths = changedEths
		hasChanges = true
	}

	// Handle object_properties
	objectPropertiesChanged := false
	var objProps openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties

	if len(plan.ObjectProperties) > 0 {
		planOP := plan.ObjectProperties[0]
		var stateOP *veritySwitchpointObjectPropertiesModel
		if len(state.ObjectProperties) > 0 {
			stateOP = &state.ObjectProperties[0]
		}

		// Prepare state values for comparison (use null if stateOP doesn't exist)
		var stateUserNotes, stateExpectedParentEndpoint, stateExpectedParentEndpointRefType types.String
		var stateNumberOfMultipoints types.Int64
		var stateAggregate, stateIsHost, stateDrawAsEdgeDevice types.Bool

		if stateOP != nil {
			stateUserNotes = stateOP.UserNotes
			stateExpectedParentEndpoint = stateOP.ExpectedParentEndpoint
			stateExpectedParentEndpointRefType = stateOP.ExpectedParentEndpointRefType
			stateNumberOfMultipoints = stateOP.NumberOfMultipoints
			stateAggregate = stateOP.Aggregate
			stateIsHost = stateOP.IsHost
			stateDrawAsEdgeDevice = stateOP.DrawAsEdgeDevice
		} else {
			stateUserNotes = types.StringNull()
			stateExpectedParentEndpoint = types.StringNull()
			stateExpectedParentEndpointRefType = types.StringNull()
			stateNumberOfMultipoints = types.Int64Null()
			stateAggregate = types.BoolNull()
			stateIsHost = types.BoolNull()
			stateDrawAsEdgeDevice = types.BoolNull()
		}

		// Handle ExpectedParentEndpoint and ExpectedParentEndpointRefType using "One ref type supported" pattern
		if !utils.HandleOneRefTypeSupported(
			planOP.ExpectedParentEndpoint, stateExpectedParentEndpoint,
			planOP.ExpectedParentEndpointRefType, stateExpectedParentEndpointRefType,
			func(v *string) { objProps.ExpectedParentEndpoint = v },
			func(v *string) { objProps.ExpectedParentEndpointRefType = v },
			"expected_parent_endpoint", "expected_parent_endpoint_ref_type_",
			&objectPropertiesChanged,
			&resp.Diagnostics,
		) {
			return
		}

		if !planOP.UserNotes.Equal(stateUserNotes) {
			if !planOP.UserNotes.IsNull() {
				objProps.UserNotes = openapi.PtrString(planOP.UserNotes.ValueString())
			}
			objectPropertiesChanged = true
		}

		if !planOP.NumberOfMultipoints.Equal(stateNumberOfMultipoints) {
			if !planOP.NumberOfMultipoints.IsNull() {
				val := int32(planOP.NumberOfMultipoints.ValueInt64())
				objProps.NumberOfMultipoints = *openapi.NewNullableInt32(&val)
			} else {
				objProps.NumberOfMultipoints = *openapi.NewNullableInt32(nil)
			}
			objectPropertiesChanged = true
		}

		if !planOP.Aggregate.Equal(stateAggregate) {
			if !planOP.Aggregate.IsNull() {
				objProps.Aggregate = openapi.PtrBool(planOP.Aggregate.ValueBool())
			}
			objectPropertiesChanged = true
		}

		if !planOP.IsHost.Equal(stateIsHost) {
			if !planOP.IsHost.IsNull() {
				objProps.IsHost = openapi.PtrBool(planOP.IsHost.ValueBool())
			}
			objectPropertiesChanged = true
		}

		if !planOP.DrawAsEdgeDevice.Equal(stateDrawAsEdgeDevice) {
			if !planOP.DrawAsEdgeDevice.IsNull() {
				objProps.DrawAsEdgeDevice = openapi.PtrBool(planOP.DrawAsEdgeDevice.ValueBool())
			}
			objectPropertiesChanged = true
		}

		// Handle object_properties.eths
		var stateEths []veritySwitchpointObjectPropertiesEthModel
		if stateOP != nil {
			stateEths = stateOP.Eths
		}

		changedEths, ethsChanged := utils.ProcessIndexedArrayUpdates(planOP.Eths, stateEths,
			utils.IndexedItemHandler[veritySwitchpointObjectPropertiesEthModel, openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner]{
				CreateNew: func(planItem veritySwitchpointObjectPropertiesEthModel) openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner {
					eth := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{}

					// Handle string fields
					utils.SetStringFields([]utils.StringFieldMapping{
						{FieldName: "EthNumIcon", APIField: &eth.EthNumIcon, TFValue: planItem.EthNumIcon},
						{FieldName: "EthNumLabel", APIField: &eth.EthNumLabel, TFValue: planItem.EthNumLabel},
					})

					// Handle int64 fields
					utils.SetInt64Fields([]utils.Int64FieldMapping{
						{FieldName: "Index", APIField: &eth.Index, TFValue: planItem.Index},
					})

					return eth
				},
				UpdateExisting: func(planItem veritySwitchpointObjectPropertiesEthModel, stateItem veritySwitchpointObjectPropertiesEthModel) (openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner, bool) {
					eth := openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{}
					fieldChanged := false

					// Handle string field changes
					utils.CompareAndSetStringField(planItem.EthNumIcon, stateItem.EthNumIcon, func(v *string) { eth.EthNumIcon = v }, &fieldChanged)
					utils.CompareAndSetStringField(planItem.EthNumLabel, stateItem.EthNumLabel, func(v *string) { eth.EthNumLabel = v }, &fieldChanged)

					// Handle index field change
					utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { eth.Index = v }, &fieldChanged)

					return eth, fieldChanged
				},
				CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner {
					return openapi.SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner{
						Index: openapi.PtrInt32(int32(index)),
					}
				},
			})

		if ethsChanged {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "switchpoint", name, spProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "switchpoint", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")
	resp.State.RemoveResource(ctx)
}

func (r *veritySwitchpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
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
		"spine_plane":                &state.SpinePlane,
		"spine_plane_ref_type_":      &state.SpinePlaneRefType,
		"pod":                        &state.Pod,
		"pod_ref_type_":              &state.PodRefType,
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
			case "spine_plane":
				if !plan.SpinePlane.IsNull() {
					*stateField = plan.SpinePlane
				}
			case "spine_plane_ref_type_":
				if !plan.SpinePlaneRefType.IsNull() {
					*stateField = plan.SpinePlaneRefType
				}
			case "pod":
				if !plan.Pod.IsNull() {
					*stateField = plan.Pod
				}
			case "pod_ref_type_":
				if !plan.PodRefType.IsNull() {
					*stateField = plan.PodRefType
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
	} else if plan != nil && !plan.SwitchRouterIdIpMask.IsNull() && !plan.SwitchRouterIdIpMask.IsUnknown() {
		state.SwitchRouterIdIpMask = plan.SwitchRouterIdIpMask
	}

	if val, ok := switchpointData["switch_router_id_ip_mask_auto_assigned_"].(bool); ok {
		state.SwitchRouterIdIpMaskAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
		state.SwitchRouterIdIpMaskAutoAssigned = plan.SwitchRouterIdIpMaskAutoAssigned
	}

	if val, ok := switchpointData["switch_vtep_id_ip_mask"].(string); ok {
		state.SwitchVtepIdIpMask = types.StringValue(val)
	} else if plan != nil && !plan.SwitchVtepIdIpMask.IsNull() && !plan.SwitchVtepIdIpMask.IsUnknown() {
		state.SwitchVtepIdIpMask = plan.SwitchVtepIdIpMask
	}

	if val, ok := switchpointData["switch_vtep_id_ip_mask_auto_assigned_"].(bool); ok {
		state.SwitchVtepIdIpMaskAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
		state.SwitchVtepIdIpMaskAutoAssigned = plan.SwitchVtepIdIpMaskAutoAssigned
	}

	if val, ok := switchpointData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
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

	if val, ok := switchpointData["is_fabric"].(bool); ok {
		state.IsFabric = types.BoolValue(val)
	} else if plan != nil && !plan.IsFabric.IsNull() {
		state.IsFabric = plan.IsFabric
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
					if plan != nil && !plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() {
						state.BgpAsNumber = plan.BgpAsNumber
					}
				}
			default:
				if plan != nil && !plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() {
					state.BgpAsNumber = plan.BgpAsNumber
				}
			}
		}
	} else if plan != nil && !plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() {
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
			if index, ok := mirror["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					mirrorModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					mirrorModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int32); ok {
					mirrorModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int64); ok {
					mirrorModel.Index = types.Int64Value(intVal)
				} else {
					mirrorModel.Index = types.Int64Null()
				}
			} else {
				mirrorModel.Index = types.Int64Null()
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
		if val, ok := objProps["draw_as_edge_device"].(bool); ok {
			op.DrawAsEdgeDevice = types.BoolValue(val)
		} else {
			op.DrawAsEdgeDevice = types.BoolNull()
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

	// Validate auto-assigned field specifications in configuration when auto-assigned
	// Check the actual configuration, not the plan
	var config veritySwitchpointResourceModel
	if !req.Config.Raw.IsNull() {
		resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !config.BgpAsNumberAutoAssigned.IsNull() && config.BgpAsNumberAutoAssigned.ValueBool() {
			if !config.BgpAsNumber.IsNull() && !config.BgpAsNumber.IsUnknown() {
				resp.Diagnostics.AddError(
					"BGP AS Number cannot be specified when auto-assigned",
					"The 'bgp_as_number' field cannot be specified in the configuration when 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}

		if !config.SwitchRouterIdIpMaskAutoAssigned.IsNull() && config.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
			if !config.SwitchRouterIdIpMask.IsNull() && !config.SwitchRouterIdIpMask.IsUnknown() && config.SwitchRouterIdIpMask.ValueString() != "" {
				resp.Diagnostics.AddError(
					"Switch Router ID IP Mask cannot be specified when auto-assigned",
					"The 'switch_router_id_ip_mask' field cannot be specified in the configuration when 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}

		if !config.SwitchVtepIdIpMaskAutoAssigned.IsNull() && config.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
			if !config.SwitchVtepIdIpMask.IsNull() && !config.SwitchVtepIdIpMask.IsUnknown() && config.SwitchVtepIdIpMask.ValueString() != "" {
				resp.Diagnostics.AddError(
					"Switch VTEP ID IP Mask cannot be specified when auto-assigned",
					"The 'switch_vtep_id_ip_mask' field cannot be specified in the configuration when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}
	}

	// For new resources (where state is null), mark auto-assigned fields as Unknown
	if req.State.Raw.IsNull() {
		if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), types.Int64Unknown())
		}
		if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), types.StringUnknown())
		}
		if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), types.StringUnknown())
		}
		return
	}

	var state veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle auto-assigned field behavior
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		if !plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
			// bgp_as_number_auto_assigned_ is changing to true, API will assign the value
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"BGP AS Number will be assigned by the API",
				"The 'bgp_as_number' field will be automatically assigned by the API because 'bgp_as_number_auto_assigned_' is being set to true.",
			)
		} else if !plan.BgpAsNumber.Equal(state.BgpAsNumber) {
			// User tried to change BgpAsNumber but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring bgp_as_number changes with auto-assignment enabled",
				"The 'bgp_as_number' field changes will be ignored because 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.BgpAsNumber.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), state.BgpAsNumber)
			}
		}
	}

	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
			// switch_router_id_ip_mask_auto_assigned_ is changing to true, API will assign the value
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch Router ID IP Mask will be assigned by the API",
				"The 'switch_router_id_ip_mask' field will be automatically assigned by the API because 'switch_router_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) {
			// User tried to change SwitchRouterIdIpMask but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring switch_router_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_router_id_ip_mask' field changes will be ignored because 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.SwitchRouterIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), state.SwitchRouterIdIpMask)
			}
		}
	}

	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
			// switch_vtep_id_ip_mask_auto_assigned_ is changing to true, API will assign the value
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch VTEP ID IP Mask will be assigned by the API",
				"The 'switch_vtep_id_ip_mask' field will be automatically assigned by the API because 'switch_vtep_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) {
			// User tried to change SwitchVtepIdIpMask but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring switch_vtep_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_vtep_id_ip_mask' field changes will be ignored because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.SwitchVtepIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), state.SwitchVtepIdIpMask)
			}
		}
	}
}

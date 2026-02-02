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
	_ resource.Resource                = &veritySwitchpointResource{}
	_ resource.ResourceWithConfigure   = &veritySwitchpointResource{}
	_ resource.ResourceWithImportState = &veritySwitchpointResource{}
	_ resource.ResourceWithModifyPlan  = &veritySwitchpointResource{}
)

const switchpointResourceType = "switchpoints"
const switchpointTerraformType = "verity_switchpoint"

func NewVeritySwitchpointResource() resource.Resource {
	return &veritySwitchpointResource{}
}

type veritySwitchpointResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
	Breakout    types.String `tfsdk:"breakout"`
	Index       types.Int64  `tfsdk:"index"`
	EthNumIcon  types.String `tfsdk:"eth_num_icon"`
	EthNumLabel types.String `tfsdk:"eth_num_label"`
	Enable      types.Bool   `tfsdk:"enable"`
	PortName    types.String `tfsdk:"port_name"`
}

func (e veritySwitchpointEthModel) GetIndex() types.Int64 {
	return e.Index
}

type veritySwitchpointObjectPropertiesModel struct {
	UserNotes                     types.String `tfsdk:"user_notes"`
	ExpectedParentEndpoint        types.String `tfsdk:"expected_parent_endpoint"`
	ExpectedParentEndpointRefType types.String `tfsdk:"expected_parent_endpoint_ref_type_"`
	NumberOfMultipoints           types.Int64  `tfsdk:"number_of_multipoints"`
	Aggregate                     types.Bool   `tfsdk:"aggregate"`
	IsHost                        types.Bool   `tfsdk:"is_host"`
	DrawAsEdgeDevice              types.Bool   `tfsdk:"draw_as_edge_device"`
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
				Computed:    true,
			},
			"device_serial_number": schema.StringAttribute{
				Description: "Device Serial Number",
				Optional:    true,
				Computed:    true,
			},
			"connected_bundle": schema.StringAttribute{
				Description: "Connected Bundle",
				Optional:    true,
				Computed:    true,
			},
			"connected_bundle_ref_type_": schema.StringAttribute{
				Description: "Object type for connected_bundle field",
				Optional:    true,
				Computed:    true,
			},
			"read_only_mode": schema.BoolAttribute{
				Description: "When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware",
				Optional:    true,
				Computed:    true,
			},
			"locked": schema.BoolAttribute{
				Description: "Permission lock",
				Optional:    true,
				Computed:    true,
			},
			"out_of_band_management": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch is managed out of band via the management port",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of Switchpoint",
				Optional:    true,
				Computed:    true,
			},
			"spine_plane": schema.StringAttribute{
				Description: "Spine Plane - subgrouping of super spine and spine",
				Optional:    true,
				Computed:    true,
			},
			"spine_plane_ref_type_": schema.StringAttribute{
				Description: "Object type for spine_plane field",
				Optional:    true,
				Computed:    true,
			},
			"pod": schema.StringAttribute{
				Description: "Pod - subgrouping of spine and leaf switches",
				Optional:    true,
				Computed:    true,
			},
			"pod_ref_type_": schema.StringAttribute{
				Description: "Object type for pod field",
				Optional:    true,
				Computed:    true,
			},
			"rack": schema.StringAttribute{
				Description: "Physical Rack location of the Switch",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask": schema.StringAttribute{
				Description: "Switch BGP Router Identifier. This field should not be specified when 'switch_router_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch BGP Router Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_router_id_ip_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask": schema.StringAttribute{
				Description: "Switch VTEP Identifier. This field should not be specified when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch VTEP Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_vtep_id_ip_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number": schema.Int64Attribute{
				Description: "BGP Autonomous System Number for the site underlay. This field should not be specified when 'bgp_as_number_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the BGP AS Number should be automatically assigned by the API. When set to true, do not specify the 'bgp_as_number' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"is_fabric": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"badge_ref_type_": schema.StringAttribute{
							Description: "Object type for badge field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"child_num_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for child_num_endpoint field",
							Optional:    true,
							Computed:    true,
						},
						"child_num_device": schema.StringAttribute{
							Description: "Device associated with the Child",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"traffic_mirror_num_source_port": schema.StringAttribute{
							Description: "Source Port for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_source_lag_indicator": schema.BoolAttribute{
							Description: "Source LAG Indicator for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_destination_port": schema.StringAttribute{
							Description: "Destination Port for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_inbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for inbound traffic",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_outbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for outbound traffic",
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
			"eths": schema.ListNestedBlock{
				Description: "Ethernet port configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"breakout": schema.StringAttribute{
							Description: "Breakout Port Override",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
						"eth_num_icon": schema.StringAttribute{
							Description: "Icon of this Eth Port",
							Optional:    true,
							Computed:    true,
						},
						"eth_num_label": schema.StringAttribute{
							Description: "Label of this Eth Port",
							Optional:    true,
							Computed:    true,
						},
						"enable": schema.BoolAttribute{
							Description: "Enable port",
							Optional:    true,
							Computed:    true,
						},
						"port_name": schema.StringAttribute{
							Description: "The name identifying the port. Used for reference only, it won't actually change the port name.",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"expected_parent_endpoint": schema.StringAttribute{
							Description: "Expected Parent Endpoint",
							Optional:    true,
							Computed:    true,
						},
						"expected_parent_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for expected_parent_endpoint field",
							Optional:    true,
							Computed:    true,
						},
						"number_of_multipoints": schema.Int64Attribute{
							Description: "Number of Multipoints",
							Optional:    true,
							Computed:    true,
						},
						"aggregate": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes switch aggregated with all of its sub switches",
							Optional:    true,
							Computed:    true,
						},
						"is_host": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes the Host Switch",
							Optional:    true,
							Computed:    true,
						},
						"draw_as_edge_device": schema.BoolAttribute{
							Description: "Turn on to display the switch as an edge device instead of as a switch",
							Optional:    true,
							Computed:    true,
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

	var config veritySwitchpointResourceModel
	diags = req.Config.Get(ctx, &config)
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

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "BgpAsNumber", APIField: &spProps.BgpAsNumber, TFValue: config.BgpAsNumber, IsConfigured: configuredAttrs.IsConfigured("bgp_as_number")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		configOp, objPropsCfg := utils.GetObjectPropertiesConfig(op, config.ObjectProperties, configuredAttrs)
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "UserNotes", TFValue: op.UserNotes, APIValue: &objProps.UserNotes},
			{Name: "ExpectedParentEndpoint", TFValue: op.ExpectedParentEndpoint, APIValue: &objProps.ExpectedParentEndpoint},
			{Name: "ExpectedParentEndpointRefType", TFValue: op.ExpectedParentEndpointRefType, APIValue: &objProps.ExpectedParentEndpointRefType},
			{Name: "Aggregate", TFValue: op.Aggregate, APIValue: &objProps.Aggregate},
			{Name: "IsHost", TFValue: op.IsHost, APIValue: &objProps.IsHost},
			{Name: "DrawAsEdgeDevice", TFValue: op.DrawAsEdgeDevice, APIValue: &objProps.DrawAsEdgeDevice},
		})
		utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
			{FieldName: "NumberOfMultipoints", APIField: &objProps.NumberOfMultipoints, TFValue: configOp.NumberOfMultipoints, IsConfigured: objPropsCfg.IsFieldConfigured("number_of_multipoints")},
		})
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
				{FieldName: "EthNumIcon", APIField: &ethItem.EthNumIcon, TFValue: eth.EthNumIcon},
				{FieldName: "EthNumLabel", APIField: &ethItem.EthNumLabel, TFValue: eth.EthNumLabel},
				{FieldName: "PortName", APIField: &ethItem.PortName, TFValue: eth.PortName},
			})

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &ethItem.Enable, TFValue: eth.Enable},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethItem.Index, TFValue: eth.Index},
			})

			eths[i] = ethItem
		}
		spProps.Eths = eths
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "switchpoint", name, *spProps, &resp.Diagnostics)
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
			state := populateSwitchpointState(ctx, minState, switchpointData, r.provCtx.mode)
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
			state = populateSwitchpointState(ctx, state, switchpointData, r.provCtx.mode)
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

	state = populateSwitchpointState(ctx, state, switchpointMap, r.provCtx.mode)
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

	// Get config for nullable field handling
	var config veritySwitchpointResourceModel
	req.Config.Get(ctx, &config)
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { spProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DeviceSerialNumber, state.DeviceSerialNumber, func(v *string) { spProps.DeviceSerialNumber = v }, &hasChanges)
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
				spProps.BgpAsNumber = *openapi.NewNullableInt32(&val)
			} else {
				spProps.BgpAsNumber = *openapi.NewNullableInt32(nil)
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
						spProps.BgpAsNumber = *openapi.NewNullableInt32(&val)
					} else if !state.BgpAsNumber.IsNull() {
						// Use current state BgpAsNumber if plan doesn't specify one
						val := int32(state.BgpAsNumber.ValueInt64())
						spProps.BgpAsNumber = *openapi.NewNullableInt32(&val)
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
					{FieldName: "EthNumIcon", APIField: &eth.EthNumIcon, TFValue: planItem.EthNumIcon},
					{FieldName: "EthNumLabel", APIField: &eth.EthNumLabel, TFValue: planItem.EthNumLabel},
					{FieldName: "PortName", APIField: &eth.PortName, TFValue: planItem.PortName},
				})

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &eth.Enable, TFValue: planItem.Enable},
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
				utils.CompareAndSetStringField(planItem.EthNumIcon, stateItem.EthNumIcon, func(v *string) { eth.EthNumIcon = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.EthNumLabel, stateItem.EthNumLabel, func(v *string) { eth.EthNumLabel = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PortName, stateItem.PortName, func(v *string) { eth.PortName = v }, &fieldChanged)

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { eth.Enable = v }, &fieldChanged)

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
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		// Get config for nullable field handling in object_properties
		configOp, objPropsCfg := utils.GetObjectPropertiesConfig(op, config.ObjectProperties, configuredAttrs)

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "UserNotes", PlanValue: op.UserNotes, StateValue: st.UserNotes, APIValue: &objProps.UserNotes},
			{Name: "ExpectedParentEndpoint", PlanValue: op.ExpectedParentEndpoint, StateValue: st.ExpectedParentEndpoint, APIValue: &objProps.ExpectedParentEndpoint},
			{Name: "ExpectedParentEndpointRefType", PlanValue: op.ExpectedParentEndpointRefType, StateValue: st.ExpectedParentEndpointRefType, APIValue: &objProps.ExpectedParentEndpointRefType},
			{Name: "Aggregate", PlanValue: op.Aggregate, StateValue: st.Aggregate, APIValue: &objProps.Aggregate},
			{Name: "IsHost", PlanValue: op.IsHost, StateValue: st.IsHost, APIValue: &objProps.IsHost},
			{Name: "DrawAsEdgeDevice", PlanValue: op.DrawAsEdgeDevice, StateValue: st.DrawAsEdgeDevice, APIValue: &objProps.DrawAsEdgeDevice},
		}, &objPropsChanged)

		// Handle nullable field in object_properties
		utils.CompareAndSetNullableInt64Field(configOp.NumberOfMultipoints, st.NumberOfMultipoints, objPropsCfg.IsFieldConfigured("number_of_multipoints"), func(v *openapi.NullableInt32) { objProps.NumberOfMultipoints = *v }, &objPropsChanged)

		if objPropsChanged {
			spProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "switchpoint", name, spProps, &resp.Diagnostics)
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
			state := populateSwitchpointState(ctx, minState, switchpointData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "switchpoint", name, nil, &resp.Diagnostics)
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

func populateSwitchpointState(ctx context.Context, state veritySwitchpointResourceModel, switchpointData map[string]interface{}, mode string) veritySwitchpointResourceModel {
	const resourceType = switchpointResourceType

	state.Name = utils.MapStringFromAPI(switchpointData["name"])

	// Int fields
	state.BgpAsNumber = utils.MapInt64WithMode(switchpointData, "bgp_as_number", resourceType, mode)

	// Bool fields
	state.Enable = utils.MapBoolWithMode(switchpointData, "enable", resourceType, mode)
	state.ReadOnlyMode = utils.MapBoolWithMode(switchpointData, "read_only_mode", resourceType, mode)
	state.Locked = utils.MapBoolWithMode(switchpointData, "locked", resourceType, mode)
	state.OutOfBandManagement = utils.MapBoolWithMode(switchpointData, "out_of_band_management", resourceType, mode)
	state.IsFabric = utils.MapBoolWithMode(switchpointData, "is_fabric", resourceType, mode)
	state.BgpAsNumberAutoAssigned = utils.MapBoolWithMode(switchpointData, "bgp_as_number_auto_assigned_", resourceType, mode)
	state.SwitchVtepIdIpMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "switch_vtep_id_ip_mask_auto_assigned_", resourceType, mode)
	state.SwitchRouterIdIpMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "switch_router_id_ip_mask_auto_assigned_", resourceType, mode)

	// String fields
	state.DeviceSerialNumber = utils.MapStringWithMode(switchpointData, "device_serial_number", resourceType, mode)
	state.ConnectedBundle = utils.MapStringWithMode(switchpointData, "connected_bundle", resourceType, mode)
	state.ConnectedBundleRefType = utils.MapStringWithMode(switchpointData, "connected_bundle_ref_type_", resourceType, mode)
	state.Type = utils.MapStringWithMode(switchpointData, "type", resourceType, mode)
	state.SpinePlane = utils.MapStringWithMode(switchpointData, "spine_plane", resourceType, mode)
	state.SpinePlaneRefType = utils.MapStringWithMode(switchpointData, "spine_plane_ref_type_", resourceType, mode)
	state.Pod = utils.MapStringWithMode(switchpointData, "pod", resourceType, mode)
	state.PodRefType = utils.MapStringWithMode(switchpointData, "pod_ref_type_", resourceType, mode)
	state.Rack = utils.MapStringWithMode(switchpointData, "rack", resourceType, mode)
	state.SwitchRouterIdIpMask = utils.MapStringWithMode(switchpointData, "switch_router_id_ip_mask", resourceType, mode)
	state.SwitchVtepIdIpMask = utils.MapStringWithMode(switchpointData, "switch_vtep_id_ip_mask", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := switchpointData["object_properties"].(map[string]interface{}); ok {
			op := veritySwitchpointObjectPropertiesModel{
				UserNotes:                     utils.MapStringWithModeNested(objProps, "user_notes", resourceType, "object_properties.user_notes", mode),
				ExpectedParentEndpoint:        utils.MapStringWithModeNested(objProps, "expected_parent_endpoint", resourceType, "object_properties.expected_parent_endpoint", mode),
				ExpectedParentEndpointRefType: utils.MapStringWithModeNested(objProps, "expected_parent_endpoint_ref_type_", resourceType, "object_properties.expected_parent_endpoint_ref_type_", mode),
				NumberOfMultipoints:           utils.MapInt64WithModeNested(objProps, "number_of_multipoints", resourceType, "object_properties.number_of_multipoints", mode),
				Aggregate:                     utils.MapBoolWithModeNested(objProps, "aggregate", resourceType, "object_properties.aggregate", mode),
				IsHost:                        utils.MapBoolWithModeNested(objProps, "is_host", resourceType, "object_properties.is_host", mode),
				DrawAsEdgeDevice:              utils.MapBoolWithModeNested(objProps, "draw_as_edge_device", resourceType, "object_properties.draw_as_edge_device", mode),
			}
			state.ObjectProperties = []veritySwitchpointObjectPropertiesModel{op}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle badges block
	if utils.FieldAppliesToMode(resourceType, "badges", mode) {
		if badgesArray, ok := switchpointData["badges"].([]interface{}); ok && len(badgesArray) > 0 {
			var badges []veritySwitchpointBadgeModel
			for _, b := range badgesArray {
				badge, ok := b.(map[string]interface{})
				if !ok {
					continue
				}
				badgeModel := veritySwitchpointBadgeModel{
					Badge:        utils.MapStringWithModeNested(badge, "badge", resourceType, "badges.badge", mode),
					BadgeRefType: utils.MapStringWithModeNested(badge, "badge_ref_type_", resourceType, "badges.badge_ref_type_", mode),
					Index:        utils.MapInt64WithModeNested(badge, "index", resourceType, "badges.index", mode),
				}
				badges = append(badges, badgeModel)
			}
			state.Badges = badges
		} else {
			state.Badges = nil
		}
	} else {
		state.Badges = nil
	}

	// Handle children block
	if utils.FieldAppliesToMode(resourceType, "children", mode) {
		if childrenArray, ok := switchpointData["children"].([]interface{}); ok && len(childrenArray) > 0 {
			var children []veritySwitchpointChildModel
			for _, c := range childrenArray {
				child, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				childModel := veritySwitchpointChildModel{
					ChildNumEndpoint:        utils.MapStringWithModeNested(child, "child_num_endpoint", resourceType, "children.child_num_endpoint", mode),
					ChildNumEndpointRefType: utils.MapStringWithModeNested(child, "child_num_endpoint_ref_type_", resourceType, "children.child_num_endpoint_ref_type_", mode),
					ChildNumDevice:          utils.MapStringWithModeNested(child, "child_num_device", resourceType, "children.child_num_device", mode),
					Index:                   utils.MapInt64WithModeNested(child, "index", resourceType, "children.index", mode),
				}
				children = append(children, childModel)
			}
			state.Children = children
		} else {
			state.Children = nil
		}
	} else {
		state.Children = nil
	}

	// Handle traffic_mirrors block
	if utils.FieldAppliesToMode(resourceType, "traffic_mirrors", mode) {
		if mirrorsArray, ok := switchpointData["traffic_mirrors"].([]interface{}); ok && len(mirrorsArray) > 0 {
			var mirrors []veritySwitchpointTrafficMirrorModel
			for _, m := range mirrorsArray {
				mirror, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				mirrorModel := veritySwitchpointTrafficMirrorModel{
					TrafficMirrorNumEnable:             utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_enable", resourceType, "traffic_mirrors.traffic_mirror_num_enable", mode),
					TrafficMirrorNumSourcePort:         utils.MapStringWithModeNested(mirror, "traffic_mirror_num_source_port", resourceType, "traffic_mirrors.traffic_mirror_num_source_port", mode),
					TrafficMirrorNumSourceLagIndicator: utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_source_lag_indicator", resourceType, "traffic_mirrors.traffic_mirror_num_source_lag_indicator", mode),
					TrafficMirrorNumDestinationPort:    utils.MapStringWithModeNested(mirror, "traffic_mirror_num_destination_port", resourceType, "traffic_mirrors.traffic_mirror_num_destination_port", mode),
					TrafficMirrorNumInboundTraffic:     utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_inbound_traffic", resourceType, "traffic_mirrors.traffic_mirror_num_inbound_traffic", mode),
					TrafficMirrorNumOutboundTraffic:    utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_outbound_traffic", resourceType, "traffic_mirrors.traffic_mirror_num_outbound_traffic", mode),
					Index:                              utils.MapInt64WithModeNested(mirror, "index", resourceType, "traffic_mirrors.index", mode),
				}
				mirrors = append(mirrors, mirrorModel)
			}
			state.TrafficMirrors = mirrors
		} else {
			state.TrafficMirrors = nil
		}
	} else {
		state.TrafficMirrors = nil
	}

	// Handle eths block
	if utils.FieldAppliesToMode(resourceType, "eths", mode) {
		if ethsArray, ok := switchpointData["eths"].([]interface{}); ok && len(ethsArray) > 0 {
			var eths []veritySwitchpointEthModel
			for _, e := range ethsArray {
				eth, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				ethModel := veritySwitchpointEthModel{
					Breakout:    utils.MapStringWithModeNested(eth, "breakout", resourceType, "eths.breakout", mode),
					Index:       utils.MapInt64WithModeNested(eth, "index", resourceType, "eths.index", mode),
					EthNumIcon:  utils.MapStringWithModeNested(eth, "eth_num_icon", resourceType, "eths.eth_num_icon", mode),
					EthNumLabel: utils.MapStringWithModeNested(eth, "eth_num_label", resourceType, "eths.eth_num_label", mode),
					Enable:      utils.MapBoolWithModeNested(eth, "enable", resourceType, "eths.enable", mode),
					PortName:    utils.MapStringWithModeNested(eth, "port_name", resourceType, "eths.port_name", mode),
				}
				eths = append(eths, ethModel)
			}
			state.Eths = eths
		} else {
			state.Eths = nil
		}
	} else {
		state.Eths = nil
	}

	return state
}

func (r *veritySwitchpointResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = switchpointResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"device_serial_number", "connected_bundle", "connected_bundle_ref_type_",
		"type", "spine_plane", "spine_plane_ref_type_",
		"pod", "pod_ref_type_", "rack",
		"switch_router_id_ip_mask", "switch_vtep_id_ip_mask",
	)

	nullifier.NullifyBools(
		"enable", "read_only_mode", "locked",
		"out_of_band_management", "is_fabric",
		"switch_router_id_ip_mask_auto_assigned_",
		"switch_vtep_id_ip_mask_auto_assigned_",
		"bgp_as_number_auto_assigned_",
	)

	nullifier.NullifyInt64s(
		"bgp_as_number",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "badges",
		ItemCount:    len(plan.Badges),
		StringFields: []string{"badge", "badge_ref_type_"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "children",
		ItemCount:    len(plan.Children),
		StringFields: []string{"child_num_endpoint", "child_num_endpoint_ref_type_", "child_num_device"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "traffic_mirrors",
		ItemCount:    len(plan.TrafficMirrors),
		StringFields: []string{"traffic_mirror_num_source_port", "traffic_mirror_num_destination_port"},
		BoolFields:   []string{"traffic_mirror_num_enable", "traffic_mirror_num_source_lag_indicator", "traffic_mirror_num_inbound_traffic", "traffic_mirror_num_outbound_traffic"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "eths",
		ItemCount:    len(plan.Eths),
		StringFields: []string{"breakout", "eth_num_icon", "eth_num_label", "port_name"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"user_notes", "expected_parent_endpoint", "expected_parent_endpoint_ref_type_"},
		BoolFields:   []string{"aggregate", "is_host", "draw_as_edge_device"},
		Int64Fields:  []string{"number_of_multipoints"},
	})

	// =========================================================================
	// CREATE operation - handle auto-assigned fields
	// =========================================================================
	if req.State.Raw.IsNull() {
		// Switchpoint-specific: auto-assignment on create
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

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySwitchpointResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "bgp_as_number", ConfigVal: config.BgpAsNumber, StateVal: state.BgpAsNumber},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	if len(config.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		configOp := config.ObjectProperties[0]
		stateOp := state.ObjectProperties[0]

		if configuredAttrs.IsBlockAttributeConfigured("object_properties.number_of_multipoints") &&
			configOp.NumberOfMultipoints.IsNull() && !stateOp.NumberOfMultipoints.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("object_properties").AtListIndex(0).AtName("number_of_multipoints"), types.Int64Null())
		}
	}

	// =========================================================================
	// Validate auto-assigned field specifications
	// =========================================================================
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

	// =========================================================================
	// Resource-specific auto-assigned field logic (BgpAsNumber)
	// =========================================================================
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		if !plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
			// bgp_as_number_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"BGP AS Number will be assigned by the API",
				"The 'bgp_as_number' field will be automatically assigned by the API because 'bgp_as_number_auto_assigned_' is being set to true.",
			)
		} else if !plan.BgpAsNumber.Equal(state.BgpAsNumber) {
			// User tried to change BgpAsNumber but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring bgp_as_number changes with auto-assignment enabled",
				"The 'bgp_as_number' field changes will be ignored because 'bgp_as_number_auto_assigned_' is set to true.",
			)
			if !state.BgpAsNumber.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), state.BgpAsNumber)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (SwitchRouterIdIpMask)
	// =========================================================================
	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
			// switch_router_id_ip_mask_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch Router ID IP Mask will be assigned by the API",
				"The 'switch_router_id_ip_mask' field will be automatically assigned by the API because 'switch_router_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) {
			// User tried to change SwitchRouterIdIpMask but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring switch_router_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_router_id_ip_mask' field changes will be ignored because 'switch_router_id_ip_mask_auto_assigned_' is set to true.",
			)
			if !state.SwitchRouterIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), state.SwitchRouterIdIpMask)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (SwitchVtepIdIpMask)
	// =========================================================================
	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
			// switch_vtep_id_ip_mask_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch VTEP ID IP Mask will be assigned by the API",
				"The 'switch_vtep_id_ip_mask' field will be automatically assigned by the API because 'switch_vtep_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) {
			// User tried to change SwitchVtepIdIpMask but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring switch_vtep_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_vtep_id_ip_mask' field changes will be ignored because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true.",
			)
			if !state.SwitchVtepIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), state.SwitchVtepIdIpMask)
			}
		}
	}
}

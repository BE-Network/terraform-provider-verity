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
	_ resource.Resource                = &verityGatewayProfileResource{}
	_ resource.ResourceWithConfigure   = &verityGatewayProfileResource{}
	_ resource.ResourceWithImportState = &verityGatewayProfileResource{}
)

func NewVerityGatewayProfileResource() resource.Resource {
	return &verityGatewayProfileResource{}
}

type verityGatewayProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityGatewayProfileResourceModel struct {
	Name               types.String            `tfsdk:"name"`
	Enable             types.Bool              `tfsdk:"enable"`
	TenantSliceManaged types.Bool              `tfsdk:"tenant_slice_managed"`
	ObjectProperties   []objectPropertiesModel `tfsdk:"object_properties"`
	ExternalGateways   []externalGatewaysModel `tfsdk:"external_gateways"`
}

type objectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

type externalGatewaysModel struct {
	Enable         types.Bool   `tfsdk:"enable"`
	Gateway        types.String `tfsdk:"gateway"`
	GatewayRefType types.String `tfsdk:"gateway_ref_type_"`
	SourceIpMask   types.String `tfsdk:"source_ip_mask"`
	PeerGw         types.Bool   `tfsdk:"peer_gw"`
	Index          types.Int64  `tfsdk:"index"`
}

func (r *verityGatewayProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway_profile"
}

func (r *verityGatewayProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *verityGatewayProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Gateway Profile",
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
			"tenant_slice_managed": schema.BoolAttribute{
				Description: "Profiles that Tenant Slice creates and manages",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the gateway profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
					},
				},
			},
			"external_gateways": schema.ListNestedBlock{
				Description: "List of external gateway configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "BGP Gateway referenced for port profile",
							Optional:    true,
						},
						"gateway_ref_type_": schema.StringAttribute{
							Description: "Object type for gateway field",
							Optional:    true,
						},
						"source_ip_mask": schema.StringAttribute{
							Description: "Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session",
							Optional:    true,
						},
						"peer_gw": schema.BoolAttribute{
							Description: "Setting for paired switches only. Flag indicating that this gateway is a peer gateway.",
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

func (r *verityGatewayProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data verityGatewayProfileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	name := data.Name.ValueString()
	profileObj := openapi.NewConfigPutRequestGatewayProfileGatewayProfileName()
	profileObj.Name = openapi.PtrString(name)

	if !data.Enable.IsNull() {
		enable := data.Enable.ValueBool()
		profileObj.SetEnable(enable)
	}

	if !data.TenantSliceManaged.IsNull() {
		tenantSliceManaged := data.TenantSliceManaged.ValueBool()
		profileObj.SetTenantSliceManaged(tenantSliceManaged)
	}

	if len(data.ObjectProperties) > 0 {
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}

		if !data.ObjectProperties[0].Group.IsNull() {
			groupVal := data.ObjectProperties[0].Group.ValueString()
			objProps.Group = openapi.PtrString(groupVal)
		} else {
			objProps.Group = nil
		}

		profileObj.ObjectProperties = &objProps
	}

	if len(data.ExternalGateways) > 0 {
		var externalGatewaysList []openapi.ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner

		for _, eg := range data.ExternalGateways {
			gatewayObj := openapi.NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner()

			if !eg.Enable.IsNull() {
				gatewayObj.SetEnable(eg.Enable.ValueBool())
			}

			if !eg.Gateway.IsNull() {
				gatewayObj.SetGateway(eg.Gateway.ValueString())
			}

			if !eg.GatewayRefType.IsNull() && eg.GatewayRefType.ValueString() != "" {
				gatewayObj.SetGatewayRefType(eg.GatewayRefType.ValueString())
			}

			if !eg.SourceIpMask.IsNull() {
				gatewayObj.SetSourceIpMask(eg.SourceIpMask.ValueString())
			}

			if !eg.PeerGw.IsNull() {
				gatewayObj.SetPeerGw(eg.PeerGw.ValueBool())
			}

			if !eg.Index.IsNull() {
				index := int32(eg.Index.ValueInt64())
				gatewayObj.SetIndex(index)
			}

			externalGatewaysList = append(externalGatewaysList, *gatewayObj)
		}

		profileObj.SetExternalGateways(externalGatewaysList)
	}

	operationID := r.bulkOpsMgr.AddGatewayProfilePut(ctx, name, *profileObj)
	r.notifyOperationAdded()

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Gateway Profile",
			fmt.Sprintf("Could not create gateway profile %s: %s", name, err),
		)
		return
	}

	clearCache(ctx, r.provCtx, "gateway_profiles")
	data.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityGatewayProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data verityGatewayProfileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	profileName := data.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentGatewayProfileOperations() {
		return
	}

	type GatewayProfileResponse struct {
		GatewayProfile map[string]map[string]interface{} `json:"gateway_profile"`
	}

	var result GatewayProfileResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch gateway profiles on attempt %d, retrying in %v", attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		profilesData, fetchErr := getCachedResponse(ctx, r.provCtx, "gateway_profiles", func() (interface{}, error) {
			req := r.client.GatewayProfilesAPI.GatewayprofilesGet(ctx)
			resp, err := req.Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading gateway profiles: %v", err)
			}
			defer resp.Body.Close()

			var result GatewayProfileResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("error decoding gateway profile response: %v", err)
			}

			return result, nil
		})

		if fetchErr == nil {
			result = profilesData.(GatewayProfileResponse)
			break
		}
		err = fetchErr
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Gateway Profile",
			fmt.Sprintf("Could not read gateway profile %s: %s", profileName, err),
		)
		return
	}

	profileData, exists := result.GatewayProfile[profileName]
	if !exists {
		for apiName, p := range result.GatewayProfile {
			if name, ok := p["name"].(string); ok && name == profileName {
				profileData = p
				profileName = apiName
				exists = true
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Gateway profile with ID '%s' not found in API response", profileName))
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(profileName)

	if v, ok := profileData["enable"].(bool); ok {
		data.Enable = types.BoolValue(v)
	}

	if v, ok := profileData["tenant_slice_managed"].(bool); ok {
		data.TenantSliceManaged = types.BoolValue(v)
	}

	if objProps, ok := profileData["object_properties"].(map[string]interface{}); ok {
		objectProps := objectPropertiesModel{}

		if group, ok := objProps["group"]; ok && group != nil {
			groupStr, isString := group.(string)
			if isString {
				objectProps.Group = types.StringValue(groupStr)
			}
		}

		data.ObjectProperties = []objectPropertiesModel{objectProps}
	}

	if ext, ok := profileData["external_gateways"].([]interface{}); ok {
		var egList []externalGatewaysModel

		for _, item := range ext {
			if m, ok := item.(map[string]interface{}); ok {
				gateway := externalGatewaysModel{}

				if v, exists := m["enable"]; exists && v != nil {
					if boolVal, ok := v.(bool); ok {
						gateway.Enable = types.BoolValue(boolVal)
					}
				}

				if v, exists := m["gateway"]; exists && v != nil {
					if strVal, ok := v.(string); ok {
						gateway.Gateway = types.StringValue(strVal)
					}
				}

				if v, exists := m["gateway_ref_type_"]; exists && v != nil {
					if strVal, ok := v.(string); ok {
						gateway.GatewayRefType = types.StringValue(strVal)
					}
				}

				if v, exists := m["source_ip_mask"]; exists && v != nil {
					if strVal, ok := v.(string); ok {
						gateway.SourceIpMask = types.StringValue(strVal)
					}
				}

				if v, exists := m["peer_gw"]; exists && v != nil {
					if boolVal, ok := v.(bool); ok {
						gateway.PeerGw = types.BoolValue(boolVal)
					}
				}

				if v, exists := m["index"]; exists && v != nil {
					var indexVal int64
					switch val := v.(type) {
					case float64:
						indexVal = int64(val)
					case int:
						indexVal = int64(val)
					case int64:
						indexVal = val
					case float32:
						indexVal = int64(val)
					case int32:
						indexVal = int64(val)
					}
					gateway.Index = types.Int64Value(indexVal)
				}

				egList = append(egList, gateway)
			}
		}

		data.ExternalGateways = egList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityGatewayProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data verityGatewayProfileResourceModel
	var state verityGatewayProfileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	name := data.Name.ValueString()
	hasChanges := false
	profileObj := openapi.ConfigPutRequestGatewayProfileGatewayProfileName{}

	if !data.Enable.Equal(state.Enable) {
		enable := data.Enable.ValueBool()
		profileObj.Enable = &enable
		hasChanges = true
	}

	if !data.TenantSliceManaged.Equal(state.TenantSliceManaged) {
		tenantSliceManaged := data.TenantSliceManaged.ValueBool()
		profileObj.TenantSliceManaged = &tenantSliceManaged
		hasChanges = true
	}

	if len(data.ObjectProperties) > 0 {
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
		objPropsChanged := false

		if len(state.ObjectProperties) == 0 ||
			!data.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objPropsChanged = true

			if !data.ObjectProperties[0].Group.IsNull() {
				groupVal := data.ObjectProperties[0].Group.ValueString()
				objProps.Group = openapi.PtrString(groupVal)
			} else {
				objProps.Group = nil
			}
		}

		if objPropsChanged {
			profileObj.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if len(data.ExternalGateways) > 0 {
		stateGatewaysByIndex := make(map[int64]externalGatewaysModel)
		for _, eg := range state.ExternalGateways {
			if !eg.Index.IsNull() {
				stateGatewaysByIndex[eg.Index.ValueInt64()] = eg
			}
		}

		var changedExternalGateways []openapi.ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner

		for _, eg := range data.ExternalGateways {
			if eg.Index.IsNull() {
				continue
			}

			index := eg.Index.ValueInt64()
			stateEg, exists := stateGatewaysByIndex[index]
			gatewayChanged := false

			if exists {
				if !eg.Enable.Equal(stateEg.Enable) ||
					!eg.Gateway.Equal(stateEg.Gateway) ||
					!eg.GatewayRefType.Equal(stateEg.GatewayRefType) ||
					!eg.SourceIpMask.Equal(stateEg.SourceIpMask) ||
					!eg.PeerGw.Equal(stateEg.PeerGw) {
					gatewayChanged = true
				}
			} else {
				gatewayChanged = true
			}

			if gatewayChanged {
				gateway := openapi.ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner{
					Index: openapi.PtrInt32(int32(index)),
				}

				if !eg.Enable.IsNull() {
					gateway.Enable = openapi.PtrBool(eg.Enable.ValueBool())
				} else {
					gateway.Enable = openapi.PtrBool(false)
				}

				if !eg.Gateway.IsNull() {
					gateway.Gateway = openapi.PtrString(eg.Gateway.ValueString())
				} else {
					gateway.Gateway = openapi.PtrString("")
				}

				if !eg.GatewayRefType.IsNull() {
					gateway.GatewayRefType = openapi.PtrString(eg.GatewayRefType.ValueString())
				} else {
					gateway.GatewayRefType = openapi.PtrString("")
				}

				if !eg.SourceIpMask.IsNull() {
					gateway.SourceIpMask = openapi.PtrString(eg.SourceIpMask.ValueString())
				} else {
					gateway.SourceIpMask = openapi.PtrString("")
				}

				if !eg.PeerGw.IsNull() {
					gateway.PeerGw = openapi.PtrBool(eg.PeerGw.ValueBool())
				} else {
					gateway.PeerGw = openapi.PtrBool(false)
				}

				changedExternalGateways = append(changedExternalGateways, gateway)
			}
		}

		if len(changedExternalGateways) > 0 {
			profileObj.ExternalGateways = changedExternalGateways
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	operationID := r.bulkOpsMgr.AddGatewayProfilePatch(ctx, name, profileObj)
	r.notifyOperationAdded()

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Gateway Profile",
			fmt.Sprintf("Could not update gateway profile %s: %s", name, err),
		)
		return
	}

	clearCache(ctx, r.provCtx, "gateway_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityGatewayProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data verityGatewayProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	name := data.Name.ValueString()
	operationID := r.bulkOpsMgr.AddGatewayProfileDelete(ctx, name)
	r.notifyOperationAdded()

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Gateway Profile",
			fmt.Sprintf("Could not delete gateway profile %s: %s", name, err),
		)
		return
	}

	clearCache(ctx, r.provCtx, "gateway_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityGatewayProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

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
	_ resource.Resource                = &verityEthPortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortProfileResource{}
	_ resource.ResourceWithImportState = &verityEthPortProfileResource{}
)

func NewVerityEthPortProfileResource() resource.Resource {
	return &verityEthPortProfileResource{}
}

type verityEthPortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityEthPortProfileResourceModel struct {
	Name               types.String                                `tfsdk:"name"`
	Enable             types.Bool                                  `tfsdk:"enable"`
	TenantSliceManaged types.Bool                                  `tfsdk:"tenant_slice_managed"`
	ObjectProperties   []verityEthPortProfileObjectPropertiesModel `tfsdk:"object_properties"`
	Services           []servicesModel                             `tfsdk:"services"`
}

type verityEthPortProfileObjectPropertiesModel struct {
	Group          types.String `tfsdk:"group"`
	PortMonitoring types.String `tfsdk:"port_monitoring"`
}

type servicesModel struct {
	RowNumEnable         types.Bool   `tfsdk:"row_num_enable"`
	RowNumService        types.String `tfsdk:"row_num_service"`
	RowNumServiceRefType types.String `tfsdk:"row_num_service_ref_type_"`
	RowNumExternalVlan   types.Int64  `tfsdk:"row_num_external_vlan"`
	Index                types.Int64  `tfsdk:"index"`
}

func (s servicesModel) GetIndex() types.Int64 {
	return s.Index
}

func (r *verityEthPortProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eth_port_profile"
}

func (r *verityEthPortProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityEthPortProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ethernet Port Profile",
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
			"tenant_slice_managed": schema.BoolAttribute{
				Description: "Profiles that Tenant Slice creates and manages",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the profile",
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
			"services": schema.ListNestedBlock{
				Description: "List of service configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"row_num_enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
						},
						"row_num_service": schema.StringAttribute{
							Description: "Choose a Service to connect",
							Optional:    true,
						},
						"row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_service field",
							Optional:    true,
						},
						"row_num_external_vlan": schema.Int64Attribute{
							Description: "Choose an external vlan. A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service.",
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

func (r *verityEthPortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityEthPortProfileResourceModel
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
	ethPortProfileProps := &openapi.EthportprofilesPutRequestEthPortProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ethPortProfileProps.Enable, TFValue: plan.Enable},
		{FieldName: "TenantSliceManaged", APIField: &ethPortProfileProps.TenantSliceManaged, TFValue: plan.TenantSliceManaged},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", TFValue: op.PortMonitoring, APIValue: &objProps.PortMonitoring},
		})
		ethPortProfileProps.ObjectProperties = &objProps
	}

	// Handle Services
	if len(plan.Services) > 0 {
		servicesItems := make([]openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, len(plan.Services))
		for i, item := range plan.Services {
			serviceItem := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowNumEnable", APIField: &serviceItem.RowNumEnable, TFValue: item.RowNumEnable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowNumService", APIField: &serviceItem.RowNumService, TFValue: item.RowNumService},
				{FieldName: "RowNumServiceRefType", APIField: &serviceItem.RowNumServiceRefType, TFValue: item.RowNumServiceRefType},
			})
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "RowNumExternalVlan", APIField: &serviceItem.RowNumExternalVlan, TFValue: item.RowNumExternalVlan},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &serviceItem.Index, TFValue: item.Index},
			})

			servicesItems[i] = serviceItem
		}
		ethPortProfileProps.Services = servicesItems
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_profile", name, *ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityEthPortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityEthPortProfileResourceModel
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

	ethPortProfileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("eth_port_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port profile %s verification – trusting recent successful API operation", ethPortProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching eth port profiles for verification of %s", ethPortProfileName))

	type EthPortProfileResponse struct {
		EthPortProfile map[string]interface{} `json:"eth_port_profile_"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "eth_port_profiles", ethPortProfileName,
		func() (EthPortProfileResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch eth port profiles")
			respAPI, err := r.client.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
			if err != nil {
				return EthPortProfileResponse{}, fmt.Errorf("error reading eth port profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res EthPortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return EthPortProfileResponse{}, fmt.Errorf("failed to decode eth port profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d eth port profiles", len(res.EthPortProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Eth Port Profile %s", ethPortProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for eth port profile with name: %s", ethPortProfileName))

	ethPortProfileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EthPortProfile,
		ethPortProfileName,
		func(data interface{}) (string, bool) {
			if ethPortProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := ethPortProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Eth Port Profile with name '%s' not found in API response", ethPortProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	ethPortProfileMap, ok := ethPortProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Eth Port Profile Data",
			fmt.Sprintf("Eth Port Profile data is not in expected format for %s", ethPortProfileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found eth port profile '%s' under API key '%s'", ethPortProfileName, actualAPIName))

	state.Name = utils.MapStringFromAPI(ethPortProfileMap["name"])

	// Handle object properties
	if objProps, ok := ethPortProfileMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityEthPortProfileObjectPropertiesModel{
			{
				Group:          utils.MapStringFromAPI(objProps["group"]),
				PortMonitoring: utils.MapStringFromAPI(objProps["port_monitoring"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":               &state.Enable,
		"tenant_slice_managed": &state.TenantSliceManaged,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ethPortProfileMap[apiKey])
	}

	// Handle services
	if services, ok := ethPortProfileMap["services"].([]interface{}); ok && len(services) > 0 {
		var servicesList []servicesModel

		for _, s := range services {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			serviceModel := servicesModel{
				RowNumEnable:         utils.MapBoolFromAPI(service["row_num_enable"]),
				RowNumService:        utils.MapStringFromAPI(service["row_num_service"]),
				RowNumServiceRefType: utils.MapStringFromAPI(service["row_num_service_ref_type_"]),
				RowNumExternalVlan:   utils.MapInt64FromAPI(service["row_num_external_vlan"]),
				Index:                utils.MapInt64FromAPI(service["index"]),
			}

			servicesList = append(servicesList, serviceModel)
		}

		state.Services = servicesList
	} else {
		state.Services = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityEthPortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityEthPortProfileResourceModel

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
	ethPortProfileProps := openapi.EthportprofilesPutRequestEthPortProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ethPortProfileProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ethPortProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TenantSliceManaged, state.TenantSliceManaged, func(v *bool) { ethPortProfileProps.TenantSliceManaged = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", PlanValue: op.PortMonitoring, StateValue: st.PortMonitoring, APIValue: &objProps.PortMonitoring},
		}, &objPropsChanged)

		if objPropsChanged {
			ethPortProfileProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle services
	servicesHandler := utils.IndexedItemHandler[servicesModel, openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner]{
		CreateNew: func(planItem servicesModel) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowNumEnable", APIField: &service.RowNumEnable, TFValue: planItem.RowNumEnable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowNumService", APIField: &service.RowNumService, TFValue: planItem.RowNumService},
				{FieldName: "RowNumServiceRefType", APIField: &service.RowNumServiceRefType, TFValue: planItem.RowNumServiceRefType},
			})

			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "RowNumExternalVlan", APIField: &service.RowNumExternalVlan, TFValue: planItem.RowNumExternalVlan},
			})

			return service
		},
		UpdateExisting: func(planItem servicesModel, stateItem servicesModel) (openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, bool) {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.RowNumEnable, stateItem.RowNumEnable, func(v *bool) { service.RowNumEnable = v }, &fieldChanged)

			// Handle nullable int64 fields
			utils.CompareAndSetNullableInt64Field(planItem.RowNumExternalVlan, stateItem.RowNumExternalVlan, func(v *openapi.NullableInt32) { service.RowNumExternalVlan = *v }, &fieldChanged)

			// Handle row_num_service and row_num_service_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowNumService, stateItem.RowNumService, planItem.RowNumServiceRefType, stateItem.RowNumServiceRefType,
				func(v *string) { service.RowNumService = v },
				func(v *string) { service.RowNumServiceRefType = v },
				"row_num_service", "row_num_service_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return service, false
			}

			return service, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: types.Int64Value(index)},
			})
			return service
		},
	}

	changedServices, servicesChanged := utils.ProcessIndexedArrayUpdates(plan.Services, state.Services, servicesHandler)
	if servicesChanged {
		ethPortProfileProps.Services = changedServices
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_profile", name, ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityEthPortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityEthPortProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

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
	_ resource.Resource                = &veritySfpBreakoutResource{}
	_ resource.ResourceWithConfigure   = &veritySfpBreakoutResource{}
	_ resource.ResourceWithImportState = &veritySfpBreakoutResource{}
	_ resource.ResourceWithModifyPlan  = &veritySfpBreakoutResource{}
)

const sfpBreakoutResourceType = "sfpbreakouts"

func NewVeritySfpBreakoutResource() resource.Resource {
	return &veritySfpBreakoutResource{}
}

type veritySfpBreakoutResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySfpBreakoutResourceModel struct {
	Name             types.String                            `tfsdk:"name"`
	Enable           types.Bool                              `tfsdk:"enable"`
	Breakout         []veritySfpBreakoutBreakoutModel        `tfsdk:"breakout"`
	ObjectProperties *veritySfpBreakoutObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySfpBreakoutBreakoutModel struct {
	Enable     types.Bool   `tfsdk:"enable"`
	Vendor     types.String `tfsdk:"vendor"`
	PartNumber types.String `tfsdk:"part_number"`
	Breakout   types.String `tfsdk:"breakout"`
	Index      types.Int64  `tfsdk:"index"`
}

func (m veritySfpBreakoutBreakoutModel) GetIndex() types.Int64 {
	return m.Index
}

type veritySfpBreakoutObjectPropertiesModel struct {
	// Empty object properties according to schema
}

func (r *veritySfpBreakoutResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sfp_breakout"
}

func (r *veritySfpBreakoutResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySfpBreakoutResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity SFP Breakout",
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
		},
		Blocks: map[string]schema.Block{
			"breakout": schema.ListNestedBlock{
				Description: "List of breakout configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "Vendor",
							Optional:    true,
							Computed:    true,
						},
						"part_number": schema.StringAttribute{
							Description: "Part Number",
							Optional:    true,
							Computed:    true,
						},
						"breakout": schema.StringAttribute{
							Description: "Breakout definition; defines number of ports of what speed this port is brokenout to.",
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
				Description: "Object properties.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						// No attributes defined - object_properties is an empty object in the schema
					},
				},
			},
		},
	}
}

func (r *veritySfpBreakoutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"SFP Breakout resources cannot be created. They represent existing hardware configurations that can only be read and updated.",
	)
}

func (r *veritySfpBreakoutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySfpBreakoutResourceModel
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

	sfpBreakoutName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if sfpBreakoutData, exists := r.bulkOpsMgr.GetResourceResponse("sfp_breakout", sfpBreakoutName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached sfp_breakout data for %s from recent operation", sfpBreakoutName))
			state = populateSfpBreakoutState(ctx, state, sfpBreakoutData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("sfp_breakout") {
		tflog.Info(ctx, fmt.Sprintf("Skipping SFP Breakout %s verification – trusting recent successful API operation", sfpBreakoutName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching SFP Breakouts for verification of %s", sfpBreakoutName))

	type SfpBreakoutsResponse struct {
		SfpBreakouts map[string]interface{} `json:"sfp_breakouts"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "sfp_breakouts", sfpBreakoutName,
		func() (SfpBreakoutsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch SFP Breakouts")
			respAPI, err := r.client.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
			if err != nil {
				return SfpBreakoutsResponse{}, fmt.Errorf("error reading SFP Breakouts: %v", err)
			}
			defer respAPI.Body.Close()

			var res SfpBreakoutsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SfpBreakoutsResponse{}, fmt.Errorf("failed to decode SFP Breakouts response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d SFP Breakouts", len(res.SfpBreakouts)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SFP Breakout %s", sfpBreakoutName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for SFP Breakout with name: %s", sfpBreakoutName))

	sfpBreakoutData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.SfpBreakouts,
		sfpBreakoutName,
		func(data interface{}) (string, bool) {
			if sfpBreakout, ok := data.(map[string]interface{}); ok {
				if name, ok := sfpBreakout["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("SFP Breakout with name '%s' not found in API response", sfpBreakoutName))
		resp.State.RemoveResource(ctx)
		return
	}

	sfpBreakoutMap, ok := sfpBreakoutData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid SFP Breakout Data",
			fmt.Sprintf("SFP Breakout data is not in expected format for %s", sfpBreakoutName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found SFP Breakout '%s' under API key '%s'", sfpBreakoutName, actualAPIName))

	state = populateSfpBreakoutState(ctx, state, sfpBreakoutMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySfpBreakoutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySfpBreakoutResourceModel

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
	sfpBreakoutProps := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { sfpBreakoutProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { sfpBreakoutProps.Enable = v }, &hasChanges)

	// Handle object properties
	if (plan.ObjectProperties == nil) != (state.ObjectProperties == nil) {
		if plan.ObjectProperties != nil {
			// SFP Breakout object properties are empty according to schema
			sfpBreakoutObjProps := map[string]interface{}{}
			sfpBreakoutProps.ObjectProperties = sfpBreakoutObjProps
		}
		hasChanges = true
	}

	// Handle breakout
	changedBreakouts, breakoutsChanged := utils.ProcessIndexedArrayUpdates(plan.Breakout, state.Breakout,
		utils.IndexedItemHandler[veritySfpBreakoutBreakoutModel, openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner]{
			CreateNew: func(planItem veritySfpBreakoutBreakoutModel) openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner {
				newBreakout := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &newBreakout.Enable, TFValue: planItem.Enable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Vendor", APIField: &newBreakout.Vendor, TFValue: planItem.Vendor},
					{FieldName: "PartNumber", APIField: &newBreakout.PartNumber, TFValue: planItem.PartNumber},
					{FieldName: "Breakout", APIField: &newBreakout.Breakout, TFValue: planItem.Breakout},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newBreakout.Index, TFValue: planItem.Index},
				})

				return newBreakout
			},
			UpdateExisting: func(planItem veritySfpBreakoutBreakoutModel, stateItem veritySfpBreakoutBreakoutModel) (openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner, bool) {
				updateBreakout := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { updateBreakout.Enable = v }, &fieldChanged)

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.Vendor, stateItem.Vendor, func(v *string) { updateBreakout.Vendor = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PartNumber, stateItem.PartNumber, func(v *string) { updateBreakout.PartNumber = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.Breakout, stateItem.Breakout, func(v *string) { updateBreakout.Breakout = v }, &fieldChanged)

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateBreakout.Index = v }, &fieldChanged)

				return updateBreakout, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner {
				return openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if breakoutsChanged {
		sfpBreakoutProps.Breakout = changedBreakouts
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "sfp_breakout", name, sfpBreakoutProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFP Breakout %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sfp_breakouts")

	var minState veritySfpBreakoutResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sfpBreakoutData, exists := bulkMgr.GetResourceResponse("sfp_breakout", name); exists {
			newState := populateSfpBreakoutState(ctx, minState, sfpBreakoutData, r.provCtx.mode)
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

func (r *veritySfpBreakoutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"SFP Breakout resources cannot be deleted. They represent existing hardware configurations that can only be read and updated.",
	)
}

func (r *veritySfpBreakoutResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSfpBreakoutState(ctx context.Context, state veritySfpBreakoutResourceModel, data map[string]interface{}, mode string) veritySfpBreakoutResourceModel {
	const resourceType = sfpBreakoutResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle breakout block
	if utils.FieldAppliesToMode(resourceType, "breakout", mode) {
		if breakoutData, ok := data["breakout"].([]interface{}); ok && len(breakoutData) > 0 {
			var breakouts []veritySfpBreakoutBreakoutModel

			for _, b := range breakoutData {
				breakout, ok := b.(map[string]interface{})
				if !ok {
					continue
				}

				breakoutModel := veritySfpBreakoutBreakoutModel{
					Enable:     utils.MapBoolWithModeNested(breakout, "enable", resourceType, "breakout.enable", mode),
					Vendor:     utils.MapStringWithModeNested(breakout, "vendor", resourceType, "breakout.vendor", mode),
					PartNumber: utils.MapStringWithModeNested(breakout, "part_number", resourceType, "breakout.part_number", mode),
					Breakout:   utils.MapStringWithModeNested(breakout, "breakout", resourceType, "breakout.breakout", mode),
					Index:      utils.MapInt64WithModeNested(breakout, "index", resourceType, "breakout.index", mode),
				}

				breakouts = append(breakouts, breakoutModel)
			}

			state.Breakout = breakouts
		} else {
			state.Breakout = nil
		}
	} else {
		state.Breakout = nil
	}

	return state
}

func (r *veritySfpBreakoutResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySfpBreakoutResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = sfpBreakoutResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "breakout",
		ItemCount:    len(plan.Breakout),
		StringFields: []string{"vendor", "part_number", "breakout"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})
}

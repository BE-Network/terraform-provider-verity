package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	_ resource.Resource                = &veritySfpBreakoutResource{}
	_ resource.ResourceWithConfigure   = &veritySfpBreakoutResource{}
	_ resource.ResourceWithImportState = &veritySfpBreakoutResource{}
)

func NewVeritySfpBreakoutResource() resource.Resource {
	return &veritySfpBreakoutResource{}
}

type veritySfpBreakoutResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
						},
						"vendor": schema.StringAttribute{
							Description: "Vendor",
							Optional:    true,
						},
						"part_number": schema.StringAttribute{
							Description: "Part Number",
							Optional:    true,
						},
						"breakout": schema.StringAttribute{
							Description: "Breakout definition; defines number of ports of what speed this port is brokenout to.",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"object_properties": schema.SingleNestedBlock{
				Description: "Object properties for the SFP Breakout",
				Attributes:  map[string]schema.Attribute{
					// Empty object properties according to schema
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("sfp_breakouts") {
		tflog.Info(ctx, fmt.Sprintf("Skipping SFP Breakout %s verification – trusting recent successful API operation", sfpBreakoutName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching SFP Breakouts for verification of %s", sfpBreakoutName))

	type SfpBreakoutsResponse struct {
		SfpBreakouts map[string]interface{} `json:"sfp_breakouts"`
	}

	var result SfpBreakoutsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		sfpBreakoutsData, fetchErr := getCachedResponse(ctx, r.provCtx, "sfp_breakouts", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch SFP Breakouts")
			respAPI, err := r.client.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading SFP Breakouts: %v", err)
			}
			defer respAPI.Body.Close()

			var res SfpBreakoutsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode SFP Breakouts response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d SFP Breakouts", len(res.SfpBreakouts)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch SFP Breakouts on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = sfpBreakoutsData.(SfpBreakoutsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SFP Breakout %s", sfpBreakoutName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for SFP Breakout with ID: %s", sfpBreakoutName))
	var sfpBreakoutData map[string]interface{}
	exists := false

	if data, ok := result.SfpBreakouts[sfpBreakoutName].(map[string]interface{}); ok {
		sfpBreakoutData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found SFP Breakout directly by ID: %s", sfpBreakoutName))
	} else {
		for apiName, g := range result.SfpBreakouts {
			sfpBreakout, ok := g.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := sfpBreakout["name"].(string); ok && name == sfpBreakoutName {
				sfpBreakoutData = sfpBreakout
				sfpBreakoutName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found SFP Breakout with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("SFP Breakout with ID '%s' not found in API response", sfpBreakoutName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", sfpBreakoutData["name"]))

	if enable, ok := sfpBreakoutData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	// Handle breakout list
	if breakoutData, ok := sfpBreakoutData["breakout"].([]interface{}); ok && len(breakoutData) > 0 {
		var breakouts []veritySfpBreakoutBreakoutModel

		for _, b := range breakoutData {
			breakout, ok := b.(map[string]interface{})
			if !ok {
				continue
			}

			breakoutModel := veritySfpBreakoutBreakoutModel{}

			if enable, ok := breakout["enable"].(bool); ok {
				breakoutModel.Enable = types.BoolValue(enable)
			} else {
				breakoutModel.Enable = types.BoolNull()
			}

			if vendor, ok := breakout["vendor"].(string); ok {
				breakoutModel.Vendor = types.StringValue(vendor)
			} else {
				breakoutModel.Vendor = types.StringNull()
			}

			if partNumber, ok := breakout["part_number"].(string); ok {
				breakoutModel.PartNumber = types.StringValue(partNumber)
			} else {
				breakoutModel.PartNumber = types.StringNull()
			}

			if breakoutVal, ok := breakout["breakout"].(string); ok {
				breakoutModel.Breakout = types.StringValue(breakoutVal)
			} else {
				breakoutModel.Breakout = types.StringNull()
			}

			if index, exists := breakout["index"]; exists && index != nil {
				switch v := index.(type) {
				case int:
					breakoutModel.Index = types.Int64Value(int64(v))
				case float64:
					breakoutModel.Index = types.Int64Value(int64(v))
				default:
					breakoutModel.Index = types.Int64Null()
				}
			} else {
				breakoutModel.Index = types.Int64Null()
			}

			breakouts = append(breakouts, breakoutModel)
		}

		state.Breakout = breakouts
	} else {
		state.Breakout = nil
	}

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

	if !plan.Name.Equal(state.Name) {
		sfpBreakoutProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		sfpBreakoutProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	oldBreakoutsByIndex := make(map[int64]veritySfpBreakoutBreakoutModel)
	for _, breakout := range state.Breakout {
		if !breakout.Index.IsNull() {
			oldBreakoutsByIndex[breakout.Index.ValueInt64()] = breakout
		}
	}

	var changedBreakouts []openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner
	breakoutsChanged := false

	for _, breakout := range plan.Breakout {
		if breakout.Index.IsNull() {
			continue
		}

		index := breakout.Index.ValueInt64()
		oldBreakout, exists := oldBreakoutsByIndex[index]

		if !exists {
			// new breakout, include all fields
			breakoutProps := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !breakout.Enable.IsNull() {
				breakoutProps.Enable = openapi.PtrBool(breakout.Enable.ValueBool())
			} else {
				breakoutProps.Enable = openapi.PtrBool(false)
			}

			if !breakout.Vendor.IsNull() {
				breakoutProps.Vendor = openapi.PtrString(breakout.Vendor.ValueString())
			} else {
				breakoutProps.Vendor = openapi.PtrString("")
			}

			if !breakout.PartNumber.IsNull() {
				breakoutProps.PartNumber = openapi.PtrString(breakout.PartNumber.ValueString())
			} else {
				breakoutProps.PartNumber = openapi.PtrString("")
			}

			if !breakout.Breakout.IsNull() {
				breakoutProps.Breakout = openapi.PtrString(breakout.Breakout.ValueString())
			} else {
				breakoutProps.Breakout = openapi.PtrString("")
			}

			changedBreakouts = append(changedBreakouts, breakoutProps)
			breakoutsChanged = true
			continue
		}

		// existing breakout, check which fields changed
		breakoutProps := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !breakout.Enable.Equal(oldBreakout.Enable) {
			breakoutProps.Enable = openapi.PtrBool(breakout.Enable.ValueBool())
			fieldChanged = true
		}

		if !breakout.Vendor.Equal(oldBreakout.Vendor) {
			if !breakout.Vendor.IsNull() {
				breakoutProps.Vendor = openapi.PtrString(breakout.Vendor.ValueString())
			} else {
				breakoutProps.Vendor = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !breakout.PartNumber.Equal(oldBreakout.PartNumber) {
			if !breakout.PartNumber.IsNull() {
				breakoutProps.PartNumber = openapi.PtrString(breakout.PartNumber.ValueString())
			} else {
				breakoutProps.PartNumber = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !breakout.Breakout.Equal(oldBreakout.Breakout) {
			if !breakout.Breakout.IsNull() {
				breakoutProps.Breakout = openapi.PtrString(breakout.Breakout.ValueString())
			} else {
				breakoutProps.Breakout = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedBreakouts = append(changedBreakouts, breakoutProps)
			breakoutsChanged = true
		}
	}

	for idx := range oldBreakoutsByIndex {
		found := false
		for _, breakout := range plan.Breakout {
			if !breakout.Index.IsNull() && breakout.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// breakout removed - include only the index for deletion
			deletedBreakout := openapi.SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedBreakouts = append(changedBreakouts, deletedBreakout)
			breakoutsChanged = true
		}
	}

	if breakoutsChanged && len(changedBreakouts) > 0 {
		sfpBreakoutProps.Breakout = changedBreakouts
		hasChanges = true
	}

	if (plan.ObjectProperties == nil) != (state.ObjectProperties == nil) {
		if plan.ObjectProperties != nil {
			// SFP Breakout object properties are empty according to schema
			sfpBreakoutObjProps := map[string]interface{}{}
			sfpBreakoutProps.ObjectProperties = sfpBreakoutObjProps
		}
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "sfp_breakout", name, sfpBreakoutProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for SFP Breakout update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update SFP Breakout %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("SFP Breakout %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sfp_breakouts")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

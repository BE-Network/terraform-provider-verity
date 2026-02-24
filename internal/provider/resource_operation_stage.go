package provider

import (
	"context"
	"fmt"
	"terraform-provider-verity/internal/bulkops"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &operationStageResource{}

func NewVerityOperationStageResource() resource.Resource {
	return &operationStageResource{}
}

type operationStageResource struct {
	bulkOpsMgr *bulkops.Manager
}

type operationStageResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (r *operationStageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operation_stage"
}

func (r *operationStageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Represents a stage in the operation sequence, used to establish dependencies between resource types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for this stage.",
			},
		},
	}
}

func (r *operationStageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.bulkOpsMgr = provCtx.bulkOpsMgr
}

func (r *operationStageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan operationStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.bulkOpsMgr != nil {
		tflog.Debug(ctx, "Operation stage barrier: waiting for and flushing all pending bulk operations")
		diags := r.bulkOpsMgr.WaitAndFlushAllOperations(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.Id = types.StringValue("stage")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *operationStageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state operationStageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *operationStageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan operationStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *operationStageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.bulkOpsMgr != nil {
		tflog.Debug(ctx, "Operation stage barrier (destroy): waiting for and flushing all pending bulk operations")
		diags := r.bulkOpsMgr.WaitAndFlushAllOperations(ctx)
		resp.Diagnostics.Append(diags...)
	}
}

func (r *operationStageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

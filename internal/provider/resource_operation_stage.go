package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &operationStageResource{}

func NewVerityOperationStageResource() resource.Resource {
	return &operationStageResource{}
}

type operationStageResource struct {
	// No client needed for this resource type
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

func (r *operationStageResource) Configure(_ context.Context, _ resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// No configuration needed
}

func (r *operationStageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan operationStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Just set a placeholder ID
	plan.Id = types.StringValue("stage")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *operationStageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state operationStageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No actual reading needed, this is just a placeholder resource
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *operationStageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan operationStageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No actual updating needed
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *operationStageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No actual deletion needed
}

func (r *operationStageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Simply set the ID of the resource to the import ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

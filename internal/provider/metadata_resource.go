package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataResourceName = "metadata"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataResource{}
	_ resource.ResourceWithImportState = &MetadataResource{}
)

var metadataFields = helpers.Fields{
	Bools: []string{"artistMetadata", "albumMetadata", "artistImages", "albumImages", "trackMetadata"},
}

func NewMetadataResource() resource.Resource {
	return &MetadataResource{}
}

// MetadataResource defines the metadata implementation.
type MetadataResource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// Metadata describes the metadata data model.
type Metadata struct {
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	ConfigContract types.String `tfsdk:"config_contract"`
	Implementation types.String `tfsdk:"implementation"`
	ID             types.Int64  `tfsdk:"id"`
	Enable         types.Bool   `tfsdk:"enable"`
	ArtistMetadata types.Bool   `tfsdk:"artist_metadata"`
	AlbumMetadata  types.Bool   `tfsdk:"album_metadata"`
	ArtistImages   types.Bool   `tfsdk:"artist_images"`
	AlbumImages    types.Bool   `tfsdk:"album_images"`
	TrackMetadata  types.Bool   `tfsdk:"track_metadata"`
}

func (m Metadata) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":            types.SetType{}.WithElementType(types.Int64Type),
			"name":            types.StringType,
			"config_contract": types.StringType,
			"implementation":  types.StringType,
			"id":              types.Int64Type,
			"enable":          types.BoolType,
			"artist_metadata": types.BoolType,
			"album_metadata":  types.BoolType,
			"artist_images":   types.BoolType,
			"album_images":    types.BoolType,
			"track_metadata":  types.BoolType,
		})
}

func (r *MetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataResourceName
}

func (r *MetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->\nGeneric Metadata resource. When possible use a specific resource instead.\nFor more information refer to [Metadata](https://wiki.servarr.com/lidarr/settings#metadata) documentation.",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Metadata configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Metadata implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"track_metadata": schema.BoolAttribute{
				MarkdownDescription: "Track metadata flag.",
				Optional:            true,
				Computed:            true,
			},
			"album_images": schema.BoolAttribute{
				MarkdownDescription: "Album images flag.",
				Optional:            true,
				Computed:            true,
			},
			"artist_images": schema.BoolAttribute{
				MarkdownDescription: "Artist images flag.",
				Optional:            true,
				Computed:            true,
			},
			"artist_metadata": schema.BoolAttribute{
				MarkdownDescription: "Artist metadata flag.",
				Optional:            true,
				Computed:            true,
			},
			"album_metadata": schema.BoolAttribute{
				MarkdownDescription: "Album metadata flag.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MetadataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *MetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *Metadata

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Metadata
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.CreateMetadata(r.auth).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Metadata

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *Metadata

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Metadata current value
	response, _, err := r.client.MetadataAPI.GetMetadataById(r.auth, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Metadata

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *Metadata

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Metadata
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.UpdateMetadata(r.auth, request.GetId()).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Metadata

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *MetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Metadata current value
	_, err := r.client.MetadataAPI.DeleteMetadata(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataResourceName+": "+req.ID)
}

func (m *Metadata) write(ctx context.Context, metadata *lidarr.MetadataResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	m.Enable = types.BoolValue(metadata.GetEnable())
	m.ID = types.Int64Value(int64(metadata.GetId()))
	m.ConfigContract = types.StringValue(metadata.GetConfigContract())
	m.Implementation = types.StringValue(metadata.GetImplementation())
	m.Name = types.StringValue(metadata.GetName())
	m.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, metadata.Tags)
	diags.Append(localDiag...)
	helpers.WriteFields(ctx, m, metadata.GetFields(), metadataFields)
}

func (m *Metadata) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.MetadataResource {
	metadata := lidarr.NewMetadataResource()
	metadata.SetEnable(m.Enable.ValueBool())
	metadata.SetId(int32(m.ID.ValueInt64()))
	metadata.SetConfigContract(m.ConfigContract.ValueString())
	metadata.SetImplementation(m.Implementation.ValueString())
	metadata.SetName(m.Name.ValueString())
	diags.Append(m.Tags.ElementsAs(ctx, &metadata.Tags, true)...)
	metadata.SetFields(helpers.ReadFields(ctx, m, metadataFields))

	return metadata
}

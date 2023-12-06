package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	metadataRoksboxResourceName   = "metadata_roksbox"
	metadataRoksboxImplementation = "RoksboxMetadata"
	metadataRoksboxConfigContract = "RoksboxMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataRoksboxResource{}
	_ resource.ResourceWithImportState = &MetadataRoksboxResource{}
)

func NewMetadataRoksboxResource() resource.Resource {
	return &MetadataRoksboxResource{}
}

// MetadataRoksboxResource defines the Roksbox metadata implementation.
type MetadataRoksboxResource struct {
	client *lidarr.APIClient
}

// MetadataRoksbox describes the Roksbox metadata data model.
type MetadataRoksbox struct {
	Tags          types.Set    `tfsdk:"tags"`
	Name          types.String `tfsdk:"name"`
	ID            types.Int64  `tfsdk:"id"`
	Enable        types.Bool   `tfsdk:"enable"`
	ArtistImages  types.Bool   `tfsdk:"artist_images"`
	AlbumImages   types.Bool   `tfsdk:"album_images"`
	TrackMetadata types.Bool   `tfsdk:"track_metadata"`
}

func (m MetadataRoksbox) toMetadata() *Metadata {
	return &Metadata{
		Tags:           m.Tags,
		Name:           m.Name,
		ID:             m.ID,
		Enable:         m.Enable,
		AlbumImages:    m.AlbumImages,
		ArtistImages:   m.ArtistImages,
		TrackMetadata:  m.TrackMetadata,
		ConfigContract: types.StringValue(metadataRoksboxConfigContract),
		Implementation: types.StringValue(metadataRoksboxImplementation),
	}
}

func (m *MetadataRoksbox) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.AlbumImages = metadata.AlbumImages
	m.ArtistImages = metadata.ArtistImages
	m.TrackMetadata = metadata.TrackMetadata
}

func (r *MetadataRoksboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataRoksboxResourceName
}

func (r *MetadataRoksboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Roksbox resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/lidarr/settings#metadata) and [ROKSBOX](https://wiki.servarr.com/lidarr/supported#roksboxmetadata).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
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
			"artist_images": schema.BoolAttribute{
				MarkdownDescription: "Artist images flag.",
				Required:            true,
			},
			"album_images": schema.BoolAttribute{
				MarkdownDescription: "Album images flag.",
				Required:            true,
			},
			"track_metadata": schema.BoolAttribute{
				MarkdownDescription: "Track metadata flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataRoksboxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataRoksboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataRoksbox current value
	response, _, err := r.client.MetadataAPI.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataRoksbox

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataRoksbox
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataRoksboxResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataRoksboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataRoksbox current value
	_, err := r.client.MetadataAPI.DeleteMetadata(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataRoksboxResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataRoksboxResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataRoksboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataRoksboxResourceName+": "+req.ID)
}

func (m *MetadataRoksbox) write(ctx context.Context, metadata *lidarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataRoksbox) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}

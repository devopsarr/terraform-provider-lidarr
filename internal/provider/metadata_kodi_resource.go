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
	metadataKodiResourceName   = "metadata_kodi"
	metadataKodiImplementation = "XbmcMetadata"
	metadataKodiConfigContract = "XbmcMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataKodiResource{}
	_ resource.ResourceWithImportState = &MetadataKodiResource{}
)

func NewMetadataKodiResource() resource.Resource {
	return &MetadataKodiResource{}
}

// MetadataKodiResource defines the Kodi metadata implementation.
type MetadataKodiResource struct {
	client *lidarr.APIClient
}

// MetadataKodi describes the Kodi metadata data model.
type MetadataKodi struct {
	Tags           types.Set    `tfsdk:"tags"`
	Name           types.String `tfsdk:"name"`
	ID             types.Int64  `tfsdk:"id"`
	Enable         types.Bool   `tfsdk:"enable"`
	ArtistMetadata types.Bool   `tfsdk:"artist_metadata"`
	AlbumMetadata  types.Bool   `tfsdk:"album_metadata"`
	ArtistImages   types.Bool   `tfsdk:"artist_images"`
	AlbumImages    types.Bool   `tfsdk:"album_images"`
}

func (m MetadataKodi) toMetadata() *Metadata {
	return &Metadata{
		Tags:           m.Tags,
		Name:           m.Name,
		ID:             m.ID,
		Enable:         m.Enable,
		ArtistMetadata: m.ArtistMetadata,
		AlbumMetadata:  m.AlbumMetadata,
		ArtistImages:   m.ArtistImages,
		AlbumImages:    m.AlbumImages,
		ConfigContract: types.StringValue(metadataKodiConfigContract),
		Implementation: types.StringValue(metadataKodiImplementation),
	}
}

func (m *MetadataKodi) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.ArtistMetadata = metadata.ArtistMetadata
	m.ArtistImages = metadata.ArtistImages
	m.AlbumMetadata = metadata.AlbumMetadata
	m.AlbumImages = metadata.AlbumImages
}

func (r *MetadataKodiResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataKodiResourceName
}

func (r *MetadataKodiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Kodi resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/lidarr/settings#metadata) and [KODI](https://wiki.servarr.com/lidarr/supported#xbmcmetadata).",
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
			"album_images": schema.BoolAttribute{
				MarkdownDescription: "Album images flag.",
				Required:            true,
			},
			"artist_images": schema.BoolAttribute{
				MarkdownDescription: "Artist images flag.",
				Required:            true,
			},
			"artist_metadata": schema.BoolAttribute{
				MarkdownDescription: "Artist metadata flag.",
				Required:            true,
			},
			"album_metadata": schema.BoolAttribute{
				MarkdownDescription: "Album metadata flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataKodiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataKodiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataKodi
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataKodi current value
	response, _, err := r.client.MetadataAPI.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataKodi
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataAPI.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataKodi current value
	_, err := r.client.MetadataAPI.DeleteMetadata(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataKodiResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataKodiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataKodiResourceName+": "+req.ID)
}

func (m *MetadataKodi) write(ctx context.Context, metadata *lidarr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataKodi) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}

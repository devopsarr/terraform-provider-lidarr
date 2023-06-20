package provider

import (
	"context"
	"fmt"
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

const artistResourceName = "artist"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ArtistResource{}
	_ resource.ResourceWithImportState = &ArtistResource{}
)

func NewArtistResource() resource.Resource {
	return &ArtistResource{}
}

// ArtistResource defines the artist implementation.
type ArtistResource struct {
	client *lidarr.APIClient
}

// Artist describes the artist data model.
type Artist struct {
	Genres            types.Set    `tfsdk:"genres"`
	Tags              types.Set    `tfsdk:"tags"`
	ArtistName        types.String `tfsdk:"artist_name"`
	ForeignArtistID   types.String `tfsdk:"foreign_artist_id"`
	Status            types.String `tfsdk:"status"`
	Path              types.String `tfsdk:"path"`
	Overview          types.String `tfsdk:"overview"`
	ID                types.Int64  `tfsdk:"id"`
	QualityProfileID  types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID types.Int64  `tfsdk:"metadata_profile_id"`
	Monitored         types.Bool   `tfsdk:"monitored"`

	// TODO: future Implementation
	// Links          types.Set    `tfsdk:"links"`
	// SortName       types.String `tfsdk:"sortName"`
	// Ended          types.Bool   `tfsdk:"ended"`
	// ArtistType     types.String `tfsdk:"artist_type"`
	// Disambiguation types.String `tfsdk:"disambiguation"`
	// RootFolderPath types.String `tfsdk:"root_folder_path"`
	// FolderName     types.String `tfsdk:"folderName"`
	// CleanName      types.String `tfsdk:"cleanName"`
	// Certification  types.String `tfsdk:"certification"`
	// Added          types.String `tfsdk:"added"`
	// Ratings        types.Object `tfsdk:"ratings"`
	// TadbId         types.Int64  `tfsdk:"tadb_id"`
	// DiscogsId      types.Int64  `tfsdk:"discogs_id"`
}

func (a Artist) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"monitored":           types.BoolType,
			"id":                  types.Int64Type,
			"quality_profile_id":  types.Int64Type,
			"metadata_profile_id": types.Int64Type,
			"artist_name":         types.StringType,
			"foreign_artist_id":   types.StringType,
			"status":              types.StringType,
			"path":                types.StringType,
			"overview":            types.StringType,
			"genres":              types.SetType{}.WithElementType(types.StringType),
			"tags":                types.SetType{}.WithElementType(types.Int64Type),
		})
}

func (r *ArtistResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + artistResourceName
}

func (r *ArtistResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Artists -->Artist resource.\nFor more information refer to [Artists](https://wiki.servarr.com/lidarr/library#artists) documentation.",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Required:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Artist ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"artist_name": schema.StringAttribute{
				MarkdownDescription: "Artist name.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full artist path.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Artist status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"foreign_artist_id": schema.StringAttribute{
				MarkdownDescription: "Foreign artist ID.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ArtistResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ArtistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var artist *Artist

	resp.Diagnostics.Append(req.Plan.Get(ctx, &artist)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Artist
	request := artist.read(ctx, &resp.Diagnostics)
	// TODO: can parametrize AddArtistOptions
	options := lidarr.NewAddArtistOptions()
	options.SetMonitor(lidarr.MONITORTYPES_ALL)
	options.SetSearchForMissingAlbums(true)

	response, _, err := r.client.ArtistApi.CreateArtist(ctx).ArtistResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, artistResourceName, err))

		return
	}

	tflog.Trace(ctx, "created artist: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	artist.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &artist)...)
}

func (r *ArtistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var artist *Artist

	resp.Diagnostics.Append(req.State.Get(ctx, &artist)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get artist current value
	response, _, err := r.client.ArtistApi.GetArtistById(ctx, int32(artist.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, artistResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+artistResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	artist.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &artist)...)
}

func (r *ArtistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var artist *Artist

	resp.Diagnostics.Append(req.Plan.Get(ctx, &artist)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Artist
	request := artist.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ArtistApi.UpdateArtist(ctx, fmt.Sprint(request.GetId())).ArtistResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, artistResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+artistResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	artist.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &artist)...)
}

func (r *ArtistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete artist current value
	_, err := r.client.ArtistApi.DeleteArtist(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, artistResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+artistResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ArtistResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+artistResourceName+": "+req.ID)
}

func (a *Artist) write(ctx context.Context, artist *lidarr.ArtistResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	a.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, artist.GetTags())
	diags.Append(localDiag...)
	a.Genres, localDiag = types.SetValueFrom(ctx, types.StringType, artist.GetGenres())
	diags.Append(localDiag...)

	a.Monitored = types.BoolValue(artist.GetMonitored())
	a.ID = types.Int64Value(int64(artist.GetId()))
	a.ArtistName = types.StringValue(artist.GetArtistName())
	a.Path = types.StringValue(artist.GetPath())
	a.QualityProfileID = types.Int64Value(int64(artist.GetQualityProfileId()))
	a.MetadataProfileID = types.Int64Value(int64(artist.GetMetadataProfileId()))
	a.ForeignArtistID = types.StringValue(artist.GetForeignArtistId())
	// Read only values
	a.Status = types.StringValue(string(artist.GetStatus()))
	a.Overview = types.StringValue(artist.GetOverview())
}

func (a *Artist) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.ArtistResource {
	artist := lidarr.NewArtistResource()
	artist.SetMonitored(a.Monitored.ValueBool())
	artist.SetArtistName(a.ArtistName.ValueString())
	artist.SetPath(a.Path.ValueString())
	artist.SetQualityProfileId(int32(a.QualityProfileID.ValueInt64()))
	artist.SetMetadataProfileId(int32(a.MetadataProfileID.ValueInt64()))
	artist.SetForeignArtistId(a.ForeignArtistID.ValueString())
	artist.SetId(int32(a.ID.ValueInt64()))
	diags.Append(a.Tags.ElementsAs(ctx, &artist.Tags, true)...)

	return artist
}

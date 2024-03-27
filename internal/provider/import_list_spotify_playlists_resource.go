package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListSpotifyPlaylistsResourceName   = "import_list_spotify_playlists"
	importListSpotifyPlaylistsImplementation = "SpotifyPlaylist"
	importListSpotifyPlaylistsConfigContract = "SpotifyPlaylistSettings"
	importListSpotifyPlaylistsType           = "spotify"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListSpotifyPlaylistsResource{}
	_ resource.ResourceWithImportState = &ImportListSpotifyPlaylistsResource{}
)

func NewImportListSpotifyPlaylistsResource() resource.Resource {
	return &ImportListSpotifyPlaylistsResource{}
}

// ImportListSpotifyPlaylistsResource defines the import list implementation.
type ImportListSpotifyPlaylistsResource struct {
	client *lidarr.APIClient
}

// ImportListSpotifyPlaylists describes the import list data model.
type ImportListSpotifyPlaylists struct {
	Tags                  types.Set    `tfsdk:"tags"`
	PlaylistIDs           types.Set    `tfsdk:"playlist_ids"`
	Name                  types.String `tfsdk:"name"`
	AccessToken           types.String `tfsdk:"access_token"`
	RefreshToken          types.String `tfsdk:"refresh_token"`
	Expires               types.String `tfsdk:"expires"`
	MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	ShouldMonitor         types.String `tfsdk:"should_monitor"`
	RootFolderPath        types.String `tfsdk:"root_folder_path"`
	QualityProfileID      types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID     types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder             types.Int64  `tfsdk:"list_order"`
	ID                    types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd    types.Bool   `tfsdk:"enable_automatic_add"`
	ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch          types.Bool   `tfsdk:"should_search"`
}

func (i ImportListSpotifyPlaylists) toImportList() *ImportList {
	return &ImportList{
		Tags:                  i.Tags,
		PlaylistIDs:           i.PlaylistIDs,
		Name:                  i.Name,
		MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:         i.ShouldMonitor,
		RootFolderPath:        i.RootFolderPath,
		AccessToken:           i.AccessToken,
		RefreshToken:          i.RefreshToken,
		Expires:               i.Expires,
		QualityProfileID:      i.QualityProfileID,
		MetadataProfileID:     i.MetadataProfileID,
		ListOrder:             i.ListOrder,
		ID:                    i.ID,
		EnableAutomaticAdd:    i.EnableAutomaticAdd,
		ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:          i.ShouldSearch,
		Implementation:        types.StringValue(importListSpotifyPlaylistsImplementation),
		ConfigContract:        types.StringValue(importListSpotifyPlaylistsConfigContract),
		ListType:              types.StringValue(importListSpotifyPlaylistsType),
	}
}

func (i *ImportListSpotifyPlaylists) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.PlaylistIDs = importList.PlaylistIDs
	i.Name = importList.Name
	i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.AccessToken = importList.AccessToken
	i.RefreshToken = importList.RefreshToken
	i.Expires = importList.Expires
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListSpotifyPlaylistsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListSpotifyPlaylistsResourceName
}

func (r *ImportListSpotifyPlaylistsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->\nImport List Spotify Playlist resource.\nFor more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists) and [Spotify Playlists](https://wiki.servarr.com/lidarr/supported#spotifyplaylist).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			"should_monitor_existing": schema.BoolAttribute{
				MarkdownDescription: "Should monitor existing flag.",
				Optional:            true,
				Computed:            true,
			},
			"should_search": schema.BoolAttribute{
				MarkdownDescription: "Should search flag.",
				Optional:            true,
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Optional:            true,
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Optional:            true,
				Computed:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "specificAlbum", "entireArtist"),
				},
			},
			"monitor_new_items": schema.StringAttribute{
				MarkdownDescription: "Monitor new items.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "all", "new"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Required:            true,
				Sensitive:           true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Required:            true,
			},
			"playlist_ids": schema.SetAttribute{
				MarkdownDescription: "Playlist IDs.",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ImportListSpotifyPlaylistsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListSpotifyPlaylistsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListSpotifyPlaylists

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListSpotifyPlaylists
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSpotifyPlaylistsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListSpotifyPlaylists

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListSpotifyPlaylists current value
	response, _, err := r.client.ImportListAPI.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSpotifyPlaylistsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListSpotifyPlaylists

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListSpotifyPlaylists
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSpotifyPlaylistsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListSpotifyPlaylists current value
	_, err := r.client.ImportListAPI.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListSpotifyPlaylistsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListSpotifyPlaylistsResourceName+": "+req.ID)
}

func (i *ImportListSpotifyPlaylists) write(ctx context.Context, importList *lidarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListSpotifyPlaylists) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}

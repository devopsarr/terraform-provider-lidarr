package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	ImportListSpotifyPlaylistsResourceName   = "import_list_spotify_playlists"
	ImportListSpotifyPlaylistsImplementation = "SpotifyPlaylist"
	ImportListSpotifyPlaylistsConfigContract = "SpotifyPlaylistSettings"
	ImportListSpotifyPlaylistsType           = "spotify"
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
	PlaylistIds           types.Set    `tfsdk:"playlist_ids"`
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
		PlaylistIds:           i.PlaylistIds,
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
	}
}

func (i *ImportListSpotifyPlaylists) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.PlaylistIds = importList.PlaylistIds
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

func (r *ImportListSpotifyPlaylistsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ImportListSpotifyPlaylistsResourceName
}

func (r *ImportListSpotifyPlaylistsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Spotify Playlist resource.\nFor more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists) and [Spotify Playlists](https://wiki.servarr.com/lidarr/supported#spotifyplaylist).",
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
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, ImportListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+ImportListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
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
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, ImportListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+ImportListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
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
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, ImportListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+ImportListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListSpotifyPlaylistsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListSpotifyPlaylists

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListSpotifyPlaylists current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, ImportListSpotifyPlaylistsResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+ImportListSpotifyPlaylistsResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListSpotifyPlaylistsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+ImportListSpotifyPlaylistsResourceName+": "+req.ID)
}

func (i *ImportListSpotifyPlaylists) write(ctx context.Context, importList *lidarr.ImportListResource) {
	genericImportList := ImportList{
		Name:                  types.StringValue(importList.GetName()),
		ShouldMonitor:         types.StringValue(string(importList.GetShouldMonitor())),
		MonitorNewItems:       types.StringValue(string(importList.GetMonitorNewItems())),
		RootFolderPath:        types.StringValue(importList.GetRootFolderPath()),
		QualityProfileID:      types.Int64Value(int64(importList.GetQualityProfileId())),
		MetadataProfileID:     types.Int64Value(int64(importList.GetMetadataProfileId())),
		ID:                    types.Int64Value(int64(importList.GetId())),
		ListOrder:             types.Int64Value(int64(importList.GetListOrder())),
		EnableAutomaticAdd:    types.BoolValue(importList.GetEnableAutomaticAdd()),
		ShouldMonitorExisting: types.BoolValue(importList.GetShouldMonitorExisting()),
		ShouldSearch:          types.BoolValue(importList.GetShouldSearch()),
	}
	genericImportList.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	genericImportList.writeFields(ctx, importList.GetFields())
	i.fromImportList(&genericImportList)
}

func (i *ImportListSpotifyPlaylists) read(ctx context.Context) *lidarr.ImportListResource {
	tags := make([]*int32, len(i.Tags.Elements()))
	tfsdk.ValueAs(ctx, i.Tags, &tags)

	list := lidarr.NewImportListResource()
	list.SetShouldMonitor(lidarr.ImportListMonitorType(i.ShouldMonitor.ValueString()))
	list.SetMonitorNewItems(lidarr.NewItemMonitorTypes(i.MonitorNewItems.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetMetadataProfileId(int32(i.MetadataProfileID.ValueInt64()))
	list.SetListOrder(int32(i.ListOrder.ValueInt64()))
	list.SetEnableAutomaticAdd(i.EnableAutomaticAdd.ValueBool())
	list.SetShouldMonitorExisting(i.ShouldMonitorExisting.ValueBool())
	list.SetShouldSearch(i.ShouldSearch.ValueBool())
	list.SetListType(ImportListSpotifyPlaylistsType)
	list.SetConfigContract(ImportListSpotifyPlaylistsConfigContract)
	list.SetImplementation(ImportListSpotifyPlaylistsImplementation)
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.toImportList().readFields(ctx))

	return list
}

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
	importListMusicBrainzResourceName   = "import_list_music_brainz"
	importListMusicBrainzImplementation = "MusicBrainzSeries"
	importListMusicBrainzConfigContract = "MusicBrainzSeriesSettings"
	importListMusicBrainzType           = "other"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListMusicBrainzResource{}
	_ resource.ResourceWithImportState = &ImportListMusicBrainzResource{}
)

func NewImportListMusicBrainzResource() resource.Resource {
	return &ImportListMusicBrainzResource{}
}

// ImportListMusicBrainzResource defines the import list implementation.
type ImportListMusicBrainzResource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// ImportListMusicBrainz describes the import list data model.
type ImportListMusicBrainz struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	ShouldMonitor         types.String `tfsdk:"should_monitor"`
	RootFolderPath        types.String `tfsdk:"root_folder_path"`
	SeriesID              types.String `tfsdk:"series_id"`
	QualityProfileID      types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID     types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder             types.Int64  `tfsdk:"list_order"`
	ID                    types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd    types.Bool   `tfsdk:"enable_automatic_add"`
	ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch          types.Bool   `tfsdk:"should_search"`
}

func (i ImportListMusicBrainz) toImportList() *ImportList {
	return &ImportList{
		Tags:                  i.Tags,
		Name:                  i.Name,
		MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:         i.ShouldMonitor,
		RootFolderPath:        i.RootFolderPath,
		SeriesID:              i.SeriesID,
		QualityProfileID:      i.QualityProfileID,
		MetadataProfileID:     i.MetadataProfileID,
		ListOrder:             i.ListOrder,
		ID:                    i.ID,
		EnableAutomaticAdd:    i.EnableAutomaticAdd,
		ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:          i.ShouldSearch,
		Implementation:        types.StringValue(importListMusicBrainzImplementation),
		ConfigContract:        types.StringValue(importListMusicBrainzConfigContract),
		ListType:              types.StringValue(importListMusicBrainzType),
	}
}

func (i *ImportListMusicBrainz) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.SeriesID = importList.SeriesID
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListMusicBrainzResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListMusicBrainzResourceName
}

func (r *ImportListMusicBrainzResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->\nImport List MusicBrainz resource.\nFor more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists) and [MusicBrainz](https://wiki.servarr.com/lidarr/supported#musicbrainzseries).",
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
			"series_id": schema.StringAttribute{
				MarkdownDescription: "Series ID.",
				Required:            true,
			},
		},
	}
}

func (r *ImportListMusicBrainzResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *ImportListMusicBrainzResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListMusicBrainz

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListMusicBrainz
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.CreateImportList(r.auth).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListMusicBrainzResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListMusicBrainzResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListMusicBrainzResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListMusicBrainz

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListMusicBrainz current value
	response, _, err := r.client.ImportListAPI.GetImportListById(r.auth, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListMusicBrainzResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListMusicBrainzResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListMusicBrainzResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListMusicBrainz

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListMusicBrainz
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.UpdateImportList(r.auth, request.GetId()).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListMusicBrainzResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListMusicBrainzResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListMusicBrainzResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListMusicBrainz current value
	_, err := r.client.ImportListAPI.DeleteImportList(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListMusicBrainzResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListMusicBrainzResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListMusicBrainzResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListMusicBrainzResourceName+": "+req.ID)
}

func (i *ImportListMusicBrainz) write(ctx context.Context, importList *lidarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListMusicBrainz) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}

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
	importListLastFMTagResourceName   = "import_list_lastfm_tag"
	importListLastFMTagImplementation = "LastFmTag"
	importListLastFMTagConfigContract = "LastFmTagSettings"
	importListLastFMTagType           = "lastFm"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListLastFMTagResource{}
	_ resource.ResourceWithImportState = &ImportListLastFMTagResource{}
)

func NewImportListLastFMTagResource() resource.Resource {
	return &ImportListLastFMTagResource{}
}

// ImportListLastFMTagResource defines the import list implementation.
type ImportListLastFMTagResource struct {
	client *lidarr.APIClient
}

// ImportListLastFMTag describes the import list data model.
type ImportListLastFMTag struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	ShouldMonitor         types.String `tfsdk:"should_monitor"`
	RootFolderPath        types.String `tfsdk:"root_folder_path"`
	TagID                 types.String `tfsdk:"tag_id"`
	Count                 types.Int64  `tfsdk:"count_list"`
	QualityProfileID      types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID     types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder             types.Int64  `tfsdk:"list_order"`
	ID                    types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd    types.Bool   `tfsdk:"enable_automatic_add"`
	ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch          types.Bool   `tfsdk:"should_search"`
}

func (i ImportListLastFMTag) toImportList() *ImportList {
	return &ImportList{
		Tags:                  i.Tags,
		Name:                  i.Name,
		MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:         i.ShouldMonitor,
		RootFolderPath:        i.RootFolderPath,
		TagID:                 i.TagID,
		Count:                 i.Count,
		QualityProfileID:      i.QualityProfileID,
		MetadataProfileID:     i.MetadataProfileID,
		ListOrder:             i.ListOrder,
		ID:                    i.ID,
		EnableAutomaticAdd:    i.EnableAutomaticAdd,
		ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:          i.ShouldSearch,
		Implementation:        types.StringValue(importListLastFMTagImplementation),
		ConfigContract:        types.StringValue(importListLastFMTagConfigContract),
		ListType:              types.StringValue(importListLastFMTagType),
	}
}

func (i *ImportListLastFMTag) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.TagID = importList.TagID
	i.Count = importList.Count
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListLastFMTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListLastFMTagResourceName
}

func (r *ImportListLastFMTagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Last.fm Tag resource.\nFor more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists) and [Last.fm Tag](https://wiki.servarr.com/lidarr/supported#lastfmtag).",
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
			"count_list": schema.Int64Attribute{
				MarkdownDescription: "Elements to pull from list.",
				Required:            true,
			},
			"tag_id": schema.StringAttribute{
				MarkdownDescription: "Tag ID.",
				Required:            true,
			},
		},
	}
}

func (r *ImportListLastFMTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListLastFMTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListLastFMTag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListLastFMTag
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListLastFMTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListLastFMTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLastFMTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListLastFMTag

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListLastFMTag current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListLastFMTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListLastFMTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLastFMTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListLastFMTag

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListLastFMTag
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListLastFMTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListLastFMTagResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLastFMTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListLastFMTag current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListLastFMTagResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListLastFMTagResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListLastFMTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListLastFMTagResourceName+": "+req.ID)
}

func (i *ImportListLastFMTag) write(ctx context.Context, importList *lidarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListLastFMTag) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}

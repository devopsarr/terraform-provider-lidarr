package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

const importListResourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListResource{}
	_ resource.ResourceWithImportState = &ImportListResource{}
)

var importListFields = helpers.Fields{
	Ints:         []string{"count"},
	Strings:      []string{"baseUrl", "apiKey", "tagId", "userId", "listId", "seriesId", "accessToken", "refreshToken", "expires"},
	IntSlices:    []string{"profileIds", "tagIds"},
	StringSlices: []string{"playlistIds"},
}

func NewImportListResource() resource.Resource {
	return &ImportListResource{}
}

// ImportListResource defines the download client implementation.
type ImportListResource struct {
	client *lidarr.APIClient
}

// ImportList describes the download client data model.
type ImportList struct {
	ProfileIds            types.Set    `tfsdk:"profile_ids"`
	TagIds                types.Set    `tfsdk:"tag_ids"`
	PlaylistIds           types.Set    `tfsdk:"playlist_ids"`
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	ConfigContract        types.String `tfsdk:"config_contract"`
	Implementation        types.String `tfsdk:"implementation"`
	MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	AccessToken           types.String `tfsdk:"access_token"`
	RefreshToken          types.String `tfsdk:"refresh_token"`
	Expires               types.String `tfsdk:"expires"`
	ShouldMonitor         types.String `tfsdk:"should_monitor"`
	ListType              types.String `tfsdk:"list_type"`
	RootFolderPath        types.String `tfsdk:"root_folder_path"`
	BaseURL               types.String `tfsdk:"base_url"`
	APIKey                types.String `tfsdk:"api_key"`
	TagID                 types.String `tfsdk:"tag_id"`
	UserID                types.String `tfsdk:"user_id"`
	ListID                types.String `tfsdk:"list_id"`
	SeriesID              types.String `tfsdk:"series_id"`
	Count                 types.Int64  `tfsdk:"count_list"`
	QualityProfileID      types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID     types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder             types.Int64  `tfsdk:"list_order"`
	ID                    types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd    types.Bool   `tfsdk:"enable_automatic_add"`
	ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch          types.Bool   `tfsdk:"should_search"`
}

func (i ImportList) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"playlist_ids":            types.SetType{}.WithElementType(types.StringType),
			"profile_ids":             types.SetType{}.WithElementType(types.Int64Type),
			"tag_ids":                 types.SetType{}.WithElementType(types.Int64Type),
			"tags":                    types.SetType{}.WithElementType(types.Int64Type),
			"name":                    types.StringType,
			"config_contract":         types.StringType,
			"implementation":          types.StringType,
			"monitor_new_items":       types.StringType,
			"access_token":            types.StringType,
			"refresh_token":           types.StringType,
			"expires":                 types.StringType,
			"should_monitor":          types.StringType,
			"list_type":               types.StringType,
			"root_folder_path":        types.StringType,
			"base_url":                types.StringType,
			"api_key":                 types.StringType,
			"tag_id":                  types.StringType,
			"user_id":                 types.StringType,
			"list_id":                 types.StringType,
			"series_id":               types.StringType,
			"count_list":              types.Int64Type,
			"quality_profile_id":      types.Int64Type,
			"metadata_profile_id":     types.Int64Type,
			"list_order":              types.Int64Type,
			"id":                      types.Int64Type,
			"enable_automatic_add":    types.BoolType,
			"should_monitor_existing": types.BoolType,
			"should_search":           types.BoolType,
		})
}

func (r *ImportListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListResourceName
}

func (r *ImportListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Generic Import List resource. When possible use a specific resource instead.\nFor more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists).",
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
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Required:            true,
			},
			"list_type": schema.StringAttribute{
				MarkdownDescription: "List type.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("program", "spotify", "lastFm", "other"),
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
				Optional:            true,
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID.",
				Optional:            true,
				Computed:            true,
			},
			"tag_id": schema.StringAttribute{
				MarkdownDescription: "Tag ID.",
				Optional:            true,
				Computed:            true,
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "List ID.",
				Optional:            true,
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "Series ID.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"playlist_ids": schema.SetAttribute{
				MarkdownDescription: "Playlist IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ImportListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportList
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportList

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportList current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportList

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportList
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state ImportList

	state.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ImportListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportList current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListResourceName+": "+req.ID)
}

func (i *ImportList) write(ctx context.Context, importList *lidarr.ImportListResource, diags *diag.Diagnostics) {
	var localDiag diag.Diagnostics

	i.Tags, localDiag = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	diags.Append(localDiag...)

	i.EnableAutomaticAdd = types.BoolValue(importList.GetEnableAutomaticAdd())
	i.ShouldMonitorExisting = types.BoolValue(importList.GetShouldMonitorExisting())
	i.ShouldSearch = types.BoolValue(importList.GetShouldSearch())
	i.QualityProfileID = types.Int64Value(int64(importList.GetQualityProfileId()))
	i.MetadataProfileID = types.Int64Value(int64(importList.GetMetadataProfileId()))
	i.ID = types.Int64Value(int64(importList.GetId()))
	i.ListOrder = types.Int64Value(int64(importList.GetListOrder()))
	i.ConfigContract = types.StringValue(importList.GetConfigContract())
	i.Implementation = types.StringValue(importList.GetImplementation())
	i.ShouldMonitor = types.StringValue(string(importList.GetShouldMonitor()))
	i.MonitorNewItems = types.StringValue(string(importList.GetMonitorNewItems()))
	i.RootFolderPath = types.StringValue(importList.GetRootFolderPath())
	i.ListType = types.StringValue(string(importList.GetListType()))
	i.Name = types.StringValue(importList.GetName())
	i.ProfileIds = types.SetValueMust(types.Int64Type, nil)
	i.TagIds = types.SetValueMust(types.Int64Type, nil)
	i.PlaylistIds = types.SetValueMust(types.StringType, nil)
	helpers.WriteFields(ctx, i, importList.GetFields(), importListFields)
}

func (i *ImportList) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.ImportListResource {
	list := lidarr.NewImportListResource()
	list.SetEnableAutomaticAdd(i.EnableAutomaticAdd.ValueBool())
	list.SetShouldMonitorExisting(i.ShouldMonitorExisting.ValueBool())
	list.SetShouldSearch(i.ShouldSearch.ValueBool())
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetMetadataProfileId(int32(i.MetadataProfileID.ValueInt64()))
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetListOrder(int32(i.ListOrder.ValueInt64()))
	list.SetShouldMonitor(lidarr.ImportListMonitorType(i.ShouldMonitor.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetMonitorNewItems(lidarr.NewItemMonitorTypes(i.MonitorNewItems.ValueString()))
	list.SetListType(lidarr.ImportListType(i.ListType.ValueString()))
	list.SetConfigContract(i.ConfigContract.ValueString())
	list.SetImplementation(i.Implementation.ValueString())
	list.SetName(i.Name.ValueString())
	diags.Append(i.Tags.ElementsAs(ctx, &list.Tags, true)...)
	list.SetFields(helpers.ReadFields(ctx, i, importListFields))

	return list
}

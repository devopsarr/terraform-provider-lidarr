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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const rootFolderResourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RootFolderResource{}
	_ resource.ResourceWithImportState = &RootFolderResource{}
)

func NewRootFolderResource() resource.Resource {
	return &RootFolderResource{}
}

// RootFolderResource defines the root folder implementation.
type RootFolderResource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// RootFolder describes the root folder data model.
type RootFolder struct {
	Tags                 types.Set    `tfsdk:"tags"`
	Path                 types.String `tfsdk:"path"`
	Name                 types.String `tfsdk:"name"`
	MonitorOption        types.String `tfsdk:"monitor_option"`
	NewItemMonitorOption types.String `tfsdk:"new_item_monitor_option"`
	ID                   types.Int64  `tfsdk:"id"`
	MetadataProfileID    types.Int64  `tfsdk:"metadata_profile_id"`
	QualityProfileID     types.Int64  `tfsdk:"quality_profile_id"`
	Accessible           types.Bool   `tfsdk:"accessible"`
}

func (r RootFolder) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"tags":                    types.SetType{}.WithElementType(types.Int64Type),
			"path":                    types.StringType,
			"name":                    types.StringType,
			"monitor_option":          types.StringType,
			"new_item_monitor_option": types.StringType,
			"id":                      types.Int64Type,
			"metadata_profile_id":     types.Int64Type,
			"quality_profile_id":      types.Int64Type,
			"accessible":              types.BoolType,
		})
}

func (r *RootFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderResourceName
}

func (r *RootFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->\nRoot Folder resource.\nFor more information refer to [Root Folders](https://wiki.servarr.com/lidarr/settings#root-folders) documentation.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Root Folder friendly name.",
				Required:            true,
			},
			"monitor_option": schema.StringAttribute{
				MarkdownDescription: "Monitor option.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "future", "missing", "existing", "latest", "first", "none", "unknown"),
				},
			},
			"new_item_monitor_option": schema.StringAttribute{
				MarkdownDescription: "New item monitor option.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "none", "new"),
				},
			},
			"accessible": schema.BoolAttribute{
				MarkdownDescription: "Access flag.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Required:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *RootFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *RootFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var folder *RootFolder

	resp.Diagnostics.Append(req.Plan.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RootFolder
	request := folder.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RootFolderAPI.CreateRootFolder(r.auth).RootFolderResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var folder *RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get rootFolder current value
	response, _, err := r.client.RootFolderAPI.GetRootFolderById(r.auth, int32(folder.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var folder *RootFolder

	resp.Diagnostics.Append(req.Plan.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := folder.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RootFolderAPI.UpdateRootFolder(r.auth, strconv.Itoa(int(request.GetId()))).RootFolderResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	_, err := r.client.RootFolderAPI.DeleteRootFolder(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+rootFolderResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *RootFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+rootFolderResourceName+": "+req.ID)
}

func (r *RootFolder) write(ctx context.Context, rootFolder *lidarr.RootFolderResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	r.Accessible = types.BoolValue(rootFolder.GetAccessible())
	r.ID = types.Int64Value(int64(rootFolder.GetId()))
	r.Path = types.StringValue(rootFolder.GetPath())
	r.MetadataProfileID = types.Int64Value(int64(rootFolder.GetDefaultMetadataProfileId()))
	r.QualityProfileID = types.Int64Value(int64(rootFolder.GetDefaultQualityProfileId()))
	r.Name = types.StringValue(rootFolder.GetName())
	r.MonitorOption = types.StringValue(string(rootFolder.GetDefaultMonitorOption()))
	r.NewItemMonitorOption = types.StringValue(string(rootFolder.GetDefaultNewItemMonitorOption()))
	r.Tags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, rootFolder.GetDefaultTags())
	diags.Append(tempDiag...)
}

func (r *RootFolder) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.RootFolderResource {
	folder := lidarr.NewRootFolderResource()
	folder.SetId(int32(r.ID.ValueInt64()))
	folder.SetDefaultMetadataProfileId(int32(r.MetadataProfileID.ValueInt64()))
	folder.SetDefaultQualityProfileId(int32(r.QualityProfileID.ValueInt64()))
	folder.SetPath(r.Path.ValueString())
	folder.SetDefaultMonitorOption(lidarr.MonitorTypes(r.MonitorOption.ValueString()))
	folder.SetDefaultNewItemMonitorOption(lidarr.NewItemMonitorTypes(r.NewItemMonitorOption.ValueString()))
	folder.SetName(r.Name.ValueString())
	diags.Append(r.Tags.ElementsAs(ctx, &folder.DefaultTags, true)...)

	return folder
}

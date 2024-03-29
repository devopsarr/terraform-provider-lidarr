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
	notificationSynologyResourceName   = "notification_synology_indexer"
	notificationSynologyImplementation = "SynologyIndexer"
	notificationSynologyConfigContract = "SynologyIndexerSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSynologyResource{}
	_ resource.ResourceWithImportState = &NotificationSynologyResource{}
)

func NewNotificationSynologyResource() resource.Resource {
	return &NotificationSynologyResource{}
}

// NotificationSynologyResource defines the notification implementation.
type NotificationSynologyResource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// NotificationSynology describes the notification data model.
type NotificationSynology struct {
	Tags            types.Set    `tfsdk:"tags"`
	Name            types.String `tfsdk:"name"`
	ID              types.Int64  `tfsdk:"id"`
	UpdateLibrary   types.Bool   `tfsdk:"update_library"`
	OnReleaseImport types.Bool   `tfsdk:"on_release_import"`
	OnAlbumDelete   types.Bool   `tfsdk:"on_album_delete"`
	OnArtistDelete  types.Bool   `tfsdk:"on_artist_delete"`
	OnTrackRetag    types.Bool   `tfsdk:"on_track_retag"`
	OnRename        types.Bool   `tfsdk:"on_rename"`
	OnUpgrade       types.Bool   `tfsdk:"on_upgrade"`
}

func (n NotificationSynology) toNotification() *Notification {
	return &Notification{
		Tags:            n.Tags,
		Name:            n.Name,
		ID:              n.ID,
		UpdateLibrary:   n.UpdateLibrary,
		OnReleaseImport: n.OnReleaseImport,
		OnRename:        n.OnRename,
		OnTrackRetag:    n.OnTrackRetag,
		OnUpgrade:       n.OnUpgrade,
		Implementation:  types.StringValue(notificationSynologyImplementation),
		ConfigContract:  types.StringValue(notificationSynologyConfigContract),
	}
}

func (n *NotificationSynology) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Name = notification.Name
	n.ID = notification.ID
	n.UpdateLibrary = notification.UpdateLibrary
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnAlbumDelete = notification.OnAlbumDelete
	n.OnArtistDelete = notification.OnArtistDelete
	n.OnTrackRetag = notification.OnTrackRetag
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
}

func (r *NotificationSynologyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSynologyResourceName
}

func (r *NotificationSynologyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Synology Indexer resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Synology](https://wiki.servarr.com/lidarr/supported#synologyindexer).",
		Attributes: map[string]schema.Attribute{
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_track_retag": schema.BoolAttribute{
				MarkdownDescription: "On track retag flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On release import flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_album_delete": schema.BoolAttribute{
				MarkdownDescription: "On album delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_artist_delete": schema.BoolAttribute{
				MarkdownDescription: "On artist delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSynology name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *NotificationSynologyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationSynologyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSynology

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSynology
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSynologyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSynologyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSynologyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSynology

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSynology current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSynologyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSynologyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSynologyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSynology

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSynology
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSynologyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSynologyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSynologyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSynology current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationSynologyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSynologyResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSynologyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSynologyResourceName+": "+req.ID)
}

func (n *NotificationSynology) write(ctx context.Context, notification *lidarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationSynology) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

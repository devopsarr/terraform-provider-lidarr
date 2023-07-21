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
	notificationSimplepushResourceName   = "notification_simplepush"
	notificationSimplepushImplementation = "Simplepush"
	notificationSimplepushConfigContract = "SimplepushSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSimplepushResource{}
	_ resource.ResourceWithImportState = &NotificationSimplepushResource{}
)

func NewNotificationSimplepushResource() resource.Resource {
	return &NotificationSimplepushResource{}
}

// NotificationSimplepushResource defines the notification implementation.
type NotificationSimplepushResource struct {
	client *lidarr.APIClient
}

// NotificationSimplepush describes the notification data model.
type NotificationSimplepush struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	Event                 types.String `tfsdk:"event"`
	Key                   types.String `tfsdk:"key"`
	ID                    types.Int64  `tfsdk:"id"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnAlbumDelete         types.Bool   `tfsdk:"on_album_delete"`
	OnArtistDelete        types.Bool   `tfsdk:"on_artist_delete"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
	OnDownloadFailure     types.Bool   `tfsdk:"on_download_failure"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
	OnImportFailure       types.Bool   `tfsdk:"on_import_failure"`
}

func (n NotificationSimplepush) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Event:                 n.Event,
		Key:                   n.Key,
		Name:                  n.Name,
		ID:                    n.ID,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		OnAlbumDelete:         n.OnAlbumDelete,
		OnArtistDelete:        n.OnArtistDelete,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnHealthRestored:      n.OnHealthRestored,
		OnDownloadFailure:     n.OnDownloadFailure,
		OnUpgrade:             n.OnUpgrade,
		OnImportFailure:       n.OnImportFailure,
		Implementation:        types.StringValue(notificationSimplepushImplementation),
		ConfigContract:        types.StringValue(notificationSimplepushConfigContract),
	}
}

func (n *NotificationSimplepush) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Event = notification.Event
	n.Key = notification.Key
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnAlbumDelete = notification.OnAlbumDelete
	n.OnArtistDelete = notification.OnArtistDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
	n.OnDownloadFailure = notification.OnDownloadFailure
	n.OnUpgrade = notification.OnUpgrade
	n.OnImportFailure = notification.OnImportFailure
}

func (r *NotificationSimplepushResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSimplepushResourceName
}

func (r *NotificationSimplepushResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Simplepush resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Simplepush](https://wiki.servarr.com/lidarr/supported#simplepush).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_import_failure": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_download_failure": schema.BoolAttribute{
				MarkdownDescription: "On download failure flag.",
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
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSimplepush name.",
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
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Optional:            true,
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationSimplepushResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationSimplepushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSimplepush
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSimplepush current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSimplepush

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSimplepush
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSimplepushResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSimplepushResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSimplepush current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationSimplepushResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSimplepushResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSimplepushResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSimplepushResourceName+": "+req.ID)
}

func (n *NotificationSimplepush) write(ctx context.Context, notification *lidarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationSimplepush) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

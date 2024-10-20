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
	notificationSubsonicResourceName   = "notification_subsonic"
	notificationSubsonicImplementation = "Subsonic"
	notificationSubsonicConfigContract = "SubsonicSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSubsonicResource{}
	_ resource.ResourceWithImportState = &NotificationSubsonicResource{}
)

func NewNotificationSubsonicResource() resource.Resource {
	return &NotificationSubsonicResource{}
}

// NotificationSubsonicResource defines the notification implementation.
type NotificationSubsonicResource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// NotificationSubsonic describes the notification data model.
type NotificationSubsonic struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Host                  types.String `tfsdk:"host"`
	Name                  types.String `tfsdk:"name"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	URLBase               types.String `tfsdk:"url_base"`
	Port                  types.Int64  `tfsdk:"port"`
	ID                    types.Int64  `tfsdk:"id"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Notify                types.Bool   `tfsdk:"notify"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnAlbumDelete         types.Bool   `tfsdk:"on_album_delete"`
	OnArtistDelete        types.Bool   `tfsdk:"on_artist_delete"`
	OnTrackRetag          types.Bool   `tfsdk:"on_track_retag"`
	OnRename              types.Bool   `tfsdk:"on_rename"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
}

func (n NotificationSubsonic) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Port:                  n.Port,
		Host:                  n.Host,
		URLBase:               n.URLBase,
		Password:              n.Password,
		Username:              n.Username,
		Name:                  n.Name,
		ID:                    n.ID,
		UseSSL:                n.UseSSL,
		Notify:                n.Notify,
		UpdateLibrary:         n.UpdateLibrary,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		OnAlbumDelete:         n.OnAlbumDelete,
		OnArtistDelete:        n.OnArtistDelete,
		OnRename:              n.OnRename,
		OnTrackRetag:          n.OnTrackRetag,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnHealthIssue:         n.OnHealthIssue,
		OnUpgrade:             n.OnUpgrade,
		Implementation:        types.StringValue(notificationSubsonicImplementation),
		ConfigContract:        types.StringValue(notificationSubsonicConfigContract),
	}
}

func (n *NotificationSubsonic) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Port = notification.Port
	n.URLBase = notification.URLBase
	n.Host = notification.Host
	n.Password = notification.Password
	n.Username = notification.Username
	n.Name = notification.Name
	n.ID = notification.ID
	n.UseSSL = notification.UseSSL
	n.Notify = notification.Notify
	n.UpdateLibrary = notification.UpdateLibrary
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnAlbumDelete = notification.OnAlbumDelete
	n.OnArtistDelete = notification.OnArtistDelete
	n.OnTrackRetag = notification.OnTrackRetag
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
}

func (r *NotificationSubsonicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSubsonicResourceName
}

func (r *NotificationSubsonicResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Subsonic resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Subsonic](https://wiki.servarr.com/lidarr/supported#xbmc).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
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
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSubsonic name.",
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
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notification flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "URL base.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Required:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationSubsonicResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationSubsonicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSubsonic
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSubsonic current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSubsonic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSubsonic
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, request.GetId()).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSubsonicResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSubsonicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSubsonic current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationSubsonicResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSubsonicResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSubsonicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSubsonicResourceName+": "+req.ID)
}

func (n *NotificationSubsonic) write(ctx context.Context, notification *lidarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationSubsonic) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

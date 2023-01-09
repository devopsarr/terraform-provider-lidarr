package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationKodiResourceName   = "notification_kodi"
	notificationKodiImplementation = "Xbmc"
	notificationKodiConfigContract = "XbmcSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationKodiResource{}
	_ resource.ResourceWithImportState = &NotificationKodiResource{}
)

func NewNotificationKodiResource() resource.Resource {
	return &NotificationKodiResource{}
}

// NotificationKodiResource defines the notification implementation.
type NotificationKodiResource struct {
	client *lidarr.APIClient
}

// NotificationKodi describes the notification data model.
type NotificationKodi struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Host                  types.String `tfsdk:"host"`
	Name                  types.String `tfsdk:"name"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	DisplayTime           types.Int64  `tfsdk:"display_time"`
	Port                  types.Int64  `tfsdk:"port"`
	ID                    types.Int64  `tfsdk:"id"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Notify                types.Bool   `tfsdk:"notify"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	CleanLibrary          types.Bool   `tfsdk:"clean_library"`
	AlwaysUpdate          types.Bool   `tfsdk:"always_update"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnTrackRetag          types.Bool   `tfsdk:"on_track_retag"`
	OnRename              types.Bool   `tfsdk:"on_rename"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
}

func (n NotificationKodi) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Port:                  n.Port,
		Host:                  n.Host,
		DisplayTime:           n.DisplayTime,
		Password:              n.Password,
		Username:              n.Username,
		Name:                  n.Name,
		ID:                    n.ID,
		UseSSL:                n.UseSSL,
		Notify:                n.Notify,
		UpdateLibrary:         n.UpdateLibrary,
		AlwaysUpdate:          n.AlwaysUpdate,
		CleanLibrary:          n.CleanLibrary,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		OnRename:              n.OnRename,
		OnTrackRetag:          n.OnTrackRetag,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnUpgrade:             n.OnUpgrade,
	}
}

func (n *NotificationKodi) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Port = notification.Port
	n.DisplayTime = notification.DisplayTime
	n.Host = notification.Host
	n.Password = notification.Password
	n.Username = notification.Username
	n.Name = notification.Name
	n.ID = notification.ID
	n.UseSSL = notification.UseSSL
	n.Notify = notification.Notify
	n.UpdateLibrary = notification.UpdateLibrary
	n.AlwaysUpdate = notification.AlwaysUpdate
	n.CleanLibrary = notification.CleanLibrary
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnTrackRetag = notification.OnTrackRetag
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
}

func (r *NotificationKodiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationKodiResourceName
}

func (r *NotificationKodiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Kodi resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Kodi](https://wiki.servarr.com/lidarr/supported#xbmc).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
				Required:            true,
			},
			"on_track_retag": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete flag.",
				Required:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationKodi name.",
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
			"clean_library": schema.BoolAttribute{
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
			},
			"always_update": schema.BoolAttribute{
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
			},
			"display_time": schema.Int64Attribute{
				MarkdownDescription: "Display time.",
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

func (r *NotificationKodiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lidarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *lidarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationKodiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationKodi
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKodiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationKodi current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(int(notification.ID.ValueInt64()))).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKodiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationKodi
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKodiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationKodi current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationKodiResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationKodiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationKodiResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationKodi) write(ctx context.Context, notification *lidarr.NotificationResource) {
	genericNotification := Notification{
		OnGrab:                types.BoolValue(notification.GetOnGrab()),
		OnUpgrade:             types.BoolValue(notification.GetOnUpgrade()),
		OnRename:              types.BoolValue(notification.GetOnRename()),
		OnTrackRetag:          types.BoolValue(notification.GetOnTrackRetag()),
		OnReleaseImport:       types.BoolValue(notification.GetOnReleaseImport()),
		OnHealthIssue:         types.BoolValue(notification.GetOnHealthIssue()),
		OnApplicationUpdate:   types.BoolValue(notification.GetOnApplicationUpdate()),
		IncludeHealthWarnings: types.BoolValue(notification.GetIncludeHealthWarnings()),
		ID:                    types.Int64Value(int64(notification.GetId())),
		Name:                  types.StringValue(notification.GetName()),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationKodi) read(ctx context.Context) *lidarr.NotificationResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := lidarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnRename(n.OnRename.ValueBool())
	notification.SetOnTrackRetag(n.OnTrackRetag.ValueBool())
	notification.SetOnReleaseImport(n.OnReleaseImport.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetConfigContract(notificationKodiConfigContract)
	notification.SetImplementation(notificationKodiImplementation)
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetTags(tags)
	notification.SetFields(n.toNotification().readFields(ctx))

	return notification
}

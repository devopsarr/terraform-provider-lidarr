package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/lidarr"
)

const (
	notificationPushbulletResourceName   = "notification_pushbullet"
	notificationPushbulletImplementation = "PushBullet"
	notificationPushbulletConfigContract = "PushBulletSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationPushbulletResource{}
	_ resource.ResourceWithImportState = &NotificationPushbulletResource{}
)

func NewNotificationPushbulletResource() resource.Resource {
	return &NotificationPushbulletResource{}
}

// NotificationPushbulletResource defines the notification implementation.
type NotificationPushbulletResource struct {
	client *lidarr.Lidarr
}

// NotificationPushbullet describes the notification data model.
type NotificationPushbullet struct {
	Tags                  types.Set    `tfsdk:"tags"`
	DeviceIds             types.Set    `tfsdk:"device_ids"`
	ChannelTags           types.Set    `tfsdk:"channel_tags"`
	SenderID              types.String `tfsdk:"sender_id"`
	Name                  types.String `tfsdk:"name"`
	APIKey                types.String `tfsdk:"api_key"`
	ID                    types.Int64  `tfsdk:"id"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnDownloadFailure     types.Bool   `tfsdk:"on_download_failure"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
	OnImportFailure       types.Bool   `tfsdk:"on_import_failure"`
}

func (n NotificationPushbullet) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		DeviceIds:             n.DeviceIds,
		ChannelTags:           n.ChannelTags,
		SenderID:              n.SenderID,
		APIKey:                n.APIKey,
		Name:                  n.Name,
		ID:                    n.ID,
		OnGrab:                n.OnGrab,
		OnReleaseImport:       n.OnReleaseImport,
		IncludeHealthWarnings: n.IncludeHealthWarnings,
		OnApplicationUpdate:   n.OnApplicationUpdate,
		OnHealthIssue:         n.OnHealthIssue,
		OnDownloadFailure:     n.OnDownloadFailure,
		OnUpgrade:             n.OnUpgrade,
		OnImportFailure:       n.OnImportFailure,
	}
}

func (n *NotificationPushbullet) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.DeviceIds = notification.DeviceIds
	n.ChannelTags = notification.ChannelTags
	n.SenderID = notification.SenderID
	n.APIKey = notification.APIKey
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnDownloadFailure = notification.OnDownloadFailure
	n.OnUpgrade = notification.OnUpgrade
	n.OnImportFailure = notification.OnImportFailure
}

func (r *NotificationPushbulletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationPushbulletResourceName
}

func (r *NotificationPushbulletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Pushbullet resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Pushbullet](https://wiki.servarr.com/lidarr/supported#pushbullet).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_import_failure": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_download_failure": schema.BoolAttribute{
				MarkdownDescription: "On movie delete flag.",
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
				MarkdownDescription: "NotificationPushbullet name.",
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
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "List of devices IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "List of channel tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationPushbulletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lidarr.Lidarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *lidarr.Lidarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationPushbulletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationPushbullet
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationPushbullet current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationPushbullet
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationPushbulletResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationPushbulletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationPushbullet

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationPushbullet current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationPushbulletResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationPushbulletResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationPushbulletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationPushbulletResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationPushbullet) write(ctx context.Context, notification *lidarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                types.BoolValue(notification.OnGrab),
		OnImportFailure:       types.BoolValue(notification.OnImportFailure),
		OnUpgrade:             types.BoolValue(notification.OnUpgrade),
		OnDownloadFailure:     types.BoolValue(notification.OnDownloadFailure),
		OnReleaseImport:       types.BoolValue(notification.OnReleaseImport),
		OnHealthIssue:         types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:   types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings: types.BoolValue(notification.IncludeHealthWarnings),
		ID:                    types.Int64Value(notification.ID),
		Name:                  types.StringValue(notification.Name),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationPushbullet) read(ctx context.Context) *lidarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &lidarr.NotificationInput{
		OnGrab:                n.OnGrab.ValueBool(),
		OnImportFailure:       n.OnImportFailure.ValueBool(),
		OnUpgrade:             n.OnUpgrade.ValueBool(),
		OnDownloadFailure:     n.OnDownloadFailure.ValueBool(),
		OnReleaseImport:       n.OnReleaseImport.ValueBool(),
		OnHealthIssue:         n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:   n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings: n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:        notificationPushbulletConfigContract,
		Implementation:        notificationPushbulletImplementation,
		ID:                    n.ID.ValueInt64(),
		Name:                  n.Name.ValueString(),
		Tags:                  tags,
		Fields:                n.toNotification().readFields(ctx),
	}
}
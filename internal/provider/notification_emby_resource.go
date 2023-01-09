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
	notificationEmbyResourceName   = "notification_emby"
	notificationEmbyImplementation = "MediaBrowser"
	notificationEmbyConfigContract = "MediaBrowserSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationEmbyResource{}
	_ resource.ResourceWithImportState = &NotificationEmbyResource{}
)

func NewNotificationEmbyResource() resource.Resource {
	return &NotificationEmbyResource{}
}

// NotificationEmbyResource defines the notification implementation.
type NotificationEmbyResource struct {
	client *lidarr.APIClient
}

// NotificationEmby describes the notification data model.
type NotificationEmby struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Host                  types.String `tfsdk:"host"`
	APIKey                types.String `tfsdk:"api_key"`
	Name                  types.String `tfsdk:"name"`
	ID                    types.Int64  `tfsdk:"id"`
	Port                  types.Int64  `tfsdk:"port"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	Notify                types.Bool   `tfsdk:"notify"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnTrackRetag          types.Bool   `tfsdk:"on_track_retag"`
	OnRename              types.Bool   `tfsdk:"on_rename"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
}

func (n NotificationEmby) toNotification() *Notification {
	return &Notification{
		Tags:                  n.Tags,
		Host:                  n.Host,
		Name:                  n.Name,
		APIKey:                n.APIKey,
		ID:                    n.ID,
		Port:                  n.Port,
		UpdateLibrary:         n.UpdateLibrary,
		Notify:                n.Notify,
		UseSSL:                n.UseSSL,
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

func (n *NotificationEmby) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Host = notification.Host
	n.Name = notification.Name
	n.APIKey = notification.APIKey
	n.ID = notification.ID
	n.UpdateLibrary = notification.UpdateLibrary
	n.Port = notification.Port
	n.Notify = notification.Notify
	n.UseSSL = notification.UseSSL
	n.OnGrab = notification.OnGrab
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnTrackRetag = notification.OnTrackRetag
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
}

func (r *NotificationEmbyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationEmbyResourceName
}

func (r *NotificationEmbyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Emby resource.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Emby](https://wiki.servarr.com/lidarr/supported#mediabrowser).",
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
				MarkdownDescription: "NotificationEmby name.",
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
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Required:            true,
			},
		},
	}
}

func (r *NotificationEmbyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationEmbyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationEmby

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationEmby
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmbyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationEmby

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationEmby current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(int(notification.ID.ValueInt64()))).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmbyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationEmby

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationEmby
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationEmbyResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationEmbyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationEmby

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationEmby current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationEmbyResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationEmbyResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationEmbyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationEmbyResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationEmby) write(ctx context.Context, notification *lidarr.NotificationResource) {
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

func (n *NotificationEmby) read(ctx context.Context) *lidarr.NotificationResource {
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
	notification.SetConfigContract(notificationEmbyConfigContract)
	notification.SetImplementation(notificationEmbyImplementation)
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetTags(tags)
	notification.SetFields(n.toNotification().readFields(ctx))

	return notification
}

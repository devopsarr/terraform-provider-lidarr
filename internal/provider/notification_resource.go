package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const notificationResourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationResource{}
	_ resource.ResourceWithImportState = &NotificationResource{}
)

var notificationFields = helpers.Fields{
	Bools:        []string{"alwaysUpdate", "cleanLibrary", "directMessage", "notify", "requireEncryption", "sendSilently", "updateLibrary", "useEuEndpoint", "useSsl"},
	Strings:      []string{"accessToken", "accessTokenSecret", "apiKey", "aPIKey", "appToken", "arguments", "author", "authToken", "authUser", "avatar", "botToken", "channel", "chatId", "consumerKey", "consumerSecret", "deviceNames", "expires", "from", "host", "icon", "mention", "password", "path", "refreshToken", "senderDomain", "senderId", "server", "signIn", "sound", "token", "urlBase", "url", "userKey", "username", "webHookUrl"},
	Ints:         []string{"method", "port", "priority", "displayTime", "retry", "expire"},
	StringSlices: []string{"channelTags", "deviceIds", "devices", "recipients", "to", "cC", "bcc"},
	IntSlices:    []string{"grabFields", "importFields"},
}

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the notification implementation.
type NotificationResource struct {
	client *lidarr.APIClient
}

// Notification describes the notification data model.
type Notification struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Recipients            types.Set    `tfsdk:"recipients"`
	Devices               types.Set    `tfsdk:"devices"`
	DeviceIds             types.Set    `tfsdk:"device_ids"`
	To                    types.Set    `tfsdk:"to"`
	Cc                    types.Set    `tfsdk:"cc"`
	Bcc                   types.Set    `tfsdk:"bcc"`
	ChannelTags           types.Set    `tfsdk:"channel_tags"`
	ImportFields          types.Set    `tfsdk:"import_fields"`
	GrabFields            types.Set    `tfsdk:"grab_fields"`
	Path                  types.String `tfsdk:"path"`
	RefreshToken          types.String `tfsdk:"refresh_token"`
	WebHookURL            types.String `tfsdk:"web_hook_url"`
	Username              types.String `tfsdk:"username"`
	UserKey               types.String `tfsdk:"user_key"`
	Mention               types.String `tfsdk:"mention"`
	Name                  types.String `tfsdk:"name"`
	Avatar                types.String `tfsdk:"avatar"`
	ConfigContract        types.String `tfsdk:"config_contract"`
	URL                   types.String `tfsdk:"url"`
	URLBase               types.String `tfsdk:"url_base"`
	Token                 types.String `tfsdk:"token"`
	Sound                 types.String `tfsdk:"sound"`
	SignIn                types.String `tfsdk:"sign_in"`
	Server                types.String `tfsdk:"server"`
	SenderID              types.String `tfsdk:"sender_id"`
	BotToken              types.String `tfsdk:"bot_token"`
	SenderDomain          types.String `tfsdk:"sender_domain"`
	Icon                  types.String `tfsdk:"icon"`
	Host                  types.String `tfsdk:"host"`
	From                  types.String `tfsdk:"from"`
	Expires               types.String `tfsdk:"expires"`
	AccessToken           types.String `tfsdk:"access_token"`
	AccessTokenSecret     types.String `tfsdk:"access_token_secret"`
	APIKey                types.String `tfsdk:"api_key"`
	AppToken              types.String `tfsdk:"app_token"`
	Arguments             types.String `tfsdk:"arguments"`
	Author                types.String `tfsdk:"author"`
	AuthToken             types.String `tfsdk:"auth_token"`
	AuthUser              types.String `tfsdk:"auth_user"`
	Implementation        types.String `tfsdk:"implementation"`
	Password              types.String `tfsdk:"password"`
	Channel               types.String `tfsdk:"channel"`
	ChatID                types.String `tfsdk:"chat_id"`
	ConsumerKey           types.String `tfsdk:"consumer_key"`
	ConsumerSecret        types.String `tfsdk:"consumer_secret"`
	DeviceNames           types.String `tfsdk:"device_names"`
	Retry                 types.Int64  `tfsdk:"retry"`
	Expire                types.Int64  `tfsdk:"expire"`
	DisplayTime           types.Int64  `tfsdk:"display_time"`
	Priority              types.Int64  `tfsdk:"priority"`
	Port                  types.Int64  `tfsdk:"port"`
	Method                types.Int64  `tfsdk:"method"`
	ID                    types.Int64  `tfsdk:"id"`
	UpdateLibrary         types.Bool   `tfsdk:"update_library"`
	UseEuEndpoint         types.Bool   `tfsdk:"use_eu_endpoint"`
	Notify                types.Bool   `tfsdk:"notify"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	SendSilently          types.Bool   `tfsdk:"send_silently"`
	RequireEncryption     types.Bool   `tfsdk:"require_encryption"`
	DirectMessage         types.Bool   `tfsdk:"direct_message"`
	CleanLibrary          types.Bool   `tfsdk:"clean_library"`
	AlwaysUpdate          types.Bool   `tfsdk:"always_update"`
	IncludeHealthWarnings types.Bool   `tfsdk:"include_health_warnings"`
	OnGrab                types.Bool   `tfsdk:"on_grab"`
	OnReleaseImport       types.Bool   `tfsdk:"on_release_import"`
	OnUpgrade             types.Bool   `tfsdk:"on_upgrade"`
	OnRename              types.Bool   `tfsdk:"on_rename"`
	OnHealthIssue         types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored      types.Bool   `tfsdk:"on_health_restored"`
	OnDownloadFailure     types.Bool   `tfsdk:"on_download_failure"`
	OnImportFailure       types.Bool   `tfsdk:"on_import_failure"`
	OnTrackRetag          types.Bool   `tfsdk:"on_track_retag"`
	OnApplicationUpdate   types.Bool   `tfsdk:"on_application_update"`
}

func (r *NotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationResourceName
}

func (r *NotificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Generic Notification resource. When possible use a specific resource instead.\nFor more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On release import flag.",
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
			"on_download_failure": schema.BoolAttribute{
				MarkdownDescription: "On download failure flag.",
				Required:            true,
			},
			"on_import_failure": schema.BoolAttribute{
				MarkdownDescription: "On import failure flag.",
				Required:            true,
			},
			"on_track_retag": schema.BoolAttribute{
				MarkdownDescription: "On track retag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
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
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Notification configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Notification implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name.",
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
			"always_update": schema.BoolAttribute{
				MarkdownDescription: "Always update flag.",
				Optional:            true,
				Computed:            true,
			},
			"clean_library": schema.BoolAttribute{
				MarkdownDescription: "Clean library flag.",
				Optional:            true,
				Computed:            true,
			},
			"direct_message": schema.BoolAttribute{
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notify flag.",
				Optional:            true,
				Computed:            true,
			},
			"require_encryption": schema.BoolAttribute{
				MarkdownDescription: "Require encryption flag.",
				Optional:            true,
				Computed:            true,
			},
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Add silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.", // TODO: add values in description
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2, 3, 4, 5, 7),
				},
			},
			"retry": schema.Int64Attribute{
				MarkdownDescription: "Retry.",
				Optional:            true,
				Computed:            true,
			},
			"expire": schema.Int64Attribute{
				MarkdownDescription: "Expire.",
				Optional:            true,
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Optional:            true,
				Computed:            true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App token.",
				Optional:            true,
				Computed:            true,
			},
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Author.",
				Optional:            true,
				Computed:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Auth token.",
				Optional:            true,
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Optional:            true,
				Computed:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "Avatar.",
				Optional:            true,
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Optional:            true,
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				MarkdownDescription: "Channel.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_key": schema.StringAttribute{
				MarkdownDescription: "Consumer key.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_secret": schema.StringAttribute{
				MarkdownDescription: "Consumer secret.",
				Optional:            true,
				Computed:            true,
			},
			"device_names": schema.StringAttribute{
				MarkdownDescription: "Device names.",
				Optional:            true,
				Computed:            true,
			},
			"display_time": schema.Int64Attribute{
				MarkdownDescription: "Display time.",
				Optional:            true,
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Optional:            true,
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon.",
				Optional:            true,
				Computed:            true,
			},
			"mention": schema.StringAttribute{
				MarkdownDescription: "Mention.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password.",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Optional:            true,
				Computed:            true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
			},
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Optional:            true,
				Computed:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "server.",
				Optional:            true,
				Computed:            true,
			},
			"sign_in": schema.StringAttribute{
				MarkdownDescription: "Sign in.",
				Optional:            true,
				Computed:            true,
			},
			"sound": schema.StringAttribute{
				MarkdownDescription: "Sound.",
				Optional:            true,
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "URL base.",
				Optional:            true,
				Computed:            true,
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "User key.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "Web hook url.",
				Optional:            true,
				Computed:            true,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "Channel tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "Device IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"grab_fields": schema.SetAttribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"import_fields": schema.SetAttribute{
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"to": schema.SetAttribute{
				MarkdownDescription: "To.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cc": schema.SetAttribute{
				MarkdownDescription: "Cc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"bcc": schema.SetAttribute{
				MarkdownDescription: "Bcc.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Notification
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Notification current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(int(notification.ID.ValueInt64()))).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *Notification

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Notification
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	// this is needed because of many empty fields are unknown in both plan and read
	var state Notification

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *Notification

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Notification current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationResourceName+": "+req.ID)
}

func (n *Notification) write(ctx context.Context, notification *lidarr.NotificationResource) {
	n.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	n.OnGrab = types.BoolValue(notification.GetOnGrab())
	n.OnImportFailure = types.BoolValue(notification.GetOnImportFailure())
	n.OnUpgrade = types.BoolValue(notification.GetOnUpgrade())
	n.OnRename = types.BoolValue(notification.GetOnRename())
	n.OnReleaseImport = types.BoolValue(notification.GetOnReleaseImport())
	n.OnTrackRetag = types.BoolValue(notification.GetOnTrackRetag())
	n.OnDownloadFailure = types.BoolValue(notification.GetOnDownloadFailure())
	n.OnHealthIssue = types.BoolValue(notification.GetOnHealthIssue())
	n.OnHealthRestored = types.BoolValue(notification.GetOnHealthRestored())
	n.OnApplicationUpdate = types.BoolValue(notification.GetOnApplicationUpdate())
	n.IncludeHealthWarnings = types.BoolValue(notification.GetIncludeHealthWarnings())
	n.ID = types.Int64Value(int64(notification.GetId()))
	n.Name = types.StringValue(notification.GetName())
	n.Implementation = types.StringValue(notification.GetImplementation())
	n.ConfigContract = types.StringValue(notification.GetConfigContract())
	n.GrabFields = types.SetValueMust(types.Int64Type, nil)
	n.ImportFields = types.SetValueMust(types.Int64Type, nil)
	n.ChannelTags = types.SetValueMust(types.StringType, nil)
	n.DeviceIds = types.SetValueMust(types.StringType, nil)
	n.Devices = types.SetValueMust(types.StringType, nil)
	n.Recipients = types.SetValueMust(types.StringType, nil)
	n.To = types.SetValueMust(types.StringType, nil)
	n.Cc = types.SetValueMust(types.StringType, nil)
	n.Bcc = types.SetValueMust(types.StringType, nil)
	helpers.WriteFields(ctx, n, notification.GetFields(), notificationFields)
}

func (n *Notification) read(ctx context.Context) *lidarr.NotificationResource {
	tags := make([]*int32, len(n.Tags.Elements()))
	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := lidarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnReleaseImport(n.OnReleaseImport.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnRename(n.OnRename.ValueBool())
	notification.SetOnTrackRetag(n.OnTrackRetag.ValueBool())
	notification.SetOnDownloadFailure(n.OnDownloadFailure.ValueBool())
	notification.SetOnImportFailure(n.OnImportFailure.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnHealthRestored(n.OnHealthRestored.ValueBool())
	notification.SetOnApplicationUpdate(n.OnApplicationUpdate.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetImplementation(n.Implementation.ValueString())
	notification.SetConfigContract(n.ConfigContract.ValueString())
	notification.SetTags(tags)
	notification.SetFields(helpers.ReadFields(ctx, n, notificationFields))

	return notification
}

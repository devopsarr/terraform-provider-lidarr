package provider

import (
	"context"
	"os"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &LidarrProvider{}

// ScaffoldingProvider defines the provider implementation.
type LidarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Lidarr describes the provider data model.
type Lidarr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *LidarrProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lidarr"
	resp.Version = p.version
}

func (p *LidarrProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Lidarr provider is used to interact with any [Lidarr](https://lidarr.audio/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Lidarr authentication. Can be specified via the `LIDARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Lidarr URL with protocol and port (e.g. `https://test.lidarr.audio:8686`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `LIDARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *LidarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Lidarr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("LIDARR_URL")
	} else {
		url = data.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("LIDARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := lidarr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", key)
	config.Servers[0].URL = url
	client := lidarr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LidarrProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Artists
		NewArtistResource,

		// Download Clients
		NewDownloadClientConfigResource,
		NewDownloadClientResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,
		NewRemotePathMappingResource,

		// Indexers
		NewIndexerResource,
		NewIndexerFilelistResource,
		NewIndexerGazelleResource,
		NewIndexerHeadphonesResource,
		NewIndexerIptorrentsResource,
		NewIndexerNewznabResource,
		NewIndexerNyaaResource,
		NewIndexerRarbgResource,
		NewIndexerRedactedResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorrentleechResource,
		NewIndexerTorznabResource,
		NewIndexerConfigResource,

		// Import Lists
		NewImportListResource,
		NewImportListLidarrResource,
		NewImportListLidarrListResource,
		NewImportListSpotifyArtistsResource,
		NewImportListSpotifyPlaylistsResource,
		NewImportListSpotifyAlbumsResource,
		NewImportListLastFMTagResource,
		NewImportListLastFMUserResource,
		NewImportListHeadphonesResource,
		NewImportListMusicBrainzResource,
		NewImportListExclusionResource,

		// Media Management
		NewMediaManagementResource,
		NewNamingResource,
		NewRootFolderResource,

		// Notifications
		NewNotificationResource,
		NewNotificationBoxcarResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSlackResource,
		NewNotificationSubsonicResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// Profiles
		NewDelayProfileResource,
		NewMetadataProfileResource,

		// Tags
		NewTagResource,
	}
}

func (p *LidarrProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Artists
		NewArtistDataSource,
		NewArtistsDataSource,

		// Download Clients
		NewDownloadClientConfigDataSource,
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,
		NewRemotePathMappingDataSource,
		NewRemotePathMappingsDataSource,

		// Indexers
		NewIndexerConfigDataSource,
		NewIndexerDataSource,
		NewIndexersDataSource,

		// Import Lists
		NewImportListDataSource,
		NewImportListsDataSource,
		NewImportListExclusionDataSource,
		NewImportListExclusionsDataSource,

		// Media Management
		NewMediaManagementDataSource,
		NewNamingDataSource,
		NewRootFolderDataSource,
		NewRootFoldersDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Profiles
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewMetadataProfileDataSource,
		NewMetadataProfilesDataSource,

		// System Status
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LidarrProvider{
			version: version,
		}
	}
}

package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
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
	ExtraHeaders types.Set    `tfsdk:"extra_headers"`
	APIKey       types.String `tfsdk:"api_key"`
	URL          types.String `tfsdk:"url"`
}

// ExtraHeader is part of Lidarr.
type ExtraHeader struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// LidarrData defines auth and client to be used when connecting to Lidarr.
type LidarrData struct {
	Auth   context.Context
	Client *lidarr.APIClient
}

func (p *LidarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lidarr"
	resp.Version = p.version
}

func (p *LidarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
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
			"extra_headers": schema.SetNestedAttribute{
				MarkdownDescription: "Extra headers to be sent along with all Lidarr requests. If this attribute is unset, it can be specified via environment variables following this pattern `LIDARR_EXTRA_HEADER_${Header-Name}=${Header-Value}`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Header name.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Header value.",
							Required:            true,
						},
					},
				},
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

	// Extract URL
	APIURL := data.URL.ValueString()
	if APIURL == "" {
		APIURL = os.Getenv("LIDARR_URL")
	}

	parsedAPIURL, err := url.Parse(APIURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find valid URL",
			"URL cannot parsed",
		)

		return
	}

	// Extract key
	key := data.APIKey.ValueString()
	if key == "" {
		key = os.Getenv("LIDARR_API_KEY")
	}

	if key == "" {
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Init config
	config := lidarr.NewConfiguration()
	// Check extra headers
	if len(data.ExtraHeaders.Elements()) > 0 {
		headers := make([]ExtraHeader, len(data.ExtraHeaders.Elements()))
		resp.Diagnostics.Append(data.ExtraHeaders.ElementsAs(ctx, &headers, false)...)

		for _, header := range headers {
			config.AddDefaultHeader(header.Name.ValueString(), header.Value.ValueString())
		}
	} else {
		env := os.Environ()
		for _, v := range env {
			if strings.HasPrefix(v, "LIDARR_EXTRA_HEADER_") {
				header := strings.Split(v, "=")
				config.AddDefaultHeader(strings.TrimPrefix(header[0], "LIDARR_EXTRA_HEADER_"), header[1])
			}
		}
	}

	// Set context for API calls
	auth := context.WithValue(
		context.Background(),
		lidarr.ContextAPIKeys,
		map[string]lidarr.APIKey{
			"X-Api-Key": {Key: key},
		},
	)
	auth = context.WithValue(auth, lidarr.ContextServerVariables, map[string]string{
		"protocol": parsedAPIURL.Scheme,
		"hostpath": parsedAPIURL.Host,
	})

	lidarrData := LidarrData{
		Auth:   auth,
		Client: lidarr.NewAPIClient(config),
	}
	resp.DataSourceData = &lidarrData
	resp.ResourceData = &lidarrData
}

func (p *LidarrProvider) Resources(_ context.Context) []func() resource.Resource {
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

		// Metadata
		NewMetadataConfigResource,
		NewMetadataResource,
		NewMetadataKodiResource,
		NewMetadataRoksboxResource,
		NewMetadataWdtvResource,

		// Notifications
		NewNotificationResource,
		NewNotificationAppriseResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationNtfyResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSignalResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationSubsonicResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// Profiles
		NewDelayProfileResource,
		NewMetadataProfileResource,
		NewQualityProfileResource,
		NewQualityDefinitionResource,
		NewReleaseProfileResource,
		NewCustomFormatResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
	}
}

func (p *LidarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
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

		// Metadata
		NewMetadataConfigDataSource,
		NewMetadataDataSource,
		NewMetadataConsumersDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Profiles
		NewCustomFormatDataSource,
		NewCustomFormatsDataSource,
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewMetadataProfileDataSource,
		NewMetadataProfilesDataSource,
		NewReleaseProfileDataSource,
		NewReleaseProfilesDataSource,
		NewQualityProfileDataSource,
		NewQualityProfilesDataSource,
		NewQualityDefinitionDataSource,
		NewQualityDefinitionsDataSource,
		NewQualityDataSource,
		NewPrimaryAlbumTypeDataSource,
		NewPrimaryAlbumTypesDataSource,
		NewSecondaryAlbumTypeDataSource,
		NewSecondaryAlbumTypesDataSource,
		NewReleaseStatusDataSource,
		NewReleaseStatusesDataSource,
		NewCustomFormatConditionDataSource,
		NewCustomFormatConditionReleaseGroupDataSource,
		NewCustomFormatConditionReleaseTitleDataSource,
		NewCustomFormatConditionSizeDataSource,

		// System
		NewHostDataSource,
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

// ResourceConfigure is a helper function to set the client for a specific resource.
func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) (context.Context, *lidarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*LidarrData)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *LidarrData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (context.Context, *lidarr.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*LidarrData)
	if !ok {
		resp.Diagnostics.AddError(
			helpers.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *LidarrData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}

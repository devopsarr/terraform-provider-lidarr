package provider

import (
	"context"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const artistDataSourceName = "artist"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ArtistDataSource{}

func NewArtistDataSource() datasource.DataSource {
	return &ArtistDataSource{}
}

// ArtistDataSource defines the artist implementation.
type ArtistDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

func (d *ArtistDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + artistDataSourceName
}

func (d *ArtistDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Artists -->\nSingle [Artist](../resources/artist).",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Artist ID.",
				Computed:            true,
			},
			"artist_name": schema.StringAttribute{
				MarkdownDescription: "Artist name.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full artist path.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Artist status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"foreign_artist_id": schema.StringAttribute{
				MarkdownDescription: "Foreign artist ID.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *ArtistDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ArtistDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Artist

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get artists current value
	response, _, err := d.client.ArtistAPI.ListArtist(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, artistDataSourceName, err))

		return
	}

	data.find(ctx, data.ForeignArtistID.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+artistDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (a *Artist) find(ctx context.Context, ID string, artists []lidarr.ArtistResource, diags *diag.Diagnostics) {
	for _, artist := range artists {
		if artist.GetForeignArtistId() == ID {
			a.write(ctx, &artist, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(artistDataSourceName, "Foreign artist ID", ID))
}

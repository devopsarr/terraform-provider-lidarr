package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const artistsDataSourceName = "artists"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ArtistsDataSource{}

func NewArtistsDataSource() datasource.DataSource {
	return &ArtistsDataSource{}
}

// ArtistsDataSource defines the artists implementation.
type ArtistsDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// Artists describes the artists data model.
type Artists struct {
	Artists types.Set    `tfsdk:"artists"`
	ID      types.String `tfsdk:"id"`
}

func (d *ArtistsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + artistsDataSourceName
}

func (d *ArtistsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Artists -->\nList all available [Artists](../resources/artist).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"artists": schema.SetNestedAttribute{
				MarkdownDescription: "Artist list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *ArtistsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ArtistsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get artists current value
	response, _, err := d.client.ArtistAPI.ListArtist(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, artistsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+artistsDataSourceName)
	// Map response body to resource schema attribute
	artists := make([]Artist, len(response))
	for i, m := range response {
		artists[i].write(ctx, &m, &resp.Diagnostics)
	}

	artistList, diags := types.SetValueFrom(ctx, Artist{}.getType(), artists)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, Artists{Artists: artistList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}

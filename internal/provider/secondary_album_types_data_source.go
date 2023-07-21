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

const secondaryAlbumTypesDataSourceName = "secondary_album_types"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecondaryAlbumTypesDataSource{}

func NewSecondaryAlbumTypesDataSource() datasource.DataSource {
	return &SecondaryAlbumTypesDataSource{}
}

// SecondaryAlbumTypesDataSource defines the secondaryAlbumType implementation.
type SecondaryAlbumTypesDataSource struct {
	client *lidarr.APIClient
}

func (d *SecondaryAlbumTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + secondaryAlbumTypesDataSourceName
}

func (d *SecondaryAlbumTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Secondary Album Types](../data-sources/secondary_album_type).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"elements": schema.SetNestedAttribute{
				MarkdownDescription: "Secondary album type list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Album type ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Album type name.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SecondaryAlbumTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SecondaryAlbumTypesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get secondary album type current value
	response, _, err := d.client.MetadataProfileSchemaApi.GetMetadataprofileSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, secondaryAlbumTypesDataSourceName, err))

		return
	}

	albumTypes := response.GetSecondaryAlbumTypes()

	tflog.Trace(ctx, "read "+secondaryAlbumTypesDataSourceName)
	// Map response body to resource schema attribute
	secondaryTypes := make([]MetadataProfileElement, len(albumTypes))
	for i, t := range albumTypes {
		secondaryTypes[i].writeSecondary(t.AlbumType)
	}

	typeList, diags := types.SetValueFrom(ctx, MetadataProfileElement{}.getType(), secondaryTypes)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, MetadataProfileElements{Elements: typeList, ID: types.StringValue(strconv.Itoa(len(albumTypes)))})...)
}

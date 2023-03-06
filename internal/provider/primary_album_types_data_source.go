package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const primaryAlbumTypesDataSourceName = "primary_album_types"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PrimaryAlbumTypesDataSource{}

func NewPrimaryAlbumTypesDataSource() datasource.DataSource {
	return &PrimaryAlbumTypesDataSource{}
}

// PrimaryAlbumTypesDataSource defines the primaryAlbumType implementation.
type PrimaryAlbumTypesDataSource struct {
	client *lidarr.APIClient
}

func (d *PrimaryAlbumTypesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + primaryAlbumTypesDataSourceName
}

func (d *PrimaryAlbumTypesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Primary Album Types](../data-sources/primary_album_type).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"elements": schema.SetNestedAttribute{
				MarkdownDescription: "Primary album type list.",
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

func (d *PrimaryAlbumTypesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *PrimaryAlbumTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *MetadataProfileElements

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get primary album type current value
	response, _, err := d.client.MetadataProfileSchemaApi.GetMetadataprofileSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, primaryAlbumTypesDataSourceName, err))

		return
	}

	albumTypes := response.GetPrimaryAlbumTypes()

	tflog.Trace(ctx, "read "+primaryAlbumTypesDataSourceName)
	// Map response body to resource schema attribute
	primaryTypes := make([]MetadataProfileElement, len(albumTypes))
	for i, t := range albumTypes {
		primaryTypes[i].writePrimary(t.AlbumType)
	}

	tfsdk.ValueFrom(ctx, primaryTypes, data.Elements.Type(ctx), &data.Elements)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(albumTypes)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

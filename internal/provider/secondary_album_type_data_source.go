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

const secondaryAlbumTypeDataSourceName = "secondary_album_type"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecondaryAlbumTypeDataSource{}

func NewSecondaryAlbumTypeDataSource() datasource.DataSource {
	return &SecondaryAlbumTypeDataSource{}
}

// SecondaryAlbumTypeDataSource defines the secondary album type implementation.
type SecondaryAlbumTypeDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

func (d *SecondaryAlbumTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + secondaryAlbumTypeDataSourceName
}

func (d *SecondaryAlbumTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->\nSingle available Secondary Album Type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "SecondaryAlbumType ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "SecondaryAlbumType.",
				Required:            true,
			},
		},
	}
}

func (d *SecondaryAlbumTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *SecondaryAlbumTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *MetadataProfileElement

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get secondary album type current value
	response, _, err := d.client.MetadataProfileSchemaAPI.GetMetadataprofileSchema(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, secondaryAlbumTypesDataSourceName, err))

		return
	}

	data.findSecondary(data.Name.ValueString(), response.GetSecondaryAlbumTypes(), &resp.Diagnostics)
	tflog.Trace(ctx, "read "+secondaryAlbumTypeDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (e *MetadataProfileElement) writeSecondary(element *lidarr.SecondaryAlbumType) {
	e.Name = types.StringValue(element.GetName())
	e.ID = types.Int64Value(int64(element.GetId()))
}

func (e *MetadataProfileElement) findSecondary(name string, types []lidarr.ProfileSecondaryAlbumTypeItemResource, diags *diag.Diagnostics) {
	for _, t := range types {
		if t.AlbumType.GetName() == name {
			e.writeSecondary(t.AlbumType)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(secondaryAlbumTypeDataSourceName, "name", name))
}

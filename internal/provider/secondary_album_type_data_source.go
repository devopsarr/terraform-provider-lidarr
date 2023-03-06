package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
}

func (d *SecondaryAlbumTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + secondaryAlbumTypeDataSourceName
}

func (d *SecondaryAlbumTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Single available Secondary Album Type.",
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
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SecondaryAlbumTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var secondaryType *MetadataProfileElement

	resp.Diagnostics.Append(req.Config.Get(ctx, &secondaryType)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get secondary album type current value
	response, _, err := d.client.MetadataProfileSchemaApi.GetMetadataprofileSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, secondaryAlbumTypesDataSourceName, err))

		return
	}

	value, err := findSecondaryAlbumType(secondaryType.Name.ValueString(), response.GetSecondaryAlbumTypes())
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", secondaryAlbumTypeDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+secondaryAlbumTypeDataSourceName)
	secondaryType.writeSecondary(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &secondaryType)...)
}

func (e *MetadataProfileElement) writeSecondary(element *lidarr.SecondaryAlbumType) {
	e.Name = types.StringValue(element.GetName())
	e.ID = types.Int64Value(int64(element.GetId()))
}

func findSecondaryAlbumType(name string, types []*lidarr.ProfileSecondaryAlbumTypeItemResource) (*lidarr.SecondaryAlbumType, error) {
	for _, t := range types {
		if t.AlbumType.GetName() == name {
			return t.AlbumType, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(secondaryAlbumTypeDataSourceName, "name", name)
}

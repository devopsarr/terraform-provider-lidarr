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

const primaryAlbumTypeDataSourceName = "primary_album_type"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PrimaryAlbumTypeDataSource{}

func NewPrimaryAlbumTypeDataSource() datasource.DataSource {
	return &PrimaryAlbumTypeDataSource{}
}

// PrimaryAlbumTypeDataSource defines the primary album type implementation.
type PrimaryAlbumTypeDataSource struct {
	client *lidarr.APIClient
}

func (d *PrimaryAlbumTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + primaryAlbumTypeDataSourceName
}

func (d *PrimaryAlbumTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Single available Primary Album Type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "PrimaryAlbumType ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "PrimaryAlbumType.",
				Required:            true,
			},
		},
	}
}

func (d *PrimaryAlbumTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *PrimaryAlbumTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var primaryType *MetadataProfileElement

	resp.Diagnostics.Append(req.Config.Get(ctx, &primaryType)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get primary album type current value
	response, _, err := d.client.MetadataProfileSchemaApi.GetMetadataprofileSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, primaryAlbumTypesDataSourceName, err))

		return
	}

	value, err := findPrimaryAlbumType(primaryType.Name.ValueString(), response.GetPrimaryAlbumTypes())
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", primaryAlbumTypeDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+primaryAlbumTypeDataSourceName)
	primaryType.writePrimary(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &primaryType)...)
}

func (e *MetadataProfileElement) writePrimary(element *lidarr.PrimaryAlbumType) {
	e.Name = types.StringValue(element.GetName())
	e.ID = types.Int64Value(int64(element.GetId()))
}

func findPrimaryAlbumType(name string, types []*lidarr.ProfilePrimaryAlbumTypeItemResource) (*lidarr.PrimaryAlbumType, error) {
	for _, t := range types {
		if t.AlbumType.GetName() == name {
			return t.AlbumType, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(primaryAlbumTypeDataSourceName, "name", name)
}

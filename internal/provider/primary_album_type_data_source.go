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

func (d *PrimaryAlbumTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + primaryAlbumTypeDataSourceName
}

func (d *PrimaryAlbumTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
	var data *MetadataProfileElement

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

	data.findPrimary(data.Name.ValueString(), response.GetPrimaryAlbumTypes(), &resp.Diagnostics)
	tflog.Trace(ctx, "read "+primaryAlbumTypeDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (e *MetadataProfileElement) writePrimary(element *lidarr.PrimaryAlbumType) {
	e.Name = types.StringValue(element.GetName())
	e.ID = types.Int64Value(int64(element.GetId()))
}

func (e *MetadataProfileElement) findPrimary(name string, types []*lidarr.ProfilePrimaryAlbumTypeItemResource, diags *diag.Diagnostics) {
	for _, t := range types {
		if t.AlbumType.GetName() == name {
			e.writePrimary(t.AlbumType)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(primaryAlbumTypeDataSourceName, "name", name))
}

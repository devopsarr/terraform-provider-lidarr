package provider

import (
	"context"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListExclusionDataSourceName = "import_list_exclusion"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListExclusionDataSource{}

func NewImportListExclusionDataSource() datasource.DataSource {
	return &ImportListExclusionDataSource{}
}

// ImportListExclusionDataSource defines the importListExclusion implementation.
type ImportListExclusionDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

func (d *ImportListExclusionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionDataSourceName
}

func (d *ImportListExclusionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->\nSingle [Import List Exclusion](../resources/import_list_exclusion).",
		Attributes: map[string]schema.Attribute{
			"foreign_id": schema.StringAttribute{
				MarkdownDescription: "Musicbrainz ID.",
				Required:            true,
			},
			"artist_name": schema.StringAttribute{
				MarkdownDescription: "Artist to be excluded.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "ImportListExclusion ID.",
				Computed:            true,
			},
		},
	}
}

func (d *ImportListExclusionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ImportListExclusionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportListExclusion

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusions current value
	response, _, err := d.client.ImportListExclusionAPI.ListImportListExclusion(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionDataSourceName, err))

		return
	}

	data.find(data.ForeignID.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+importListExclusionDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *ImportListExclusion) find(foreignID string, importListExclusions []lidarr.ImportListExclusionResource, diags *diag.Diagnostics) {
	for _, t := range importListExclusions {
		if t.GetForeignId() == foreignID {
			i.write(&t)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(importListExclusionDataSourceName, "foreign_id", foreignID))
}

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

const qualityDataSourceName = "quality"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityDataSource{}

func NewQualityDataSource() datasource.DataSource {
	return &QualityDataSource{}
}

// QualityDataSource defines the quality implementation.
type QualityDataSource struct {
	client *lidarr.APIClient
}

// Quality is part of QualityGroup.
type Quality struct {
	Name types.String `tfsdk:"name"`
	ID   types.Int64  `tfsdk:"id"`
}

func (d *QualityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDataSourceName
}

func (d *QualityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single Quality.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality  ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Name.",
				Required:            true,
			},
		},
	}
}

func (d *QualityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *QualityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Quality

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualitys current value
	response, _, err := d.client.QualityDefinitionApi.ListQualityDefinition(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDataSourceName, err))

		return
	}

	quality, err := findQuality(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", qualityDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityDataSourceName)
	data.writeFromDefinition(quality)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findQuality(name string, s []*lidarr.QualityDefinitionResource) (*lidarr.QualityDefinitionResource, error) {
	for _, p := range s {
		if p.Quality.GetName() == name {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(qualityDataSourceName, "name", name)
}

func (q *Quality) writeFromDefinition(quality *lidarr.QualityDefinitionResource) {
	q.ID = types.Int64Value(int64(quality.Quality.GetId()))
	q.Name = types.StringValue(quality.Quality.GetName())
}

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

const releaseStatusDataSourceName = "release_status"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseStatusDataSource{}

func NewReleaseStatusDataSource() datasource.DataSource {
	return &ReleaseStatusDataSource{}
}

// ReleaseStatusDataSource defines the release status type implementation.
type ReleaseStatusDataSource struct {
	client *lidarr.APIClient
}

func (d *ReleaseStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseStatusDataSourceName
}

func (d *ReleaseStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Single available Release Status.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Release Status ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Release Status name.",
				Required:            true,
			},
		},
	}
}

func (d *ReleaseStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ReleaseStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var releaseType *MetadataProfileElement

	resp.Diagnostics.Append(req.Config.Get(ctx, &releaseType)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get release status type current value
	response, _, err := d.client.MetadataProfileSchemaApi.GetMetadataprofileSchema(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseStatusesDataSourceName, err))

		return
	}

	value, err := findReleaseStatus(releaseType.Name.ValueString(), response.GetReleaseStatuses())
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", releaseStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseStatusDataSourceName)
	releaseType.writeRelease(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &releaseType)...)
}

func (e *MetadataProfileElement) writeRelease(element *lidarr.ReleaseStatus) {
	e.Name = types.StringValue(element.GetName())
	e.ID = types.Int64Value(int64(element.GetId()))
}

func findReleaseStatus(name string, types []*lidarr.ProfileReleaseStatusItemResource) (*lidarr.ReleaseStatus, error) {
	for _, t := range types {
		if t.ReleaseStatus.GetName() == name {
			return t.ReleaseStatus, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(releaseStatusDataSourceName, "name", name)
}
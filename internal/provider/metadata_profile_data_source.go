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

const metadataProfileDataSourceName = "metadata_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataProfileDataSource{}

func NewMetadataProfileDataSource() datasource.DataSource {
	return &MetadataProfileDataSource{}
}

// MetadataProfileDataSource defines the metadata profile implementation.
type MetadataProfileDataSource struct {
	client *lidarr.APIClient
}

func (d *MetadataProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataProfileDataSourceName
}

func (d *MetadataProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the metadata server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Metadata Profile](../resources/metadata_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Profile ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata Profile name.",
				Required:            true,
			},
			"primary_album_types": schema.SetAttribute{
				MarkdownDescription: "Primary album types.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"secondary_album_types": schema.SetAttribute{
				MarkdownDescription: "Secondary album types.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"release_statuses": schema.SetAttribute{
				MarkdownDescription: "Release statuses.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *MetadataProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MetadataProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var metadataProfile *MetadataProfile

	resp.Diagnostics.Append(req.Config.Get(ctx, &metadataProfile)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get metadataprofiles current value
	response, _, err := d.client.MetadataProfileApi.ListMetadataProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataProfileDataSourceName, err))

		return
	}

	profile, err := findMetadataProfile(metadataProfile.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", metadataProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataProfileDataSourceName)
	metadataProfile.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadataProfile)...)
}

func findMetadataProfile(name string, profiles []*lidarr.MetadataProfileResource) (*lidarr.MetadataProfileResource, error) {
	for _, p := range profiles {
		if p.GetName() == name {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(metadataProfileDataSourceName, "name", name)
}

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const releaseProfileDataSourceName = "release_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseProfileDataSource{}

func NewReleaseProfileDataSource() datasource.DataSource {
	return &ReleaseProfileDataSource{}
}

// ReleaseProfileDataSource defines the release profile implementation.
type ReleaseProfileDataSource struct {
	client *lidarr.APIClient
}

func (d *ReleaseProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfileDataSourceName
}

func (d *ReleaseProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the release server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Release Profile](../resources/release_profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Release Profile ID.",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled.",
				Computed:            true,
			},
			"include_preferred_when_renaming": schema.BoolAttribute{
				MarkdownDescription: "Include preferred when renaming flag.",
				Computed:            true,
			},
			"indexer_id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID. Set `0` for all.",
				Computed:            true,
			},
			"required": schema.StringAttribute{
				MarkdownDescription: "Required terms. Comma separated list. At least one of `required` and `ignored` must be set.",
				Computed:            true,
			},
			"ignored": schema.StringAttribute{
				MarkdownDescription: "Ignored terms. Comma separated list. At least one of `required` and `ignored` must be set.",
				Computed:            true,
			},
			"preferred": schema.SetNestedAttribute{
				MarkdownDescription: "Preferred terms.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"score": schema.Int64Attribute{
							MarkdownDescription: "Score.",
							Computed:            true,
						},
						"term": schema.StringAttribute{
							MarkdownDescription: "Term.",
							Computed:            true,
						},
					},
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *ReleaseProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ReleaseProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var releaseProfile *ReleaseProfile

	resp.Diagnostics.Append(req.Config.Get(ctx, &releaseProfile)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get releaseprofiles current value
	response, _, err := d.client.ReleaseProfileApi.ListReleaseProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileDataSourceName, err))

		return
	}

	profile, err := findReleaseProfile(releaseProfile.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", releaseProfileDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileDataSourceName)
	releaseProfile.write(ctx, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &releaseProfile)...)
}

func findReleaseProfile(id int64, profiles []*lidarr.ReleaseProfileResource) (*lidarr.ReleaseProfileResource, error) {
	for _, p := range profiles {
		if int64(p.GetId()) == id {
			return p, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(releaseProfileDataSourceName, "id", strconv.Itoa(int(id)))
}

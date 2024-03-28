package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const releaseStatusesDataSourceName = "release_statuses"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseStatusesDataSource{}

func NewReleaseStatusesDataSource() datasource.DataSource {
	return &ReleaseStatusesDataSource{}
}

// ReleaseStatusesDataSource defines the releaseStatus implementation.
type ReleaseStatusesDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

func (d *ReleaseStatusesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseStatusesDataSourceName
}

func (d *ReleaseStatusesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->\nList all available [Release Status](../data-sources/release_status).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"elements": schema.SetNestedAttribute{
				MarkdownDescription: "Release status list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Release status ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Release status name.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ReleaseStatusesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ReleaseStatusesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get release status type current value
	response, _, err := d.client.MetadataProfileSchemaAPI.GetMetadataprofileSchema(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseStatusesDataSourceName, err))

		return
	}

	statuses := response.GetReleaseStatuses()

	tflog.Trace(ctx, "read "+releaseStatusesDataSourceName)
	// Map response body to resource schema attribute
	releaseTypes := make([]MetadataProfileElement, len(statuses))
	for i, t := range statuses {
		releaseTypes[i].writeRelease(t.ReleaseStatus)
	}

	releaseList, diags := types.SetValueFrom(ctx, MetadataProfileElement{}.getType(), releaseTypes)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, MetadataProfileElements{Elements: releaseList, ID: types.StringValue(strconv.Itoa(len(statuses)))})...)
}

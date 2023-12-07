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

const metadataProfilesDataSourceName = "metadata_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataProfilesDataSource{}

func NewMetadataProfilesDataSource() datasource.DataSource {
	return &MetadataProfilesDataSource{}
}

// MetadataProfilesDataSource defines the metadata profiles implementation.
type MetadataProfilesDataSource struct {
	client *lidarr.APIClient
}

// MetadataProfiles describes the metadata profiles data model.
type MetadataProfiles struct {
	MetadataProfiles types.Set    `tfsdk:"metadata_profiles"`
	ID               types.String `tfsdk:"id"`
}

func (d *MetadataProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataProfilesDataSourceName
}

func (d *MetadataProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the metadata server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Metadata Profiles](../resources/metadata_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Metadata Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Metadata Profile ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Metadata Profile name.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *MetadataProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MetadataProfilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get metadataprofiles current value
	response, _, err := d.client.MetadataProfileAPI.ListMetadataProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataProfileResourceName)
	// Map response body to resource schema attribute
	profiles := make([]MetadataProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p, &resp.Diagnostics)
	}

	profileList, diags := types.SetValueFrom(ctx, MetadataProfile{}.getType(), profiles)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, MetadataProfiles{MetadataProfiles: profileList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}

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

const rootFoldersDataSourceName = "root_folders"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFoldersDataSource{}

func NewRootFoldersDataSource() datasource.DataSource {
	return &RootFoldersDataSource{}
}

// RootFoldersDataSource defines the root folders implementation.
type RootFoldersDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

// RootFolders describes the root folders data model.
type RootFolders struct {
	RootFolders types.Set    `tfsdk:"root_folders"`
	ID          types.String `tfsdk:"id"`
}

func (d *RootFoldersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFoldersDataSourceName
}

func (d *RootFoldersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->\nList all available [Root Folders](../resources/root_folder).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"root_folders": schema.SetNestedAttribute{
				MarkdownDescription: "Root Folder list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							MarkdownDescription: "Root Folder absolute path.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Root Folder friendly name.",
							Computed:            true,
						},
						"monitor_option": schema.StringAttribute{
							MarkdownDescription: "Monitor option.",
							Computed:            true,
						},
						"new_item_monitor_option": schema.StringAttribute{
							MarkdownDescription: "New item monitor option.",
							Computed:            true,
						},
						"accessible": schema.BoolAttribute{
							MarkdownDescription: "Access flag.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Root Folder ID.",
							Computed:            true,
						},
						"metadata_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Metadata profile ID.",
							Computed:            true,
						},
						"quality_profile_id": schema.Int64Attribute{
							MarkdownDescription: "Quality profile ID.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *RootFoldersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *RootFoldersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get rootfolders current value
	response, _, err := d.client.RootFolderAPI.ListRootFolder(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFoldersDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFoldersDataSourceName)
	// Map response body to resource schema attribute
	rootFolders := make([]RootFolder, len(response))
	for i, f := range response {
		rootFolders[i].write(ctx, &f, &resp.Diagnostics)
	}

	folderList, diags := types.SetValueFrom(ctx, RootFolder{}.getType(), rootFolders)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, RootFolders{RootFolders: folderList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}

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

const rootFolderDataSourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFolderDataSource{}

func NewRootFolderDataSource() datasource.DataSource {
	return &RootFolderDataSource{}
}

// RootFolderDataSource defines the root folders implementation.
type RootFolderDataSource struct {
	client *lidarr.APIClient
}

func (d *RootFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderDataSourceName
}

func (d *RootFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->Single [Root Folder](../resources/root_folder).",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
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
	}
}

func (d *RootFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RootFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var folder *RootFolder

	resp.Diagnostics.Append(req.Config.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, _, err := d.client.RootFolderApi.ListRootFolder(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderDataSourceName, err))

		return
	}

	// Map response body to resource schema attribute
	rootFolder, err := findRootFolder(folder.Path.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", rootFolderDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFolderDataSourceName)
	folder.write(ctx, rootFolder)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func findRootFolder(path string, folders []*lidarr.RootFolderResource) (*lidarr.RootFolderResource, error) {
	for _, f := range folders {
		if f.GetPath() == path {
			return f, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(rootFolderDataSourceName, "path", path)
}

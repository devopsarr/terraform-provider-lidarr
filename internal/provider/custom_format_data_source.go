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

const customFormatDataSourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatDataSource{}

func NewCustomFormatDataSource() datasource.DataSource {
	return &CustomFormatDataSource{}
}

// CustomFormatDataSource defines the custom_format implementation.
type CustomFormatDataSource struct {
	client *lidarr.APIClient
	auth   context.Context
}

func (d *CustomFormatDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatDataSourceName
}

func (d *CustomFormatDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->\nSingle [Custom Format](../resources/custom_format).",
		Attributes: map[string]schema.Attribute{
			"include_custom_format_when_renaming": schema.BoolAttribute{
				MarkdownDescription: "Include custom format when renaming flag.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Custom Format name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom Format ID.",
				Computed:            true,
			},
			"specifications": schema.SetNestedAttribute{
				MarkdownDescription: "Specifications.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"negate": schema.BoolAttribute{
							MarkdownDescription: "Negate flag.",
							Computed:            true,
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Computed flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Specification name.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Implementation.",
							Computed:            true,
						},
						// Field values
						"value": schema.StringAttribute{
							MarkdownDescription: "Value.",
							Computed:            true,
						},
						"min": schema.Int64Attribute{
							MarkdownDescription: "Min.",
							Computed:            true,
						},
						"max": schema.Int64Attribute{
							MarkdownDescription: "Max.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *CustomFormatDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *CustomFormatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormat

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get customFormat current value
	response, _, err := d.client.CustomFormatAPI.ListCustomFormat(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, customFormatDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatDataSourceName)
	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (c *CustomFormat) find(ctx context.Context, name string, customFormats []lidarr.CustomFormatResource, diags *diag.Diagnostics) {
	for _, i := range customFormats {
		if i.GetName() == name {
			c.write(ctx, &i, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(customFormatDataSourceName, "name", name))
}

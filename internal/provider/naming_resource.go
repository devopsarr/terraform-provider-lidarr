package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const namingResourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NamingResource{}
	_ resource.ResourceWithImportState = &NamingResource{}
)

func NewNamingResource() resource.Resource {
	return &NamingResource{}
}

// NamingResource defines the naming implementation.
type NamingResource struct {
	client *lidarr.APIClient
}

// Naming describes the naming data model.
type Naming struct {
	ArtistFolderFormat       types.String `tfsdk:"artist_folder_format"`
	MultiDiscTrackFormat     types.String `tfsdk:"multi_disc_track_format"`
	StandardTrackFormat      types.String `tfsdk:"standard_track_format"`
	ID                       types.Int64  `tfsdk:"id"`
	RenameTracks             types.Bool   `tfsdk:"rename_tracks"`
	ReplaceIllegalCharacters types.Bool   `tfsdk:"replace_illegal_characters"`
}

func (r *NamingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingResourceName
}

func (r *NamingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Naming resource.\nFor more information refer to [Naming](https://wiki.servarr.com/lidarr/settings#community-naming-suggestions) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Naming ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"rename_tracks": schema.BoolAttribute{
				MarkdownDescription: "Lidarr will use the existing file name if false.",
				Required:            true,
			},
			"replace_illegal_characters": schema.BoolAttribute{
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Required:            true,
			},
			"artist_folder_format": schema.StringAttribute{
				MarkdownDescription: "Artist folder format.",
				Required:            true,
			},
			"multi_disc_track_format": schema.StringAttribute{
				MarkdownDescription: "Multi disc track format.",
				Required:            true,
			},
			"standard_track_format": schema.StringAttribute{
				MarkdownDescription: "Standard track formatss.",
				Required:            true,
			},
		},
	}
}

func (r *NamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init call if we remove this it the very first update on a brand new instance will fail
	if _, _, err := r.client.NamingConfigApi.GetNamingConfig(ctx).Execute(); err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError("init", namingResourceName, err))

		return
	}

	// Build Create resource
	request := naming.read()
	request.SetId(1)

	// Create new Naming
	response, _, err := r.client.NamingConfigApi.UpdateNamingConfig(ctx, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var naming *Naming

	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get naming current value
	response, _, err := r.client.NamingConfigApi.GetNamingConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var naming *Naming

	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := naming.read()

	// Update Naming
	response, _, err := r.client.NamingConfigApi.UpdateNamingConfig(ctx, strconv.Itoa(int(request.GetId()))).NamingConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, namingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+namingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	naming.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (r *NamingResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Naming cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+namingResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *NamingResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+namingResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (n *Naming) write(naming *lidarr.NamingConfigResource) {
	n.RenameTracks = types.BoolValue(naming.GetRenameTracks())
	n.ReplaceIllegalCharacters = types.BoolValue(naming.GetReplaceIllegalCharacters())
	n.ID = types.Int64Value(int64(naming.GetId()))
	n.ArtistFolderFormat = types.StringValue(naming.GetArtistFolderFormat())
	n.MultiDiscTrackFormat = types.StringValue(naming.GetMultiDiscTrackFormat())
	n.StandardTrackFormat = types.StringValue(naming.GetStandardTrackFormat())
}

func (n *Naming) read() *lidarr.NamingConfigResource {
	naming := lidarr.NewNamingConfigResource()
	naming.SetId(int32(n.ID.ValueInt64()))
	naming.SetRenameTracks(n.RenameTracks.ValueBool())
	naming.SetReplaceIllegalCharacters(n.ReplaceIllegalCharacters.ValueBool())
	naming.SetArtistFolderFormat(n.ArtistFolderFormat.ValueString())
	naming.SetMultiDiscTrackFormat(n.MultiDiscTrackFormat.ValueString())
	naming.SetStandardTrackFormat(n.StandardTrackFormat.ValueString())

	return naming
}

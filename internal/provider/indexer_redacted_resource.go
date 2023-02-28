package provider

import (
	"context"
	"fmt"
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

const (
	indexerRedactedResourceName   = "indexer_redacted"
	indexerRedactedImplementation = "Redacted"
	indexerRedactedConfigContract = "RedactedSettings"
	indexerRedactedProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerRedactedResource{}
	_ resource.ResourceWithImportState = &IndexerRedactedResource{}
)

func NewIndexerRedactedResource() resource.Resource {
	return &IndexerRedactedResource{}
}

// IndexerRedactedResource defines the Redacted indexer implementation.
type IndexerRedactedResource struct {
	client *lidarr.APIClient
}

// IndexerRedacted describes the Redacted indexer data model.
type IndexerRedacted struct {
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	Passkey                 types.String  `tfsdk:"passkey"`
	APIKey                  types.String  `tfsdk:"api_key"`
	Priority                types.Int64   `tfsdk:"priority"`
	ID                      types.Int64   `tfsdk:"id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	EarlyReleaseLimit       types.Int64   `tfsdk:"early_release_limit"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	DiscographySeedTime     types.Int64   `tfsdk:"discography_seed_time"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	UseFreeleechToken       types.Bool    `tfsdk:"use_freeleech_token"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerRedacted) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		ID:                      i.ID,
		Name:                    i.Name,
		UseFreeleechToken:       i.UseFreeleechToken,
		MinimumSeeders:          i.MinimumSeeders,
		EarlyReleaseLimit:       i.EarlyReleaseLimit,
		SeedTime:                i.SeedTime,
		DiscographySeedTime:     i.DiscographySeedTime,
		SeedRatio:               i.SeedRatio,
		Passkey:                 i.Passkey,
		APIKey:                  i.APIKey,
		Tags:                    i.Tags,
		Implementation:          types.StringValue(indexerRedactedImplementation),
		ConfigContract:          types.StringValue(indexerRedactedConfigContract),
		Protocol:                types.StringValue(indexerRedactedProtocol),
	}
}

func (i *IndexerRedacted) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.UseFreeleechToken = indexer.UseFreeleechToken
	i.EarlyReleaseLimit = indexer.EarlyReleaseLimit
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeedTime = indexer.SeedTime
	i.DiscographySeedTime = indexer.DiscographySeedTime
	i.SeedRatio = indexer.SeedRatio
	i.Passkey = indexer.Passkey
	i.APIKey = indexer.APIKey
	i.Tags = indexer.Tags
}

func (r *IndexerRedactedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerRedactedResourceName
}

func (r *IndexerRedactedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Redacted resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/lidarr/settings#indexers) and [Redacted](https://wiki.servarr.com/lidarr/supported#redacted).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "IndexerRedacted name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerRedacted ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"use_freeleech_token": schema.BoolAttribute{
				MarkdownDescription: "Use freeleech token flag.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
			},
			"early_release_limit": schema.Int64Attribute{
				MarkdownDescription: "Early release limit.",
				Optional:            true,
				Computed:            true,
			},
			"seed_time": schema.Int64Attribute{
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
			},
			"discography_seed_time": schema.Int64Attribute{
				MarkdownDescription: "Discography seed time.",
				Optional:            true,
				Computed:            true,
			},
			"seed_ratio": schema.Float64Attribute{
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"passkey": schema.StringAttribute{
				MarkdownDescription: "passkey.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *IndexerRedactedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerRedactedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerRedacted

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerRedacted
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerRedactedResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerRedactedResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerRedactedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerRedacted

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerRedacted current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerRedactedResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerRedactedResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerRedactedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerRedacted

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerRedacted
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+indexerRedactedResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerRedactedResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerRedactedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerRedacted

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerRedacted current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerRedactedResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerRedactedResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerRedactedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerRedactedResourceName+": "+req.ID)
}

func (i *IndexerRedacted) write(ctx context.Context, indexer *lidarr.IndexerResource) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerRedacted) read(ctx context.Context) *lidarr.IndexerResource {
	return i.toIndexer().read(ctx)
}

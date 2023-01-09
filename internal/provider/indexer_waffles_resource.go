package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerWafflesResourceName   = "indexer_waffles"
	indexerWafflesImplementation = "Waffles"
	indexerWafflesConfigContract = "WafflesSettings"
	indexerWafflesProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerWafflesResource{}
	_ resource.ResourceWithImportState = &IndexerWafflesResource{}
)

func NewIndexerWafflesResource() resource.Resource {
	return &IndexerWafflesResource{}
}

// IndexerWafflesResource defines the Waffles indexer implementation.
type IndexerWafflesResource struct {
	client *lidarr.APIClient
}

// IndexerWaffles describes the Waffles indexer data model.
type IndexerWaffles struct {
	Tags                    types.Set     `tfsdk:"tags"`
	Name                    types.String  `tfsdk:"name"`
	UserID                  types.String  `tfsdk:"user_id"`
	RSSPasskey              types.String  `tfsdk:"rss_passkey"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	Priority                types.Int64   `tfsdk:"priority"`
	ID                      types.Int64   `tfsdk:"id"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	EarlyReleaseLimit       types.Int64   `tfsdk:"early_release_limit"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	DiscographySeedTime     types.Int64   `tfsdk:"discography_seed_time"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
}

func (i IndexerWaffles) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		ID:                      i.ID,
		Name:                    i.Name,
		MinimumSeeders:          i.MinimumSeeders,
		EarlyReleaseLimit:       i.EarlyReleaseLimit,
		SeedTime:                i.SeedTime,
		DiscographySeedTime:     i.DiscographySeedTime,
		SeedRatio:               i.SeedRatio,
		RSSPasskey:              i.RSSPasskey,
		UserID:                  i.UserID,
		BaseURL:                 i.BaseURL,
		Tags:                    i.Tags,
	}
}

func (i *IndexerWaffles) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.EarlyReleaseLimit = indexer.EarlyReleaseLimit
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeedTime = indexer.SeedTime
	i.DiscographySeedTime = indexer.DiscographySeedTime
	i.SeedRatio = indexer.SeedRatio
	i.RSSPasskey = indexer.RSSPasskey
	i.UserID = indexer.UserID
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
}

func (r *IndexerWafflesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerWafflesResourceName
}

func (r *IndexerWafflesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Waffles resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/lidarr/settings#indexers) and [Waffles](https://wiki.servarr.com/lidarr/supported#waffles).",
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
				MarkdownDescription: "IndexerWaffles name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerWaffles ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID.",
				Required:            true,
			},
			"rss_passkey": schema.StringAttribute{
				MarkdownDescription: "RSS Passkey.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *IndexerWafflesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*lidarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *lidarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IndexerWafflesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerWaffles

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerWaffles
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerWafflesResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerWafflesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerWafflesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerWaffles

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerWaffles current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerWafflesResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerWafflesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerWafflesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerWaffles

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerWaffles
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update "+indexerWafflesResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerWafflesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerWafflesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerWaffles

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerWaffles current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerWafflesResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerWafflesResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerWafflesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerWafflesResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *IndexerWaffles) write(ctx context.Context, indexer *lidarr.IndexerResource) {
	genericIndexer := Indexer{
		EnableAutomaticSearch:   types.BoolValue(indexer.GetEnableAutomaticSearch()),
		EnableInteractiveSearch: types.BoolValue(indexer.GetEnableInteractiveSearch()),
		EnableRss:               types.BoolValue(indexer.GetEnableRss()),
		Priority:                types.Int64Value(int64(indexer.GetPriority())),
		ID:                      types.Int64Value(int64(indexer.GetId())),
		Name:                    types.StringValue(indexer.GetName()),
	}
	genericIndexer.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	genericIndexer.writeFields(ctx, indexer.Fields)
	i.fromIndexer(&genericIndexer)
}

func (i *IndexerWaffles) read(ctx context.Context) *lidarr.IndexerResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	indexer := lidarr.NewIndexerResource()
	indexer.SetEnableAutomaticSearch(i.EnableAutomaticSearch.ValueBool())
	indexer.SetEnableInteractiveSearch(i.EnableInteractiveSearch.ValueBool())
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(indexerWafflesConfigContract)
	indexer.SetImplementation(indexerWafflesImplementation)
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(indexerWafflesProtocol)
	indexer.SetTags(tags)
	indexer.SetFields(i.toIndexer().readFields(ctx))

	return indexer
}

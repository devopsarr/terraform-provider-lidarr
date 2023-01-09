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
	indexerHeadphonesResourceName   = "indexer_headphones"
	indexerHeadphonesImplementation = "Headphones"
	indexerHeadphonesConfigContract = "HeadphonesSettings"
	indexerHeadphonesProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerHeadphonesResource{}
	_ resource.ResourceWithImportState = &IndexerHeadphonesResource{}
)

func NewIndexerHeadphonesResource() resource.Resource {
	return &IndexerHeadphonesResource{}
}

// IndexerHeadphonesResource defines the Headphones indexer implementation.
type IndexerHeadphonesResource struct {
	client *lidarr.APIClient
}

// IndexerHeadphones describes the Headphones indexer data model.
type IndexerHeadphones struct {
	Tags                    types.Set    `tfsdk:"tags"`
	Categories              types.Set    `tfsdk:"categories"`
	Name                    types.String `tfsdk:"name"`
	Username                types.String `tfsdk:"username"`
	Password                types.String `tfsdk:"password"`
	Priority                types.Int64  `tfsdk:"priority"`
	ID                      types.Int64  `tfsdk:"id"`
	EarlyReleaseLimit       types.Int64  `tfsdk:"early_release_limit"`
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
}

func (i IndexerHeadphones) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		ID:                      i.ID,
		Name:                    i.Name,
		EarlyReleaseLimit:       i.EarlyReleaseLimit,
		Username:                i.Username,
		Password:                i.Password,
		Categories:              i.Categories,
		Tags:                    i.Tags,
	}
}

func (i *IndexerHeadphones) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.EarlyReleaseLimit = indexer.EarlyReleaseLimit
	i.Username = indexer.Username
	i.Password = indexer.Password
	i.Categories = indexer.Categories
	i.Tags = indexer.Tags
}

func (r *IndexerHeadphonesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerHeadphonesResourceName
}

func (r *IndexerHeadphonesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Headphones resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/lidarr/settings#indexers) and [Headphones](https://wiki.servarr.com/lidarr/supported#headphones).",
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
				MarkdownDescription: "IndexerHeadphones name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerHeadphones ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"early_release_limit": schema.Int64Attribute{
				MarkdownDescription: "Early release limit.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Required:            true,
				Sensitive:           true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Series list.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerHeadphonesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IndexerHeadphonesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerHeadphones

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerHeadphones
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", indexerHeadphonesResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerHeadphonesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerHeadphonesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerHeadphones

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerHeadphones current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerHeadphonesResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerHeadphonesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerHeadphonesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerHeadphones

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerHeadphones
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update "+indexerHeadphonesResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerHeadphonesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerHeadphonesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerHeadphones

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerHeadphones current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerHeadphonesResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerHeadphonesResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerHeadphonesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+indexerHeadphonesResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (i *IndexerHeadphones) write(ctx context.Context, indexer *lidarr.IndexerResource) {
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

func (i *IndexerHeadphones) read(ctx context.Context) *lidarr.IndexerResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, i.Tags, &tags)

	indexer := lidarr.NewIndexerResource()
	indexer.SetEnableAutomaticSearch(i.EnableAutomaticSearch.ValueBool())
	indexer.SetEnableInteractiveSearch(i.EnableInteractiveSearch.ValueBool())
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(indexerHeadphonesConfigContract)
	indexer.SetImplementation(indexerHeadphonesImplementation)
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(indexerHeadphonesProtocol)
	indexer.SetTags(tags)
	indexer.SetFields(i.toIndexer().readFields(ctx))

	return indexer
}

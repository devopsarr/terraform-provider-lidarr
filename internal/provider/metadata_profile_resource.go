package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataProfileResourceName = "metadata_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataProfileResource{}
	_ resource.ResourceWithImportState = &MetadataProfileResource{}
)

func NewMetadataProfileResource() resource.Resource {
	return &MetadataProfileResource{}
}

// MetadataProfileResource defines the metadata profile implementation.
type MetadataProfileResource struct {
	client *lidarr.APIClient
}

// MetadataProfile describes the metadata profile data model.
type MetadataProfile struct {
	PrimaryAlbumTypes   types.Set    `tfsdk:"primary_album_types"`
	SecondaryAlbumTypes types.Set    `tfsdk:"secondary_album_types"`
	ReleaseStatuses     types.Set    `tfsdk:"release_statuses"`
	Name                types.String `tfsdk:"name"`
	ID                  types.Int64  `tfsdk:"id"`
}

func (p MetadataProfile) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"id":                    types.Int64Type,
			"name":                  types.StringType,
			"release_statuses":      types.SetType{}.WithElementType(types.Int64Type),
			"secondary_album_types": types.SetType{}.WithElementType(types.Int64Type),
			"primary_album_types":   types.SetType{}.WithElementType(types.Int64Type),
		})
}

// MetadataProfileElement describes the metadata profile element data model.
type MetadataProfileElement struct {
	Name types.String `tfsdk:"name"`
	ID   types.Int64  `tfsdk:"id"`
}

func (m MetadataProfileElement) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"id":   types.Int64Type,
			"name": types.StringType,
		})
}

// MetadataProfileElements describes the metadata profile elements data model.
type MetadataProfileElements struct {
	Elements types.Set    `tfsdk:"elements"`
	ID       types.String `tfsdk:"id"`
}

func (r *MetadataProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Metadata Profile resource.\nFor more information refer to [Metadata Profile](https://wiki.servarr.com/lidarr/settings#metadata-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata Profile name.",
				Required:            true,
			},
			"primary_album_types": schema.SetAttribute{
				MarkdownDescription: "Primary album types.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"secondary_album_types": schema.SetAttribute{
				MarkdownDescription: "Secondary album types.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"release_statuses": schema.SetAttribute{
				MarkdownDescription: "Release statuses.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *MetadataProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataProfileResourceName
}

func (r *MetadataProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataProfile
	request := profile.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataProfileApi.CreateMetadataProfile(ctx).MetadataProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get metadataProfile current value
	response, _, err := r.client.MetadataProfileApi.GetMetadataProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataProfile
	request := profile.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataProfileApi.UpdateMetadataProfile(ctx, strconv.Itoa(int(request.GetId()))).MetadataProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete metadataProfile current value
	_, err := r.client.MetadataProfileApi.DeleteMetadataProfile(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataProfileResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataProfileResourceName+": "+req.ID)
}

func (p *MetadataProfile) write(ctx context.Context, profile *lidarr.MetadataProfileResource, diags *diag.Diagnostics) {
	var (
		tempDiag                    diag.Diagnostics
		primary, secondary, release []*int32
	)

	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())

	for _, p := range profile.GetPrimaryAlbumTypes() {
		if p.GetAllowed() {
			primary = append(primary, p.GetAlbumType().Id)
		}
	}

	for _, p := range profile.GetSecondaryAlbumTypes() {
		if p.GetAllowed() {
			secondary = append(secondary, p.GetAlbumType().Id)
		}
	}

	for _, p := range profile.GetReleaseStatuses() {
		if p.GetAllowed() {
			release = append(release, p.GetReleaseStatus().Id)
		}
	}

	p.PrimaryAlbumTypes, tempDiag = types.SetValueFrom(ctx, types.Int64Type, primary)
	diags.Append(tempDiag...)
	p.SecondaryAlbumTypes, tempDiag = types.SetValueFrom(ctx, types.Int64Type, secondary)
	diags.Append(tempDiag...)
	p.ReleaseStatuses, tempDiag = types.SetValueFrom(ctx, types.Int64Type, release)
	diags.Append(tempDiag...)
}

func (p *MetadataProfile) read(ctx context.Context, diags *diag.Diagnostics) *lidarr.MetadataProfileResource {
	var primary, secondary, release []*int64

	diags.Append(p.PrimaryAlbumTypes.ElementsAs(ctx, &primary, true)...)
	diags.Append(p.SecondaryAlbumTypes.ElementsAs(ctx, &secondary, true)...)
	diags.Append(p.ReleaseStatuses.ElementsAs(ctx, &release, true)...)

	primaryTypes := make([]*lidarr.ProfilePrimaryAlbumTypeItemResource, len(primary))
	for i, e := range primary {
		primaryTypes[i] = lidarr.NewProfilePrimaryAlbumTypeItemResource()
		element := lidarr.NewPrimaryAlbumType()
		element.SetId(int32(*e))
		primaryTypes[i].SetAlbumType(*element)
		primaryTypes[i].SetAllowed(true)
	}

	secondaryTypes := make([]*lidarr.ProfileSecondaryAlbumTypeItemResource, len(secondary))
	for i, e := range secondary {
		secondaryTypes[i] = lidarr.NewProfileSecondaryAlbumTypeItemResource()
		element := lidarr.NewSecondaryAlbumType()
		element.SetId(int32(*e))
		secondaryTypes[i].SetAlbumType(*element)
		secondaryTypes[i].SetAllowed(true)
	}

	releaseStatuses := make([]*lidarr.ProfileReleaseStatusItemResource, len(release))
	for i, e := range release {
		releaseStatuses[i] = lidarr.NewProfileReleaseStatusItemResource()
		element := lidarr.NewReleaseStatus()
		element.SetId(int32(*e))
		releaseStatuses[i].SetReleaseStatus(*element)
		releaseStatuses[i].SetAllowed(true)
	}

	profile := lidarr.NewMetadataProfileResource()
	profile.SetName(p.Name.ValueString())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetPrimaryAlbumTypes(primaryTypes)
	profile.SetSecondaryAlbumTypes(secondaryTypes)
	profile.SetReleaseStatuses(releaseStatuses)

	return profile
}

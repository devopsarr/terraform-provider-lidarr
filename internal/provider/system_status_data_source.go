package provider

import (
	"context"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/devopsarr/terraform-provider-lidarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const systemStatusDataSourceName = "system_status"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *lidarr.APIClient
}

// SystemStatus describes the system status data model.
type SystemStatus struct {
	RuntimeName                   types.String `tfsdk:"runtime_name"`
	AppData                       types.String `tfsdk:"app_data"`
	OsName                        types.String `tfsdk:"os_name"`
	PackageVersion                types.String `tfsdk:"package_version"`
	Mode                          types.String `tfsdk:"mode"`
	PackageUpdateMechanism        types.String `tfsdk:"package_update_mechanism"`
	PackageAuthor                 types.String `tfsdk:"package_author"`
	InstanceName                  types.String `tfsdk:"instance_name"`
	AppName                       types.String `tfsdk:"app_name"`
	Branch                        types.String `tfsdk:"branch"`
	StartupPath                   types.String `tfsdk:"startup_path"`
	RuntimeVersion                types.String `tfsdk:"runtime_version"`
	StartTime                     types.String `tfsdk:"start_time"`
	BuildTime                     types.String `tfsdk:"build_time"`
	Version                       types.String `tfsdk:"version"`
	Authentication                types.String `tfsdk:"authentication"`
	OsVersion                     types.String `tfsdk:"os_version"`
	DatabaseVersion               types.String `tfsdk:"database_version"`
	DatabaseType                  types.String `tfsdk:"database_type"`
	URLBase                       types.String `tfsdk:"url_base"`
	PackageUpdateMechanismMessage types.String `tfsdk:"package_update_mechanism_message"`
	MigrationVersion              types.Int64  `tfsdk:"migration_version"`
	ID                            types.Int64  `tfsdk:"id"`
	IsDocker                      types.Bool   `tfsdk:"is_docker"`
	IsDebug                       types.Bool   `tfsdk:"is_debug"`
	IsNetCore                     types.Bool   `tfsdk:"is_net_core"`
	IsAdmin                       types.Bool   `tfsdk:"is_admin"`
	IsProduction                  types.Bool   `tfsdk:"is_production"`
	IsWindows                     types.Bool   `tfsdk:"is_windows"`
	IsOsx                         types.Bool   `tfsdk:"is_osx"`
	IsLinux                       types.Bool   `tfsdk:"is_linux"`
	IsUserInteractive             types.Bool   `tfsdk:"is_user_interactive"`
}

func (d *SystemStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_status"
}

func (d *SystemStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:System -->System Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/lidarr/system#status) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"is_debug": schema.BoolAttribute{
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
			},
			"is_production": schema.BoolAttribute{
				MarkdownDescription: "Is production flag.",
				Computed:            true,
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
			},
			"is_user_interactive": schema.BoolAttribute{
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
			},
			"is_net_core": schema.BoolAttribute{
				MarkdownDescription: "Is net core flag.",
				Computed:            true,
			},
			"is_linux": schema.BoolAttribute{
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
			},
			"is_osx": schema.BoolAttribute{
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
			},
			"is_windows": schema.BoolAttribute{
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
			},
			"is_docker": schema.BoolAttribute{
				MarkdownDescription: "Is docker flag.",
				Computed:            true,
			},
			"migration_version": schema.Int64Attribute{
				MarkdownDescription: "Migration version.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version.",
				Computed:            true,
			},
			"startup_path": schema.StringAttribute{
				MarkdownDescription: "Startup path.",
				Computed:            true,
			},
			"app_data": schema.StringAttribute{
				MarkdownDescription: "App data folder.",
				Computed:            true,
			},
			"os_name": schema.StringAttribute{
				MarkdownDescription: "OS name.",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode.",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch.",
				Computed:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version.",
				Computed:            true,
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name.",
				Computed:            true,
			},
			"build_time": schema.StringAttribute{
				MarkdownDescription: "Build time.",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Start time.",
				Computed:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "App name.",
				Computed:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Computed:            true,
			},
			"package_author": schema.StringAttribute{
				MarkdownDescription: "Package author.",
				Computed:            true,
			},
			"package_update_mechanism": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
			},
			"package_version": schema.StringAttribute{
				MarkdownDescription: "Package version.",
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				MarkdownDescription: "OS version.",
				Computed:            true,
			},
			"database_version": schema.StringAttribute{
				MarkdownDescription: "Database version.",
				Computed:            true,
			},
			"database_type": schema.StringAttribute{
				MarkdownDescription: "Database type.",
				Computed:            true,
			},
			"package_update_mechanism_message": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism message.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SystemStatusDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get system status current value
	response, _, err := d.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}

func (s *SystemStatus) write(status *lidarr.SystemResource) {
	s.IsDebug = types.BoolValue(status.GetIsDebug())
	s.IsProduction = types.BoolValue(status.GetIsProduction())
	s.IsAdmin = types.BoolValue(status.GetIsAdmin())
	s.IsUserInteractive = types.BoolValue(status.GetIsUserInteractive())
	s.IsNetCore = types.BoolValue(status.GetIsNetCore())
	s.IsLinux = types.BoolValue(status.GetIsLinux())
	s.IsOsx = types.BoolValue(status.GetIsOsx())
	s.IsWindows = types.BoolValue(status.GetIsWindows())
	s.IsDocker = types.BoolValue(status.GetIsDocker())
	s.ID = types.Int64Value(int64(1))
	s.MigrationVersion = types.Int64Value(int64(status.GetMigrationVersion()))
	s.Version = types.StringValue(status.GetVersion())
	s.StartupPath = types.StringValue(status.GetStartupPath())
	s.AppData = types.StringValue(status.GetAppData())
	s.OsName = types.StringValue(status.GetOsName())
	s.Mode = types.StringValue(string(status.GetMode()))
	s.Branch = types.StringValue(status.GetBranch())
	s.Authentication = types.StringValue(string(status.GetAuthentication()))
	s.RuntimeVersion = types.StringValue(status.GetRuntimeVersion())
	s.RuntimeName = types.StringValue(status.GetRuntimeName())
	s.BuildTime = types.StringValue(status.GetBuildTime().String())
	s.StartTime = types.StringValue(status.GetStartTime().String())
	s.AppName = types.StringValue(status.GetAppName())
	s.InstanceName = types.StringValue(status.GetInstanceName())
	s.PackageAuthor = types.StringValue(status.GetPackageAuthor())
	s.PackageUpdateMechanism = types.StringValue(string(status.GetPackageUpdateMechanism()))
	s.PackageVersion = types.StringValue(status.GetPackageVersion())
	s.OsVersion = types.StringValue(status.GetOsVersion())
	s.DatabaseVersion = types.StringValue(status.GetDatabaseVersion())
	s.DatabaseType = types.StringValue(string(status.GetDatabaseType()))
	s.URLBase = types.StringValue(status.GetUrlBase())
	s.PackageUpdateMechanismMessage = types.StringValue(status.GetPackageUpdateMechanismMessage())
}

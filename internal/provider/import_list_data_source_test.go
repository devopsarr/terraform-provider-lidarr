package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_import_list.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_import_list.test", "should_monitor", "none")),
			},
		},
	})
}

const testAccImportListDataSourceConfig = `
resource "lidarr_import_list" "test" {
	enable_automatic_add = false
	should_monitor = "none"
	should_search = false
	list_type = "program"
	root_folder_path = "/config"
	monitor_new_items = "all"
	quality_profile_id = 1
	metadata_profile_id = 1
	name = "importListDataTest"
	implementation = "LidarrImport"
	config_contract = "LidarrSettings"
	base_url = "http://127.0.0.1:8686"
	api_key = "testAPIKey"
}

data "lidarr_import_list" "test" {
	name = lidarr_import_list.test.name
}
`

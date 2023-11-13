package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListResourceConfig("importListResourceTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListResourceConfig("importListResourceTest", "entireArtist"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list.test", "should_monitor", "entireArtist"),
					resource.TestCheckResourceAttrSet("lidarr_import_list.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListResourceConfig("importListResourceTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListResourceConfig("importListResourceTest", "specificAlbum"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list.test", "should_monitor", "specificAlbum"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lidarr_import_list.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "lidarr_import_list" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		should_search = false
		list_type = "program"
		root_folder_path = "/config"
		monitor_new_items = "all"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		implementation = "LidarrImport"
    	config_contract = "LidarrSettings"
		base_url = "http://127.0.0.1:8686"
		api_key = "testAPIKey"
		tag_ids = [1,2]
		profile_ids = [1]
		tags = []
	}`, monitor, name)
}

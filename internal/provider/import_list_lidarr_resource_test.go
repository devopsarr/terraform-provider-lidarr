package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListLidarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListLidarrResourceConfig("resourceLidarrTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListLidarrResourceConfig("resourceLidarrTest", "entireArtist"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_lidarr.test", "should_monitor", "entireArtist"),
					resource.TestCheckResourceAttrSet("lidarr_import_list_lidarr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListLidarrResourceConfig("resourceLidarrTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListLidarrResourceConfig("resourceLidarrTest", "specificAlbum"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_lidarr.test", "should_monitor", "specificAlbum"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lidarr_import_list_lidarr.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListLidarrResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "lidarr_import_list_lidarr" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		should_search = false
		root_folder_path = "/config"
		monitor_new_items = "all"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		base_url = "http://127.0.0.1:8686"
		api_key = "testAPIKey"
	}`, folder, name)
}

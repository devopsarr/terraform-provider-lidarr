package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListLastFMUserResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListLastFMUserResourceConfig("resourceLastFMUserTest", "entireArtist"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_lastfm_user.test", "should_monitor", "entireArtist"),
					resource.TestCheckResourceAttrSet("lidarr_import_list_lastfm_user.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListLastFMUserResourceConfig("resourceLastFMUserTest", "specificAlbum"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_lastfm_user.test", "should_monitor", "specificAlbum"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_import_list_lastfm_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListLastFMUserResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "lidarr_import_list_lastfm_user" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		should_search = false
		root_folder_path = "/config"
		monitor_new_items = "all"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		count_list = 25
		user_id = "testUser"
	}`, folder, name)
}

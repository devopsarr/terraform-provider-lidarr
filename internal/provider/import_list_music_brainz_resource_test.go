package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListMusicBrainzResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListMusicBrainzResourceConfig("resourceMusicBrainzTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListMusicBrainzResourceConfig("resourceMusicBrainzTest", "entireArtist"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_music_brainz.test", "should_monitor", "entireArtist"),
					resource.TestCheckResourceAttrSet("lidarr_import_list_music_brainz.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListMusicBrainzResourceConfig("resourceMusicBrainzTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListMusicBrainzResourceConfig("resourceMusicBrainzTest", "specificAlbum"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_music_brainz.test", "should_monitor", "specificAlbum"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_import_list_music_brainz.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListMusicBrainzResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "lidarr_import_list_music_brainz" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		should_search = false
		root_folder_path = "/config"
		monitor_new_items = "all"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		series_id = "testMusicBrainzSeriesID"
	}`, folder, name)
}

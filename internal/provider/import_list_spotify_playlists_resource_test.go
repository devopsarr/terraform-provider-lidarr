package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListSpotifyPlaylistsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListSpotifyPlaylistsResourceConfig("resourceSpotifyPlaylistTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListSpotifyPlaylistsResourceConfig("resourceSpotifyPlaylistTest", "entireArtist"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_spotify_playlists.test", "should_monitor", "entireArtist"),
					resource.TestCheckResourceAttrSet("lidarr_import_list_spotify_playlists.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListSpotifyPlaylistsResourceConfig("resourceSpotifyPlaylistTest", "entireArtist") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListSpotifyPlaylistsResourceConfig("resourceSpotifyPlaylistTest", "specificAlbum"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_spotify_playlists.test", "should_monitor", "specificAlbum"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_import_list_spotify_playlists.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListSpotifyPlaylistsResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "lidarr_import_list_spotify_playlists" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		should_search = false
		root_folder_path = "/config"
		monitor_new_items = "all"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		access_token = "accessToken"
		refresh_token = "refreshToken"
		expires = "0001-01-01T00:01:00Z"
		playlist_ids = ["play1", "play2"]
	}`, folder, name)
}

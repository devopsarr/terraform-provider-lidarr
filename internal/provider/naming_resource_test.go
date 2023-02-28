package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNamingResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNamingResourceConfig("{Artist Name}") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNamingResourceConfig("{Artist Name}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_naming.test", "artist_folder_format", "{Artist Name}"),
					resource.TestCheckResourceAttrSet("lidarr_naming.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNamingResourceConfig("{Artist Name}") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNamingResourceConfig("{Artist_Name}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_naming.test", "artist_folder_format", "{Artist_Name}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_naming.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNamingResourceConfig(artist string) string {
	return fmt.Sprintf(`
	resource "lidarr_naming" "test" {
		rename_tracks              = true
		replace_illegal_characters = true
		standard_track_format      = "{Album Title} ({Release Year})/{Artist Name} - {Album Title} - {track:00} - {Track Title}"
		multi_disc_track_format    = "{Album Title} ({Release Year})/{Medium Format} {medium:00}/{Artist Name} - {Album Title} - {track:00} - {Track Title}"
		artist_folder_format       = "%s"
	}`, artist)
}

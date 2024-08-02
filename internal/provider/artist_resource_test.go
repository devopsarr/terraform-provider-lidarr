package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccArtistResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccArtistResourceConfig("Error", "test", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccArtistResourceConfig("Queen", "test", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_artist.test", "path", "/config/test"),
					resource.TestCheckResourceAttrSet("lidarr_artist.test", "id"),
					resource.TestCheckResourceAttr("lidarr_artist.test", "artist_name", "Queen"),
					resource.TestCheckResourceAttr("lidarr_artist.test", "status", "continuing"),
					resource.TestCheckResourceAttr("lidarr_artist.test", "monitored", "false"),
					resource.TestCheckResourceAttrSet("lidarr_artist.test", "genres.0"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccArtistResourceConfig("Error", "test", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccArtistResourceConfig("Queen", "test123", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_artist.test", "path", "/config/test123"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_artist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccArtistResourceConfig(title, path, foreignID string) string {
	return fmt.Sprintf(`
		resource "lidarr_artist" "test" {
			monitored = false
			artist_name = "%s"
			path = "/config/%s"
			quality_profile_id = 1
			metadata_profile_id = 1
			foreign_artist_id = "%s"
		}
	`, title, path, foreignID)
}

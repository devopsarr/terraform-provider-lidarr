package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataKodiResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_kodi.test", "artist_metadata", "false"),
					resource.TestCheckResourceAttrSet("lidarr_metadata_kodi.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataKodiResourceConfig("kodiResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_kodi.test", "artist_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata_kodi.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataKodiResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "lidarr_metadata_kodi" "test" {
		enable = false
		name = "%s"
		artist_metadata = %s
		album_images = true
		artist_images = true
		album_metadata = false
	}`, name, metadata)
}

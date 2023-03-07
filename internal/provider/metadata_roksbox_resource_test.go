package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataRoksboxResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_roksbox.test", "track_metadata", "false"),
					resource.TestCheckResourceAttrSet("lidarr_metadata_roksbox.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_roksbox.test", "track_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata_roksbox.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataRoksboxResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "lidarr_metadata_roksbox" "test" {
		enable = false
		name = "%s"
		track_metadata = %s
		artist_images = false
		album_images = true
	}`, name, metadata)
}

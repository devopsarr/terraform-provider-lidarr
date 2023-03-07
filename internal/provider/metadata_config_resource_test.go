package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataConfigResourceConfig("allFiles") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataConfigResourceConfig("no"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_config.test", "write_audio_tags", "no"),
					resource.TestCheckResourceAttrSet("lidarr_metadata_config.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataConfigResourceConfig("allFiles") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataConfigResourceConfig("allFiles"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_config.test", "write_audio_tags", "allFiles"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataConfigResourceConfig(country string) string {
	return fmt.Sprintf(`
	resource "lidarr_metadata_config" "test" {
		write_audio_tags = "%s"
		scrub_audio_tags = false
	}`, country)
}

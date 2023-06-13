package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataProfileResourceConfig("error", "1,2") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataProfileResourceConfig("remotemapResourceTest", "1,2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_profile.test", "primary_album_types.1", "2"),
					resource.TestCheckResourceAttrSet("lidarr_metadata_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataProfileResourceConfig("error", "1,2") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataProfileResourceConfig("profileResourceTest", "1,3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_profile.test", "primary_album_types.1", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataProfileResourceConfig(name, primary string) string {
	return fmt.Sprintf(`
		resource "lidarr_metadata_profile" "test" {
  			name = "%s"
			primary_album_types = [%s]
			secondary_album_types = [1]
			release_statuses = [3]
		}
	`, name, primary)
}

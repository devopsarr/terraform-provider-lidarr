package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataResourceConfig("resourceTest", "true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata.test", "track_metadata", "true"),
					resource.TestCheckResourceAttrSet("lidarr_metadata.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataResourceConfig("resourceTest", "true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata.test", "track_metadata", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "lidarr_metadata" "test" {
		enable = true
		name = "%s"
		implementation = "WdtvMetadata"
    	config_contract = "WdtvMetadataSettings"
		track_metadata = %s
	}`, name, metadata)
}

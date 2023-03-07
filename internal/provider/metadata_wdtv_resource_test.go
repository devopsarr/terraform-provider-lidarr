package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataWdtvResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataWdtvResourceConfig("wdtvResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataWdtvResourceConfig("wdtvResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_wdtv.test", "track_metadata", "false"),
					resource.TestCheckResourceAttrSet("lidarr_metadata_wdtv.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataWdtvResourceConfig("wdtvResourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataWdtvResourceConfig("wdtvResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_metadata_wdtv.test", "track_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_metadata_wdtv.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataWdtvResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "lidarr_metadata_wdtv" "test" {
		enable = false
		name = "%s"
		track_metadata = %s
	}`, name, metadata)
}

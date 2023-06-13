package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReleaseProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccReleaseProfileResourceConfig("\"test1\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("\"test1\""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_release_profile.test", "required.0", "test1"),
					resource.TestCheckResourceAttrSet("lidarr_release_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccReleaseProfileResourceConfig("\"test1\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("\"test2\",\"test3\""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_release_profile.test", "required.1", "test3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_release_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccReleaseProfileResourceConfig(required string) string {
	return fmt.Sprintf(`
	resource "lidarr_release_profile" "test" {
		enabled = true
		indexer_id = 0
		required = [%s]
	}`, required)
}

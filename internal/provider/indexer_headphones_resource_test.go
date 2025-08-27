package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerHeadphonesResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerHeadphonesResourceConfig("headphonesResourceTest", "User1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerHeadphonesResourceConfig("headphonesResourceTest", "User1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_headphones.test", "username", "User1"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_headphones.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerHeadphonesResourceConfig("headphonesResourceTest", "User1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerHeadphonesResourceConfig("headphonesResourceTest", "User2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_headphones.test", "username", "User2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer_headphones.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerHeadphonesResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_headphones" "test" {
		enable_automatic_search = false
		priority = 1
		name = "%s"
		username = "%s"
		password = "Pass"
		categories = [ 3000, 3010, 3020, 3030, 3040 ]
	}`, name, user)
}

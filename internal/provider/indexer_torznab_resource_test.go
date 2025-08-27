package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerTorznabResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerTorznabResourceConfig("torznabResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_torznab.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_torznab.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerTorznabResourceConfig("torznabResourceTest", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_torznab.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer_torznab.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerTorznabResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_torznab" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://feed.animetosho.org"
		api_path = "/nabapi"
		minimum_seeders = %d
		categories = [2000,2010]
		priority = 1
	}`, name, seeders)
}

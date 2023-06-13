package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerIptorrentsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_iptorrents.test", "base_url", "https://iptorrents.org"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_iptorrents.test", "id"),
				),
			},
			// Unauthorized Create
			{
				Config:      testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.net"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_iptorrents.test", "base_url", "https://iptorrents.net"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer_iptorrents.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerIptorrentsResourceConfig(name, url string) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_iptorrents" "test" {
		enable_rss = false
		name = "%s"
		base_url = "%s"
		minimum_seeders = 1
	}`, name, url)
}

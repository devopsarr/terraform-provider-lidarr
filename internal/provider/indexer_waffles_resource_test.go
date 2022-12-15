package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerWafflesResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerWafflesResourceConfig("wafflesResourceTest", "User1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_waffles.test", "user_id", "User1"),
					resource.TestCheckResourceAttr("lidarr_indexer_waffles.test", "base_url", "https://www.waffles.ch"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_waffles.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerWafflesResourceConfig("wafflesResourceTest", "User2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_waffles.test", "user_id", "User2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer_waffles.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerWafflesResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_waffles" "test" {
		enable_automatic_search = false
		name = "%s"
		user_id = "%s"
		rss_passkey = "Pass"
		base_url = "https://www.waffles.ch"
		minimum_seeders = 1
	}`, name, user)
}

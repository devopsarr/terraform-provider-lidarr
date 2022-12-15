package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerGazelleResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerGazelleResourceConfig("gazelleResourceTest", "User1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_gazelle.test", "username", "User1"),
					resource.TestCheckResourceAttr("lidarr_indexer_gazelle.test", "base_url", "https://orpheus.network"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_gazelle.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerGazelleResourceConfig("gazelleResourceTest", "User2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_gazelle.test", "username", "User2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer_gazelle.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerGazelleResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_gazelle" "test" {
		enable_automatic_search = false
		name = "%s"
		username = "%s"
		password = "Pass"
		base_url = "https://orpheus.network"
		use_freeleech_token = false
		minimum_seeders = 1
	}`, name, user)
}

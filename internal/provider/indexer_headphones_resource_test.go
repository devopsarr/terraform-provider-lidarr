package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerHeadphonesResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerHeadphonesResourceConfig("headphonesResourceTest", "User1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_headphones.test", "username", "User1"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_headphones.test", "id"),
				),
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
		name = "%s"
		username = "%s"
		password = "Pass"
		categories = [ 3000, 3010, 3020, 3030, 3040 ]
	}`, name, user)
}

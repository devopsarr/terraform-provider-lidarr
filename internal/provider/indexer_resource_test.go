package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerResourceConfig("resourceTest", 25) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", 25),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer.test", "priority", "25"),
					resource.TestCheckResourceAttr("lidarr_indexer.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("lidarr_indexer.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerResourceConfig("resourceTest", 25) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer.test", "priority", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer" "test" {
		priority = %d
		name = "%s"
		implementation = "Newznab"
		protocol = "usenet"
    	config_contract = "NewznabSettings"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [8000, 5000]
	}`, priority, name)
}

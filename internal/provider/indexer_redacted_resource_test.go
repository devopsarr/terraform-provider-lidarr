package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerRedactedResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerRedactedResourceConfig("redactedResourceTest", "Key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerRedactedResourceConfig("redactedResourceTest", "Key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_redacted.test", "api_key", "Key1"),
					resource.TestCheckResourceAttrSet("lidarr_indexer_redacted.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerRedactedResourceConfig("redactedResourceTest", "Key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerRedactedResourceConfig("redactedResourceTest", "Key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_indexer_redacted.test", "api_key", "Key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lidarr_indexer_redacted.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerRedactedResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "lidarr_indexer_redacted" "test" {
		enable_automatic_search = false
		name = "%s"
		api_key = "%s"
		use_freeleech_token = false
		minimum_seeders = 1
		priority = 1
	}`, name, user)
}

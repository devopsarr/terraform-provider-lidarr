package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDelayProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay_profile_resource") + testAccDelayProfileResourceConfig("usenet", "lidarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_delay_profile.test", "preferred_protocol", "usenet"),
					resource.TestCheckResourceAttrSet("lidarr_delay_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay_profile_resource") + testAccDelayProfileResourceConfig("torrent", "lidarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_delay_profile.test", "preferred_protocol", "torrent"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_delay_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDelayProfileResourceConfig(protocol, tag string) string {
	return fmt.Sprintf(`
	resource "lidarr_delay_profile" "test" {
		enable_usenet = true
		enable_torrent = true
		order = 100
		usenet_delay = 0
		torrent_delay = 0
		preferred_protocol= "%s"
		tags = [%s]
	}`, protocol, tag)
}

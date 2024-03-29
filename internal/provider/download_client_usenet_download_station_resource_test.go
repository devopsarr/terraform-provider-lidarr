package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientUsenetDownloadStationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_usenet_download_station.test", "use_ssl", "false"),
					resource.TestCheckResourceAttrSet("lidarr_download_client_usenet_download_station.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientUsenetDownloadStationResourceConfig("resourceUsenetDownloadStationTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_usenet_download_station.test", "use_ssl", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_download_client_usenet_download_station.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientUsenetDownloadStationResourceConfig(name, ssl string) string {
	return fmt.Sprintf(`
	resource "lidarr_download_client_usenet_download_station" "test" {
		enable = false
		use_ssl = %s
		priority = 1
		name = "%s"
		host = "usenet-download-station"
		port = 9091
	}`, ssl, name)
}

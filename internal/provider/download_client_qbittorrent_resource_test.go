package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientQbittorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_qbittorrent.test", "host", "qbittorrent"),
					resource.TestCheckResourceAttr("lidarr_download_client_qbittorrent.test", "url_base", "/qbittorrent/"),
					resource.TestCheckResourceAttrSet("lidarr_download_client_qbittorrent.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientQbittorrentResourceConfig("resourceQbittorrentTest", "qbittorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_qbittorrent.test", "host", "qbittorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_download_client_qbittorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientQbittorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "lidarr_download_client_qbittorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/qbittorrent/"
		port = 9091
		music_category = "tv-lidarr"
		first_and_last = true
	}`, name, host)
}

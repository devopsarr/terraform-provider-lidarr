package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientVuzeResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientVuzeResourceConfig("resourceVuzeTest", "vuze"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_vuze.test", "host", "vuze"),
					resource.TestCheckResourceAttr("lidarr_download_client_vuze.test", "url_base", "/vuze/"),
					resource.TestCheckResourceAttrSet("lidarr_download_client_vuze.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientVuzeResourceConfig("resourceVuzeTest", "vuze-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_download_client_vuze.test", "host", "vuze-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_download_client_vuze.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientVuzeResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "lidarr_download_client_vuze" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/vuze/"
		port = 9091
	}`, name, host)
}

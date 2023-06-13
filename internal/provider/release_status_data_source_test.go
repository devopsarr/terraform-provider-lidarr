package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReleaseStatusDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccReleaseStatusDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccReleaseStatusDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find release_status"),
			},
			// Read testing
			{
				Config: testAccReleaseStatusDataSourceConfig("Official"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_release_status.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_release_status.test", "name", "Official"),
				),
			},
		},
	})
}

func testAccReleaseStatusDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_release_status" "test" {
		name = "%s"
	}
	`, name)
}

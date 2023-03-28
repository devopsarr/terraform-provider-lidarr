package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccReleaseProfileDataSourceConfig("\"999\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccReleaseProfileResourceConfig("\"notreally\"") + testAccReleaseProfileDataSourceConfig("lidarr_release_profile.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_release_profile.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_release_profile.test", "required.0", "notreally")),
			},
		},
	})
}

func testAccReleaseProfileDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "lidarr_release_profile" "test" {
	id = %s
}
`, id)
}

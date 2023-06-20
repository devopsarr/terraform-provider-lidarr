package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityProfileDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccQualityProfileDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find quality_profile"),
			},
			// Read testing
			{
				Config: testAccQualityProfileDataSourceConfig("Lossless"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_quality_profile.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_quality_profile.test", "cutoff", "1005")),
			},
		},
	})
}

func testAccQualityProfileDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_quality_profile" "test" {
		name = "%s"
	}
	`, name)
}

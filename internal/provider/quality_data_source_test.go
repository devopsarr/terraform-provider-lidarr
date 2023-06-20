package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccQualityDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find quality"),
			},
			// Read testing
			{
				Config: testAccQualityDataSourceConfig("MP3-64"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_quality.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_quality.test", "id", "25")),
			},
		},
	})
}

func testAccQualityDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_quality" "test" {
		name = "%s"
	}
	`, name)
}

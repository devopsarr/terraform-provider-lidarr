package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMetadataProfileDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccMetadataProfileDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find metadata_profile"),
			},
			// Read testing
			{
				Config: testAccMetadataProfileDataSourceConfig("Standard"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_metadata_profile.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_metadata_profile.test", "release_statuses.0", "0")),
			},
		},
	})
}

func testAccMetadataProfileDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_metadata_profile" "test" {
		name = "%s"
	}
	`, name)
}

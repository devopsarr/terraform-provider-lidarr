package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityProfilesDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccQualityProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lidarr_quality_profiles.test", "quality_profiles.*", map[string]string{"name": "Any"}),
				),
			},
		},
	})
}

const testAccQualityProfilesDataSourceConfig = `
data "lidarr_quality_profiles" "test" {
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionReleaseGroupDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionReleaseGroupDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_custom_format_condition_release_group.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_custom_format_condition_release_group.test", "name", "HDBits"),
					resource.TestCheckResourceAttr("lidarr_custom_format.test", "specifications.0.value", ".*HDBits.*")),
			},
		},
	})
}

const testAccCustomFormatConditionReleaseGroupDataSourceConfig = `
data  "lidarr_custom_format_condition_release_group" "test" {
	name = "HDBits"
	negate = false
	required = false
	value = ".*HDBits.*"
}

resource "lidarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSReleaseGroup"
	
	specifications = [data.lidarr_custom_format_condition_release_group.test]	
}`

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListExclusionsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccImportListExclusionsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccImportListExclusionResourceConfig("testList", "53b106e7-0cc6-42cc-ac95-ed8d30a3a98e"),
			},
			// Read testing
			{
				Config: testAccImportListExclusionsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lidarr_import_list_exclusions.test", "import_list_exclusions.*", map[string]string{"foreign_id": "53b106e7-0cc6-42cc-ac95-ed8d30a3a98e"}),
				),
			},
		},
	})
}

const testAccImportListExclusionsDataSourceConfig = `
data "lidarr_import_list_exclusions" "test" {
}
`

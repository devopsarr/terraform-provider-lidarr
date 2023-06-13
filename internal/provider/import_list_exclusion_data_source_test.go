package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListExclusionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccImportListExclusionDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccImportListExclusionDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find import_list_exclusion"),
			},
			// Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", "810068af-2b3c-3e9c-b2ab-68a3f3e3787d") + testAccImportListExclusionDataSourceConfig("lidarr_import_list_exclusion.test.foreign_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_import_list_exclusion.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_import_list_exclusion.test", "artist_name", "Queen"),
				),
			},
		},
	})
}

func testAccImportListExclusionDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "lidarr_import_list_exclusion" "test" {
		foreign_id = %s
	}
	`, id)
}

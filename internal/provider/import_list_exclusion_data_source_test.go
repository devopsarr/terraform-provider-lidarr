package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImportListExclusionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_import_list_exclusion.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_import_list_exclusion.test", "artist_name", "testDS"),
				),
			},
		},
	})
}

const testAccImportListExclusionDataSourceConfig = `
resource "lidarr_import_list_exclusion" "test" {
	artist_name = "testDS"
	foreign_id = "8b8a38a9-a290-4560-84f6-3d4466e8d791"
}

data "lidarr_import_list_exclusion" "test" {
	foreign_id = lidarr_import_list_exclusion.test.foreign_id
}
`

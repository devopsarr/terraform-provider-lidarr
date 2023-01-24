package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListExclusionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", "b1a9c0e9-d987-4042-ae91-78d6a3267d69"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_exclusion.test", "foreign_id", "b1a9c0e9-d987-4042-ae91-78d6a3267d69"),
					resource.TestCheckResourceAttrSet("lidarr_import_list_exclusion.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_import_list_exclusion.test", "foreign_id", "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_import_list_exclusion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListExclusionResourceConfig(name, ID string) string {
	return fmt.Sprintf(`
		resource "lidarr_import_list_exclusion" "%s" {
			artist_name = "Queen"
			foreign_id = "%s"
		}
	`, name, ID)
}

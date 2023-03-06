package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecondaryAlbumTypeDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccSecondaryAlbumTypeDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccSecondaryAlbumTypeDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find secondary_album_type"),
			},
			// Read testing
			{
				Config: testAccSecondaryAlbumTypeDataSourceConfig("Demo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_secondary_album_type.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_secondary_album_type.test", "name", "Demo"),
				),
			},
		},
	})
}

func testAccSecondaryAlbumTypeDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_secondary_album_type" "test" {
		name = "%s"
	}
	`, name)
}

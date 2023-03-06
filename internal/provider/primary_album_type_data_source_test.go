package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPrimaryAlbumTypeDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccPrimaryAlbumTypeDataSourceConfig("Error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccPrimaryAlbumTypeDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find primary_album_type"),
			},
			// Read testing
			{
				Config: testAccPrimaryAlbumTypeDataSourceConfig("Album"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_primary_album_type.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_primary_album_type.test", "name", "Album"),
				),
			},
		},
	})
}

func testAccPrimaryAlbumTypeDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_primary_album_type" "test" {
		name = "%s"
	}
	`, name)
}

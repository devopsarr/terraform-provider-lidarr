package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccArtistDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccArtistDataSourceConfig("\"999\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccArtistDataSourceConfig("\"999\""),
				ExpectError: regexp.MustCompile("Unable to find artist"),
			},
			// Read testing
			{
				Config: testAccArtistResourceConfig("Ludwig van Beethoven", "Ludwig_Van_Beethoven", "1f9df192-a621-4f54-8850-2c5373b7eac9") + testAccArtistDataSourceConfig("lidarr_artist.test.foreign_artist_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_artist.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_artist.test", "artist_name", "Ludwig van Beethoven"),
				),
			},
		},
	})
}

func testAccArtistDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "lidarr_artist" "test" {
		foreign_artist_id = %s
	}
	`, id)
}

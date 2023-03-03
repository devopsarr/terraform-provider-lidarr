package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArtistsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccArtistResourceConfig("Error", "error", "error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccArtistResourceConfig("Lucio Battisti", "Lucio_Battisti", "c0c0de23-d9c1-4776-97e0-0c2529402622") + testAccArtistsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lidarr_artists.test", "artists.*", map[string]string{"artist_name": "Lucio Battisti"}),
				),
			},
		},
	})
}

const testAccArtistsDataSourceConfig = `
data "lidarr_artists" "test" {
	depends_on = [lidarr_artist.test]
}
`

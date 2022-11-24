package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIndexerDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_indexer.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_indexer.test", "protocol", "usenet")),
			},
		},
	})
}

const testAccIndexerDataSourceConfig = `
resource "lidarr_indexer" "test" {
	enable_automatic_search = false
	name = "indexerdata"
	implementation = "Newznab"
	protocol = "usenet"
	config_contract = "NewznabSettings"
	base_url = "https://lolo.sickbeard.com"
	api_path = "/api"
	categories = [5030, 5040]
}

data "lidarr_indexer" "test" {
	name = lidarr_indexer.test.name
}
`

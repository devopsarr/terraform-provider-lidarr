package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexerDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccIndexerDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find indexer"),
			},
			// Read testing
			{
				Config: testAccIndexerResourceConfig("indexerdata", 20) + testAccIndexerDataSourceConfig("lidarr_indexer.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_indexer.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_indexer.test", "protocol", "usenet")),
			},
		},
	})
}

func testAccIndexerDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_indexer" "test" {
		name = %s
	}
	`, name)
}

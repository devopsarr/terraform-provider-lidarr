package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a tag to have a value to check
			{
				Config: testAccTagResourceConfig("test-1", "spotify") + testAccTagResourceConfig("test-2", "deezer"),
			},
			// Read testing
			{
				Config: testAccTagsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.lidarr_tags.test", "tags.*", map[string]string{"label": "spotify"}),
				),
			},
		},
	})
}

const testAccTagsDataSourceConfig = `
data "lidarr_tags" "test" {
}
`

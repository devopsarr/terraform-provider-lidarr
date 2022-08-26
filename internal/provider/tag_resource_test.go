package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "mp3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_tag.test", "label", "mp3"),
					resource.TestCheckResourceAttrSet("lidarr_tag.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "hvec"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_tag.test", "label", "hvec"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTagResourceConfig(name, label string) string {
	return fmt.Sprintf(`
		resource "lidarr_tag" "%s" {
  			label = "%s"
		}
	`, name, label)
}

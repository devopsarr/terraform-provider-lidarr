package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRootFolderResourceConfig("all", "/config/asp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_root_folder.test", "monitor_option", "all"),
					resource.TestCheckResourceAttrSet("lidarr_root_folder.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRootFolderResourceConfig("future", "/config/asp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_root_folder.test", "monitor_option", "future"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_root_folder.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRootFolderResourceConfig(monitor, path string) string {
	return fmt.Sprintf(`
		resource "lidarr_root_folder" "test" {
			name = "test"
			quality_profile_id = 1
			monitor_option = "%s"
			new_item_monitor_option = "all"
  			path = "%s"
		}
	`, monitor, path)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationEmbyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationEmbyResourceConfig("resourceEmbyTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_emby.test", "api_key", "token123"),
					resource.TestCheckResourceAttrSet("lidarr_notification_emby.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationEmbyResourceConfig("resourceEmbyTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_emby.test", "api_key", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_emby.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationEmbyResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_emby" "test" {
		on_grab               = false
		on_upgrade            = false
		on_rename             = false
		on_track_retag        = false
		on_release_import     = false
		on_health_issue       = false
		on_application_update = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		host = "emby.lcl"
		port = 8096
		api_key = "%s"
	}`, name, token)
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationBoxcarResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationBoxcarResourceConfig("resourceBoxcarTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_boxcar.test", "token", "token123"),
					resource.TestCheckResourceAttrSet("lidarr_notification_boxcar.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationBoxcarResourceConfig("resourceBoxcarTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_boxcar.test", "token", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_boxcar.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationBoxcarResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_boxcar" "test" {
		on_grab                 = false
		on_import_failure       = false
		on_upgrade              = false
		on_download_failure     = false
		on_release_import   	= false
		on_health_issue         = false
		on_application_update   = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		token = "%s"
	}`, name, token)
}
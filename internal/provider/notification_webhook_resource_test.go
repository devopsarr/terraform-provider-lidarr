package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationWebhookResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_webhook.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("lidarr_notification_webhook.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_webhook.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationWebhookResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_webhook" "test" {
		on_grab                            = false
		on_track_retag                     = true
		on_upgrade                         = %s
		on_rename                          = false
		on_release_import                  = false
		on_download_failure                = false
		on_import_failure 				   = true
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		url = "http://transmission:9091"
		method = 1
	}`, upgrade, name)
}

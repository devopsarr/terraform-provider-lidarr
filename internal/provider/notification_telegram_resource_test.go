package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationTelegramResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_telegram.test", "chat_id", "chat01"),
					resource.TestCheckResourceAttrSet("lidarr_notification_telegram.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTelegramResourceConfig("resourceTelegramTest", "chat02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_telegram.test", "chat_id", "chat02"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_telegram.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTelegramResourceConfig(name, chat string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_telegram" "test" {
		on_grab               	= false
		on_import_failure		= false
		on_upgrade           	= false
		on_download_failure    	= false
		on_release_import   	= false
		on_health_issue       	= false
		on_application_update  	= false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		chat_id = "%s"
		bot_token = "Token"
	}`, name, chat)
}

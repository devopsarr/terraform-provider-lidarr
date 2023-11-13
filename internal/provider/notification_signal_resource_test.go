package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationSignalResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSignalResourceConfig("resourceSignalTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "chat01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_signal.test", "auth_password", "chat01"),
					resource.TestCheckResourceAttrSet("lidarr_notification_signal.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSignalResourceConfig("resourceSignalTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "chat02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_signal.test", "auth_password", "chat02"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lidarr_notification_signal.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password", "sender_number"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSignalResourceConfig(name, chat string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_signal" "test" {
		on_grab               	= false
		on_import_failure		= false
		on_upgrade           	= false
		on_download_failure    	= false
		on_release_import   	= false
		on_health_issue       	= false
		on_application_update  	= false

		include_health_warnings = false
		name                    = "%s"

		auth_username = "User"
		auth_password = "%s"

		host = "localhost"
		port = 8080
		use_ssl = true
		sender_number = "1234"
		receiver_id = "4321"
	}`, name, chat)
}

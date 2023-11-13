package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationTwitterResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationTwitterResourceConfig("resourceTwitterTest", "me") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationTwitterResourceConfig("resourceTwitterTest", "me"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_twitter.test", "mention", "me"),
					resource.TestCheckResourceAttrSet("lidarr_notification_twitter.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationTwitterResourceConfig("resourceTwitterTest", "me") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTwitterResourceConfig("resourceTwitterTest", "myself"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_twitter.test", "mention", "myself"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "lidarr_notification_twitter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_token", "access_token_secret", "consumer_key", "consumer_secret"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTwitterResourceConfig(name, mention string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_twitter" "test" {
		on_grab                	= false
		on_import_failure      	= false
		on_upgrade            	= false
		on_download_failure   	= false
		on_release_import 		= false
		on_health_issue        	= false
		on_application_update	= false

		include_health_warnings = false
		name                    = "%s"

		access_token = "Token"
		access_token_secret = "TokenSecret"
		consumer_key = "Key"
		consumer_secret = "Secret"
		mention = "%s"
	}`, name, mention)
}

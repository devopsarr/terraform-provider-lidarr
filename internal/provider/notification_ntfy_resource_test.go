package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationNtfyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_ntfy.test", "password", "key1"),
					resource.TestCheckResourceAttrSet("lidarr_notification_ntfy.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_ntfy.test", "password", "key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_ntfy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationNtfyResourceConfig(name, key string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_ntfy" "test" {
		on_grab                 = false
		on_import_failure       = false
		on_upgrade              = false
		on_download_failure     = false
		on_release_import   	= false
		on_health_issue   	    = false
		on_application_update   = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		priority = 1
		server_url = "https://ntfy.sh"
		username = "User"
		password = "%s"
		topics = ["Topic1234","Topic4321"]
		field_tags = ["warning","skull"]
	}`, name, key)
}

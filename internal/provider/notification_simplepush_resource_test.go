package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSimplepushResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "chat01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_simplepush.test", "key", "chat01"),
					resource.TestCheckResourceAttrSet("lidarr_notification_simplepush.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "chat01") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "chat02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_simplepush.test", "key", "chat02"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_simplepush.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSimplepushResourceConfig(name, chat string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_simplepush" "test" {
		on_grab               	= false
		on_import_failure		= false
		on_upgrade           	= false
		on_download_failure    	= false
		on_release_import   	= false
		on_health_issue       	= false
		on_application_update  	= false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		key = "%s"
		event = "Test"
	}`, name, chat)
}

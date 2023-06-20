package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationEmailResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationEmailResourceConfig("resourceEmailTest", "test@email.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationEmailResourceConfig("resourceEmailTest", "test@email.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_email.test", "from", "test@email.com"),
					resource.TestCheckResourceAttrSet("lidarr_notification_email.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationEmailResourceConfig("resourceEmailTest", "test@email.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationEmailResourceConfig("resourceEmailTest", "test123@email.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_email.test", "from", "test123@email.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_email.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationEmailResourceConfig(name, from string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_email" "test" {
		on_grab                            = false
		on_import_failure                        = false
		on_upgrade                         = false
		on_download_failure                    = false
		on_release_import   = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		server = "http://email-server.net"
		port = 587
		from = "%s"
		to = ["test@test.com", "test1@test.com"]
	}`, name, from)
}

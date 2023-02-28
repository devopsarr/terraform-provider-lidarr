package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationMailgunResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_mailgun.test", "from", "test@mailgun.com"),
					resource.TestCheckResourceAttrSet("lidarr_notification_mailgun.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test@mailgun.com") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationMailgunResourceConfig("resourceMailgunTest", "test123@mailgun.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_mailgun.test", "from", "test123@mailgun.com"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_mailgun.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationMailgunResourceConfig(name, from string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_mailgun" "test" {
		on_grab           		= false
		on_upgrade        	 	= false
		on_release_import   	= false
		on_health_issue    		= false
		on_application_update   = false
	  
		include_health_warnings = false
		name                    = "%s"
		
		api_key = "APIkey"
		from = "%s"
		recipients = ["test@test.com", "test1@test.com"]
	}`, name, from)
}

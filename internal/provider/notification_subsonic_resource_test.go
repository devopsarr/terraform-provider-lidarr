package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSubsonicResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationSubsonicResourceConfig("resourceSubsonicTest", "pass1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_subsonic.test", "password", "pass1"),
					resource.TestCheckResourceAttrSet("lidarr_notification_subsonic.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSubsonicResourceConfig("resourceSubsonicTest", "pass2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_notification_subsonic.test", "password", "pass2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_notification_subsonic.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSubsonicResourceConfig(name, avatar string) string {
	return fmt.Sprintf(`
	resource "lidarr_notification_subsonic" "test" {
		on_grab                            = false
		on_upgrade                         = false
		on_rename                          = false
		on_track_retag               	   = false
		on_release_import   	 		   = false
		on_health_issue                    = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		host = "http://subsonic.com"
		port = 8080
		username = "User"
		password = "%s"
		notify = true
	}`, name, avatar)
}

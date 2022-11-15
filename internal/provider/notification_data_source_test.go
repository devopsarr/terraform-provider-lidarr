package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotificationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

const testAccNotificationDataSourceConfig = `
resource "lidarr_notification" "test" {
	on_grab                            = false
	on_track_retag                     = true
	on_upgrade                         = true
	on_rename                          = false
	on_release_import                  = false
	on_download_failure                = false
	on_import_failure 				   = true
	on_health_issue                    = false
	on_application_update              = false
  
	include_health_warnings = false
	name                    = "notificationData"
  
	implementation  = "CustomScript"
	config_contract = "CustomScriptSettings"
  
	path = "/scripts/test.sh"
}

data "lidarr_notification" "test" {
	name = lidarr_notification.test.name
}
`

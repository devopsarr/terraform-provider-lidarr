package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMediaManagementResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMediaManagementResourceConfig("none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMediaManagementResourceConfig("none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_media_management.test", "file_date", "none"),
					resource.TestCheckResourceAttrSet("lidarr_media_management.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMediaManagementResourceConfig("none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMediaManagementResourceConfig("albumReleaseDate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_media_management.test", "file_date", "albumReleaseDate"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_media_management.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMediaManagementResourceConfig(date string) string {
	return fmt.Sprintf(`
	resource "lidarr_media_management" "test" {
		unmonitor_previous_tracks   = true
		hardlinks_copy              = true
		create_empty_folders        = true
		delete_empty_folders        = true
		watch_library_for_changes   = true
		import_extra_files          = true
		set_permissions             = true
		skip_free_space_check       = true
		minimum_free_space          = 100
		recycle_bin_days            = 7
		chmod_folder                = "755"
		chown_group                 = "arrs"
		download_propers_repacks    = "preferAndUpgrade"
		allow_fingerprinting        = "never"
		extra_file_extensions       = "info"
		file_date                   = "%s"
		recycle_bin_path            = ""
		rescan_after_refresh        = "always"
	}`, date)
}

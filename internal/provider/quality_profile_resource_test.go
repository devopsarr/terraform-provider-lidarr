package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-flac"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_quality_profile.test", "name", "example-flac"),
					resource.TestCheckResourceAttrSet("lidarr_quality_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-alac"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_quality_profile.test", "name", "example-alac"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_quality_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccQualityProfileResourceError = `
resource "lidarr_quality_profile" "test" {
	name            = "Error"
	upgrade_allowed = true
	cutoff          = 2000
	quality_groups = []
}
`

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	data "lidarr_quality" "flac" {
		name = "FLAC"
	}

	data "lidarr_quality" "alac" {
		name = "ALAC"
	}

	data "lidarr_quality" "ogg" {
		name = "OGG Vorbis Q10"
	}

	resource "lidarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 2000

		quality_groups = [
			{
				id   = 2000
				name = "lossless"
				qualities = [
					data.lidarr_quality.alac,
					data.lidarr_quality.flac,
				]
			},
			{
				qualities = [data.lidarr_quality.ogg]
			}
		]
	}
	`, name)
}

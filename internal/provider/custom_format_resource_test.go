package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccCustomFormatResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_custom_format.test", "include_custom_format_when_renaming", "false"),
					resource.TestCheckResourceAttrSet("lidarr_custom_format.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccCustomFormatResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lidarr_custom_format.test", "include_custom_format_when_renaming", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "lidarr_custom_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomFormatResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "lidarr_custom_format" "test" {
		include_custom_format_when_renaming = %s
		name = "%s"
		
		specifications = [
			{
				name = "Preferred Words"
				implementation = "ReleaseTitleSpecification"
				negate = false
				required = false
				value = "\\b(SPARKS|Framestor)\\b"
			},
			{
				name = "Size"
				implementation = "SizeSpecification"
				negate = false
				required = false
				min = 0
				max = 100
			}
		]
	}`, enable, name)
}

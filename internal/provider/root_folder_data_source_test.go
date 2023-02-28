package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/devopsarr/lidarr-go/lidarr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRootFolderDataSourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRootFolderDataSourceConfig("/error"),
				ExpectError: regexp.MustCompile("Unable to find root_folder"),
			},
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccRootFolderDataSourceConfig("/config"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.lidarr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.lidarr_root_folder.test", "path", "/config")),
			},
		},
	})
}

func testAccRootFolderDataSourceConfig(path string) string {
	return fmt.Sprintf(`
	data "lidarr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

func rootFolderDSInit() {
	// ensure a /config root path is configured
	client := testAccAPIClient()
	folder := lidarr.NewRootFolderResource()
	folder.SetPath("/config")
	folder.SetName("config")
	folder.SetDefaultQualityProfileId(1)
	folder.SetDefaultMetadataProfileId(1)
	folder.SetDefaultTags([]*int32{})
	_, _, _ = client.RootFolderApi.CreateRootFolder(context.TODO()).RootFolderResource(*folder).Execute()
}

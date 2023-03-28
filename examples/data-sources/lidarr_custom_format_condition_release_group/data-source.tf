data "lidarr_custom_format_condition_release_group" "example" {
  name     = "HDBits"
  negate   = false
  required = false
  value    = ".*HDBits.*"
}

resource "lidarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.lidarr_custom_format_condition_release_group.example]
}
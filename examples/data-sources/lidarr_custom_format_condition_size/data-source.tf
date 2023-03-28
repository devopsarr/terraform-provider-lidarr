data "lidarr_custom_format_condition_size" "example" {
  name     = "Example"
  negate   = false
  required = false
  min      = 5
  max      = 50
}

resource "lidarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.lidarr_custom_format_condition_size.example]
}
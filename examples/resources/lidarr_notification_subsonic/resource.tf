resource "lidarr_notification_subsonic" "example" {
  on_grab           = false
  on_upgrade        = false
  on_rename         = false
  on_track_retag    = false
  on_release_import = true
  on_health_issue   = false

  include_health_warnings = false
  name                    = "Example"

  host     = "http://subsonic.com"
  port     = 8080
  username = "User"
  password = "MyPass"
  notify   = true
}
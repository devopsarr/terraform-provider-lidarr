resource "lidarr_notification_emby" "example" {
  on_grab               = false
  on_upgrade            = true
  on_rename             = false
  on_track_retag        = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  host    = "emby.lcl"
  port    = 8096
  api_key = "API_Key"
}
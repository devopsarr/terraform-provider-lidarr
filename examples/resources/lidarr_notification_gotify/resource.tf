resource "lidarr_notification_gotify" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_rename             = false
  on_download_failure   = false
  on_track_retag        = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  server    = "http://gotify-server.net"
  app_token = "Token"
  priority  = 5
}
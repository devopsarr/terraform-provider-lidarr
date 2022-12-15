resource "lidarr_notification_notifiarr" "example" {
  on_grab               = false
  on_upgrade            = true
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key = "Token"
}